package main

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"hash/fnv"
	"log"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"time"
)

func hashing(key string) uint32 {
	h := fnv.New32()
	h.Write([]byte(key))
	return h.Sum32()
}

var m *ConcurrentMap[string, string]
var cm = cmap.New[string]()

func main() {
	rand.Seed(time.Now().UnixNano())
	m = NewConcurrentMap[string, string](hashing)
	mux := http.NewServeMux()
	mux.Handle("/map", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		for i := 0; i < 10000; i++ {
			b := make([]rune, 5)
			for i := range b {
				b[i] = letters[rand.Intn(len(letters))]
			}
			m.Set(string(b), "value")
		}
		w.Write([]byte("lock and load.."))
		return
	}))
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
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
