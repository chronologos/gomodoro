package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// logger functions log and read pomodoro stats from a file, specified by the flag statsFileName.

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// checkD is Debug-only version of check()
func checkD(e error) {
	if e != nil {
		log.Print(e.Error)
	}
}

func isToday(t time.Time) bool {
	var year, day, nowYear, nowDay int
	var month, nowMonth time.Month
	year, month, day = t.Date()
	nowYear, nowMonth, nowDay = time.Now().Date()
	if year != nowYear || day != nowDay || month != nowMonth {
		return false
	}
	return true
}

func allDates(t time.Time) bool {
	return true
}

type dateChecker func(time.Time) bool

func genericLoadStats(r *csv.Reader, f dateChecker) {
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		checkD(err)
		var state pomodoroState
		var stateCount int
		for i, x := range record {
			switch i {
			case 0:
				d, err := time.Parse("Jan 2 15:04:05 2006", x)
				checkD(err)
				// Only load stats from today
				if !isToday(d) {
					continue
				}
			case 1:
				state, err = stateInfoMap.findShortName(x)
				if err != nil {
					continue
				}
			case 3:
				count, err := strconv.Atoi(x)
				if err != nil {
					continue
				} else {
					stateCount = count
				}
			}
		}
		stats[state] += stateCount
	}
}

func loadStats(r *csv.Reader) {
	genericLoadStats(r, isToday)
}

// readStats() reads from file into stats variable.
func readStats() {
	f, err := os.Open(*statFileName)
	defer f.Close()
	checkD(err)
	r := csv.NewReader(bufio.NewReader(f))
	r.FieldsPerRecord = 4
	loadStats(r)
}

// writeStats() writes from stats variable into file.
// TODO(chronologos) currently the whole file is read into memory
func writeStats() {
	return
}
func writeStat(ps pomodoroState, l time.Duration, num int) {
	record := []string{time.Now().Format("Jan 2 15:04:05 2006"), stateInfoMap[ps].shortName, l.String(), strconv.Itoa(num)}
	//for _, r := range record {
	//fmt.Print(r + " ")
	//}
	f, err := os.OpenFile(*statFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	checkD(err)
	w := csv.NewWriter(f)
	err = w.Write(record)
	checkD(err)
	w.Flush()
	err = w.Error()
	checkD(err)
}
