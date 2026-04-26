package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// appLogger mirrors messages to stdout and a log file for daemon mode.
type appLogger struct {
	path   string
	stdout bool
}

func newLogger(path string, stdout bool) *appLogger {
	return &appLogger{path: strings.TrimSpace(path), stdout: stdout}
}

func (l *appLogger) Printf(format string, args ...any) {
	line := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
	if l.stdout {
		fmt.Println(line)
	}
	if l.path == "" {
		return
	}
	if err := ensureParentDir(l.path); err != nil {
		if l.stdout {
			fmt.Printf("log mkdir failed: %v\n", err)
		}
		return
	}
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		if l.stdout {
			fmt.Printf("log open failed: %v\n", err)
		}
		return
	}
	defer f.Close()
	_, _ = f.WriteString(line + "\n")
}
