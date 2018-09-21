package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/giantswarm/micrologger"
	yaml "gopkg.in/yaml.v2"
)

const (
	// Port which listens the server
	Port = "80"
)

var (
	description = "Application run for health checking."
	gitCommit   = "n/a"
	name        = "liveness"
	source      = "https://github.com/giantswarm/liveness"
)

func main() {
	if (len(os.Args) > 1) && (os.Args[1] == "version") {
		d, err := yaml.Marshal(newVersionResponse())
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", d)

		return
	}

	var lHelp bool
	var sHelp bool
	flag.BoolVar(&lHelp, "help", false, "Print help usage.")
	flag.BoolVar(&sHelp, "h", false, "Print help usage.")
	flag.Parse()
	if lHelp || sHelp {
		flag.Usage()
		return
	}

	ctx := context.Background()
	mux := http.NewServeMux()
	logger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		panic(err)
	}

	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	}))

	server := &http.Server{
		Addr:    ":" + Port,
		Handler: mux,
	}

	go func() {
		logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("running server at http://0.0.0.0:%s", Port))
		err := server.ListenAndServe()
		if IsServerClosed(err) {
			// fall through
		} else if err != nil {
			panic(err)
		}
	}()

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	logger.LogCtx(ctx, "level", "debug", "message", "received termination signal")
	logger.LogCtx(ctx, "level", "debug", "message", "draining server connections")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	logger.LogCtx(ctx, "level", "debug", "message", "shutting down")
}
