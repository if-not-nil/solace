package orm

import (
	"gorm.io/gorm"
)

func (u *User) GetMembership(db *gorm.DB, serverID string) (Membership, error) {
	var membership Membership
	if err := db.Preload("User").Preload("Server").First(&membership, "user_id = ? AND server_id = ?", u.ID, serverID).Error; err != nil {
		return Membership{}, err
	}
	return membership, nil
}

func (u *User) InServer(db *gorm.DB, serverID string) bool {
	membership, err := u.GetMembership(db, serverID)
	if err != nil || membership.Banned {
		return false
	}
	return true
}

func (u *User) ByID(db *gorm.DB, id string) error {
	if err := db.First(&u, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) ByName(db *gorm.DB, name string) error {
	if err := db.First(&u, "name = ?", name).Error; err != nil {
		return err
	}
	return nil
}

func (s *Server) CreateChannel(db *gorm.DB, channel *Channel) error {
	if err := db.Create(channel).Error; err != nil {
		return err
	}
	s.Channels = append(s.Channels, *channel)
	return db.Save(s).Error
}
