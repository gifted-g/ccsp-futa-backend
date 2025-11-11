package handlers

import (
	"net/http"
	"time"
	"database/sql"

	"github.com/gin-gonic/gin"
)

var DB *sql.DB

func DeleteMessage(c *gin.Context) {
	msgID := c.Param("id")
	adminID := c.GetString("user_id")
	_, err := DB.Exec(`UPDATE messages SET is_hidden=true, deleted_by=$1, deleted_at=$2 WHERE id=$3`,
		adminID, time.Now(), msgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete message"})
		return
	}
	
	addAudit(adminID, "delete", msgID, "Message hidden")
	c.JSON(http.StatusOK, gin.H{"status": "message deleted"})
}

func addAudit(adminID, action, msgID, details string) {
	_, _ = DB.Exec(`INSERT INTO audit_logs (admin_id, action, message_id, details) VALUES ($1,$2,$3,$4)`,
		adminID, action, msgID, details)
}

func PinMessage(c *gin.Context) {
    msgID := c.Param("id")
    adminID := c.GetString("user_id")
    _, err := DB.Exec(`UPDATE messages SET is_pinned=true WHERE id=$1`, msgID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not pin message"})
        return
    }
   

    addAudit(adminID, "pin", msgID, "Message pinned")
    c.JSON(http.StatusOK, gin.H{"status": "message pinned"})
}


func Broadcast(c *gin.Context) {
    var req struct {
        Content string `json:"content"`
        Type    string `json:"type"` 
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }
    
    // FIX: Define adminID
    adminID := c.GetString("user_id")

    _, err := DB.Exec(`INSERT INTO messages (user_id, content, broadcast_type) VALUES ($1,$2,$3)`,
        adminID, req.Content, req.Type)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not broadcast"})
        return
    }
    // REMOVE THE EXTRA '}' HERE

    addAudit(adminID, "broadcast", "", "Broadcast: "+req.Content)
    c.JSON(http.StatusCreated, gin.H{"status": "broadcast sent"})
}


