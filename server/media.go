package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DavidKlz/homeserver-backend/types"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) HandleMediaRequest(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetMedia(w, r)
	case "POST":
		return s.handleSaveMedia(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *Server) handleGetMedia(w http.ResponseWriter, r *http.Request) error {
	media, err := s.Storage.GetAllMedia()
	if err != nil {
		return err
	}
	mediaRes := []types.MediaReturn{}

	for _, elem := range media {
		mediaRes = append(mediaRes, types.MediaReturn{
			ID:       elem.ID,
			Name:     elem.Name,
			Favorite: elem.Favorite,
			Type:     elem.Type,
		})
	}

	return WriteJson(w, http.StatusOK, mediaRes)
}

func (s *Server) handleSaveMedia(w http.ResponseWriter, r *http.Request) error {
	media := types.Media{}
	err := json.NewDecoder(r.Body).Decode(&media)
	if err != nil {
		return err
	}

	err = s.Storage.SaveMedia(&media)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusCreated, media)
}

func (s *Server) HandleMediaWithGetParamsRequest(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetMediaById(w, r)
	case "DELETE":
		return s.handleDeleteMediaById(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *Server) handleGetMediaById(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}

	media, err := s.Storage.GetMedia(&id)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, &types.MediaReturn{
		ID:       media.ID,
		Name:     media.Name,
		Favorite: media.Favorite,
		Type:     media.Type,
	})
}

func (s *Server) handleDeleteMediaById(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		return err
	}

	err = s.Storage.DeleteMedia(&id)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, types.DefaultResponse{Success: true, Message: "media file removed"})
}

func (s *Server) HandleSearchMediaRequest(w http.ResponseWriter, r *http.Request) error {
	return nil
}
