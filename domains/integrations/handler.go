package integrations

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct{ svc Service }

func NewHandler(svc Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) Register(g *echo.Group, authMW echo.MiddlewareFunc) {
	pub := g.Group("/integrations/google")
	pub.GET("/callback", h.callback)

	priv := g.Group("/integrations/google", authMW)
	priv.GET("/connect", h.connect)
}

func (h *Handler) connect(c echo.Context) error {
	uid := c.Get("userID").(int)
	url, err := h.svc.ConnectURL(c.Request().Context(), uid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.Redirect(http.StatusFound, url)
}

func (h *Handler) callback(c echo.Context) error {
	if err := h.svc.ExchangeCallback(c.Request().Context(), c.QueryParam("state"), c.QueryParam("code")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.String(http.StatusOK, "Google connected ✓")
}
