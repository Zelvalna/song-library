package controllers

import (
	"github.com/Zelvalna/song-library/models"
	"github.com/Zelvalna/song-library/services"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SongController ...
type SongController struct {
	service services.SongService
}

// NewSongController ...
func NewSongController(service services.SongService) *SongController {
	return &SongController{service}
}

// AddSong godoc
// @Summary      Добавить новую песню
// @Description  Добавление новой песни
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        song  body      models.Song  true  "Песня"
// @Success      200   {object}  models.Song
// @Failure      400   {string}  string "Bad Request"
// @Router       /songs [post]
func (c *SongController) AddSong(ctx *gin.Context) {
	var input struct {
		Group string `json:"group" binding:"required"`
		Song  string `json:"song" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		logrus.Debug("Неверный запрос: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.AddSong(input.Group, input.Song)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Песня успешно добавлена"})
}

// GetSongs godoc
// @Summary      Получить список песен
// @Description  Получение списка песен с фильтрацией и пагинацией
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        group   query     string  false  "Группа"
// @Param        title   query     string  false  "Название песни"
// @Param        limit   query     int     false  "Лимит"
// @Param        offset  query     int     false  "Смещение"
// @Success      200     {array}   models.Song
// @Failure      400     {string}  string "Bad Request"
// @Router       /songs [get]
func (c *SongController) GetSongs(ctx *gin.Context) {
	filters := make(map[string]interface{})
	if group := ctx.Query("group"); group != "" {
		filters["group"] = group
	}
	if title := ctx.Query("title"); title != "" {
		filters["title"] = title
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	songs, err := c.service.GetSongs(filters, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, songs)
}

// GetSongText godoc
// @Summary      Получить текст песни
// @Description  Получение текста песни с пагинацией по куплетам
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id      path      int  true  "ID песни"
// @Param        page    query     int  false "Страница"
// @Param        size    query     int  false "Размер страницы"
// @Success      200     {array}   string
// @Failure      400     {string}  string "Bad Request"
// @Router       /songs/{id}/text [get]
func (c *SongController) GetSongText(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "1"))

	song, err := c.service.GetSongByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	verses := paginateText(song.Text, page, size)
	ctx.JSON(http.StatusOK, verses)
}

// UpdateSong godoc
// @Summary      Обновить данные песни
// @Description  Обновление данных песни
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id    path      int          true  "ID песни"
// @Param        song  body      models.Song  true  "Песня"
// @Success      200   {object}  models.Song
// @Failure      400   {string}  string "Bad Request"
// @Router       /songs/{id} [put]
func (c *SongController) UpdateSong(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var song models.Song
	if err := ctx.ShouldBindJSON(&song); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	song.ID = uint(id)

	err = c.service.UpdateSong(&song)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, song)
}

// DeleteSong godoc
// @Summary      Удалить песню
// @Description  Удаление песни
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id    path      int  true  "ID песни"
// @Success      200   {string}  string "OK"
// @Failure      400   {string}  string "Bad Request"
// @Router       /songs/{id} [delete]
func (c *SongController) DeleteSong(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = c.service.DeleteSong(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Песня успешно удалена"})
}

func paginateText(text string, page, size int) []string {
	verses := splitIntoVerses(text)
	start := (page - 1) * size
	end := start + size

	if start > len(verses) {
		return []string{}
	}
	if end > len(verses) {
		end = len(verses)
	}

	return verses[start:end]
}

func splitIntoVerses(text string) []string {
	re := regexp.MustCompile(`\n\s*\n`)
	verses := re.Split(strings.TrimSpace(text), -1)
	return verses // Реализуйте логику разделения текста на куплеты
}
