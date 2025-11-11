package test

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/gin-gonic/gin"
    "ccsp-futa-alumni/handlers"   // FIXED import prefix
    "ccsp-futa-alumni/middleware"
)

func TestAdminMiddlewareBlocksNonAdmin(t *testing.T) {
    r := gin.Default()
    r.Use(func(c *gin.Context) { c.Set("role", "user") }) // fake user
    r.GET("/protected", middleware.IsAdmin(), func(c *gin.Context) {
        c.JSON(200, gin.H{"ok": true})
    })

    req, _ := http.NewRequest("GET", "/protected", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusForbidden {
        t.Errorf("expected 403, got %d", w.Code)
    }
}

func TestBroadcastEndpoint(t *testing.T) {
    r := gin.Default()
    r.Use(func(c *gin.Context) {
        c.Set("role", "admin")
        c.Set("user_id", "admin-123")
    })
    r.POST("/broadcast", handlers.Broadcast)

    body := `{"content":"Hello world","type":"announcement"}`
    req, _ := http.NewRequest("POST", "/broadcast", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    if w.Code != http.StatusCreated {
        t.Errorf("expected 201, got %d", w.Code)
    }
}
