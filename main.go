package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

var chiLambda *chiadapter.ChiLambda

// todo - could put slack client creation in init function "your handler may declare an init function that is executed when your handler is loaded."

// handler is the function called when the lambda is invoked
func handler(ctx context.Context, req events.APIGatewayProxyRequest) (string, error) { //(events.APIGatewayProxyResponse, error)
	return fmt.Sprintf("hai %s", req.RequestContext.RequestID), nil
}

// main is called when a new lambda starts, so don't
// expect to have something done for every query here.
func main() {
	r := chi.NewRouter()
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		// q := r.URL.Query()
		// q.Get(key string)
		// get query params
		_ = render.Render(w, r, &apiResponse{
			Status:      http.StatusOK,
			URL:         r.URL.String(),
		})
	})

	// start the lambda with a context
	lambda.StartWithContext(context.Background(), handler)
}

// apiResponse is the response to the API. 
// TODO decide what format the codeowners response should be - probably string or error?
type apiResponse struct {
	Status      int    `json:"status_code,omitempty"`
	URL         string `json:"url,omitempty"`
}

// Render is used by go-chi-render to render the JSON response.
func (a apiResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, a.Status)
	return nil
}
