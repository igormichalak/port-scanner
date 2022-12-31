package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

var (
	ErrPortClosedOrFiltered = errors.New("the port is closed or filtered")
)

func main() {
	var wg sync.WaitGroup
	for i := 1; i <= 1024; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			res := make(chan error)
			//go scanPort("scanme.nmap.org", j, res)
			go scanPort("127.0.0.1", j, res)
			if <-res == nil {
				fmt.Printf("%d open\n", j)
			}
		}(i)
	}
	wg.Wait()
}

func scanPort(host string, port int, out chan error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		out <- ErrPortClosedOrFiltered
		return
	}
	conn.Close()
	out <- nil
}
