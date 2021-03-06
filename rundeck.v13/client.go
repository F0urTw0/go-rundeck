package rundeck

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/jmcvetta/napping.v2"
	"net/http"
	"os"
)

type ClientConfig struct {
	BaseURL    string
	Token      string
	VerifySSL  bool
	Username   string
	Password   string
	AuthMethod string
}

type RundeckClient struct {
	Client *napping.Session
	Config *ClientConfig
}

func NewClient(config *ClientConfig) (c RundeckClient) {
	verifySSL := func() bool {
		if config.VerifySSL != true {
			return false
		} else {
			return true
		}
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: verifySSL()},
	}
	client := &http.Client{
		Transport: tr,
	}
	s := napping.Session{
		Client: client,
	}
	return RundeckClient{Client: &s, Config: config}
}

func clientConfigFrom(from string) (c *ClientConfig) {
	config := ClientConfig{}

	switch from {
	case "environment":
		if os.Getenv("RUNDECK_TOKEN") == "" {
			if os.Getenv("RUNDECK_USER") == "" && os.Getenv("RUNDECK_PASSWORD") == "" {
				fmt.Printf("You must set either RUNDECK_TOKEN or RUNDECK_USERNAME and RUNDECK_PASSWORD\n")
				os.Exit(1)
			} else {
				config.AuthMethod = "basic"
			}
		} else {
			config.AuthMethod = "token"
		}

		if os.Getenv("RUNDECK_URL") == "" {
			fmt.Printf("You must set the environment variable RUNDECK_URL\n")
			os.Exit(1)
		} else {
			config.BaseURL = os.Getenv("RUNDECK_URL")
		}
	}
	if config.AuthMethod == "token" {
		config.Token = os.Getenv("RUNDECK_TOKEN")
	} else {
		config.Username = os.Getenv("RUNDECK_USERNAME")
		config.Password = os.Getenv("RUNDECK_PASSWORD")
	}
	return &config
}

func NewClientFromEnv() (c RundeckClient) {
	config := clientConfigFrom("environment")
	client := NewClient(config)
	return client
}
