package goconsist

import (
	"net/netip"

	"github.com/twmb/murmur3"
)

// Ring is not concurrency safety implementation of consistency hash ring with
// virtual shards.
type Ring struct {
	// A module value to prevent going beyond the ring.
	//
	// This value is calculated as the minimum value of the last section plus
	// the section factor from configuration.
	//
	// Example:
	//	last section - start: 48, end 0, server: 2
	//	sectionFactor - 15
	//	hashMod - 48 + 15 = 63
	hashMod uint32

	// Ring sections.
	//
	// Example:
	//
	//	start: 0, end 15, server: 1,
	//	start: 16, end 31, server: 2,
	//	start: 32, end 47, server: 1,
	//	start: 48, end 0, server: 2
	sections []section
	// Servers that are distributed in a ring.
	servers []netip.AddrPort
}

func NewRing(config Config, servers ...netip.AddrPort) *Ring {
	defaultCfg := defaultConfig()

	if config.SectionFactor > 0 {
		defaultCfg.SectionFactor = config.SectionFactor
	}

	if config.SectionCount > 0 {
		defaultCfg.SectionCount = config.SectionCount
	}

	distrubution := distribute(defaultCfg.SectionFactor, defaultCfg.SectionCount)
	lastShard := distrubution[len(distrubution)-1]
	hashMod := lastShard.min + defaultCfg.SectionFactor

	s := &Ring{
		servers:  make([]netip.AddrPort, 0),
		sections: distrubution,
		hashMod:  hashMod,
	}

	s.AddServers(servers...)

	return s
}

// AddServers adds new servers and redistribute available
// servers across the ring.
//
// Example:
//
//	start: 0, end 15, server: 1,
//	start: 16, end 31, server: 1,
//	start: 32, end 47, server: 1,
//	start: 48, end 0, server: 1
//
//	Add server2.
//
//	start: 0, end 15, server: 1,
//	start: 16, end 31, server: 2,
//	start: 32, end 47, server: 1,
//	start: 48, end 0, server: 2
func (s *Ring) AddServers(servers ...netip.AddrPort) {
	s.servers = append(s.servers, servers...)

	s.DistributeServers()
}

// DistributeServers starts server balancing across the ring.
// If no servers were previously distributed, then they
// will be distributed.
func (s *Ring) DistributeServers() {
	if len(s.servers) == 0 {
		return
	}

	for i, curr := range s.sections {
		servIndex := i % len(s.servers)

		curr.serverAddr = s.servers[servIndex]

		s.sections[i] = curr
	}
}

// Acquire searches for server which falls within the range
// calculated by hash function. The function uses murmur3 hash algorithm
// for calculating hash.
//
// Example:
//
//	sections:
//		start: 0, end 15, server: 1,
//		start: 16, end 31, server: 2,
//		start: 32, end 47, server: 1,
//		start: 48, end 0, server: 2
//
//	key:
//	  []byte("4")
//
//	hashMod = 63 (48 + 15)
//	hash(50) = murmur3([]byte("4")) % hashMod
func (s *Ring) Acquire(key []byte) netip.AddrPort {
	hash := murmur3.Sum32(key) % s.hashMod

	result, ok := search(s.sections, hash)
	if !ok {
		return netip.AddrPort{}
	}

	return result.serverAddr
}

// RemoveServer removes the server with the provided IP address.
// Redistribute servers after removal.
//
// Example:
//
//	start: 0, end 15, server: 1,
//	start: 16, end 31, server: 2,
//	start: 32, end 47, server: 1,
//	start: 48, end 0, server: 2
//
//	Remove server2.
//
//	start: 0, end 15, server: 1,
//	start: 16, end 31, server: 1,
//	start: 32, end 47, server: 1,
//	start: 48, end 0, server: 1
func (s *Ring) RemoveServer(server netip.AddrPort) {
	var index int

	for i := 0; i < len(s.servers); i++ {
		if s.servers[i] != server {
			continue
		}

		index = i

		break
	}

	s.servers = append(s.servers[:index], s.servers[index+1:]...)

	s.DistributeServers()
}
