## Godoc coverage:

- Package documentation
- Tất cả exported types: App, Section, Renderer, Terminal
- Tất cả exported functions với mô tả params và returns
- ANSI color constants documented

## Unit tests:

```bash
# Chạy tests
go test -v ./...

# Chạy với coverage

go test -cover ./...

# Chạy benchmarks

go test -bench=. ./...
```

## Test categories:

- TestNewApp, TestParseSections - khởi tạo
- TestNextSection, TestPrevSection, TestGotoSection - navigation
- TestSearchSections - tìm kiếm
- TestToggleCheckbox, TestGetCheckboxLines - checkbox
- TestAddNote - ghi chú
- TestGetProgress, TestGetTotalProgress - tiến độ
- TestRenderLine\* - markdown rendering
- TestNavigationFlow, TestCheckboxWorkflow - integration
- TestEmptyFile, TestSpecialCharacters - edge cases
- BenchmarkParseSections, BenchmarkRenderLine - performance

## Build & Run:

```bash
go build -o sre-learn .
./sre-learn
```
