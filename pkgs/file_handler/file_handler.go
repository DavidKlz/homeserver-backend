package filehandler

import (
	"bufio"
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	vidio "github.com/AlexEidt/Vidio"
	"github.com/DavidKlz/homeserver-backend/pkgs/logger"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/liujiawm/graphics-go/graphics"
)

type FileHandler struct {
	File   fs.DirEntry
	Path   string
	FileID string
	Mime   *mimetype.MIME
}

func CreateFileHandler(file fs.DirEntry) (*FileHandler, error) {
	id := uuid.New().String()
	path := filepath.Join(os.Getenv("INPUT_FOLDER"), file.Name())
	mime, err := mimetype.DetectFile(path)
	if err != nil {
		return nil, err
	}
	return &FileHandler{
		File:   file,
		Path:   path,
		FileID: id,
		Mime:   mime,
	}, nil
}

func (fh *FileHandler) CreateThumbnail(dir string) (string, error) {
	if strings.Contains(fh.Mime.String(), "image") {
		return fh.createImageThumbnail(dir)
	} else if strings.Contains(fh.Mime.String(), "video") {
		video, _ := vidio.NewVideo(fh.Path)

		thumbDir := filepath.Join(dir, ".thumb")
		err := os.MkdirAll(thumbDir, os.ModePerm)

		if err != nil {
			logger.Error("The subfolders could not be created: %s", err.Error())
		}

		return createVideoThumbnail(fh.Path, filepath.Join(thumbDir, generateUniqueFilename("THUMB", fh.FileID, ".jpg")), strconv.FormatInt(int64(video.Duration()/2), 10))
	}
	return "", fmt.Errorf("Unsupported file type")
}

func (fh *FileHandler) createImageThumbnail(dstDir string) (string, error) {
	img, err := os.Open(fh.Path)
	if err != nil {
		return "", err
	}
	defer img.Close()

	var imData image.Image

	if fh.Mime.Is("image/gif") {
		imData, err = gif.Decode(bufio.NewReader(img))
		if err != nil {
			return "", err
		}
	} else {
		imData, _, err = image.Decode(img)
		if err != nil {
			return "", err
		}
	}

	height := float64(imData.Bounds().Dy()) * (400.0 / float64(imData.Bounds().Dx()))
	dstImage := image.NewRGBA(image.Rect(0, 0, 400, int(height)))
	graphics.Thumbnail(dstImage, imData)

	thumbDir := filepath.Join(dstDir, ".thumb")
	err = os.MkdirAll(thumbDir, os.ModePerm)

	if err != nil {
		logger.Error("The subfolders could not be created: %s", err.Error())
	}

	thumbFile := filepath.Join(thumbDir, generateUniqueFilename("THUMB", fh.FileID, ".jpg"))
	newImage, err := os.Create(thumbFile)
	if err != nil {
		return "", err
	}
	defer newImage.Close()

	err = png.Encode(newImage, imData)
	return thumbFile, err
}

func createVideoThumbnail(srcDir, dstDir, time string) (string, error) {
	ffCmd := exec.Command("ffmpeg", "-i", srcDir, "-vf", "thumbnail,scale=400:-1", "-frames:v", "1", "-ss", time, dstDir)

	_, err := ffCmd.CombinedOutput()

	return dstDir, err
}

func (fh *FileHandler) MoveFile(dir string) (string, error) {
	dstDir := filepath.Join(dir, generateUniqueFilename("FILE", fh.FileID, fh.Mime.Extension()))
	err := os.Rename(fh.Path, dstDir)

	if err != nil {
		return "", err
	}

	fh.Path = dstDir

	return dstDir, nil
}

func generateUniqueFilename(prefix, id, suffix string) string {
	return fmt.Sprintf("%s_%s_%s%s", prefix, time.Now().Format("20060102150405"), id, suffix)
}
