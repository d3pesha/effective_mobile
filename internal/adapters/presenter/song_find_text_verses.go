package presenter

import (
	"em/internal/entities"
	"em/internal/usecase"
)

type songFindTextVersesPresenter struct{}

func NewSongFindTextVersesPresenter() usecase.SongFindTextVersesPresenter {
	return &songFindTextVersesPresenter{}
}

func (p *songFindTextVersesPresenter) Output(song entities.Song, verses []string) usecase.SongFindTextVersesOutput {
	return usecase.SongFindTextVersesOutput{
		Group:       song.Group,
		Song:        song.Song,
		ReleaseDate: song.ReleaseDate,
		Link:        song.Link,
		Verses:      verses,
	}
}
