package daemon

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/kurocifer/rivulet/rivulet-base/serverUtils"
	"github.com/kurocifer/rivulet/rivulet-cli/pkg/utils"
	"github.com/sevlyar/go-daemon"
)

const (
	pidFile = "rivulet.pid" // File that holds the daemon's PID.
	logFile = "rivulet.log" // File to which daemon writes it's logs.
)

// StartDaemon, starts up a daemon that runs the server
func StartDaemon() {
	workDir := utils.GetWorkDir()

	cntxt := &daemon.Context{
		PidFileName: workDir + "/" + pidFile,
		PidFilePerm: 0644, // PID file permissions (read/write for owner, read for others)

		LogFileName: workDir + "/" + logFile,
		LogFilePerm: 0640, // Log file permissions (read/write for owner, read for group)

		WorkDir: workDir,
		Umask:   027,
		Args:    os.Args,
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalf("Error starting rivulet daemon: %v", err)
	}

	if d != nil {
		log.Printf("rivulet started with PID: %d", d.Pid)
		return
	}

	log.Printf("Daemon reborn with PID: %d", os.Getegid())

	defer func() {
		if err := cntxt.Release(); err != nil {
			log.Printf("Error releasing daemon: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	server := serverUtils.MakeServer(":3000")
	server.Start()

	s := <-sig

	log.Printf("Revieved signal: %d, shutting down...", s)
	log.Println("Daemon exited")
}

// StopDaemon, ends the server daemon
func StopDaemon() {
	pidFileLoc := utils.GetWorkDir() + "/" + pidFile
	// Get PID from file
	pidBytes, err := os.ReadFile(pidFileLoc)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Rivulet daemon is not running (no PID file related to it exits).")
			return
		}
		log.Fatalf("Error reading PID file %s: %d", pidFileLoc, err)
	}

	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		log.Fatalf("Error parsing PID from file %s: %v", pidFileLoc, err)
	}

	// Find the process
	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Fatalf("Error finding rivulet process with PID %d: %v", pid, err)
	}

	// Send SIGTERM signal
	err = proc.Signal(syscall.SIGTERM)
	if err != nil {
		if err.Error() == "os: process already finished" || err.Error() == "no such process" {
			fmt.Printf("Rivulet daemon (PID %d) is already stopped or not running.\n", pid)
			// Clean up stale PID file if process is gone
			if err := os.Remove(pidFile); err != nil && !os.IsNotExist(err) {
				log.Printf("Warning: Could not remove stale PID file %s: %v", pidFile, err)
			}
			return
		}
		log.Fatalf("Error sending SIGTERM to PID %d: %v", pid, err)
	}

	log.Printf("Sent SIGTERM to rivulet daemon (PID %d). Waiting for it to stop...\n", pid)

	// Ensure the PID file has disappeared before we exit
	timeout := 10 * time.Second
	startTime := time.Now()

	for {
		_, err := os.Stat(pidFileLoc)
		if os.IsNotExist(err) {
			log.Println("Rivulet daemon stopped successfully")
			return
		}
		if time.Since(startTime) > timeout {
			fmt.Printf("Rivulet daemon (PID %d) did not stop within %s timeout. PID file still exists.\n", pid, timeout)
			return
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// GetDaemonStatus, checks if the rivulet daemon is still running
func GetDaemonStatus() {
	cntxt := daemon.Context{
		PidFileName: utils.GetWorkDir() + "/" + pidFile,
	}

	pid, err := cntxt.Search()
	if err != nil {
		log.Println("rivulet daemon is not running")
	} else {
		log.Printf("rivulet daemon running with PID: %d", pid.Pid)
	}
}
