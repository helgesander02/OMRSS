// Package logger provides a lightweight execution-flow logger that mirrors
// fmt.Print* semantics while persisting output to a per-experiment log file.
//
// Typical usage:
//
//	cfg, _ := config.Load(*configFile)
//	if err := logger.Init(cfg.GetExperimentName()); err != nil {
//	    log.Fatalf("failed to initialise logger: %v", err)
//	}
//	defer logger.Close()
//
//	logger.Println("Hello, world!")
//	logger.Printf("Topology: %s\n", cfg.Network.Topology)
//
// All Print/Println/Printf calls are written both to stdout and to
// ./log/<experimentName>.log so existing terminal output stays intact while
// the full execution trace is captured on disk.
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogDir is the directory where log files are stored. It is created on Init
// if it does not already exist.
const LogDir = "./log"

var (
	mu     sync.Mutex
	file   *os.File
	writer io.Writer = os.Stdout
)

// Init opens (or creates) the log file ./log/<experimentName>.log and routes
// all logger output to both stdout and that file. It also redirects the
// standard library log package so that log.Print*/log.Fatal* are captured.
//
// Init may be called more than once; the previous log file (if any) is closed
// before a new one is opened.
func Init(experimentName string) error {
	if experimentName == "" {
		return fmt.Errorf("logger: experiment name must not be empty")
	}

	if err := os.MkdirAll(LogDir, 0o755); err != nil {
		return fmt.Errorf("logger: failed to create log directory %q: %w", LogDir, err)
	}

	logPath := filepath.Join(LogDir, experimentName+".log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("logger: failed to open log file %q: %w", logPath, err)
	}

	mu.Lock()
	defer mu.Unlock()

	if file != nil {
		_ = file.Close()
	}
	file = f
	writer = io.MultiWriter(os.Stdout, f)

	// Route the standard library log package into the same destination so
	// existing log.Printf / log.Fatalf calls in main.go also land in the
	// file.
	log.SetOutput(writer)

	// Write a small session header to make multiple runs easy to tell apart
	// when appending to an existing file.
	fmt.Fprintf(f, "\n========================================\n")
	fmt.Fprintf(f, "Session started : %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(f, "Experiment      : %s\n", experimentName)
	fmt.Fprintf(f, "========================================\n")

	return nil
}

// Close flushes and closes the underlying log file. Safe to call multiple
// times; subsequent calls are no-ops.
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if file == nil {
		return nil
	}

	err := file.Close()
	file = nil
	writer = os.Stdout
	log.SetOutput(os.Stderr) // restore log's default destination
	return err
}

// Writer returns the current logger writer. Before Init is called this is
// just os.Stdout; afterwards it is an io.MultiWriter targeting stdout and
// the active log file.
func Writer() io.Writer {
	mu.Lock()
	defer mu.Unlock()
	return writer
}

// Print writes its arguments to the logger using fmt.Print formatting rules.
func Print(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	_, _ = fmt.Fprint(writer, a...)
}

// Println writes its arguments followed by a newline using fmt.Println rules.
func Println(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	_, _ = fmt.Fprintln(writer, a...)
}

// Printf writes a formatted string to the logger using fmt.Printf rules.
func Printf(format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	_, _ = fmt.Fprintf(writer, format, a...)
}

// Fatal writes its arguments to the logger using fmt.Print rules and then
// terminates the program with os.Exit(1). The log file is flushed and closed
// before exit so the final message is not lost.
func Fatal(a ...any) {
	Print(a...)
	fatalExit()
}

// Fatalln writes its arguments to the logger using fmt.Println rules and then
// terminates the program with os.Exit(1).
func Fatalln(a ...any) {
	Println(a...)
	fatalExit()
}

// Fatalf writes a formatted string to the logger using fmt.Printf rules and
// then terminates the program with os.Exit(1).
func Fatalf(format string, a ...any) {
	Printf(format, a...)
	fatalExit()
}

// fatalExit closes the log file (if any) and exits the process with status 1.
// It must not be called with the package mutex held.
func fatalExit() {
	_ = Close()
	os.Exit(1)
}
