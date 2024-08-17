package goconsist

import "net/netip"

type section struct {
	min        uint32
	max        uint32
	serverAddr netip.AddrPort
}

func (v section) more(x uint32) bool {
	return x > v.max
}

func (v section) less(x uint32) bool {
	return x < v.min
}

func (v section) between(x uint32) bool {
	// NOTE: corner case where maxFactor is end of ring and equals to zero.
	if (v.min > v.max) && x >= v.min {
		return true
	}

	return x >= v.min && x <= v.max
}

// Search for specific section in the ring.
//
// NOTE: If target valued does not fit into any section of the ring
// the last vhsard will be returned.
func search(list []section, target uint32) (section, bool) {
	l, r := 0, len(list)-1

	for l <= r {
		mid := l + (r-l)/2
		curr := list[mid]

		if curr.between(target) {
			return curr, true
		}

		if curr.less(target) {
			r = mid - 1

			continue
		}

		if curr.more(target) {
			l = mid + 1

			continue
		}
	}

	return section{}, false
}

// Distribute ranges across ring.
func distribute(sectionFactor, sectionCount uint32) []section {
	ranges := make([]section, 0, sectionCount)
	start := uint32(0)

	for i := 0; i < int(sectionCount)-1; i++ {
		next := start + sectionFactor
		vshard := section{min: start, max: next}
		start = next + 1

		ranges = append(ranges, vshard)
	}

	lastShard := section{min: start, max: 0}
	ranges = append(ranges, lastShard)

	return ranges
}
