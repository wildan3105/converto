package test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	config "github.com/wildan3105/converto/configs"
	"github.com/wildan3105/converto/pkg/api"
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
