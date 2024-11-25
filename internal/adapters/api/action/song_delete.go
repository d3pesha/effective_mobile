package action

import (
	"em/internal/adapters/api/logging"
	"em/internal/adapters/api/response"
	"em/internal/adapters/logger"
	"em/internal/usecase"
	"fmt"
	"net/http"
	"strconv"
)

type SongDeleteAction struct {
	uc  usecase.SongDeleteUseCase
	log logger.Logger
}

func NewSongDeleteAction(uc usecase.SongDeleteUseCase, log logger.Logger) *SongDeleteAction {
	return &SongDeleteAction{uc: uc, log: log}
}

func (a SongDeleteAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "song_delete"

	var input usecase.SongDeleteInput

	var idStr = r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusBadRequest,
		).Log(fmt.Sprintf("invalid id format"))

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	input.ID = id

	err = a.uc.Execute(r.Context(), input)
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("failed to delete song")
		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	logging.NewInfo(a.log, logKey, http.StatusOK).Log("success delete song")
	response.NewSuccess(nil, http.StatusOK).Send(w)
}
