package story

import (
	"encoding/json"
	"github.com/gorilla/mux"
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

		resp, err := service.Create(req.Context(), createRequest)

		api.RespondWithJSON(rw, http.StatusOK, api.Response{
			Data: resp,
		})
	})
}

func HandleGetStory(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var err error
		ctx := req.Context()
		vars := mux.Vars(req)
		id := vars["id"]
		status, err := service.GetStory(ctx, id)
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
		status := req.URL.Query().Get("status")
		resp, err := service.List(req.Context(), status)
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

func HandlePublishStory(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		var request []UpdateSceneOrderReq
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			logger.Warnw(req.Context(), "error reading request body", "error", err.Error())
			api.RespondWithError(rw, http.StatusBadRequest, api.Response{
				Error: "error reading request body",
			})
			return
		}

		path, err := service.Publish(req.Context(), request, id)
		if err != nil {
			logger.Warnw(req.Context(), "error listing stories", "error", err.Error())
			api.RespondWithError(rw, http.StatusBadRequest, api.Response{
				Error: "error listing stories",
			})
			return
		}

		api.RespondWithJSON(rw, http.StatusOK, api.Response{
			Data: path,
		})
	})
}
