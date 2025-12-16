package system

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/u00io/nuiforms/ui"
)

type System struct {
	mtx sync.Mutex

	events []Event

	started  bool
	aborting bool
	host     string

	result Result

	resultTableText string

	processNamesById map[uint32]string
}

type Event struct {
	Name      string
	Parameter string
}

var Instance *System

func NewSystem() *System {
	var c System
	return &c
}

func (c *System) Start() {
	go c.thWork()
}

func (c *System) Stop() {
}

func (c *System) IsStarted() bool {
	return c.started
}

func (c *System) IsAborting() bool {
	return c.aborting
}

func (c *System) Run(host string) {
	c.mtx.Lock()
	if c.started {
		c.mtx.Unlock()
		return
	}
	c.started = true
	c.host = host
	c.mtx.Unlock()
}

func (c *System) Abort() {
	c.mtx.Lock()
	c.aborting = true
	c.mtx.Unlock()
}

func (c *System) GetResult() Result {
	var res Result
	c.mtx.Lock()
	res = c.result
	c.mtx.Unlock()
	return res
}

func (c *System) thWork() {
	for {
		time.Sleep(100 * time.Millisecond)
		if !c.started {
			continue
		}

		c.mtx.Lock()
		c.result = *NewResult()
		c.result.Status = "Running"
		c.mtx.Unlock()

		target := c.host
		maxTTL := 30

		if target == "" {
			c.mtx.Lock()
			c.result.Status = "ERR: No target specified"
			c.started = false
			c.aborting = false
			c.mtx.Unlock()
			continue
		}

		ipAddr, err := net.ResolveIPAddr("ip4", target)
		if err != nil {
			fmt.Println("ResolveIPAddr error:", err)
			c.mtx.Lock()
			c.result.Status = "ERR: " + err.Error()
			c.started = false
			c.aborting = false
			c.mtx.Unlock()
			continue
		}

		if ipAddr.IP == nil {
			fmt.Println("Could not resolve IP address for [" + target + "]")
			c.mtx.Lock()
			c.result.Status = "ERR: Could not resolve IP address"
			c.started = false
			c.aborting = false
			c.mtx.Unlock()
			continue
		}

		dst := binary.LittleEndian.Uint32(ipAddr.IP.To4())

		c.mtx.Lock()
		c.result.IP = ipAddr.IP.String()
		c.result.CountryName, _ = GetCountryByIP(c.result.IP)
		c.result.CountryISO, _ = GetCountryISOCodeByIP(c.result.IP)
		c.mtx.Unlock()

		icmp := icmpCreateFile()

		fmt.Printf("ICMP traceroute to %s (%s)\n", target, ipAddr.IP)

		for ttl := 1; ttl <= maxTTL; ttl++ {
			start := time.Now()

			c.mtx.Lock()
			aborting := c.aborting
			c.mtx.Unlock()
			if aborting {
				break
			}

			hop := &ResultHop{}
			hop.IP = ""
			hop.TimeMs = -1
			c.mtx.Lock()
			c.result.Hops = append(c.result.Hops, hop)
			c.mtx.Unlock()

			reply, ok := icmpSendEcho(icmp, dst, ttl, 1)
			elapsed := time.Since(start)

			fmt.Printf("%2d  ", ttl)

			if !ok {
				// No reply - go to next hop
				c.mtx.Lock()
				c.result.Hops[len(c.result.Hops)-1].IP = ""
				c.result.Hops[len(c.result.Hops)-1].TimeMs = -1
				c.mtx.Unlock()
				continue
			}

			hopIP := net.IPv4(
				byte(reply.Address),
				byte(reply.Address>>8),
				byte(reply.Address>>16),
				byte(reply.Address>>24),
			)

			c.mtx.Lock()
			c.result.Hops[len(c.result.Hops)-1].IP = hopIP.String()
			c.result.Hops[len(c.result.Hops)-1].TimeMs = int64(float64(elapsed.Microseconds()) / 1000)
			c.result.Hops[len(c.result.Hops)-1].CountryName, _ = GetCountryByIP(hopIP.String())
			c.result.Hops[len(c.result.Hops)-1].CountryISO, _ = GetCountryISOCodeByIP(hopIP.String())
			c.mtx.Unlock()

			fmt.Printf("%s  %.2f ms\n", hopIP, float64(elapsed.Microseconds())/1000)
			if reply.Status == IP_SUCCESS {
				// Reached destination
				break
			}
		}
		icmpClose(icmp)

		c.mtx.Lock()
		c.started = false
		c.aborting = false
		c.result.Status = "Finished"
		c.mtx.Unlock()
	}
}

func (c *System) EmitEvent(event string, parameter string) {
	c.mtx.Lock()
	c.events = append(c.events, Event{Name: event, Parameter: parameter})
	c.mtx.Unlock()
}

func (c *System) GetAndClearEvents() []Event {
	c.mtx.Lock()
	events := c.events
	c.events = make([]Event, 0)
	c.mtx.Unlock()
	return events
}

var (
	iphlpapi         = syscall.NewLazyDLL("iphlpapi.dll")
	procIcmpCreate   = iphlpapi.NewProc("IcmpCreateFile")
	procIcmpClose    = iphlpapi.NewProc("IcmpCloseHandle")
	procIcmpSendEcho = iphlpapi.NewProc("IcmpSendEcho")
)

type ICMP_ECHO_REPLY struct {
	Address       uint32
	Status        uint32
	RoundTripTime uint32
	DataSize      uint16
	Reserved      uint16
	Data          uintptr
	Options       uintptr
}

const (
	IP_SUCCESS     = 0
	IP_TTL_EXPIRED = 11013
)

func icmpCreateFile() syscall.Handle {
	h, _, _ := procIcmpCreate.Call()
	return syscall.Handle(h)
}

func icmpClose(h syscall.Handle) {
	procIcmpClose.Call(uintptr(h))
}

func icmpSendEcho(
	h syscall.Handle,
	dst uint32,
	ttl int,
	timeout int,
) (*ICMP_ECHO_REPLY, bool) {

	type ipOptions struct {
		Ttl         byte
		Tos         byte
		Flags       byte
		OptionsSize byte
		OptionsData uintptr
	}

	opts := ipOptions{Ttl: byte(ttl)}
	buf := make([]byte, 64)

	r, _, _ := procIcmpSendEcho.Call(
		uintptr(h),
		uintptr(dst),
		0,
		0,
		uintptr(unsafe.Pointer(&opts)),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
		uintptr(timeout),
	)

	if r == 0 {
		return nil, false
	}

	reply := *(*ICMP_ECHO_REPLY)(unsafe.Pointer(&buf[0]))
	return &reply, true
}

func (c *System) SetResultTableText(text string) {
	c.mtx.Lock()
	c.resultTableText = text
	c.mtx.Unlock()
}

func (c *System) CopyResultsToClipboard() {
	c.mtx.Lock()
	txt := c.resultTableText
	c.mtx.Unlock()
	ui.ClipboardSetText(txt)
}
