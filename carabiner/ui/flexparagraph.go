// Package ui implements UI logic for carabiner.
package ui

import (
	"github.com/gizak/termui/v3/widgets"
)

// FlexParagraph represents a flexible Paragraph.
type FlexParagraph struct {
	widgets.List
}

// NewFlexParagraph returns a new FlexParagraph.
func NewFlexParagraph() *FlexParagraph {
	return &FlexParagraph{
		List: *widgets.NewList(),
	}
}

// AddLine adds a line of text to the Paragraph.
func (f *FlexParagraph) AddLine(l string) {
	f.Rows = append(f.Rows, l)
	f.ScrollDown()
}
