package cache

import (
	"config"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/mdns"
)

const gadLocalDiscoveryAddress = "_gad._clickyab.local"

type PeerFinder interface {
	Self() string
	FindPeer(key string) (string, bool)
	AllPeers() map[string]string
	Change() <-chan map[string]string
}

type peerPicker struct {
	lock    *sync.RWMutex
	peers   map[string]string
	self    string
	port    int
	service *mdns.MDNSService
	server  *mdns.Server

	ticker *time.Ticker
	change chan map[string]string
}

//func getInterfaceIPv4(in *net.Interface) (net.IP, error) {
//	addrs, err := in.Addrs()
//	if err != nil {
//		return nil, err
//	}
//
//	for _, addr := range addrs {
//		var ip net.IP
//		switch v := addr.(type) {
//		case *net.IPNet:
//			ip = v.IP
//		case *net.IPAddr:
//			ip = v.IP
//		}
//
//		if ip.To4() != nil {
//			return ip.To4(), nil
//		}
//	}
//
//	return nil, errors.New("no ip v4 on this interface")
//}

func NewPeerSelector(ip net.IP, port int) (PeerFinder, error) {
	service, err := mdns.NewMDNSService(
		config.Config.MachineName,
		gadLocalDiscoveryAddress,
		"",
		"",
		port,
		[]net.IP{ip},
		[]string{"the gad local server"},
	)

	if err != nil {
		return nil, err
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return nil, err
	}

	res := peerPicker{
		lock:    &sync.RWMutex{},
		peers:   make(map[string]string),
		self:    fmt.Sprintf("http://%s:%d", ip.String(), port),
		port:    port,
		service: service,
		server:  server,
		ticker:  time.NewTicker(time.Minute),
		change:  make(chan map[string]string),
	}

	go func() {
		defer res.server.Shutdown()
		time.Sleep(time.Second)
		res.refresh()
		for range res.ticker.C {
			res.refresh()
		}

	}()
	return &res, nil
}

func (p *peerPicker) refresh() {
	services := make(chan *mdns.ServiceEntry, 5)
	defer close(services)
	go p.reload(services)
	mdns.Lookup(gadLocalDiscoveryAddress, services)
}

func (p *peerPicker) signal(m map[string]string) {
	select {
	case p.change <- m:
	default:
	}
}

func (p *peerPicker) reload(d chan *mdns.ServiceEntry) {
	p.lock.Lock()
	defer p.lock.Unlock()
	newp := make(map[string]string)
	for s := range d {
		parts := strings.SplitN(s.Name, ".", 2)
		if len(parts) != 2 {
			continue
		}

		newp[parts[0]] = fmt.Sprintf("http://%s:%d", s.AddrV4, s.Port)
	}

	if len(newp) != len(p.peers) {
		// changed
		p.peers = newp
		p.signal(newp)
		return
	}

	for i := range newp {
		if newp[i] != p.peers[i] {
			p.peers = newp
			p.signal(newp)
			return
		}
	}
}

func (p *peerPicker) FindPeer(key string) (string, bool) {
	parts := strings.SplitN(key, "-", 2)
	if len(parts) != 2 {
		return "", false
	}
	// No need to find a peer if this is me
	if parts[0] == config.Config.MachineName {
		return "", false
	}
	p.lock.RLock()
	defer p.lock.RUnlock()

	res, ok := p.peers[parts[0]]
	return res, ok
}

func (p *peerPicker) AllPeers() map[string]string {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.peers
}

func (p *peerPicker) Self() string {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.self
}

func (p *peerPicker) Change() <-chan map[string]string {
	return p.change
}
