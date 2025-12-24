package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const dataFile = "learning-path-full.md"

var reader *bufio.Reader

func main() {
	reader = bufio.NewReader(os.Stdin)

	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		initMarkdown()
	}

	clearScreen()
	printBanner()

	for {
		showMenu()
		choice := readInput("\n→ ")

		switch choice {
		case "1":
			viewFile()
		case "2":
			toggleCheckbox()
		case "3":
			addNote()
		case "4":
			addDiscussion()
		case "5":
			openEditor()
		case "q", "Q":
			fmt.Println("\n👋 Tạm biệt!")
			os.Exit(0)
		}
	}
}

func printBanner() {
	fmt.Println(`
╔══════════════════════════════════════════════╗
║     SRE LEARNING PATH - CLI                  ║
║     File: learning-path.md                   ║
╚══════════════════════════════════════════════╝`)
}

func showMenu() {
	fmt.Println("\n1. Xem nội dung")
	fmt.Println("2. Tick/Untick checkbox")
	fmt.Println("3. Thêm ghi chú")
	fmt.Println("4. Thêm thảo luận")
	fmt.Println("5. Mở editor (vi)")
	fmt.Println("q. Thoát")
}

func viewFile() {
	clearScreen()
	content, _ := os.ReadFile(dataFile)
	fmt.Println(string(content))
	readInput("\n[Enter để quay lại]")
}

func toggleCheckbox() {
	content, _ := os.ReadFile(dataFile)
	lines := strings.Split(string(content), "\n")

	// Find all checkbox lines
	checkboxes := []struct {
		idx  int
		line string
	}{}

	for i, line := range lines {
		if strings.Contains(line, "- [ ]") || strings.Contains(line, "- [x]") {
			checkboxes = append(checkboxes, struct {
				idx  int
				line string
			}{i, line})
		}
	}

	clearScreen()
	fmt.Println("═══ CHECKBOX ═══\n")
	for i, cb := range checkboxes {
		fmt.Printf("%2d. %s\n", i+1, strings.TrimSpace(cb.line))
	}

	fmt.Println("\nNhập số để toggle, 0 để quay lại")
	choice := readInput("→ ")
	idx, err := strconv.Atoi(choice)
	if err != nil || idx < 1 || idx > len(checkboxes) {
		return
	}

	lineIdx := checkboxes[idx-1].idx
	if strings.Contains(lines[lineIdx], "- [ ]") {
		lines[lineIdx] = strings.Replace(lines[lineIdx], "- [ ]", "- [x]", 1)
	} else {
		lines[lineIdx] = strings.Replace(lines[lineIdx], "- [x]", "- [ ]", 1)
	}

	os.WriteFile(dataFile, []byte(strings.Join(lines, "\n")), 0o644)
	fmt.Println("✅ Đã cập nhật!")
}

func addNote() {
	fmt.Println("\n📝 Nhập ghi chú (END để kết thúc):")
	note := readMultiline()
	if note == "" {
		return
	}

	content, _ := os.ReadFile(dataFile)
	timestamp := time.Now().Format("2006-01-02 15:04")

	newContent := string(content) + fmt.Sprintf("\n### Ghi chú - %s\n\n%s\n", timestamp, note)
	os.WriteFile(dataFile, []byte(newContent), 0o644)
	fmt.Println("✅ Đã thêm ghi chú!")
}

func addDiscussion() {
	topic := readInput("\n📌 Chủ đề: ")
	if topic == "" {
		return
	}

	fmt.Println("💬 Nội dung (END để kết thúc):")
	content := readMultiline()
	if content == "" {
		return
	}

	fileContent, _ := os.ReadFile(dataFile)
	timestamp := time.Now().Format("2006-01-02 15:04")

	newContent := string(fileContent) + fmt.Sprintf("\n### Thảo luận: %s - %s\n\n%s\n", topic, timestamp, content)
	os.WriteFile(dataFile, []byte(newContent), 0o644)
	fmt.Println("✅ Đã thêm thảo luận!")
}

func openEditor() {
	fmt.Println("\n📄 Mở file bằng editor yêu thích:")
	fmt.Printf("   vi %s\n", dataFile)
	fmt.Printf("   nano %s\n", dataFile)
	fmt.Printf("   code %s\n", dataFile)
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func readMultiline() string {
	var lines []string
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimRight(line, "\n\r")
		if strings.ToUpper(line) == "END" {
			break
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func initMarkdown() {
	fmt.Println("⚠️  File learning-path-full.md không tìm thấy!")
	fmt.Println("   Đặt file vào cùng thư mục với CLI tool.")
	os.Exit(1)
}
