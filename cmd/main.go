package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/keratomalacian/dumbbin/internal/routes"
)

var cmdlineArgs = struct {
	address     string
	binPath     string
	logger      bool
	requestSize int64

	ratelimit    bool
	requestLimit int
	timeLimit    time.Duration
}{}

func main() {
	log.Println("starting dumbbin...")

	// parsing command line arguments
	flag.StringVar(&cmdlineArgs.address, "address", ":3232", "sets the address in which the server will run")
	flag.StringVar(&cmdlineArgs.binPath, "binpath", "bins", "sets the directory where bins are stored")
	flag.BoolVar(&cmdlineArgs.logger, "logger", true, "enables or disables the request logger")
	flag.Int64Var(&cmdlineArgs.requestSize, "requestsize", 1, "sets the maximum number of megabytes that will be read from the request body")

	flag.BoolVar(&cmdlineArgs.ratelimit, "ratelimit", true, "enables or disables ratelimiting")
	flag.IntVar(&cmdlineArgs.requestLimit, "requestlimit", 100, "sets the limit of requests that can be done in the 'timelimit' time window")
	flag.DurationVar(&cmdlineArgs.timeLimit, "timelimit", time.Hour, "sets the time window for 'requestlimit'. eg: 1h30s")

	flag.Parse()

	// creating directory for storing bins
	binPath, err := filepath.Abs(cmdlineArgs.binPath)
	if err != nil {
		log.Fatalf("could not get absolute path for '%s': %s\n", cmdlineArgs.binPath, err)
	}

	err = os.MkdirAll(binPath, 0777)
	if err != nil {
		log.Fatalf("error creating folder for storing bins at '%s': %s\n", binPath, err)
	}
	log.Printf("directory for storing bins set to '%s'", binPath)

	// creating chi router
	var router = chi.NewRouter()

	// setting up muddlewares from chi
	if cmdlineArgs.logger {
		router.Use(middleware.Logger)
	}

	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.RedirectSlashes)
	router.Use(middleware.CleanPath)
	router.Use(middleware.RequestSize(cmdlineArgs.requestSize * 1048576)) // will read a maximum of 'requestSize' megabytes from the request body

	if cmdlineArgs.ratelimit {
		router.Use(httprate.Limit(cmdlineArgs.requestLimit, cmdlineArgs.timeLimit, httprate.WithKeyByIP(),
			httprate.WithLimitHandler(routes.RateLimited()))) // will ratelimit IP's to 'requestLimit' requests per 'timeLimit', if enabled.
	}

	// setting up the route handlers
	router.Get("/", routes.Root())
	router.Post("/", routes.CreateBin(binPath))
	router.Get("/{binID}", routes.GetBin(binPath))

	// running the server
	log.Println("running dumbbin on", cmdlineArgs.address)
	if err := http.ListenAndServe(cmdlineArgs.address, router); err != nil {
		log.Fatalln("error while starting the web server:", err)
	}
}
