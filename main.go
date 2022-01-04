package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

// main is called when a new lambda starts, so don't
// expect to have something done for every query here.
func main() {
	http.HandleFunc("/", Handler)

	lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
}
