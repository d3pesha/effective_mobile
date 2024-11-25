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

type SongFindTextVersesAction struct {
	uc  usecase.SongFindTextVersesUseCase
	log logger.Logger
}

func NewSongFindTextVersesAction(uc usecase.SongFindTextVersesUseCase, log logger.Logger) *SongFindTextVersesAction {
	return &SongFindTextVersesAction{uc: uc, log: log}
}

func (a SongFindTextVersesAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "song_find_text_verses"

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
	id, err := strconv.Atoi(keysStr["id"][0])
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

	input := usecase.SongFindTextVersesInput{
		ID:    id,
		Page:  page,
		Limit: limit,
	}

	output, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		logging.NewError(a.log, err, logKey, http.StatusInternalServerError).Log("failed to find song verses")
		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	logging.NewInfo(a.log, logKey, http.StatusOK).Log("success when returning song verses")
	response.NewSuccess(output, http.StatusOK).Send(w)
}
