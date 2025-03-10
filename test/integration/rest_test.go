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

	collection := mongoClient.Database(config.AppConfig.MongoDbName).Collection("conversions")
	if _, err := collection.DeleteMany(ctx, bson.M{}); err != nil {
		log.Fatalf("Failed to cleanup conversions collection: %v", err)
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
	// POST /api/v1/conversions - Happy Path
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

	req := httptest.NewRequest("POST", "/api/v1/conversions", buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createResp schema.CreateConversionResponse
	err = json.NewDecoder(resp.Body).Decode(&createResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, createResp.ID)
	assert.Equal(t, "pending", string(createResp.Status))

	// GET /api/v1/conversions - Happy Path
	req = httptest.NewRequest("GET", "/api/v1/conversions?status=pending&page=1&limit=10", nil)
	resp, _ = app.Test(req, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var listResp schema.ListConversionsResponse
	err = json.NewDecoder(resp.Body).Decode(&listResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, listResp)
	assert.Equal(t, createResp.ID, listResp.Data[0].ID)

	// GET /api/v1/conversions/:id - Happy Path
	req = httptest.NewRequest("GET", "/api/v1/conversions/"+createResp.ID, nil)
	resp, _ = app.Test(req, -1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var getResp schema.ConversionResponse
	err = json.NewDecoder(resp.Body).Decode(&getResp)
	assert.NoError(t, err)
	assert.Equal(t, createResp.ID, getResp.ID)

	// Negative Scenarios
	invalidReq := httptest.NewRequest("POST", "/api/v1/conversions", nil)
	resp, _ = app.Test(invalidReq, -1)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Invalid Status Query
	req = httptest.NewRequest("GET", "/api/v1/conversions?status=invalid&page=0&limit=50", nil)
	resp, _ = app.Test(req, -1)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Invalid Object ID
	req = httptest.NewRequest("GET", "/api/v1/conversions/invalid_id", nil)
	resp, _ = app.Test(req, -1)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Not Found Scenario
	req = httptest.NewRequest("GET", "/api/v1/conversions/60b8d6f5f9c4b8b8b8b8b8b8", nil)
	resp, _ = app.Test(req, -1)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// GET /api/v1/conversions/:id/files - 	does not provided any query, return 400
	invalidGetFileByConversionId := httptest.NewRequest("GET", "/api/v1/conversions/"+createResp.ID+"/files", nil)
	resp, _ = app.Test(invalidGetFileByConversionId, -1)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Happy path - return the expected response
	validGetFileByConversionId := httptest.NewRequest("GET", "/api/v1/conversions/"+createResp.ID+"/files?type=original", nil)
	respGetFileByConversionId, _ := app.Test(validGetFileByConversionId, -1)

	assert.Equal(t, http.StatusOK, respGetFileByConversionId.StatusCode)

	contentDisposition := respGetFileByConversionId.Header.Get("Content-Disposition")
	assert.Contains(t, contentDisposition, "attachment; filename=")

	contentType := respGetFileByConversionId.Header.Get("Content-Type")
	assert.Equal(t, "application/octet-stream", contentType)

	body, err := io.ReadAll(respGetFileByConversionId.Body)
	assert.NoError(t, err)
	assert.NotEmpty(t, body)
}
