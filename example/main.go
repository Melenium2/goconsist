package main

import (
	"log/slog"
	"net/netip"

	"github.com/google/uuid"

	"github.com/Melenium2/goconsist"
)

var servers = []netip.AddrPort{
	netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 10),
	netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 20),
	netip.AddrPortFrom(netip.AddrFrom4([4]byte{10, 1, 1, 1}), 30),
}

func main() {
	ring := goconsist.NewRing(goconsist.Config{}, servers...)

	for i := 0; i < 100; i++ {
		id, _ := uuid.New().MarshalBinary()

		server := ring.Acquire(id)

		slog.Info("server acquired", "server", server)
	}
}
