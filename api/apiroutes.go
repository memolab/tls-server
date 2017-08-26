package api

import (
	"time"
	"tls-server/api/middlewares"
	"tls-server/api/types"
)

// initRoutes init middlewares and handlres return routes
func initRoutes(c *APICtl, config map[string]string) *[]route {

	// new FrontMiddleware
	frontMiddleware := middlewares.NewFrontMiddleware(c, config["accessLogsDump"],
		"Content-Type,"+config["headerTokenKey"])
	middFront := frontMiddleware.Handler()
	c.regMidd["front"] = frontMiddleware
	//

	// new CacheMiddleware
	c.cache = middlewares.NewCacheMiddleware(c)
	d5min, _ := time.ParseDuration("5m")
	//

	// new AuthMiddleware
	c.auth = middlewares.NewAuthMiddleware(c, config["headerTokenKey"], config["rateLimiteAPI"],
		config["secretKey1"], config["secretKey2"], config["rateLimiteLogsDump"])
	middAuth := c.auth.Handler()
	middRateLimitIP := c.auth.RateLimitIPHandler(config["headerClientIPs"], config["rateLimiteIP"])
	//

	// New isadminMiddleware handler
	isAdmin := middlewares.IsAdmin(c, config["adminID"])

	// return routes
	return &[]route{
		// API routes
		route{
			url:         "/",
			handler:     c.indexHandler,
			middlewares: []types.MiddlewareHandler{middFront},
		},
		route{
			url:         "/signin",
			handler:     c.signInHandler,
			middlewares: []types.MiddlewareHandler{middRateLimitIP, middFront},
		},
		route{
			url:         "/signup",
			handler:     c.signUpHandler,
			middlewares: []types.MiddlewareHandler{middRateLimitIP, middFront},
		},

		route{
			url:         "/initdb",
			handler:     c.initDBHandler,
			middlewares: []types.MiddlewareHandler{middAuth, middFront},
		},
		route{
			url:     "/user",
			handler: c.userIndexHandler,
			middlewares: []types.MiddlewareHandler{
				c.cache.CacheHandler("/user", map[string]string{"tokenUID": ""}, d5min),
				middAuth,
				middFront},
		},
		route{
			url:     "/user2",
			handler: c.user2IndexHandler,
			middlewares: []types.MiddlewareHandler{
				//c.cache.CacheHandler("/user2", map[string]string{"tokenUID": ""}, d5min),
				middAuth, middFront},
		},

		//Admin routes
		route{
			url:         "/admin",
			handler:     c.adminIndexHandler,
			middlewares: []types.MiddlewareHandler{middFront},
		},
		route{
			url:         "/admin/overview",
			handler:     c.adminOverviewHandler,
			middlewares: []types.MiddlewareHandler{middFront},
		},
		route{
			url:         "/admin/accesslogs",
			handler:     c.adminAccesslogsHandler,
			middlewares: []types.MiddlewareHandler{isAdmin, middAuth, middFront},
		},
	}

}
