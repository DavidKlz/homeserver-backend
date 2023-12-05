package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DavidKlz/homeserver-backend/types"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) HandleMetaInfoRequest(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetMetaInfo(w, r)
	case "POST":
		return s.handleSaveMetaInfo(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *Server) handleGetMetaInfo(w http.ResponseWriter, r *http.Request) error {
	mi, err := s.Storage.GetAllMetaInfo()
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, mi)
}

func (s *Server) handleSaveMetaInfo(w http.ResponseWriter, r *http.Request) error {
	mi := types.MetaInfo{}
	err := json.NewDecoder(r.Body).Decode(&mi)
	if err != nil {
		return err
	}

	err = s.Storage.SaveMetaInfo(&mi)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusCreated, mi)
}

func (s *Server) HandleMetaInfoWithGetParamsRequest(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetMetaInfoById(w, r)
	case "DELETE":
		return s.handleDeleteMetaInfoById(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *Server) handleGetMetaInfoById(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}

	mi, err := s.Storage.GetMetaInfo(&id)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, mi)
}

func (s *Server) handleDeleteMetaInfoById(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}

	err = s.Storage.DeleteMetaInfo(&id)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, types.DefaultResponse{Success: true, Message: "meta info removed"})
}

func (s *Server) HandleMediaToMetaInfoRequest(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) HandleMetaInfoToMediaRequest(w http.ResponseWriter, r *http.Request) error {
	return nil
}
