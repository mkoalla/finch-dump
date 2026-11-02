package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"finchd/lib/api"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var OUTPUT string
var COUNT int
var MAJOR string

func init() {
	if OUTPUT = os.Getenv("OUTPUT_FILE"); OUTPUT == "" {
		OUTPUT = "stdout"
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
	MAJOR = os.Getenv("MAJOR")
	log.Printf("Output to: %s", OUTPUT)
	log.Printf("Max Fetched: %d", COUNT)
	log.Printf("Major: '%s'", MAJOR)
}

func write_data(ctx context.Context, buf *bufio.Writer, count *int) (err error) {
	jobs := api.Jobs(MAJOR, COUNT)
	first := true
	for data, err := range jobs {
		if err != nil {
			break
		}
		if ctx.Err() != nil {
			return ctx.Err() // stops the iterator; yield returns false
		}
		marshalled, err := json.Marshal(data)
		if err != nil {
			break
		}
		if first {
			first = false
		} else {
			buf.WriteString(",")
		}
		buf.Write(marshalled)
		*count++
		if *count%50 == 0 {
			log.Printf("recv %d", *count)
		}
	}
	return nil
}

func main() {
	var buf *bufio.Writer
	if OUTPUT != "stdout" {
		f, err := os.OpenFile(OUTPUT, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
		if err != nil {
			log.Fatal(fmt.Errorf("While opening %s: %w", OUTPUT, err))
		}
		buf = bufio.NewWriter(f)
		defer f.Close()
	} else {
		buf = bufio.NewWriter(os.Stdout)
	}
	defer buf.Flush()

	// ctrl c
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	count := 0
	buf.WriteString("[")
	err := write_data(ctx, buf, &count)
	buf.WriteString("]")
	if errors.Is(err, context.Canceled) {
		log.Println("Interrupted")
		err = nil
	}
	log.Printf("Total received: %d", count)

	if err != nil {
		log.Fatal(fmt.Errorf("GET Jobs: err"))
	}
}
