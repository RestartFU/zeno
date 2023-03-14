package specs

// PingRange represents a ping range that a user may prefer.
type PingRange struct {
	unrestricted bool
	ping         int
}

// PingRangeUnrestricted returns the ping group with no restrictions.
func PingRangeUnrestricted() PingRange {
	return PingRange{unrestricted: true}
}

// PingRangeTwentyFive returns the ping group of twenty-five or below.
func PingRangeTwentyFive() PingRange {
	return PingRange{ping: 25}
}

// PingRangeFifty returns the ping group of fifty or below.
func PingRangeFifty() PingRange {
	return PingRange{ping: 50}
}

// PingRangeSeventyFive returns the ping group of seventy-five or below.
func PingRangeSeventyFive() PingRange {
	return PingRange{ping: 75}
}

// PingRangeHundred returns the ping group of a hundred or below.
func PingRangeHundred() PingRange {
	return PingRange{ping: 100}
}

// PingRangeHundredTwentyFive returns the ping group of a hundred and twenty-five or below.
func PingRangeHundredTwentyFive() PingRange {
	return PingRange{ping: 125}
}

// PingRangeHundredFifty returns the ping group of a hundred and fifty or below.
func PingRangeHundredFifty() PingRange {
	return PingRange{ping: 150}
}

// PingRanges returns all possible ping groups.
func PingRanges() []PingRange {
	return []PingRange{
		PingRangeUnrestricted(),
		PingRangeTwentyFive(),
		PingRangeFifty(),
		PingRangeSeventyFive(),
		PingRangeHundred(),
		PingRangeHundredTwentyFive(),
		PingRangeHundredFifty(),
	}
}

// Name ...
func (p PingRange) Name() string {
	switch {
	case p.ping == 25:
		return "0-25"
	case p.ping == 50:
		return "25-50"
	case p.ping == 75:
		return "50-75"
	case p.ping == 100:
		return "175-100"
	case p.ping == 125:
		return "100-125"
	case p.ping == 150:
		return "125-150"
	}
	return "Unrestricted"
}

// Unrestricted ...
func (p PingRange) Unrestricted() bool {
	return p.unrestricted
}

// Compare ...
func (p PingRange) Compare(ping int, unrestricted bool) bool {
	if p.unrestricted && unrestricted {
		// If both groups are unrestricted, we can match them together.
		return true
	}
	return ping <= p.ping
}
