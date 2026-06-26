package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AestheticAutonomy/justctx/internal/doctor"
)

func TestRunDoctorCmd_HealthyProject(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ".jctx", "rules"), 0755)

	var buf bytes.Buffer
	allPass, err := runDoctorCmd(root, "", false, &buf)
	if err != nil {
		t.Fatalf("runDoctorCmd: %v", err)
	}
	if !allPass {
		t.Errorf("expected all checks to pass:\n%s", buf.String())
	}
	if !strings.Contains(buf.String(), "OK") {
		t.Errorf("expected OK in output:\n%s", buf.String())
	}
}

func TestRunDoctorCmd_MissingJctx(t *testing.T) {
	root := t.TempDir()
	// No .jctx/ directory

	var buf bytes.Buffer
	allPass, err := runDoctorCmd(root, "", false, &buf)
	if err != nil {
		t.Fatalf("runDoctorCmd: %v", err)
	}
	if allPass {
		t.Error("expected failure when .jctx/ is missing")
	}
	if !strings.Contains(buf.String(), "FAIL") {
		t.Errorf("expected FAIL in output:\n%s", buf.String())
	}
}

func TestRunDoctorCmd_JSON(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ".jctx", "rules"), 0755)

	var buf bytes.Buffer
	_, err := runDoctorCmd(root, "", true, &buf)
	if err != nil {
		t.Fatalf("runDoctorCmd JSON: %v", err)
	}

	var res doctor.Result
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	if len(res.Checks) == 0 {
		t.Error("expected at least one check in JSON result")
	}
}
