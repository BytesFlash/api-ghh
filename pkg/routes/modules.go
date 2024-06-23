package routes

import (
	"fmt"
	"net/http"

	"github.com/Unleash/unleash-client-go/v3"
	"github.com/gin-gonic/gin"
	"github.com/imedcl/manager-api/pkg/data"
	"github.com/imedcl/manager-api/pkg/events"
	"github.com/sirupsen/logrus"
)

type moduleResponse struct {
	Data *data.Module `json:"data"`
}

// @Summary modules
// @Description modules
// @Tags module
// @security barerToken
// @Accept json
// @Produce json
// @Success 200 {object} []data.Module
// @failure 404 {object} MessageResponse
// @Router /modules [get]
func ModulesGetRoute(router *gin.RouterGroup, db *data.DB, client *unleash.Client) {
	router.GET("/modules", func(c *gin.Context) {
		modules, err := db.ListAllModules()
		if err != nil {
			c.JSON(http.StatusNotFound, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusNotFound,
			})
			logrus.Println("error al obtener el módulos")
			return
		}
		c.JSON(http.StatusOK, modules)
	})

	router.POST("/modules", func(c *gin.Context) {

	})
	router.PUT("/modules/:module", func(c *gin.Context) {

	})
}

// @Summary delete modules
// @Description delete modules
// @Tags module
// @security barerToken
// @Accept json
// @Produce json
// @Success 200 {object} MessageResponse
// @failure 400 {object} MessageResponse
// @Router /modules/{name} [delete]
// @param name path string true "name"
func ModulesDeleteRoute(router *gin.RouterGroup, db *data.DB, client *unleash.Client) {
	// Delete Module
	router.DELETE("/modules/:name", func(c *gin.Context) {

		name := c.Param("name")

		module, err := db.GetModuleByName(name)
		if err != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "Módulo no encontrado",
				Code:    http.StatusBadRequest,
			})
			return
		}

		if db.DeleteModule(module) != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: "Error al eliminar módulo, intente nuevamente",
				Code:    http.StatusBadRequest,
			})
			return
		}
		PrtyParams, _ := events.PrettyParams(c.Params)
		event := &events.EventLog{
			UserNickname: currentUser.NickName,
			Resource:     "Módulos",
			Event:        fmt.Sprintf("Se ha eliminado al módulo %s", name),
			Params:       PrtyParams,
		}
		event.Write()
		c.JSON(http.StatusOK, MessageResponse{Details: "Módulo eliminado exitosamente", Code: http.StatusOK})
	})

}
