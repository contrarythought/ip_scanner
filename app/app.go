package app

import (
	"database/sql"
	"math/rand"
	"net"
	"runtime"
	"sync"
	"time"
)

func genRandInt(randIntCh chan<- uint32) {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	randIntCh <- r.Uint32()
}

/*
	ports
	20 and 21: FTP
	22: SSH'
	23: Telnet
	53: DNS
	80: HTTP
	123: network time
	179: BGS
	443: HTTPS
	500: ISAKMP
	587: SMTP w/ encryption
	3389: remote desktop
*/

const (
	NUM_PORTS = 12
)

// TODO
func scanPorts(ip string) ([]string, error) {
	var wg sync.WaitGroup
	var opened_ports []string
	scan := func() {
		ports := [NUM_PORTS]string{"20", "21", "22", "23", "53", "80", "123", "179", "443", "500", "587", "3389"}
		for _, port := range ports {
			wg.Add(1)
			go func(port string) {
				defer wg.Done()
				ipPort := ip + ":" + port
				if _, err := net.Dial("tcp", ipPort); err == nil {
					opened_ports = append(opened_ports, port)
				} else {
					// log error
				}
			}(port)
		}
		wg.Wait()
	}

	scan()

	return opened_ports, nil
}

// TODO
func Scan(db *sql.DB, errCh chan<- error) {
	randIntCh := make(chan (uint32))
	numWorkers := runtime.NumCPU()

	for i := 0; i < numWorkers; i++ {
		go func() {
			for r := range randIntCh {
				// create the ip address
				var b1, b2, b3, b4 uint8

				b1 = uint8(r & 0xFF)
				b2 = uint8((r >> 8) & 0xFF)
				b3 = uint8((r >> 16) & 0xFF)
				b4 = uint8(r >> 24)

				ipv4 := net.IPv4(b4, b3, b2, b1)

				// check to see if ip was already found (don't want to make API request if already found)

				// scan ip address
				opened_ports, err := scanPorts(ipv4.String())
				if err != nil {
					errCh <- err
				}

				// find geo location

				// insert ip and opened ports into db

			}
		}()
	}

	for {
		genRandInt(randIntCh)
	}
}
