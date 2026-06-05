package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitWritesToFile(t *testing.T) {
	// Run the test inside a temp dir so we don't pollute ./log.
	tmp := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	const name = "test_experiment"
	if err := Init(name); err != nil {
		t.Fatalf("Init: %v", err)
	}

	Println("hello", "world")
	Printf("answer=%d\n", 42)
	Print("done\n")

	if err := Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	logPath := filepath.Join(tmp, LogDir, name+".log")
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile %s: %v", logPath, err)
	}
	out := string(data)

	for _, want := range []string{"hello world", "answer=42", "done", "Experiment      : " + name} {
		if !strings.Contains(out, want) {
			t.Errorf("log file missing %q\n--- log ---\n%s", want, out)
		}
	}
}

func TestInitRejectsEmptyName(t *testing.T) {
	if err := Init(""); err == nil {
		t.Fatal("expected error for empty experiment name, got nil")
	}
}

func TestPrintBeforeInitGoesToStdout(t *testing.T) {
	// Just make sure the package-level functions don't panic when used
	// before Init. We can't easily capture stdout in a unit test without
	// extra plumbing, so this just exercises the code path.
	Print("")
	Println()
	Printf("%s", "")
}
