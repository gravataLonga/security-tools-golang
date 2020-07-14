package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sort"
)

var workerQuantity = flag.Int("worker",20, "How many worker do you want")
var domain = flag.String("domain", "scanme.nmap.org", "Which domain you want scan")

func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%v:%d", *domain, p)
		log.Println(address)
		conn, err := net.Dial("tcp", address)

		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func PortScan() {
	flag.Parse()
	ports := make(chan int, *workerQuantity)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 1; i <= 65535; i++ {
			log.Println("Scan Port")
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