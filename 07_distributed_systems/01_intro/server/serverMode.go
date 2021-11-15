package main

type serverMode int

// abozhenko for oz:
// Is there a way to declare this enum
// in a way that will enforce consitency, like in oCaml
// you can do variant types and patter matching,
// and compiler will bark at you for non-exhaustive matches
//  https://www.cs.cornell.edu/courses/cs3110/2020fa/textbook/data/variants.html
func (mode serverMode) String() string {
	switch mode {
	case PRIMARY:
		return "primary"
	case PRIMARY_PARTITION:
		return "primary_partition"
	case SYNCHRONOUS_FOLLOWER:
		return "sync_follower"
	case ASYNCHRONOUS_FOLLOWER:
		return "async_follower"
	default:
		return "UNDEFINED"
	}
}

const (
	PRIMARY = serverMode(iota)
	SYNCHRONOUS_FOLLOWER
	ASYNCHRONOUS_FOLLOWER
	PRIMARY_PARTITION
)
