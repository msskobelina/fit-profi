package programs

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct{ service Service }

func NewHandler(service Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) Register(g *echo.Group, authMW echo.MiddlewareFunc) {
	pg := g.Group("/programs", authMW)

	pg.POST("", h.create)
	pg.DELETE("/:id", h.delete)

	pg.GET("/:id", h.get)
	pg.GET("", h.listMine)

	pg.POST("/progress", h.addProgress)
}

// swagger:route POST /programs Programs programsCreate
// Create training program (admin)
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:trainingProgram
//	400: body:errorResponse
//	401: body:errorResponse
//	403: body:errorResponse
func (h *Handler) create(c echo.Context) error {
	role := c.Get("userRole").(string)
	coachID := c.Get("userID").(int)
	in := new(CreateProgramRequest)
	if err := c.Bind(in); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "bad body"})
	}
	out, err := h.service.Create(c.Request().Context(), coachID, role, in)
	if err != nil {
		return c.JSON(http.StatusForbidden, ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, out)
}

// swagger:route GET /programs/{id} Programs programsGet
// Get program by ID (owner or admin)
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:trainingProgram
//	401: body:errorResponse
//	403: body:errorResponse
func (h *Handler) get(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	role := c.Get("userRole").(string)
	uid := c.Get("userID").(int)
	out, err := h.service.Get(c.Request().Context(), uid, role, id)
	if err != nil {
		return c.JSON(http.StatusForbidden, ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, out)
}

// swagger:route GET /programs Programs programsListMine
// List my programs (client)
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:[]trainingProgram
//	401: body:errorResponse
func (h *Handler) listMine(c echo.Context) error {
	uid := c.Get("userID").(int)
	out, err := h.service.ListByUser(c.Request().Context(), uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed"})
	}

	return c.JSON(http.StatusOK, out)
}

// swagger:route DELETE /programs/{id} Programs programsDelete
// Delete program (admin)
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
	role := c.Get("userRole").(string)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(c.Request().Context(), role, id); err != nil {
		return c.JSON(http.StatusForbidden, ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// swagger:route POST /programs/progress Programs programsAddProgress
// Add exercise progress for current week (owner)
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:exerciseProgress
//	400: body:errorResponse
//	401: body:errorResponse
//	403: body:errorResponse
func (h *Handler) addProgress(c echo.Context) error {
	uid := c.Get("userID").(int)
	in := new(AddProgressRequest)
	if err := c.Bind(in); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "bad body"})
	}
	out, err := h.service.AddProgress(c.Request().Context(), uid, in)
	if err != nil {
		return c.JSON(http.StatusForbidden, ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, out)
}
