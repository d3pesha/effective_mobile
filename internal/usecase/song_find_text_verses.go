package usecase

import (
	"context"
	"em/internal/entities"
	"time"
)

type (
	SongFindTextVersesUseCase interface {
		Execute(context.Context, SongFindTextVersesInput) (SongFindTextVersesOutput, error)
	}

	SongFindTextVersesPresenter interface {
		Output(entities.Song, []string) SongFindTextVersesOutput
	}

	SongFindTextVersesInput struct {
		ID    int
		Page  int
		Limit int
	}

	SongFindTextVersesOutput struct {
		Group       string     `json:"group"`
		Song        string     `json:"song"`
		ReleaseDate *time.Time `json:"releaseDate,omitempty"`
		Link        *string    `json:"link,omitempty"`
		Verses      []string   `json:"verses"`
	}

	songFindTextVersesInteractor struct {
		songRepo   entities.SongRepository
		presenter  SongFindTextVersesPresenter
		ctxTimeout time.Duration
	}
)

func NewSongFindTextVersesInteractor(
	songRepo entities.SongRepository,
	presenter SongFindTextVersesPresenter,
	ctxTimeout time.Duration,
) SongFindTextVersesUseCase {
	return songFindTextVersesInteractor{
		songRepo:   songRepo,
		presenter:  presenter,
		ctxTimeout: ctxTimeout,
	}
}

func (uc songFindTextVersesInteractor) Execute(ctx context.Context, input SongFindTextVersesInput) (SongFindTextVersesOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	song, err := uc.songRepo.FindByID(ctx, input.ID)
	if err != nil {
		return SongFindTextVersesOutput{}, err
	}

	verses, err := uc.songRepo.FindSongText(ctx, input.ID, input.Page, input.Limit)
	if err != nil {
		return SongFindTextVersesOutput{}, err
	}

	return uc.presenter.Output(song, verses), nil
}
