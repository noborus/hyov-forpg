package internal

import (
	"io"

	"github.com/noborus/ov/oviewer"
)

func setupOviewer(r io.Reader) (*oviewer.Root, error) {
	oviewer.MemoryLimit = -1 // Disable memory limit for oviewer
	doc, err := oviewer.NewDocument()
	if err != nil {
		return nil, err
	}
	doc.ControlReader(r, nil)

	ov, err := oviewer.NewOviewer(doc)
	if err != nil {
		return nil, err
	}
	return ov, nil
}

func configureOviewer(ov *oviewer.Root) {
	header := 1
	delimiter := "|"
	columnMode := true
	columnRainbow := true
	align := true

	ov.Config.General = oviewer.General{
		Header:          &header,
		ColumnDelimiter: &delimiter,
		ColumnMode:      &columnMode,
		ColumnRainbow:   &columnRainbow,
		Align:           &align,
	}
	ov.SetConfig(ov.Config)
}
