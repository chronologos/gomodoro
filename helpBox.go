package main

import (
	ui "github.com/gizak/termui"
)

func makeHelpBox() (list *ui.List) {
	strs := []string{"[s] start timer", "[p] pause timer", "[r] reset timer"}
	list = ui.NewList()
	list.Items = strs
	list.ItemFgColor = ui.ColorYellow
	list.BorderLabel = "Commands"
	list.Height = 7
	list.Width = 25
	list.Y = 4
	return
}
