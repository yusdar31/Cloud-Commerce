package gateway

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// ServiceConfig defines a backend service route.
type ServiceConfig struct {
	Name      string
	BasePath  string // e.g. "/api/v1/auth"
	TargetURL string // e.g. "http://user-service:8081"
	StripPath bool   // if true, strips the base path before forwarding
	Protected bool   // if true, requires JWT authentication
}

// Router holds reverse proxies for all backend services.
type Router struct {
	proxies map[string]*httputil.ReverseProxy
	configs []ServiceConfig
}

// NewRouter creates a new gateway router from service configs.
func NewRouter(configs []ServiceConfig) *Router {
	r := &Router{
		proxies: make(map[string]*httputil.ReverseProxy),
		configs: configs,
	}

	for _, cfg := range configs {
		target, err := url.Parse(cfg.TargetURL)
		if err != nil {
			panic("invalid target URL for " + cfg.Name + ": " + err.Error())
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = target.Host

			if cfg.StripPath {
				req.URL.Path = strings.TrimPrefix(req.URL.Path, cfg.BasePath)
				if !strings.HasPrefix(req.URL.Path, "/") {
					req.URL.Path = "/" + req.URL.Path
				}
			}
		}
		proxy.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Del("Access-Control-Allow-Origin")
			resp.Header.Del("Access-Control-Allow-Methods")
			resp.Header.Del("Access-Control-Allow-Headers")
			resp.Header.Del("Access-Control-Allow-Credentials")
			resp.Header.Del("Access-Control-Expose-Headers")
			resp.Header.Del("Access-Control-Max-Age")
			return nil
		}
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"success":false,"error":{"type":"gateway_error","title":"Service Unavailable","status":502,"detail":"` + cfg.Name + ` is not responding"}}`))
		}

		r.proxies[cfg.BasePath] = proxy
	}

	return r
}

// HandleProxy returns a gin handler that proxies to the matching backend service.
func (r *Router) HandleProxy(cfg ServiceConfig) gin.HandlerFunc {
	proxy, ok := r.proxies[cfg.BasePath]
	if !ok {
		return func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"type":   "gateway_error",
					"title":  "Route Not Found",
					"status": 500,
					"detail": "No proxy configured for " + cfg.BasePath,
				},
			})
		}
	}

	return func(c *gin.Context) {
		// Forward X-Forwarded headers
		c.Request.Header.Set("X-Forwarded-For", c.ClientIP())
		c.Request.Header.Set("X-Forwarded-Proto", "http")
		if c.Request.Header.Get("X-Request-ID") == "" {
			if requestID, exists := c.Get("request_id"); exists {
				c.Request.Header.Set("X-Request-ID", requestID.(string))
			}
		}

		// Forward tenant ID from context if present
		if tenantID, exists := c.Get("tenant_id"); exists {
			c.Request.Header.Set("X-Tenant-ID", tenantID.(string))
		}
		// Forward user ID from context if present
		if userID, exists := c.Get("user_id"); exists {
			c.Request.Header.Set("X-User-ID", userID.(string))
		}
		// Forward user role from context if present
		if role, exists := c.Get("user_role"); exists {
			c.Request.Header.Set("X-User-Role", role.(string))
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// Configs returns the default service routing configuration.
func DefaultConfigs() []ServiceConfig {
	return []ServiceConfig{
		{
			Name:      "identity-service",
			BasePath:  "/api/v1/auth",
			TargetURL: getServiceURL("USER_SERVICE_URL", "http://localhost:8081"),
			Protected: false,
		},
		{
			Name:      "identity-service",
			BasePath:  "/api/v1/users",
			TargetURL: getServiceURL("USER_SERVICE_URL", "http://localhost:8081"),
			Protected: true,
		},
		{
			Name:      "store-service",
			BasePath:  "/api/v1/stores",
			TargetURL: getServiceURL("STORE_SERVICE_URL", "http://localhost:8082"),
			Protected: false,
		},
		{
			Name:      "catalog-service",
			BasePath:  "/api/v1/products",
			TargetURL: getServiceURL("PRODUCT_SERVICE_URL", "http://localhost:8083"),
			Protected: true,
		},
		{
			Name:      "catalog-service",
			BasePath:  "/api/v1/categories",
			TargetURL: getServiceURL("PRODUCT_SERVICE_URL", "http://localhost:8083"),
			Protected: true,
		},
	}
}

func getServiceURL(envKey, defaultVal string) string {
	if val := os.Getenv(envKey); val != "" {
		return val
	}
	return defaultVal
}
