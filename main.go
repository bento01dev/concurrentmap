package main

import (
	"hash/fnv"
	"log"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"
)

func hashing(key string) uint32 {
	h := fnv.New32()
	h.Write([]byte(key))
	return h.Sum32()
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/map", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second)
		m := NewConcurrentMap[string, string](hashing)
		for i := 0; i < 10000000; i++ {
			m.Set(strconv.Itoa(i), "value")
			if (i % 100000) == 0 {
				time.Sleep(500 * time.Millisecond)
			}
		}
		w.Write([]byte("lock and load.."))
		return
	}))
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/{action}", pprof.Index)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
