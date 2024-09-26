package repositories

import (
	"github.com/Zelvalna/song-library/models"
	"gorm.io/gorm"
)

type SongRepository interface {
	Create(song *models.Song) error
	GetAll(filters map[string]interface{}, limit, offset int) ([]models.Song, error)
	GetByID(id uint) (*models.Song, error)
	Update(song *models.Song) error
	Delete(id uint) error
}

type songRepository struct {
	db *gorm.DB
}

func NewSongRepository(db *gorm.DB) SongRepository {
	return &songRepository{db}
}

func (r *songRepository) Create(song *models.Song) error {
	return r.db.Create(song).Error
}

func (r *songRepository) GetAll(filters map[string]interface{}, limit, offset int) ([]models.Song, error) {
	var songs []models.Song
	query := r.db.Model(&models.Song{}).Where(filters).Limit(limit).Offset(offset)
	err := query.Find(&songs).Error
	return songs, err
}

func (r *songRepository) GetByID(id uint) (*models.Song, error) {
	var song models.Song
	err := r.db.First(&song, id).Error
	return &song, err
}

func (r *songRepository) Update(song *models.Song) error {
	return r.db.Save(song).Error
}

func (r *songRepository) Delete(id uint) error {
	return r.db.Delete(&models.Song{}, id).Error
}
