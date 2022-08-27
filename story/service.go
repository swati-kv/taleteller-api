package story

import (
	"context"
	"github.com/gojektech/heimdall/v6/httpclient"
	"github.com/gojektech/heimdall/v6/plugins"
	"taleteller/api"
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
}

type service struct {
	httpClient     *httpclient.Client
	store          store.StoryStorer
	generatorUtils utils.IDGeneratorUtils
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

	go func() {
		s.generateImage(ctx, createSceneRequest, sceneID)
		wg.Done()
	}()

	fmt.Println("generating ...")
	wg.Wait()
	fmt.Println("returning--")
	return
}

func (s *service) generateImage(ctx context.Context, createSceneRequest CreateSceneRequest, sceneID string) (response PyImageResponse, err error) {
	wg.Add(1)
	ctx = context.Background()
	//defer s.processGenerateImage(response, sceneID, err)

	url := fmt.Sprintf("%s/%s", s.pyServerBaseURL, createSceneEndPoint)
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	//header["request-id"] = sceneID
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

	for _, image := range response.Data.GeneratedImage {
		fmt.Println("inside for")
		imageID, err := utils.NewGeneratorUtils().GenerateIDWithPrefix("image_")
		if err != nil {
			logger.Errorw(context.Background(), "error while generating image id", "error", err.Error())
			return
		}

		awsService := utils.NewAWSService()
		awsRequest := utils.UploadS3{
			File:       image,
			FileType:   "image",
			FileFormat: response.Data.GeneratedImageFormat,
		}
		config := app.InitServiceConfig()
		bucket := config.GetAWSGeneratedAssetsBucket()
		link, err := awsService.UploadFile(bucket, awsRequest, true)
		if err != nil {
			logger.Errorw(context.Background(), "error while uploading image into aws", "error", err.Error())
			return
		}
		fmt.Println("link - ", link)

		insertImageStoreRequest := store.InsertImage{
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
	return
}

func (s *service) createSceneEntry(id string, request CreateSceneRequest, storyID string) (err error) {
	storeRequest := store.CreateSceneRequest{
		SceneID:     id,
		Status:      statusStarted,
		StoryID:     storyID,
		SceneNumber: request.SceneNumber,
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

func (s *service) List(ctx context.Context, status string) (stories []store.Story, err error) {
	stories, err = s.store.List(ctx, status)
	if err != nil {
		logger.Error(ctx, "error getting stories", err.Error())
		return
	}
	return
}

func (s *service) Publish(ctx context.Context, req []UpdateSceneOrderReq, storyID string) (path string, err error) {
	timeout := 30000 * time.Millisecond
	client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

	requestLogger := plugins.NewRequestLogger(nil, nil)
	client.AddPlugin(requestLogger)

	for _, scene := range req {
		err = s.store.UpdateScene(ctx, scene.SceneID, scene.SceneNumber, storyID)
		if err != nil {
			logger.Error(ctx, "error updating scene", err.Error())
			return
		}
	}

	httpResponse, err := api.Post(ctx, url, requestBody, headers, client)
	if err != nil {
		logger.Errorw(ctx, "error in publishing", "error", err.Error(), "url", url, "request_body", requestBody, "headers", headers)
		return
	}
	return "test", nil
}
