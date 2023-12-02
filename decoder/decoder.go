// Package decoder implements a simple decoder for wiki pages.
package decoder

import (
	"encoding/xml"
	"fmt"
	"io"
)

// Decoder is an XML decoder tailored to the Wikipedia dataset.
type Decoder struct {
	err   error
	start *xml.StartElement

	d *xml.Decoder
}

// New instanciates a new Decoder.
//
// It fail if it cannot find the mediawiki and siteinfo elements from the dataset.
func New(r io.Reader) (*Decoder, error) {
	// Instanciates the decoder.
	d := &Decoder{
		d: xml.NewDecoder(r),
	}

	// Read the mediawiki element.
	err := d.checkToken("mediawiki")
	if err != nil {
		return nil, fmt.Errorf("check mediawiki element: %v", err)
	}

	// Skip siteinfo element.
	err = d.consumeToken("siteinfo")
	if err != nil {
		return nil, fmt.Errorf("consume siteinfo element: %v", err)
	}

	return d, nil
}

// Next moves to the next element.
func (d *Decoder) Next() bool {
	if d.err != nil {
		return false
	}

	for {
		t, err := d.d.Token()
		if err != nil {
			d.start, d.err = nil, err
			return false
		}

		start, ok := t.(xml.StartElement)
		if ok {
			d.start = &start
			break
		}
	}

	return true
}

// Err returns any error encountered by Next.
func (d *Decoder) Err() error { return d.err }

// Scan scans an element.
func (d *Decoder) Scan(v any) error {
	if d.err != nil {
		return d.err
	}
	if d.start == nil {
		return fmt.Errorf("Scan called without calling Next")
	}
	return d.d.DecodeElement(v, d.start)
}

func (d *Decoder) checkToken(name string) error {
	if !d.Next() {
		return d.err
	}

	n := d.start.Name.Local
	if n != name {
		return fmt.Errorf("unexpected name: expects=%v got=%v", n, name)
	}
	return nil
}

func (d *Decoder) consumeToken(name string) error {
	if err := d.checkToken(name); err != nil {
		return err
	}
	if err := d.d.Skip(); err != nil {
		return fmt.Errorf("skip: %v", err)
	}
	return nil
}
