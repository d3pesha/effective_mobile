package action

import (
	"em/internal/adapters/api/logging"
	"em/internal/adapters/api/response"
	"em/internal/adapters/logger"
	"em/internal/entities"
	"em/internal/usecase"
	"encoding/json"
	"net/http"
)

type SongCreateAction struct {
	uc  usecase.SongCreateUseCase
	log logger.Logger
}

func NewSongCreateAction(uc usecase.SongCreateUseCase, log logger.Logger) *SongCreateAction {
	return &SongCreateAction{uc: uc, log: log}
}

func (a SongCreateAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "song_create"

	var input usecase.SongCreateInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusBadRequest,
		).Log("invalid JSON format")

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}

	if input.Song == "" {
		logging.NewError(a.log, entities.ErrSongNameIsEmpty, logKey, http.StatusBadRequest).Log("song_name_is_empty")
		response.NewError(entities.ErrSongNameIsEmpty, http.StatusBadRequest).Send(w)
		return
	}

	if input.Group == "" {
		logging.NewError(a.log, entities.ErrGroupIsEmpty, logKey, http.StatusBadRequest).Log("song_group_is_empty")
		response.NewError(entities.ErrGroupIsEmpty, http.StatusBadRequest).Send(w)
		return
	}

	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("failed to create song")
		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	logging.NewInfo(a.log, logKey, http.StatusCreated).Log("success create song")
	response.NewSuccess(output, http.StatusCreated).Send(w)
}
