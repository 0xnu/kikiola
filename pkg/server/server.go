package server

import (
	"encoding/json"
	"net/http"

	"github.com/0xnu/kikiola/pkg/db"
	"github.com/0xnu/kikiola/pkg/index"
	"github.com/gorilla/mux"
)

type Server struct {
	storage *db.Storage
	index   *index.Index
}

func NewServer(storage *db.Storage, index *index.Index) *Server {
	return &Server{
		storage: storage,
		index:   index,
	}
}

func (s *Server) Start(addr string) error {
	router := s.Router()
	return http.ListenAndServe(addr, router)
}

func (s *Server) Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/vectors", s.handleInsertVector).Methods("POST")
	router.HandleFunc("/vectors/{id}", s.handleGetVector).Methods("GET")
	router.HandleFunc("/vectors/{id}", s.handleDeleteVector).Methods("DELETE")
	router.HandleFunc("/search", s.handleSearchVectors).Methods("POST")

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

func (s *Server) handleGetVector(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	vector, err := s.storage.Get(id)
	if err != nil {
		http.Error(w, "Vector not found", http.StatusNotFound)
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

type SearchRequest struct {
	Vector *db.Vector `json:"vector"`
	K      int        `json:"k"`
}
