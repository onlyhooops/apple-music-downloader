package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"main/internal/api"
	"main/internal/constants"
	"main/internal/core"
	"main/internal/downloader"
	"main/internal/logger"
	"main/internal/network"
	"main/internal/parser"
	"main/internal/progress"
	"main/internal/ui"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

// 版本信息（编译时通过 ldflags 注入）
var (
	Version   = "v1.3.0"  // 版本号
	BuildTime = "unknown" // 编译时间
	GitCommit = "unknown" // Git提交哈希
)

// loadDevEnv 自动加载 dev.env 文件中的环境变量
func loadDevEnv() {
	envFile := "dev.env"
	data, err := os.ReadFile(envFile)
	if err != nil {
		// dev.env 不存在是正常的，不报错
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// 移除引号
			value = strings.Trim(value, "\"'")
			os.Setenv(key, value)
		}
	}
}

func handleSingleMV(urlRaw string) {
	if core.Debug_mode {
		return
	}
	storefront, albumId := parser.CheckUrlMv(urlRaw)
	accountForMV, err := core.GetAccountForStorefront(storefront)
	if err != nil {
		logger.Error("MV 下载失败: %v", err)
		core.SharedLock.Lock()
		core.Counter.Error++
		core.SharedLock.Unlock()
		return
	}

	core.SharedLock.Lock()
	core.Counter.Total++
	core.SharedLock.Unlock()

	if len(accountForMV.MediaUserToken) < constants.MinTokenLength {
		logger.Error("MV 下载失败: MediaUserToken 无效或过短")
		logger.Info("提示: 请确保在 dev.env 中配置了有效的 APPLE_MUSIC_MEDIA_USER_TOKEN_CN")
		core.SharedLock.Lock()
		core.Counter.Error++
		core.SharedLock.Unlock()
		return
	}

	if _, err := exec.LookPath("mp4decrypt"); err != nil {
		logger.Error("MV 下载失败: 未找到 mp4decrypt 工具")
		logger.Info("提示: 请安装 mp4decrypt (https://www.bento4.com/downloads/)")
		core.SharedLock.Lock()
		core.Counter.Error++
		core.SharedLock.Unlock()
		return
	}

	mvInfo, err := api.GetMVInfoFromAdam(albumId, accountForMV, storefront)
	if err != nil {
		logger.Error("获取 MV 信息失败: %v", err)
		core.SharedLock.Lock()
		core.Counter.Error++
		core.SharedLock.Unlock()
		return
	}

	// 输出MV信息
	core.SafePrintf("🎤 歌手: %s\n", mvInfo.Data[0].Attributes.ArtistName)
	core.SafePrintf("🎬 MV: %s\n", mvInfo.Data[0].Attributes.Name)

	// 提取发行年份
	var releaseYear string
	if len(mvInfo.Data[0].Attributes.ReleaseDate) >= 4 {
		releaseYear = mvInfo.Data[0].Attributes.ReleaseDate[:4]
		core.SafePrintf("📅 发行年份: %s\n", releaseYear)
	}

	var artistFolder string
	if core.Config.ArtistFolderFormat != "" {
		artistFolder = strings.NewReplacer(
			"{UrlArtistName}", core.LimitString(mvInfo.Data[0].Attributes.ArtistName),
			"{ArtistName}", core.LimitString(mvInfo.Data[0].Attributes.ArtistName),
			"{ArtistId}", "",
		).Replace(core.Config.ArtistFolderFormat)
	}
	sanitizedArtistFolder := core.ForbiddenNames.ReplaceAllString(artistFolder, "_")

	// Use MVSaveFolder if configured, otherwise fallback to AlacSaveFolder
	mvSaveFolder := core.Config.MVSaveFolder
	if mvSaveFolder == "" {
		mvSaveFolder = core.Config.AlacSaveFolder
	}

	// 应用缓存机制
	cachePath, finalPath, usingCache := downloader.GetCacheBasePath(mvSaveFolder, albumId)

	mvOutPath, mvResolution, err := downloader.MvDownloader(albumId, cachePath, sanitizedArtistFolder, "", storefront, nil, accountForMV)

	// 分辨率信息已在 MvDownloader 内部显示，这里不再重复显示
	_ = mvResolution

	// 如果使用缓存且下载成功，移动文件到最终位置
	if err == nil && usingCache && mvOutPath != "" {
		// 计算最终路径
		relPath, _ := filepath.Rel(cachePath, mvOutPath)
		finalMvPath := filepath.Join(finalPath, relPath)

		// 移动文件
		core.SafePrintf("\n📤 正在从缓存转移MV文件到目标位置...\n")
		if moveErr := downloader.SafeMoveFile(mvOutPath, finalMvPath); moveErr != nil {
			// 检查是否是文件已存在的情况
			if strings.Contains(moveErr.Error(), "目标文件已存在") {
				logger.Info("✅ MV 文件已存在，跳过下载")
				logger.Info("💾 保存路径: %s", finalMvPath)
				// 文件已存在视为成功，清理缓存
				os.RemoveAll(cachePath)
			} else {
				logger.Error("从缓存移动MV文件失败: %v", moveErr)
				err = moveErr
			}
		} else {
			core.SafePrintf("📥 MV文件转移完成！\n")
			core.SafePrintf("💾 保存路径: %s\n", finalMvPath)

			// 清理缓存目录
			mvCacheDir := filepath.Dir(mvOutPath)
			for mvCacheDir != cachePath && mvCacheDir != "." && mvCacheDir != "/" {
				if os.Remove(mvCacheDir) != nil {
					break
				}
				mvCacheDir = filepath.Dir(mvCacheDir)
			}
		}
	} else if err == nil && !usingCache && mvOutPath != "" {
		// 未使用缓存，直接保存
		core.SafePrintf("\n📥 MV下载完成！\n")
		core.SafePrintf("💾 保存路径: %s\n", mvOutPath)
	}

	// 如果出错且使用了缓存，清理缓存
	if err != nil && usingCache {
		os.RemoveAll(cachePath)
	}

	if err != nil {
		core.SharedLock.Lock()
		core.Counter.Error++
		core.SharedLock.Unlock()
		return
	}
	core.SharedLock.Lock()
	core.Counter.Success++
	core.SharedLock.Unlock()
}

func processURL(ctx context.Context, urlRaw string, wg *sync.WaitGroup, semaphore chan struct{}, currentTask int, totalTasks int, notifier *progress.ProgressNotifier) (string, string, error) {
	if wg != nil {
		defer wg.Done()
	}
	if semaphore != nil {
		defer func() { <-semaphore }()
	}

	// 检查 context 是否已取消
	select {
	case <-ctx.Done():
		logger.Info("任务已取消: %s", urlRaw)
		return "", "", ctx.Err()
	default:
	}

	if totalTasks > 1 {
		core.SafePrintf("🧾 [%d/%d] 开始处理: %s\n", currentTask, totalTasks, urlRaw)
	}

	var storefront, albumId string
	var albumName string

	if strings.Contains(urlRaw, "/music-video/") {
		handleSingleMV(urlRaw)
		return "", "", nil
	}

	if strings.Contains(urlRaw, "/song/") {
		tempStorefront, _ := parser.CheckUrlSong(urlRaw)
		accountForSong, err := core.GetAccountForStorefront(tempStorefront)
		if err != nil {
			logger.Error("获取歌曲信息失败 for %s: %v", urlRaw, err)
			return "", "", err
		}
		urlRaw, err = api.GetUrlSong(urlRaw, accountForSong)
		if err != nil {
			logger.Error("获取歌曲链接失败 for %s: %v", urlRaw, err)
			return "", "", err
		}
		core.Dl_song = true
	}

	if strings.Contains(urlRaw, "/playlist/") {
		storefront, albumId = parser.CheckUrlPlaylist(urlRaw)
	} else {
		storefront, albumId = parser.CheckUrl(urlRaw)
	}

	if albumId == "" {
		err := fmt.Errorf("无效的URL")
		logger.Warn("无效的URL: %s", urlRaw)
		return "", "", err
	}

	// 获取专辑信息
	mainAccount, err := core.GetAccountForStorefront(storefront)
	if err == nil {
		meta, err := api.GetMeta(albumId, mainAccount, storefront)
		if err == nil && len(meta.Data) > 0 {
			albumName = meta.Data[0].Attributes.Name
		}
	}

	parse, err := url.Parse(urlRaw)
	if err != nil {
		log.Printf("解析URL失败 %s: %v", urlRaw, err)
		return albumId, albumName, err
	}
	var urlArg_i = parse.Query().Get("i")
	err = downloader.Rip(albumId, storefront, urlArg_i, urlRaw, notifier)
	if err != nil {
		core.SafePrintf("专辑下载失败: %s -> %v\n", urlRaw, err)
		return albumId, albumName, err
	} else {
		if totalTasks > 1 {
			core.SafePrintf("✅ [%d/%d] 任务完成: %s\n", currentTask, totalTasks, urlRaw)
		}
		return albumId, albumName, nil
	}
}

// parseTxtFile 从TXT文件中解析URL列表
func parseTxtFile(filePath string) ([]string, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	lines := strings.Split(string(fileBytes), "\n")
	var urls []string
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		// 跳过空行和注释行（以#开头）
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}
		// 支持一行多个链接（空格分隔）
		linksInLine := strings.Fields(trimmedLine)
		for _, link := range linksInLine {
			link = strings.TrimSpace(link)
			if link != "" {
				urls = append(urls, link)
			}
		}
	}
	return urls, nil
}

func runDownloads(ctx context.Context, initialUrls []string, isBatch bool, taskFile string, notifier *progress.ProgressNotifier) {
	var finalUrls []string

	// 显示输入链接统计
	if isBatch && len(initialUrls) > 0 {
		core.SafePrintf("📋 初始链接总数: %d\n", len(initialUrls))
		core.SafePrintf("🔄 开始预处理链接...\n\n")
	}

	for _, urlRaw := range initialUrls {
		if strings.Contains(urlRaw, "/artist/") {
			core.SafePrintf("🔍 正在解析歌手页面: %s\n", urlRaw)
			artistAccount := &core.Config.Accounts[0]
			urlArtistName, urlArtistID, err := api.GetUrlArtistName(urlRaw, artistAccount)
			if err != nil {
				core.SafePrintf("获取歌手名称失败 for %s: %v\n", urlRaw, err)
				continue
			}

			core.Config.ArtistFolderFormat = strings.NewReplacer(
				"{UrlArtistName}", core.LimitString(urlArtistName),
				"{ArtistId}", urlArtistID,
			).Replace(core.Config.ArtistFolderFormat)

			albumArgs, err := api.CheckArtist(urlRaw, artistAccount, "albums")
			if err != nil {
				core.SafePrintf("获取歌手专辑失败 for %s: %v\n", urlRaw, err)
			} else {
				finalUrls = append(finalUrls, albumArgs...)
				core.SafePrintf("📀 从歌手 %s 页面添加了 %d 张专辑到队列。\n", urlArtistName, len(albumArgs))
			}

			mvArgs, err := api.CheckArtist(urlRaw, artistAccount, "music-videos")
			if err != nil {
				core.SafePrintf("获取歌手MV失败 for %s: %v\n", urlRaw, err)
			} else {
				finalUrls = append(finalUrls, mvArgs...)
				core.SafePrintf("🎬 从歌手 %s 页面添加了 %d 个MV到队列。\n", urlArtistName, len(mvArgs))
			}
		} else {
			finalUrls = append(finalUrls, urlRaw)
		}
	}

	if len(finalUrls) == 0 {
		logger.Warn("队列中没有有效的链接可供下载。")
		return
	}

	// 如果最终链接数量>1，也应该视为批量模式（支持工作-休息循环）
	if len(finalUrls) > 1 {
		isBatch = true
	}

	totalTasks := len(finalUrls)

	// 处理 --start 参数
	startIndex := 0 // 实际数组索引（从0开始）
	if core.StartFrom > 0 {
		if core.StartFrom > totalTasks {
			core.SafePrintf("⚠️  起始位置 %d 超过了总任务数 %d，将从第 1 个开始\n", core.StartFrom, totalTasks)
			core.StartFrom = 1
		} else {
			startIndex = core.StartFrom - 1 // 用户输入从1开始，转换为0开始的索引
			skippedCount := startIndex
			core.SafePrintf("⏭️  跳过前 %d 个任务，从第 %d 个开始下载\n", skippedCount, core.StartFrom)
			finalUrls = finalUrls[startIndex:] // 跳过前面的链接
		}
	}

	// 准备下载任务
	totalTasks = len(finalUrls)

	// 保存原始总数用于显示
	originalTotalTasks := len(initialUrls)

	if isBatch {
		core.SafePrintf("\n📋 ===== 开始下载任务 =====\n")
		if len(initialUrls) != totalTasks {
			core.SafePrintf("📝 预处理完成: %d → %d 任务\n", len(initialUrls), originalTotalTasks)
		} else {
			core.SafePrintf("📝 任务总数: %d\n", originalTotalTasks)
		}
		if core.StartFrom > 0 {
			core.SafePrintf("📝 实际下载: 第 %d-%d 个（共 %d 个）\n", core.StartFrom, originalTotalTasks, totalTasks)
		}
		core.SafePrintf("⚡ 执行模式: 串行模式\n")
		core.SafePrintf("📦 专辑内并发: 由配置控制\n")
		core.SafePrintf("=============================\n")
	} else {
		core.SafePrintf("📋 开始下载任务\n📝 总数: %d\n", originalTotalTasks)
	}

	// 批量模式：串行执行（按链接顺序依次下载）
	// 专辑内歌曲并发数由配置文件控制 (lossless_downloadthreads 等)

	// 工作-休息循环机制
	var workStartTime time.Time
	if isBatch && core.Config.WorkRestEnabled {
		workStartTime = time.Now()
		logger.Debug("[工作-休息] 循环已启用: 工作=%d分钟, 休息=%d分钟, 任务数=%d",
			core.Config.WorkDurationMinutes, core.Config.RestDurationMinutes, len(finalUrls))
		core.SafePrintf("⏰ 工作-休息循环: 工作 %d 分钟 / 休息 %d 分钟\n",
			core.Config.WorkDurationMinutes,
			core.Config.RestDurationMinutes)
		core.SafePrintf("⏱️  工作开始: %s\n", workStartTime.Format("15:04:05"))
	} else if isBatch {
		logger.Debug("[工作-休息] 循环未启用: WorkRestEnabled=%v, 任务数=%d",
			core.Config.WorkRestEnabled, len(finalUrls))
	}

	for i, urlToProcess := range finalUrls {
		// 检查 context 是否已取消
		select {
		case <-ctx.Done():
			logger.Warn("下载已中断，已完成 %d/%d 个任务", i, len(finalUrls))
			return
		default:
		}

		// 计算实际的任务编号（考虑 --start 参数）
		actualTaskNum := i + 1 + startIndex    // 实际编号 = 当前索引 + 1 + 跳过的数量
		originalTotalTasks := len(initialUrls) // 原始总数（包括被跳过的）

		_, _, _ = processURL(ctx, urlToProcess, nil, nil, actualTaskNum, originalTotalTasks, notifier)

		// 任务之间添加视觉间隔（最后一个任务不需要）
		if isBatch && i < len(finalUrls)-1 {
			core.SafePrintf("\n%s\n", strings.Repeat("=", constants.VisualSeparatorLength))
		}

		// 工作-休息循环检查（在任务完成后）
		if isBatch && core.Config.WorkRestEnabled && i < len(finalUrls)-1 {
			elapsed := time.Since(workStartTime)
			workDuration := time.Duration(core.Config.WorkDurationMinutes) * time.Minute

			logger.Debug("[工作-休息] 检查点: 已工作 %.1f 分钟 / 阈值 %d 分钟, 任务进度 %d/%d",
				elapsed.Minutes(), core.Config.WorkDurationMinutes, i+1, len(finalUrls))

			if elapsed >= workDuration {
				logger.Info("[工作-休息] 达到工作时长阈值，准备进入休息")
				// 工作时间已到，需要休息
				restDuration := time.Duration(core.Config.RestDurationMinutes) * time.Minute

				cyan := color.New(color.FgCyan, color.Bold)
				yellow := color.New(color.FgYellow)
				green := color.New(color.FgGreen)

				core.SafePrintf("\n%s\n", strings.Repeat("=", constants.VisualSeparatorLength))
				cyan.Printf("⏸️  已工作 %d 分钟，进入休息\n", core.Config.WorkDurationMinutes)
				yellow.Printf("😴 休息 %d 分钟\n", core.Config.RestDurationMinutes)
				core.SafePrintf("📊 已完成: %d/%d\n", i+1, totalTasks)
				core.SafePrintf("⏰ 当前时间: %s\n", time.Now().Format("15:04:05"))
				core.SafePrintf("⏱️  恢复时间: %s\n", time.Now().Add(restDuration).Format("15:04:05"))
				core.SafePrintf("%s\n", strings.Repeat("=", constants.VisualSeparatorLength))

				// 休息倒计时（每30秒提示一次）
				restTicker := time.NewTicker(constants.RestTickerInterval)
				restTimer := time.NewTimer(restDuration)
				restStartTime := time.Now()

				restDone := false
				for !restDone {
					select {
					case <-restTimer.C:
						// 休息时间结束
						restDone = true
					case <-restTicker.C:
						// 显示剩余时间
						remainingTime := restDuration - time.Since(restStartTime)
						if remainingTime > 0 {
							core.SafePrintf("⏳ 休息中... 剩余时间: %.0f 分钟 %.0f 秒\n",
								remainingTime.Minutes(),
								remainingTime.Seconds()-remainingTime.Minutes()*60)
						}
					}
				}
				restTicker.Stop()

				// 休息结束，重新开始计时
				workStartTime = time.Now()
				core.SafePrintf("\n%s\n", strings.Repeat("=", constants.VisualSeparatorLength))
				green.Printf("✅ 休息完毕，继续任务\n")
				core.SafePrintf("⏱️  工作开始: %s\n", workStartTime.Format("15:04:05"))
				core.SafePrintf("%s\n", strings.Repeat("=", constants.VisualSeparatorLength))
			}
		}
	}

}

func main() {
	// 自动加载 dev.env 文件（如果存在）
	loadDevEnv()

	// 打印版本信息
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	fmt.Println(strings.Repeat("=", constants.BannerSeparatorLength)) // OK: 程序启动横幅
	cyan.Printf("🎵 Apple Music Downloader %s\n", Version)

	// 显示编译时间（本地时间）
	if BuildTime != "unknown" && BuildTime != "" {
		yellow.Printf("📅 编译时间: %s\n", BuildTime)
	} else {
		yellow.Printf("📅 编译时间: %s\n", BuildTime)
	}

	if GitCommit != "unknown" {
		yellow.Printf("🔖 Git提交: %s\n", GitCommit)
	}
	fmt.Println(strings.Repeat("=", constants.BannerSeparatorLength)) // OK: 程序启动横幅
	fmt.Println()                                                     // OK: 程序启动横幅

	core.InitFlags()

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: %s [选项] [url1 url2 ... | file.txt ...]\n", os.Args[0])
		logger.Info("如果没有提供URL或文件，程序将进入交互模式。")
		logger.Info("")
		logger.Info("支持的启动方式:")
		logger.Info("  1. 交互模式: 运行程序后输入链接或TXT文件路径")
		logger.Info("  2. 单链接模式: ./程序名 <url>")
		logger.Info("  3. 多链接模式: ./程序名 <url1> <url2> ...")
		logger.Info("  4. TXT文件模式: ./程序名 <file.txt>")
		logger.Info("  5. 混合模式: ./程序名 <url1> <file.txt> <url2> ...")
		logger.Info("")
		logger.Info("TXT文件格式:")
		logger.Info("  - 支持单行单链接（传统格式）")
		logger.Info("  - 支持单行多链接（空格分隔）")
		logger.Info("  - 支持注释行（以#开头）")
		logger.Info("  - 空行会被自动跳过")
		logger.Info("")
		logger.Info("选项:")
		pflag.PrintDefaults()
	}

	pflag.Parse()

	err := core.LoadConfig(core.ConfigPath)
	if err != nil {
		if os.IsNotExist(err) && core.ConfigPath == "config.yaml" {
			// OK: logger还未初始化，必须使用fmt
			fmt.Println("错误: 默认配置文件 config.yaml 未找到。")
			pflag.Usage()
			return
		}
		// OK: logger还未初始化，必须使用fmt
		fmt.Printf("加载配置文件 %s 失败: %v\n", core.ConfigPath, err)
		return
	}

	// 初始化logger系统
	loggerCfg := logger.Config{
		Level:         core.Config.Logging.Level,
		Output:        core.Config.Logging.Output,
		ShowTimestamp: core.Config.Logging.ShowTimestamp,
	}
	if err := logger.InitFromConfig(loggerCfg); err != nil {
		// OK: 这里不能用logger.Error，因为logger初始化失败
		fmt.Printf("初始化logger失败: %v\n", err)
		return
	}

	// 初始化网络客户端（包括本地 wrapper 优化）
	network.InitializeClients(&core.Config)

	// 创建可取消的 context 用于优雅退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 设置信号处理，确保程序退出时清理资源
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigChan
		yellow := color.New(color.FgYellow)
		yellow.Printf("\n\n⚠️  收到中断信号，正在安全退出...\n")

		// 取消所有进行中的任务
		cancel()

		// 等待清理完成
		time.Sleep(constants.CleanupWaitSeconds * time.Second)

		yellow.Printf("✅ 清理完成\n")
		yellow.Printf("👋 再见！\n")
		os.Exit(0)
	}()

	// 创建进度通知器并注册UI监听器
	progressNotifier := progress.NewNotifier()
	uiListener := ui.NewUIProgressListener()
	progressNotifier.AddListener(uiListener)
	logger.Debug("Progress notifier initialized with UI listener")

	if core.OutputPath != "" {
		core.Config.AlacSaveFolder = core.OutputPath
		core.Config.AtmosSaveFolder = core.OutputPath
	}

	token, err := api.GetToken()
	if err != nil {
		if len(core.Config.Accounts) > 0 && core.Config.Accounts[0].AuthorizationToken != "" && core.Config.Accounts[0].AuthorizationToken != "your-authorization-token" {
			token = strings.Replace(core.Config.Accounts[0].AuthorizationToken, "Bearer ", "", -1)
		} else {
			logger.Error("获取开发者 token 失败。")
			return
		}
	}
	core.DeveloperToken = token

	args := pflag.Args()
	if len(args) == 0 {
		logger.Info("请输入专辑链接或TXT文件路径: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			logger.Info("未输入内容，程序退出。")
			return
		}

		if strings.HasSuffix(strings.ToLower(input), ".txt") {
			if _, err := os.Stat(input); err == nil {
				urls, err := parseTxtFile(input)
				if err != nil {
					logger.Error("读取文件 %s 失败: %v", input, err)
					return
				}
				logger.Info("📊 从文件 %s 中解析到 %d 个链接\n", input, len(urls))
				runDownloads(ctx, urls, true, input, progressNotifier)
			} else {
				logger.Error("错误: 文件不存在 %s", input)
				return
			}
		} else {
			runDownloads(ctx, []string{input}, false, "", progressNotifier)
		}
	} else {
		// 处理命令行参数：支持TXT文件或直接的URL列表
		var urls []string
		isBatch := false
		var taskFile string

		for _, arg := range args {
			if strings.HasSuffix(strings.ToLower(arg), ".txt") {
				// 参数是TXT文件
				if _, err := os.Stat(arg); err == nil {
					fileUrls, err := parseTxtFile(arg)
					if err != nil {
						logger.Error("读取文件 %s 失败: %v", arg, err)
						continue
					}
					logger.Info("📊 从文件 %s 中解析到 %d 个链接", arg, len(fileUrls))
					urls = append(urls, fileUrls...)
					isBatch = true
					// 记录第一个txt文件作为任务文件
					if taskFile == "" {
						taskFile = arg
					}
				} else {
					logger.Error("错误: 文件不存在 %s", arg)
				}
			} else {
				// 参数是URL
				urls = append(urls, arg)
			}
		}

		if len(urls) > 1 {
			isBatch = true
		}

		if len(urls) > 0 {
			if isBatch {
				logger.Info("")
			}
			runDownloads(ctx, urls, isBatch, taskFile, progressNotifier)
		} else {
			logger.Warn("没有有效的链接可供处理。")
		}
	}

	logger.Info("\n📦 已完成: %d/%d | 警告: %d | 错误: %d", core.Counter.Success, core.Counter.Total, core.Counter.Unavailable+core.Counter.NotSong, core.Counter.Error)
	if core.Counter.Error > 0 {
		logger.Warn("部分任务在执行过程中出错，请检查上面的日志记录。")
	}
}
