package main
import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

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

	err = downloader.MvDownloader(albumId, core.Config.AlacSaveFolder, sanitizedArtistFolder, "", storefront, nil, accountForMV)

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

func main() {
	core.InitFlags()

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] url1 url2 ...\n", "[main | main.exe | go run main.go]")
		fmt.Println("Options:")
		pflag.PrintDefaults()
	}

	pflag.Parse()

	err := core.LoadConfig(core.ConfigPath)
	if err != nil {
		if os.IsNotExist(err) && core.ConfigPath == "config.yaml" {
			fmt.Println("错误: 默认配置文件 config.yaml 未找到，请在程序同目录下创建它，或通过 --config 参数指定一个有效的配置文件。")
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
			fmt.Println("Failed to get developer token.")
			return
		}
	}
	core.DeveloperToken = token

	args := pflag.Args()
	if len(args) == 0 {
		fmt.Println("No URLs provided. Please provide at least one URL.")
		pflag.Usage()
		return
	}
	os.Args = args

	if strings.Contains(os.Args[0], "/artist/") {
		artistAccount := &core.Config.Accounts[0]
		urlArtistName, urlArtistID, err := api.GetUrlArtistName(os.Args[0], artistAccount)
		if err != nil {
			fmt.Println("Failed to get artistname.")
			return
		}
		core.Config.ArtistFolderFormat = strings.NewReplacer(
			"{UrlArtistName}", core.LimitString(urlArtistName),
			"{ArtistId}", urlArtistID,
		).Replace(core.Config.ArtistFolderFormat)
		albumArgs, err := api.CheckArtist(os.Args[0], artistAccount, "albums")
		if err != nil {
			fmt.Println("Failed to get artist albums.")
			return
		}
		mvArgs, err := api.CheckArtist(os.Args[0], artistAccount, "music-videos")
		if err != nil {
			fmt.Println("Failed to get artist music-videos.")
		}
		os.Args = append(albumArgs, mvArgs...)
	}

	albumTotal := len(os.Args)
	for {
		for albumNum, urlRaw := range os.Args {
			fmt.Printf("正在处理专辑 %d of %d: %s\n", albumNum+1, albumTotal, urlRaw)
			var storefront, albumId string

			if strings.Contains(urlRaw, "/music-video/") {
				handleSingleMV(urlRaw)
				continue
			}

			if strings.Contains(urlRaw, "/song/") {
				tempStorefront, _ := parser.CheckUrlSong(urlRaw)
				accountForSong, err := core.GetAccountForStorefront(tempStorefront)
				if err != nil {
					fmt.Printf("获取歌曲信息失败: %v\n", err)
					continue
				}
				urlRaw, err = api.GetUrlSong(urlRaw, accountForSong)
				if err != nil {
					fmt.Println("Failed to get Song info.")
					continue
				}
				core.Dl_song = true
			}

			if strings.Contains(urlRaw, "/playlist/") {
				storefront, albumId = parser.CheckUrlPlaylist(urlRaw)
			} else {
				storefront, albumId = parser.CheckUrl(urlRaw)
			}
			if albumId == "" {
				fmt.Printf("Invalid URL: %s\n", urlRaw)
				continue
			}
			parse, err := url.Parse(urlRaw)
			if err != nil {
				log.Fatalf("Invalid URL: %v", err)
			}
			var urlArg_i = parse.Query().Get("i")
			err = downloader.Rip(albumId, storefront, urlArg_i, urlRaw)
			if err != nil {
				fmt.Println("Album failed:", err)
			}
		}
		fmt.Printf("=======  [\u2714 ] Completed: %d/%d  |  [\u26A0 ] Warnings: %d  |  [\u2716 ] Errors: %d  =======\n", core.Counter.Success, core.Counter.Total, core.Counter.Unavailable+core.Counter.NotSong, core.Counter.Error)
		if core.Counter.Error == 0 {
			break
		}
		fmt.Println("错误已发生，正在自动重试...")
		fmt.Println("开始重试...")
		core.SharedLock.Lock()
		core.Counter = core.InitCounter()
		core.SharedLock.Unlock()
	}
}

