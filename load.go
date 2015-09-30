package rrd

import (
	"encoding/xml"
	"os/exec"
)

func (rrd *RRD) Load(path string) error {
	cmd := exec.Command("rrdtool", "dump", path)
	b, err := cmd.Output()
	if err != nil {
		return err
	}
	return xml.Unmarshal(b, rrd)
}
