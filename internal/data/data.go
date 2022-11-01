package data

import (
	"ult/pkg/global"
	"ult/pkg/logger"
)

type DataRepo struct {}

func NewDataRepo(logger *logger.Logger, repo global.DataRepo) *DataRepo {
	return &DataRepo{}
}