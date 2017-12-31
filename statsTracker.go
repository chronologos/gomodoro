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

// statsTracker reads and writes stats to a file.
func readStats() {
	f, err := os.Open(*statFileName)
	defer f.Close()
	check(err)
	r := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record)
	}
}

func writeStats(ps pomodoroState, l time.Duration, num int) {
	record := []string{time.Now().Format("Jan 2 15:04:05 MST 2006"), nameStateMap[ps], l.String(), strconv.Itoa(num)}
	for _, r := range record {
		fmt.Print(r + " ")
	}
	f, err := os.OpenFile(*statFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	check(err)
	w := csv.NewWriter(f)
	err = w.Write(record)
	check(err)
	w.Flush()
	err = w.Error()
	check(err)
	fmt.Print("stats written...")
}
