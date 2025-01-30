package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"githum.com/leebrouse/urlshortener/internal/model"
)

type URLservice interface {
	CreateURL(ctx context.Context, req model.CreateURLRequest) (*model.CreateURLResponse, error)

	GetURLByShortCode(ctx context.Context, shortCode string) (string, error)
}

type URLHandler struct {
	urlservice URLservice
}

// POST /api/url original_url,custom_code,duration=> 短链接，过期时间
func (h *URLHandler) CreateURL(c echo.Context) error {
	//数据提取
	var req model.CreateURLRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	//验证数据格式
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	//调用业务函数
	resp, err := h.urlservice.CreateURL(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//返回响应
	return c.JSON(http.StatusCreated, resp)
}

// GET /:code =>把短url重定向到长url
func (h *URLHandler) RedirectURL(c echo.Context) error {
	//把code 取出来
	shortCode := c.Param("code")
	//根据shortcode => url调用业务函数
	originalURL, err := h.urlservice.GetURLByShortCode(c.Request().Context(), shortCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Redirect(http.StatusPermanentRedirect, originalURL)
	return nil
}
