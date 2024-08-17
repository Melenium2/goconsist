package goconsist

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShards(t *testing.T) {
	cfg := Config{
		SectionFactor: 1,
		SectionCount:  4,
	}
	server1 := netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 10)
	server2 := netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 20)
	servers := []netip.AddrPort{server1, server2}

	t.Run("should return new struct with config, configured sections and servers", func(t *testing.T) {
		expected := &Ring{
			servers: servers,
			sections: []section{
				{min: 0, max: 1, serverAddr: server1},
				{min: 2, max: 3, serverAddr: server2},
				{min: 4, max: 5, serverAddr: server1},
				{min: 6, max: 0, serverAddr: server2},
			},
			// Last section min factor plus section factor.
			hashMod: 7,
		}

		res := NewRing(cfg, servers...)
		assert.Equal(t, expected, res)
	})

	t.Run("should return new struct with config, sections and empty servers", func(t *testing.T) {
		expected := &Ring{
			servers: make([]netip.AddrPort, 0),
			sections: []section{
				{min: 0, max: 1},
				{min: 2, max: 3},
				{min: 4, max: 5},
				{min: 6, max: 0},
			},
			// Last section min factor plus section factor.
			hashMod: 7,
		}

		res := NewRing(cfg)
		assert.Equal(t, expected, res)
	})

	t.Run("should return struct with default config, if specified config is empty", func(t *testing.T) {
		expected := &Ring{
			servers: make([]netip.AddrPort, 0),
			sections: []section{
				{min: 0, max: 1},
				{min: 2, max: 3},
				{min: 4, max: 5},
				{min: 6, max: 7},
				{min: 8, max: 9},
				{min: 10, max: 11},
				{min: 12, max: 13},
				{min: 14, max: 15},
				{min: 16, max: 17},
				{min: 18, max: 0},
			},
			// Last section min factor plus section factor.
			hashMod: 19,
		}

		res := NewRing(Config{})
		assert.Equal(t, expected, res)
	})
}

func TestAddServers(t *testing.T) {
	var (
		server1 = netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 10)
		server2 = netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 20)
		servers = []netip.AddrPort{server1, server2}
	)

	t.Run("should add new server and distribute new servers across the ring", func(t *testing.T) {
		distrubution := []section{
			{min: 0, max: 15},
			{min: 16, max: 31},
			{min: 32, max: 47},
			{min: 48, max: 63},
			{min: 64, max: 79},
			{min: 80, max: 0},
		}

		expected := &Ring{
			servers: servers,
			sections: []section{
				{min: 0, max: 15, serverAddr: server1},
				{min: 16, max: 31, serverAddr: server2},
				{min: 32, max: 47, serverAddr: server1},
				{min: 48, max: 63, serverAddr: server2},
				{min: 64, max: 79, serverAddr: server1},
				{min: 80, max: 0, serverAddr: server2},
			},
		}

		ring := &Ring{sections: distrubution}

		ring.AddServers(servers...)

		assert.Equal(t, expected, ring)
	})

	t.Run("should do nothing if servers not provided", func(t *testing.T) {
		distrubution := []section{
			{min: 0, max: 15},
			{min: 16, max: 31},
			{min: 32, max: 47},
			{min: 48, max: 63},
			{min: 64, max: 79},
			{min: 80, max: 0},
		}

		expected := &Ring{
			sections: []section{
				{min: 0, max: 15},
				{min: 16, max: 31},
				{min: 32, max: 47},
				{min: 48, max: 63},
				{min: 64, max: 79},
				{min: 80, max: 0},
			},
		}

		ring := &Ring{sections: distrubution}

		ring.AddServers()

		assert.Equal(t, expected, ring)
	})
}

func TestAcquire(t *testing.T) {
	var (
		server1 = netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 10)
		server2 = netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 20)
		servers = []netip.AddrPort{server1, server2}
	)
	ring := &Ring{
		servers: servers,
		sections: []section{
			{min: 0, max: 15, serverAddr: server1},
			{min: 16, max: 31, serverAddr: server2},
			{min: 32, max: 47, serverAddr: server1},
			{min: 48, max: 63, serverAddr: server2},
			{min: 64, max: 79, serverAddr: server1},
			{min: 80, max: 0, serverAddr: server2},
		},
		// Last section min + step.
		hashMod: 80 + 15,
	}

	t.Run("should receive server address from first section", func(t *testing.T) {
		key := []byte("5")

		res := ring.Acquire(key)
		assert.Equal(t, server1, res)
	})

	t.Run("should receive server address from last section", func(t *testing.T) {
		key := []byte("9")

		res := ring.Acquire(key)
		assert.Equal(t, server2, res)
	})

	t.Run("should receive server address from the middle of the ring", func(t *testing.T) {
		key := []byte("4")

		res := ring.Acquire(key)
		assert.Equal(t, server2, res)
	})
}

func TestRemoveServer(t *testing.T) {
	t.Run("should remove the server from the ring and distribute remaining servers across the ring", func(t *testing.T) {
		var (
			server1 = netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 10)
			server2 = netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 20)
			servers = []netip.AddrPort{server1, server2}
		)
		ring := &Ring{
			servers: servers,
			sections: []section{
				{min: 0, max: 15, serverAddr: server1},
				{min: 16, max: 31, serverAddr: server2},
				{min: 32, max: 47, serverAddr: server1},
				{min: 48, max: 63, serverAddr: server2},
				{min: 64, max: 79, serverAddr: server1},
				{min: 80, max: 0, serverAddr: server2},
			},
		}

		expected := &Ring{
			servers: []netip.AddrPort{server2},
			sections: []section{
				{min: 0, max: 15, serverAddr: server2},
				{min: 16, max: 31, serverAddr: server2},
				{min: 32, max: 47, serverAddr: server2},
				{min: 48, max: 63, serverAddr: server2},
				{min: 64, max: 79, serverAddr: server2},
				{min: 80, max: 0, serverAddr: server2},
			},
		}

		ring.RemoveServer(server1)

		assert.Equal(t, expected, ring)
	})
}

func TestDistributeServers(t *testing.T) {
	t.Run("should distribute servers between every section", func(t *testing.T) {
		var (
			server1 = netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 10)
			server2 = netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 20)
			server3 = netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 20)
			servers = []netip.AddrPort{server1, server2, server3}
		)
		ring := &Ring{
			servers: servers,
			sections: []section{
				{min: 0, max: 15},
				{min: 16, max: 31},
				{min: 32, max: 47},
				{min: 48, max: 63},
				{min: 64, max: 79},
				{min: 80, max: 0},
			},
		}

		expected := &Ring{
			servers: servers,
			sections: []section{
				{min: 0, max: 15, serverAddr: server1},
				{min: 16, max: 31, serverAddr: server2},
				{min: 32, max: 47, serverAddr: server3},
				{min: 48, max: 63, serverAddr: server1},
				{min: 64, max: 79, serverAddr: server2},
				{min: 80, max: 0, serverAddr: server3},
			},
		}

		ring.DistributeServers()

		assert.Equal(t, expected, ring)
	})
}
