package config

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"
)

type config struct {
	Name           string   `yaml:"name"`
	Authentication string   `yaml:"authentication"`
	SigningMethod  string   `yaml:"signingMethod"`
	Routes         []*route `yaml:"routes"`
}

type route struct {
	Endpoint      string   `yaml:"endpoint"`
	Host          string   `yaml:"host"`
	Authenticated bool     `yaml:"authenticated"`
	Methods       []string `yaml:"methods"`
}

func (c *config) GetRoutes() []*route {
	return c.Routes
}

func (c *config) GetName() string {
	return c.Name
}
func (c *config) GetAuth() string {
	return c.Authentication
}

func (c *config) GetSigningMethod() string {
	return c.SigningMethod
}

func ReadConfigFile() *config {
	conf := config{}
	log.Info("Reading flurry config file.")
	source, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(source, &conf)
	if err != nil {
		log.Error(err)
	}

	log.Infof("Config file found. Continuing with configuration: %v.", conf.Name)
	return &conf
}

/// Route helpers
func (r *route) GetEndpoint() string {
	return r.Endpoint
}
func (r *route) GetHost() string {
	return r.Host
}
func (r *route) GetAuthenticated() bool {
	return r.Authenticated
}
func (r *route) GetMethods() []string {
	return r.Methods
}
