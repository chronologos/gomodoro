package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
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

func loadStats(r *csv.Reader) (v statsDisplayT) {
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("%d\n", len(record))
		var state pomodoroState
		var stateCount int
		for i, x := range record {
			fmt.Println(x)
			switch i {
			case 0:
				// Mon Jan 2 15:04:05 -0700 MST 2006
				d, err := time.Parse("Jan 2 15:04:05 2006", x)
				if err != nil {
					fmt.Print(err)
					fmt.Print(" gg ")
				}
				// Only load stats from today
				if isToday(d) {
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
		fmt.Println("---")
		// fmt.Println(record)
	}
	v = make(statsDisplayT)
	return v
}

// statsTracker reads and writes stats to a file.
func readStats() {
	f, err := os.Open(*statFileName)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	r := csv.NewReader(bufio.NewReader(f))
	r.FieldsPerRecord = 4
	stats = loadStats(r)
}

func writeStats(ps pomodoroState, l time.Duration, num int) {
	record := []string{time.Now().Format("Jan 2 15:04:05 2006"), stateInfoMap[ps].shortName, l.String(), strconv.Itoa(num)}
	for _, r := range record {
		fmt.Print(r + " ")
	}
	f, err := os.OpenFile(*statFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	w := csv.NewWriter(f)
	err = w.Write(record)
	if err != nil {
		fmt.Println(err)
	}
	w.Flush()
	err = w.Error()
	check(err)
	fmt.Print("stats written...")
}
