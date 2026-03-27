package service

import "strings"

func buildAIReportSummaryV2(report string) string {
	report = normalizeAIReportMarkdown(report)
	lines := strings.Split(report, "\n")
	for _, line := range lines {
		text := strings.TrimSpace(line)
		if isMarkdownNoiseLine(text) {
			continue
		}
		text = strings.TrimLeft(text, "#-*0123456789. ")
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		runes := []rune(text)
		if len(runes) > 42 {
			return string(runes[:42]) + "..."
		}
		return text
	}
	return "社区运营简报（近7日）"
}

func normalizeAIReportMarkdown(content string) string {
	text := strings.ReplaceAll(content, "\r\n", "\n")
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	lines := strings.Split(text, "\n")
	if len(lines) > 0 && strings.HasPrefix(strings.TrimSpace(lines[0]), "```") {
		lines = lines[1:]
		for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
			lines = lines[:len(lines)-1]
		}
		if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "```" {
			lines = lines[:len(lines)-1]
		}
		text = strings.TrimSpace(strings.Join(lines, "\n"))
	}

	lines = strings.Split(text, "\n")
	if len(lines) > 0 {
		first := strings.ToLower(strings.TrimSpace(lines[0]))
		if first == "markdown" || first == "md" {
			text = strings.TrimSpace(strings.Join(lines[1:], "\n"))
		}
	}
	return text
}

func normalizeAIReportSummaryText(summary string) string {
	text := strings.TrimSpace(summary)
	text = strings.TrimLeft(text, "`#-*0123456789. ")
	text = strings.TrimSpace(text)
	if isMarkdownNoiseLine(text) || text == "" {
		return "社区运营简报（近7日）"
	}

	runes := []rune(text)
	if len(runes) > 42 {
		return string(runes[:42]) + "..."
	}
	return text
}

func isMarkdownNoiseLine(text string) bool {
	line := strings.TrimSpace(text)
	if line == "" {
		return true
	}

	lower := strings.ToLower(line)
	if lower == "markdown" || lower == "md" {
		return true
	}
	if strings.HasPrefix(line, "```") || strings.HasPrefix(line, ">") {
		return true
	}
	if strings.Contains(line, "|") {
		return true
	}
	return false
}
