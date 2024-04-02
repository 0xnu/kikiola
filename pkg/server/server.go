package server

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/0xnu/kikiola/pkg/db"
	"github.com/0xnu/kikiola/pkg/index"
	"github.com/gorilla/mux"
)

type Server struct {
	storage *db.DistributedStorage
	index   *index.Index
	server  *http.Server
}

type Object struct {
	ID       string            `json:"id"`
	Object   []byte            `json:"object"`
	Metadata map[string]string `json:"metadata"`
}

func NewServer(storage *db.DistributedStorage, index *index.Index) *Server {
	return &Server{
		storage: storage,
		index:   index,
	}
}

func (s *Server) Start(addr string) error {
	router := s.Router()
	return http.ListenAndServe(addr, router)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/vectors", s.handleInsertVector).Methods("POST")
	router.HandleFunc("/vectors/{id}", s.handleGetVector).Methods("GET")
	router.HandleFunc("/vectors/{id}", s.handleDeleteVector).Methods("DELETE")
	router.HandleFunc("/vectors/{id}/metadata", s.handleUpdateVectorMetadata).Methods("PATCH")
	router.HandleFunc("/query/{id}", s.handleQueryVector).Methods("GET")
	router.HandleFunc("/search", s.handleSearchVectors).Methods("POST")
	router.HandleFunc("/objects", s.handleInsertObject).Methods("POST")
	router.HandleFunc("/objects/{id}", s.handleGetObject).Methods("GET")
	router.HandleFunc("/objects/{id}", s.handleDeleteObject).Methods("DELETE")
	router.HandleFunc("/objects/{id}/metadata", s.handleUpdateObjectMetadata).Methods("PATCH")
	router.HandleFunc("/objects/{id}/content", s.handleUpdateObjectContent).Methods("PATCH")

	return router
}

func (s *Server) handleInsertVector(w http.ResponseWriter, r *http.Request) {
	var vector db.Vector
	err := json.NewDecoder(r.Body).Decode(&vector)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = s.index.Insert(&vector)
	if err != nil {
		http.Error(w, "Failed to insert vector", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleQueryVector(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	vector, err := s.storage.GetVector(id)
	if err != nil {
		if errors.Is(err, db.ErrVectorNotFound) {
			http.Error(w, "Vector not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve vector", http.StatusInternalServerError)
		}
		return
	}

	response := struct {
		Text string `json:"text"`
	}{
		Text: vector.Text,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleUpdateVectorMetadata(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var updateReq struct {
		Metadata map[string]string `json:"metadata"`
	}
	err := json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = s.storage.UpdateVectorMetadata(id, updateReq.Metadata)
	if err != nil {
		if errors.Is(err, db.ErrVectorNotFound) {
			http.Error(w, "Vector not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update vector metadata", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleGetVector(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	vector, err := s.storage.GetVector(id)
	if err != nil {
		if errors.Is(err, db.ErrVectorNotFound) {
			http.Error(w, "Vector not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve vector", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(vector)
}

func (s *Server) handleDeleteVector(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := s.index.Delete(id)
	if err != nil {
		http.Error(w, "Failed to delete vector", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleSearchVectors(w http.ResponseWriter, r *http.Request) {
	var searchReq SearchRequest
	err := json.NewDecoder(r.Body).Decode(&searchReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	results, err := s.index.Search(searchReq.Vector, searchReq.K)
	if err != nil {
		http.Error(w, "Failed to search vectors", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}

func (s *Server) handleInsertObject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	jsonData := r.FormValue("data")
	var object db.Object

	err = json.Unmarshal([]byte(jsonData), &object)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("object")
	if err != nil {
		http.Error(w, "Failed to retrieve object file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	objectData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read object file", http.StatusBadRequest)
		return
	}

	object.Object = objectData

	err = s.storage.InsertObject(&object)
	if err != nil {
		http.Error(w, "Failed to insert object", http.StatusInternalServerError)
		log.Printf("Error inserting object: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleGetObject(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	object, err := s.storage.GetObject(id)
	if err != nil {
		if errors.Is(err, db.ErrObjectNotFound) {
			http.Error(w, "Object not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve object", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(object)
}

func (s *Server) handleDeleteObject(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := s.storage.DeleteObject(id)
	if err != nil {
		if errors.Is(err, db.ErrObjectNotFound) {
			http.Error(w, "Object not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete object from storage", http.StatusInternalServerError)
			log.Printf("Error deleting object from storage: %v", err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleUpdateObjectMetadata(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var updateReq struct {
		Metadata map[string]string `json:"metadata"`
	}

	err := json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = s.storage.UpdateObjectMetadata(id, updateReq.Metadata)
	if err != nil {
		if errors.Is(err, db.ErrObjectNotFound) {
			http.Error(w, "Object not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update object metadata", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleUpdateObjectContent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("object")
	if err != nil {
		http.Error(w, "Failed to retrieve object file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	objectData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read object file", http.StatusBadRequest)
		return
	}

	object, err := s.storage.GetObject(id)
	if err != nil {
		if errors.Is(err, db.ErrObjectNotFound) {
			http.Error(w, "Object not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve object", http.StatusInternalServerError)
		}
		return
	}

	object.Object = objectData

	err = s.storage.InsertObject(object)
	if err != nil {
		http.Error(w, "Failed to update object content", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type SearchRequest struct {
	Vector *db.Vector `json:"vector"`
	K      int        `json:"k"`
}
