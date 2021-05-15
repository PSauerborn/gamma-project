package utils

import "github.com/PSauerborn/gamma-project/internal/pkg/utils"

func NewBasePersistence(url string) *utils.BasePostgresPersistence {
	return &utils.BasePostgresPersistence{
		DatabaseURL: url,
	}
}
