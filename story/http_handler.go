package story

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
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

func HandleCreateScene(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		vars := mux.Vars(req)
		id := vars["id"]
		if len(id) == 0 {
			logger.Errorw(ctx, "error getting id from url")
			api.RespondWithError(rw, http.StatusBadRequest, api.Response{
				Error:     "bad request",
				ErrorCode: "BAD_REQUEST",
			})
			return
		}

		reqByte, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.Errorw(ctx, "error while reading request body", "error", err.Error())
			api.RespondWithError(rw, http.StatusBadRequest, api.Response{
				Error:     "bad request",
				ErrorCode: "BAD_REQUEST",
			})
			return
		}

		var createSceneRequest CreateSceneRequest
		err = json.Unmarshal(reqByte, &createSceneRequest)
		if err != nil {
			logger.Errorw(ctx, "error while reading request body", "error", err.Error())
			api.RespondWithError(rw, http.StatusInternalServerError, api.Response{
				Error:     "error unmarshalling request",
				ErrorCode: "INTERNAL_SERVER_ERROR",
			})
			return
		}

		response, err := service.CreateScene(ctx, createSceneRequest)
		if err != nil {
			logger.Errorw(ctx, "error generating scene", "error", err.Error())
			api.RespondWithError(rw, http.StatusInternalServerError, api.Response{
				Error:     "error generating scene",
				ErrorCode: "INTERNAL_SERVER_ERROR",
			})
			return
		}

		api.RespondWithJSON(rw, http.StatusOK, api.Response{
			Data: response,
		})
	})
}
