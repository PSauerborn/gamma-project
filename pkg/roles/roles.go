package roles

import (
	"github.com/PSauerborn/gamma-project/internal/pkg/roles"
	db "github.com/PSauerborn/gamma-project/internal/pkg/roles/persistence"
	"github.com/PSauerborn/gamma-project/pkg/utils"
	"github.com/gin-gonic/gin"
)

// function used to generate new instance of roles API
func NewRolesAPI(p roles.Persistence) *gin.Engine {
	// set persistence as global within module
	roles.SetPersistence(p)
	// generate new instance of gin router and assign routes
	r := gin.Default()
	r.Use(utils.UserHeaderMiddleware())

	r.GET("/roles/health_check", roles.HealthCheckHandler)
	r.GET("/roles/:uid", roles.GetUserRolesHandler)
	r.PUT("/roles/set", roles.SetUserRolesHandler)

	return r
}

// function used to generate new instance of postgres persistence
func NewPostgresPersistence(url string) *db.PostgresPersistence {
	// generate base persistence layer
	base := utils.NewBasePersistence(url)
	return &db.PostgresPersistence{
		BasePostgresPersistence: base,
	}
}
