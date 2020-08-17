package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/walln/flurry2/flurry/auth"
)

type Proxy struct {
	host          string
	endpoint      string
	methods       []string
	authenticated bool
	authMethod    string
}

func NewProxy(host, endpoint string, methods []string, authMethod string, authenticated bool) Proxy {
	p := Proxy{host: host, endpoint: endpoint, methods: methods, authMethod: authMethod, authenticated: authenticated}
	return p
}

func (p *Proxy) GetAuthenticated() bool {
	return p.authenticated
}

func (p *Proxy) GetEndpoint() string {
	return p.endpoint
}

func (p *Proxy) GetMethods() []string {
	return p.methods
}

func (p *Proxy) Handle(rw http.ResponseWriter, req *http.Request) {

	logger := log.WithFields(log.Fields{
		"Proxy Route":    req.URL.Host + req.URL.RequestURI(),
		"Client IP":      req.RemoteAddr,
		"Request Method": req.Method,
	})

	validated := true

	// All routes that require authentication middleware must be validated before sending the request through the proxy.
	if (p.GetAuthenticated()) && p.authMethod == "FIREBASE" {
		validated = auth.AuthenticateWithFirebase(rw, req)
	} else if (p.GetAuthenticated()) && p.authMethod == "JWT" {
		validated = auth.AuthenticateWithJWT(rw, req)
	}

	if validated {
		demourl, err := url.Parse(p.host)
		if err != nil {
			log.Fatal(err)
		}

		req.Host = demourl.Host
		req.URL.Host = demourl.Host
		req.URL.Scheme = demourl.Scheme
		req.RequestURI = ""

		proxyReq, err := http.NewRequest(req.Method, req.URL.RequestURI(), req.Body)
		proxyReq.URL.Scheme = demourl.Scheme
		proxyReq.URL.Host = demourl.Host
		if err != nil {
			log.Fatal(err)
		}

		proxyReq.Header.Set("Host", req.Host)
		proxyReq.Header.Set("X-Forwarded-For", req.RemoteAddr)

		for header, values := range req.Header {
			for _, value := range values {
				proxyReq.Header.Add(header, value)
			}
		}

		response, err := http.DefaultClient.Do(proxyReq)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(rw, err)
			return
		}
		rw.WriteHeader(response.StatusCode)
		for key, values := range response.Header {
			for _, value := range values {
				rw.Header().Set(key, value)
			}
		}

		_, err = io.Copy(rw, response.Body)
		if err != nil {
			log.Fatal(err)
		}

		logger.Info("Successfully forwarded request.")

		defer response.Body.Close()
	} else {
		rw.WriteHeader(http.StatusUnauthorized)

		type errorBody struct {
			Error string `json:"error"`
		}

		responseBody := errorBody{"Failed to Authenticate"}

		js, err := json.Marshal(responseBody)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		_, err = rw.Write(js)
		if err != nil {
			log.Println(err)
		}

		logger.Info("Request denied. Invalid authorization credentials.")

	}

}
