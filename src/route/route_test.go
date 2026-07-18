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
		"GET /api/v1/appointments": false,
		"GET /api/v1/decoration":   false,
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

func TestLegacyRoutesAreRemoved(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	SetupRoutes(router)
	removed := map[string]bool{
		"GET /api/v1/orders":        false,
		"GET /api/v1/shops":         false,
		"GET /api/v1/service-types": false,
		"GET /api/mini/orders":      false,
		"GET /api/mini/addresses":   false,
		"GET /api/mini/meal/dishes": false,
	}
	for _, route := range router.Routes() {
		key := route.Method + " " + route.Path
		if _, legacy := removed[key]; legacy {
			t.Errorf("legacy route still registered: %s", key)
		}
	}
}
