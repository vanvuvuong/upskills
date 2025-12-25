// Package main implements an interactive CLI tool for the SRE Learning Path.
//
// This tool provides a terminal-based markdown viewer with navigation,
// checkbox toggling, note-taking, and progress tracking capabilities.
// It renders markdown with ANSI colors and supports keyboard navigation.
//
// Usage:
//
//	go build -o sre-learn .
//	./sre-learn
//
// The tool expects a file named "learning-path-full.md" in the current directory.
//
// Keyboard shortcuts:
//
// Content navigation:
//   - j/↓: Scroll down within section
//   - k/↑: Scroll up within section
//
// Section navigation:
//   - n: Next section
//   - p: Previous section
//   - Enter: Next section
//   - t: Open interactive TOC
//   - g: Go to section by number
//   - G: Go to last section
//   - /: Search sections
//
// Features:
//   - x: Toggle checkbox
//   - a: Add note
//   - s: Save file
//
// Display:
//   - +: Increase visible lines
//   - -: Decrease visible lines
//   - ?: Show help
//   - q: Quit
package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//go:embed templates/default.md
var defaultTemplate string

// ANSI escape codes for terminal styling.
// These constants provide color and formatting for terminal output.
const (
	// Text formatting
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	// Foreground colors
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// Background colors
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

// Section represents a markdown section parsed from the document.
// Each section corresponds to a header (# through ####) and its content.
type Section struct {
	// Title is the text after the # symbols
	Title string
	// Content is all text until the next header
	Content string
	// Level indicates header depth (1 = #, 2 = ##, etc.)
	Level int
	// Line is the line number in the source file (0-indexed)
	Line int
}

// App holds the application state.
// It encapsulates all mutable state for easier testing and management.
type App struct {
	// Sections contains all parsed markdown sections
	Sections []Section
	// CurrentIdx is the index of the currently displayed section
	CurrentIdx int
	// FilePath is the path to the markdown file
	FilePath string
	// FileContent is the raw file content
	FileContent string
	// FileLines is the file split by newlines
	FileLines []string
	// TermWidth is the terminal width in columns
	TermWidth int
	// TermHeight is the terminal height in rows
	TermHeight int
	// StateFile is the path to save/load state
	StateFile string
}

// NewApp creates a new App instance with default values.
// It initializes terminal dimensions and sets the default file path.
func NewApp() *App {
	return &App{
		FilePath:   "learning-path-full.md",
		StateFile:  ".sre-learn-state",
		TermWidth:  80,
		TermHeight: 24,
	}
}

// SaveState saves current reading position and settings to state file.
func (a *App) SaveState(pageSize int) error {
	content := fmt.Sprintf("current_section=%d\npage_size=%d\nfile_path=%s\n",
		a.CurrentIdx, pageSize, a.FilePath)
	return os.WriteFile(a.StateFile, []byte(content), 0o644)
}

// LoadState restores reading position and settings from state file.
// Returns (pageSize, error). If file doesn't exist, returns defaults.
func (a *App) LoadState() (int, error) {
	data, err := os.ReadFile(a.StateFile)
	if err != nil {
		return 0, err // File doesn't exist, use defaults
	}

	pageSize := 0
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]
		switch key {
		case "current_section":
			if idx, err := strconv.Atoi(value); err == nil {
				a.CurrentIdx = idx
			}
		case "page_size":
			if ps, err := strconv.Atoi(value); err == nil {
				pageSize = ps
			}
		case "file_path":
			// Only use saved file_path if current one is default
			if a.FilePath == "learning-path-full.md" && value != "" {
				a.FilePath = value
			}
		}
	}

	return pageSize, nil
}

// LoadFile reads the markdown file into memory.
// It populates FileContent and FileLines fields.
// Returns an error if the file cannot be read.
func (a *App) LoadFile() error {
	data, err := os.ReadFile(a.FilePath)
	if err != nil {
		return fmt.Errorf("cannot read file %s: %w", a.FilePath, err)
	}
	a.FileContent = string(data)
	a.FileLines = strings.Split(a.FileContent, "\n")
	return nil
}

// ParseSections extracts sections from the loaded markdown content.
// A section starts with a header (# to ####) and includes all content
// until the next header of any level.
func (a *App) ParseSections() {
	a.Sections = []Section{}
	var currentSection *Section
	var contentLines []string

	headerRegex := regexp.MustCompile(`^(#{1,4})\s+(.+)$`)

	for i, line := range a.FileLines {
		if matches := headerRegex.FindStringSubmatch(line); matches != nil {
			// Save previous section
			if currentSection != nil {
				currentSection.Content = strings.Join(contentLines, "\n")
				a.Sections = append(a.Sections, *currentSection)
			}

			// Start new section
			currentSection = &Section{
				Title: matches[2],
				Level: len(matches[1]),
				Line:  i,
			}
			contentLines = []string{}
		} else if currentSection != nil {
			contentLines = append(contentLines, line)
		}
	}

	// Save last section
	if currentSection != nil {
		currentSection.Content = strings.Join(contentLines, "\n")
		a.Sections = append(a.Sections, *currentSection)
	}
}

// GetCurrentSection returns the currently selected section.
// Returns nil if no sections exist or index is out of bounds.
func (a *App) GetCurrentSection() *Section {
	if len(a.Sections) == 0 || a.CurrentIdx < 0 || a.CurrentIdx >= len(a.Sections) {
		return nil
	}
	return &a.Sections[a.CurrentIdx]
}

// NextSection moves to the next section if possible.
// Returns true if the move was successful, false if already at the end.
func (a *App) NextSection() bool {
	if a.CurrentIdx < len(a.Sections)-1 {
		a.CurrentIdx++
		return true
	}
	return false
}

// PrevSection moves to the previous section if possible.
// Returns true if the move was successful, false if already at the beginning.
func (a *App) PrevSection() bool {
	if a.CurrentIdx > 0 {
		a.CurrentIdx--
		return true
	}
	return false
}

// GotoSection moves to the section at the given index.
// Returns true if the index is valid, false otherwise.
func (a *App) GotoSection(idx int) bool {
	if idx >= 0 && idx < len(a.Sections) {
		a.CurrentIdx = idx
		return true
	}
	return false
}

// SearchSections finds all sections matching the query string.
// The search is case-insensitive and matches both title and content.
// Returns a slice of indices for matching sections.
func (a *App) SearchSections(query string) []int {
	query = strings.ToLower(query)
	matches := []int{}

	for i, sec := range a.Sections {
		if strings.Contains(strings.ToLower(sec.Title), query) ||
			strings.Contains(strings.ToLower(sec.Content), query) {
			matches = append(matches, i)
		}
	}

	return matches
}

// GetCheckboxLines returns the line indices of all checkboxes in the current section.
// A checkbox is either "- [ ]" (unchecked) or "- [x]" (checked).
func (a *App) GetCheckboxLines() []int {
	sec := a.GetCurrentSection()
	if sec == nil {
		return nil
	}

	lines := strings.Split(sec.Content, "\n")
	checkboxLines := []int{}

	for i, line := range lines {
		if strings.Contains(line, "- [ ]") || strings.Contains(line, "- [x]") {
			checkboxLines = append(checkboxLines, i)
		}
	}

	return checkboxLines
}

// ToggleCheckbox toggles the checkbox at the given content line index.
// Returns true if a checkbox was toggled, false if the line has no checkbox.
func (a *App) ToggleCheckbox(contentLineIdx int) bool {
	sec := a.GetCurrentSection()
	if sec == nil {
		return false
	}

	lines := strings.Split(sec.Content, "\n")
	if contentLineIdx < 0 || contentLineIdx >= len(lines) {
		return false
	}

	line := lines[contentLineIdx]
	if strings.Contains(line, "- [ ]") {
		lines[contentLineIdx] = strings.Replace(line, "- [ ]", "- [x]", 1)
	} else if strings.Contains(line, "- [x]") {
		lines[contentLineIdx] = strings.Replace(line, "- [x]", "- [ ]", 1)
	} else {
		return false
	}

	a.Sections[a.CurrentIdx].Content = strings.Join(lines, "\n")
	return true
}

// AddNote appends a timestamped note to the current section.
// The note is formatted as a blockquote with the current timestamp.
func (a *App) AddNote(note string) {
	if note == "" {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04")
	noteText := fmt.Sprintf("\n\n> **Ghi chú [%s]:** %s", timestamp, note)
	a.Sections[a.CurrentIdx].Content += noteText
}

// GetProgress calculates the completion progress for a section.
// Returns (checked, total) where checked is the number of checked boxes
// and total is the total number of checkboxes.
func (a *App) GetProgress(sectionIdx int) (checked, total int) {
	if sectionIdx < 0 || sectionIdx >= len(a.Sections) {
		return 0, 0
	}

	content := a.Sections[sectionIdx].Content
	checked = strings.Count(content, "- [x]")
	total = checked + strings.Count(content, "- [ ]")
	return
}

// GetTotalProgress calculates the overall progress across all sections.
// Returns (checked, total) aggregated from all sections.
func (a *App) GetTotalProgress() (checked, total int) {
	for i := range a.Sections {
		c, t := a.GetProgress(i)
		checked += c
		total += t
	}
	return
}

// UpdateFileSection updates the file lines to reflect changes in a section.
// This syncs the in-memory section changes back to the file lines array.
func (a *App) UpdateFileSection(idx int) {
	if idx < 0 || idx >= len(a.Sections) {
		return
	}

	sec := a.Sections[idx]
	startLine := sec.Line

	// Find end line
	endLine := len(a.FileLines)
	if idx < len(a.Sections)-1 {
		endLine = a.Sections[idx+1].Line
	}

	// Rebuild section content
	headerLine := strings.Repeat("#", sec.Level) + " " + sec.Title
	newLines := []string{headerLine}
	newLines = append(newLines, strings.Split(sec.Content, "\n")...)

	// Replace in fileLines
	newFileLines := append(a.FileLines[:startLine], newLines...)
	newFileLines = append(newFileLines, a.FileLines[endLine:]...)
	a.FileLines = newFileLines

	// Update file content
	a.FileContent = strings.Join(a.FileLines, "\n")
}

// SaveFile writes the current file content to disk.
// Returns an error if the file cannot be written.
func (a *App) SaveFile() error {
	a.FileContent = strings.Join(a.FileLines, "\n")
	return os.WriteFile(a.FilePath, []byte(a.FileContent), 0o644)
}

// RenderLine converts a markdown line to ANSI-styled terminal output.
// It handles checkboxes, bold, italic, code, bullets, and blockquotes.
func RenderLine(line string, termWidth int) string {
	// Checkbox: - [ ] or - [x]
	if strings.Contains(line, "- [ ]") {
		line = strings.Replace(line, "- [ ]", Red+"☐"+Reset, 1)
	}
	if strings.Contains(line, "- [x]") {
		line = strings.Replace(line, "- [x]", Green+"☑"+Reset, 1)
	}

	// Bold: **text**
	boldRegex := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	line = boldRegex.ReplaceAllString(line, Bold+"$1"+Reset)

	// Italic: *text* (but not **)
	italicRegex := regexp.MustCompile(`(?:^|[^*])\*([^*]+)\*(?:[^*]|$)`)
	line = italicRegex.ReplaceAllString(line, Italic+"$1"+Reset)

	// Inline code: `code`
	codeRegex := regexp.MustCompile("`([^`]+)`")
	line = codeRegex.ReplaceAllString(line, BgBlack+Cyan+"$1"+Reset)

	// Bullet points (but not checkboxes)
	if strings.HasPrefix(strings.TrimSpace(line), "- ") &&
		!strings.Contains(line, "☐") &&
		!strings.Contains(line, "☑") {
		line = strings.Replace(line, "- ", Yellow+"• "+Reset, 1)
	}

	// Numbered lists
	numRegex := regexp.MustCompile(`^(\s*)(\d+)\.\s`)
	line = numRegex.ReplaceAllString(line, "$1"+Cyan+"$2."+Reset+" ")

	// Quote blocks: > text
	if strings.HasPrefix(strings.TrimSpace(line), ">") {
		line = Dim + "│ " + strings.TrimPrefix(strings.TrimSpace(line), "> ") + Reset
	}

	// Horizontal rule
	if strings.TrimSpace(line) == "---" {
		line = Dim + strings.Repeat("─", termWidth-4) + Reset
	}

	// Table separator
	if strings.Contains(line, "|") && strings.Contains(line, "---") {
		line = Dim + line + Reset
	}

	return line
}

// Renderer handles all terminal output operations.
type Renderer struct {
	App          *App
	TermWidth    int
	TermHeight   int
	ScrollOffset int // Track scroll within section content
	PageSize     int // Number of lines per page (user adjustable)
}

// NewRenderer creates a new Renderer for the given App.
func NewRenderer(app *App) *Renderer {
	// Default to showing more content - user can adjust with +/-
	pageSize := app.TermHeight - 6
	if pageSize < 15 {
		pageSize = 15
	}
	return &Renderer{
		App:          app,
		TermWidth:    app.TermWidth,
		TermHeight:   app.TermHeight,
		ScrollOffset: 0,
		PageSize:     pageSize,
	}
}

// ResetScroll resets the content scroll position.
func (r *Renderer) ResetScroll() {
	r.ScrollOffset = 0
}

// ScrollDown scrolls content down.
// Returns true if scrolled, false if already at bottom.
func (r *Renderer) ScrollDown() bool {
	sec := r.App.GetCurrentSection()
	if sec == nil {
		return false
	}

	lines := strings.Split(sec.Content, "\n")

	if r.ScrollOffset+r.PageSize < len(lines) {
		r.ScrollOffset += 3 // Scroll by 3 lines for smoother navigation
		return true
	}
	return false
}

// ScrollUp scrolls content up.
// Returns true if scrolled, false if already at top.
func (r *Renderer) ScrollUp() bool {
	if r.ScrollOffset > 0 {
		r.ScrollOffset -= 3 // Scroll by 3 lines
		if r.ScrollOffset < 0 {
			r.ScrollOffset = 0
		}
		return true
	}
	return false
}

// AdjustPageSize changes the number of visible lines.
// Minimum is 5 lines, no upper limit (content will scroll in terminal if needed).
func (r *Renderer) AdjustPageSize(delta int) {
	r.PageSize += delta
	if r.PageSize < 5 {
		r.PageSize = 5
	}
	// No upper limit - let user decide how much to show
}

// ClearScreen clears the terminal screen.
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// Render displays the current section with header and footer.
func (r *Renderer) Render() {
	ClearScreen()

	if len(r.App.Sections) == 0 {
		fmt.Println("Không có sections.")
		return
	}

	sec := r.App.GetCurrentSection()
	if sec == nil {
		return
	}

	r.printHeader(sec)
	r.printContent(sec.Content)
	r.printFooter()
}

// printHeader renders the top bar with progress and section title.
func (r *Renderer) printHeader(sec *Section) {
	// Progress bar
	progress := float64(r.App.CurrentIdx+1) / float64(len(r.App.Sections)) * 100
	barWidth := 20
	filled := int(float64(barWidth) * float64(r.App.CurrentIdx+1) / float64(len(r.App.Sections)))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	fmt.Printf("%s%s", BgBlue+White+Bold, strings.Repeat(" ", r.TermWidth))
	fmt.Print("\r")
	fmt.Printf(" 📖 SRE Learning Path  [%s] %.0f%%  (%d/%d)", bar, progress, r.App.CurrentIdx+1, len(r.App.Sections))
	fmt.Printf("%s\n", Reset)

	// Section title
	levelColors := []string{White, Cyan, Yellow, Green}
	levelColor := levelColors[min(sec.Level-1, 3)]
	prefix := strings.Repeat("  ", sec.Level-1)
	fmt.Printf("\n%s%s%s %s%s\n", prefix, Bold+levelColor, strings.Repeat("#", sec.Level), sec.Title, Reset)
	fmt.Println(Dim + strings.Repeat("─", r.TermWidth-4) + Reset)
}

// printContent renders the section content with markdown styling.
func (r *Renderer) printContent(content string) {
	lines := strings.Split(content, "\n")

	rendered := make([]string, len(lines))
	for i, line := range lines {
		rendered[i] = RenderLine(line, r.TermWidth)
	}

	// Apply scroll offset
	startIdx := r.ScrollOffset
	if startIdx >= len(rendered) {
		startIdx = 0
		r.ScrollOffset = 0
	}

	endIdx := min(startIdx+r.PageSize, len(rendered))
	displayLines := rendered[startIdx:endIdx]

	for _, line := range displayLines {
		fmt.Println(line)
	}

	// Show position indicator
	if len(rendered) > r.PageSize {
		above := startIdx
		below := len(rendered) - endIdx

		posInfo := fmt.Sprintf("[%d-%d/%d]", startIdx+1, endIdx, len(rendered))
		scrollHint := ""

		if above > 0 && below > 0 {
			scrollHint = fmt.Sprintf("↑%d ↓%d", above, below)
		} else if above > 0 {
			scrollHint = fmt.Sprintf("↑%d (k lên đầu)", above)
		} else if below > 0 {
			scrollHint = fmt.Sprintf("↓%d (j xem tiếp)", below)
		}

		fmt.Printf("\n%s%s %s  [%d dòng/trang, +/- chỉnh]%s", Dim, posInfo, scrollHint, r.PageSize, Reset)
	}
}

// printFooter renders the bottom navigation bar.
func (r *Renderer) printFooter() {
	fmt.Println()
	fmt.Printf("%s%s", BgBlack+White, strings.Repeat(" ", r.TermWidth))
	fmt.Print("\r")
	fmt.Printf(" %sj%s/%sk%s scroll %sn%s/%sp%s section %st%s toc %sx%s tick %sa%s note %s?%s help %sq%s quit",
		Bold+Cyan, Reset+BgBlack+White,
		Bold+Cyan, Reset+BgBlack+White,
		Bold+Cyan, Reset+BgBlack+White,
		Bold+Cyan, Reset+BgBlack+White,
		Bold+Cyan, Reset+BgBlack+White,
		Bold+Cyan, Reset+BgBlack+White,
		Bold+Cyan, Reset+BgBlack+White,
		Bold+Cyan, Reset+BgBlack+White,
		Bold+Cyan, Reset+BgBlack+White)
	fmt.Printf("%s\n", Reset)
}

// Terminal provides terminal manipulation utilities.
type Terminal struct{}

// GetSize returns the terminal dimensions (width, height).
// Falls back to 80x24 if unable to determine.
func (t *Terminal) GetSize() (width, height int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err == nil {
		fmt.Sscanf(string(out), "%d %d", &height, &width)
		return width, height
	}
	return 80, 24
}

// SetRawMode enables or disables raw terminal mode.
// In raw mode, input is read character by character without echo.
func (t *Terminal) SetRawMode(enable bool) {
	if enable {
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1", "-echo").Run()
	} else {
		exec.Command("stty", "-F", "/dev/tty", "-cbreak", "echo").Run()
	}
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Global instances for main program
var (
	app      *App
	renderer *Renderer
	terminal *Terminal
	reader   *bufio.Reader
)

func main() {
	app = NewApp()
	terminal = &Terminal{}

	// Get terminal size
	app.TermWidth, app.TermHeight = terminal.GetSize()

	// Check if file exists, prompt if not
	if !fileExists(app.FilePath) {
		handleFileNotFound()
	}

	// Load file
	if err := app.LoadFile(); err != nil {
		fmt.Printf("❌ Lỗi: %v\n", err)
		os.Exit(1)
	}
	app.ParseSections()

	// Create renderer with default settings
	renderer = NewRenderer(app)
	reader = bufio.NewReader(os.Stdin)

	// Load saved state (position, page size)
	if savedPageSize, err := app.LoadState(); err == nil {
		if savedPageSize > 0 {
			renderer.PageSize = savedPageSize
		}
		// Validate CurrentIdx
		if app.CurrentIdx >= len(app.Sections) {
			app.CurrentIdx = 0
		}
	}

	// Enable raw mode for keyboard input
	terminal.SetRawMode(true)
	defer func() {
		terminal.SetRawMode(false)
		// Save state on exit
		app.SaveState(renderer.PageSize)
	}()

	// Main loop
	for {
		renderer.Render()
		handleInput()
	}
}

// fileExists checks if a file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// handleFileNotFound prompts user when the markdown file doesn't exist.
func handleFileNotFound() {
	fmt.Printf("%s📚 SRE Learning Path CLI%s\n\n", Bold+Cyan, Reset)
	fmt.Printf("File %s%s%s không tồn tại.\n\n", Yellow, app.FilePath, Reset)
	fmt.Println("Chọn:")
	fmt.Printf("  %s1%s. Tạo file mới với template mặc định\n", Bold+Cyan, Reset)
	fmt.Printf("  %s2%s. Nhập đường dẫn file khác\n", Bold+Cyan, Reset)
	fmt.Printf("  %s3%s. Thoát\n", Bold+Cyan, Reset)
	fmt.Printf("\nLựa chọn (1/2/3): ")

	inputReader := bufio.NewReader(os.Stdin)
	input, _ := inputReader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		createDefaultFile()
	case "2":
		fmt.Printf("Nhập đường dẫn file: ")
		path, _ := inputReader.ReadString('\n')
		path = strings.TrimSpace(path)
		if path == "" {
			fmt.Println("Đường dẫn trống. Thoát.")
			os.Exit(1)
		}
		app.FilePath = path
		if !fileExists(app.FilePath) {
			fmt.Printf("File %s không tồn tại. Thoát.\n", app.FilePath)
			os.Exit(1)
		}
	default:
		fmt.Println("Thoát.")
		os.Exit(0)
	}
}

// createDefaultFile creates a new markdown file with default template.
func createDefaultFile() {
	if err := os.WriteFile(app.FilePath, []byte(defaultTemplate), 0o644); err != nil {
		fmt.Printf("❌ Không thể tạo file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s✅ Đã tạo file %s%s\n", Green, app.FilePath, Reset)
	time.Sleep(time.Second)
}

// handleInput reads and processes a single keyboard input.
func handleInput() {
	b := make([]byte, 3)
	os.Stdin.Read(b)

	switch {
	// Content scrolling within section
	case b[0] == 'j' || (b[0] == 27 && b[1] == 91 && b[2] == 66): // j or down arrow
		renderer.ScrollDown()
	case b[0] == 'k' || (b[0] == 27 && b[1] == 91 && b[2] == 65): // k or up arrow
		renderer.ScrollUp()

	// Section navigation
	case b[0] == 'n': // next section
		if app.NextSection() {
			renderer.ResetScroll()
		}
	case b[0] == 'p': // previous section
		if app.PrevSection() {
			renderer.ResetScroll()
		}
	case b[0] == 13 || b[0] == 10: // Enter - next section
		if app.NextSection() {
			renderer.ResetScroll()
		}

	// Features
	case b[0] == 't' || b[0] == 'T': // TOC
		handleTOC()
		renderer.ResetScroll()
	case b[0] == 'x' || b[0] == 'X': // toggle checkbox (x = check)
		handleToggle()
	case b[0] == 'g': // goto section
		handleGoto()
		renderer.ResetScroll()
	case b[0] == 'G': // goto last section
		app.GotoSection(len(app.Sections) - 1)
		renderer.ResetScroll()
	case b[0] == '/': // search
		handleSearch()
		renderer.ResetScroll()
	case b[0] == 'a' || b[0] == 'A': // add note
		handleNote()

	// Display settings
	case b[0] == '+' || b[0] == '=': // increase visible lines
		renderer.AdjustPageSize(10)
	case b[0] == '-' || b[0] == '_': // decrease visible lines
		renderer.AdjustPageSize(-10)

	// System
	case b[0] == 's' || b[0] == 'S': // save
		app.SaveFile()
		app.SaveState(renderer.PageSize)
	case b[0] == 'q' || b[0] == 'Q' || b[0] == 3: // quit or Ctrl+C
		terminal.SetRawMode(false)
		app.SaveState(renderer.PageSize)
		ClearScreen()
		fmt.Println("👋 Tạm biệt! Tiến độ đã lưu.")
		os.Exit(0)
	case b[0] == '?': // help
		handleHelp()
	}
}

// handleGoto displays section list and jumps to selected section.
func handleGoto() {
	terminal.SetRawMode(false)
	ClearScreen()

	fmt.Println(Bold + "📑 DANH SÁCH SECTIONS" + Reset)
	fmt.Println(Dim + strings.Repeat("─", 60) + Reset)

	for i, sec := range app.Sections {
		prefix := strings.Repeat("  ", sec.Level-1)
		marker := ""
		if i == app.CurrentIdx {
			marker = Green + " ◀" + Reset
		}

		checked, total := app.GetProgress(i)
		progress := ""
		if total > 0 {
			progress = fmt.Sprintf(" %s[%d/%d]%s", Dim, checked, total, Reset)
		}

		fmt.Printf("%s%3d. %s%s%s%s\n", Cyan, i+1, Reset, prefix, sec.Title, progress+marker)
	}

	fmt.Printf("\n%sNhập số (1-%d) hoặc Enter để hủy:%s ", Bold, len(app.Sections), Reset)

	inputReader := bufio.NewReader(os.Stdin)
	input, _ := inputReader.ReadString('\n')
	input = strings.TrimSpace(input)

	if num, err := strconv.Atoi(input); err == nil {
		app.GotoSection(num - 1)
	}

	terminal.SetRawMode(true)
}

// handleSearch prompts for search query and shows matching sections.
func handleSearch() {
	terminal.SetRawMode(false)
	ClearScreen()

	fmt.Printf("%s🔍 Tìm kiếm:%s ", Bold, Reset)

	inputReader := bufio.NewReader(os.Stdin)
	query, _ := inputReader.ReadString('\n')
	query = strings.TrimSpace(query)

	if query == "" {
		terminal.SetRawMode(true)
		return
	}

	matches := app.SearchSections(query)

	if len(matches) == 0 {
		fmt.Println(Red + "Không tìm thấy." + Reset)
		time.Sleep(time.Second)
		terminal.SetRawMode(true)
		return
	}

	fmt.Printf("\n%sTìm thấy %d kết quả:%s\n\n", Green, len(matches), Reset)
	for j, i := range matches {
		fmt.Printf("%s%2d.%s %s\n", Cyan, j+1, Reset, app.Sections[i].Title)
	}

	fmt.Printf("\n%sChọn số hoặc Enter để hủy:%s ", Bold, Reset)
	input, _ := inputReader.ReadString('\n')
	input = strings.TrimSpace(input)

	if num, err := strconv.Atoi(input); err == nil && num >= 1 && num <= len(matches) {
		app.GotoSection(matches[num-1])
	}

	terminal.SetRawMode(true)
}

// handleToggle displays checkboxes and toggles the selected one.
func handleToggle() {
	checkboxLines := app.GetCheckboxLines()
	if len(checkboxLines) == 0 {
		return
	}

	terminal.SetRawMode(false)
	ClearScreen()

	sec := app.GetCurrentSection()
	lines := strings.Split(sec.Content, "\n")

	fmt.Printf("%s☑ TOGGLE CHECKBOX%s\n", Bold, Reset)
	fmt.Println(Dim + strings.Repeat("─", 60) + Reset)

	for j, lineIdx := range checkboxLines {
		line := lines[lineIdx]
		status := Red + "☐" + Reset
		if strings.Contains(line, "- [x]") {
			status = Green + "☑" + Reset
		}
		text := strings.TrimSpace(line)
		text = strings.TrimPrefix(text, "- [ ]")
		text = strings.TrimPrefix(text, "- [x]")
		text = strings.TrimSpace(text)

		fmt.Printf("%s%2d.%s %s %s\n", Cyan, j+1, Reset, status, text)
	}

	fmt.Printf("\n%sNhập số để toggle (hoặc Enter để hủy):%s ", Bold, Reset)

	inputReader := bufio.NewReader(os.Stdin)
	input, _ := inputReader.ReadString('\n')
	input = strings.TrimSpace(input)

	if num, err := strconv.Atoi(input); err == nil && num >= 1 && num <= len(checkboxLines) {
		lineIdx := checkboxLines[num-1]
		if app.ToggleCheckbox(lineIdx) {
			app.UpdateFileSection(app.CurrentIdx)
			app.ParseSections() // Re-parse to update line numbers
			app.SaveFile()
		}
	}

	terminal.SetRawMode(true)
}

// handleNote provides a menu for note management.
func handleNote() {
	terminal.SetRawMode(false)
	// Reset terminal to sane state for proper input
	exec.Command("stty", "sane").Run()

	sec := app.GetCurrentSection()
	existingNotes := extractNotes(sec.Content)

	for {
		ClearScreen()
		fmt.Printf("%s📝 GHI CHÚ - %s%s\n", Bold+Cyan, sec.Title, Reset)
		fmt.Println(Dim + strings.Repeat("─", 60) + Reset)

		if len(existingNotes) > 0 {
			fmt.Printf("\n%sGhi chú hiện có (%d):%s\n\n", Yellow, len(existingNotes), Reset)
			for i, note := range existingNotes {
				// Truncate long notes for display
				displayNote := note
				if len(displayNote) > 200 {
					displayNote = displayNote[:200] + "..."
				}
				// Clean up for display
				displayNote = strings.ReplaceAll(displayNote, "\n", " ")
				fmt.Printf("  %s%d.%s %s\n", Cyan, i+1, Reset, displayNote)
			}
		} else {
			fmt.Printf("\n%sChưa có ghi chú nào.%s\n", Dim, Reset)
		}

		fmt.Println()
		fmt.Printf("%sChọn:%s\n", Bold, Reset)
		fmt.Printf("  %sa%s - Thêm ghi chú mới\n", Cyan, Reset)
		if len(existingNotes) > 0 {
			fmt.Printf("  %sv%s - Xem chi tiết ghi chú\n", Cyan, Reset)
			fmt.Printf("  %se%s - Sửa ghi chú\n", Cyan, Reset)
			fmt.Printf("  %sd%s - Xóa ghi chú\n", Cyan, Reset)
			fmt.Printf("  %sc%s - Xóa TẤT CẢ ghi chú (clean)\n", Cyan, Reset)
		}
		fmt.Printf("  %sq%s - Quay lại\n", Cyan, Reset)
		fmt.Printf("\nLựa chọn: ")

		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(strings.ToLower(choice))

		switch choice {
		case "a":
			addNewNote(reader)
			// Refresh notes list
			sec = app.GetCurrentSection()
			existingNotes = extractNotes(sec.Content)
		case "v":
			if len(existingNotes) > 0 {
				viewNoteDetail(existingNotes, reader)
			}
		case "e":
			if len(existingNotes) > 0 {
				if editNote(reader, existingNotes) {
					// Refresh after edit
					sec = app.GetCurrentSection()
					existingNotes = extractNotes(sec.Content)
				}
			}
		case "d":
			if len(existingNotes) > 0 {
				if deleteNote(reader, existingNotes) {
					// Refresh after delete
					sec = app.GetCurrentSection()
					existingNotes = extractNotes(sec.Content)
				}
			}
		case "c":
			if len(existingNotes) > 0 {
				if cleanAllNotes(reader) {
					// Refresh after clean
					sec = app.GetCurrentSection()
					existingNotes = extractNotes(sec.Content)
				}
			}
		case "q", "":
			terminal.SetRawMode(true)
			return
		}
	}
}

// addNewNote handles adding a new note using an external editor.
// This ensures proper UTF-8 support and cursor navigation.
func addNewNote(reader *bufio.Reader) {
	ClearScreen()
	fmt.Printf("%s📝 THÊM GHI CHÚ MỚI%s\n", Bold+Cyan, Reset)
	fmt.Println(Dim + strings.Repeat("─", 60) + Reset)
	fmt.Println()

	// Create temp file for editing
	tmpFile, err := os.CreateTemp("", "sre-note-*.txt")
	if err != nil {
		fmt.Printf("%s❌ Lỗi tạo file tạm: %v%s\n", Red, err, Reset)
		fmt.Printf("\n%s[Enter để quay lại]%s", Dim, Reset)
		reader.ReadString('\n')
		return
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)
	tmpFile.Close()

	// Find editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		// Try common editors
		for _, e := range []string{"nano", "vim", "vi", "notepad"} {
			if _, err := exec.LookPath(e); err == nil {
				editor = e
				break
			}
		}
	}

	if editor == "" {
		// Fallback to simple stdin input
		fmt.Println("Không tìm thấy editor (nano/vim). Dùng input đơn giản:")
		fmt.Println("(Nhập ghi chú, dòng trống để kết thúc)")
		fmt.Println()

		var lines []string
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			line = strings.TrimRight(line, "\r\n")
			if line == "" {
				break
			}
			lines = append(lines, line)
		}

		note := strings.TrimSpace(strings.Join(lines, "\n"))
		if note != "" {
			saveNote(note)
		}
		return
	}

	fmt.Printf("Mở %s%s%s để soạn ghi chú...\n", Bold+Cyan, editor, Reset)
	fmt.Printf("%s(Lưu và thoát editor để hoàn thành)%s\n", Dim, Reset)
	time.Sleep(500 * time.Millisecond)

	// Open editor
	cmd := exec.Command(editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("\n%s❌ Lỗi mở editor: %v%s\n", Red, err, Reset)
		fmt.Printf("\n%s[Enter để quay lại]%s", Dim, Reset)
		reader.ReadString('\n')
		return
	}

	// Read the edited content
	content, err := os.ReadFile(tmpPath)
	if err != nil {
		fmt.Printf("\n%s❌ Lỗi đọc file: %v%s\n", Red, err, Reset)
		fmt.Printf("\n%s[Enter để quay lại]%s", Dim, Reset)
		reader.ReadString('\n')
		return
	}

	note := strings.TrimSpace(string(content))
	if note == "" {
		fmt.Printf("\n%sGhi chú trống - đã hủy.%s\n", Yellow, Reset)
		time.Sleep(time.Second)
		return
	}

	saveNote(note)
}

// saveNote saves a note to the current section.
func saveNote(note string) {
	app.AddNote(note)
	app.UpdateFileSection(app.CurrentIdx)
	app.ParseSections()
	if err := app.SaveFile(); err != nil {
		fmt.Printf("\n%s❌ Lỗi lưu: %v%s\n", Red, err, Reset)
	} else {
		fmt.Printf("\n%s✅ Đã lưu ghi chú!%s\n", Green, Reset)
	}
	time.Sleep(time.Second)
}

// viewNoteDetail shows full content of a specific note.
func viewNoteDetail(notes []string, reader *bufio.Reader) {
	ClearScreen()
	fmt.Printf("%s📖 XEM GHI CHÚ%s\n", Bold+Cyan, Reset)
	fmt.Println(Dim + strings.Repeat("─", 60) + Reset)
	fmt.Println()

	for i := range notes {
		fmt.Printf("  %s%d%s. Ghi chú #%d\n", Cyan, i+1, Reset, i+1)
	}

	fmt.Printf("\nNhập số (1-%d) hoặc Enter để quay lại: ", len(notes))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return
	}

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(notes) {
		return
	}

	// Show full note
	ClearScreen()
	fmt.Printf("%s📖 GHI CHÚ #%d%s\n", Bold+Cyan, idx, Reset)
	fmt.Println(Dim + strings.Repeat("─", 60) + Reset)
	fmt.Println()
	fmt.Println(notes[idx-1])
	fmt.Println()
	fmt.Printf("%s[Enter để quay lại]%s", Dim, Reset)
	reader.ReadString('\n')
}

// editNote opens an editor to modify an existing note.
func editNote(reader *bufio.Reader, notes []string) bool {
	ClearScreen()
	fmt.Printf("%s✏️ SỬA GHI CHÚ%s\n", Bold+Cyan, Reset)
	fmt.Println(Dim + strings.Repeat("─", 60) + Reset)
	fmt.Println()

	for i, note := range notes {
		displayNote := note
		if len(displayNote) > 100 {
			displayNote = displayNote[:100] + "..."
		}
		displayNote = strings.ReplaceAll(displayNote, "\n", " ")
		fmt.Printf("  %s%d%s. %s\n", Cyan, i+1, Reset, displayNote)
	}

	fmt.Printf("\nNhập số để sửa (1-%d) hoặc Enter để hủy: ", len(notes))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return false
	}

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(notes) {
		return false
	}

	oldNote := notes[idx-1]

	// Extract just the note content (remove timestamp prefix)
	noteContent := oldNote
	if strings.HasPrefix(noteContent, "> **Ghi chú [") {
		// Find the end of timestamp
		if endIdx := strings.Index(noteContent, ":**"); endIdx != -1 {
			noteContent = strings.TrimSpace(noteContent[endIdx+3:])
		}
	}
	// Remove leading > from subsequent lines
	lines := strings.Split(noteContent, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimPrefix(strings.TrimPrefix(line, "> "), ">")
	}
	noteContent = strings.Join(lines, "\n")

	// Create temp file with existing content
	tmpFile, err := os.CreateTemp("", "sre-note-edit-*.txt")
	if err != nil {
		fmt.Printf("%s❌ Lỗi tạo file tạm: %v%s\n", Red, err, Reset)
		fmt.Printf("\n%s[Enter để quay lại]%s", Dim, Reset)
		reader.ReadString('\n')
		return false
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	tmpFile.WriteString(noteContent)
	tmpFile.Close()

	// Find editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		for _, e := range []string{"nano", "vim", "vi"} {
			if _, err := exec.LookPath(e); err == nil {
				editor = e
				break
			}
		}
	}

	if editor == "" {
		fmt.Printf("%s❌ Không tìm thấy editor%s\n", Red, Reset)
		fmt.Printf("\n%s[Enter để quay lại]%s", Dim, Reset)
		reader.ReadString('\n')
		return false
	}

	fmt.Printf("\nMở %s%s%s để sửa...\n", Bold+Cyan, editor, Reset)
	time.Sleep(500 * time.Millisecond)

	// Open editor
	cmd := exec.Command(editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("\n%s❌ Lỗi mở editor: %v%s\n", Red, err, Reset)
		fmt.Printf("\n%s[Enter để quay lại]%s", Dim, Reset)
		reader.ReadString('\n')
		return false
	}

	// Read edited content
	content, err := os.ReadFile(tmpPath)
	if err != nil {
		fmt.Printf("\n%s❌ Lỗi đọc file: %v%s\n", Red, err, Reset)
		fmt.Printf("\n%s[Enter để quay lại]%s", Dim, Reset)
		reader.ReadString('\n')
		return false
	}

	newNote := strings.TrimSpace(string(content))
	if newNote == "" {
		fmt.Printf("\n%sGhi chú trống - đã hủy.%s\n", Yellow, Reset)
		time.Sleep(time.Second)
		return false
	}

	// Replace old note with new one
	sec := app.GetCurrentSection()
	newContent := removeNoteFromContent(sec.Content, oldNote)
	app.Sections[app.CurrentIdx].Content = newContent

	// Add the edited note
	app.AddNote(newNote)
	app.UpdateFileSection(app.CurrentIdx)
	app.ParseSections()

	if err := app.SaveFile(); err != nil {
		fmt.Printf("\n%s❌ Lỗi lưu: %v%s\n", Red, err, Reset)
		time.Sleep(time.Second)
		return false
	}

	fmt.Printf("\n%s✅ Đã cập nhật ghi chú!%s\n", Green, Reset)
	time.Sleep(time.Second)
	return true
}

// deleteNote removes a note from the section.
func deleteNote(reader *bufio.Reader, notes []string) bool {
	ClearScreen()
	fmt.Printf("%s🗑️ XÓA GHI CHÚ%s\n", Bold+Red, Reset)
	fmt.Println(Dim + strings.Repeat("─", 60) + Reset)
	fmt.Println()

	for i, note := range notes {
		displayNote := note
		if len(displayNote) > 100 {
			displayNote = displayNote[:100] + "..."
		}
		displayNote = strings.ReplaceAll(displayNote, "\n", " ")
		fmt.Printf("  %s%d%s. %s\n", Cyan, i+1, Reset, displayNote)
	}

	fmt.Printf("\nNhập số để xóa (1-%d) hoặc Enter để hủy: ", len(notes))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return false
	}

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(notes) {
		return false
	}

	// Confirm delete
	fmt.Printf("\n%sXác nhận xóa ghi chú #%d? (y/N): %s", Yellow, idx, Reset)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" && confirm != "yes" {
		return false
	}

	// Remove note from content
	noteToDelete := notes[idx-1]
	sec := app.GetCurrentSection()
	newContent := removeNoteFromContent(sec.Content, noteToDelete)
	app.Sections[app.CurrentIdx].Content = newContent

	app.UpdateFileSection(app.CurrentIdx)
	app.ParseSections()
	if err := app.SaveFile(); err != nil {
		fmt.Printf("\n%s❌ Lỗi: %v%s\n", Red, err, Reset)
		time.Sleep(time.Second)
		return false
	}

	fmt.Printf("\n%s✅ Đã xóa ghi chú!%s\n", Green, Reset)
	time.Sleep(time.Second)
	return true
}

// removeNoteFromContent removes a specific note from section content.
func removeNoteFromContent(content, noteToRemove string) string {
	// Find and remove the note block
	lines := strings.Split(content, "\n")
	var result []string
	skipUntilNonNote := false
	noteLines := strings.Split(noteToRemove, "\n")
	firstNoteLine := ""
	if len(noteLines) > 0 {
		firstNoteLine = strings.TrimSpace(noteLines[0])
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Check if this is the start of the note to delete
		if strings.Contains(trimmed, "**Ghi chú [") && strings.Contains(firstNoteLine, trimmed[2:]) {
			skipUntilNonNote = true
			continue
		}

		if skipUntilNonNote {
			// Skip lines that are part of the note (start with > or are empty after note)
			if strings.HasPrefix(trimmed, ">") {
				continue
			}
			// Also skip empty lines immediately after note
			if trimmed == "" && i+1 < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i+1]), ">") {
				continue
			}
			skipUntilNonNote = false
		}

		result = append(result, line)
	}

	// Clean up multiple consecutive empty lines
	return strings.TrimSpace(strings.Join(result, "\n"))
}

// cleanAllNotes removes all notes from current section.
func cleanAllNotes(reader *bufio.Reader) bool {
	fmt.Printf("\n%s⚠️ Xác nhận xóa TẤT CẢ ghi chú trong section này? (y/N): %s", Yellow, Reset)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" && confirm != "yes" {
		return false
	}

	// Remove all notes from content
	sec := app.GetCurrentSection()
	lines := strings.Split(sec.Content, "\n")
	var result []string
	inNote := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if this is start of a note
		if strings.HasPrefix(trimmed, "> **Ghi chú [") {
			inNote = true
			continue
		}

		if inNote {
			if strings.HasPrefix(trimmed, ">") {
				continue // Skip note content
			}
			if trimmed == "" {
				continue // Skip empty lines after note
			}
			inNote = false
		}

		result = append(result, line)
	}

	app.Sections[app.CurrentIdx].Content = strings.TrimSpace(strings.Join(result, "\n"))
	app.UpdateFileSection(app.CurrentIdx)
	app.ParseSections()

	if err := app.SaveFile(); err != nil {
		fmt.Printf("\n%s❌ Lỗi: %v%s\n", Red, err, Reset)
		time.Sleep(time.Second)
		return false
	}

	fmt.Printf("\n%s✅ Đã xóa tất cả ghi chú!%s\n", Green, Reset)
	time.Sleep(time.Second)
	return true
}

// extractNotes extracts existing notes from section content.
func extractNotes(content string) []string {
	var notes []string
	lines := strings.Split(content, "\n")
	var currentNote strings.Builder
	inNote := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "> **Ghi chú [") {
			// Save previous note if exists
			if currentNote.Len() > 0 {
				notes = append(notes, strings.TrimSpace(currentNote.String()))
			}
			currentNote.Reset()
			inNote = true
			currentNote.WriteString(trimmed)
		} else if inNote && strings.HasPrefix(trimmed, ">") {
			currentNote.WriteString("\n")
			currentNote.WriteString(trimmed)
		} else if inNote && trimmed == "" {
			// Empty line might be part of note or end of note
			// Look ahead logic would be complex, so just end the note
			if currentNote.Len() > 0 {
				notes = append(notes, strings.TrimSpace(currentNote.String()))
				currentNote.Reset()
			}
			inNote = false
		} else {
			// Non-note line
			if inNote && currentNote.Len() > 0 {
				notes = append(notes, strings.TrimSpace(currentNote.String()))
				currentNote.Reset()
			}
			inNote = false
		}
	}

	// Don't forget last note
	if currentNote.Len() > 0 {
		notes = append(notes, strings.TrimSpace(currentNote.String()))
	}

	return notes
}

// handleHelp displays all keyboard shortcuts.
func handleHelp() {
	ClearScreen()

	fmt.Printf("%s%s", BgCyan+Black+Bold, strings.Repeat(" ", app.TermWidth))
	fmt.Print("\r")
	fmt.Printf(" ❓ KEYBOARD SHORTCUTS")
	fmt.Printf("%s\n\n", Reset)

	helpItems := []struct {
		key  string
		desc string
	}{
		{"j / ↓", "Scroll xuống trong section"},
		{"k / ↑", "Scroll lên trong section"},
		{"n", "Section tiếp theo (next)"},
		{"p", "Section trước (previous)"},
		{"Enter", "Section tiếp theo"},
		{"", ""},
		{"t", "Mở Table of Contents"},
		{"g", "Goto - nhảy đến section"},
		{"G", "Goto section cuối"},
		{"/", "Tìm kiếm section"},
		{"", ""},
		{"x", "Toggle checkbox (tick/untick)"},
		{"a", "Ghi chú (thêm/xem/sửa/xóa)"},
		{"s", "Lưu file & tiến độ"},
		{"", ""},
		{"+", "Tăng 10 dòng hiển thị"},
		{"-", "Giảm 10 dòng hiển thị"},
		{"", ""},
		{"?", "Hiển thị help này"},
		{"q", "Thoát"},
	}

	for _, item := range helpItems {
		if item.key == "" {
			fmt.Println()
		} else {
			fmt.Printf("  %s%-10s%s %s\n", Bold+Cyan, item.key, Reset, item.desc)
		}
	}

	fmt.Printf("\n%sTrong TOC:%s\n", Bold+Magenta, Reset)
	fmt.Printf("  %s%-10s%s %s\n", Bold+Cyan, "j/k", Reset, "Di chuyển lên/xuống")
	fmt.Printf("  %s%-10s%s %s\n", Bold+Cyan, "Enter", Reset, "Chọn section")
	fmt.Printf("  %s%-10s%s %s\n", Bold+Cyan, "q/Esc", Reset, "Đóng TOC")

	fmt.Printf("\n%sGhi chú (nhấn a):%s\n", Bold+Magenta, Reset)
	fmt.Printf("  %s%-10s%s %s\n", Bold+Cyan, "a", Reset, "Thêm mới (mở editor)")
	fmt.Printf("  %s%-10s%s %s\n", Bold+Cyan, "v", Reset, "Xem chi tiết")
	fmt.Printf("  %s%-10s%s %s\n", Bold+Cyan, "e", Reset, "Sửa ghi chú")
	fmt.Printf("  %s%-10s%s %s\n", Bold+Cyan, "d", Reset, "Xóa")
	fmt.Printf("  %sDùng nano/vim, set EDITOR env để đổi editor%s\n", Dim, Reset)

	fmt.Printf("\n%sHiện tại: %d dòng/trang (nhấn +/- để chỉnh, không giới hạn)%s\n", Dim, renderer.PageSize, Reset)

	fmt.Printf("\n%s[Nhấn phím bất kỳ để quay lại]%s", Dim, Reset)

	// Wait for any key
	b := make([]byte, 1)
	os.Stdin.Read(b)
}

// handleTOC displays an interactive table of contents.
// Supports j/k navigation, Enter to select, q to quit.
func handleTOC() {
	// Build list of navigable sections (skip phase headers)
	type tocItem struct {
		idx   int
		title string
		level int
	}

	items := []tocItem{}
	for i, sec := range app.Sections {
		items = append(items, tocItem{i, sec.Title, sec.Level})
	}

	if len(items) == 0 {
		return
	}

	// Find current position in TOC
	tocIdx := 0
	for i, item := range items {
		if item.idx == app.CurrentIdx {
			tocIdx = i
			break
		}
	}

	// Scrolling state
	scrollOffset := 0
	maxVisible := app.TermHeight - 6

	for {
		ClearScreen()

		// Header
		fmt.Printf("%s%s", BgMagenta+White+Bold, strings.Repeat(" ", app.TermWidth))
		fmt.Print("\r")
		fmt.Printf(" 📚 MỤC LỤC  (j/k: di chuyển, Enter: chọn, q: đóng)")
		fmt.Printf("%s\n\n", Reset)

		// Adjust scroll to keep selection visible
		if tocIdx < scrollOffset {
			scrollOffset = tocIdx
		}
		if tocIdx >= scrollOffset+maxVisible {
			scrollOffset = tocIdx - maxVisible + 1
		}

		// Display items
		endIdx := min(scrollOffset+maxVisible, len(items))
		for i := scrollOffset; i < endIdx; i++ {
			item := items[i]

			// Selection indicator
			selector := "  "
			if i == tocIdx {
				selector = Green + "▶ " + Reset
			}

			// Indentation based on level
			indent := strings.Repeat("  ", item.level-1)

			// Progress indicator
			checked, total := app.GetProgress(item.idx)
			progress := ""
			if total > 0 {
				pct := float64(checked) / float64(total) * 100
				if pct == 100 {
					progress = Green + " ✓" + Reset
				} else if pct > 0 {
					progress = fmt.Sprintf(" %s%.0f%%%s", Yellow, pct, Reset)
				} else {
					progress = Dim + " ○" + Reset
				}
			}

			// Current section marker
			current := ""
			if item.idx == app.CurrentIdx {
				current = Cyan + " (hiện tại)" + Reset
			}

			// Title styling based on level
			title := item.title
			if len(title) > 50 {
				title = title[:47] + "..."
			}

			titleStyle := ""
			switch item.level {
			case 1:
				titleStyle = Bold + White
			case 2:
				titleStyle = Bold + Magenta
			case 3:
				titleStyle = Cyan
			default:
				titleStyle = Dim
			}

			// Print row
			fmt.Printf("%s%s%s%s%s%s%s\n", selector, indent, titleStyle, title, Reset, progress, current)
		}

		// Scroll indicators
		if scrollOffset > 0 {
			fmt.Printf("\n%s  ↑ còn %d mục phía trên%s", Dim, scrollOffset, Reset)
		}
		if endIdx < len(items) {
			if scrollOffset == 0 {
				fmt.Println()
			}
			fmt.Printf("\n%s  ↓ còn %d mục phía dưới%s", Dim, len(items)-endIdx, Reset)
		}

		// Footer with total progress
		fmt.Println()
		checked, total := app.GetTotalProgress()
		if total > 0 {
			pct := float64(checked) / float64(total) * 100
			barWidth := 20
			filled := int(float64(barWidth) * pct / 100)
			bar := Green + strings.Repeat("█", filled) + Dim + strings.Repeat("░", barWidth-filled) + Reset
			fmt.Printf("\n  Tiến độ: [%s] %d/%d (%.0f%%)\n", bar, checked, total, pct)
		}

		// Read input
		b := make([]byte, 3)
		os.Stdin.Read(b)

		switch {
		case b[0] == 'j' || (b[0] == 27 && b[1] == 91 && b[2] == 66): // j or down
			if tocIdx < len(items)-1 {
				tocIdx++
			}
		case b[0] == 'k' || (b[0] == 27 && b[1] == 91 && b[2] == 65): // k or up
			if tocIdx > 0 {
				tocIdx--
			}
		case b[0] == 'g': // go to top
			tocIdx = 0
			scrollOffset = 0
		case b[0] == 'G': // go to bottom
			tocIdx = len(items) - 1
		case b[0] == 13 || b[0] == 10: // Enter - select
			app.GotoSection(items[tocIdx].idx)
			return
		case b[0] == 'q' || b[0] == 'Q' || b[0] == 27: // q or Escape - close
			return
		case b[0] == ' ': // Space - page down
			tocIdx = min(tocIdx+maxVisible, len(items)-1)
		}
	}
}
