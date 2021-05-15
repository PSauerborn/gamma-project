package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func UserHeaderMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// extract user ID from request header and parse
		userid := ctx.Request.Header.Get("X-Authenticated-Userid")
		if len(userid) == 0 || userid == "undefined" {
			log.Warn(("cannot extract user ID from header"))
			status := http.StatusForbidden
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Forbidden"})
			return
		}
		ctx.Set("uid", userid)
		ctx.Next()
	}
}
