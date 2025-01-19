package nmea

import (
	"fmt"
	"time"
)

/*
	ZDA Time & Date – UTC, Day, Month, Year and Local Time Zone
	1 2 3 4 5 6 7
	| | | | | | |
	$--ZDA,hhmmss.ss,xx,xx,xxxx,xx,xx*hh
	1) Local zone minutes description, same sign as local hours
	2) Local zone description, 00 to +/- 13 hours
	3) Year
	4) Month, 01 to 12
	5) Day, 01 to 31
	6) Time (UTC)
	7) Checksum
*/

// NewGPZDA allocate GPZDA struct for ZDA sentence (Satellites in view)
func NewGPZDA(m Message) *GPZDA {
	return &GPZDA{Message: m}
}

// GPZDA struct
type GPZDA struct {
	Message
	DateTimeUTC  time.Time // Aggregation of TimeUTC data field
}


func (m *GPZDA) parse() (err error) {
	//log.Printf("ZDA: %d fields\n", len(m.Fields))
	if len(m.Fields) != 7 {
		return m.Error(fmt.Errorf("Incomplete GPZDA message, not enougth data fields (got: %d, wanted: %d)", len(m.Fields), 7))
	}

	datetime := fmt.Sprintf("%s%s%s %s", m.Fields[4], m.Fields[3], m.Fields[2], m.Fields[5])
	if m.DateTimeUTC, err = time.Parse("020106 150405.000", datetime); err != nil {
		return m.Error(fmt.Errorf("Unable to parse datetime UTC from data field (got: %s)", datetime))
	}

	return nil
}

// Serialize return a valid sentence ZDA as string
func (m GPZDA) Serialize() string { // Implement NMEA interface
	hdr := TypeIDs["GPZDA"]
	fields := make([]string, 0)

	fields = append(fields,
		m.DateTimeUTC.Format("020106 150405.000"))
	msg := Message{Type: hdr, Fields: fields}
	msg.Checksum = msg.ComputeChecksum()

	return msg.Serialize()
}
