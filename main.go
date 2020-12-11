package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"net"
	"time"
	"github.com/cenk/backoff"
)

type ServerInfo struct {
	State string `json:"state"`
}

func main() {
	// Should be in format `http://127.0.0.1:9010`
	host, ok := os.LookupEnv("CLOUDSQL_PROXY_API")
	if ok && os.Getenv("START_WITHOUT_CLOUDSQL_PROXY_API") != "true" {
		block(host)
	}

	if len(os.Args) < 2 {
		return
	}

	binary, err := exec.LookPath(os.Args[1])
	if err != nil {
		panic(err)
	}

	var proc *os.Process

	// Pass signals to the child process
	go func() {
		stop := make(chan os.Signal, 2)
		signal.Notify(stop)
		for sig := range stop {
			if proc != nil {
				proc.Signal(sig)
			} else {
				// Signal received before the process even started. Let's just exit.
				os.Exit(1)
			}
		}
	}()

	proc, err = os.StartProcess(binary, os.Args[1:], &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})
	if err != nil {
		panic(err)
	}

	state, err := proc.Wait()
	if err != nil {
		panic(err)
	}

	exitCode := state.ExitCode()

	switch {
	case !ok:
		// We don't have an CLOUDSQL_PROXY_API env var, do nothing
	case !strings.Contains(host, "127.0.0.1") && !strings.Contains(host, "localhost"):
		// Cloud SQL Proxy is not local; do nothing
	}

	os.Exit(exitCode)
}

func block(host string) {
	if os.Getenv("START_WITHOUT_CLOUDSQL_PROXY_API") == "true" {
		return
	}

	b := backoff.NewExponentialBackOff()
	// We wait forever for Cloud SQL Proxy to start. In practice k8s will kill the pod if we take too long.
	b.MaxElapsedTime = 0



	_ = backoff.Retry(func() error {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", host, timeout)

		if conn != nil {
				defer conn.Close()
				fmt.Println("Opened", host)
		} else {
				fmt.Println("Cloud SQL Proxy not live yet")
		}

		if err != nil {
			return errors.New("Cloud SQL Proxy not live yet")
		}

		return nil
	}, b)
}

func raw_connect(host string, port string) {

    timeout := time.Second
    conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
    if err != nil {
        fmt.Println("Connecting error:", err)
    }
    if conn != nil {
        defer conn.Close()
        fmt.Println("Opened", net.JoinHostPort(host, port))
    }
}
