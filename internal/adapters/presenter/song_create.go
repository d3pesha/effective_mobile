package presenter

import (
	"em/internal/entities"
	"em/internal/usecase"
)

type songCreatePresenter struct{}

func NewSongCreatePresenter() usecase.SongCreatePresenter {
	return &songCreatePresenter{}
}

func (p *songCreatePresenter) Output(song entities.Song) usecase.SongCreateOutput {
	return usecase.SongCreateOutput{
		Group:       song.Group,
		Song:        song.Song,
		ReleaseDate: song.ReleaseDate,
		Text:        song.Text,
		Link:        song.Link,
	}
}
