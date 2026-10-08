package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

var version = "dev"

var keyActions = map[string]func() error{
	"f15":    sendF15,
	"alttab": sendAltTab,
}

func main() {
	interval := 10 * time.Minute
	var until string
	keyName := "f15"
	flag.DurationVar(&interval, "interval", interval, "")
	flag.DurationVar(&interval, "i", interval, "")
	flag.StringVar(&until, "until", "", "")
	flag.StringVar(&until, "u", "", "")
	flag.StringVar(&keyName, "key", keyName, "")
	flag.StringVar(&keyName, "k", keyName, "")
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "")
	flag.BoolVar(&showVersion, "v", false, "")
	flag.Usage = printUsage

	if wantsHelp(os.Args[1:]) {
		printUsage()
		return
	}
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		return
	}

	if interval <= 0 {
		log.Fatal("interval must be positive")
	}

	keyName = strings.ToLower(keyName)
	sendKey, ok := keyActions[keyName]
	if !ok {
		log.Fatalf("unknown -key %q: want f15 or alttab", keyName)
	}

	var deadline <-chan time.Time
	if until != "" {
		stopAt, err := parseUntil(until, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		deadline = time.After(time.Until(stopAt))
		log.Printf("will stop at %s", stopAt.Format("15:04"))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("sending %s every %s; press Ctrl+C to stop", keyName, interval)
	for {
		select {
		case <-ticker.C:
			if err := sendKey(); err != nil {
				log.Printf("%s failed: %v", keyName, err)
				continue
			}
			log.Printf("%s sent", keyName)
		case <-deadline:
			log.Print("reached stop time, exiting")
			return
		case <-stop:
			log.Print("interrupted, exiting")
			return
		}
	}
}

func parseUntil(s string, now time.Time) (time.Time, error) {
	t, err := time.ParseInLocation("15:04", s, now.Location())
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid -until %q: want HH:MM", s)
	}
	stopAt := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
	if !stopAt.After(now) {
		return time.Time{}, fmt.Errorf("-until %s is already past", s)
	}
	return stopAt, nil
}

func wantsHelp(args []string) bool {
	for _, a := range args {
		switch a {
		case "?", "/?", "-?", "--?":
			return true
		}
	}
	return false
}

func printUsage() {
	fmt.Fprint(flag.CommandLine.Output(), `vdi_poke - sends a keypress on a timer to keep the session active

Usage: vdi_poke.exe [options]

Options:
  -i, -interval <duration>  time between keypresses, e.g. 10m, 90s (default 10m)
  -u, -until <HH:MM>        stop at this local time (24h), e.g. 17:30
  -k, -key <f15|alttab>     key to send (default f15)
  -v, -version              print the version
  -h, -?, /?, ?             show this help

Examples:
  vdi_poke.exe
  vdi_poke.exe -i 5m
  vdi_poke.exe -i 5m -u 17:30
  vdi_poke.exe -k alttab

Press Ctrl+C to stop.
`)
}
