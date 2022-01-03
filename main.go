package main

import (
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

// todo - could put slack client creation in init function "your handler may declare an init function that is executed when your handler is loaded."

// TODO decide what format the codeowners response should be - probably string or error?

// handler is the function called when the lambda is invoked
// handler function signature (context, TIn where TIn can be unmarshalled by Go)
// return value must be either error or TOut where TOut can marshalled

// main is called when a new lambda starts, so don't
// expect to have something done for every query here.
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Test")
		// w.Write([]byte(`test`))
	})

	lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
}
