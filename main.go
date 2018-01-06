package main

import (
	"errors"
	"flag"
	"time"

	"github.com/0xAX/notificator"
	ui "github.com/gizak/termui"
)

var notify *notificator.Notificator

var stats = statsDisplayT{
	work:      0,
	shortRest: 0,
	longRest:  0,
}

var timeMapx = map[pomodoroState]*time.Duration{
	work:      workPeriod,
	shortRest: shortRestPeriod,
	longRest:  longRestPeriod,
}

type stateStrings struct {
	longName      string
	shortName     string
	completionMsg string
	period        time.Duration
}

type stateInfoMapT map[pomodoroState]stateStrings

// find searches for exact matches in
func (s stateInfoMapT) findShortName(sn string) (pomodoroState, error) {
	for k, v := range s {
		if v.shortName == sn {
			return k, nil
		}
	}
	return work, errors.New("Not found.")
}

var stateInfoMap = stateInfoMapT{
	work:      stateStrings{"Pomodoro", "p", "Go take a break!", time.Second},
	shortRest: stateStrings{"Short rest", "sr", "Get back to work!", time.Second},
	longRest:  stateStrings{"Long rest", "lr", "Get back to work!", time.Second},
}

var debugDisplay = makeTextBox(0, 10, 25, 5, "debug") // TODO(iantay) remove

func main() {
	flag.Parse()
	for ps, ss := range stateInfoMap {
		if ps == work {
			ss.period = *workPeriod
		} else if ps == shortRest {
			ss.period = *shortRestPeriod
		} else if ps == longRest {
			ss.period = *longRestPeriod
		}
	}
	notify = notificator.New(notificator.Options{
		DefaultIcon: "gomodoro-small.png",
		AppName:     "gomodoro",
	})

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	sd := makeStatsDisplay(51, 0, 16, 10)
	p := makePomodoro(0, 7, 50, 3)
	helpStrs := []string{"[s] start", "[p] pause", "[r] reset"}
	mtb := makeTextBox(0, 2, 25, 5, "help")
	mtb.updateText(helpStrs)

	draw := func() {
		ui.Render(p.gauge, mtb.list, sd.barChart, debugDisplay.list)
	}

	ui.Handle("/timer/1s", func(e ui.Event) {
		draw()
	})

	defineHandlers(p, sd)
	ui.Loop()
}
