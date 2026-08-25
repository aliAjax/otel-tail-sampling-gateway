package observability

import (
	"encoding/json"
	"log"
	"os"
	"sync/atomic"
)

type Metrics struct {
	Ingested atomic.Uint64
	Kept     atomic.Uint64
	Dropped  atomic.Uint64
	Errors   atomic.Uint64
}

func Log(event string, fields map[string]any) {
	fields["event"] = event
	json.NewEncoder(os.Stdout).Encode(fields)
}

var Logger = log.New(os.Stdout, "", 0)
