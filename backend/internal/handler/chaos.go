package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
)

var (
	failPercent atomic.Int64
	slowMs      atomic.Int64
)

func envInt(key string, def int64) int64 {
	if v, err := strconv.ParseInt(os.Getenv(key), 10, 64); err == nil {
		return v
	}
	return def
}

func InitChaos() {
	failPercent.Store(envInt("CHAOS_FAIL_PERCENT", 0))
	slowMs.Store(envInt("CHAOS_SLOW_MS", 0))
}

func ChaosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "use POST", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	if v := q.Get("fail"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n < 0 || n > 100 {
			http.Error(w, "fail must be 0..100", http.StatusBadRequest)
			return
		}
		failPercent.Store(n)
	}
	if v := q.Get("slow"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n < 0 {
			http.Error(w, "slow must be >= 0", http.StatusBadRequest)
			return
		}
		slowMs.Store(n)
	}
	fmt.Fprintf(w, "fail=%d%% slow=%dms\n", failPercent.Load(), slowMs.Load())
}
