package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//
// ─── USER AND PROFILE MODELS ────────────────────────────────────────────────────
//

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Phone        string    `json:"phone"`
	Role         string    `json:"role"`
	Active       bool      `gorm:"default:true" json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	
    FullName     string // <-- ADD THIS
  


	Profile Profile `gorm:"constraint:OnDelete:CASCADE" json:"profile"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

type Profile struct {
	UserID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	DisplayName   string    `json:"display_name"`
	Bio           string    `gorm:"type:text" json:"bio"`
	Location      string    `json:"location"`
	Phone        string    `json:"phone"`
	ProfileImgURL string    `json:"profile_image_url"`
	UpdatedAt     time.Time `json:"updated_at"`
}

//
// ─── POSTS AND EVENTS ───────────────────────────────────────────────────────────
//

type Post struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	AuthorID  *uuid.UUID `gorm:"type:uuid;index" json:"author_id"`
	Title     string     `gorm:"not null" json:"title"`
	Body      string     `gorm:"type:text" json:"body"`
	Published bool       `gorm:"default:false" json:"published"`
	CreatedAt time.Time  `json:"created_at"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

type Event struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	StartsAt    string     `json:"starts_at"`
	EndsAt      string     `json:"ends_at"`
	IsPublished bool       `gorm:"default:false" json:"is_published"`
	OrganizerID *uuid.UUID `gorm:"type:uuid;index" json:"organizer_id"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return
}

type RSVP struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	EventID uuid.UUID `gorm:"type:uuid;index" json:"event_id"`
	UserID  uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Status  string    `gorm:"default:'pending'" json:"status"`
}

func (r *RSVP) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return
}

//
// ─── CHAT SYSTEM ────────────────────────────────────────────────────────────────
//

type ChatChannel struct {
	ID        uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string                 `json:"name"`
	IsGroup   bool                   `gorm:"default:false" json:"is_group"`
	CreatedBy uuid.UUID              `gorm:"type:uuid;index" json:"created_by"`
	CreatedAt time.Time              `json:"created_at"`
	Metadata  map[string]interface{} `gorm:"type:jsonb;default:'{}'" json:"metadata"`
	Members   []ChatMember           `gorm:"foreignKey:ChannelID" json:"members"`
	Messages  []Message              `gorm:"foreignKey:ChannelID" json:"messages"`
}

func (c *ChatChannel) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}

type ChatMember struct {
	ChannelID uuid.UUID `gorm:"type:uuid;primaryKey" json:"channel_id"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	Role      string    `gorm:"default:'member'" json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
}

type Message struct {
	ID          uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	ChannelID   uuid.UUID              `gorm:"type:uuid;index" json:"channel_id"`
	SenderID    uuid.UUID              `gorm:"type:uuid;index" json:"sender_id"`
	Body        string                 `gorm:"type:text" json:"body"`
	Attachments map[string]interface{} `gorm:"type:jsonb" json:"attachments,omitempty"`
	Status      string                 `gorm:"default:'sent'" json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}

type AuditLog struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	AdminID   uuid.UUID  `gorm:"type:uuid;index" json:"admin_id"`
	Action    string     `json:"action"`
	MessageID *uuid.UUID `gorm:"type:uuid;index" json:"message_id,omitempty"`
	Details   string     `json:"details"`
	CreatedAt time.Time  `json:"created_at"`
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}

//
// ─── GROUPS AND MEMBERS ─────────────────────────────────────────────────────────
//

type SetGroup struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Year int       `json:"year"`
	Name string    `json:"name"`
}

func (s *SetGroup) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

type SetMember struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SetID  uuid.UUID `gorm:"type:uuid;index" json:"set_id"`
	UserID uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Role   string    `gorm:"default:'member'" json:"role"`
}

func (sm *SetMember) BeforeCreate(tx *gorm.DB) (err error) {
	if sm.ID == uuid.Nil {
		sm.ID = uuid.New()
	}
	return
}

//
// ─── PUSH NOTIFICATIONS ─────────────────────────────────────────────────────────
//

type PushToken struct {
	Token     string    `gorm:"primaryKey" json:"token"`
	UserID    uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Platform  string    `json:"platform"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
}
