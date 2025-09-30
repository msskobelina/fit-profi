package nutrition

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct{ service Service }

func NewHandler(svc Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) Register(g *echo.Group, authMW echo.MiddlewareFunc) {
	ng := g.Group("/nutrition", authMW)

	ng.POST("/entries", h.create)
	ng.GET("/entries", h.list)
	ng.GET("/entries/:id", h.get)
	ng.PUT("/entries/:id", h.update)
	ng.DELETE("/entries/:id", h.delete)

	ng.GET("/summary", h.summary)
}

// swagger:route POST /nutrition/entries Nutrition nutritionCreate
// Create diary entry (self)
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:diaryEntry
//	400: body:errorResponse
//	401: body:errorResponse
func (h *Handler) create(c echo.Context) error {
	req := new(CreateDiaryEntryRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
	}
	uid := c.Get("userID").(int)
	out, err := h.service.Create(c.Request().Context(), uid, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}

// swagger:route GET /nutrition/entries Nutrition nutritionList
// List diary entries (self).
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:[]diaryEntry
//	400: body:errorResponse
//	401: body:errorResponse
func (h *Handler) list(c echo.Context) error {
	uid := c.Get("userID").(int)
	date := c.QueryParam("date")
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	out, err := h.service.List(c.Request().Context(), uid, date, from, to)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}

// swagger:route GET /nutrition/entries/{id} Nutrition nutritionGet
// Get diary entry by id (self or admin)
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:diaryEntry
//	401: body:errorResponse
//	403: body:errorResponse
func (h *Handler) get(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	role := c.Get("userRole").(string)
	uid := c.Get("userID").(int)
	out, err := h.service.Get(c.Request().Context(), uid, id, role)
	if err != nil {
		code := http.StatusForbidden
		if err.Error() == "record not found" {
			code = http.StatusNotFound
		}
		return c.JSON(code, ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}

// swagger:route PUT /nutrition/entries/{id} Nutrition nutritionUpdate
// Update diary entry (self)
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:diaryEntry
//	400: body:errorResponse
//	401: body:errorResponse
//	403: body:errorResponse
func (h *Handler) update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	req := new(UpdateDiaryEntryRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
	}
	uid := c.Get("userID").(int)
	out, err := h.service.Update(c.Request().Context(), uid, id, req)
	if err != nil {
		code := http.StatusBadRequest
		if err.Error() == "forbidden" {
			code = http.StatusForbidden
		}
		return c.JSON(code, ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}

// swagger:route DELETE /nutrition/entries/{id} Nutrition nutritionDelete
// Delete diary entry (self or admin)
//
// security:
//   - Bearer: []
//
// responses:
//
//	204: description: no content
//	401: body:errorResponse
//	403: body:errorResponse
func (h *Handler) delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	role := c.Get("userRole").(string)
	uid := c.Get("userID").(int)
	if err := h.service.Delete(c.Request().Context(), uid, id, role); err != nil {
		code := http.StatusForbidden
		return c.JSON(code, ErrorResponse{Error: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// swagger:route GET /nutrition/summary Nutrition nutritionSummary
// Summary for date range
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:diarySummary
//	400: body:errorResponse
//	401: body:errorResponse
func (h *Handler) summary(c echo.Context) error {
	uid := c.Get("userID").(int)
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	out, err := h.service.Summary(c.Request().Context(), uid, from, to)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}
