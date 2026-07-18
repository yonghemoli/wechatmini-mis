package route

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetupRoutesIncludesMiniBusinessEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	SetupRoutes(router)
	want := map[string]bool{
		"GET /api/mini/services":   false,
		"GET /api/mini/caregivers": false,
		"POST /api/mini/demands":   false,
		"GET /api/v1/caregivers":   false,
		"GET /api/v1/demands":      false,
	}
	for _, route := range router.Routes() {
		key := route.Method + " " + route.Path
		if _, ok := want[key]; ok {
			want[key] = true
		}
	}
	for route, found := range want {
		if !found {
			t.Errorf("route not registered: %s", route)
		}
	}
}
