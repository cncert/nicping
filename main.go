package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

var (
	hosts       string
	packetCount int
	concurrent  = 10
	wg          sync.WaitGroup
	tokens      = make(chan struct{}, concurrent)
)

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func runPing(host string, packetCount int) {
	// fmt.Println(host)
	defer wg.Done()
	m := make(map[string]interface{})
	pinger, err := ping.NewPinger(host)
	pinger.SetPrivileged(true)
	pinger.Timeout = 3 * time.Second // 程序最长ping时间，如果时间设置为3秒，但是却设置想要发送20个包，那么程序只能发送三个包
	if err != nil {
		panic(err)
	}
	pinger.Count = packetCount //超时时间内最多发送包数量
	err = pinger.Run()         // Blocks until finished.
	if err != nil {
		panic(err)
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	//转换为map
	elem := reflect.ValueOf(stats).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		m[relType.Field(i).Name] = elem.Field(i).Interface()
	}
	// fmt.Println(m)
	fmt.Printf("%-20s	packetLoss: %3.0f%%,	RTT: %s\n", m["Addr"], m["PacketLoss"], m["AvgRtt"])
	<-tokens
}

func mutiPing(host string, packetCount int) {
	wg.Add(1)
	tokens <- struct{}{}
	go runPing(host, packetCount)
}

func init() {
	flag.StringVar(&hosts, "hosts", "hosts.txt", "host's ip or domain in this file, line by line. if not exist file, can replace by ip ")
	flag.IntVar(&packetCount, "p", 2, "number of send packets per ping")
	flag.IntVar(&concurrent, "c", 10, "concurrent number")
}

func main() {
	flag.Parse()
	if fileExist(hosts) {
		fmt.Printf("read ip from %s\n", hosts)
		f, err := os.Open(hosts)
		if err != nil {
			fmt.Println("read file fail", err)
			return
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			host := scanner.Text()
			mutiPing(host, packetCount)
		}
	} else {
		fmt.Printf("read ip %s from command line args\n", hosts)
		mutiPing(hosts, packetCount)
	}
	wg.Wait()

}
