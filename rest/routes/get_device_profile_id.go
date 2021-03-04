package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/go_cache/model/profile"
	"github.com/nmarsollier/go_cache/utils/errors"
)

// Servicio REST que nos retorna información de un dialogo a mostrar en pantalla
// Vamos a usar el contexto como un Builder Pattern
func init() {
	router().GET(
		"/profile",
		getProfile,
	)
}

func getProfile(c *gin.Context) {
	data := profile.FetchProfile("123")

	if data == nil {
		c.Error(errors.NewCustomError(500, "Internal Server Error"))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"login": data.Login,
		"web":   data.Web,
		"name":  data.Name,
	})
}
