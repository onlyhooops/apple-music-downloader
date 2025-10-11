package ui

import (
	"fmt"
	"main/internal/core"
	"regexp"
	"strings"
)

// DisplayMode 定义显示模式
type DisplayMode int

const (
	FullMode    DisplayMode = iota // 完整模式（120+字符终端）
	CompactMode                    // 紧凑模式（80-119字符）
	MinimalMode                    // 极简模式（<80字符）
)

// GetDisplayMode 根据终端宽度自动选择显示模式
func GetDisplayMode(termWidth int) DisplayMode {
	if termWidth >= 120 {
		return FullMode
	} else if termWidth >= 80 {
		return CompactMode
	}
	return MinimalMode
}

// FormatTrackLine 格式化曲目状态行（智能自适应）
// 保证输出永远不超过终端宽度，避免自动换行导致的光标错位
func FormatTrackLine(ts core.TrackStatus, termWidth int) string {
	mode := GetDisplayMode(termWidth)

	var line string
	switch mode {
	case FullMode:
		line = formatFull(ts)
	case CompactMode:
		line = formatCompact(ts)
	case MinimalMode:
		line = formatMinimal(ts)
	}

	// 最终安全检查：强制截断到终端宽度-2（留margin）
	visualLen := getVisualLength(line)
	if visualLen > termWidth-2 {
		line = truncateToWidth(line, termWidth-2)
	}

	return line
}

// formatFull 完整模式（120+字符）
// 格式: [3/7] Spanish Key (feat. Wayne Shorter, ...) (16bit/96.0kHz) ↓54% 4.4MB/s
func formatFull(ts core.TrackStatus) string {
	simpleName := simplifyTrackName(ts.TrackName, 50)
	statusStr := formatStatusDetailed(ts.Status, true)

	return fmt.Sprintf("[%d/%d] %s (%s) %s",
		ts.TrackNum, ts.TrackTotal,
		simpleName,
		ts.Quality,
		statusStr)
}

// formatCompact 紧凑模式（80-119字符）
// 格式: [3/7] Spanish Key (feat. ...) 16/96 ↓54% 4.4MB/s
func formatCompact(ts core.TrackStatus) string {
	simpleName := simplifyTrackName(ts.TrackName, 35)
	shortQuality := shortQualityFormat(ts.Quality)
	statusStr := formatStatusDetailed(ts.Status, true)

	return fmt.Sprintf("[%d/%d] %s %s %s",
		ts.TrackNum, ts.TrackTotal,
		simpleName,
		shortQuality,
		statusStr)
}

// formatMinimal 极简模式（<80字符）
// 格式: [3/7] Spanish Key... ↓54%
func formatMinimal(ts core.TrackStatus) string {
	simpleName := simplifyTrackName(ts.TrackName, 20)
	statusStr := formatStatusDetailed(ts.Status, false)

	return fmt.Sprintf("[%d/%d] %s %s",
		ts.TrackNum, ts.TrackTotal,
		simpleName,
		statusStr)
}

// simplifyTrackName 简化曲目名称
// 处理过长的feat艺术家列表，并截断到指定长度
func simplifyTrackName(name string, maxLen int) string {
	// 处理feat艺术家列表（通常很长）
	if idx := strings.Index(name, "(feat."); idx > 0 {
		mainTitle := strings.TrimSpace(name[:idx])

		// 如果主标题本身就太长，直接截断
		if len([]rune(mainTitle)) > maxLen-10 {
			return truncateString(mainTitle, maxLen)
		}

		// 提取feat部分
		featPart := name[idx:]

		// 只保留第一个艺术家
		if commaIdx := strings.Index(featPart, ","); commaIdx > 0 {
			// (feat. Artist1, Artist2, ...) → (feat. Artist1, ...)
			featPart = featPart[:commaIdx] + ", ...)"
		} else if andIdx := strings.Index(featPart, " & "); andIdx > 0 {
			// (feat. Artist1 & Artist2) → (feat. Artist1, ...)
			featPart = featPart[:andIdx] + ", ...)"
		}

		name = mainTitle + " " + featPart
	}

	// 最终截断
	return truncateString(name, maxLen)
}

// shortQualityFormat 简化音质格式
// (24bit/96.0kHz) → 24/96
// (16bit/44.1kHz) → 16/44
func shortQualityFormat(quality string) string {
	// 匹配 (XXbit/YY.YkHz) 格式
	re := regexp.MustCompile(`\((\d+)bit/([\d.]+)kHz\)`)
	matches := re.FindStringSubmatch(quality)

	if len(matches) >= 3 {
		bit := matches[1]
		khz := strings.TrimSuffix(matches[2], ".0") // 96.0 → 96
		return fmt.Sprintf("%s/%s", bit, khz)
	}

	// 降级处理：直接去掉括号
	return strings.Trim(quality, "()")
}

// formatStatusDetailed 格式化状态信息
// showSpeed=true: ↓54% 4.4MB/s
// showSpeed=false: ↓54%
func formatStatusDetailed(status string, showSpeed bool) string {
	// 解析百分比
	percentRe := regexp.MustCompile(`(\d+)%`)
	percentMatches := percentRe.FindStringSubmatch(status)

	// 解析速度
	speedRe := regexp.MustCompile(`\(?([\d.]+)\s*MB/s\)?`)
	speedMatches := speedRe.FindStringSubmatch(status)

	// 判断状态类型
	if strings.Contains(status, "下载中") || strings.Contains(status, "Downloading") {
		symbol := "↓"
		if percentMatches != nil {
			if showSpeed && speedMatches != nil {
				return fmt.Sprintf("%s%s%% %sMB/s", symbol, percentMatches[1], speedMatches[1])
			}
			return fmt.Sprintf("%s%s%%", symbol, percentMatches[1])
		}
		return symbol + "..."

	} else if strings.Contains(status, "解密中") || strings.Contains(status, "Decrypt") {
		symbol := "🔓"
		if percentMatches != nil {
			return fmt.Sprintf("%s%s%%", symbol, percentMatches[1])
		}
		return symbol + "..."

	} else if strings.Contains(status, "检测") || strings.Contains(status, "Check") {
		return "🔍..."

	} else if strings.Contains(status, "等待") || strings.Contains(status, "Wait") {
		return "⏸"

	} else if strings.Contains(status, "完成") || strings.Contains(status, "Complete") || strings.Contains(status, "Done") {
		return "✓"

	} else if strings.Contains(status, "错误") || strings.Contains(status, "Error") || strings.Contains(status, "失败") {
		// 错误信息截取前15字符
		errorMsg := strings.ReplaceAll(status, "错误: ", "")
		errorMsg = strings.ReplaceAll(errorMsg, "Error: ", "")
		if len([]rune(errorMsg)) > 15 {
			errorMsg = string([]rune(errorMsg)[:15]) + "..."
		}
		return "✗ " + errorMsg
	}

	// 默认：截断原始状态到25字符
	return truncateString(status, 25)
}

// truncateString 截断字符串到指定长度（考虑中文）
func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	if maxLen <= 3 {
		return "..."
	}

	return string(runes[:maxLen-3]) + "..."
}

// getVisualLength 计算字符串的视觉长度（去除ANSI颜色码）
func getVisualLength(s string) int {
	// 去除ANSI颜色码
	colorRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	plain := colorRegex.ReplaceAllString(s, "")

	// 计算rune数（正确处理中文等多字节字符）
	return len([]rune(plain))
}

// truncateToWidth 截断字符串到指定视觉宽度（保留颜色码）
func truncateToWidth(s string, maxWidth int) string {
	colorRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)

	// 提取所有颜色码及其位置
	type colorCode struct {
		pos  int
		code string
	}
	var colors []colorCode
	for _, match := range colorRegex.FindAllStringIndex(s, -1) {
		colors = append(colors, colorCode{
			pos:  match[0],
			code: s[match[0]:match[1]],
		})
	}

	// 去除颜色码
	plain := colorRegex.ReplaceAllString(s, "")
	runes := []rune(plain)

	// 如果不需要截断
	if len(runes) <= maxWidth {
		return s
	}

	// 截断
	if maxWidth <= 3 {
		return "..."
	}
	truncated := string(runes[:maxWidth-3]) + "..."

	// 尝试恢复颜色码（简化处理：只保留开头的颜色）
	if len(colors) > 0 && colors[0].pos == 0 {
		truncated = colors[0].code + truncated
	}

	return truncated
}
