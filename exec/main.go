package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {
	errCh := make(chan error)
	go LaunchOgpApp(errCh)
	go LaunchCaddy(errCh)
	e := <-errCh
	log.Fatalf("%v\n", e)
}

func LaunchOgpApp(e chan<- error) {
	cmd := exec.Command("/app/ogp-app", "-c", "/app/prd.toml")
	stdout, stderr, err := pipe(cmd)
	if err != nil {
		e <- err
		return
	}
	if err := cmd.Start(); err != nil {
		e <- err
		return
	}
	go forwardOutput(stdout)
	go forwardOutput(stderr)
	if err := cmd.Wait(); err != nil {
		e <- err
		return
	}
}

func LaunchCaddy(e chan<- error) {
	cmd := exec.Command("/app/caddy", "run", "--config", "/app/Caddyfile", "--adapter", "caddyfile")
	stdout, stderr, err := pipe(cmd)
	if err != nil {
		e <- err
		return
	}
	if err := cmd.Start(); err != nil {
		e <- err
		return
	}
	go forwardOutput(stdout)
	go forwardOutput(stderr)
	if err := cmd.Wait(); err != nil {
		e <- err
		return
	}
}

func pipe(cmd *exec.Cmd) (io.ReadCloser, io.ReadCloser, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}
	return stdout, stderr, nil
}

func forwardOutput(rc io.ReadCloser) {
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		l := scanner.Text()
		fmt.Println(l)
	}
}
