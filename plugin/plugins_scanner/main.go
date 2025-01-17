package main

import (
	"fmt"
	"log"
	"net/http"

	"sharp-lib/plugin/scanner"
)

var Users = []string{"admin", "manager", "tomcat"}
var Passwords = []string{"admin", "manager", "tomcat", "password"}

type TomcatChecker struct{}

func (c *TomcatChecker) Check(host string, port uint64) *scanner.Result {
	var (
		resp   *http.Response
		err    error
		url    string
		res    *scanner.Result
		client *http.Client
		req    *http.Request
	)

	log.Println("Checking for Tomcat  Manager...")
	res = new(scanner.Result)
	url = fmt.Sprintf("http://%s:%d/manager/html", host, port)
	if resp, err = http.Head(url); err != nil {
		log.Printf("Head request failed:%s\n", err)
		return res
	}

	log.Println("Host response to /manager/html request")

	if resp.StatusCode != http.StatusUnauthorized || resp.Header.Get("WWW-Authenticate") == "" {
		log.Println("Target doesn't appear to require basic auth.")
		return res
	}

	log.Println("Host requires authentication. Proceeding with password guessing...")
	client = new(http.Client)
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		log.Println("Unable to build Get request")
		return res
	}

	for _, user := range Users {
		for _, password := range Passwords {
			req.SetBasicAuth(user, password)
			if resp, err = client.Do(req); err != nil {
				log.Println("Unable to send Get request")
				continue
			}

			if resp.StatusCode == http.StatusOK {
				res.Vulnerable = true
				res.Details = fmt.Sprintf("Valid credentials found %s:%s", user, password)
				return res
			}
		}
	}
	return res
}

func New() scanner.Checker {
	return new(TomcatChecker)
}
