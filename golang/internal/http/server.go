package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/riskiramdan/efishery/golang/config"
	"github.com/riskiramdan/efishery/golang/internal/concurrency"
	"github.com/riskiramdan/efishery/golang/internal/data"
	"github.com/riskiramdan/efishery/golang/internal/hosts"
	"github.com/riskiramdan/efishery/golang/internal/http/controller"
	"github.com/riskiramdan/efishery/golang/internal/user"
	"github.com/riskiramdan/efishery/golang/util"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

// Server represents the http server that handles the requests
type Server struct {
	dataManager           *data.Manager
	utility               *util.Utility
	userService           user.ServiceInterface
	userController        *controller.UserController
	concurrencyService    concurrency.ServiceInterface
	concurrencyController *controller.ConcurrencyController
	httpManager           *hosts.HTTPManager
}

func (hs *Server) authMethod(r chi.Router, method string, path string, handler http.HandlerFunc) {
	r.With(
		hs.instrument(method, "/v1"+path),
	).Method(method, path, handler)
}

func (hs *Server) compileRouter() chi.Router {
	r := chi.NewRouter()

	// Base middlewares
	//

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Basic CORS
	//Routes()
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Access-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	// Prometheus handler
	//

	// Add routes

	r.HandleFunc("/login", hs.userController.PostLogin)
	r.HandleFunc("/register", hs.userController.PostCreateUser)
	r.HandleFunc("/verify", hs.userController.PostVerifyToken)

	r.Route("/auth", func(r chi.Router) {
		r.Use(hs.authorizedOnly(hs.userService))

		hs.authMethod(r, "GET", "/users", hs.userController.GetListUser)
		hs.authMethod(r, "GET", "/prices", hs.concurrencyController.GetListPrices)
		hs.authMethod(r, "POST", "/verify", hs.userController.PostVerifyToken)

		r.Mount("/admin", adminRouter(hs))
	})

	return r
}

// Serve serves http requests
func (hs *Server) Serve() {
	// Compile all the routes
	//

	r := hs.compileRouter()

	// Run the server + gracefully shutdown mechanism
	//

	log.Printf("About to listen on 8083. Go to http://127.0.0.1:8083")
	srv := http.Server{Addr: ":8083", Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

// NewServer creates a new http server
func NewServer(
	userService user.ServiceInterface,
	concurrencyService concurrency.ServiceInterface,
	dataManager *data.Manager,
	config *config.Config,
	utility *util.Utility,
	httpManager *hosts.HTTPManager,
) *Server {
	userController := controller.NewUserController(userService, dataManager, utility)
	concurrencyController := controller.NewConcurrencyController(concurrencyService, utility)
	return &Server{
		dataManager:           dataManager,
		userService:           userService,
		userController:        userController,
		utility:               utility,
		concurrencyService:    concurrencyService,
		concurrencyController: concurrencyController,
		httpManager:           httpManager,
	}
}
