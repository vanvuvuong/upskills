package main

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

// ============================================================================
// Test fixtures
// ============================================================================

// sampleMarkdown provides test markdown content with various elements.
const sampleMarkdown = `# Main Title

Introduction text.

## Giai đoạn 1: Learning

Some content here.

### Chapter 1: Basics

- [ ] Task one
- [x] Task two completed
- [ ] Task three

**Bold text** and *italic text*.

### Chapter 2: Advanced

More content.

- [ ] Advanced task

## Giai đoạn 2: Practice

Practice section.

### Exercise 1

- [x] Done
- [x] Also done
`

// createTestApp creates an App with sample markdown loaded.
func createTestApp() *App {
	app := NewApp()
	app.FileContent = sampleMarkdown
	app.FileLines = strings.Split(sampleMarkdown, "\n")
	app.ParseSections()
	return app
}

// ============================================================================
// App Tests
// ============================================================================

func TestNewApp(t *testing.T) {
	app := NewApp()

	if app.FilePath != "learning-path-full.md" {
		t.Errorf("Expected default FilePath 'learning-path-full.md', got '%s'", app.FilePath)
	}

	if app.StateFile != ".sre-learn-state" {
		t.Errorf("Expected default StateFile '.sre-learn-state', got '%s'", app.StateFile)
	}

	if app.TermWidth != 80 {
		t.Errorf("Expected default TermWidth 80, got %d", app.TermWidth)
	}

	if app.TermHeight != 24 {
		t.Errorf("Expected default TermHeight 24, got %d", app.TermHeight)
	}

	if app.CurrentIdx != 0 {
		t.Errorf("Expected CurrentIdx 0, got %d", app.CurrentIdx)
	}
}

func TestParseSections(t *testing.T) {
	app := createTestApp()

	if len(app.Sections) == 0 {
		t.Fatal("Expected sections to be parsed, got 0")
	}

	// Check first section
	if app.Sections[0].Title != "Main Title" {
		t.Errorf("Expected first section title 'Main Title', got '%s'", app.Sections[0].Title)
	}

	if app.Sections[0].Level != 1 {
		t.Errorf("Expected first section level 1, got %d", app.Sections[0].Level)
	}
}

func TestParseSectionsLevels(t *testing.T) {
	app := createTestApp()

	// Should have multiple levels
	levels := make(map[int]bool)
	for _, sec := range app.Sections {
		levels[sec.Level] = true
	}

	if !levels[1] {
		t.Error("Expected level 1 sections")
	}
	if !levels[2] {
		t.Error("Expected level 2 sections")
	}
	if !levels[3] {
		t.Error("Expected level 3 sections")
	}
}

// ============================================================================
// Navigation Tests
// ============================================================================

func TestNextSection(t *testing.T) {
	app := createTestApp()
	app.CurrentIdx = 0

	app.NextSection()

	if app.CurrentIdx != 1 {
		t.Errorf("Expected CurrentIdx 1 after NextSection, got %d", app.CurrentIdx)
	}
}

func TestNextSectionAtEnd(t *testing.T) {
	app := createTestApp()
	app.CurrentIdx = len(app.Sections) - 1
	lastIdx := app.CurrentIdx

	app.NextSection()

	if app.CurrentIdx != lastIdx {
		t.Errorf("Expected CurrentIdx to stay at %d when at end, got %d", lastIdx, app.CurrentIdx)
	}
}

func TestPrevSection(t *testing.T) {
	app := createTestApp()
	app.CurrentIdx = 2

	app.PrevSection()

	if app.CurrentIdx != 1 {
		t.Errorf("Expected CurrentIdx 1 after PrevSection, got %d", app.CurrentIdx)
	}
}

func TestPrevSectionAtStart(t *testing.T) {
	app := createTestApp()
	app.CurrentIdx = 0

	app.PrevSection()

	if app.CurrentIdx != 0 {
		t.Errorf("Expected CurrentIdx to stay at 0 when at start, got %d", app.CurrentIdx)
	}
}

func TestGotoSection(t *testing.T) {
	app := createTestApp()

	app.GotoSection(3)

	if app.CurrentIdx != 3 {
		t.Errorf("Expected CurrentIdx 3 after GotoSection(3), got %d", app.CurrentIdx)
	}
}

func TestGotoSectionOutOfBounds(t *testing.T) {
	app := createTestApp()
	app.CurrentIdx = 0

	app.GotoSection(999)

	if app.CurrentIdx != 0 {
		t.Errorf("Expected CurrentIdx to stay at 0 for out of bounds, got %d", app.CurrentIdx)
	}

	app.GotoSection(-1)

	if app.CurrentIdx != 0 {
		t.Errorf("Expected CurrentIdx to stay at 0 for negative index, got %d", app.CurrentIdx)
	}
}

func TestGetCurrentSection(t *testing.T) {
	app := createTestApp()
	app.CurrentIdx = 0

	sec := app.GetCurrentSection()

	if sec.Title != "Main Title" {
		t.Errorf("Expected 'Main Title', got '%s'", sec.Title)
	}
}

// ============================================================================
// Search Tests
// ============================================================================

func TestSearchSections(t *testing.T) {
	app := createTestApp()

	results := app.SearchSections("Chapter")

	if len(results) == 0 {
		t.Fatal("Expected search results for 'Chapter'")
	}

	// Verify results contain Chapter
	for _, idx := range results {
		if !strings.Contains(strings.ToLower(app.Sections[idx].Title), "chapter") {
			t.Errorf("Search result '%s' doesn't contain 'chapter'", app.Sections[idx].Title)
		}
	}
}

func TestSearchSectionsCaseInsensitive(t *testing.T) {
	app := createTestApp()

	resultsLower := app.SearchSections("chapter")
	resultsUpper := app.SearchSections("CHAPTER")

	if len(resultsLower) != len(resultsUpper) {
		t.Error("Search should be case insensitive")
	}
}

func TestSearchSectionsNoResults(t *testing.T) {
	app := createTestApp()

	results := app.SearchSections("nonexistent12345")

	if len(results) != 0 {
		t.Errorf("Expected no results for nonexistent query, got %d", len(results))
	}
}

// ============================================================================
// Checkbox Tests
// ============================================================================

func TestToggleCheckbox(t *testing.T) {
	app := createTestApp()

	// Find a section with checkboxes
	for i, sec := range app.Sections {
		if strings.Contains(sec.Content, "- [ ]") {
			app.CurrentIdx = i
			break
		}
	}

	// Get actual checkbox line indices
	checkboxLines := app.GetCheckboxLines()
	if len(checkboxLines) == 0 {
		t.Skip("No checkboxes found in test content")
	}

	sec := app.GetCurrentSection()
	initialUnchecked := strings.Count(sec.Content, "- [ ]")

	// Toggle the first actual checkbox line
	app.ToggleCheckbox(checkboxLines[0])

	sec = app.GetCurrentSection()
	newUnchecked := strings.Count(sec.Content, "- [ ]")

	if newUnchecked >= initialUnchecked {
		t.Error("Expected checkbox to be toggled from unchecked to checked")
	}
}

func TestGetCheckboxLines(t *testing.T) {
	app := createTestApp()

	// Find a section with checkboxes
	for i, sec := range app.Sections {
		if strings.Contains(sec.Content, "- [ ]") || strings.Contains(sec.Content, "- [x]") {
			app.CurrentIdx = i
			break
		}
	}

	lines := app.GetCheckboxLines()

	if len(lines) == 0 {
		t.Error("Expected checkbox lines in test content")
	}
}

// ============================================================================
// Note Tests
// ============================================================================

func TestAddNote(t *testing.T) {
	app := createTestApp()
	app.CurrentIdx = 0

	sec := app.GetCurrentSection()
	initialContent := sec.Content

	app.AddNote("Test note content")

	sec = app.GetCurrentSection()

	if !strings.Contains(sec.Content, "Test note content") {
		t.Error("Expected note to be added to content")
	}

	if !strings.Contains(sec.Content, "**Ghi chú [") {
		t.Error("Expected note to have timestamp header")
	}

	if len(sec.Content) <= len(initialContent) {
		t.Error("Expected content to be longer after adding note")
	}
}

func TestAddNoteEmpty(t *testing.T) {
	app := createTestApp()

	sec := app.GetCurrentSection()
	initialContent := sec.Content

	app.AddNote("")

	newSec := app.GetCurrentSection()
	if newSec.Content != initialContent {
		t.Error("Expected no change for empty note")
	}
}

func TestExtractNotes(t *testing.T) {
	content := `Some content here.

> **Ghi chú [2025-01-01 10:00]:** First note
> continues here

More content.

> **Ghi chú [2025-01-02 11:00]:** Second note
`

	notes := extractNotes(content)

	if len(notes) != 2 {
		t.Errorf("Expected 2 notes, got %d", len(notes))
	}

	if len(notes) > 0 && !strings.Contains(notes[0], "First note") {
		t.Error("Expected first note to contain 'First note'")
	}
}

func TestExtractNotesEmpty(t *testing.T) {
	content := "Some content without any notes."

	notes := extractNotes(content)

	if len(notes) != 0 {
		t.Errorf("Expected 0 notes for content without notes, got %d", len(notes))
	}
}

func TestRemoveNoteFromContent(t *testing.T) {
	content := `Some content here.

> **Ghi chú [2025-01-01 10:00]:** First note

More content.

> **Ghi chú [2025-01-02 11:00]:** Second note
`

	noteToRemove := "> **Ghi chú [2025-01-01 10:00]:** First note"

	result := removeNoteFromContent(content, noteToRemove)

	if strings.Contains(result, "First note") {
		t.Error("Expected 'First note' to be removed")
	}

	if !strings.Contains(result, "Second note") {
		t.Error("Expected 'Second note' to remain")
	}

	if !strings.Contains(result, "Some content here") {
		t.Error("Expected other content to remain")
	}
}

// ============================================================================
// Progress Tests
// ============================================================================

func TestGetProgress(t *testing.T) {
	app := createTestApp()

	// Find section with checkboxes
	for i, sec := range app.Sections {
		if strings.Contains(sec.Content, "- [x]") {
			app.CurrentIdx = i
			checked, total := app.GetProgress(i)

			if total == 0 {
				t.Error("Expected total > 0 for section with checkboxes")
			}

			if checked == 0 {
				t.Error("Expected some checked items in test content")
			}

			if checked > total {
				t.Error("Checked cannot exceed total")
			}
			break
		}
	}
}

func TestGetTotalProgress(t *testing.T) {
	app := createTestApp()

	checked, total := app.GetTotalProgress()

	if total == 0 {
		t.Error("Expected total checkboxes > 0 in test content")
	}

	if checked > total {
		t.Error("Checked cannot exceed total")
	}
}

// ============================================================================
// State Persistence Tests
// ============================================================================

func TestSaveAndLoadState(t *testing.T) {
	app := createTestApp()
	app.StateFile = "/tmp/test-sre-state"
	app.CurrentIdx = 5

	// Clean up
	defer os.Remove(app.StateFile)

	// Save state
	err := app.SaveState(30)
	if err != nil {
		t.Fatalf("SaveState failed: %v", err)
	}

	// Create new app and load state
	app2 := NewApp()
	app2.StateFile = app.StateFile

	pageSize, err := app2.LoadState()
	if err != nil {
		t.Fatalf("LoadState failed: %v", err)
	}

	if app2.CurrentIdx != 5 {
		t.Errorf("Expected CurrentIdx 5, got %d", app2.CurrentIdx)
	}

	if pageSize != 30 {
		t.Errorf("Expected pageSize 30, got %d", pageSize)
	}
}

func TestLoadStateFileNotExists(t *testing.T) {
	app := NewApp()
	app.StateFile = "/tmp/nonexistent-state-file"

	pageSize, err := app.LoadState()

	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	if pageSize != 0 {
		t.Errorf("Expected pageSize 0 for non-existent file, got %d", pageSize)
	}
}

// ============================================================================
// Rendering Tests
// ============================================================================

func TestRenderLineCheckboxUnchecked(t *testing.T) {
	result := RenderLine("- [ ] Test item", 80)

	if !strings.Contains(result, "☐") {
		t.Error("Expected unchecked box symbol")
	}
}

func TestRenderLineCheckboxChecked(t *testing.T) {
	result := RenderLine("- [x] Completed item", 80)

	if !strings.Contains(result, "☑") {
		t.Error("Expected checked box symbol")
	}

	if !strings.Contains(result, Green) {
		t.Error("Expected green color for checked item")
	}
}

func TestRenderLineBold(t *testing.T) {
	result := RenderLine("Some **bold text** here", 80)

	if !strings.Contains(result, "bold text") {
		t.Error("Expected bold text to be preserved")
	}

	if !strings.Contains(result, Bold) {
		t.Error("Expected bold formatting")
	}
}

func TestRenderLineCode(t *testing.T) {
	result := RenderLine("Use `code here` for example", 80)

	if !strings.Contains(result, "code here") {
		t.Error("Expected code text to be preserved")
	}
}

func TestRenderLineBullet(t *testing.T) {
	result := RenderLine("- List item", 80)

	if !strings.Contains(result, "•") {
		t.Error("Expected bullet point")
	}
}

func TestRenderLineBlockquote(t *testing.T) {
	result := RenderLine("> Quoted text", 80)

	if !strings.Contains(result, "│") {
		t.Error("Expected blockquote indicator")
	}

	if !strings.Contains(result, Dim) {
		t.Error("Expected dim formatting for blockquote")
	}
}

// ============================================================================
// Renderer Tests
// ============================================================================

func TestNewRenderer(t *testing.T) {
	app := createTestApp()
	renderer := NewRenderer(app)

	if renderer.App != app {
		t.Error("Expected renderer to reference app")
	}

	if renderer.PageSize < 15 {
		t.Errorf("Expected PageSize >= 15, got %d", renderer.PageSize)
	}

	if renderer.ScrollOffset != 0 {
		t.Errorf("Expected initial ScrollOffset 0, got %d", renderer.ScrollOffset)
	}
}

func TestRendererScrollDown(t *testing.T) {
	app := NewApp()
	// Create content with many lines to ensure scrolling works
	var longContent strings.Builder
	longContent.WriteString("# Test\n\n")
	for i := 0; i < 100; i++ {
		longContent.WriteString("Line " + strconv.Itoa(i) + "\n")
	}
	app.FileContent = longContent.String()
	app.FileLines = strings.Split(app.FileContent, "\n")
	app.ParseSections()

	renderer := NewRenderer(app)
	renderer.PageSize = 10 // Small page size to ensure we can scroll

	initialOffset := renderer.ScrollOffset
	scrolled := renderer.ScrollDown()

	if !scrolled {
		t.Error("Expected ScrollDown to return true")
	}

	if renderer.ScrollOffset <= initialOffset {
		t.Error("Expected ScrollOffset to increase after ScrollDown")
	}
}

func TestRendererScrollUp(t *testing.T) {
	app := createTestApp()
	renderer := NewRenderer(app)

	renderer.ScrollOffset = 10
	renderer.ScrollUp()

	if renderer.ScrollOffset >= 10 {
		t.Error("Expected ScrollOffset to decrease after ScrollUp")
	}
}

func TestRendererScrollUpAtTop(t *testing.T) {
	app := createTestApp()
	renderer := NewRenderer(app)

	renderer.ScrollOffset = 0
	renderer.ScrollUp()

	if renderer.ScrollOffset != 0 {
		t.Error("Expected ScrollOffset to stay at 0 when already at top")
	}
}

func TestRendererResetScroll(t *testing.T) {
	app := createTestApp()
	renderer := NewRenderer(app)

	renderer.ScrollOffset = 100
	renderer.ResetScroll()

	if renderer.ScrollOffset != 0 {
		t.Errorf("Expected ScrollOffset 0 after reset, got %d", renderer.ScrollOffset)
	}
}

func TestRendererAdjustPageSize(t *testing.T) {
	app := createTestApp()
	renderer := NewRenderer(app)

	initialSize := renderer.PageSize

	// Increase - no upper limit
	renderer.AdjustPageSize(50)
	if renderer.PageSize != initialSize+50 {
		t.Errorf("Expected PageSize %d after +50, got %d", initialSize+50, renderer.PageSize)
	}

	// Decrease back
	renderer.AdjustPageSize(-50)
	if renderer.PageSize != initialSize {
		t.Errorf("Expected PageSize %d after -50, got %d", initialSize, renderer.PageSize)
	}
}

func TestRendererAdjustPageSizeMinimum(t *testing.T) {
	app := createTestApp()
	renderer := NewRenderer(app)

	// Try to decrease below minimum
	renderer.PageSize = 5
	renderer.AdjustPageSize(-10)

	if renderer.PageSize != 5 {
		t.Errorf("PageSize should not go below 5, got %d", renderer.PageSize)
	}
}

// ============================================================================
// Utility Tests
// ============================================================================

func TestClearScreen(t *testing.T) {
	// ClearScreen just prints escape codes, hard to test
	// Just verify it doesn't panic
	ClearScreen()
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestNavigationFlow(t *testing.T) {
	app := createTestApp()

	// Start at beginning
	if app.CurrentIdx != 0 {
		t.Fatal("Expected to start at index 0")
	}

	// Navigate forward
	app.NextSection()
	app.NextSection()
	if app.CurrentIdx != 2 {
		t.Errorf("Expected index 2 after two NextSection calls, got %d", app.CurrentIdx)
	}

	// Navigate back
	app.PrevSection()
	if app.CurrentIdx != 1 {
		t.Errorf("Expected index 1 after PrevSection, got %d", app.CurrentIdx)
	}

	// Jump to specific
	app.GotoSection(0)
	if app.CurrentIdx != 0 {
		t.Errorf("Expected index 0 after GotoSection(0), got %d", app.CurrentIdx)
	}
}

func TestCheckboxWorkflow(t *testing.T) {
	app := createTestApp()

	// Find section with unchecked items
	for i, sec := range app.Sections {
		if strings.Contains(sec.Content, "- [ ]") {
			app.CurrentIdx = i

			// Get actual checkbox line indices
			checkboxLines := app.GetCheckboxLines()
			if len(checkboxLines) == 0 {
				continue
			}

			// Get initial state
			checked1, total1 := app.GetProgress(i)

			// Toggle the first actual checkbox
			app.ToggleCheckbox(checkboxLines[0])

			// Verify change
			checked2, total2 := app.GetProgress(i)

			if total2 != total1 {
				t.Error("Total should not change after toggle")
			}

			if checked2 == checked1 {
				t.Error("Checked count should change after toggle")
			}
			return
		}
	}
	t.Skip("No section with checkboxes found")
}

func TestSearchAndGoto(t *testing.T) {
	app := createTestApp()

	results := app.SearchSections("Exercise")

	if len(results) > 0 {
		app.GotoSection(results[0])
		sec := app.GetCurrentSection()

		if !strings.Contains(strings.ToLower(sec.Title), "exercise") {
			t.Error("Expected to navigate to Exercise section")
		}
	}
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestEmptyFile(t *testing.T) {
	app := NewApp()
	app.FileContent = ""
	app.FileLines = []string{}
	app.ParseSections()

	if len(app.Sections) != 0 {
		t.Errorf("Expected 0 sections for empty file, got %d", len(app.Sections))
	}
}

func TestSpecialCharactersInContent(t *testing.T) {
	app := NewApp()
	app.FileContent = "# Test\n\nContent with émojis 🎉 and Việt Nam"
	app.FileLines = strings.Split(app.FileContent, "\n")
	app.ParseSections()

	if len(app.Sections) == 0 {
		t.Fatal("Expected section to be parsed")
	}

	if !strings.Contains(app.Sections[0].Content, "🎉") {
		t.Error("Expected emoji to be preserved")
	}

	if !strings.Contains(app.Sections[0].Content, "Việt Nam") {
		t.Error("Expected Vietnamese characters to be preserved")
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkParseSections(b *testing.B) {
	app := NewApp()
	app.FileContent = sampleMarkdown
	app.FileLines = strings.Split(sampleMarkdown, "\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ParseSections()
	}
}

func BenchmarkRenderLine(b *testing.B) {
	line := "- [ ] **Bold task** with `code` and *italic*"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RenderLine(line, 80)
	}
}

func BenchmarkSearchSections(b *testing.B) {
	app := createTestApp()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.SearchSections("Chapter")
	}
}

func BenchmarkToggleCheckbox(b *testing.B) {
	app := createTestApp()
	app.CurrentIdx = 2 // Section with checkboxes

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ToggleCheckbox(0)
	}
}
