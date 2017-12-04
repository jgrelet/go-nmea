package nmea

import (
	"fmt"
	"strconv"
)

// Examples:
// $MXDBT,108.34,f,33.02,M,18.06,F,*09
// $INDBT,,,000033.0,M,,*06
//         1  2  3  4  5  6 7
//         |  |  |  |  |  | |
// $--DBT,x.x,f,x.x,M,x.x,F*hh
//
// 1) Depth, feet
// 2) f = feet
// 3) Depth, meters
// 4) M = meters
// 5) Depth, Fathoms
// 6) F = Fathoms
// 7) Checksum

// NewGPDBT allocate echo-sounder sentence DBT (Depth Below Transducer)
func NewGPDBT(m Message) *GPDBT {
	return &GPDBT{Message: m}
}

// GPDBT struct
type GPDBT struct {
	Message

	DepthInFeet    float64
	DepthInMeters  float64
	DepthInFathoms float64
}

func (m *GPDBT) parse() (err error) {
	if len(m.Fields) != 6 {
		return m.Error(fmt.Errorf("Incomplete GPDBT message, not enougth data fields (got: %d, wanted: %d)", len(m.Fields), 6))
	}

	// Validate fixed field
	for i, v := range map[int]string{1: "f", 3: "M", 5: "F"} {
		if m.Fields[i] != v {
			return m.Error(fmt.Errorf("Invalid fixed field at %d (got: %s, wanted: %s)", i+1, m.Fields[i], v))
		}
	}

	if m.DepthInFeet, err = strconv.ParseFloat(m.Fields[0], 64); err != nil {
		return m.Error(fmt.Errorf("Unable to parse depth in feet from data field (got: %s)", m.Fields[0]))
	}

	if m.DepthInMeters, err = strconv.ParseFloat(m.Fields[2], 64); err != nil {
		return m.Error(fmt.Errorf("Unable to parse depth in meters from data field (got: %s)", m.Fields[4]))
	}

	if m.DepthInFathoms, err = strconv.ParseFloat(m.Fields[4], 64); err != nil {
		return m.Error(fmt.Errorf("Unable to parse depth in fathoms from data field (got: %s)", m.Fields[6]))
	}

	return nil
}
