package story

import (
	"encoding/json"
	"net/http"
	"taleteller/api"
	"taleteller/logger"
)

func HandleStoryCreate() http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var createRequest CreateStoryRequest

		err := json.NewDecoder(req.Body).Decode(&createRequest)
		if err != nil {
			logger.Warnw(req.Context(), "error reading request body", "error", err.Error())
			api.RespondWithError(rw, http.StatusBadRequest, api.Response{
				Error: "error reading request body",
			})
			return
		}

		api.RespondWithJSON(rw, http.StatusOK, api.Response{
			Data: "done",
		})
	})
}
