package main

import "time"
import ui "github.com/gizak/termui"

func defineHandlers(p *pomodoro, sd *statsDisplay) {
	ui.Handle("/sys/kbd/s", func(ui.Event) {
		p.start()
	})

	ui.Handle("/sys/kbd/p", func(ui.Event) {
		p.pause()
	})

	ui.Handle("/sys/kbd/r", func(ui.Event) {
		p.reset()
	})

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/v", func(ui.Event) {
		writeStats(0, time.Second*5, 5)
	})

	ui.Handle("/sys/kbd/c", func(ui.Event) {
		readStats()
	})

	ui.Handle("/gomodoro/sdupdate", sd.refreshStatsDisplay)

}
