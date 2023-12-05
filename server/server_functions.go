package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	filehandler "github.com/DavidKlz/homeserver-backend/pkgs/file_handler"
	"github.com/DavidKlz/homeserver-backend/pkgs/logger"
	"github.com/DavidKlz/homeserver-backend/types"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) HandleSynchronizationRequest(w http.ResponseWriter, r *http.Request) error {
	logger.Info("The synchronization process is started")
	files, err := os.ReadDir(os.Getenv("INPUT_FOLDER"))
	if err != nil {
		return err
	}

	for _, file := range files {
		media := &types.Media{
			Favorite: false,
		}
		fileHandler, err := filehandler.CreateFileHandler(file)
		if err != nil {
			return err
		}

		mType, err := getFileType(fileHandler.Mime)
		if err != nil {
			return err
		}
		media.Type = mType

		dstDir := filepath.Join(os.Getenv("OUTPUT_FOLDER"), mType)
		err = os.MkdirAll(dstDir, os.ModePerm)

		if err != nil {
			logger.Error("The subfolders could not be created: %s", err.Error())
		}

		if fileHandler.Mime.Is("image/webp") || fileHandler.Mime.Is("image/heic") || fileHandler.Mime.Is("image/heic-sequence") {
			pathToFile, err := fileHandler.MoveFile(dstDir)
			if err != nil {
				return err
			}
			media.PathToFile = pathToFile
			media.PathToThumbnail = pathToFile
			media.Name = pathToFile[strings.LastIndex(pathToFile, "\\")+1:]
		} else {
			pathToThumb, err := fileHandler.CreateThumbnail(dstDir)
			if err != nil {
				return err
			}
			media.PathToThumbnail = pathToThumb

			pathToFile, err := fileHandler.MoveFile(dstDir)
			if err != nil {
				return err
			}
			media.PathToFile = pathToFile
			media.Name = pathToFile[strings.LastIndex(pathToFile, "\\")+1:]
		}

		if err = s.Storage.SaveMedia(media); err != nil {
			return err
		}
	}

	logger.Success("The synchronization was successfully completed!")
	return WriteJson(w, http.StatusOK, types.DefaultResponse{Success: true, Message: "Media files synched"})
}

func getFileType(mime *mimetype.MIME) (string, error) {
	if strings.Contains(mime.String(), "gif") {
		return types.ANIMATION, nil
	} else if strings.Contains(mime.String(), "image") {
		return types.IMAGE, nil
	} else if strings.Contains(mime.String(), "video") {
		return types.VIDEO, nil
	}
	return "", fmt.Errorf("Unsupported File type")
}

func (s *Server) HandleFileUploadRequest(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) HandleServeFileRequest(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	media, err := s.Storage.GetMedia(&id)
	if err != nil {
		return err
	}

	return ServeFile(w, r, http.StatusOK, media.PathToFile)
}

func (s *Server) HandleServeThumbnailRequest(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	media, err := s.Storage.GetMedia(&id)
	if err != nil {
		return err
	}

	return ServeFile(w, r, http.StatusOK, media.PathToThumbnail)
}
