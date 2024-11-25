package repo

import (
	"context"
	"em/internal/entities"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type LibrarySQL struct {
	db *gorm.DB
}

type songGORM struct {
	ID          int        `gorm:"primaryKey"`
	Group       string     `gorm:"column:group_name"`
	Song        string     `gorm:"column:song"`
	ReleaseDate *time.Time `gorm:"column:release_date"`
	Text        *string    `gorm:"column:text"`
	Link        *string    `gorm:"column:link"`
}

func (songGORM) TableName() string {
	return "library"
}

func NewLibraryRepository(db *gorm.DB) *LibrarySQL {
	return &LibrarySQL{db: db}
}

func (r *LibrarySQL) Create(ctx context.Context, inputSong entities.Song) (entities.Song, error) {
	song := songGORM{
		Group:       inputSong.Group,
		Song:        inputSong.Song,
		ReleaseDate: inputSong.ReleaseDate,
		Text:        inputSong.Text,
		Link:        inputSong.Link,
	}

	if err := r.db.WithContext(ctx).Create(&song).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return entities.Song{}, entities.ErrSongNotFound
		default:
			return entities.Song{}, errors.New("error_creating_song")
		}
	}

	inputSong.ID = song.ID

	return inputSong, nil
}

func (r *LibrarySQL) FindAll(ctx context.Context, input entities.SongFilter) ([]entities.Song, int64, error) {
	if input.Page <= 0 || input.Limit <= 0 {
		return nil, 0, errors.New("invalid pagination parameters")
	}

	var songs []songGORM
	var total int64
	offset := (input.Page - 1) * input.Limit

	query := r.db.WithContext(ctx).Model(&songGORM{})
	if input.Group != "" {
		query = query.Where("group_name ILIKE ?", "%"+input.Group+"%")
	}
	if input.Song != "" {
		query = query.Where("song ILIKE ?", "%"+input.Song+"%")
	}
	if input.Text != "" {
		query = query.Where("text ILIKE ?", "%"+input.Text+"%")
	}

	orderBys, err := parseOrderBy(input.OrderBy)
	if err != nil {
		return nil, 0, err
	}
	for _, orderBy := range orderBys {
		query = query.Order(fmt.Sprintf("%s %s", orderBy.Field, orderBy.Direction))
	}

	if err = query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("error_count_songs")
	}
	if err = query.Offset(offset).
		Limit(input.Limit).
		Find(&songs).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, 0, entities.ErrSongNotFound
		default:
			return nil, 0, errors.New("error_finding_songs")
		}
	}

	result := make([]entities.Song, len(songs))
	for i, song := range songs {
		result[i] = entities.Song{
			ID:          song.ID,
			Group:       song.Group,
			Song:        song.Song,
			ReleaseDate: song.ReleaseDate,
			Text:        song.Text,
			Link:        song.Link,
		}
	}

	return result, total, nil
}

func (r *LibrarySQL) FindByID(ctx context.Context, id int) (entities.Song, error) {
	var song songGORM
	if err := r.db.WithContext(ctx).First(&song, id).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return entities.Song{}, errors.New("song_not_found")
		default:
			return entities.Song{}, errors.New("error_finding_song_text")
		}
	}

	return entities.Song{
		ID:          song.ID,
		Group:       song.Group,
		Song:        song.Song,
		ReleaseDate: song.ReleaseDate,
		Text:        song.Text,
		Link:        song.Link,
	}, nil
}

func (r *LibrarySQL) FindSongText(ctx context.Context, id, page, limit int) ([]string, error) {
	var song songGORM
	if err := r.db.WithContext(ctx).First(&song, id).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, errors.New("song_not_found")
		default:
			return nil, errors.New("error_finding_song_text")
		}
	}

	if song.Text == nil {
		return nil, entities.ErrSongTextNotFound
	}

	verses := strings.Split(*song.Text, "\n\n")
	start := (page - 1) * limit
	end := start + limit

	if start >= len(verses) {
		return nil, nil
	}
	if end > len(verses) {
		end = len(verses)
	}

	return verses[start:end], nil
}

func (r *LibrarySQL) Delete(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&songGORM{}, id).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return errors.New("song_not_found")
		default:
			return errors.New("error_deleting_song")
		}
	}
	return nil
}

func (r *LibrarySQL) UpdateSongText(ctx context.Context, id int, text string) error {
	if err := r.db.WithContext(ctx).
		Model(&songGORM{}).
		Where("id = ?", id).
		Update("text", text).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return entities.ErrSongNotFound
		default:
			return errors.New("error_updating_song")
		}
	}
	return nil
}

func (r *LibrarySQL) FindByGroupAndSong(ctx context.Context, group, song string) (entities.Song, error) {
	var existingSong songGORM
	if err := r.db.WithContext(ctx).
		Where("group_name = ? AND song = ?", group, song).
		First(&existingSong).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return entities.Song{}, entities.ErrSongNotFound
		default:
			return entities.Song{}, errors.New("error_finding_song")
		}
	}

	return entities.Song{
		ID:          existingSong.ID,
		Group:       existingSong.Group,
		Song:        existingSong.Song,
		ReleaseDate: existingSong.ReleaseDate,
		Text:        existingSong.Text,
		Link:        existingSong.Link,
	}, nil
}
