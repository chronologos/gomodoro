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

func (p *pomodoro) nextState() {
	if p.pStateIdx >= cap(p.stateSeq)-1 {
		p.pStateIdx = 0
	} else {
		p.pStateIdx++
	}
}
func (p *pomodoro) getStateLongname() string {
	return stateInfoMap[p.stateSeq[p.pStateIdx]].longName
}

// We view pomodoro as n-state fsm, with states defined in main.go.
// This transitions the pomodoro state and performs side-effects.
func (p *pomodoro) stateTx() {
	stateName := p.getStateLongname()
	title := fmt.Sprintf("%s Complete", stateName)
	text := stateInfoMap[p.stateSeq[p.pStateIdx]].completionMsg
	notify.Push(title, text, "gomodoro-small.png", notificator.UR_CRITICAL)
	fmt.Println("\a") // \a is the bell literal.
	debugDisplay.updateText([]string{stateInfoMap[p.stateSeq[p.pStateIdx]].shortName})
	stats[p.stateSeq[p.pStateIdx]]++
	ui.SendCustomEvt("/gomodoro/sdupdate", statsDisplayUpdate{})
	// newStateName := nameStateMap[p.stateSeq[p.pStateIdx]]
	p.nextState()
}

func (p *pomodoro) start() {
	p.tState = started
	p.render()
}

func (p *pomodoro) pause() {
	p.tState = paused
	p.render()
}

func (p *pomodoro) reset() {
	p.t = 0
	p.render()
}

func (p *pomodoro) render() {
	p.gauge.Percent = int(float64(p.maxDuration-p.t) / float64(p.maxDuration) * 100)
	stateName := p.getStateLongname()
	p.gauge.BorderLabel = fmt.Sprintf("%s of %s || Elapsed: %s", p.maxDuration, stateName, p.t.String())
}

func (p *pomodoro) handleTick() {
	if p.tState != started {
		return
	}
	p.t += time.Duration(1) * time.Second
	if p.t >= p.maxDuration {
		p.stateTx()
		p.pause()
		p.reset()
		p.maxDuration = stateInfoMap[p.stateSeq[p.pStateIdx]].period
		return
	}

}

func makePomodoro(x, y, w, h int) *pomodoro {
	g := ui.NewGauge()
	g.X = x
	g.Y = y
	g.Width = w
	g.Height = h
	g.BorderLabel = ""
	g.BarColor = ui.ColorRed
	g.BorderFg = ui.ColorWhite
	g.BorderLabelFg = ui.ColorCyan
	var p = pomodoro{g, 0, 0, paused, 0, [4]pomodoroState{work, shortRest, work, longRest}}
	p.maxDuration = stateInfoMap[p.stateSeq[p.pStateIdx]].period
	p.render()
	return &p
}
