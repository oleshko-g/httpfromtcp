package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/oleshko-g/httpfromtcp/internal/headers"
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

func convertStdStatusCode(stdStatusCode int) response.StatusCode {
	return response.StatusCode([]rune(strconv.Itoa(stdStatusCode)))
}

func httpBinStreamHandler(res *response.Writer, req *request.Request) {
	if !strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/stream/") {
		res.WriteStatusLine(response.StatusCodeNotFound())
		headers := response.GetDefaultHeaders(0)
		res.WriteHeaders(headers)
		return
	}

	numberOfResponces := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/stream/")
	const httpBinStreamAddress = "https://httpbin.org/stream/%s"
	fullAddress := fmt.Sprintf(httpBinStreamAddress, numberOfResponces)

	httpBinResp, err := http.Get(fullAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error making request to %s", fullAddress)
		return
	}

	if httpBinResp.StatusCode < 200 || httpBinResp.StatusCode >= 300 {
		statusCode := convertStdStatusCode(httpBinResp.StatusCode)
		res.WriteStatusLine(statusCode)
		headers := response.GetDefaultHeaders(0)
		res.WriteHeaders(headers)
		return
	}

	res.WriteStatusLine(response.StatusCodeOK())
	headers := make(headers.Headers)
	for header, values := range httpBinResp.Header {
		if header == "Content-Length" {
			continue
		}
		for _, value := range values {
			headers.Set(header, value)
		}
	}
	headers.Set("Transfer-Encoding", "chunked")
	res.WriteHeaders(headers)
	fmt.Printf("%#v\n", headers)
	buf := make([]byte, 32)
	for i := 0; ; i++ {
		bytesRead, errRead := httpBinResp.Body.Read(buf)
		if bytesRead == 0 && errRead == io.EOF {
			break
		}

		if errRead != nil {
			fmt.Fprintf(os.Stderr, "error reading response Body of %s\n", fullAddress)
			return
		}

		_, errWrite := res.WriteChunkedBody(buf)
		if errWrite != nil {
			fmt.Fprintf(os.Stderr, "error writing response Body of %s: %s\n", fullAddress, errWrite)
			return
		}
	}
	res.WriteChunkedBodyDone()
}

func main() {
	server, err := server.Serve(port, httpBinStreamHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
