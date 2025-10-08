package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"

	"main/internal/api"
	"main/internal/core"
	"main/internal/downloader"
	"main/internal/logger"
	"main/internal/parser"

	"github.com/spf13/pflag"
)

func handleSingleMV(urlRaw string) {
	if core.Debug_mode {
		return
	}
	
	logger.Section("开始下载MV")
	
	storefront, albumId := parser.CheckUrlMv(urlRaw)
	accountForMV, err := core.GetAccountForStorefront(storefront)
	if err != nil {
		logger.Error("MV 下载失败: %v", err)
		core.SharedLock.Lock()
		core.Counter.Error++
		core.SharedLock.Unlock()
		logger.SectionEnd()
		return
	}

	core.SharedLock.Lock()
	core.Counter.Total++
	core.SharedLock.Unlock()
	if len(accountForMV.MediaUserToken) <= 50 {
		logger.Error("账户 MediaUserToken 无效")
		core.SharedLock.Lock()
		core.Counter.Error++
		core.SharedLock.Unlock()
		logger.SectionEnd()
		return
	}
	if _, err := exec.LookPath("mp4decrypt"); err != nil {
		logger.Error("未找到 mp4decrypt 工具")
		core.SharedLock.Lock()
		core.Counter.Error++
		core.SharedLock.Unlock()
		logger.SectionEnd()
		return
	}

	mvInfo, err := api.GetMVInfoFromAdam(albumId, accountForMV, storefront)
	if err != nil {
		logger.Error("获取 MV 信息失败: %v", err)
		core.SharedLock.Lock()
		core.Counter.Error++
		core.SharedLock.Unlock()
		logger.SectionEnd()
		return
	}
	
	// 显示MV信息
	logger.Plain("🎤 歌手: %s\n", mvInfo.Data[0].Attributes.ArtistName)
	logger.Plain("🎬 MV: %s\n", mvInfo.Data[0].Attributes.Name)
	if len(mvInfo.Data[0].Attributes.ReleaseDate) >= 4 {
		logger.Plain("📅 发行: %s\n", mvInfo.Data[0].Attributes.ReleaseDate[:4])
	}
	logger.Plain("\n")

	var artistFolder string
	if core.Config.ArtistFolderFormat != "" {
		artistFolder = strings.NewReplacer(
			"{UrlArtistName}", core.LimitString(mvInfo.Data[0].Attributes.ArtistName),
			"{ArtistName}", core.LimitString(mvInfo.Data[0].Attributes.ArtistName),
			"{ArtistId}", "",
		).Replace(core.Config.ArtistFolderFormat)
	}
	sanitizedArtistFolder := core.ForbiddenNames.ReplaceAllString(artistFolder, "_")

	// 使用 MvSaveFolder 配置（如果有），否则回退到 AlacSaveFolder
	mvBasePath := core.Config.MvSaveFolder
	if mvBasePath == "" {
		mvBasePath = core.Config.AlacSaveFolder
	}

	_, err = downloader.MvDownloader(albumId, mvBasePath, sanitizedArtistFolder, "", storefront, nil, accountForMV)

	if err != nil {
		logger.Error("MV 下载失败: %v", err)
		core.SharedLock.Lock()
		core.Counter.Error++
		core.SharedLock.Unlock()
		logger.SectionEnd()
		return
	}
	
	logger.Success("MV 下载完成")
	core.SharedLock.Lock()
	core.Counter.Success++
	core.SharedLock.Unlock()
	logger.SectionEnd()
	logger.Plain("\n")
}

func processURL(urlRaw string, wg *sync.WaitGroup, semaphore chan struct{}, currentTask int, totalTasks int) {
	if wg != nil {
		defer wg.Done()
	}
	if semaphore != nil {
		defer func() { <-semaphore }()
	}

	var storefront, albumId string

	logger.TaskProgress(currentTask, totalTasks, fmt.Sprintf("开始处理: %s", urlRaw))

	if strings.Contains(urlRaw, "/music-video/") {
		handleSingleMV(urlRaw)
		return
	}

	if strings.Contains(urlRaw, "/song/") {
		tempStorefront, _ := parser.CheckUrlSong(urlRaw)
		accountForSong, err := core.GetAccountForStorefront(tempStorefront)
		if err != nil {
			logger.Error("获取歌曲信息失败 for %s: %v", urlRaw, err)
			return
		}
		urlRaw, err = api.GetUrlSong(urlRaw, accountForSong)
		if err != nil {
			logger.Error("获取歌曲链接失败 for %s: %v", urlRaw, err)
			return
		}
		core.Dl_song = true
	}

	if strings.Contains(urlRaw, "/playlist/") {
		storefront, albumId = parser.CheckUrlPlaylist(urlRaw)
	} else {
		storefront, albumId = parser.CheckUrl(urlRaw)
	}

	if albumId == "" {
		logger.Error("无效的URL: %s", urlRaw)
		return
	}

	parse, err := url.Parse(urlRaw)
	if err != nil {
		logger.Error("解析URL失败 %s: %v", urlRaw, err)
		return
	}
	var urlArg_i = parse.Query().Get("i")
	err = downloader.Rip(albumId, storefront, urlArg_i, urlRaw)
	if err != nil {
		logger.Error("专辑下载失败: %s -> %v", urlRaw, err)
	} else {
		logger.TaskProgress(currentTask, totalTasks, fmt.Sprintf("任务完成: %s", urlRaw))
	}
}

func runDownloads(initialUrls []string, isBatch bool) {
	var finalUrls []string

	for _, urlRaw := range initialUrls {
		if strings.Contains(urlRaw, "/artist/") {
			logger.Info("正在解析歌手页面: %s", urlRaw)
			artistAccount := &core.Config.Accounts[0]
			urlArtistName, urlArtistID, err := api.GetUrlArtistName(urlRaw, artistAccount)
			if err != nil {
				logger.Error("获取歌手名称失败 for %s: %v", urlRaw, err)
				continue
			}

			core.Config.ArtistFolderFormat = strings.NewReplacer(
				"{UrlArtistName}", core.LimitString(urlArtistName),
				"{ArtistId}", urlArtistID,
			).Replace(core.Config.ArtistFolderFormat)

			albumArgs, err := api.CheckArtist(urlRaw, artistAccount, "albums")
			if err != nil {
				logger.Error("获取歌手专辑失败 for %s: %v", urlRaw, err)
			} else {
				finalUrls = append(finalUrls, albumArgs...)
				logger.Info("从歌手 %s 页面添加了 %d 张专辑到队列。", urlArtistName, len(albumArgs))
			}

			mvArgs, err := api.CheckArtist(urlRaw, artistAccount, "music-videos")
			if err != nil {
				logger.Error("获取歌手MV失败 for %s: %v", urlRaw, err)
			} else {
				finalUrls = append(finalUrls, mvArgs...)
				logger.Info("从歌手 %s 页面添加了 %d 个MV到队列。", urlArtistName, len(mvArgs))
			}
		} else {
			finalUrls = append(finalUrls, urlRaw)
		}
	}

	if len(finalUrls) == 0 {
		logger.Warn("队列中没有有效的链接可供下载。")
		return
	}

	numThreads := 1
	if isBatch && core.Config.TxtDownloadThreads > 1 {
		numThreads = core.Config.TxtDownloadThreads
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, numThreads)
	totalTasks := len(finalUrls)

	logger.Section("开始下载任务")
	logger.TotalCount(totalTasks, numThreads)

	for i, urlToProcess := range finalUrls {
		wg.Add(1)
		semaphore <- struct{}{}
		go processURL(urlToProcess, &wg, semaphore, i+1, totalTasks)
	}

	wg.Wait()
}

func main() {
	core.InitFlags()

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: %s [选项] [url1 url2 ...]\n", os.Args[0])
		fmt.Println("如果没有提供URL，程序将进入交互模式。")
		fmt.Println("选项:")
		pflag.PrintDefaults()
	}

	pflag.Parse()

	err := core.LoadConfig(core.ConfigPath)
	if err != nil {
		if os.IsNotExist(err) && core.ConfigPath == "config.yaml" {
			logger.Error("默认配置文件 config.yaml 未找到。")
			pflag.Usage()
			return
		}
		logger.Error("加载配置文件 %s 失败: %v", core.ConfigPath, err)
		return
	}

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
		fmt.Print("请输入专辑链接或TXT文件路径: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			logger.Warn("未输入内容，程序退出。")
			return
		}

		if strings.HasSuffix(strings.ToLower(input), ".txt") {
			if _, err := os.Stat(input); err == nil {
				fileBytes, err := os.ReadFile(input)
				if err != nil {
					logger.Error("读取文件 %s 失败: %v", input, err)
					return
				}
				lines := strings.Split(string(fileBytes), "\n")
				var urls []string
				for _, line := range lines {
					trimmedLine := strings.TrimSpace(line)
					if trimmedLine != "" {
						urls = append(urls, trimmedLine)
					}
				}
				runDownloads(urls, true)
			} else {
				logger.Error("文件不存在 %s", input)
				return
			}
		} else {
			runDownloads([]string{input}, false)
		}
	} else {
		runDownloads(args, false)
	}

	logger.Plain("\n")
	logger.Success(fmt.Sprintf("已完成: %d/%d  |  警告: %d  |  错误: %d",
		core.Counter.Success, core.Counter.Total,
		core.Counter.Unavailable+core.Counter.NotSong, core.Counter.Error))
	if core.Counter.Error > 0 {
		logger.Warn("部分任务在执行过程中出错，请检查上面的日志记录。")
	}
}
