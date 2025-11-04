package handlers

import (
	"net/http"
	"time"

	"ccsp-futa-alumni/db"
	"ccsp-futa-alumni/models"
	"ccsp-futa-alumni/notifications"
	"ccsp-futa-alumni/ws"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateChannelReq struct {
	Name    string   `json:"name"`
	IsGroup bool     `json:"is_group"`
	Members []string `json:"members"` // user ids
}

func CreateChannelHandler(c *gin.Context) {
	var req CreateChannelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	creator := c.GetString("sub")
	ch := models.ChatChannel{
		ID:        uuid.New(),
		Name:      req.Name,
		IsGroup:   req.IsGroup,
		CreatedBy: uuid.MustParse(creator),
		CreatedAt: time.Now(),
	}
	if err := db.DB.Create(&ch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"}); return
	}
	// add members
	for _, m := range req.Members {
		member := models.ChatMember{
			ChannelID: ch.ID,
			UserID:    uuid.MustParse(m),
			Role:      "member",
			JoinedAt:  time.Now(),
		}
		db.DB.Create(&member)
	}
	// add creator as member if not present
	db.DB.FirstOrCreate(&models.ChatMember{}, models.ChatMember{ChannelID: ch.ID, UserID: uuid.MustParse(creator), Role: "admin", JoinedAt: time.Now()})
	c.JSON(http.StatusCreated, gin.H{"channel": ch})
}

func GetChannelHandler(c *gin.Context) {
	cid := c.Param("channel_id")
	var ch models.ChatChannel
	if err := db.DB.Preload("Members").Where("id = ?", cid).First(&ch).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "channel not found"}); return
	}
	c.JSON(http.StatusOK, gin.H{"channel": ch})
}

type AddMemberReq struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role"`
}

func AddChannelMemberHandler(c *gin.Context) {
	cid := c.Param("channel_id")
	var req AddMemberReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	member := models.ChatMember{
		ChannelID: uuid.MustParse(cid),
		UserID:    uuid.MustParse(req.UserID),
		Role:      req.Role,
		JoinedAt:  time.Now(),
	}
	if err := db.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add"}); return
	}
	c.JSON(http.StatusCreated, gin.H{"member": member})
}

type PostMessageReq struct {
	Body        string                 `json:"body"`
	Attachments map[string]interface{} `json:"attachments"`
}

func PostMessageHandler(c *gin.Context) {
	cid := c.Param("channel_id")
	var req PostMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	sender := c.GetString("sub")
	msg := models.Message{
		ID:          uuid.New(),
		ChannelID:   uuid.MustParse(cid),
		SenderID:    uuid.MustParse(sender),
		Body:        req.Body,
		Attachments: req.Attachments,
		CreatedAt:   time.Now(),
		Status:      "sent",
	}
	if err := db.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed"}); return
	}

	// Publish realtime via ws hub
	go ws.Hub.BroadcastMessage(msg)

	// Send notifications to offline users (worker)
	go notifications.NotifyChannelMembers(msg)

	c.JSON(http.StatusCreated, gin.H{"message": msg})
}

func ListMessagesHandler(c *gin.Context) {
	cid := c.Param("channel_id")
	limit := 50
	var msgs []models.Message
	if err := db.DB.Where("channel_id = ?", cid).Order("created_at desc").Limit(limit).Find(&msgs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch"}); return
	}
	c.JSON(http.StatusOK, gin.H{"messages": msgs})
}
