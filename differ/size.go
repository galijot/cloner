package differ

import ()

// SizeUnit represents a unit of size of a file on a system.
type SizeUnit int8

const (
	// B represents a byte
	B SizeUnit = iota
	// KB represents a kilobyte
	KB
	// MB represents a megabyte
	MB
	// GB represents a gigabyte
	GB
)

type sizeValue struct {
	unit  SizeUnit
	value float64
}

func (sv sizeValue) description() string {
	return ""
}

// Size represents a size of a file in different units.
type Size struct {
	b, kb, mb, gb sizeValue
}

// SizeWithBytes creates a Size from the file's bytes.
func SizeWithBytes(bytes int64) Size {

	var kilobytes int64 = (bytes / 1024)
	var megabytes float64 = float64(kilobytes / 1024) // cast to type float64
	var gigabytes float64 = megabytes / 1024

	return Size{
		b: sizeValue{
			unit:  B,
			value: float64(bytes),
		},
		kb: sizeValue{
			unit:  KB,
			value: float64(kilobytes),
		},
		mb: sizeValue{
			unit:  MB,
			value: megabytes,
		},
		gb: sizeValue{
			unit:  GB,
			value: gigabytes,
		},
	}
}
