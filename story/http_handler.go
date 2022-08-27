package story

import (
	"encoding/json"
	"net/http"
	"taleteller/api"
	"taleteller/logger"
)

func HandleStoryCreate(service Service) http.HandlerFunc {
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

		err = service.Create(req.Context(), createRequest)

		api.RespondWithJSON(rw, http.StatusOK, api.Response{
			Data: "done",
		})
	})
}

func HandleGetStoryStatus(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var err error
		ctx := req.Context()
		id := req.URL.Query().Get("id")
		status, err := service.GetStoryStatus(ctx, id)
		if err != nil {
			logger.Errorw(req.Context(), "error getting status", "error", err.Error(), "storyID", id)
			api.RespondWithError(rw, http.StatusInternalServerError, api.Response{
				Error: "error getting status",
			})
			return
		}

		api.RespondWithJSON(rw, http.StatusOK, api.Response{
			Data: status,
		})
	})
}

func HandleListStories(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	   resp, err := service.List(req.Context())
	   if err != nil {
		  logger.Warnw(req.Context(), "error listing stories", "error", err.Error())
		  api.RespondWithError(rw, http.StatusBadRequest, api.Response{
			 Error: "error listing stories",
		  })
		  return
	   }
 
	   api.RespondWithJSON(rw, http.StatusOK, api.Response{
		  Data: resp,
	   })
	})
 }