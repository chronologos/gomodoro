package main

import (
	"flag"
	"fmt"
	"time"

	ui "github.com/gizak/termui"
)

var workPeriod = flag.Duration("work", 10*time.Second, "work period")
var shortRestPeriod = flag.Duration("rest", 5*time.Second, "rest period")
var longRestPeriod = flag.Duration("long_rest", 6*time.Second, "long rest period")
var statFileName = flag.String("stats_file", "stats.csv", "name of csv file in which stats should be/are tracked")
var statsGlobal = statsDisplayT{
	work:      0,
	shortRest: 0,
	longRest:  0,
}

var timeMap = map[pomodoroState]*time.Duration{
	work:      workPeriod,
	shortRest: shortRestPeriod,
	longRest:  longRestPeriod,
}

var nameStateMap = map[pomodoroState]string{
	work:      "p",
	shortRest: "sr",
	longRest:  "lr",
}

var debugDisplay = makeTextBox(0, 10, 25, 5, "debug") // TODO(iantay) remove

func main() {
	flag.Parse()
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()
	sd := makeStatsDisplay()
	p := makePomodoro(*sd)
	helpStrs := []string{"[s] start timer", "[p] pause timer", "[r] reset timer"}
	mtb := makeTextBox(0, 2, 25, 5, "help")
	mtb.updateText(helpStrs)
	draw := func() {
		ui.Render(p.gauge, mtb.list, sd.barChart, debugDisplay.list)
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

	ui.Handle("/sys/kbd/v", func(ui.Event) {
		writeStats(0, time.Second*5, 5)
	})

	ui.Handle("/sys/kbd/r", func(ui.Event) {
		readStats()
	})

	ui.Handle("/timer/1s", func(e ui.Event) {
		draw()
	})

	ui.Loop()
}
