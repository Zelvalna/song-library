package migrations

import (
	"github.com/Zelvalna/song-library/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Song{})
}
