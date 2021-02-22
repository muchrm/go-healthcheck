package lineauth

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/muchrm/go-healthcheck/config"
)

// CodeResponse represents the code received by the local server's callback handler.
type CodeResponse struct {
	Code  string
	State string
}

// bindLocalServer initializes a LocalServer that will listen on a randomly available TCP port.
func bindLocalServer() (*localServer, error) {
	hostAddr, err := config.GetString(config.AuthHostAddr)
	if err != nil {
		return nil, fmt.Errorf("BindLocalServer error %w", err)
	}
	listener, err := net.Listen("tcp", hostAddr)
	if err != nil {
		return nil, fmt.Errorf("BindLocalServer error %w", err)
	}

	return &localServer{
		listener:   listener,
		resultChan: make(chan CodeResponse, 1),
	}, nil
}

type localServer struct {
	CallbackPath     string
	WriteSuccessHTML func(w io.Writer)

	resultChan chan (CodeResponse)
	listener   net.Listener
}

func (s *localServer) Port() int {
	return s.listener.Addr().(*net.TCPAddr).Port
}

func (s *localServer) Close() error {
	return s.listener.Close()
}

func (s *localServer) Serve() error {
	return http.Serve(s.listener, s)
}

func (s *localServer) WaitForCode() CodeResponse {
	return <-s.resultChan
}

func (s *localServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.CallbackPath != "" && r.URL.Path != s.CallbackPath {
		w.WriteHeader(404)
		return
	}
	defer func() {
		_ = s.Close()
	}()

	params := r.URL.Query()
	s.resultChan <- CodeResponse{
		Code:  params.Get("code"),
		State: params.Get("state"),
	}

	w.Header().Add("content-type", "text/html")
	successHTML(w)
}

func successHTML(w io.Writer) {
	fmt.Fprintf(w, `<html xmlns:th="http://www.thymeleaf.org">
	<head>
	  <meta http-equiv='Content-type' content='text/html; charset=utf-8' />
	  <meta name="viewport" content="width=device-width, initial-scale=1" />
	  <link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" />
	  <link rel="stylesheet" href="css/line-login.css" />
	  <script src="https://code.jquery.com/jquery-2.2.4.min.js" integrity="sha256-BbhdlvQf/xTY9gja0Dq3HiwQF8LaCRTXxZKRutelT44=" crossorigin="anonymous"></script>
	  <script src="js/success.js"></script>
	  <title>LINE Web Login Success</title>
	</head>
	<body>
	  <div class="container">
		<div class="row">
		  <div class="col-md-4 col-md-offset-4">
			<div class="area">
			  <div class="center-block profile-margin">
				<p>You may now close this page and return to the client app</p>
			  </div>
			</div>
		  </div>
		</div>
	  </div>
	</body>
  </html>`)
}
