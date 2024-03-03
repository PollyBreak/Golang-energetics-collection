package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/Golang-energetics-collection/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func TestGetEnergetics(t *testing.T) {

	setupTestDatabase()

	request, _ := http.NewRequest("GET", "/energetix", nil)
	response := httptest.NewRecorder()
	getEnergetics(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Incorrect status code. Expected: %d, Got: %d", http.StatusOK, response.Code)
	}

}

func setupTestDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	err = models.DB.AutoMigrate()
	if err != nil {
		log.Fatalf("Error migrating database schema: %v", err)
	}
}

func TestGetEnergeticsById(t *testing.T) {
	req, err := http.NewRequest("GET", "/energetix/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getEnergeticsById)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestPostEnergetic(t *testing.T) {
	newEnergetic := Energetic{
		Name:               "Test Energetic",
		Taste:              "Test Taste",
		Description:        "Test Description",
		ManufacturerName:   "Test Manufacturer",
		ManufactureCountry: "Test Country",
	}

	jsonData, err := json.Marshal(newEnergetic)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/energetix", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postEnergetic)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestUpdateEnergeticsById(t *testing.T) {
	updatedEnergetic := Energetic{
		Name:               "Updated Test Energetic",
		Taste:              "Updated Test Taste",
		Description:        "Updated Test Description",
		ManufacturerName:   "Updated Test Manufacturer",
		ManufactureCountry: "Updated Test Country",
	}

	jsonData, err := json.Marshal(updatedEnergetic)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/energetix/1", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateEnergeticsById)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteEnergeticById(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/energetix/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteEnergeticById)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetNumberOfPages(t *testing.T) {
	req, err := http.NewRequest("GET", "/pages", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getNumberOfPages)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

// TestEndToEndScenario tests the entire application flow from client request to server response.
func TestEndToEndScenario(t *testing.T) {
	// Create a request to simulate a client request
	requestBody := map[string]interface{}{
		"name":               "Test Energetic",
		"taste":              "Test Taste",
		"description":        "Test Description",
		"manufacturerName":   "Test Manufacturer",
		"manufactureCountry": "Test Country",
	}
	requestJSON, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/energetix", bytes.NewReader(requestJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call your HTTP handler directly to handle the request
	handler := http.HandlerFunc(postEnergetic)
	handler.ServeHTTP(rr, req)

	// Check the status code of the response
	assert.Equal(t, http.StatusCreated, rr.Code, "status code is not Created")

	// Check the response body to ensure it contains the expected data
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Perform assertions on the response body
	assert.Equal(t, "Test Energetic", response["name"], "name does not match expected")
	assert.Equal(t, "Test Taste", response["taste"], "taste does not match expected")
}
