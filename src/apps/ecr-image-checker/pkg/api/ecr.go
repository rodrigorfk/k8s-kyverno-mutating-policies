package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

var clientsMu sync.RWMutex
var cacheRepoNotFound *cache.Cache = cache.New(15*time.Minute, 1*time.Minute)
var cacheImages *cache.Cache = cache.New(15*time.Minute, 1*time.Minute)
var ecrClients map[string]*ecr.Client = make(map[string]*ecr.Client)

func (s *Server) ecrHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	var req RuntimeRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		s.ErrorResponse(w, r, "error decoding request body", http.StatusBadRequest)
		return
	}

	if req.RegistryId == "" {
		s.ErrorResponse(w, r, "The registryId query string is empty", http.StatusBadRequest)
		return
	}
	if req.Region == "" {
		s.ErrorResponse(w, r, "The region query string is empty", http.StatusBadRequest)
		return
	}
	if req.RepositoryName == "" {
		s.ErrorResponse(w, r, "The repositoryName query string is empty", http.StatusBadRequest)
		return
	}
	if req.ImageTag == "" {
		s.ErrorResponse(w, r, "The imageTag query string is empty", http.StatusBadRequest)
		return
	}

	data := RuntimeResponse{
		Exists:         false,
		RegistryId:     req.RegistryId,
		Region:         req.Region,
		RepositoryName: req.RepositoryName,
		ImageTag:       req.ImageTag,
	}

	defer s.JSONResponse(w, r, &data)

	// dev mode, skip AWS calls and return a predefined response
	if req.RegistryId == "000000000000" {
		data.Exists = true
		return
	} else if req.RegistryId == "111111111111" {
		data.Exists = false
		return
	}

	var repositoryKey = req.Region + "/" + req.RepositoryName
	_, ok := cacheRepoNotFound.Get(repositoryKey)
	if ok {
		s.logger.Warn("repository was cached as not found", zap.String("repositoryKey", repositoryKey))
		return
	}

	var imageKey = repositoryKey + ":" + req.ImageTag
	_, ok = cacheImages.Get(imageKey)
	if ok {
		s.logger.Debug("image found in the cache", zap.String("imageKey", imageKey))
		data.Exists = true
		return
	}

	clientsMu.RLock()
	client, ok := ecrClients[req.Region]
	clientsMu.RUnlock()
	if !ok {
		clientsMu.Lock()
		client, err = buildEcrClient(req.Region, s)
		if err != nil {
			s.logger.Error("error creating ECR client", zap.Error(err))
			return
		} else {
			ecrClients[req.Region] = client
		}
		clientsMu.Unlock()
	}

	input := &ecr.DescribeImagesInput{
		RepositoryName: &req.RepositoryName,
		RegistryId:     &req.RegistryId,
		ImageIds: []types.ImageIdentifier{
			{
				ImageTag: &req.ImageTag,
			},
		},
	}
	_, err = client.DescribeImages(context.TODO(), input)
	if err != nil {
		var rnfe *types.RepositoryNotFoundException
		if errors.As(err, &rnfe) {
			s.logger.Warn("repository not found", zap.String("repositoryKey", repositoryKey))
			cacheRepoNotFound.Add(repositoryKey, true, 15*time.Minute)
		} else {
			s.logger.Error("error describing ECR image", zap.Error(err))
		}
		return
	}
	cacheImages.Add(imageKey, true, 15*time.Minute)
	data.Exists = true
}

func buildEcrClient(region string, s *Server) (*ecr.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		s.logger.Error("error loading AWS configuration", zap.Error(err))
		return nil, err
	}
	client := ecr.NewFromConfig(cfg)
	return client, nil
}

type RuntimeResponse struct {
	Exists         bool   `json:"exists"`
	RegistryId     string `json:"registryId"`
	Region         string `json:"region"`
	RepositoryName string `json:"repositoryName"`
	ImageTag       string `json:"imageTag"`
}

type RuntimeRequest struct {
	RegistryId     string `json:"registryId"`
	Region         string `json:"region"`
	RepositoryName string `json:"repositoryName"`
	ImageTag       string `json:"imageTag"`
}
