package homemate

import (
	"database/sql/driver"
	"fmt"
)

type Role int

// Define the enum values
const (
	Mate  Role = iota // 0 (Catch-all/default)
	Admin             // 1
)

// String implements the fmt.Stringer interface for clean logging/printing
func (r Role) String() string {
	switch r {
	case Mate:
		return "Mate"
	case Admin:
		return "Admin"
	default:
		return "Mate"
	}
}

// Value converts the Go enum into a driver.Value to save into SQLite
func (r Role) Value() (driver.Value, error) {
	if r < Mate || r > Admin {
		return nil, fmt.Errorf("invalid role value: %d", r)
	}
	return int64(r), nil
}

// Scan converts the SQLite value back into the Go custom Role type
func (r *Role) Scan(value interface{}) error {
	if value == nil {
		*r = Mate
		return nil
	}

	intVal, ok := value.(int64)
	if !ok {
		return fmt.Errorf("invalid type for Role: %T", value)
	}

	castedRole := Role(intVal)
	if castedRole < Mate || castedRole > Admin {
		return fmt.Errorf("invalid role database value: %d", intVal)
	}

	*r = castedRole
	return nil
}
