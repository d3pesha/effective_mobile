package usecase

import (
	"context"
	"em/internal/entities"
	"time"
)

type (
	SongDeleteUseCase interface {
		Execute(context.Context, SongDeleteInput) error
	}

	SongDeleteInput struct {
		ID int
	}

	songDeleteInteractor struct {
		songRepo   entities.SongRepository
		ctxTimeout time.Duration
	}
)

func NewSongDeleteInteractor(
	songRepo entities.SongRepository,
	ctxTimeout time.Duration,
) SongDeleteUseCase {
	return songDeleteInteractor{
		songRepo:   songRepo,
		ctxTimeout: ctxTimeout,
	}
}

func (uc songDeleteInteractor) Execute(ctx context.Context, input SongDeleteInput) error {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	err := uc.songRepo.Delete(ctx, input.ID)
	if err != nil {
		return err
	}

	return nil
}
