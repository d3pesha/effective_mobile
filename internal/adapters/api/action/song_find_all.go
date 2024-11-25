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

type SongFindAllAction struct {
	uc  usecase.SongFindAllUseCase
	log logger.Logger
}

func NewSongFindAllAction(uc usecase.SongFindAllUseCase, log logger.Logger) *SongFindAllAction {
	return &SongFindAllAction{uc: uc, log: log}
}

func (a SongFindAllAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "song_find_all"

	var input usecase.SongFindAllInput
	var keysStr = r.URL.Query()

	page, err := strconv.Atoi(keysStr["page"][0])
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusBadRequest,
		).Log(fmt.Sprintf("invalid page format"))

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}

	limit, err := strconv.Atoi(keysStr["limit"][0])
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusBadRequest,
		).Log(fmt.Sprintf("invalid limit format"))

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	input.Limit = limit
	input.Page = page
	input.Group = keysStr.Get("group")
	input.Song = keysStr.Get("song")
	input.OrderBy = keysStr.Get("orderBy")
	input.Text = keysStr.Get("text")

	output, total, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		logging.NewError(a.log, err, logKey, http.StatusInternalServerError).Log("failed to find songs")
		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	logging.NewInfo(a.log, logKey, http.StatusOK).Log("success when returning songs")
	response.NewSuccessList(output, total, http.StatusOK).Send(w)
}
