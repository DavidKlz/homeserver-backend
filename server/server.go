package server

import (
	"encoding/json"
	"net/http"

	"github.com/DavidKlz/homeserver-backend/storage"
	"github.com/DavidKlz/homeserver-backend/types"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gorilla/mux"
)

type Server struct {
	Port    string
	Router  *mux.Router
	Storage storage.Storage
}

func CreateNewServer(port string, store storage.Storage) *Server {
	router := mux.NewRouter()
	s := &Server{
		Port:    port,
		Router:  router,
		Storage: store,
	}
	setupRoutes(router, s)
	return s
}

func (s *Server) Run() {
	http.ListenAndServe(s.Port, s.Router)
}

func setupRoutes(r *mux.Router, s *Server) {
	r.HandleFunc("/metainfo", s.makeHttpHandlerFunc(s.HandleMetaInfoRequest, true))
	r.HandleFunc("/metainfo/{id}", s.makeHttpHandlerFunc(s.HandleMetaInfoWithGetParamsRequest, true))
	r.HandleFunc("/media", s.makeHttpHandlerFunc(s.HandleMediaRequest, true))
	r.HandleFunc("/media/{id}", s.makeHttpHandlerFunc(s.HandleMediaWithGetParamsRequest, true))

	r.HandleFunc("/search", s.makeHttpHandlerFunc(s.HandleSearchMediaRequest, true))

	r.HandleFunc("/media/metainfo", s.makeHttpHandlerFunc(s.HandleMediaToMetaInfoRequest, true))
	r.HandleFunc("/metainfo/media", s.makeHttpHandlerFunc(s.HandleMetaInfoToMediaRequest, true))

	r.HandleFunc("/authenticate", s.makeHttpHandlerFunc(s.HandleAuthRequest, false))
	r.HandleFunc("/synchronize", s.makeHttpHandlerFunc(s.HandleSynchronizationRequest, true))
	r.HandleFunc("/upload", s.makeHttpHandlerFunc(s.HandleFileUploadRequest, true))

	r.HandleFunc("/file/{id}", s.makeHttpHandlerFunc(s.HandleServeFileRequest, false))
	r.HandleFunc("/thumb/{id}", s.makeHttpHandlerFunc(s.HandleServeThumbnailRequest, true))
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func ServeFile(w http.ResponseWriter, r *http.Request, status int, path string) error {
	mimetype, err := mimetype.DetectFile(path)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", mimetype.String())
	w.WriteHeader(status)
	http.ServeFile(w, r, path)
	return nil
}

type serverFunc func(http.ResponseWriter, *http.Request) error

func (s *Server) makeHttpHandlerFunc(f serverFunc, jwtSecured bool) http.HandlerFunc {
	if jwtSecured {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("jwt-token")
			err := s.validateJWT(tokenString)
			if err != nil {
				WriteJson(w, http.StatusUnauthorized, types.DefaultResponse{Success: false, Message: err.Error()})
			} else {
				if err := f(w, r); err != nil {
					WriteJson(w, http.StatusBadRequest, types.DefaultResponse{Success: false, Message: err.Error()})
				}
			}
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, types.DefaultResponse{Success: false, Message: err.Error()})
		}
	}
}

func (s *Server) validateJWT(tokenString string) error {
	_, err := s.Storage.FindUserByToken(tokenString)
	return err
}
