package usecase

import (
	"context"
	"em/internal/entities"
	"errors"
	"time"
)

type (
	SongCreateUseCase interface {
		Execute(context.Context, SongCreateInput) (SongCreateOutput, error)
	}

	SongCreatePresenter interface {
		Output(entities.Song) SongCreateOutput
	}

	SongCreateInput struct {
		Group string `json:"group"`
		Song  string `json:"song"`
	}

	SongCreateOutput struct {
		Group       string     `json:"group"`
		Song        string     `json:"song"`
		ReleaseDate *time.Time `json:"releaseDate"`
		Text        *string    `json:"text,omitempty"`
		Link        *string    `json:"link,omitempty"`
	}

	songCreateInteractor struct {
		songRepo   entities.SongRepository
		presenter  SongCreatePresenter
		ctxTimeout time.Duration
	}
)

func NewSongCreateInteractor(
	songRepo entities.SongRepository,
	presenter SongCreatePresenter,
	ctxTimeout time.Duration,
) SongCreateUseCase {
	return songCreateInteractor{
		songRepo:   songRepo,
		presenter:  presenter,
		ctxTimeout: ctxTimeout,
	}
}

func (uc songCreateInteractor) Execute(ctx context.Context, input SongCreateInput) (SongCreateOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	existSong, err := uc.songRepo.FindByGroupAndSong(ctx, input.Group, input.Song)
	if err == nil && existSong.ID != 0 {
		return SongCreateOutput{}, entities.ErrSongAlreadyExists
	} else if !errors.Is(err, entities.ErrSongNotFound) {
		return SongCreateOutput{}, err
	}

	song := entities.Song{
		Group: input.Group,
		Song:  input.Song,
	}

	newSong, err := uc.songRepo.Create(ctx, song)
	if err != nil {
		return SongCreateOutput{}, err
	}

	return uc.presenter.Output(newSong), nil
}
