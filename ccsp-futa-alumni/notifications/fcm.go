package notifications

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"ccsp-futa-alumni/db"
	"ccsp-futa-alumni/models"
)

const fcmURLTemplate = "https://fcm.googleapis.com/v1/projects/%s/messages:send"

// If you set FCM_PROJECT_ID and FCM_SERVICE_ACCOUNT (not shipped here) you can implement proper FCM HTTP v1 with OAuth2.
// For simplicity this function logs and attempts to send if FCM_SERVER_KEY env var is set (legacy server key).
func sendFCMLegacy(token, title, body string, data map[string]string) error {
	serverKey := "" // os.Getenv("FCM_SERVER_KEY") // legacy
	if serverKey == "" {
		// Not configured -> log and skip
		log.Printf("[FCM] not configured. would-send -> token=%s title=%s body=%s data=%v\n", token, title, body, data)
		return nil
	}
	payload := map[string]interface{}{
		"to": token,
		"notification": map[string]string{
			"title": title,
			"body":  body,
		},
		"data": data,
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewReader(b))
	req.Header.Set("Authorization", "key="+serverKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Printf("[FCM] sent %s -> status=%s", token, resp.Status)
	return nil
}

// NotifyChannelMembers fetches push tokens for channel members (naive) and sends a notification
func NotifyChannelMembers(msg models.Message) {
	// Fetch members of the channel
	var members []models.ChatMember
	if err := db.DB.Where("channel_id = ?", msg.ChannelID).Find(&members).Error; err != nil {
		log.Println("notify: could not fetch members", err)
		return
	}
	// Build list of user ids excluding sender
	var userIDs []string
	for _, m := range members {
		if m.UserID.String() == msg.SenderID.String() {
			continue
		}
		userIDs = append(userIDs, m.UserID.String())
	}
	// Query push tokens
	var tokens []models.PushToken
	if err := db.DB.Where("user_id IN ?", userIDs).Find(&tokens).Error; err != nil {
		log.Println("notify: could not fetch tokens", err)
	}
	for _, t := range tokens {
		_ = sendFCMLegacy(t.Token, "New message", msg.Body, map[string]string{
			"channel_id": msg.ChannelID.String(),
			"message_id": msg.ID.String(),
		})
	}
}
