package main

import (
	"fmt"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

type config struct {
	hostname string
	ports    string
	workers  int
}

func worker(ports <-chan int, results chan<- int, hostname string) {
	for p := range ports {
		addr := fmt.Sprintf("%s:%d", hostname, p)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	var cfg config
	flag.StringVarP(&cfg.hostname, "hostname", "h", "", "e.g. scanme.nmap.org")
	flag.StringVarP(&cfg.ports, "ports", "p", "", "e.g. 1-1024")
	flag.IntVarP(&cfg.workers, "workers", "w", 100, "number of workers in the worker pool")
	flag.Parse()

	if cfg.hostname == "" || cfg.ports == "" {
		log.Fatalln("Please specify additional flags")
	}

	firstPort, err := strconv.Atoi(strings.Split(cfg.ports, "-")[0])
	if err != nil {
		log.Fatalln("Wrong flag format")
	}

	lastPort, err := strconv.Atoi(strings.Split(cfg.ports, "-")[1])
	if err != nil {
		log.Fatalln("Wrong flag format")
	}

	ports := make(chan int, cfg.workers)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results, cfg.hostname)
	}

	go func() {
		for i := firstPort; i <= lastPort; i++ {
			ports <- i
		}
	}()

	for i := firstPort; i <= lastPort; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
