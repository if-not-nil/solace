package orm

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	db.AutoMigrate(
		&Message{}, &User{}, &Server{}, &Channel{},
		&Membership{}, &Role{}, &UserRole{},
		&ChannelPermissionOverride{}, &File{},
		&InviteLink{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
