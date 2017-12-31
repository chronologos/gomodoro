package main

import (
	ui "github.com/gizak/termui"
)

// statsDisplay implements the cli component that tracks pomodoros/breaks.
type statsDisplayT map[pomodoroState]int

type statsDisplay struct {
	barChart *ui.BarChart
}

type statsDisplayUpdate struct{}

func makeStatsDisplay(x, y, w, h int) (sd *statsDisplay) {
	sd = new(statsDisplay)
	bc := ui.NewBarChart()
	sd.barChart = bc
	var names []string
	var vals []int
	for name, val := range stats {
		names = append(names, stateInfoMap[name].shortName)
		vals = append(vals, val)
	}
	bc.BorderLabel = "Bar Chart"
	bc.X = x
	bc.Y = y
	bc.Width = w
	bc.Height = h
	bc.DataLabels = names
	bc.Data = vals
	bc.BarColor = ui.ColorGreen
	bc.NumColor = ui.ColorBlack
	return
}

func (sd *statsDisplay) refreshStatsDisplay(e ui.Event) {
	// sdu := e.Data.(statsDisplayUpdate)
	var names []string
	var vals []int
	for name, val := range stats {
		names = append(names, stateInfoMap[name].shortName)
		vals = append(vals, val)
	}
	sd.barChart.DataLabels = names
	sd.barChart.Data = vals
}
