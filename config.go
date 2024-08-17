package goconsist

const (
	defaultSectionFactor = uint32(1)
	defaultSectionCount  = uint32(10)
)

// Config is configuration to Ring structure.
type Config struct {
	// SectionFactor is a range of numbers included to single ring section.
	//
	// Example:
	//  Given a ring of 3 ranges:
	//  0 - 2, 3 - 5, 6 - 0.
	//  In this case, shard factor equals to 2.
	//
	// Be default: defaultSectionFactor.
	SectionFactor uint32
	// SectionCount is a number of ranges located in the ring.
	//
	// Example:
	//  0 - 1, 2 - 3, 4 - 5, 6 - 0.
	//  In this case ranges count equals to 4.
	//
	// By default: defaultSectionCount.
	SectionCount uint32
}

func defaultConfig() Config {
	return Config{
		SectionFactor: defaultSectionFactor,
		SectionCount:  defaultSectionCount,
	}
}
