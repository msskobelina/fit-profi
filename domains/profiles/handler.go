package profiles

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
	pr := g.Group("/profiles", authMW)

	pr.GET("/user", h.getUserSelf)
	pr.POST("/user", h.createUserSelf)
	pr.PUT("/user", h.updateUserSelf)

	pr.GET("/coach", h.getCoachSelf)
	pr.POST("/coach", h.createCoachSelf)
	pr.PUT("/coach", h.updateCoachSelf)

	pr.GET("/user/:userId", h.adminGetUser)
}

// swagger:route GET /profiles/user Profiles profilesGetUserSelf
// Get my user profile
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:userProfile
//	401: body:errorResponse
func (h *Handler) getUserSelf(c echo.Context) error {
	uid := c.Get("userID").(int)
	role := c.Get("userRole").(string)
	p, err := h.service.GetUserProfile(c.Request().Context(), uid, role, uid)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "not found"})
	}

	return c.JSON(http.StatusOK, p)
}

// swagger:route POST /profiles/user Profiles profilesCreateUserSelf
// Create my user profile
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:userProfile
//	400: body:errorResponse
//	401: body:errorResponse
func (h *Handler) createUserSelf(c echo.Context) error {
	req := new(CreateUserProfileRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
	}
	uid := c.Get("userID").(int)
	out, err := h.service.CreateUserProfile(c.Request().Context(), uid, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, out)
}

// swagger:route PUT /profiles/user Profiles profilesUpdateUserSelf
// Update my user profile
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:userProfile
//	400: body:errorResponse
//	401: body:errorResponse
func (h *Handler) updateUserSelf(c echo.Context) error {
	req := new(UpdateUserProfileRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
	}
	uid := c.Get("userID").(int)
	out, err := h.service.UpdateUserProfile(c.Request().Context(), uid, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, out)
}

// swagger:route GET /profiles/coach Profiles profilesGetCoachSelf
// Get my coach profile (admin only)
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:coachProfile
//	401: body:errorResponse
//	403: body:errorResponse
func (h *Handler) getCoachSelf(c echo.Context) error {
	role := c.Get("userRole").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, ErrorResponse{Error: "admin only"})
	}
	uid := c.Get("userID").(int)
	p, err := h.service.GetCoachProfile(c.Request().Context(), role, uid)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "not found"})
	}

	return c.JSON(http.StatusOK, p)
}

// swagger:route POST /profiles/coach Profiles profilesCreateCoachSelf
// Create my coach profile (admin only)
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:coachProfile
//	400: body:errorResponse
//	401: body:errorResponse
//	403: body:errorResponse
func (h *Handler) createCoachSelf(c echo.Context) error {
	role := c.Get("userRole").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, ErrorResponse{Error: "admin only"})
	}
	req := new(CreateCoachProfileRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
	}
	uid := c.Get("userID").(int)
	out, err := h.service.CreateCoachProfile(c.Request().Context(), uid, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, out)
}

// swagger:route PUT /profiles/coach Profiles profilesUpdateCoachSelf
// Update my coach profile (admin only)
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:coachProfile
//	400: body:errorResponse
//	401: body:errorResponse
//	403: body:errorResponse
func (h *Handler) updateCoachSelf(c echo.Context) error {
	role := c.Get("userRole").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, ErrorResponse{Error: "admin only"})
	}
	req := new(UpdateCoachProfileRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
	}
	uid := c.Get("userID").(int)
	out, err := h.service.UpdateCoachProfile(c.Request().Context(), uid, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, out)
}

// swagger:route GET /profiles/user/{userId} Profiles profilesAdminGetUser
// Admin: get user profile by userId
//
// security:
//   - Bearer: []
//
// responses:
//
//	200: body:userProfile
//	401: body:errorResponse
//	403: body:errorResponse
func (h *Handler) adminGetUser(c echo.Context) error {
	role := c.Get("userRole").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, ErrorResponse{Error: "admin only"})
	}
	id, _ := strconv.Atoi(c.Param("userId"))
	p, err := h.service.GetUserProfile(c.Request().Context(), 0, role, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "not found"})
	}

	return c.JSON(http.StatusOK, p)
}
