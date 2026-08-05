package main

import (
	"os"

	"zsb_tmux_agent_notification/internal/notification"
)

func main() { os.Exit(notification.Run(os.Args[1:])) }
