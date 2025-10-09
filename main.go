package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"main/internal/api"
	"main/internal/core"
	"main/internal/downloader"
	"main/internal/parser"

	"github.com/spf13/pflag"
)

func handleSingleMV(urlRaw string) {
	if core.Debug_mode {
		return
	}
	storefront, albumId := parser.CheckUrlMv(urlRaw)
	accountForMV, err := core.GetAccountForStorefront(storefront)
	if err != nil {
		fmt.Printf("MV 下载失败: %v\n", err)
		core.SharedLock.Lock()
		core.Counter.Error++
		core.SharedLock.Unlock()
		return
	}

	core.SharedLock.Lock()
	core.Counter.Total++
	core.SharedLock.Unlock()
	if len(accountForMV.MediaUserToken) <= 50 {
		core.SharedLock.Lock()
		core.Counter.Error++
		core.SharedLock.Unlock()
		return
	}
	if _, err := exec.LookPath("mp4decrypt"); err != nil {
		core.SharedLock.Lock()
		core.Counter.Error++
		core.SharedLock.Unlock()
		return
	}

	mvInfo, err := api.GetMVInfoFromAdam(albumId, accountForMV, storefront)
	if err != nil {
		fmt.Printf("获取 MV 信息失败: %v\n", err)
		core.SharedLock.Lock()
		core.Counter.Error++
		core.SharedLock.Unlock()
		return
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

	mvOutPath, err := downloader.MvDownloader(albumId, cachePath, sanitizedArtistFolder, "", storefront, nil, accountForMV)

	// 如果使用缓存且下载成功，移动文件到最终位置
	if err == nil && usingCache && mvOutPath != "" {
		// 计算最终路径
		relPath, _ := filepath.Rel(cachePath, mvOutPath)
		finalMvPath := filepath.Join(finalPath, relPath)

		// 移动文件
		if moveErr := downloader.SafeMoveFile(mvOutPath, finalMvPath); moveErr != nil {
			fmt.Printf("从缓存移动MV文件失败: %v\n", moveErr)
			err = moveErr
		} else {
			// 清理缓存目录
			mvCacheDir := filepath.Dir(mvOutPath)
			for mvCacheDir != cachePath && mvCacheDir != "." && mvCacheDir != "/" {
				if os.Remove(mvCacheDir) != nil {
					break
				}
				mvCacheDir = filepath.Dir(mvCacheDir)
			}
		}
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

func processURL(urlRaw string, wg *sync.WaitGroup, semaphore chan struct{}, currentTask int, totalTasks int) {
	if wg != nil {
		defer wg.Done()
	}
	if semaphore != nil {
		defer func() { <-semaphore }()
	}

	if totalTasks > 1 {
		core.SafePrintf("🧾 [%d/%d] 开始处理: %s\n", currentTask, totalTasks, urlRaw)
	}

	var storefront, albumId string

	if strings.Contains(urlRaw, "/music-video/") {
		handleSingleMV(urlRaw)
		return
	}

	if strings.Contains(urlRaw, "/song/") {
		tempStorefront, _ := parser.CheckUrlSong(urlRaw)
		accountForSong, err := core.GetAccountForStorefront(tempStorefront)
		if err != nil {
			fmt.Printf("获取歌曲信息失败 for %s: %v\n", urlRaw, err)
			return
		}
		urlRaw, err = api.GetUrlSong(urlRaw, accountForSong)
		if err != nil {
			fmt.Printf("获取歌曲链接失败 for %s: %v\n", urlRaw, err)
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
		fmt.Printf("无效的URL: %s\n", urlRaw)
		return
	}

	parse, err := url.Parse(urlRaw)
	if err != nil {
		log.Printf("解析URL失败 %s: %v", urlRaw, err)
		return
	}
	var urlArg_i = parse.Query().Get("i")
	err = downloader.Rip(albumId, storefront, urlArg_i, urlRaw)
	if err != nil {
		core.SafePrintf("专辑下载失败: %s -> %v\n", urlRaw, err)
	} else {
		if totalTasks > 1 {
			core.SafePrintf("✅ [%d/%d] 任务完成: %s\n", currentTask, totalTasks, urlRaw)
		}
	}
}

func runDownloads(initialUrls []string, isBatch bool) {
	var finalUrls []string

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
		fmt.Println("队列中没有有效的链接可供下载。")
		return
	}

	numThreads := 1
	if isBatch && core.Config.TxtDownloadThreads > 1 {
		numThreads = core.Config.TxtDownloadThreads
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, numThreads)
	totalTasks := len(finalUrls)

	core.SafePrintf("📋 开始下载任务\n📝 总数: %d, 并发数: %d\n--------------------\n", totalTasks, numThreads)

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
			fmt.Println("错误: 默认配置文件 config.yaml 未找到。")
			pflag.Usage()
			return
		}
		fmt.Printf("加载配置文件 %s 失败: %v\n", core.ConfigPath, err)
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
			fmt.Println("获取开发者 token 失败。")
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
			fmt.Println("未输入内容，程序退出。")
			return
		}

		if strings.HasSuffix(strings.ToLower(input), ".txt") {
			if _, err := os.Stat(input); err == nil {
				fileBytes, err := os.ReadFile(input)
				if err != nil {
					fmt.Printf("读取文件 %s 失败: %v\n", input, err)
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
				fmt.Printf("错误: 文件不存在 %s\n", input)
				return
			}
		} else {
			runDownloads([]string{input}, false)
		}
	} else {
		runDownloads(args, false)
	}

	fmt.Printf("\n📦 已完成: %d/%d | 警告: %d | 错误: %d\n", core.Counter.Success, core.Counter.Total, core.Counter.Unavailable+core.Counter.NotSong, core.Counter.Error)
	if core.Counter.Error > 0 {
		fmt.Println("部分任务在执行过程中出错，请检查上面的日志记录。")
	}
}
