package repository

import (
	cfg "github.com/andreipimenov/golang-training-2021/internal/config"
	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/rs/zerolog/log"
)

type Repository interface {
	Load(string) (model.Price, bool)
	Store(string, model.Price)
}

func Get() (Repository, func()) {
	if cfg.Get().RepositoryType == "Db" {
		return GetDbRepository()
	}
	if cfg.Get().RepositoryType == "Cache" {
		return GetCacheRepository()
	}
	if cfg.Get().RepositoryType == "Mongo" {
		return GetMongoRepository()
	}
	//TODO return err
	panic("unknown repository type")
}

func GetMongoRepository() (Repository, func()) {
	log.Debug().Msg("Getting 'Mongo' repository")
	return NewMongo()
}

func GetCacheRepository() (Repository, func()) {
	log.Debug().Msg("Getting 'Cache' repository")
	return NewCache(), func() {}
}

func GetDbRepository() (Repository, func()) {
	log.Debug().Msg("Getting 'Db' repository")
	db := dbInit()
	dbMigrations(db)
	return NewDB(db), closeDb(db)
}
