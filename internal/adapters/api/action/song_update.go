package action

import (
	"em/internal/adapters/api/logging"
	"em/internal/adapters/api/response"
	"em/internal/adapters/logger"
	"em/internal/usecase"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type SongUpdateAction struct {
	uc  usecase.SongUpdateUseCase
	log logger.Logger
}

func NewSongUpdateAction(uc usecase.SongUpdateUseCase, log logger.Logger) *SongUpdateAction {
	return &SongUpdateAction{uc: uc, log: log}
}

func (a SongUpdateAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "song_update"

	var input usecase.SongUpdateInput

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

	err = a.uc.Execute(r.Context(), input, id)
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusInternalServerError,
		).Log("failed to update song")
		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	logging.NewInfo(a.log, logKey, http.StatusOK).Log("success update song")
	response.NewSuccess(nil, http.StatusOK).Send(w)
}
