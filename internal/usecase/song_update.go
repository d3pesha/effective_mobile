package usecase

import (
	"context"
	"em/internal/entities"
	"errors"
	"time"
)

type (
	SongUpdateUseCase interface {
		Execute(context.Context, SongUpdateInput, int) error
	}

	SongUpdateInput struct {
		Text string `json:"text"`
	}

	songUpdateInteractor struct {
		songRepo   entities.SongRepository
		ctxTimeout time.Duration
	}
)

func NewSongUpdateInteractor(
	songRepo entities.SongRepository,
	ctxTimeout time.Duration,
) SongUpdateUseCase {
	return songUpdateInteractor{
		songRepo:   songRepo,
		ctxTimeout: ctxTimeout,
	}
}

func (uc songUpdateInteractor) Execute(ctx context.Context, input SongUpdateInput, id int) error {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	existSong, err := uc.songRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, entities.ErrSongNotFound) {
			return entities.ErrSongNotFound
		}
		return err
	}

	if input.Text == "" {
		return entities.ErrTextIsEmpty
	}

	err = uc.songRepo.UpdateSongText(ctx, existSong.ID, input.Text)
	if err != nil {
		return err
	}

	return nil
}
