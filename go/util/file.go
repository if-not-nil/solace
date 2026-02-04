package util

import (
	"errors"
	"fmt"
	"mime"
	"os"
	"solace/orm"

	"gorm.io/gorm"
)

const UploadDir = "data/uploads"

func GetFile(id string) []byte {
	path := fmt.Sprintf("%s/%s", UploadDir, id)
	file, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return file
}

const MaxFileSize = 5 * 1024 * 1024 // 5MB

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

// TODO: compress files
func UploadFile(db *gorm.DB, bytes []byte, contentType string) (id string, err error) {
	if len(bytes) > MaxFileSize {
		return "", errors.New("file size exceeds 5MB limit")
	}
	if !allowedImageTypes[contentType] {
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}

	exts, _ := mime.ExtensionsByType(contentType)
	ext := ""
	if len(exts) > 0 {
		ext = exts[0] // like ".png"
	}

	// create dummy record first to get ID
	record := orm.File{Type: contentType}
	if err := db.Create(&record).Error; err != nil {
		return "", err
	}

	newID := record.ID + ext                            //  "abc123.png"
	if err := db.Model(&record).Update("id", newID).Error; err != nil {
		return "", err
	}

	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		return "", err
	}
	path := fmt.Sprintf("%s/%s", UploadDir, newID)
	if err := os.WriteFile(path, bytes, 0644); err != nil {
		return "", err
	}

	return newID, nil
}
