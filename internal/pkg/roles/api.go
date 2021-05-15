package roles

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var persistence Persistence

func SetPersistence(p Persistence) {
	persistence = p
}

// API handler used to serve health check routes
func HealthCheckHandler(ctx *gin.Context) {
	log.Info("received request for health check route")
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"message": "Service running"})
}

func GetUserRolesHandler(ctx *gin.Context) {
	log.Info("received request to retrieve user roles")
	uid := ctx.Param("uid")
	role, err := persistence.GetUserRole(uid)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve roles"))
		status := http.StatusInternalServerError
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"role": role.String()})
}

func SetUserRolesHandler(ctx *gin.Context) {
	log.Info("received request to retrieve set roles")
	uid := ctx.MustGet("uid").(string)

	role, err := persistence.GetUserRole(uid)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve user role: %+v", err))
		status := http.StatusInternalServerError
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Internal server error"})
		return
	}

	// only allow admin users to set roles in database
	if role < Admin {
		log.Warn(fmt.Sprintf("received request to set roles without permissions from user %s", uid))
		status := http.StatusForbidden
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Forbidden"})
		return
	}

	var r struct {
		Uid      string `json:"uid" binding:"required"`
		UserRole Role   `json:"role" binding:"required"`
	}
	if err := ctx.ShouldBind(&r); err != nil {
		log.Error(fmt.Errorf("unable to parse request body: %+v", err))
		status := http.StatusBadRequest
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Invalid request body"})
		return
	}
	// check that role is valid else return 400
	if !r.UserRole.IsValid() {
		log.Error("cannot set roles for user: received invalid role")
		status := http.StatusBadRequest
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Invalid role"})
		return
	}

	if err := persistence.SetUserRole(r.Uid, r.UserRole); err != nil {
		log.Error(fmt.Errorf("unable to set user role: %+v", err))
		status := http.StatusInternalServerError
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"message": "Successfully set user role"})
}
