package testutils

import (
	"net"
	"os"
	"os/exec"
	"strings"
)

// RunNucleiAndGetResults returns a list of results for a template
func RunNucleiAndGetResults(template, url string, debug bool, extra ...string) ([]string, error) {
	cmd := exec.Command("./nuclei", "-t", template, "-target", url)
	if debug {
		cmd = exec.Command("./nuclei", "-t", template, "-target", url, "-debug")
		cmd.Stderr = os.Stderr
	}
	cmd.Args = append(cmd.Args, extra...)

	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	parts := []string{}
	items := strings.Split(string(data), "\n")
	for _, i := range items {
		if i != "" {
			parts = append(parts, i)
		}
	}
	return parts, nil
}

// RunNucleiWorkflowAndGetResults returns a list of results for a workflow
func RunNucleiWorkflowAndGetResults(template, url string, debug bool, extra ...string) ([]string, error) {
	cmd := exec.Command("./nuclei", "-w", template, "-target", url)
	if debug {
		cmd = exec.Command("./nuclei", "-w", template, "-target", url, "-debug")
		cmd.Stderr = os.Stderr
	}
	cmd.Args = append(cmd.Args, extra...)

	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	parts := []string{}
	items := strings.Split(string(data), "\n")
	for _, i := range items {
		if i != "" {
			parts = append(parts, i)
		}
	}
	return parts, nil
}

// TestCase is a single integration test case
type TestCase interface {
	// Execute executes a test case and returns any errors if occurred
	Execute(filePath string) error
}

// TCPServer creates a new tcp server that returns a response
type TCPServer struct {
	URL      string
	listener net.Listener
}

// NewTCPServer creates a new TCP server from a handler
func NewTCPServer(handler func(conn net.Conn)) *TCPServer {
	server := &TCPServer{}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	server.URL = l.Addr().String()
	server.listener = l

	go func() {
		for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				continue
			}
			// Handle connections in a new goroutine.
			go handler(conn)
		}
	}()
	return server
}

// Close closes the TCP server
func (s *TCPServer) Close() {
	s.listener.Close()
}
