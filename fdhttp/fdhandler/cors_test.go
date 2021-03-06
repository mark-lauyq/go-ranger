package fdhandler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/foodora/go-ranger/fdhttp"
	"github.com/foodora/go-ranger/fdhttp/fdhandler"
	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware(t *testing.T) {
	origin := "https://api.foodora.com"
	corsMiddleware := fdhandler.NewCORSMiddleware(origin)

	router := fdhttp.NewRouter()
	router.Use(corsMiddleware)
	router.StdGET("/foo", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "bar")
	})

	req := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, origin, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_IsIgnoredIfHandlerSetted(t *testing.T) {
	origin := "https://api.foodora.com"
	corsMiddleware := fdhandler.NewCORSMiddleware(origin)

	router := fdhttp.NewRouter()
	router.Use(corsMiddleware)
	router.GET("/foo", func(ctx context.Context) (int, interface{}) {
		fdhttp.SetResponseHeaderValue(ctx, "Access-Control-Allow-Origin", "*")
		return http.StatusOK, nil
	})

	// standard handler
	req := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, fdhandler.CORSOriginAll, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_IsIgnoredIfStdHandlerSetted(t *testing.T) {
	origin := "https://api.foodora.com"
	corsMiddleware := fdhandler.NewCORSMiddleware(origin)

	router := fdhttp.NewRouter()
	router.Use(corsMiddleware)
	router.StdGET("/foo", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "bar")
	})

	req := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, fdhandler.CORSOriginAll, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestNewCORS(t *testing.T) {
	corsHandler := fdhandler.NewCORS()
	corsHandler.Origin = "https://api.foodora.com"
	corsHandler.Credentials = true
	corsHandler.Methods = []string{
		"OPTIONS",
		"GET",
		"PUT",
	}
	corsHandler.ExposeHeaders = []string{
		"X-Personal-One",
		"X-Personal-Two",
	}
	corsHandler.MaxAge = 25 * time.Minute

	router := fdhttp.NewRouter()
	router.Register(corsHandler)

	req := httptest.NewRequest(http.MethodOptions, "/foo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, corsHandler.Origin, w.Header().Get("Access-Control-Allow-Origin"))
	assert.ElementsMatch(t, []string{
		"OPTIONS",
		"GET",
		"PUT",
	}, strings.Split(w.Header().Get("Access-Control-Allow-Methods"), ", "))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "X-Personal-One,X-Personal-Two", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "1500", w.Header().Get("Access-Control-Max-Age"))
}

func TestNewCORS_LoadMethodsFromRouter(t *testing.T) {
	corsHandler := fdhandler.NewCORS()
	corsHandler.MaxAge = 0

	router := fdhttp.NewRouter()
	router.Register(corsHandler)
	router.StdGET("/foo", func(w http.ResponseWriter, req *http.Request) {})
	router.PUT("/foo", func(ctx context.Context) (int, interface{}) {
		return http.StatusOK, nil
	})

	req := httptest.NewRequest(http.MethodOptions, "/foo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, corsHandler.Origin, w.Header().Get("Access-Control-Allow-Origin"))
	assert.ElementsMatch(t, []string{
		"OPTIONS",
		"GET",
		"PUT",
	}, strings.Split(w.Header().Get("Access-Control-Allow-Methods"), ", "))
	_, ok := w.HeaderMap["Access-Control-Allow-Credentials"]
	assert.False(t, ok)
	_, ok = w.HeaderMap["Access-Control-Allow-Headers"]
	assert.False(t, ok)
	_, ok = w.HeaderMap["Access-Control-Max-Age"]
	assert.False(t, ok)
}

func TestNewCORS_WithCredentialsDisabled(t *testing.T) {
	corsHandler := fdhandler.NewCORS()
	corsHandler.MaxAge = 0

	router := fdhttp.NewRouter()
	router.Register(corsHandler)

	req := httptest.NewRequest(http.MethodOptions, "/foo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, fdhandler.CORSOriginAll, w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	_, ok := w.HeaderMap["Access-Control-Allow-Credentials"]
	assert.False(t, ok)
	_, ok = w.HeaderMap["Access-Control-Allow-Headers"]
	assert.False(t, ok)
	_, ok = w.HeaderMap["Access-Control-Max-Age"]
	assert.False(t, ok)
}

func TestNewCORS_SubRouterMethodsAreReturned(t *testing.T) {
	corsHandler := fdhandler.NewCORS()

	router := fdhttp.NewRouter()
	router.Register(corsHandler)
	router.StdGET("/foo", func(w http.ResponseWriter, req *http.Request) {})
	router.PUT("/foo", func(ctx context.Context) (int, interface{}) {
		return http.StatusOK, nil
	})

	subrouter := router.SubRouter()
	subrouter.Prefix = "/foo"
	subrouter.StdPOST("/bar", func(w http.ResponseWriter, req *http.Request) {})
	subrouter.DELETE("/bar", func(ctx context.Context) (int, interface{}) {
		return http.StatusOK, nil
	})

	req := httptest.NewRequest(http.MethodOptions, "/foo", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.ElementsMatch(t, []string{
		"OPTIONS",
		"GET",
		"PUT",
		"POST",
		"DELETE",
	}, strings.Split(w.Header().Get("Access-Control-Allow-Methods"), ", "))
}
