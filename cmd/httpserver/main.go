package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/oleshko-g/httpfromtcp/internal/request"
	"github.com/oleshko-g/httpfromtcp/internal/response"
	"github.com/oleshko-g/httpfromtcp/internal/server"
)

const port = 42069

func myProblemYourProblem(w io.Writer, r *request.Request) *server.HandlerError {
	switch {
	case r.RequestLine.RequestTarget == "/yourproblem":
		return &server.HandlerError{
			StatusCode: response.StatusCodeBadRequest(),
			Message:    "Your problem is not my problem\n",
		}
	case r.RequestLine.RequestTarget == "/myproblem":
		return &server.HandlerError{
			StatusCode: response.StatusCodeInternalServerError(),
			Message:    "Woopsie, my bad\n",
		}

	default:
		responseBody := "All good, frfr\n"
		w.Write([]byte(responseBody))
		return nil
	}
}

func myProblemYourProblemV2(res *response.Writer, req *request.Request) {
	switch {
	case req.RequestLine.RequestTarget == "/yourproblem":
		res.WriteStatusLine(response.StatusCodeBadRequest())
		responseBody := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>`)
		headers := response.GetDefaultHeaders(len(responseBody))
		headers.Set("Content-Type", "text/html")
		res.WriteHeaders(headers)
		res.WriteBody(responseBody)
	case req.RequestLine.RequestTarget == "/myproblem":
		res.WriteStatusLine(response.StatusCodeInternalServerError())
		responseBody := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>`)
		headers := response.GetDefaultHeaders(len(responseBody))
		headers.Set("Content-Type", "text/html")
		res.WriteHeaders(headers)
		res.WriteBody(responseBody)
	default:
		res.WriteStatusLine(response.StatusCodeOK())
		responseBody := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>`)
		headers := response.GetDefaultHeaders(len(responseBody))
		headers.Set("Content-Type", "text/html")
		res.WriteHeaders(headers)
		res.WriteBody(responseBody)
	}
}

func main() {
	server, err := server.Serve(port, myProblemYourProblemV2)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
