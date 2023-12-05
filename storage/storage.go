package storage

import (
	"fmt"
	"os"

	"github.com/DavidKlz/homeserver-backend/pkgs/logger"
	"github.com/DavidKlz/homeserver-backend/types"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gl "gorm.io/gorm/logger"
)

type Storage interface {
	SaveUser(*types.User) error
	SaveMedia(*types.Media) error
	SaveMetaInfo(*types.MetaInfo) error
	SaveMediaType(*types.MediaType) error
	SaveMediaToMetaInfo(*types.MediaToMetaInfo) error
	DeleteUser(*uuid.UUID) error
	DeleteMedia(*uuid.UUID) error
	DeleteMetaInfo(*uuid.UUID) error
	DeleteMediaType(string) error
	DeleteMediaToMetaInfo(*uuid.UUID) error
	GetUser(*uuid.UUID) (*types.User, error)
	GetMedia(*uuid.UUID) (*types.Media, error)
	GetMetaInfo(*uuid.UUID) (*types.MetaInfo, error)
	GetMediaType(string) (*types.MediaType, error)
	GetMediaToMetaInfo(*uuid.UUID) (*types.MediaToMetaInfo, error)
	GetAllUser() ([]types.User, error)
	GetAllMedia() ([]types.Media, error)
	GetAllMetaInfo() ([]types.MetaInfo, error)
	GetAllMediaTypes() ([]types.MediaType, error)
	GetAllMediaToMetaInfos() ([]types.MediaToMetaInfo, error)
	GetMetaInfosOfMedia(*uuid.UUID) ([]types.MetaInfo, error)
	GetMediaOfMetaInfo(*uuid.UUID) ([]types.Media, error)
	FindUserByToken(string) (*types.User, error)
	FindUser(string, string) (string, error)
}

type PostgresStore struct {
	DB *gorm.DB
}

func CreatePostgresStore() (*PostgresStore, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))
	logger.Info("Connecting to database")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gl.Default.LogMode(gl.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database:\n %s", err.Error())
	}
	logger.Info("Connection to database established")

	logger.Info("Running migrations")
	db.AutoMigrate(&types.User{})
	db.AutoMigrate(&types.MetaInfo{})
	db.AutoMigrate(&types.MediaType{})
	db.AutoMigrate(&types.Media{})
	db.AutoMigrate(&types.MediaToMetaInfo{})
	logger.Info("Migrations finished")

	logger.Info("Inserting default records")
	user, err := types.CreateUser(os.Getenv("DEFAULT_USERNAME"), os.Getenv("DEFAULT_PASSWORD"))
	if err != nil {
		return nil, fmt.Errorf("Couldn't updated or insert User: %s", err.Error())
	}
	if err := db.Where(user).
		FirstOrCreate(&user).Error; err != nil {
		return nil, fmt.Errorf("Couldn't updated or insert User: %s", err.Error())
	}
	if err := db.Save(types.MediaType{Name: types.IMAGE}).Error; err != nil {
		return nil, fmt.Errorf("Couldn't updated or insert type: %s", err.Error())
	}
	if err := db.Save(types.MediaType{Name: types.VIDEO}).Error; err != nil {
		return nil, fmt.Errorf("Couldn't updated or insert type: %s", err.Error())
	}
	if err := db.Save(types.MediaType{Name: types.ANIMATION}).Error; err != nil {
		return nil, fmt.Errorf("Couldn't updated or insert type: %s", err.Error())
	}
	logger.Info("Inserted default records")

	return &PostgresStore{
		DB: db,
	}, nil
}
