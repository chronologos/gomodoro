package main

import (
	"flag"
	"time"
)

var workPeriod = flag.Duration("work", 10*time.Second, "work period")
var shortRestPeriod = flag.Duration("rest", 5*time.Second, "rest period")
var longRestPeriod = flag.Duration("long_rest", 6*time.Second, "long rest period")
var statFileName = flag.String("stats_file", "stats.csv", "name of csv file in which stats should be/are tracked")
