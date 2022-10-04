package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qulaz/gas-price-test/internal/usecase"
	"github.com/qulaz/gas-price-test/pkg/logging"
)

type gasGraphRoutes struct {
	gas    usecase.GasGraph
	logger logging.ContextLogger
}

func newGasGraphRoutes(handler *gin.RouterGroup, g usecase.GasGraph, l logging.ContextLogger) {
	r := &gasGraphRoutes{g, l}

	h := handler.Group("/")
	{
		h.GET("/graph", r.getGraph)
	}
}

// @Summary     Get gas statistic graph
// @ID          get-graph
// @Accept      json
// @Produce     json
// @Success     200 {object} entity.GasResult
// @Failure     500 {object} response
// @Router      /graph [get].
func (g *gasGraphRoutes) getGraph(c *gin.Context) {
	res, err := g.gas.Calculate(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}
