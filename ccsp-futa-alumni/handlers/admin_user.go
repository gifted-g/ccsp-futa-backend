package handlers

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"

	 "github.com/ccspfutaalumnitech/backend/internal/util""

)

type AdminUserCreate struct { Email string `json:"email"` PasswordHash string `json:"password_hash"` FirstName string `json:"first_name"` LastName string `json:"last_name"` Role string `json:"role"` }

type AdminUserUpdate struct { 