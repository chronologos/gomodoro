package main

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui"
)

type pomodoroState int

const (
	started pomodoroState = iota
	paused
)

type pomodoro struct {
	gauge       *ui.Gauge
	t           time.Duration
	maxDuration time.Duration
	state       pomodoroState
}

func makePomodoro(maxDuration time.Duration) *pomodoro {
	g := ui.NewGauge()
	g.Percent = 100 // updated periodically
	g.Width = 50
	g.Height = 3
	g.Y = 11
	g.BorderLabel = "Time Left: " // updated periodically
	g.BarColor = ui.ColorRed
	g.BorderFg = ui.ColorWhite
	g.BorderLabelFg = ui.ColorCyan
	var p = pomodoro{g, 0, maxDuration, paused}

	g.Handle("/timer/1s", func(e ui.Event) {
		if p.state != started {
			return
		}
		_ = e.Data.(ui.EvtTimer)
		p.t += time.Duration(1) * time.Second
		g.Percent = int(float64(maxDuration-p.t) / float64(maxDuration) * 100)
		g.BorderLabel = fmt.Sprintf("Time Elapsed: %s", p.t.String())
	})
	g.Handle("/sys/kbd/s", func(ui.Event) {
		p.state = started
	})
	g.Handle("/sys/kbd/p", func(ui.Event) {
		p.state = paused
	})
	g.Handle("/sys/kbd/r", func(ui.Event) {
		p.t = 0
		g.Percent = int(float64(maxDuration-p.t) / float64(maxDuration) * 100)
		g.BorderLabel = fmt.Sprintf("Time Elapsed: %s", p.t.String())
	})
	return &p
}

func main() {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()
	p := makePomodoro(time.Duration(25) * time.Second) // TODO(iantay) change duration later.
	draw := func() {
		ui.Render(p.gauge)
	}

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/timer/1s", func(e ui.Event) {
		draw()
	})

	ui.Loop()
}
