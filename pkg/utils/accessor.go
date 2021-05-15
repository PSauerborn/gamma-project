package utils

import "github.com/PSauerborn/gamma-project/internal/pkg/utils"

func NewBaseAccessor(host, protocol string, port int) *utils.BaseAPIAccessor {
	return &utils.BaseAPIAccessor{
		Host:     host,
		Port:     &port,
		Protocol: protocol,
	}
}
