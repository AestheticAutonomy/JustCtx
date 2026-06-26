package scanner

import (
	"testing"

	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

func TestDetectConflicts_DuplicateHeading(t *testing.T) {
	chunks := []schema.Chunk{
		{Content: "## Core Principles\nUse spaces.\n", SourceID: "s1"},
		{Content: "## Core Principles\nUse tabs.\n", SourceID: "s2"},
	}
	sources := []schema.Source{{ID: "s1"}, {ID: "s2"}}

	conflicts := DetectConflicts(chunks, sources)

	found := false
	for _, c := range conflicts {
		if c.Type == "duplicate_heading" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected duplicate_heading conflict, got: %v", conflicts)
	}
}

func TestDetectConflicts_NearDuplicateParagraph(t *testing.T) {
	// Two long paragraphs with >85% word overlap from different sources
	para1 := "Always write unit tests for all your functions and methods in the codebase"
	para2 := "Always write unit tests for all your functions and methods across the codebase"
	chunks := []schema.Chunk{
		{Content: para1, SourceID: "s1"},
		{Content: para2, SourceID: "s2"},
	}
	sources := []schema.Source{{ID: "s1"}, {ID: "s2"}}

	conflicts := DetectConflicts(chunks, sources)

	found := false
	for _, c := range conflicts {
		if c.Type == "near_duplicate_paragraph" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected near_duplicate_paragraph conflict, got: %v", conflicts)
	}
}

func TestDetectConflicts_ContradictingImperative(t *testing.T) {
	chunks := []schema.Chunk{
		{Content: "Always use spaces for indentation in all Go files\n", SourceID: "s1"},
		{Content: "Never use spaces for indentation in all Go files\n", SourceID: "s2"},
	}
	sources := []schema.Source{{ID: "s1"}, {ID: "s2"}}

	conflicts := DetectConflicts(chunks, sources)

	found := false
	for _, c := range conflicts {
		if c.Type == "contradicting_imperative" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected contradicting_imperative conflict, got: %v", conflicts)
	}
}

func TestDetectConflicts_NoFalsePositives(t *testing.T) {
	chunks := []schema.Chunk{
		{Content: "## Formatting\nUse gofmt for Go code formatting.\n", SourceID: "s1"},
		{Content: "## Testing\nWrite integration tests covering API endpoints.\n", SourceID: "s2"},
	}
	sources := []schema.Source{{ID: "s1"}, {ID: "s2"}}

	conflicts := DetectConflicts(chunks, sources)
	if len(conflicts) != 0 {
		t.Errorf("expected no conflicts for distinct content, got: %v", conflicts)
	}
}
