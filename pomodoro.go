package main

import (
	"fmt"
	"time"

	"github.com/0xAX/notificator"
	ui "github.com/gizak/termui"
)

// pomodoro implements the cli component that displays and counts down for
// the current pomodoro.

type timerState int

const (
	started timerState = iota
	paused
) // controls whether timer is running

type pomodoroState int

const (
	work pomodoroState = iota
	shortRest
	longRest
) // controls duration of timer

type pomodoro struct {
	gauge       *ui.Gauge
	t           time.Duration
	maxDuration time.Duration
	tState      timerState
	pStateIdx   int // points into stateSeq
	stateSeq    [4]pomodoroState
}

// We view pomodoro as n-state fsm, with states defined in main.go.
// This transitions the pomodoro state and performs side-effects.
func (p *pomodoro) stateTx() {
	stateName := stateInfoMap[p.stateSeq[p.pStateIdx]].longName
	title := fmt.Sprintf("%s Complete", stateName)
	text := stateInfoMap[p.stateSeq[p.pStateIdx]].completionMsg
	notify.Push(title, text, "gomodoro-small.png", notificator.UR_CRITICAL)
	fmt.Println("\a") // \a is the bell literal.
	debugDisplay.updateText([]string{stateInfoMap[p.stateSeq[p.pStateIdx]].shortName})
	stats[p.stateSeq[p.pStateIdx]]++
	ui.SendCustomEvt("/gomodoro/sdupdate", statsDisplayUpdate{})
	if p.pStateIdx >= cap(p.stateSeq)-1 {
		p.pStateIdx = 0
	} else {
		p.pStateIdx++
	}
	// newStateName := nameStateMap[p.stateSeq[p.pStateIdx]]
}

func (p *pomodoro) reset() {
	p.t = 0
	p.gauge.Percent = int(float64(p.maxDuration-p.t) / float64(p.maxDuration) * 100)
	p.gauge.BorderLabel = fmt.Sprintf("Time Elapsed: %s", p.t.String())
}

func makePomodoro(x, y, w, h int) *pomodoro {
	g := ui.NewGauge()
	g.X = x
	g.Y = y
	g.Width = w
	g.Height = h
	g.BorderLabel = "Duration || Elapsed: || pState: " // updated periodically
	g.BarColor = ui.ColorRed
	g.BorderFg = ui.ColorWhite
	g.BorderLabelFg = ui.ColorCyan
	var p = pomodoro{g, 0, *workPeriod, paused, 0, [4]pomodoroState{work, shortRest, work, longRest}}
	g.Handle("/timer/1s", func(e ui.Event) {
		if p.tState != started {
			return
		}
		_ = e.Data.(ui.EvtTimer)
		p.t += time.Duration(1) * time.Second
		g.Percent = int(float64(p.maxDuration-p.t) / float64(p.maxDuration) * 100)
		g.BorderLabel = fmt.Sprintf("Duration %s || Elapsed: %s || pState: %d", p.maxDuration, p.t.String(), p.stateSeq[p.pStateIdx])
		if p.t >= p.maxDuration {
			p.stateTx()
			p.tState = paused
			p.t = 0
			p.maxDuration = stateInfoMap[p.stateSeq[p.pStateIdx]].period
			return
		}

	})
	return &p
}
