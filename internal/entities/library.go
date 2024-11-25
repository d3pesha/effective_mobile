package entities

import (
	"context"
	"errors"
	"time"
)

type Song struct {
	ID          int
	Group       string
	Song        string
	ReleaseDate *time.Time
	Text        *string
	Link        *string
}

type SongFilter struct {
	Group   string
	Song    string
	Text    string
	OrderBy string
	Page    int
	Limit   int
}

var (
	ErrSongNotFound      = errors.New("song_not_found")
	ErrSongAlreadyExists = errors.New("song_already_exists")
	ErrSongTextNotFound  = errors.New("song_text_not_found")
	ErrTextIsEmpty       = errors.New("song_text_is_empty")
	ErrGroupIsEmpty      = errors.New("song_group_is_empty")
	ErrSongNameIsEmpty   = errors.New("song_name_is_empty")
)

type SongRepository interface {
	FindByID(ctx context.Context, id int) (Song, error)
	FindByGroupAndSong(ctx context.Context, group, song string) (Song, error)
	FindAll(ctx context.Context, filter SongFilter) ([]Song, int64, error)
	FindSongText(ctx context.Context, id, page, limit int) ([]string, error)
	Delete(ctx context.Context, id int) error
	UpdateSongText(ctx context.Context, id int, text string) error
	Create(ctx context.Context, song Song) (Song, error)
}
