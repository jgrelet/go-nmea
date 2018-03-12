package nmea

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Allowed cardinal points
const (
	// North as CardinalPoint and value "N"
	North CardinalPoint = "N"
	// South as CardinalPoint and value "S"
	South CardinalPoint = "S"
	// East as CardinalPoint and value "E"
	East CardinalPoint = "E"
	// West as CardinalPoint and value "W"
	West CardinalPoint = "W"
)

// CardinalPoint type as string
type CardinalPoint string

func (c CardinalPoint) String() string {
	return string(c)
}

// ParseCardinalPoint check CardinalPoint validity, return an error
// "unknow value" if not
func ParseCardinalPoint(raw string) (cp CardinalPoint, err error) {
	cp = CardinalPoint(raw)
	switch cp {
	case North, South, East, West:
	default:
		err = fmt.Errorf("unknow value")
	}
	return
}

const (
	// LatLong Thresholds (ie: spherical degrees)
	// Min is the minimum value allowed for a LatLong
	//MinLong LatLong = -180
	//MinLat LatLong = -90

	// MaxLong is the maximum value allowed for a longitude
	MaxLong float64 = 180
	// MaxLat is the maximum value allowed for a latitude
	MaxLat float64 = 90
)

// LatLong type as float64
type LatLong float64

// NewLatLong parses input has coordinate or return error
// Allowed format:
// - DMS (Degrees, Minutes, Secondes), ie: "N 31° 50' 72.38'"
// - DD (Decimal Degree), ie: "31.8534389" "22.870216666666668"
func NewLatLong(raw string) (l LatLong, err error) {

	if strings.TrimSpace(raw) == "" {
		err = fmt.Errorf("Invalid LatLong, can't be empty")
		return
	}

	if l, err = ParseDM(raw); err != nil {
		return
	}

	return
}

// ParseDM return LatLong from provided format from GPS module (in format ‘ddmm.mmmm’: degree and minutes)
// Allowed format: "3150.7238N" or "3150.7238 N"
// @see https://fr.wikipedia.org/wiki/Coordonn%C3%A9es_g%C3%A9ographiques
// => 1 degree = 60 minutes
// => 1 minute = 60 secondes
// Example: Baltimore (United state) => latitude = 39,28° N, longitude = 76,60° O (39° 17′ N, 76° 36′ O).
// 0.28° = (0.28°*60min)/1° = 16.8min => ~17 minutes
func ParseDM(raw string) (LatLong, error) {

	var (
		dir CardinalPoint
		dm  float64
		err error
	)

	if len(raw) < 2 {
		return LatLong(0), fmt.Errorf("nmea.ParseDM() Wrong DM format, got: \"%s\"", raw)
	}

	// Explode data
	if dm, err = strconv.ParseFloat(strings.TrimSpace(raw[:len(raw)-2]), 64); err != nil {
		return LatLong(0), err
	}

	if dir, err = ParseCardinalPoint(string(raw[len(raw)-1])); err != nil {
		return LatLong(0), err
	}

	// Compute LatLong
	d := math.Floor(dm / 100) // div dm by 100 and truncate decimal value to get only degrees
	m := dm - (d * 100)       // Sub degrees to dm value
	dm = d + m/60             // switch minute to degree to get value in the same referential

	switch dir {
	case North, South:
		if math.Abs(dm) > MaxLat {
			return 0, fmt.Errorf("nmea.ParseDM() invalid range (got: %f)", dm)
		}
	case East, West:
		if math.Abs(dm) > MaxLong {
			return 0, fmt.Errorf("nmea.ParseDM() invalid range (got: %f)", dm)
		}
	default:
		return 0, fmt.Errorf("nmea.ParseDM() Wrong direction (got: %s)", dir.String())
	}

	switch dir {
	case North, East:
		return LatLong(dm), nil
	case South, West:
		return LatLong(0 - dm), nil
	default:
		return 0, fmt.Errorf("nmea.ParseDM() Wrong direction (got: %s)", dir.String())
	}
}

// CardinalPoint return the cardinal point related to the kind of coordinate (long or lat)
func (l LatLong) CardinalPoint(isLatitude bool) CardinalPoint {
	if l == 0 {
		return ""
	}

	if l < 0 {
		if isLatitude {
			return South
		}
		return West
	}

	if isLatitude {
		return North
	}

	return East
}

// DM extract degrees and minutes
func (l LatLong) DM() (int, float64) {
	var d, m float64
	if l >= 0 {
		d = math.Floor(float64(l))
		m = (float64(l) - d) * 60
	} else {
		d = math.Ceil(float64(l))
		m = (math.Abs(float64(l)) - math.Abs(d)) * 60
	}
	// fmt.Println(d, m)
	return int(d), m
}

// DMS extract degrees, minutes and secondes
func (l LatLong) DMS() (int, int, float64) {
	var s float64
	d, m := l.DM()
	m = math.Floor(m)
	if l >= 0 {
		s = (float64(l) - (float64(d) + (m / 60))) * 60 * 60 // TODO: round secondes
	} else {
		s = (float64(l) + (math.Abs(float64(d)) + (m / 60))) * 60 * 60 // TODO: round secondes
	}
	return d, int(m), s
}

// ToDM return string like ‘ddmm.mmmm’: degree and minutes as GPS module provide
func (l LatLong) ToDM() string {
	if l == 0 {
		return ""
	}
	d, m := l.DM()
	return strings.Trim(fmt.Sprintf("%02d%09.6f", d, m), "0")
}

// PrintDMS return string like: dd° mm' ss.ss" to be human readable
func (l LatLong) PrintDMS() string {
	degrees, minutes, secondes := l.DMS()
	return fmt.Sprintf("%d° %d' %f\"", degrees, minutes, secondes)
}
