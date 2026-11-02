package json

import (
	"bufio"
	"context"
	"encoding/json"
	"finchd/lib/api"
	"fmt"
	"iter"
	"log"
	"os"
)

// ctx optional
func WriteData(output string, jobs iter.Seq2[*api.JobEntry, error], count *int, ctx *context.Context) (err error) {
	var buf *bufio.Writer
	if output != "stdout" {
		f, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
		if err != nil {
			log.Fatal(fmt.Errorf("while opening %s: %w", output, err))
		}
		buf = bufio.NewWriter(f)
		defer f.Close() //nolint
	} else {
		buf = bufio.NewWriter(os.Stdout)
	}
	defer buf.Flush() // nolint
	first := true
	if _, err = buf.WriteString("["); err != nil {
		return
	}
	for data, err := range jobs {
		if err != nil {
			break
		}
		if ctx != nil && (*ctx).Err() != nil {
			return (*ctx).Err()
		}
		marshalled, err := json.Marshal(data)
		if err != nil {
			break
		}
		if first {
			first = false
		} else if _, err := buf.WriteString(","); err != nil {
			break
		}
		if _, err := buf.Write(marshalled); err != nil {
			break
		}
		*count++
		if *count%50 == 0 {
			log.Printf("recv %d", *count)
		}
	}
	_, err = buf.WriteString("]")
	return
}
