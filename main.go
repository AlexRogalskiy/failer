package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	failRandom = "random"

	// given success rate 80%, succeed 80 times, then fail 20 times
	failContiguous = "contiguous"

	// given success rate 80%, succeed 4 times, then fail 1 time
	failEvenly = "evenly"
)

type handler struct {
	sr       int
	dist     string
	requests int
	sync.Mutex

	// specific to failEvenly
	loopSuccesses int
	loopTotal     int
}

func gcd(a, b int) int {
	c := a % b
	for c > 0 {
		a = b
		b = c
		c = a % b
	}
	return b
}

func mkHandler(sr int, dist string) *handler {
	h := &handler{
		sr:   sr,
		dist: dist,
	}

	if dist == failEvenly {
		g := gcd(sr, 100)
		h.loopSuccesses = sr / g
		h.loopTotal = 100 / g
	}

	return h
}

func (h *handler) handle(w http.ResponseWriter, r *http.Request) {
	h.Lock()
	defer h.Unlock()

	success := false
	msg := ""

	switch h.dist {
	case failRandom:
		r := rand.Float32()
		success = int(r*100) < h.sr
		msg = fmt.Sprintf("%s [%f]", h.dist, r)
	case failContiguous:
		num := h.requests % 100
		success = num < h.sr
		msg = fmt.Sprintf("%s [%d]", h.dist, num)
	case failEvenly:
		num := h.requests % h.loopTotal
		success = num < h.loopSuccesses
		msg = fmt.Sprintf("%s [%d/%d] [%d]", h.dist, h.loopSuccesses, h.loopTotal, num)
	}

	h.requests++

	if success {
		fmt.Fprintf(w, "SR: [%d%%] success: %s\n", h.sr, msg)
	} else {
		http.Error(w, fmt.Sprintf("SR [%d%%] fail: %s", h.sr, msg), http.StatusInternalServerError)
	}
}

func main() {
	addr := flag.String("addr", ":8080", "address to serve on")
	sr := flag.Int("success-rate", 100, "server succcess rate percentage [0,100]")
	dist := flag.String("distribution", failRandom,
		fmt.Sprintf("failure distribution, must be one of: %s, %s, %s", failRandom, failContiguous, failEvenly),
	)

	flag.Parse()

	if *sr < 0 || *sr > 100 {
		log.Fatalf("invalid success-rate: %d", *sr)
	}

	if *dist != failRandom && *dist != failContiguous && *dist != failEvenly {
		log.Fatalf("invalid distribution: %s", *dist)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	http.HandleFunc("/", mkHandler(*sr, *dist).handle)

	go func() {
		fmt.Printf("listening on %s\n", *addr)
		err := http.ListenAndServe(*addr, nil)
		if err != nil {
			log.Fatalf("failed to listen on %s: %s", *addr, err)
		}
	}()

	<-stop
}
