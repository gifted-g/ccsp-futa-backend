package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"ccsp-futa-alumni/db"
	"ccsp-futa-alumni/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileUpdateReq struct {
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
	Phone       string `json:"phone"`
	Location    string `json:"location"`
}

func GetProfileHandler(c *gin.Context) {
	uid := c.Param("user_id")
	var p models.Profile
	if err := db.DB.Where("user_id = ?", uid).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"profile": p})
}

func UpdateProfileHandler(c *gin.Context) {
	uid := c.Param("user_id")
	sub := c.GetString("sub")
	if sub != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	var req ProfileUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	if err := db.DB.Model(&models.Profile{}).Where("user_id = ?", uid).Updates(models.Profile{
		DisplayName: req.DisplayName,
		Bio:         req.Bio,
		Phone:       req.Phone,
		Location:    req.Location,
		UpdatedAt:   time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"}); return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "updated"})
}

// RequestAvatarUploadHandler returns a server-side upload URL (for simplicity, server handles upload)
func RequestAvatarUploadHandler(c *gin.Context) {
	uid := c.Param("user_id")
	sub := c.GetString("sub")
	if sub != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"}); return
	}
	// for S3: return presigned PUT; here we return the server upload endpoint
	uploadURL := fmt.Sprintf("/api/v1/profile/%s/avatar/upload", uid)
	c.JSON(http.StatusOK, gin.H{
		"upload_url": uploadURL,
		"method":     "POST",
		"fields":     nil,
	})
}

// UploadAvatarHandler accepts multipart upload and stores file under ./uploads/<userID>/
func UploadAvatarHandler(c *gin.Context) {
	uid := c.Param("user_id")
	sub := c.GetString("sub")
	if sub != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"}); return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"}); return
	}
	// validate extension
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create uploads dir"}); return
	}
	filename := uuid.New().String() + ext
	dest := filepath.Join(uploadsDir, filename)
	if err := c.SaveUploadedFile(file, dest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed"}); return
	}
	publicURL := fmt.Sprintf("/%s/%s", uploadsDir, filename) // static serve in production should be via CDN
	// Save into profile (temporary until confirm)
	if err := db.DB.Model(&models.Profile{}).Where("user_id = ?", uid).Update("profile_img_url", publicURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db update failed"}); return
	}
	c.JSON(http.StatusOK, gin.H{"url": publicURL})
}

// ConfirmAvatarHandler can be used if you use presigned S3 flow (here it simply re-checks)
type ConfirmAvatarReq struct {
	ProfileImageURL string `json:"profile_image_url" binding:"required"`
}

func ConfirmAvatarHandler(c *gin.Context) {
	uid := c.Param("user_id")
	sub := c.GetString("sub")
	if sub != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"}); return
	}
	var req ConfirmAvatarReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	if err := db.DB.Model(&models.Profile{}).Where("user_id = ?", uid).Update("profile_img_url", req.ProfileImageURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db update failed"}); return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "confirmed"})
}
