package main

import (
	"context"
	"errors"
	"finchd/lib/api"
	"finchd/lib/db"
	"finchd/lib/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var AVAILABLE_TYPES = map[string]struct{}{
	"json":   {},
	"sqlite": {},
}

var OUTPUT string
var COUNT int
var MAJOR string
var TYPE string
var USONLY bool

func init() {
	if OUTPUT = os.Getenv("OUTPUT_FILE"); OUTPUT == "" {
		OUTPUT = "stdout"
	}

	if os.Getenv("USONLY") != "0" {
		USONLY = true
	}

	if cnt := os.Getenv("COUNT"); cnt != "" {
		tmp, err := strconv.Atoi(cnt)
		if err != nil {
			log.Fatal(err)
		} else if tmp <= 0 && tmp != -1 {
			log.Fatal("Invalid count: ", tmp)
		}
		COUNT = tmp
	} else {
		COUNT = 50
	}
	TYPE = strings.ToLower(os.Getenv("TYPE"))
	if TYPE == "" {
		TYPE = "json"
	} else if _, ok := AVAILABLE_TYPES[TYPE]; !ok {
		log.Fatalf("Unavailable filetype: %s", TYPE)
	}

	MAJOR = os.Getenv("MAJOR")

	log.Printf("Output to: %s", OUTPUT)
	log.Printf("Output type: %s", TYPE)
	log.Printf("Max Fetched: %d", COUNT)
	log.Printf("Major: '%s'", MAJOR)
}

func main() {
	// ctrl c
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var err error
	counter := 0
	jobs := api.Jobs(MAJOR, COUNT, USONLY)
	switch TYPE {
	case "json":
		err = json.WriteData(OUTPUT, jobs, &counter, &ctx)
	case "sqlite":
		err = db.WriteData(OUTPUT, jobs, &counter, &ctx)
	}
	if errors.Is(err, context.Canceled) {
		log.Println("Interrupted")
		err = nil
	}
	log.Printf("Total received: %d", counter)

	if err != nil {
		panic(fmt.Errorf("GET Jobs: %w", err))
	}
}
