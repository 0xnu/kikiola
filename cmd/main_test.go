package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xnu/kikiola/pkg/db"
	"github.com/0xnu/kikiola/pkg/index"
	"github.com/0xnu/kikiola/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestVectorDatabase(t *testing.T) {
	storage, err := db.NewStorage("data/test_vectors.db")
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
}
