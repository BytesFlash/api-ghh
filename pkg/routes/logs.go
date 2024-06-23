package routes

import (
	"net/http"

	"github.com/Unleash/unleash-client-go/v3"
	"github.com/gin-gonic/gin"
	"github.com/imedcl/manager-api/pkg/data"
	"github.com/imedcl/manager-api/pkg/events"
	"github.com/sirupsen/logrus"
)

type LogResponse struct {
	Data *data.Log `json:"data"`
}

type logParams struct {
	Resource string `form:"resource" binding:"required"`
	Input    string `form:"input" binding:"required"`
	Output   string `form:"output" binding:"required"`
}

// @Summary get Logs
// @Description get logs
// @Tags log
// @security barerToken
// @Accept json
// @Produce json
// @Success 200 {object} []data.Log
// @failure 404 {object} MessageResponse
// @Router /logs/{date} [get]
// @param date path string true "date"
// @param starthour query string true "starthour"
// @param finishhour query string true "finishhour"
func LogGetRoute(router *gin.RouterGroup, db *data.DB, client *unleash.Client) {
	router.GET("/logs/:date", func(c *gin.Context) {
		date := c.Param("date")
		startHour := c.Query("starthour")
		finishHour := c.Query("finishhour")
		logs, err := db.GetLogs(date, startHour, finishHour)
		if err != nil {
			c.JSON(http.StatusNotFound, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusNotFound,
			})
			logrus.Println("error al obtener logs")
			return
		}
		c.JSON(http.StatusOK, logs)
	})
}

// @Summary post Logs
// @Description post logs
// @Tags log
// @security barerToken
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @failure 400 {object} MessageResponse
// @Router /logs/trx [post]
// @Param log body logParams true "log"
func LogPostRoute(router *gin.RouterGroup, db *data.DB, client *unleash.Client) {
	router.POST("/logs/trx", func(c *gin.Context) {
		var params logParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, MessageResponse{
				Details: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		event := &events.EventLog{
			UserNickname: currentUser.NickName,
			Resource:     params.Resource,
			Event:        params.Output,
			Params:       params.Input,
		}
		event.Write()
		c.JSON(http.StatusOK, true)
	})
}
