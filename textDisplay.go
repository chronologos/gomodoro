package main

import (
	ui "github.com/gizak/termui"
)

// TODO(chronologos) I am encapsulating ui.List so we can define methods on it. Is this correct?
type mutableTextBox struct {
	list *ui.List
}

func (mtb *mutableTextBox) updateText(t []string) {
	mtb.list.Items = t
}

// makeHelpBox generates the cli component that lists the simple commands
// that this app supports.
func makeTextBox(x, y, w, h int, name string) *mutableTextBox {
	var mtb mutableTextBox
	l := ui.NewList()
	mtb.list = l
	mtb.list.ItemFgColor = ui.ColorYellow
	mtb.list.BorderLabel = name
	mtb.list.X = x
	mtb.list.Y = y
	mtb.list.Width = w
	mtb.list.Height = h
	return &mtb
}
