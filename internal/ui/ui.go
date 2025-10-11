package ui

import (
	"bufio"
	"fmt"
	"main/internal/logger"
	"main/utils/structs"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"main/internal/core"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"
)

// UI暂停/恢复控制通道
var (
	suspendChan = make(chan struct{}, 1)
	resumeChan  = make(chan struct{}, 1)
	isSuspended = false
)

// Suspend 暂停UI更新（用于需要直接输出或交互的场景）
func Suspend() {
	select {
	case suspendChan <- struct{}{}:
		isSuspended = true
	default:
		// 已经在暂停状态，忽略
	}
}

// Resume 恢复UI更新
func Resume() {
	if isSuspended {
		select {
		case resumeChan <- struct{}{}:
			isSuspended = false
		default:
			// 已经在运行状态，忽略
		}
	}
}

// getTerminalWidth 获取终端宽度，如果获取失败则返回默认值
func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		// 获取失败或无效，使用保守的默认值
		return 80
	}
	return width
}

func RenderUI(done <-chan struct{}) {
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	// 首次更新标志：延迟初始化，避免预先打印换行符导致光标位置错误
	firstUpdate := true

	for {
		select {
		case <-done:
			return
		case <-suspendChan:
			// UI暂停，等待恢复信号
			<-resumeChan
		case <-ticker.C:
			PrintUI(firstUpdate)
			firstUpdate = false
		}
	}
}

func PrintUI(isFirstUpdate bool) {
	core.UiMutex.Lock()
	defer core.UiMutex.Unlock()

	if len(core.TrackStatuses) == 0 {
		return
	}

	var builder strings.Builder

	// 首次更新时打印占位换行符，后续更新时向上移动光标
	if isFirstUpdate {
		builder.WriteString(strings.Repeat("\n", len(core.TrackStatuses)))
	}

	builder.WriteString(fmt.Sprintf("\033[%dA", len(core.TrackStatuses)))

	// 动态获取终端宽度，确保内容不会因换行而产生额外的行
	terminalWidth := getTerminalWidth()
	colorRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)

	for _, ts := range core.TrackStatuses {
		displayName := ts.TrackName
		statusStrWithColor := ts.StatusColor(ts.Status)
		plainStatusStr := colorRegex.ReplaceAllString(statusStrWithColor, "")

		// 智能分级显示：根据终端宽度选择不同的显示格式
		var line string
		var prefixStr string
		var qualityStr string

		// 级别1：完整格式 (需要至少60字符)
		// Track 1 of 14: SongName (24bit/96.0kHz) - Status
		if terminalWidth >= 60 {
			prefixStr = fmt.Sprintf("Track %d of %d: ", ts.TrackNum, ts.TrackTotal)
			qualityStr = ts.Quality

			prefixRunes := len([]rune(prefixStr))
			suffixRunes := len([]rune(qualityStr)) + len([]rune(" - ")) + len([]rune(plainStatusStr))
			availableRunesForName := terminalWidth - prefixRunes - suffixRunes - 1

			if availableRunesForName < 10 {
				availableRunesForName = 10
			}

			displayNameRunes := []rune(displayName)
			if len(displayNameRunes) > availableRunesForName {
				if availableRunesForName > 3 {
					displayName = string(displayNameRunes[:availableRunesForName-3]) + "..."
				} else {
					displayName = "..."
				}
			}

			line = fmt.Sprintf("%s%s %s - %s", prefixStr, displayName, qualityStr, statusStrWithColor)

			// 级别2：紧凑格式，去掉音质信息 (需要至少40字符)
			// Track 1 of 14: SongName - Status
		} else if terminalWidth >= 40 {
			prefixStr = fmt.Sprintf("Track %d of %d: ", ts.TrackNum, ts.TrackTotal)

			prefixRunes := len([]rune(prefixStr))
			suffixRunes := len([]rune(" - ")) + len([]rune(plainStatusStr))
			availableRunesForName := terminalWidth - prefixRunes - suffixRunes - 1

			if availableRunesForName < 8 {
				availableRunesForName = 8
			}

			displayNameRunes := []rune(displayName)
			if len(displayNameRunes) > availableRunesForName {
				if availableRunesForName > 3 {
					displayName = string(displayNameRunes[:availableRunesForName-3]) + "..."
				} else {
					displayName = "..."
				}
			}

			line = fmt.Sprintf("%s%s - %s", prefixStr, displayName, statusStrWithColor)

			// 级别3：极简格式 (需要至少25字符)
			// [1/14] SongName - Status
		} else if terminalWidth >= 25 {
			prefixStr = fmt.Sprintf("[%d/%d] ", ts.TrackNum, ts.TrackTotal)

			prefixRunes := len([]rune(prefixStr))
			suffixRunes := len([]rune(" - ")) + len([]rune(plainStatusStr))
			availableRunesForName := terminalWidth - prefixRunes - suffixRunes - 1

			if availableRunesForName < 5 {
				availableRunesForName = 5
			}

			displayNameRunes := []rune(displayName)
			if len(displayNameRunes) > availableRunesForName {
				if availableRunesForName > 3 {
					displayName = string(displayNameRunes[:availableRunesForName-3]) + "..."
				} else {
					displayName = "..."
				}
			}

			line = fmt.Sprintf("%s%s - %s", prefixStr, displayName, statusStrWithColor)

			// 级别4：最小格式 (宽度小于25字符)
			// [1/14] Status
		} else {
			prefixStr = fmt.Sprintf("[%d/%d] ", ts.TrackNum, ts.TrackTotal)
			line = fmt.Sprintf("%s%s", prefixStr, statusStrWithColor)
		}

		// 使用 \r\033[K 清除当前行，然后打印内容
		// 最后做一次安全检查，确保不超过终端宽度
		lineRunes := []rune(colorRegex.ReplaceAllString(line, ""))
		if len(lineRunes) > terminalWidth {
			// 如果仍然超过（理论上不应该发生），强制截断
			line = string([]rune(line)[:terminalWidth-3]) + "..."
		}

		builder.WriteString(fmt.Sprintf("\r\033[K%s\n", line))
	}
	fmt.Print(builder.String()) // OK: UI渲染核心，必须使用fmt.Print输出到stdout
}

func UpdateStatus(index int, status string, sColor func(a ...interface{}) string) {
	core.UiMutex.Lock()
	defer core.UiMutex.Unlock()
	if index < len(core.TrackStatuses) {
		// 去重：只有当状态真正改变时才更新
		// 这避免了重复的进度更新导致日志刷屏
		if core.TrackStatuses[index].Status != status {
			core.TrackStatuses[index].Status = status
			core.TrackStatuses[index].StatusColor = sColor
		}
	}
}

func SelectTracks(meta *structs.AutoGenerated, storefront, urlArg_i string) []int {
	trackTotal := len(meta.Data[0].Relationships.Tracks.Data)
	arr := make([]int, trackTotal)
	for i := 0; i < trackTotal; i++ {
		arr[i] = i + 1
	}
	selected := []int{}

	if core.Dl_song {
		found := false
		for i, track := range meta.Data[0].Relationships.Tracks.Data {
			if urlArg_i == track.ID {
				selected = append(selected, i+1)
				found = true
				break
			}
		}
		if !found {
			logger.Error("指定的单曲ID未在专辑中找到")
			return nil
		}
	} else if !core.Dl_select {
		selected = arr
	} else {
		var data [][]string
		for trackNum, track := range meta.Data[0].Relationships.Tracks.Data {
			trackNum++
			var trackName string
			if meta.Data[0].Type == "albums" {
				trackName = fmt.Sprintf("%02d. %s", track.Attributes.TrackNumber, track.Attributes.Name)
			} else {
				trackName = fmt.Sprintf("%s - %s", track.Attributes.Name, track.Attributes.ArtistName)
			}
			data = append(data, []string{fmt.Sprint(trackNum),
				trackName,
				track.Attributes.ContentRating,
				track.Type})
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"", "Track Name", "Rating", "Type"})
		table.SetRowLine(false)
		table.SetCaption(meta.Data[0].Type == "albums", fmt.Sprintf("Storefront: %s, %d tracks missing", strings.ToUpper(storefront), meta.Data[0].Attributes.TrackCount-trackTotal))
		table.SetHeaderColor(tablewriter.Colors{},
			tablewriter.Colors{tablewriter.FgRedColor, tablewriter.Bold},
			tablewriter.Colors{tablewriter.FgBlackColor, tablewriter.Bold},
			tablewriter.Colors{tablewriter.FgBlackColor, tablewriter.Bold})

		table.SetColumnColor(tablewriter.Colors{tablewriter.FgCyanColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgRedColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor})
		for _, row := range data {
			if row[2] == "explicit" {
				row[2] = "E"
			} else if row[2] == "clean" {
				row[2] = "C"
			} else {
				row[2] = "None"
			}
			if row[3] == "music-videos" {
				row[3] = "MV"
			} else if row[3] == "songs" {
				row[3] = "SONG"
			}
			table.Append(row)
		}
		table.Render()
		logger.Info("Please select from the track options above (multiple options separated by commas, ranges supported, or type 'all' to select all)")
		cyanColor := color.New(color.FgCyan)
		cyanColor.Print("select: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("读取输入错误: %v", err)
		}
		input = strings.TrimSpace(input)
		if input == "all" {
			selected = arr
		} else {
			selectedOptions := [][]string{}
			parts := strings.Split(input, ",")
			for _, part := range parts {
				if strings.Contains(part, "-") {
					rangeParts := strings.Split(part, "-")
					selectedOptions = append(selectedOptions, rangeParts)
				} else {
					selectedOptions = append(selectedOptions, []string{part})
				}
			}
			for _, opt := range selectedOptions {
				if len(opt) == 1 {
					num, err := strconv.Atoi(opt[0])
					if err != nil {
						continue
					}
					if num > 0 && num <= len(arr) {
						selected = append(selected, num)
					}
				} else if len(opt) == 2 {
					start, err1 := strconv.Atoi(opt[0])
					end, err2 := strconv.Atoi(opt[1])
					if err1 != nil || err2 != nil {
						continue
					}
					if start < 1 || end > len(arr) || start > end {
						continue
					}
					for i := start; i <= end; i++ {
						selected = append(selected, i)
					}
				}
			}
		}
	}
	return selected
}
