package server

import (
	"encoding/json"
	"net/http"

	"github.com/DavidKlz/homeserver-backend/types"
)

func (s *Server) HandleAuthRequest(w http.ResponseWriter, r *http.Request) error {
	user := &types.User{}
	json.NewDecoder(r.Body).Decode(&user)

	token, err := s.Storage.FindUser(user.Username, user.Password)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, &types.DefaultResponse{Success: true, Message: token})
}
