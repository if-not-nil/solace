package orm

import (
	"crypto/rand"
	"errors"
	"regexp"
	"time"

	"gorm.io/gorm"
)

var valid_name = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type Model struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// id generator
func GenerateID() (string, error) {
	const length = 8
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i := range bytes {
		bytes[i] = chars[int(bytes[i])%len(chars)]
	}
	return string(bytes), nil
}

// hook for auto-generating string id
func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID, err = GenerateID()
	}
	return
}

type File struct {
	Model
	Type string `gorm:"not null"`
}

type InviteLink struct {
	Model     `tstype:",extends"`
	ServerID  string `gorm:"not null" json:"server_id"`
	JoinsLeft uint   `gorm:"not null" json:"joins_left"`
}

type Server struct {
	Model
	Name     string    `gorm:"unique;not null" json:"name"`
	OwnerID  string    `gorm:"not null" json:"owner_id"`
	Owner    User      `gorm:"foreignKey:OwnerID" json:"-"`
	AvatarID string    `json:"avatar"`
	Avatar   File      `gorm:"foreignKey:AvatarID;references:ID" json:"-"`
	Channels []Channel `json:"channels,omitempty"`
}

func (s *Server) BeforeSave(tx *gorm.DB) error {
	if tx.Statement.Changed("Name") && !valid_name.MatchString(s.Name) {
		return errors.New("invalid server name: must be a-z, A-Z, 0-9, underscore or hyphen only")
	}

	return nil
}

type Channel struct {
	Model
	ServerID string `gorm:"not null;index:idx_server_name" json:"server_id"`
	Server   Server `gorm:"foreignKey:ServerID" json:"-"`
	Name     string `gorm:"not null;index:idx_server_name,unique" json:"name"`
	Type     string `gorm:"not null" json:"type"`
}

func (c *Channel) BeforeSave(tx *gorm.DB) error {
	if !valid_name.MatchString(c.Name) {
		return errors.New("invalid channel name: must be a-z, A-Z, 0-9, underscore or hyphen only")
	}
	return nil
}

type Membership struct {
	Model
	UserID   string `gorm:"not null;index:idx_user_server,unique" json:"user_id"`
	User     User   `gorm:"foreignKey:UserID" json:"-"`
	ServerID string `gorm:"not null;index:idx_user_server,unique" json:"server_id"`
	Server   Server `gorm:"foreignKey:ServerID" json:"-"`
	Banned   bool   `gorm:"default:false" json:"banned"`
}

type Role struct {
	Model       `tstype:",extends"`
	ServerID    string   `gorm:"not null;index" json:"server_id"`
	Server      Server   `gorm:"foreignKey:ServerID" json:"-"`
	Name        string   `gorm:"not null" json:"name"`
	Color       string   `json:"color"`
	Permissions []string `gorm:"type:text;serializer:json" json:"permissions"`
}

type UserRole struct {
	Model
	UserID   string `gorm:"not null;index:idx_user_role_server,unique" json:"user_id"`
	User     User   `gorm:"foreignKey:UserID" json:"-"`
	RoleID   string `gorm:"not null;index:idx_user_role_server,unique" json:"role_id"`
	Role     Role   `gorm:"foreignKey:RoleID" json:"-"`
	ServerID string `gorm:"not null;index:idx_user_role_server,unique" json:"server_id"`
	Server   Server `gorm:"foreignKey:ServerID" json:"-"`
}

type ChannelPermissionOverride struct {
	Model
	ServerID  string   `gorm:"not null;index" json:"server_id"`
	Server    Server   `gorm:"foreignKey:ServerID" json:"-"`
	ChannelID string   `gorm:"not null;index" json:"channel_id"`
	Channel   Channel  `gorm:"foreignKey:ChannelID" json:"-"`
	RoleID    string   `gorm:"not null;index" json:"role_id"`
	Role      Role     `gorm:"foreignKey:RoleID" json:"-"`
	AllowMask []string `gorm:"type:text;serializer:json" json:"allow_mask"`
	DenyMask  []string `gorm:"type:text;serializer:json" json:"deny_mask"`
}

type User struct {
	Model
	Name     string `gorm:"unique;not null" json:"name"`
	Password string `gorm:"not null" json:"-"`
	AvatarID string `json:"avatar"`
	Avatar   File   `gorm:"foreignKey:AvatarID;references:ID" json:"-"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `gorm:"unique;not null" json:"name"`
	AvatarID string `json:"avatar"`
	Avatar   File   `gorm:"foreignKey:AvatarID;references:ID" json:"-"`
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	if !valid_name.MatchString(u.Name) {
		return errors.New("invalid user name: must be a-z, A-Z, 0-9, underscore or hyphen only")
	}
	return nil
}

type ChannelResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ServerID string `json:"server_id"`
	Type     string `json:"type"`
}

type Message struct {
	Model
	Content   string  `gorm:"type:text;not null" json:"content"`
	UserID    string  `gorm:"not null" json:"user_id"`
	ChannelID string  `gorm:"not null" json:"channel_id"`
	User      User    `gorm:"foreignKey:UserID" json:"user"`
	Channel   Channel `gorm:"foreignKey:ChannelID" json:"-"`
}

type MessageResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	ChannelID string    `json:"channel_id"`
	CreatedAt time.Time `json:"created_at"`
}

type MembershipResponse struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	ServerID string `json:"server_id"`
	Role     string `json:"role"`
}

type ServerResponse struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	OwnerID  string            `json:"owner_id"`
	AvatarID string            `json:"avatar"`
	Channels []ChannelResponse `json:"channels"`
}

type MessageHistoryResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChannelHistoryResponse struct {
	Messages []MessageHistoryResponse `json:"messages"`
	Users    []UserResponse           `json:"users"`
}
