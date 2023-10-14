package api

import "github.com/labstack/echo/v4"

type apiV1 struct {
}

func NewApiV1() *apiV1 {
	return &apiV1{}
}

func (s *apiV1) RegisterRoutes(e *echo.Echo) {
	e.GET("/health", s.handleHealth)
}
