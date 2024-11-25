package presenter

import (
	"em/internal/entities"
	"em/internal/usecase"
)

type songFindAllPresenter struct{}

func NewSongFindAllPresenter() usecase.SongFindAllPresenter {
	return &songFindAllPresenter{}
}

func (p songFindAllPresenter) Output(songs []entities.Song) []usecase.SongFindAllOutput {
	outputs := make([]usecase.SongFindAllOutput, len(songs))
	for i, song := range songs {
		outputs[i] = usecase.SongFindAllOutput{
			ID:          song.ID,
			Group:       song.Group,
			Song:        song.Song,
			ReleaseDate: song.ReleaseDate,
			Link:        song.Link,
			Text:        song.Text,
		}
	}
	return outputs
}
