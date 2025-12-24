package main

import (
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

	// Find sections by level
	levelCounts := make(map[int]int)
	for _, sec := range app.Sections {
		levelCounts[sec.Level]++
	}

	if levelCounts[1] != 1 {
		t.Errorf("Expected 1 level-1 section, got %d", levelCounts[1])
	}

	if levelCounts[2] < 2 {
		t.Errorf("Expected at least 2 level-2 sections, got %d", levelCounts[2])
	}

	if levelCounts[3] < 3 {
		t.Errorf("Expected at least 3 level-3 sections, got %d", levelCounts[3])
	}
}

func TestGetCurrentSection(t *testing.T) {
	app := createTestApp()

	sec := app.GetCurrentSection()
	if sec == nil {
		t.Fatal("Expected current section, got nil")
	}

	if sec.Title != "Main Title" {
		t.Errorf("Expected title 'Main Title', got '%s'", sec.Title)
	}

	// Test with empty sections
	emptyApp := NewApp()
	if emptyApp.GetCurrentSection() != nil {
		t.Error("Expected nil for empty sections")
	}
}

func TestNextSection(t *testing.T) {
	app := createTestApp()

	// Move to next
	if !app.NextSection() {
		t.Error("Expected NextSection to return true")
	}

	if app.CurrentIdx != 1 {
		t.Errorf("Expected CurrentIdx 1, got %d", app.CurrentIdx)
	}

	// Move to end
	for app.NextSection() {
	}

	lastIdx := len(app.Sections) - 1
	if app.CurrentIdx != lastIdx {
		t.Errorf("Expected CurrentIdx %d, got %d", lastIdx, app.CurrentIdx)
	}

	// Try to go past end
	if app.NextSection() {
		t.Error("Expected NextSection to return false at end")
	}
}

func TestPrevSection(t *testing.T) {
	app := createTestApp()

	// At beginning, should return false
	if app.PrevSection() {
		t.Error("Expected PrevSection to return false at beginning")
	}

	// Move forward then back
	app.NextSection()
	app.NextSection()

	if !app.PrevSection() {
		t.Error("Expected PrevSection to return true")
	}

	if app.CurrentIdx != 1 {
		t.Errorf("Expected CurrentIdx 1, got %d", app.CurrentIdx)
	}
}

func TestGotoSection(t *testing.T) {
	app := createTestApp()

	tests := []struct {
		name     string
		idx      int
		expected bool
	}{
		{"valid index 0", 0, true},
		{"valid index 2", 2, true},
		{"negative index", -1, false},
		{"index too large", 1000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.GotoSection(tt.idx)
			if result != tt.expected {
				t.Errorf("GotoSection(%d) = %v, expected %v", tt.idx, result, tt.expected)
			}
			if result && app.CurrentIdx != tt.idx {
				t.Errorf("CurrentIdx = %d, expected %d", app.CurrentIdx, tt.idx)
			}
		})
	}
}

func TestSearchSections(t *testing.T) {
	app := createTestApp()

	tests := []struct {
		name     string
		query    string
		minCount int
	}{
		{"search title", "Chapter", 2},
		{"search content", "task", 2},
		{"case insensitive", "CHAPTER", 2},
		{"no results", "xyznonexistent", 0},
		{"search giai đoạn", "Giai đoạn", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := app.SearchSections(tt.query)
			if len(matches) < tt.minCount {
				t.Errorf("SearchSections(%q) returned %d results, expected at least %d",
					tt.query, len(matches), tt.minCount)
			}
		})
	}
}

func TestGetCheckboxLines(t *testing.T) {
	app := createTestApp()

	// Go to Chapter 1 which has checkboxes
	for i, sec := range app.Sections {
		if sec.Title == "Chapter 1: Basics" {
			app.CurrentIdx = i
			break
		}
	}

	checkboxes := app.GetCheckboxLines()

	if len(checkboxes) != 3 {
		t.Errorf("Expected 3 checkboxes, got %d", len(checkboxes))
	}
}

func TestGetCheckboxLinesEmpty(t *testing.T) {
	app := createTestApp()

	// Go to main title which has no checkboxes
	app.CurrentIdx = 0

	checkboxes := app.GetCheckboxLines()

	if len(checkboxes) != 0 {
		t.Errorf("Expected 0 checkboxes for main title, got %d", len(checkboxes))
	}
}

func TestToggleCheckbox(t *testing.T) {
	app := createTestApp()

	// Go to Chapter 1
	for i, sec := range app.Sections {
		if sec.Title == "Chapter 1: Basics" {
			app.CurrentIdx = i
			break
		}
	}

	checkboxes := app.GetCheckboxLines()
	if len(checkboxes) == 0 {
		t.Fatal("No checkboxes found")
	}

	// Get initial state
	sec := app.GetCurrentSection()
	initialContent := sec.Content
	hadUnchecked := strings.Contains(initialContent, "- [ ]")

	// Toggle first checkbox
	lineIdx := checkboxes[0]
	result := app.ToggleCheckbox(lineIdx)

	if !result {
		t.Error("Expected ToggleCheckbox to return true")
	}

	// Verify toggle happened
	newSec := app.GetCurrentSection()
	if hadUnchecked {
		// Should now have one more checked
		if strings.Count(newSec.Content, "- [x]") <= strings.Count(initialContent, "- [x]") {
			t.Error("Expected checkbox to be checked")
		}
	}
}

func TestToggleCheckboxInvalidLine(t *testing.T) {
	app := createTestApp()

	// Invalid line index
	if app.ToggleCheckbox(-1) {
		t.Error("Expected false for negative line index")
	}

	if app.ToggleCheckbox(10000) {
		t.Error("Expected false for line index out of bounds")
	}
}

func TestAddNote(t *testing.T) {
	app := createTestApp()

	sec := app.GetCurrentSection()
	initialLen := len(sec.Content)

	app.AddNote("Test note content")

	newSec := app.GetCurrentSection()
	if len(newSec.Content) <= initialLen {
		t.Error("Expected content to grow after adding note")
	}

	if !strings.Contains(newSec.Content, "Test note content") {
		t.Error("Expected note content to be present")
	}

	if !strings.Contains(newSec.Content, "Ghi chú") {
		t.Error("Expected 'Ghi chú' label in note")
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

func TestGetProgress(t *testing.T) {
	app := createTestApp()

	// Find Chapter 1 which has 2 unchecked, 1 checked
	for i, sec := range app.Sections {
		if sec.Title == "Chapter 1: Basics" {
			checked, total := app.GetProgress(i)
			if total != 3 {
				t.Errorf("Expected 3 total checkboxes, got %d", total)
			}
			if checked != 1 {
				t.Errorf("Expected 1 checked checkbox, got %d", checked)
			}
			return
		}
	}

	t.Error("Chapter 1 not found")
}

func TestGetProgressInvalidIndex(t *testing.T) {
	app := createTestApp()

	checked, total := app.GetProgress(-1)
	if checked != 0 || total != 0 {
		t.Errorf("Expected (0, 0) for invalid index, got (%d, %d)", checked, total)
	}

	checked, total = app.GetProgress(10000)
	if checked != 0 || total != 0 {
		t.Errorf("Expected (0, 0) for index out of bounds, got (%d, %d)", checked, total)
	}
}

func TestGetTotalProgress(t *testing.T) {
	app := createTestApp()

	checked, total := app.GetTotalProgress()

	if total < 5 {
		t.Errorf("Expected at least 5 total checkboxes, got %d", total)
	}

	if checked < 3 {
		t.Errorf("Expected at least 3 checked checkboxes, got %d", checked)
	}

	if checked > total {
		t.Errorf("Checked (%d) cannot be greater than total (%d)", checked, total)
	}
}

// ============================================================================
// RenderLine Tests
// ============================================================================

func TestRenderLineCheckboxUnchecked(t *testing.T) {
	result := RenderLine("- [ ] Task", 80)

	if strings.Contains(result, "- [ ]") {
		t.Error("Expected '- [ ]' to be replaced")
	}

	if !strings.Contains(result, "☐") {
		t.Error("Expected '☐' symbol in output")
	}
}

func TestRenderLineCheckboxChecked(t *testing.T) {
	result := RenderLine("- [x] Done", 80)

	if strings.Contains(result, "- [x]") {
		t.Error("Expected '- [x]' to be replaced")
	}

	if !strings.Contains(result, "☑") {
		t.Error("Expected '☑' symbol in output")
	}
}

func TestRenderLineBold(t *testing.T) {
	result := RenderLine("This is **bold** text", 80)

	if strings.Contains(result, "**") {
		t.Error("Expected ** markers to be removed")
	}

	if !strings.Contains(result, Bold) {
		t.Error("Expected bold ANSI code in output")
	}
}

func TestRenderLineCode(t *testing.T) {
	result := RenderLine("Use `code` here", 80)

	if strings.Contains(result, "`") {
		t.Error("Expected backticks to be removed")
	}

	if !strings.Contains(result, Cyan) {
		t.Error("Expected cyan ANSI code for code")
	}
}

func TestRenderLineBullet(t *testing.T) {
	result := RenderLine("- Item", 80)

	if !strings.Contains(result, "•") {
		t.Error("Expected bullet symbol '•' in output")
	}
}

func TestRenderLineNumberedList(t *testing.T) {
	result := RenderLine("1. First item", 80)

	if !strings.Contains(result, Cyan) {
		t.Error("Expected cyan ANSI code for number")
	}
}

func TestRenderLineBlockquote(t *testing.T) {
	result := RenderLine("> Quote text", 80)

	if !strings.Contains(result, "│") {
		t.Error("Expected blockquote marker '│' in output")
	}

	if !strings.Contains(result, Dim) {
		t.Error("Expected dim ANSI code for blockquote")
	}
}

func TestRenderLineHorizontalRule(t *testing.T) {
	result := RenderLine("---", 80)

	if !strings.Contains(result, "─") {
		t.Error("Expected horizontal line character in output")
	}
}

func TestRenderLinePreservesIndentation(t *testing.T) {
	result := RenderLine("   - Indented item", 80)

	// Should preserve leading spaces
	if !strings.HasPrefix(result, "   ") {
		t.Error("Expected indentation to be preserved")
	}
}

// ============================================================================
// Renderer Tests
// ============================================================================

func TestNewRenderer(t *testing.T) {
	app := createTestApp()
	renderer := NewRenderer(app)

	if renderer.App != app {
		t.Error("Expected renderer to have same app reference")
	}

	if renderer.TermWidth != app.TermWidth {
		t.Errorf("Expected TermWidth %d, got %d", app.TermWidth, renderer.TermWidth)
	}

	if renderer.ScrollOffset != 0 {
		t.Errorf("Expected ScrollOffset 0, got %d", renderer.ScrollOffset)
	}

	if renderer.PageSize < 15 {
		t.Errorf("Expected PageSize >= 15, got %d", renderer.PageSize)
	}
}

func TestRendererResetScroll(t *testing.T) {
	app := createTestApp()
	renderer := NewRenderer(app)

	renderer.ScrollOffset = 10
	renderer.ResetScroll()

	if renderer.ScrollOffset != 0 {
		t.Errorf("Expected ScrollOffset 0 after reset, got %d", renderer.ScrollOffset)
	}
}

func TestRendererScrollDown(t *testing.T) {
	app := createTestApp()
	app.TermHeight = 20 // Small height to force pagination
	renderer := NewRenderer(app)

	// Find a section with lots of content
	for i, sec := range app.Sections {
		lines := strings.Split(sec.Content, "\n")
		if len(lines) > 20 {
			app.CurrentIdx = i
			break
		}
	}

	initialOffset := renderer.ScrollOffset
	canScroll := renderer.ScrollDown()

	if canScroll && renderer.ScrollOffset <= initialOffset {
		t.Error("Expected ScrollOffset to increase when scrolling down")
	}
}

func TestRendererScrollUp(t *testing.T) {
	app := createTestApp()
	renderer := NewRenderer(app)

	// First scroll down
	renderer.ScrollOffset = 10

	// Then scroll up
	canScroll := renderer.ScrollUp()

	if !canScroll {
		t.Error("Expected ScrollUp to return true when not at top")
	}

	if renderer.ScrollOffset >= 10 {
		t.Error("Expected ScrollOffset to decrease after scrolling up")
	}
}

func TestRendererScrollUpAtTop(t *testing.T) {
	app := createTestApp()
	renderer := NewRenderer(app)

	// Already at top
	renderer.ScrollOffset = 0
	canScroll := renderer.ScrollUp()

	if canScroll {
		t.Error("Expected ScrollUp to return false when at top")
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

func TestMin(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{1, 2, 1},
		{2, 1, 1},
		{5, 5, 5},
		{-1, 1, -1},
		{0, 0, 0},
	}

	for _, tt := range tests {
		result := min(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("min(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
		}
	}
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestNavigationFlow(t *testing.T) {
	app := createTestApp()

	// Simulate navigation: next -> next -> prev -> goto
	app.NextSection()
	app.NextSection()

	if app.CurrentIdx != 2 {
		t.Errorf("After 2 nexts, expected idx 2, got %d", app.CurrentIdx)
	}

	app.PrevSection()

	if app.CurrentIdx != 1 {
		t.Errorf("After prev, expected idx 1, got %d", app.CurrentIdx)
	}

	app.GotoSection(4)

	if app.CurrentIdx != 4 {
		t.Errorf("After goto(4), expected idx 4, got %d", app.CurrentIdx)
	}
}

func TestCheckboxWorkflow(t *testing.T) {
	app := createTestApp()

	// Find section with checkboxes
	for i, sec := range app.Sections {
		if strings.Contains(sec.Content, "- [ ]") {
			app.CurrentIdx = i
			break
		}
	}

	// Get initial progress
	initialChecked, total := app.GetProgress(app.CurrentIdx)

	// Toggle an unchecked box
	checkboxes := app.GetCheckboxLines()
	for _, lineIdx := range checkboxes {
		lines := strings.Split(app.GetCurrentSection().Content, "\n")
		if strings.Contains(lines[lineIdx], "- [ ]") {
			app.ToggleCheckbox(lineIdx)
			break
		}
	}

	// Verify progress changed
	newChecked, newTotal := app.GetProgress(app.CurrentIdx)

	if newTotal != total {
		t.Errorf("Total changed from %d to %d", total, newTotal)
	}

	if newChecked != initialChecked+1 {
		t.Errorf("Expected checked to increase by 1, was %d now %d", initialChecked, newChecked)
	}
}

func TestSearchAndGoto(t *testing.T) {
	app := createTestApp()

	// Search for "Advanced"
	matches := app.SearchSections("Advanced")

	if len(matches) == 0 {
		t.Fatal("Expected to find 'Advanced' section")
	}

	// Goto first match
	app.GotoSection(matches[0])

	sec := app.GetCurrentSection()
	if !strings.Contains(strings.ToLower(sec.Title), "advanced") &&
		!strings.Contains(strings.ToLower(sec.Content), "advanced") {
		t.Error("Expected current section to contain 'advanced'")
	}
}

// ============================================================================
// Edge Case Tests
// ============================================================================

func TestEmptyFile(t *testing.T) {
	app := NewApp()
	app.FileContent = ""
	app.FileLines = []string{}
	app.ParseSections()

	if len(app.Sections) != 0 {
		t.Errorf("Expected 0 sections for empty file, got %d", len(app.Sections))
	}

	if app.GetCurrentSection() != nil {
		t.Error("Expected nil current section for empty file")
	}
}

func TestFileWithOnlyHeaders(t *testing.T) {
	app := NewApp()
	app.FileContent = "# Title\n## Subtitle\n### Sub-subtitle"
	app.FileLines = strings.Split(app.FileContent, "\n")
	app.ParseSections()

	if len(app.Sections) != 3 {
		t.Errorf("Expected 3 sections, got %d", len(app.Sections))
	}

	// All should have empty content
	for _, sec := range app.Sections {
		if strings.TrimSpace(sec.Content) != "" {
			t.Errorf("Expected empty content for section '%s', got '%s'", sec.Title, sec.Content)
		}
	}
}

func TestSpecialCharactersInContent(t *testing.T) {
	app := NewApp()
	app.FileContent = "# Title\n\nContent with émojis 🎉 and spëcial çhars"
	app.FileLines = strings.Split(app.FileContent, "\n")
	app.ParseSections()

	sec := app.GetCurrentSection()
	if !strings.Contains(sec.Content, "🎉") {
		t.Error("Expected emoji to be preserved")
	}
}

func TestMultipleCheckboxesOnSameLine(t *testing.T) {
	// This is unusual but should be handled
	line := "- [ ] First - [ ] Second"
	result := RenderLine(line, 80)

	// Should replace at least the first one
	if strings.Contains(result, "- [ ]") && !strings.Contains(result, "☐") {
		t.Error("Expected at least one checkbox to be rendered")
	}
}

// ============================================================================
// Benchmark Tests
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

func BenchmarkGetTotalProgress(b *testing.B) {
	app := createTestApp()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.GetTotalProgress()
	}
}
