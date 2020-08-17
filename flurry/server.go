package flurry

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"github.com/walln/flurry2/flurry/config"
	"github.com/walln/flurry2/flurry/global"
	"github.com/walln/flurry2/flurry/proxy"
)

type server struct {
	router        *mux.Router
	FirebaseApp   *firebase.App
	authType      string
	signingMethod string
}

func (s *server) setupAuth(authMethod string) {
	s.authType = authMethod

	if s.authType == "FIREBASE" {
		app, err := firebase.NewApp(context.Background(), nil)
		if err != nil {
			log.Fatalf("Error initializing app: %v", err)
		}
		s.FirebaseApp = app
		global.SetFirebaseApp(app)
	} else if s.authType == "JWT" {
		global.SetSigningMethod(s.signingMethod)
	}
}

func Initialize() server {
	s := server{}
	s.router = mux.NewRouter()

	flurryConfig := config.ReadConfigFile()

	log.Infof("Initializing %v utilizing %v authentication.", flurryConfig.GetName(), flurryConfig.GetAuth())

	if flurryConfig.Authentication == "FIREBASE" {
		s.setupAuth(flurryConfig.Authentication)
	}

	for _, route := range flurryConfig.GetRoutes() {
		proxyRoute := proxy.NewProxy(route.GetHost(), route.GetEndpoint(), route.GetMethods(), flurryConfig.Authentication, route.GetAuthenticated())
		log.Infof("Creating proxy route for %v.", route.GetHost()+route.GetEndpoint())
		s.registerRoute(&proxyRoute)
	}

	return s

}

func (s *server) registerRoute(p *proxy.Proxy) {
	s.router.HandleFunc(p.GetEndpoint(), p.Handle).Methods(p.GetMethods()...)
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	if !strings.Contains(port, ":") {
		port = ":" + port
	}
	return port
}

func (s *server) ListenAndServe() {
	muxWithMiddlewares := http.TimeoutHandler(s.router, time.Second*30, "Timeout!")

	log.Infof("Flurry serving traffic on port %v.", getPort())
	err := http.ListenAndServe(getPort(), muxWithMiddlewares)
	if err != nil {
		log.Error(err)
	}
}
