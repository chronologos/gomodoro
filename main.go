package main

import (
	"fmt"

	ui "github.com/gizak/termui"
)

func main() {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()
	p := makePomodoro()
	l := makeHelpBox()
	draw := func() {
		ui.Render(p.gauge, l)
	}
	p.gauge.Handle("/sys/kbd/s", func(ui.Event) {
		p.tState = started
	})
	p.gauge.Handle("/sys/kbd/p", func(ui.Event) {
		p.tState = paused
	})
	p.gauge.Handle("/sys/kbd/r", func(ui.Event) {
		p.t = 0
		p.gauge.Percent = int(float64(p.maxDuration-p.t) / float64(p.maxDuration) * 100)
		p.gauge.BorderLabel = fmt.Sprintf("Time Elapsed: %s", p.t.String())
	})
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/timer/1s", func(e ui.Event) {
		draw()
	})

	ui.Loop()
}
