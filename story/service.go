package story

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gojektech/heimdall/v6/httpclient"
	"github.com/gojektech/heimdall/v6/plugins"
	"io/ioutil"
	"net/http"
	"sync"
	"taleteller/api"
	"taleteller/app"
	"taleteller/logger"
	"taleteller/store"
	"taleteller/utils"
	"time"
)

type Service interface {
	Create(ctx context.Context, createRequest CreateStoryRequest) (response CreateStoryResponse, err error)
	GetStory(ctx context.Context, storyID string) (storyDetails store.Story, err error)
	List(ctx context.Context, status string) (stories []store.Story, err error)
	Publish(ctx context.Context, req []UpdateSceneOrderReq, storyID string) (path string, err error)
	UpdateScene(ctx context.Context, storyID string, sceneID string, selectedImage string) (updatedScene store.Scene, err error)
	GetScene(ctx context.Context) (response GetSceneResponse, err error)
	CreateScene(ctx context.Context, createSceneRequest CreateSceneRequest) (response CreateSceneResponse, err error)
}

type service struct {
	httpClient      *httpclient.Client
	store           store.StoryStorer
	generatorUtils  utils.IDGeneratorUtils
	pyServerBaseURL string
}

func NewService(store store.StoryStorer, pyServerBaseURL string, generatorUtils utils.IDGeneratorUtils) Service {
	timeout := 3000 * time.Minute
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))
	requestLogger := plugins.NewRequestLogger(nil, nil)
	client.AddPlugin(requestLogger)
	return &service{
		store:           store,
		generatorUtils:  generatorUtils,
		pyServerBaseURL: pyServerBaseURL,
		httpClient:      client,
	}
}

func (s *service) Create(ctx context.Context, createRequest CreateStoryRequest) (response CreateStoryResponse, err error) {

	storyID, err := s.generatorUtils.GenerateIDWithPrefix("sto_")
	if err != nil {
		logger.Error(ctx, "error generating ID", err.Error())
		return
	}

	req := store.Story{
		StoryID:     storyID,
		Name:        createRequest.Name,
		Description: createRequest.Description,
		Mood:        createRequest.Mood,
		Category:    createRequest.Category,
		CustomerID:  "cus_123",
		Status:      "processing",
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}

	err = s.store.Create(ctx, req)
	if err != nil {
		logger.Error(ctx, "error creating story", err.Error())
		return
	}
	response.StoryID = storyID
	return
}

var wg sync.WaitGroup

func (s *service) CreateScene(ctx context.Context, createSceneRequest CreateSceneRequest) (response CreateSceneResponse, err error) {
	idGenerator := utils.NewGeneratorUtils()

	sceneID, err := idGenerator.GenerateIDWithPrefix("scene_")
	if err != nil {
		logger.Errorw(ctx, "error generating scene id", "error", err.Error())
		return
	}

	storyID := ctx.Value("story-id").(string)

	err = s.createSceneEntry(sceneID, createSceneRequest, storyID)
	if err != nil {
		logger.Errorw(ctx, "error while insertinf new scene", "error", err.Error())
		return
	}

	//TODO handle defer
	//wg.Add(1)
	go func() {
		s.generateImage(ctx, createSceneRequest, sceneID)
		s.generateAudio(ctx, createSceneRequest, sceneID)
		//wg.Done()
	}()

	response.Status = statusProcessing
	response.SceneID = sceneID
	wg.Wait()
	return
}

func (s *service) generateImage(ctx context.Context, createSceneRequest CreateSceneRequest, sceneID string) (response PyImageResponse, err error) {
	ctx = context.Background()
	defer s.processGenerateImage(response, sceneID, err)

	url := fmt.Sprintf("%s/%s", s.pyServerBaseURL, createImageCraiyonEndPoint)
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["request-id"] = sceneID
	pyImageRequest := PyImageRequest{
		Prompt: createSceneRequest.Prompt,
		Count:  createSceneRequest.ImageCount,
	}

	fmt.Println("body - ", pyImageRequest)
	requestJSON, err := json.Marshal(pyImageRequest)
	if err != nil {
		logger.Errorw(ctx, "error marshalling downstream request", "error", err.Error())
		return
	}
	fmt.Println("body json- ", string(requestJSON))

	httpResponse, err := api.Post(ctx, url, requestJSON, header, s.httpClient)
	if err != nil {
		logger.Errorw(ctx, "error while api request create scene", "error", err.Error())
		return
	}

	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		logger.Errorw(ctx, "error while reading py server response body", "error", err.Error())
		return
	}

	if len(response.Error) != 0 {
		logger.Errorw(ctx, "error while generating image from py server", "error", response.Error)
		err = errors.New("error generating image from py server")
		return
	}

	logger.Infow(ctx, "generated a scene")
	s.processGenerateImage(response, sceneID, err)
	return
}

func (s *service) processGenerateImage(response PyImageResponse, sceneID string, err error) {
	if err != nil {
		return
	}
	if len(response.Error) != 0 {
		return
	}
	logger.Infow(context.Background(), "processing generated image - ")

	for x, image := range response.Data.GeneratedImage {
		imageID, err := utils.NewGeneratorUtils().GenerateIDWithPrefix("image_")
		if err != nil {
			logger.Errorw(context.Background(), "error while generating image id", "error", err.Error())
			return
		}

		if x == 0 {
			_, err := s.store.UpdateScene(context.Background(), "", sceneID, imageID)
			if err != nil {
				logger.Errorw(context.Background(), "error while generating image id", "error", err.Error())
				return
			}
		}

		awsService := utils.NewAWSService()
		decodedImage, _ := base64.StdEncoding.DecodeString(image)

		awsRequest := utils.UploadS3{
			File:       decodedImage,
			FileType:   "image",
			FileFormat: "jpeg",
			FileName:   imageID,
		}
		config := app.InitServiceConfig()
		bucket := config.GetAWSGeneratedAssetsBucket()
		link, err := awsService.UploadFile(bucket, awsRequest, true)
		if err != nil {
			logger.Errorw(context.Background(), "error while uploading image into aws", "error", err.Error())
			return
		}
		fmt.Println("link - ", link)

		insertImageStoreRequest := store.InsertImageRequest{
			ID:        imageID,
			ImagePath: link,
			SceneID:   sceneID,
		}

		err = s.store.InsertImage(context.Background(), insertImageStoreRequest)
		if err != nil {
			logger.Errorw(context.Background(), "error while inserting image link", "error", err.Error())
			return
		}

	}

	//s.store.UpdateSceneStatus(context.Background(), "image", sceneID, statusImageDone)
	return
}

func (s *service) generateAudio(ctx context.Context, createSceneRequest CreateSceneRequest, sceneID string) (response PyAudioResponse, err error) {
	//wg.Add(1)
	//createSceneRequest.Audio = "Test audio 123"
	ctx = context.Background()
	//defer s.processGenerateImage(response, sceneID, err)
	logger.Infow(ctx, "generating audio --- ")
	url := fmt.Sprintf("%s/%s", s.pyServerBaseURL, createAudioEndPoint)
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["request-id"] = sceneID
	pyAudioRequest := PyAudioRequest{
		Prompt:   createSceneRequest.Audio,
		Language: "en",
	}

	requestJSON, err := json.Marshal(pyAudioRequest)
	if err != nil {
		logger.Errorw(ctx, "error marshalling downstream request", "error", err.Error())
		return
	}
	fmt.Println("body json- ", string(requestJSON))

	httpResponse, err := api.Post(ctx, url, requestJSON, header, s.httpClient)
	if err != nil {
		logger.Errorw(ctx, "error while api request create scene", "error", err.Error())
		return
	}

	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		logger.Errorw(ctx, "error while reading py server response body", "error", err.Error())
		return
	}

	if len(response.Error) != 0 {
		logger.Errorw(ctx, "error while generating image from py server", "error", response.Error)
		err = errors.New("error generating image from py server")
		return
	}

	s.processGeneratedAudio(response, sceneID, err)

	return
}

func (s *service) processGeneratedAudio(response PyAudioResponse, sceneID string, err error) {
	if err != nil {
		return
	}
	if len(response.Error) != 0 {
		return
	}
	logger.Infow(context.Background(), "processing generated audio - ")

	audioID, err := utils.NewGeneratorUtils().GenerateIDWithPrefix("audio_")
	if err != nil {
		logger.Errorw(context.Background(), "error while generating image id", "error", err.Error())
		return
	}

	awsService := utils.NewAWSService()
	decodedAudio, _ := base64.StdEncoding.DecodeString(response.Data)
	awsRequest := utils.UploadS3{
		File:       decodedAudio,
		FileType:   "audio",
		FileFormat: "mp3",
		FileName:   audioID,
	}
	config := app.InitServiceConfig()
	bucket := config.GetAWSGeneratedAssetsBucket()
	link, err := awsService.UploadFile(bucket, awsRequest, true)
	if err != nil {
		logger.Errorw(context.Background(), "error while uploading image into aws", "error", err.Error())
		return
	}
	fmt.Println("link - ", link)

	insertAudioRequest := store.InsertAudioRequest{
		ID:        audioID,
		SceneID:   sceneID,
		AudioPath: link,
	}

	err = s.store.InsertAudio(context.Background(), insertAudioRequest)
	if err != nil {
		logger.Errorw(context.Background(), "error while inserting image link", "error", err.Error())
		return
	}

	err = s.store.UpdateSceneAudio(context.Background(), audioID, sceneID)
	if err != nil {
		logger.Errorw(context.Background(), "error updating audio in scene", "error", err.Error())
		return
	}

	s.store.UpdateSceneStatus(context.Background(), "audio", sceneID, "completed")

	return
}

func (s *service) createSceneEntry(id string, request CreateSceneRequest, storyID string) (err error) {
	storeRequest := store.CreateSceneRequest{
		SceneID:         id,
		Status:          statusStarted,
		StoryID:         storyID,
		SceneNumber:     request.SceneNumber,
		BackgroundMusic: request.BackgroundMusic,
	}

	err = s.store.CreateScene(context.Background(), storeRequest)
	if err != nil {
		logger.Errorw(context.Background(), "error while inserting new story into db", "error", err.Error())
		return
	}
	return
}

func (s *service) GetStory(ctx context.Context, storyID string) (storyDetails store.Story, err error) {
	storyDetails, err = s.store.GetStoryByID(ctx, storyID)
	if err != nil {
		logger.Errorw(ctx, "error getting story by story ID", "error", err.Error())
		return
	}
	return
}

func (s *service) GetScene(ctx context.Context) (response GetSceneResponse, err error) {
	storyID := ctx.Value("story-id").(string)
	sceneID := ctx.Value("scene-id").(string)

	dbResponse, err := s.store.GetSceneByID(ctx, sceneID, storyID)
	if err != nil {
		logger.Errorw(ctx, "error while getting scene", "error", err.Error(), "sceneID", sceneID, "storyID", storyID)
		return
	}
	if len(dbResponse) == 0 {
		response.Status = "processing"
		logger.Warnw(ctx, "no rowns selected", "sceneID", sceneID, "storyID", storyID)
		return
	}
	for _, resp := range dbResponse {
		fmt.Println("in for")
		if resp.Status != "completed" {
			response.Status = resp.Status
			fmt.Println("return here - ", resp.Status)
			return
		}
		var imageDetails ImageDetails
		imageDetails = ImageDetails{
			ImageID:   resp.ImageID,
			ImagePath: resp.ImagePath,
		}

		response.Images = append(response.Images, imageDetails)
	}
	response.Status = "completed"
	return
}

func (s *service) List(ctx context.Context, status string) (stories []store.Story, err error) {
	stories, err = s.store.List(ctx, status)
	if err != nil {
		logger.Error(ctx, "error getting stories", err.Error())
		return
	}
	return
}

func (s *service) Publish(ctx context.Context, req []UpdateSceneOrderReq, storyID string) (path string, err error) {
	var images []string
	var generatedAudios []string
	var backgroundAudios []string
	for _, scene := range req {
		resp, sceneErr := s.store.UpdateSceneOrder(ctx, scene.SceneID, scene.SceneNumber, storyID)
		if sceneErr != nil {
			err = sceneErr
			logger.Error(ctx, "error updating scene", sceneErr.Error())
			return
		}
		imageResp, awsErr := http.Get(resp.SelectedImagePath)
		logger.Info(ctx, "ivde ", imageResp)
		if awsErr != nil {
			return
		}
		img, _ := ioutil.ReadAll(imageResp.Body)
		image64 := base64.StdEncoding.EncodeToString(img)
		images = append(images, image64)
		generatedAudioResp, awsErr := http.Get(resp.BackgroundAudioPath)
		if awsErr != nil {
			return
		}
		ga, _ := ioutil.ReadAll(generatedAudioResp.Body)
		ga64 := base64.StdEncoding.EncodeToString(ga)
		generatedAudios = append(generatedAudios, ga64)
		backgroundAudioResp, awsErr := http.Get(resp.GeneratedAudioPath)
		if awsErr != nil {
			return
		}
		ba, _ := ioutil.ReadAll(backgroundAudioResp.Body)
		ba64 := base64.StdEncoding.EncodeToString(ba)
		backgroundAudios = append(backgroundAudios, ba64)
	}
	r := PublishRequest{
		Images:      images,
		ImageFormat: "jpeg",
		Audios:      generatedAudios,
		AudioFormat: "mp3",
		BGM:         backgroundAudios,
	}
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	url := fmt.Sprintf("%s/%s", s.pyServerBaseURL, publishEndPoint)
	requestJSON, err := json.Marshal(r)
	if err != nil {
		logger.Errorw(ctx, "error marshalling downstream request", "error", err.Error())
		return
	}

	httpResponse, err := api.Post(ctx, url, requestJSON, header, s.httpClient)
	if err != nil {
		logger.Errorw(ctx, "error in publishing", "error", err.Error(), "url", url, "request_body", requestJSON, "headers", header)
		return
	}
	var vid PyVideoResponse
	err = json.NewDecoder(httpResponse.Body).Decode(&vid)
	if err != nil {
		logger.Errorw(ctx, "error while reading py server response body", "error", err.Error())
		return
	}

	data, _ := base64.StdEncoding.DecodeString(vid.Data)

	awsService := utils.NewAWSService()
	awsRequest := utils.UploadS3{
		FileName:   "",
		FileType:   "video",
		FileFormat: "mp4",
		FileBytes:  data,
	}
	config := app.InitServiceConfig()
	bucket := config.GetAWSGeneratedAssetsBucket()
	link, err := awsService.UploadFileV2(bucket, awsRequest, true)
	if err != nil {
		logger.Errorw(context.Background(), "error while uploading image into aws", "error", err.Error())
		return
	}

	return link, nil

}

func (s *service) UpdateScene(ctx context.Context, storyID string, sceneID string, selectedImage string) (updatedScene store.Scene, err error) {

	updatedScene, err = s.store.UpdateScene(ctx, storyID, sceneID, selectedImage)
	if err != nil {
		logger.Error(ctx, "error updating scene", err.Error())
		return
	}
	return
}
