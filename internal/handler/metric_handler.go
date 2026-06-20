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

	result, err :=
		h.service.GetTop50ByCriteria(c.Param("criteria"))

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
		result,
	)
}
