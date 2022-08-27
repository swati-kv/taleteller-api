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

func HandleGetStory(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var err error
		ctx := req.Context()
		vars := mux.Vars(req)
		id := vars["id"]
		status, err := service.GetStory(ctx, id)
		if err != nil {
			logger.Errorw(req.Context(), "error getting story", "error", err.Error(), "storyID", id)
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
func HandleUpdateScene(service Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var err error
		var image Image
		ctx := req.Context()

		vars := mux.Vars(req)
		storyID := vars["story_id"]
		sceneID := vars["scene_id"]

		reqByte, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.Errorw(req.Context(), "error reading request body", "error", err.Error())
			return
		}
		err = json.Unmarshal(reqByte, &image)
		if err != nil {
			logger.Errorw(req.Context(), "error unmarshalling request body", "error", err.Error())
			return
		}
		if err != nil {
			logger.Errorw(req.Context(), "error reading update scene request body", "error", err.Error())
			api.RespondWithError(rw, http.StatusBadRequest, api.Response{
				Error: "error reading update scene request body",
			})
			return
		}

		scene, err := service.UpdateScene(ctx, storyID, sceneID, image.SelectedImage)
		if err != nil {
			logger.Errorw(req.Context(), "error updating scene", "error", err.Error(), "storyID", storyID, "sceneID", sceneID)
			api.RespondWithError(rw, http.StatusInternalServerError, api.Response{
				Error: "error updating scene",
			})
			return
		}
		api.RespondWithJSON(rw, http.StatusOK, api.Response{
			Data: scene,
		})
	})
}
