package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PSauerborn/gamma-project/internal/pkg/roles"
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

func RoleMiddelware(required roles.Role, host string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get user id from request headers
		userid := ctx.Request.Header.Get("X-Authenticated-Userid")
		if len(userid) == 0 || userid == "undefined" {
			log.Warn(("cannot extract user ID from header"))
			status := http.StatusForbidden
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Forbidden"})
			return
		}
		// generate new http request
		url := fmt.Sprintf("%s/roles/%s", host, userid)
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Error(fmt.Errorf("unable to generate new HTTP request: %+v", err))
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
			return
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-Authenticated-Userid", "roles-lookup")
		// generate http client and get roles from roles API
		client := &http.Client{}

		resp, err := client.Do(request)
		if err != nil {
			log.Error(fmt.Errorf("unable to execute HTTP request: %+v", err))
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
			return
		}
		defer resp.Body.Close()

		// handle response from roles API
		switch resp.StatusCode {
		case 200:
			var r struct {
				HTTPCode int    `json:"http_code"`
				Role     string `json:"role"`
			}
			// decode API JSON response to struct
			if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
				log.Error(fmt.Errorf("unable to parse API response: %+v", err))
				status := http.StatusInternalServerError
				ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
					"message": "Internal server error"})
				return
			}

			// convert role returned from API into Role instance
			role, err := roles.StringToRole(r.Role)
			if err != nil {
				log.Errorf(fmt.Sprintf("received invalid role %s from API", role))
				status := http.StatusInternalServerError
				ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
					"message": "Internal server error"})
				return
			}

			// if user role is large or equal to required role,
			// executer request. else return 403
			if role >= required {
				ctx.Next()
			} else {
				log.Warn(fmt.Errorf("user %s does not have required roles to access route", userid))
				status := http.StatusForbidden
				ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
					"message": "Forbidden"})
				return
			}
		default:
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("received invalid API response: %+v", err)
			} else {
				log.Error(fmt.Errorf("unable to retrieve user roles: received response %s", string(body)))
			}
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
			return
		}
	}
}
