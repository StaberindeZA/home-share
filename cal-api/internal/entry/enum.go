package entry

import (
	"database/sql/driver"
	"fmt"
)

// Role represents the user role enum
type EntryValue int

// Define the enum values
const (
	Open       EntryValue = iota // 0 (Catch-all/default)
	Booked                       // 1
	Talking                      // 2
	NotTalking                   // 3
)

// String implements the fmt.Stringer interface for clean logging/printing
func (ev EntryValue) String() string {
	switch ev {
	case Booked:
		return "Booked"
	case Talking:
		return "Talking"
	case NotTalking:
		return "Not Talking"
	default:
		return "Open"
	}
}

// Value converts the Go enum into a driver.Value to save into SQLite
func (ev EntryValue) Value() (driver.Value, error) {
	if ev < Open || ev > NotTalking {
		return nil, fmt.Errorf("invalid role value: %d", ev)
	}
	return int64(ev), nil
}

// Scan converts the SQLite value back into the Go custom Role type
func (ev *EntryValue) Scan(value interface{}) error {
	if value == nil {
		*ev = Open
		return nil
	}

	intVal, ok := value.(int64)
	if !ok {
		return fmt.Errorf("invalid type for Role: %T", value)
	}

	castedRole := EntryValue(intVal)
	if castedRole < Open || castedRole > NotTalking {
		return fmt.Errorf("invalid role database value: %d", intVal)
	}

	*ev = castedRole
	return nil
}
