package actions

import (
	"net/http"
	"sync"

	"projectcollection/locales"
	"projectcollection/models"
	"projectcollection/notification"
	"projectcollection/public"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v3/pop/popmw"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/middleware/csrf"
	"github.com/gobuffalo/middleware/forcessl"
	"github.com/gobuffalo/middleware/i18n"
	"github.com/gobuffalo/middleware/paramlogger"
	"github.com/unrolled/secure"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

var (
	app             *buffalo.App
	appOnce         sync.Once
	T               *i18n.Translator
	NotificationHub *notification.Hub
)

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	appOnce.Do(func() {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_projectcollection_session",
		})

		NotificationHub = notification.NewHub()
		// Start the notification hub
		go NotificationHub.Run()

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//   c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))
		// Setup and use translations:
		app.Use(translations())

		app.GET("/", HomeHandler)

		// setup the active menu
		app.Use(activeMenu)
		// add the notification hub to the context
		app.Use(notificationHub)

		// NOTE: this block should go before any resources
		// that need to be protected by buffalo-goth!
		//AuthMiddlewares
		app.Use(SetCurrentUser)
		app.Use(Authorize)

		//Routes for Auth
		auth := app.Group("/auth")
		auth.GET("/", AuthLanding)
		auth.GET("/new", AuthNew)
		auth.POST("/", AuthCreate)
		auth.DELETE("/", AuthDestroy)
		auth.Middleware.Skip(Authorize, AuthLanding, AuthNew, AuthCreate)

		//Routes for User registration
		users := app.Group("/users")
		users.GET("/new", UsersNew)
		users.POST("/", UsersCreate)
		users.Middleware.Remove(Authorize)

		app.Resource("/dbtypes", DbtypesResource{})
		app.Resource("/runtimes", RuntimesResource{})

		app.Resource("/environments", EnvironmentsResource{})

		app.Resource("/hosts", HostsResource{})
		app.Resource("/clients", ClientsResource{})
		app.Resource("/projects", ProjectsResource{})
		app.Resource("/applications", ApplicationsResource{})
		app.GET("/job_applications", JobApplicationsResource{}.List)
		app.GET("/job_applications/{job_application_id}", JobApplicationsResource{}.Show)
		app.DELETE("/job_applications/{job_application_id}", JobApplicationsResource{}.Destroy)
		app.GET("/ws", JobApplicationsResource{}.Websocket)
		app.GET("/wsbc", JobApplicationsResource{}.WebsocketBroadcast)
		// the application job has to be able to broadcast to the websocket without authorization
		app.Middleware.Skip(Authorize, JobApplicationsResource{}.WebsocketBroadcast)

		app.Resource("/aliases", AliasesResource{})
		app.Resource("/pools", PoolsResource{})
		app.ServeFiles("/", http.FS(public.FS())) // serve files from the public directory
	})

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(locales.FS(), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}

// activeMenu is a helper function to determine which menu item to highlight
// based on the current request path.
func activeMenu(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		firstPath := c.Request().URL.Path[1:]
		c.Set("activeMenu", firstPath)
		c.Set("activeClass", func(path, menuName string) string {
			if path == menuName {
				return "active"
			}
			return ""
		})
		return next(c)
	}
}

// notificationHub adds the notification hub to the context
func notificationHub(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		c.Set("hub", NotificationHub)
		return next(c)
	}
}
