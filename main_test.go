package main

import "testing"

func TestReadFileExists(t *testing.T) {
	filename := "test.lock"
	pid := 123

	// Setup test file
	if err := createFile(filename, pid); err != nil {
		t.Fatal(err)
	}

	result := readFile(filename)
	if result != pid {
		t.Fatalf("Expected %d, got %d", pid, result)
	}

	// Cleanup
	err := cleanupFiles([]string{filename})
	if err != nil {
		t.Fatalf("Failed to cleanup test files: %s, due to error: %s", filename, err)
	}

}

func TestReadFileMissing(t *testing.T) {
	filename := "doesnt.exist"
	pid := 0

	result := readFile(filename)
	if result != pid {
		t.Fatalf("Expected %d, got %d", pid, result)
	}

}

func TestFileExists(t *testing.T) {
	filename := "test.lock"
	pid := 123

	// Setup test file
	if err := createFile(filename, pid); err != nil {
		t.Fatal(err)
	}

	if !fileExists(filename) {
		t.Fatalf("Expected file \"%s\" to exist", filename)
	}

	// Cleanup
	err := cleanupFiles([]string{filename})
	if err != nil {
		t.Fatalf("Failed to cleanup test files: %s, due to error: %s", filename, err)
	}

}

func TestResolveFilesWithPidfile(t *testing.T) {
	pid := 123
	lockfile := "test.lock"
	pidfile := "test.pid"

	// Setup test lockfile
	if err := createFile(lockfile, pid); err != nil {
		t.Fatal(err)
	}
	// Setup test pidfile
	if err := createFile(pidfile, pid); err != nil {
		t.Fatal(err)
	}

	resultPid, resultFiles := resolveFiles(Opts{pidfile, lockfile, true, "", []string{}})

	if resultPid != pid {
		t.Fatalf("Expected %d, got %d", pid, resultPid)
	}

	if resultFiles[0] != pidfile && resultFiles[1] != lockfile {
		t.Fatalf("Expected files: %s, %s. Got: %s, %s", pidfile, lockfile, resultFiles[0], resultFiles[1])
	}

	// Cleanup
	err := cleanupFiles([]string{lockfile, pidfile})
	if err != nil {
		t.Fatalf("Failed to cleanup test files, due to error: %s", err)
	}
}

func TestResolveFilesWithoutPidfile(t *testing.T) {
	pid := 123
	lockfile := "test.lock"

	// Setup test lockfile
	if err := createFile(lockfile, pid); err != nil {
		t.Fatal(err)
	}

	resultPid, resultFiles := resolveFiles(Opts{"", lockfile, false, "", []string{}})

	if resultPid != pid {
		t.Fatalf("Expected %d, got %d", pid, resultPid)
	}

	if resultFiles[0] != lockfile {
		t.Fatalf("Expected file: %s. Got: %s", lockfile, resultFiles[0])
	}

	if len(resultFiles) > 1 {
		t.Fatalf("Expected 1 file, got %d", len(resultFiles))
	}

	// Cleanup
	err := cleanupFiles([]string{lockfile})
	if err != nil {
		t.Fatalf("Failed to cleanup test files, due to error: %s", err)
	}
}
