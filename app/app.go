package app

import (
	"database/sql"
	"log"
	"math/rand"
	"net"
	"runtime"
	"sync"
	"time"
)

func genRandInt(randIntCh chan<- uint32, source *rand.Source) {
	randNum := rand.New(*source)
	randIntCh <- randNum.Uint32()
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
func scanPorts(ip string, logger *log.Logger) ([]string, error) {
	var wg sync.WaitGroup
	var openedPorts []string
	scan := func() {
		ports := [NUM_PORTS]string{"20", "21", "22", "23", "53", "80", "123", "179", "443", "500", "587", "3389"}
		for _, port := range ports {
			wg.Add(1)
			go func(port string) {
				defer wg.Done()
				ipPort := ip + ":" + port
				if _, err := net.Dial("tcp", ipPort); err == nil {
					openedPorts = append(openedPorts, port)
				} else {
					logger.Println(err)
				}
			}(port)
		}
		wg.Wait()
	}

	scan()

	return openedPorts, nil
}

// TODO
func Scan(db *sql.DB, errCh chan<- error, logger *log.Logger) {
	randIntCh := make(chan uint32)
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
				var ipToCmp string
				row := db.QueryRow("SELECT ip_addr FROM ip_addresses WHERE ip_addr = $1", ipv4.String())
				if err := row.Scan(&ipToCmp); err == sql.ErrNoRows {
					// scan ip address
					/*
						openedPorts, err := scanPorts(ipv4.String(), logger)
						if err != nil {
							errCh <- err
						}
					*/

					// find geo location

					// insert ip and opened ports into db
				} else {
					logger.Println(err)
				}

			}
		}()
	}

	source := rand.NewSource(time.Now().UnixMicro())
	for {
		genRandInt(randIntCh, &source)
	}
}
