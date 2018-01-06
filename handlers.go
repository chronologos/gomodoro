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

	ui.Handle("/sys/kbd/n", func(ui.Event) {
		p.nextState()
		p.maxDuration = stateInfoMap[p.stateSeq[p.pStateIdx]].period
		p.render()
	})

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/v", func(ui.Event) {
		writeStat(0, time.Second*5, 5)
	})

	ui.Handle("/sys/kbd/c", func(e ui.Event) {
		readStats()
		sd.refreshStatsDisplay(e)
	})

	ui.Handle("/gomodoro/sdupdate", sd.refreshStatsDisplay)

}
