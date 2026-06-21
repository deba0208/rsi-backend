package handler

import (
	"net/http"

	"github.com/deba0208/stock-rsi-dashboard/internal/service"
	"github.com/gin-gonic/gin"
)

type MetricHandler struct {
	service *service.MetricService
}

func NewMetricHandler(
	service *service.MetricService,
) *MetricHandler {

	return &MetricHandler{
		service: service,
	}
}

func (h *MetricHandler) GetTop50ByCriteria(
	c *gin.Context,
) {

	timeFrame :=
		c.Query("timeFrame")

	if timeFrame == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "missing required query param: timeFrame",
			},
		)
		return
	}

	metrics, err :=
		h.service.GetTopByTimeFrame(
			timeFrame,
			50,
		)

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		metrics,
	)
}
