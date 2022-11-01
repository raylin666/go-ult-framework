package data

import (
	"ult/pkg/logger"
	"ult/pkg/repositories"
)

type DataRepo struct {}

func NewDataRepo(logger *logger.Logger, repo repositories.DataRepo) *DataRepo {
	return &DataRepo{}
}