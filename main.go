package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Kingson4Wu/distributed_ratelimiter/limiter"
)

func main() {
	// Read startup parameters
	var (
		selfID    = flag.String("id", "", "Unique ID for this instance")
		globalQPS = flag.Float64("qps", 1000, "Global QPS limit")
		weight    = flag.Float64("weight", 1, "Weight for this instance")
		port      = flag.Int("port", 8080, "Listen port")
	)
	flag.Parse()
	if *selfID == "" {
		fmt.Println("--id is required")
		os.Exit(1)
	}

	nm := limiter.NewNodeManager(*selfID)
	nm.UpdateNodes([]limiter.Node{{ID: *selfID, Weight: *weight}})

	selfRate := nm.CalcSelfRate(*globalQPS)
	lim := limiter.NewLimiter(selfRate)

	api := &limiter.API{
		NodeMgr:   nm,
		Limiter:   lim,
		GlobalQPS: *globalQPS,
	}

	http.HandleFunc("/update_nodes", api.UpdateNodesHandler)

	log.Printf("Service started: id=%s, qps=%.2f, weight=%.2f, port=%d", *selfID, *globalQPS, *weight, *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
