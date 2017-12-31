package main

import (
	"flag"
	"time"

	"github.com/0xAX/notificator"
	ui "github.com/gizak/termui"
)

var notify *notificator.Notificator
var workPeriod = flag.Duration("work", 10*time.Second, "work period")
var shortRestPeriod = flag.Duration("rest", 5*time.Second, "rest period")
var longRestPeriod = flag.Duration("long_rest", 6*time.Second, "long rest period")
var statFileName = flag.String("stats_file", "stats.csv", "name of csv file in which stats should be/are tracked")
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

var stateInfoMap = map[pomodoroState]stateStrings{
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
	helpStrs := []string{"[s] start timer", "[p] pause timer", "[r] reset timer"}
	mtb := makeTextBox(0, 2, 25, 5, "help")
	mtb.updateText(helpStrs)
	draw := func() {
		ui.Render(p.gauge, mtb.list, sd.barChart, debugDisplay.list)
	}

	ui.Handle("/sys/kbd/s", func(ui.Event) {
		p.tState = started
	})

	ui.Handle("/sys/kbd/p", func(ui.Event) {
		p.tState = paused
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

	ui.Handle("/timer/1s", func(e ui.Event) {
		draw()
	})

	ui.Handle("/gomodoro/sdupdate", sd.refreshStatsDisplay)

	ui.Loop()
}
