package main

import (
	"fmt"
	"net"
	"sort"

	flag "github.com/spf13/pflag"
)

type config struct {
	hostname string
	ports    string
	workers  int
}

func worker(ports <-chan int, results chan<- int) {
	for p := range ports {
		addr := fmt.Sprintf("scanme.nmap.org:%d", p)
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

	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ {
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
