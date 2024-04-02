package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/0xnu/kikiola/pkg/db"
	"github.com/0xnu/kikiola/pkg/index"
	"github.com/0xnu/kikiola/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestDistributedVectorDatabase(t *testing.T) {
	nodeAddresses := generateNodeAddresses("localhost", 3401, 3420)

	storage, err := db.NewDistributedStorage(nodeAddresses)
	assert.NoError(t, err)
	defer storage.Close()

	index, err := index.NewIndex(storage)
	assert.NoError(t, err)

	server := server.NewServer(storage, index)
	ts := httptest.NewServer(server.Router())
	defer ts.Close()

	entries := []db.Vector{
		{ID: "vector1", Embedding: []float64{0.1, 0.2, 0.3}, Metadata: map[string]string{"name": "Vector 1", "category": "sample"}, Text: "This is the text content for vector1."},
		{ID: "vector2", Embedding: []float64{0.4, 0.5, 0.6}, Metadata: map[string]string{"name": "Vector 2", "category": "sample"}, Text: "This is the text content for vector2."},
		{ID: "vector3", Embedding: []float64{0.7, 0.8, 0.9}, Metadata: map[string]string{"name": "Vector 3", "category": "sample"}, Text: "This is the text content for vector3."},
	}

	for _, entry := range entries {
		jsonData, _ := json.Marshal(entry)
		resp, err := http.Post(ts.URL+"/vectors", "application/json", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	}

	resp, err := http.Get(ts.URL + "/vectors/vector2")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var retrievedEntry db.Vector
	err = json.NewDecoder(resp.Body).Decode(&retrievedEntry)
	assert.NoError(t, err)
	assert.Equal(t, entries[1], retrievedEntry)

	newEntry := db.Vector{ID: "vector4", Embedding: []float64{0.2, 0.4, 0.6}, Metadata: map[string]string{"name": "Vector 4", "category": "sample"}, Text: "This is the text content for vector4."}
	jsonData, _ := json.Marshal(newEntry)
	resp, err = http.Post(ts.URL+"/vectors", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/vectors/vector1", nil)
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	searchReq := struct {
		Vector *db.Vector `json:"vector"`
		K      int        `json:"k"`
	}{
		Vector: &db.Vector{ID: "query_vector", Embedding: []float64{0.5, 0.6, 0.7}},
		K:      2,
	}
	jsonData, _ = json.Marshal(searchReq)
	resp, err = http.Post(ts.URL+"/search", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var searchResults []*db.Vector
	err = json.NewDecoder(resp.Body).Decode(&searchResults)
	assert.NoError(t, err)
	assert.Len(t, searchResults, 2)

	resp, err = http.Get(ts.URL + "/query/vector2")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var queryResult struct {
		Text string `json:"text"`
	}
	err = json.NewDecoder(resp.Body).Decode(&queryResult)
	assert.NoError(t, err)
	assert.Equal(t, "This is the text content for vector2.", queryResult.Text)

	updateReq := struct {
		Metadata map[string]string `json:"metadata"`
	}{
		Metadata: map[string]string{
			"name":     "Updated Vector 2",
			"category": "updated",
		},
	}
	jsonData, _ = json.Marshal(updateReq)
	req, err = http.NewRequest(http.MethodPatch, ts.URL+"/vectors/vector2/metadata", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Get(ts.URL + "/vectors/vector2")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedVector db.Vector
	err = json.NewDecoder(resp.Body).Decode(&updatedVector)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Vector 2", updatedVector.Metadata["name"])
	assert.Equal(t, "updated", updatedVector.Metadata["category"])

	jsonData, _ = json.Marshal(updateReq)
	req, err = http.NewRequest(http.MethodPatch, ts.URL+"/vectors/nonexistent_vector/metadata", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	file, err := os.Open("oxford.jpg")
	assert.NoError(t, err)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("object", "oxford.jpg")
	assert.NoError(t, err)
	_, err = io.Copy(part, file)
	assert.NoError(t, err)

	metadata := map[string]string{
		"id":       "0539f0ac-6771-47c6-8f5e-2cdf272a6de0",
		"name":     "Oxford",
		"category": "Images",
	}
	metadataJSON, _ := json.Marshal(metadata)
	_ = writer.WriteField("data", string(metadataJSON))
	err = writer.Close()
	assert.NoError(t, err)

	resp, err = http.Post(ts.URL+"/objects", writer.FormDataContentType(), body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = http.Get(ts.URL + "/objects/0539f0ac-6771-47c6-8f5e-2cdf272a6de0")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// req, err = http.NewRequest(http.MethodDelete, ts.URL+"/objects/0539f0ac-6771-47c6-8f5e-2cdf272a6de0", nil)
	// assert.NoError(t, err)
	// resp, err = http.DefaultClient.Do(req)
	// assert.NoError(t, err)
	// assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	updateMetadata := map[string]string{
		"name":     "Oxford High Street",
		"category": "Image",
	}
	updateMetadataJSON, _ := json.Marshal(updateMetadata)
	req, err = http.NewRequest(http.MethodPatch, ts.URL+"/objects/0539f0ac-6771-47c6-8f5e-2cdf272a6de0/metadata", bytes.NewBuffer(updateMetadataJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	file, err = os.Open("oxford_high_street.webp")
	assert.NoError(t, err)
	defer file.Close()

	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	part, err = writer.CreateFormFile("object", "oxford_high_street.webp")
	assert.NoError(t, err)
	_, err = io.Copy(part, file)
	assert.NoError(t, err)
	err = writer.Close()
	assert.NoError(t, err)

	req, err = http.NewRequest(http.MethodPatch, ts.URL+"/objects/0539f0ac-6771-47c6-8f5e-2cdf272a6de0/content", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
