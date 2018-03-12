package nmea

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

/*
RMC Recommended Minimum Navigation Information
       1         2 3       4 5        6 7   8   9   10  11 12
       |         | |       | |        | |   |   |    |   | |
$--RMC,hhmmss.ss,A,llll.ll,a,yyyyy.yy,a,x.x,x.x,xxxx,x.x,a*hh

1) Time (UTC)
2) Status, V = Navigation receiver warning
3) Latitude
4) N or S
5) Longitude
6) E or W
7) Speed over ground, knots
8) Track made good, degrees true
9) Date, ddmmyy
10) Magnetic Variation, degrees
11) E or W
12) Checksum

 Examples:
 $GPRMC,013732.000,A,3150.7238,N,11711.7278,E,0.00,0.00,220413,,,A*68
 $GPRMC,081836,A,3751.65,S,14507.36,E,000.0,360.0,130998,011.3,E*62
 $GPRMC,225446,A,4916.45,N,12311.12,W,000.5,054.7,191194,020.3,E*68
 $GPRMC,220516,A,5133.82,N,00042.24,W,173.8,231.8,130694,004.2,W*70
*/

// NewGPRMC allocate GPRMC struct for RMC sentence
// (RMC Recommended Minimum Navigation Information)
func NewGPRMC(m Message) *GPRMC {
	return &GPRMC{Message: m}
}

// GPRMC struct
type GPRMC struct {
	Message

	DateTimeUTC       time.Time // Aggregation of TimeUTC+Date data field
	IsValid           DataValid // 'V' =Invalid / 'A' = Valid
	Latitude          LatLong   // In decimal format
	Longitude         LatLong   // In decimal format
	Speed             float64   // Speed over ground in knots
	COG               float64   // Course over ground in degree
	MagneticVariation float64   // Magnetic variation in degree, not being output
	PositioningMode   PositioningMode
}

func (m *GPRMC) parse() (err error) {
	if len(m.Fields) != 12 {
		return m.Error(fmt.Errorf("Incomplete GPRMC message, not enougth data fields (got: %d, wanted: %d)", len(m.Fields), 12))
	}

	datetime := fmt.Sprintf("%s %s", m.Fields[8], m.Fields[0])
	if m.DateTimeUTC, err = time.Parse("020106 150405.000", datetime); err != nil {
		return m.Error(fmt.Errorf("Unable to parse datetime UTC from data field (got: %s)", datetime))
	}

	m.IsValid = (m.Fields[1] == "A")

	if latitude := strings.TrimSpace(strings.Join(m.Fields[2:4], " ")); len(latitude) > 0 {
		if m.Latitude, err = NewLatLong(latitude); err != nil {
			return m.Error(err)
		}
	}
	if longitude := strings.TrimSpace(strings.Join(m.Fields[4:6], " ")); len(longitude) > 0 {
		if m.Longitude, err = NewLatLong(longitude); err != nil {
			return m.Error(err)
		}
	}

	if m.Speed, err = strconv.ParseFloat(m.Fields[6], 64); err != nil {
		return m.Error(fmt.Errorf("Unable to parse speed from data field (got: %s)", m.Fields[6]))
	}

	if m.COG, err = strconv.ParseFloat(m.Fields[7], 64); err != nil {
		return m.Error(fmt.Errorf("Unable to parse course over ground from data field (got: %s)", m.Fields[7]))
	}

	if len(m.Fields[9]) > 0 {
		if m.MagneticVariation, err = strconv.ParseFloat(m.Fields[9], 64); err != nil {
			return m.Error(fmt.Errorf("Unable to parse magnetic variation from data field (got: %s)", m.Fields[9]))
		}

		if len(m.Fields[10]) > 0 {
			magneticVariationDir, err := ParseCardinalPoint(m.Fields[10])
			if err != nil {
				return m.Error(fmt.Errorf("Unable to parse magnetic variation indicator from data field (got: %s)", m.Fields[10]))
			}

			switch magneticVariationDir {
			case West:
				m.MagneticVariation = 0 - m.MagneticVariation
			case East:
				// Allowed direction
			default:
				return m.Error(fmt.Errorf("Wrong magnetic variation direction (got: %s)", m.Fields[10]))
			}
		}
	}

	if m.PositioningMode, err = ParsePositioningMode(m.Fields[11]); err != nil {
		return m.Error(fmt.Errorf("Unable to parse GPS positioning mode from data field (got: %s)", m.Fields[11]))
	}

	return nil
}
