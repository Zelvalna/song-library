package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Zelvalna/song-library/models"
	"github.com/Zelvalna/song-library/repositories"
	"github.com/sirupsen/logrus"
)

type SongService interface {
	AddSong(group, title string) error
	GetSongs(filters map[string]interface{}, limit, offset int) ([]models.Song, error)
	GetSongByID(id uint) (*models.Song, error)
	UpdateSong(song *models.Song) error
	DeleteSong(id uint) error
}

type songService struct {
	repo   repositories.SongRepository
	apiURL string
}

func NewSongService(repo repositories.SongRepository, apiURL string) SongService {
	return &songService{repo, apiURL}
}

func (s *songService) AddSong(group, title string) error {
	// Запрос к внешнему API
	resp, err := http.Get(fmt.Sprintf("%s?group=%s&song=%s", s.apiURL, group, title))
	if err != nil {
		logrus.Error("Ошибка при запросе к внешнему API: ", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Error("Внешний API вернул статус: ", resp.StatusCode)
		return fmt.Errorf("external API error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Ошибка при чтении ответа внешнего API: ", err)
		return err
	}

	var songDetail struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}

	err = json.Unmarshal(body, &songDetail)
	if err != nil {
		logrus.Error("Ошибка при разборе JSON ответа: ", err)
		return err
	}

	song := &models.Song{
		Group:       group,
		Title:       title,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	err = s.repo.Create(song)
	if err != nil {
		logrus.Error("Ошибка при сохранении песни в БД: ", err)
	}

	return err
}

func (s *songService) GetSongs(filters map[string]interface{}, limit, offset int) ([]models.Song, error) {
	return s.repo.GetAll(filters, limit, offset)
}

func (s *songService) GetSongByID(id uint) (*models.Song, error) {
	return s.repo.GetByID(id)
}

func (s *songService) UpdateSong(song *models.Song) error {
	return s.repo.Update(song)
}

func (s *songService) DeleteSong(id uint) error {
	return s.repo.Delete(id)
}
