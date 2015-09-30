package rrd

import (
	"encoding/xml"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type Value float64
type SpacedInt int
type SpacedString string
type UnixTimestamp time.Time
type DurationSeconds time.Duration

type RRD struct {
	Version    string          `xml:"version"`
	Step       DurationSeconds `xml:"step"`
	LastUpdate UnixTimestamp   `xml:"lastupdate"`
	DS         []DS            `xml:"ds"`
	RRA        []RRA           `xml:"rra"`
}

type DS struct {
	Name             SpacedString `xml:"name"`
	Type             SpacedString `xml:"type"`
	MinimalHeartbeat int          `xml:"minimal_heartbeat"`
	Min              Value        `xml:"min"`
	Max              Value        `xml:"max"`
	LastDS           int          `xml:"last_ds"`
	Value            Value        `xml:"value"`
	UnknownSec       SpacedInt    `xml:"unknown_sec"`
}

type RRA struct {
	CF        string   `xml:"cf"`
	PDPPerRow int      `xml:"pdp_per_row"`
	Params    *Params  `xml:"params"`
	CDPPrep   CDPPrep  `xml:"cdp_prep"`
	Database  Database `xml:"database"`
}

type Params struct {
	XFF float64 `xml:"xff"`
}

type CDPPrep struct {
	DS []CDPPrepDS `xml:"ds"`
}

type CDPPrepDS struct {
	PrimaryValue      Value `xml:"primary_value"`
	SecondaryValue    Value `xml:"secondary_value"`
	Value             Value `xml:"value"`
	UnknownDatapoints Value `xml:"unknown_datapoints"`
}

type Database struct {
	Row []Row `xml:"row"`
}

type Row struct {
	V []Value `xml:"v"`
}

func (i *SpacedInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var str string
	err := d.DecodeElement(&str, &start)
	if err != nil {
		return err
	}
	parsed, err := strconv.Atoi(strings.Trim(str, " \n\t"))
	*i = SpacedInt(parsed)
	return err
}

func (t *UnixTimestamp) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var ts int64
	err := d.DecodeElement(&ts, &start)
	if err != nil {
		return err
	}
	*t = UnixTimestamp(time.Unix(ts, 0))
	return nil
}

func (s *SpacedString) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var str string
	err := d.DecodeElement(&str, &start)
	if err != nil {
		return err
	}
	*s = SpacedString(strings.Trim(str, " \n\t"))
	return nil
}

func (dur *DurationSeconds) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var ds int64
	err := d.DecodeElement(&ds, &start)
	if err != nil {
		return err
	}
	*dur = DurationSeconds(time.Second * time.Duration(ds))
	return nil
}

func (v Value) MarshalJSON() ([]byte, error) {
	if math.IsNaN(float64(v)) {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%f"`, v)), nil
}
