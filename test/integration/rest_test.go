package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	config "github.com/wildan3105/converto/configs"
	"github.com/wildan3105/converto/pkg/api"
	"github.com/wildan3105/converto/pkg/api/schema"
	"github.com/wildan3105/converto/pkg/domain"
	"github.com/wildan3105/converto/pkg/infrastructure/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var app *fiber.App
var mongoClient *mongo.Client

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	config.LoadConfig()

	mongoClient, err = mongodb.Connect(config.AppConfig.MongoURI)
	if err != nil {
		panic("Failed to connect to MongoDB: " + err.Error())
	}

	cleanup()

	app = api.Setup()

	code := m.Run()

	log.Println("Cleaning up test environment...")
	cleanup()

	os.Exit(code)
}

func cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mongoClient.Database(config.AppConfig.MongoDbName).Collection(config.AppConfig.MongoDbCollection)
	if _, err := collection.DeleteMany(ctx, bson.M{}); err != nil {
		log.Fatalf("Failed to cleanup conversions collection: %v", err)
	}

	if err := os.RemoveAll(config.AppConfig.BaseDirectory); err != nil {
		log.Fatalf("Failed to remove base directory %s: %v", config.AppConfig.BaseDirectory, err)
	}
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/health", nil)
	resp, _ := app.Test(req, -1)
	defer resp.Body.Close()

	var responseBody map[string]any
	json.NewDecoder(resp.Body).Decode(&responseBody)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "ok", responseBody["message"])
}

func TestConversionEndpoint(t *testing.T) {
	// Start: POST /api/v1/conversions
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	filePath := "../data/file.shapr"
	file, err := os.Open(filePath)
	assert.NoError(t, err)
	defer file.Close()

	formFile, err := writer.CreateFormFile("file", filepath.Base(filePath))
	assert.NoError(t, err)

	_, err = io.Copy(formFile, file)
	assert.NoError(t, err)

	writer.WriteField("target_format", ".stl")
	writer.Close()

	createConversionRequest := httptest.NewRequest("POST", "/api/v1/conversions", buffer)
	createConversionRequest.Header.Set("Content-Type", writer.FormDataContentType())

	rawCreateConversionResponse, _ := app.Test(createConversionRequest, -1)
	assert.Equal(t, http.StatusCreated, rawCreateConversionResponse.StatusCode)

	var createConversionResponse schema.CreateConversionResponse
	err = json.NewDecoder(rawCreateConversionResponse.Body).Decode(&createConversionResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, createConversionResponse.ID)
	assert.Equal(t, string(domain.ConversionPending), string(createConversionResponse.Status))
	assert.Equal(t, "Conversion created successfully", createConversionResponse.Message)
	// End: POST /api/v1/conversions

	// Start: GET /api/v1/conversions
	getConversionsRequest := httptest.NewRequest("GET", "/api/v1/conversions?status=pending&page=1&limit=10", nil)
	rawGetConversationsResponse, _ := app.Test(getConversionsRequest, -1)
	assert.Equal(t, http.StatusOK, rawGetConversationsResponse.StatusCode)

	var listResp schema.ListConversionsResponse
	err = json.NewDecoder(rawGetConversationsResponse.Body).Decode(&listResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, listResp)
	assert.Equal(t, 1, listResp.Page)
	assert.Equal(t, 10, listResp.Limit)
	assert.Equal(t, createConversionResponse.ID, listResp.Data[0].ID)
	assert.Len(t, listResp.Data, 1)
	// End: GET /api/v1/conversions

	// add sleep to simulate the "background job" to upload and convert file
	time.Sleep(5 * time.Second)

	// Start: GET /api/v1/conversions/:id
	getConversationRequest := httptest.NewRequest("GET", "/api/v1/conversions/"+createConversionResponse.ID, nil)
	rawGetConversationResponse, _ := app.Test(getConversationRequest, -1)
	assert.Equal(t, http.StatusOK, rawGetConversationResponse.StatusCode)

	var getConversationResponse schema.ConversionResponse
	err = json.NewDecoder(rawGetConversationResponse.Body).Decode(&getConversationResponse)
	assert.NoError(t, err)
	assert.Equal(t, createConversionResponse.ID, getConversationResponse.ID)
	assert.Equal(t, string(domain.ConversionCompleted), string(getConversationResponse.Status))
	assert.Equal(t, 100, getConversationResponse.Progress)
	assert.Contains(t, getConversationResponse.OriginalFilePath, "/original", "OriginalFilePath should contain a directory path")
	assert.Contains(t, getConversationResponse.ConvertedFilePath, "/converted", "ConvertedFilePath should contain a directory path")
	// End: GET /api/v1/conversions/:id

	// Start: GET /api/v1/conversions/:id/files
	getFileByConversionIdRequest := httptest.NewRequest("GET", "/api/v1/conversions/"+createConversionResponse.ID+"/files?type=original", nil)
	respGetFileByConversionId, _ := app.Test(getFileByConversionIdRequest, -1)

	assert.Equal(t, http.StatusOK, respGetFileByConversionId.StatusCode)

	contentDisposition := respGetFileByConversionId.Header.Get("Content-Disposition")
	assert.Contains(t, contentDisposition, "attachment; filename=")

	contentType := respGetFileByConversionId.Header.Get("Content-Type")
	assert.Equal(t, "application/octet-stream", contentType)

	body, err := io.ReadAll(respGetFileByConversionId.Body)
	assert.NoError(t, err)
	assert.NotEmpty(t, body)
	// End: GET /api/v1/conversions/:id/files
}
