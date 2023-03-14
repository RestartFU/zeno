package module

// Priority represents a priority of a module.
type Priority struct {
	priority
}

type priority uint8

// LowPriority ...
func LowPriority() Priority {
	return Priority{priority(0)}
}

// MediumPriority ...
func MediumPriority() Priority {
	return Priority{priority(1)}
}

// HighPriority ...
func HighPriority() Priority {
	return Priority{priority(2)}
}

// Priorities ...
func Priorities() []Priority {
	return []Priority{
		HighPriority(),
		MediumPriority(),
		LowPriority(),
	}
}

// Uint8 returns the double plant as a uint8.
func (d priority) Uint8() uint8 {
	return uint8(d)
}
