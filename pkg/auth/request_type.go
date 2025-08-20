package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type requestType string

const (
	requestTypeWeb requestType = "web"
	requestTypeApi requestType = "api"
	requestTypeUnknown requestType = "unknown"
)

func getRequestType(c *fiber.Ctx) requestType {
	// Проверяем различные признаки API запроса
	accept := c.Get("Accept")
	contentType := c.Get("Content-Type")
	userAgent := c.Get("User-Agent")
	requestedWith := c.Get("X-Requested-With")
	
	if strings.Contains(accept, "application/json") ||
		strings.Contains(contentType, "application/json") ||
		requestedWith == "XMLHttpRequest" {
		return requestTypeApi
	}

	// Проверяем на браузер
	if strings.Contains(userAgent, "Mozilla") ||
		strings.Contains(userAgent, "Chrome") ||
		strings.Contains(userAgent, "Safari") ||
		strings.Contains(userAgent, "Firefox") ||
		strings.Contains(accept, "text/html") {
		return requestTypeWeb
	}

	return requestTypeUnknown
}