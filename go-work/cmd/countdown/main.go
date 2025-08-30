package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"example.com/workspace/lib/args"
	u "example.com/workspace/lib/utils"
)

const (
	MaxSeconds = 215999 // 59:59:59 in seconds
)

func main() {
	d := u.Must1(args.Parse())
	u.Expect(d.Len == 1, "Usage: countdown <time>\n  Examples:\n    countdown 30      (30 seconds)\n    countdown 5m      (5 minutes)\n    countdown 1h      (1 hour)\n    countdown 05:30   (5 minutes 30 seconds)\n    countdown 1:05:30 (1 hour 5 minutes 30 seconds)")

	inputTime := d.Args[0]
	totalSeconds := parseTimeInput(inputTime)
	validateSeconds(totalSeconds)
	runCountdown(totalSeconds)
}

func parseTimeInput(input string) int {
	// Try short format first (e.g., 5m, 30s, 1h)
	if seconds, ok := parseShortFormat(input); ok {
		return seconds
	}

	// Try time format (e.g., 05:30, 1:05:30)
	if seconds, ok := parseTimeFormat(input); ok {
		return seconds
	}

	// Try raw seconds
	if seconds, err := strconv.Atoi(input); err == nil {
		return seconds
	}

	u.Expect(false, "Invalid time format. Use formats like: 30 (seconds), 5m, 1h, 05:30, or 1:05:30")
	return 0
}

func parseShortFormat(input string) (int, bool) {
	shortTimeRegex := regexp.MustCompile(`^([0-9]+)([hHmMsS])$`)
	matches := shortTimeRegex.FindStringSubmatch(input)

	if len(matches) != 3 {
		return 0, false
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, false
	}

	unit := strings.ToLower(matches[2])
	switch unit {
	case "h":
		return value * 3600, true
	case "m":
		return value * 60, true
	case "s":
		return value, true
	default:
		return 0, false
	}
}

func parseTimeFormat(input string) (int, bool) {
	timeRegex := regexp.MustCompile(`^([0-5]?[0-9]):([0-5]?[0-9])(?::([0-5]?[0-9]))?$`)
	matches := timeRegex.FindStringSubmatch(input)

	if len(matches) < 3 {
		return 0, false
	}

	// Parse components
	var hours, minutes, seconds int
	var err error

	if len(matches) == 4 && matches[3] != "" {
		// Format: h:m:s
		hours, err = strconv.Atoi(matches[1])
		if err != nil {
			return 0, false
		}
		minutes, err = strconv.Atoi(matches[2])
		if err != nil {
			return 0, false
		}
		seconds, err = strconv.Atoi(matches[3])
		if err != nil {
			return 0, false
		}
	} else {
		// Format: m:s
		minutes, err = strconv.Atoi(matches[1])
		if err != nil {
			return 0, false
		}
		seconds, err = strconv.Atoi(matches[2])
		if err != nil {
			return 0, false
		}
	}

	return hours*3600 + minutes*60 + seconds, true
}

func validateSeconds(totalSeconds int) {
	u.Expect(totalSeconds > 0 && totalSeconds <= MaxSeconds,
		fmt.Sprintf("Bad argument.\nTry with (hh:)?mm:ss (min 1, max 59:59:59)\nOR {s : s ∈ Z and 1 ≤ s ≤ %d}\nOR ^n[hHmMsS]$ (min 1s, max 59h)", MaxSeconds))
}

func runCountdown(totalSeconds int) {
	for i := totalSeconds; i > 0; i-- {
		timeStr := formatTime(i)
		fmt.Printf("\r⏳ %s   ", timeStr) // extra spaces keep output clean
		time.Sleep(1 * time.Second)
	}

	// Clear the countdown line and show completion message
	fmt.Print("\r")

	// Play notification sound in background
	go playNotificationSound()

	// Show completion message
	customTimeMessage := getCustomTimeMessage(totalSeconds)
	currentTime := time.Now().Format("15:04:05")
	fmt.Printf("✅ The timer for %s was up at %s\n", customTimeMessage, currentTime)

	// Show system notification
	showNotification(customTimeMessage)
}

func formatTime(totalSeconds int) string {
	hours := totalSeconds / 3600
	minutes := (totalSeconds / 60) % 60
	seconds := totalSeconds % 60

	if hours == 0 {
		if minutes == 0 {
			return fmt.Sprintf("%d", seconds)
		}
		return fmt.Sprintf("%02d:%02d", minutes, seconds)
	}

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func getCustomTimeMessage(totalSeconds int) string {
	if totalSeconds == 1 {
		return "1 second"
	} else if totalSeconds < 60 {
		return fmt.Sprintf("%d seconds", totalSeconds)
	} else {
		return formatTime(totalSeconds)
	}
}

func playNotificationSound() {
	soundFile := os.Getenv("HOME") + "/.zsh-spell-book/src/media/sounds/xylofon.wav"
	exec.Command("afplay", soundFile).Run()
}

func showNotification(timeMessage string) {
	message := fmt.Sprintf("The timer for %s is over", timeMessage)

	// Try macOS notification
	if _, err := exec.LookPath("osascript"); err == nil {
		script := fmt.Sprintf(`display notification "%s" with title "Countdown Timer"`, message)
		exec.Command("osascript", "-e", script).Run()
		return
	}

	// Try Linux notification
	if _, err := exec.LookPath("notify-send"); err == nil {
		exec.Command("notify-send", "Countdown Timer", message).Run()
		return
	}
}
