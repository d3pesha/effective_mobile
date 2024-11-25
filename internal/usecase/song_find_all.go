package usecase

import (
	"context"
	"em/internal/entities"
	"time"
)

type (
	SongFindAllUseCase interface {
		Execute(context.Context, SongFindAllInput) ([]SongFindAllOutput, int64, error)
	}

	SongFindAllPresenter interface {
		Output([]entities.Song) []SongFindAllOutput
	}

	SongFindAllInput struct {
		Group   string
		Song    string
		OrderBy string
		Text    string
		Page    int
		Limit   int
	}

	SongFindAllOutput struct {
		ID          int        `json:"id"`
		Group       string     `json:"group"`
		Song        string     `json:"song"`
		ReleaseDate *time.Time `json:"releaseDate,omitempty"`
		Link        *string    `json:"link,omitempty"`
		Text        *string    `json:"text,omitempty"`
	}

	songFindAllInteractor struct {
		songRepo   entities.SongRepository
		presenter  SongFindAllPresenter
		ctxTimeout time.Duration
	}
)

func NewSongFindAllInteractor(
	songRepo entities.SongRepository,
	presenter SongFindAllPresenter,
	ctxTimeout time.Duration,
) SongFindAllUseCase {
	return songFindAllInteractor{
		songRepo:   songRepo,
		presenter:  presenter,
		ctxTimeout: ctxTimeout,
	}
}

func (uc songFindAllInteractor) Execute(ctx context.Context, input SongFindAllInput) ([]SongFindAllOutput, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ctxTimeout)
	defer cancel()

	filter := entities.SongFilter{
		Group:   input.Group,
		Song:    input.Song,
		Text:    input.Text,
		OrderBy: input.OrderBy,
		Page:    input.Page,
		Limit:   input.Limit,
	}

	songs, total, err := uc.songRepo.FindAll(ctx, filter)
	if err != nil {
		return []SongFindAllOutput{}, 0, err
	}

	return uc.presenter.Output(songs), total, nil
}
