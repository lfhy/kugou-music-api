package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type daemonInfo struct {
	PID       int    `json:"pid"`
	StartedAt string `json:"started_at"`
}

func daemonState(pidPath string) (daemonInfo, bool) {
	info, err := readDaemonInfo(pidPath)
	if err != nil || info.PID <= 0 {
		return daemonInfo{}, false
	}
	return info, !processExists(info.PID)
}

func readDaemonInfo(pidPath string) (daemonInfo, error) {
	body, err := os.ReadFile(pidPath)
	if err != nil {
		return daemonInfo{}, err
	}
	var info daemonInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return daemonInfo{}, err
	}
	return info, nil
}

func writeDaemonInfo(pidPath string) error {
	info := daemonInfo{PID: os.Getpid(), StartedAt: time.Now().Format(time.RFC3339)}
	if err := ensureParentDir(pidPath); err != nil {
		return err
	}
	body, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(pidPath, body, 0o644)
}

func acquireDaemonLock(pidPath string) error {
	state, stale := daemonState(pidPath)
	if state.PID > 0 && !stale {
		return fmt.Errorf("instance already running: pid=%d", state.PID)
	}
	if stale {
		_ = os.Remove(pidPath)
	}
	return writeDaemonInfo(pidPath)
}

func releaseDaemonLock(pidPath string) {
	state, err := readDaemonInfo(pidPath)
	if err != nil {
		_ = os.Remove(pidPath)
		return
	}
	if state.PID == os.Getpid() {
		_ = os.Remove(pidPath)
	}
}

func processExists(pid int) bool {
	if pid <= 0 {
		return false
	}
	return syscall.Kill(pid, 0) == nil
}

func runDaemonCommand(paths appPaths, interval time.Duration) error {
	state, stale := daemonState(paths.PID)
	if state.PID > 0 && !stale {
		return fmt.Errorf("daemon already running: pid=%d", state.PID)
	}
	if err := ensureParentDir(paths.Log); err != nil {
		return err
	}
	logFile, err := os.OpenFile(paths.Log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer logFile.Close()

	cmd, err := buildDetachedStartCommand(paths, interval)
	if err != nil {
		return err
	}
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Stdin = nil
	if err := cmd.Start(); err != nil {
		return err
	}

	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		state, stale = daemonState(paths.PID)
		if state.PID > 0 && !stale {
			fmt.Printf("daemon started: pid=%d\n", state.PID)
			fmt.Printf("log: %s\n", paths.Log)
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("daemon start timed out, check log: %s", paths.Log)
}

func buildDetachedStartCommand(paths appPaths, interval time.Duration) (*exec.Cmd, error) {
	args := []string{
		"start",
		"-config", paths.Config,
		"-source-session-file", paths.SourceSession,
		"-log", paths.Log,
		"-pid", paths.PID,
		"-interval", interval.String(),
	}
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(exe, args...)
	cmd.Env = append(os.Environ(), "KUGOU_AUTO_CHECKIN_DAEMON=1")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	return cmd, nil
}

func stopDaemon(pidPath string) error {
	state, stale := daemonState(pidPath)
	if state.PID <= 0 {
		return fmt.Errorf("daemon not running")
	}
	if stale {
		_ = os.Remove(pidPath)
		fmt.Println("stale pid file removed")
		return nil
	}
	if err := syscall.Kill(state.PID, syscall.SIGTERM); err != nil {
		return err
	}
	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		if !processExists(state.PID) {
			_ = os.Remove(pidPath)
			fmt.Printf("daemon stopped: pid=%d\n", state.PID)
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("daemon still running: pid=%d", state.PID)
}
