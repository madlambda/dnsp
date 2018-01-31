package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	wd, err := os.Getwd()
	abortonerr(err, "getting working dir")

	var dir string
	var domain string

	flag.StringVar(&dir, "dir", wd, "dir with files that will be provided by fake server")
	flag.StringVar(&domain, "domain", "", "domain name that fill be faked")

	flag.Parse()

	if domain == "" {
		flag.Usage()
		abortonerr(errors.New("domain argument is obligatory"), "checking obligatory args")
	}

	cleanup := poisonDNS(domain)
	defer cleanup()

	server := startFakeServer(dir)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	<-signals
	fmt.Println("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	server.Shutdown(ctx)
	fmt.Println("Cleaning up DNS entries")
}

func poisonDNS(domain string) func() {
	// FIXME: not portable at all =P
	const hostsfile = "/etc/hosts"

	stat, err := os.Stat(hostsfile)
	abortonerr(err, "stating hosts file")
	original, err := ioutil.ReadFile(hostsfile)
	abortonerr(err, "reading hosts file")

	newhosts := fmt.Sprintf("%s\n127.0.0.1  %s\n", string(original), domain)
	err = ioutil.WriteFile(hostsfile, []byte(newhosts), stat.Mode())
	abortonerr(err, "writing new tampered hosts file")

	return func() {
		err := ioutil.WriteFile(hostsfile, original, stat.Mode())
		abortonerr(err, "recovering original hosts file")
	}
}

func startFakeServer(dir string) *http.Server {
	handler := http.FileServer(http.Dir(dir))
	server := &http.Server{Addr: "", Handler: handler}
	go func() {
		server.ListenAndServe()
	}()
	return server
}

func abortonerr(err error, detail string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error[%s] %s\n", err, detail)
		os.Exit(1)
	}
}
