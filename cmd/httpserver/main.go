package main

import (
	"crypto/sha256"
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
	if !strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/html") {
		res.WriteStatusLine(response.StatusCodeNotFound())
		headers := response.GetDefaultHeaders(0)
		res.WriteHeaders(headers)
		return
	}

	fullAddress := "https://httpbin.org/html"

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
	resHeaders := make(headers.Headers)
	for header, values := range httpBinResp.Header {
		if header == "Content-Length" {
			continue
		}
		for _, value := range values {
			resHeaders.Set(header, value)
		}
	}
	resHeaders.Set("Transfer-Encoding", "chunked")
	resHeaders.Set("Trailer", "X-Content-SHA256")
	resHeaders.Set("Trailer", "X-Content-Length")
	res.WriteHeaders(resHeaders)
	fmt.Printf("%#v\n", resHeaders)
	buf := make([]byte, 32)

	var resBody []byte
	var i int
	for ; ; i++ {
		bytesRead, errRead := httpBinResp.Body.Read(buf)

		if errRead != nil && errRead != io.EOF {
			fmt.Fprintf(os.Stderr, "error reading response Body of %s\n", fullAddress)
			return
		}
		resBody = append(resBody, buf[:bytesRead]...)
		_, errWrite := res.WriteChunkedBody(buf[:bytesRead])
		if errWrite != nil {
			fmt.Fprintf(os.Stderr, "error writing response Body of %s: %s\n", fullAddress, errWrite)
			return
		}
		if bytesRead == 0 && errRead == io.EOF {
			break
		}
	}
	fmt.Printf("%d\n", i)
	trailers := make(headers.Headers)
	trailers.Set("X-Content-Sha256", fmt.Sprintf("%x", sha256.Sum256(resBody)))
	trailers.Set("X-Content-Length", strconv.Itoa(len(resBody)))
	fmt.Printf("%+v\n", trailers)
	err = res.WriteTrailers(trailers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "errr: %s", err)
		return
	}

	_, err = res.WriteDone()
	if err != nil {
		fmt.Fprintf(os.Stderr, "errr: %s", err)
		return
	}
}

func videoHandler(res *response.Writer, req *request.Request) {
	if !strings.HasPrefix(req.RequestLine.RequestTarget, "/video") {
		res.WriteStatusLine(response.StatusCodeNotFound())
		headers := response.GetDefaultHeaders(0)
		res.WriteHeaders(headers)
		return
	}
	const fileName = "../../assets/vim.mp4"
	video, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading the file %s: %s", fileName, err)
		return
	}

	res.WriteStatusLine(response.StatusCodeOK())
	resHeaders := headers.NewHeaders()
	resHeaders.Set("Connection", "close")
	resHeaders.Set("Content-Length", fmt.Sprintf("%d", len(video)))
	resHeaders.Set("Content-Type", "video/mp4")
	res.WriteHeaders(resHeaders)
	res.WriteBody(video)
}

func main() {
	server, err := server.Serve(port, videoHandler)
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
