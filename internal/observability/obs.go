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
	out := make(map[string]any, len(fields)+1)
	for k, v := range fields {
		out[k] = v
	}
	out["event"] = event
	json.NewEncoder(os.Stdout).Encode(out)
}

var Logger = log.New(os.Stdout, "", 0)
