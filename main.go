package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

func ipToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func uint32ToIP(n uint32) net.IP {
	return net.IPv4(byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
}

// 检查是否是私有地址或保留地址
func isPrivateIP(ip uint32) bool {
	// 10.0.0.0/8
	if ip>>24 == 10 {
		return true
	}
	// 172.16.0.0/12
	if ip>>20 == 0xAC1 {
		return true
	}
	// 192.168.0.0/16
	if ip>>16 == 0xC0A8 {
		return true
	}
	// 127.0.0.0/8
	if ip>>24 == 127 {
		return true
	}
	// 224.0.0.0/4 (multicast)
	if ip>>28 == 0xE {
		return true
	}
	// 240.0.0.0/4 (reserved)
	if ip>>28 == 0xF {
		return true
	}
	return false
}

func ipRangeWorker(start, end, step uint32, ch chan<- string, wg *sync.WaitGroup, count *int64, dbFile string) {
	defer wg.Done()

	searcher, err := xdb.NewWithFileOnly(dbFile)
	if err != nil {
		log.Printf("failed to open searcher: %v", err)
		return
	}
	defer searcher.Close()

	for ip := start; ip < end; ip += step {
		if isPrivateIP(ip) {
			continue
		}
		region, err := searcher.Search(ip)
		if err != nil {
			log.Println("search error:", err)
			continue
		}
		ch <- region
		atomic.AddInt64(count, 1)
	}
}

func main() {
	const (
		totalIPs uint32 = 1<<32 - 1
		step     uint32 = 0x10000
		workers  int    = 100
		dbFile   string = "ip2region.xdb"
	)

	ch := make(chan string, 1000)
	var wg sync.WaitGroup
	var count int64

	rangeSize := totalIPs / uint32(workers)

	for i := 0; i < workers; i++ {
		start := uint32(i) * rangeSize
		end := start + rangeSize
		if i == workers-1 {
			end = totalIPs
		}
		wg.Add(1)
		go ipRangeWorker(start, end, step, ch, &wg, &count, dbFile)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	provinces := make(map[string]struct{})
	for region := range ch {
		provinces[region] = struct{}{}
	}

	fmt.Printf("Processed %d IPs, found %d unique regions:\n", count, len(provinces))
	for province := range provinces {
		fmt.Println(province)
	}
}
