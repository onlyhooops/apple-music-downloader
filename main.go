package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"main/utils/lyrics"
	"main/utils/runv10"
	"main/utils/runv3"
	"main/utils/structs"

	"github.com/fatih/color"
	"github.com/grafov/m3u8"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/pflag"
	"github.com/zhaarey/go-mp4tag"
	"gopkg.in/yaml.v2"
)

var (
	forbiddenNames = regexp.MustCompile(`[/\\<>:"|?*]`)
	dl_atmos       bool
	dl_aac         bool
	dl_select      bool
	dl_song        bool
	artist_select  bool
	debug_mode     bool
	alac_max       *int
	atmos_max      *int
	mv_max         *int
	mv_audio_type  *string
	aac_type       *string
	Config         structs.ConfigSet
	counter        structs.Counter
	okDict         = make(map[string][]int)
	configPath     string
	outputPath     string
	sharedLock     sync.Mutex
	developerToken string
	maxPathLength  int
)

type TrackStatus struct {
	Index       int
	TrackNum    int
	TrackTotal  int
	TrackName   string
	Quality     string
	Status      string
	StatusColor func(a ...interface{}) string
}

var uiMutex sync.Mutex
var trackStatuses []TrackStatus

func ensureSafePath(basePath, artistDir, albumDir, fileName string) (string, string, string) {
	truncate := func(s string, n int) string {
		if n <= 0 {
			return s
		}
		runes := []rune(s)
		if len(runes) <= n {
			return ""
		}
		return string(runes[:len(runes)-n])
	}

	for {
		currentPath := filepath.Join(basePath, artistDir, albumDir, fileName)
		if len(currentPath) <= maxPathLength {
			break
		}

		overage := len(currentPath) - maxPathLength
		ext := filepath.Ext(fileName)
		stem := strings.TrimSuffix(fileName, ext)

		var prefixPart string
		var namePart string
		re := regexp.MustCompile(`^(\d+[\s.-]*)`)
		matches := re.FindStringSubmatch(stem)

		if len(matches) > 1 {
			prefixPart = matches[1]
			namePart = strings.TrimPrefix(stem, prefixPart)
		} else {
			prefixPart = ""
			namePart = stem
		}

		if len(namePart) > 0 {
			canShorten := len(namePart)
			shortenAmount := overage
			if shortenAmount > canShorten {
				shortenAmount = canShorten
			}
			namePart = truncate(namePart, shortenAmount)

			if namePart == "" {
				prefixPart = strings.TrimRight(prefixPart, " .-")
			}

			fileName = prefixPart + namePart + ext
			continue
		}

		if len(albumDir) > 1 {
			canShorten := len(albumDir)
			shortenAmount := overage
			if shortenAmount > canShorten {
				shortenAmount = canShorten
			}
			albumDir = truncate(albumDir, shortenAmount)
			continue
		}

		if len(artistDir) > 1 {
			canShorten := len(artistDir)
			shortenAmount := overage
			if shortenAmount > canShorten {
				shortenAmount = canShorten
			}
			artistDir = truncate(artistDir, shortenAmount)
			continue
		}
		break
	}

	return artistDir, albumDir, fileName
}

func loadConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		return err
	}

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	if len(Config.Accounts) == 0 {
		return errors.New(red("配置错误: 'accounts' 列表为空，请在 config.yaml 中至少配置一个账户"))
	}

	useAutoDetect := true
	if Config.MaxPathLength > 0 {
		maxPathLength = Config.MaxPathLength
		useAutoDetect = false
		fmt.Printf("%s%s\n",
			green("从配置文件强制使用专辑文件夹路径最大字符长度: "),
			red(fmt.Sprintf("%d", maxPathLength)),
		)
	}

	if useAutoDetect {
		if runtime.GOOS == "windows" {
			maxPathLength = 255
			fmt.Printf("%s%d\n",
				green("检测到 Windows 系统, 已自动设置最大路径长度限制为: "),
				maxPathLength,
			)
		} else {
			maxPathLength = 4096
			fmt.Printf("%s%s%s%d\n",
				green("检测到 "),
				red(runtime.GOOS),
				green(" 系统, 已自动设置最大路径长度限制为: "),
				maxPathLength,
			)
		}
	}
	return nil
}

func getAccountForStorefront(storefront string) (*structs.Account, error) {
	if len(Config.Accounts) == 0 {
		return nil, errors.New("无可用账户")
	}

	for i := range Config.Accounts {
		acc := &Config.Accounts[i]
		if strings.ToLower(acc.Storefront) == strings.ToLower(storefront) {
			return acc, nil
		}
	}

	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf(
		"%s 未找到与 %s 匹配的账户,将尝试使用 %s 等区域进行下载\n",
		red("警告:"),
		red(storefront),
		yellow(Config.Accounts[0].Name),
	)
	return &Config.Accounts[0], nil
}

func LimitString(s string) string {
	if len([]rune(s)) > Config.LimitMax {
		return string([]rune(s)[:Config.LimitMax])
	}
	return s
}

func isInArray(arr []int, target int) bool {
	for _, num := range arr {
		if num == target {
			return true
		}
	}
	return false
}

func fileExists(path string) (bool, error) {
	f, err := os.Stat(path)
	if err == nil {
		return !f.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func formatSpeed(bytesPerSecond float64) string {
	if bytesPerSecond < 1024 {
		return fmt.Sprintf("%.1f B/s", bytesPerSecond)
	}
	kbps := bytesPerSecond / 1024
	if kbps < 1024 {
		return fmt.Sprintf("%.1f KB/s", kbps)
	}
	mbps := kbps / 1024
	return fmt.Sprintf("%.1f MB/s", mbps)
}

func checkUrl(url string) (string, string) {
	pat := regexp.MustCompile(`^(?:https:\/\/(?:beta\.music|music|classical\.music)\.apple\.com\/(\w{2})(?:\/album|\/album\/.+))\/(?:id)?(\d[^\D]+)(?:$|\?)`)
	matches := pat.FindAllStringSubmatch(url, -1)

	if matches == nil {
		return "", ""
	} else {
		return matches[0][1], matches[0][2]
	}
}
func checkUrlMv(url string) (string, string) {
	pat := regexp.MustCompile(`^(?:https:\/\/(?:beta\.music|music)\.apple\.com\/(\w{2})(?:\/music-video|\/music-video\/.+))\/(?:id)?(\d[^\D]+)(?:$|\?)`)
	matches := pat.FindAllStringSubmatch(url, -1)

	if matches == nil {
		return "", ""
	} else {
		return matches[0][1], matches[0][2]
	}
}
func checkUrlSong(url string) (string, string) {
	pat := regexp.MustCompile(`^(?:https:\/\/(?:beta\.music|music)\.apple\.com\/(\w{2})(?:\/song|\/song\/.+))\/(?:id)?(\d[^\D]+)(?:$|\?)`)
	matches := pat.FindAllStringSubmatch(url, -1)

	if matches == nil {
		return "", ""
	} else {
		return matches[0][1], matches[0][2]
	}
}
func checkUrlPlaylist(url string) (string, string) {
	pat := regexp.MustCompile(`^(?:https:\/\/(?:beta\.music|music)\.apple\.com\/(\w{2})(?:\/playlist|\/playlist\/.+))\/(?:id)?(pl\.[\w-]+)(?:$|\?)`)
	matches := pat.FindAllStringSubmatch(url, -1)

	if matches == nil {
		return "", ""
	} else {
		return matches[0][1], matches[0][2]
	}
}

func checkUrlArtist(url string) (string, string) {
	pat := regexp.MustCompile(`^(?:https:\/\/(?:beta\.music|music)\.apple\.com\/(\w{2})(?:\/artist|\/artist\/.+))\/(?:id)?(\d[^\D]+)(?:$|\?)`)
	matches := pat.FindAllStringSubmatch(url, -1)

	if matches == nil {
		return "", ""
	} else {
		return matches[0][1], matches[0][2]
	}
}
func getUrlSong(songUrl string, account *structs.Account) (string, error) {
	storefront, songId := checkUrlSong(songUrl)
	manifest, err := getInfoFromAdam(songId, account, storefront)
	if err != nil {
		fmt.Println("\u26A0 Failed to get manifest:", err)
		sharedLock.Lock()
		counter.NotSong++
		sharedLock.Unlock()
		return "", err
	}
	if len(manifest.Relationships.Albums.Data) == 0 {
		return "", errors.New("song does not appear to be part of an album")
	}
	albumId := manifest.Relationships.Albums.Data[0].ID
	songAlbumUrl := fmt.Sprintf("https://music.apple.com/%s/album/1/%s?i=%s", storefront, albumId, songId)
	return songAlbumUrl, nil
}
func getUrlArtistName(artistUrl string, account *structs.Account) (string, string, error) {
	storefront, artistId := checkUrlArtist(artistUrl)
	req, err := http.NewRequest("GET", fmt.Sprintf("https://amp-api.music.apple.com/v1/catalog/%s/artists/%s", storefront, artistId), nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", developerToken))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Origin", "https://music.apple.com")
	query := url.Values{}
	query.Set("l", Config.Language)
	req.URL.RawQuery = query.Encode()
	do, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return "", "", errors.New(do.Status)
	}
	obj := new(structs.AutoGeneratedArtist)
	err = json.NewDecoder(do.Body).Decode(&obj)
	if err != nil {
		return "", "", err
	}
	return obj.Data[0].Attributes.Name, obj.Data[0].ID, nil
}

func checkArtist(artistUrl string, account *structs.Account, relationship string) ([]string, error) {
	storefront, artistId := checkUrlArtist(artistUrl)
	Num := 0
	var args []string
	var urls []string
	var options [][]string
	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("https://amp-api.music.apple.com/v1/catalog/%s/artists/%s/%s?limit=100&offset=%d&l=%s", storefront, artistId, relationship, Num, Config.Language), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", developerToken))
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Origin", "https://music.apple.com")
		do, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer do.Body.Close()
		if do.StatusCode != http.StatusOK {
			return nil, errors.New(do.Status)
		}
		obj := new(structs.AutoGeneratedArtist)
		err = json.NewDecoder(do.Body).Decode(&obj)
		if err != nil {
			return nil, err
		}
		for _, album := range obj.Data {
			options = append(options, []string{album.Attributes.Name, album.Attributes.ReleaseDate, album.ID, album.Attributes.URL})
		}
		Num = Num + 100
		if len(obj.Next) == 0 {
			break
		}
	}
	sort.Slice(options, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", options[i][1])
		dateJ, _ := time.Parse("2006-01-02", options[j][1])
		return dateI.Before(dateJ)
	})

	table := tablewriter.NewWriter(os.Stdout)
	if relationship == "albums" {
		table.SetHeader([]string{"", "Album Name", "Date", "Album ID"})
	} else if relationship == "music-videos" {
		table.SetHeader([]string{"", "MV Name", "Date", "MV ID"})
	}
	table.SetRowLine(false)
	table.SetHeaderColor(tablewriter.Colors{},
		tablewriter.Colors{tablewriter.FgRedColor, tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor})

	table.SetColumnColor(tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgRedColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor})
	for i, v := range options {
		urls = append(urls, v[3])
		options[i] = append([]string{fmt.Sprint(i + 1)}, v[:3]...)
		table.Append(options[i])
	}
	table.Render()
	if artist_select {
		fmt.Println("You have selected all options:")
		return urls, nil
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please select from the " + relationship + " options above (multiple options separated by commas, ranges supported, or type 'all' to select all)")
	cyanColor := color.New(color.FgCyan)
	cyanColor.Print("Enter your choice: ")
	input, _ := reader.ReadString('\n')

	input = strings.TrimSpace(input)
	if input == "all" {
		fmt.Println("You have selected all options:")
		return urls, nil
	}

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

	fmt.Println("You have selected the following options:")
	for _, opt := range selectedOptions {
		if len(opt) == 1 {
			num, err := strconv.Atoi(opt[0])
			if err != nil {
				fmt.Println("Invalid option:", opt[0])
				continue
			}
			if num > 0 && num <= len(options) {
				fmt.Println(options[num-1])
				args = append(args, urls[num-1])
			} else {
				fmt.Println("Option out of range:", opt[0])
			}
		} else if len(opt) == 2 {
			start, err1 := strconv.Atoi(opt[0])
			end, err2 := strconv.Atoi(opt[1])
			if err1 != nil || err2 != nil {
				fmt.Println("Invalid range:", opt)
				continue
			}
			if start < 1 || end > len(options) || start > end {
				fmt.Println("Range out of range:", opt)
				continue
			}
			for i := start; i <= end; i++ {
				fmt.Println(options[i-1])
				args = append(args, urls[i-1])
			}
		} else {
			fmt.Println("Invalid option:", opt)
		}
	}
	return args, nil
}

func getMeta(albumId string, account *structs.Account, storefront string) (*structs.AutoGenerated, error) {
	var mtype string
	var next string
	if strings.Contains(albumId, "pl.") {
		mtype = "playlists"
	} else {
		mtype = "albums"
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://amp-api.music.apple.com/v1/catalog/%s/%s/%s", storefront, mtype, albumId), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", developerToken))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Origin", "https://music.apple.com")
	query := url.Values{}
	query.Set("omit[resource]", "autos")
	query.Set("include", "tracks,artists,record-labels")
	query.Set("include[songs]", "artists,albums")
	query.Set("fields[artists]", "name,artwork")
	query.Set("fields[albums:albums]", "artistName,artwork,name,releaseDate,url")
	query.Set("fields[record-labels]", "name")
	query.Set("extend", "editorialVideo")
	query.Set("l", Config.Language)
	req.URL.RawQuery = query.Encode()
	do, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return nil, errors.New(do.Status)
	}
	obj := new(structs.AutoGenerated)
	err = json.NewDecoder(do.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}
	if strings.Contains(albumId, "pl.") {
		obj.Data[0].Attributes.ArtistName = "Apple Music"
	}
	if len(obj.Data[0].Relationships.Tracks.Next) > 0 {
		next = obj.Data[0].Relationships.Tracks.Next
		for {
			req, err := http.NewRequest("GET", fmt.Sprintf("https://amp-api.music.apple.com/%s&l=%s&include=albums", next, Config.Language), nil)
			if err != nil {
				return nil, err
			}
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", developerToken))
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
			req.Header.Set("Origin", "https://music.apple.com")
			do, err := http.DefaultClient.Do(req)
			if err != nil {
				return nil, err
			}
			defer do.Body.Close()
			if do.StatusCode != http.StatusOK {
				return nil, errors.New(do.Status)
			}
			obj2 := new(structs.AutoGeneratedTrack)
			err = json.NewDecoder(do.Body).Decode(&obj2)
			if err != nil {
				return nil, err
			}
			for _, value := range obj2.Data {
				obj.Data[0].Relationships.Tracks.Data = append(obj.Data[0].Relationships.Tracks.Data, value)
			}
			next = obj2.Next
			if len(next) == 0 {
				break
			}
		}
	}
	return obj, nil
}

func writeCover(sanAlbumFolder, name string, url string) (string, error) {
	covPath := filepath.Join(sanAlbumFolder, name+"."+Config.CoverFormat)
	if Config.CoverFormat == "original" {
		ext := strings.Split(url, "/")[len(strings.Split(url, "/"))-2]
		ext = ext[strings.LastIndex(ext, ".")+1:]
		covPath = filepath.Join(sanAlbumFolder, name+"."+ext)
	}
	exists, err := fileExists(covPath)
	if err != nil {
		return "", err
	}
	if exists {
		_ = os.Remove(covPath)
	}
	if Config.CoverFormat == "png" {
		re := regexp.MustCompile(`\{w\}x\{h\}`)
		parts := re.Split(url, 2)
		url = parts[0] + "{w}x{h}" + strings.Replace(parts[1], ".jpg", ".png", 1)
	}
	url = strings.Replace(url, "{w}x{h}", Config.CoverSize, 1)
	if Config.CoverFormat == "original" {
		url = strings.Replace(url, "is1-ssl.mzstatic.com/image/thumb", "a5.mzstatic.com/us/r1000/0", 1)
		url = url[:strings.LastIndex(url, "/")]
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	do, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return "", errors.New(do.Status)
	}
	f, err := os.Create(covPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, do.Body)
	if err != nil {
		return "", err
	}
	return covPath, nil
}

func writeLyrics(sanAlbumFolder, filename string, lrc string) error {
	lyricspath := filepath.Join(sanAlbumFolder, filename)
	f, err := os.Create(lyricspath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(lrc)
	if err != nil {
		return err
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func downloadTrackWithFallback(track structs.TrackData, meta *structs.AutoGenerated, albumId, storefront, baseSaveFolder, Codec, covPath string, workingAccounts []structs.Account, initialAccountIndex int, statusIndex int, updateStatus func(index int, status string, sColor func(a ...interface{}) string), progressChan chan runv10.ProgressUpdate) error {
	maxRetries := 2
	var lastError error

	for i := 0; i < len(workingAccounts); i++ {
		accountIndex := (initialAccountIndex + i) % len(workingAccounts)
		account := &workingAccounts[accountIndex]

		for attempt := 0; attempt <= maxRetries; attempt++ {
			err := downloadTrackSilently(track, meta, albumId, storefront, baseSaveFolder, Codec, covPath, account, progressChan)
			if err == nil {
				return nil
			}
			lastError = err
			if attempt < maxRetries {
				time.Sleep(2 * time.Second)
			}
		}
		warningMsg := fmt.Sprintf("账户 %s 失败, 尝试下一个...", account.Name)
		updateStatus(statusIndex, warningMsg, color.New(color.FgRed).SprintFunc())
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("所有可用账户均尝试失败: %w", lastError)
}

func downloadTrackSilently(track structs.TrackData, meta *structs.AutoGenerated, albumId, storefront, baseSaveFolder, Codec, covPath string, account *structs.Account, progressChan chan runv10.ProgressUpdate) error {
	if track.Type == "music-videos" {
		if len(account.MediaUserToken) <= 50 {
			return errors.New("media-user-token is not set, skip MV dl")
		}
		if _, err := exec.LookPath("mp4decrypt"); err != nil {
			return errors.New("mp4decrypt is not found, skip MV dl")
		}
		err := mvDownloader(track.ID, baseSaveFolder, storefront, meta, account)
		if err != nil {
			return fmt.Errorf("failed to dl MV: %w", err)
		}
		return nil
	}

	manifest, err := getInfoFromAdam(track.ID, account, storefront)
	if err != nil {
		return fmt.Errorf("failed to get manifest with account %s: %w", account.Name, err)
	}

	needDlAacLc := false
	if dl_aac && Config.AacType == "aac-lc" {
		needDlAacLc = true
	}
	if manifest.Attributes.ExtendedAssetUrls.EnhancedHls == "" {
		if dl_atmos {
			return errors.New("atmos unavailable")
		}
		needDlAacLc = true
	}
	needCheck := false

	if Config.GetM3u8Mode == "all" {
		needCheck = true
	} else if Config.GetM3u8Mode == "hires" && contains(track.Attributes.AudioTraits, "hi-res-lossless") {
		needCheck = true
	}
	var EnhancedHls_m3u8 string
	if needCheck && !needDlAacLc {
		EnhancedHls_m3u8, _ = checkM3u8(track.ID, "song", account)
		if strings.HasSuffix(EnhancedHls_m3u8, ".m3u8") {
			manifest.Attributes.ExtendedAssetUrls.EnhancedHls = EnhancedHls_m3u8
		}
	}
	var Quality string
	if strings.Contains(Config.SongFileFormat, "Quality") {
		if dl_atmos {
			Quality = fmt.Sprintf("%dkbps", Config.AtmosMax-2000)
		} else if needDlAacLc {
			Quality = "256kbps"
		} else {
			_, Quality, _, err = extractMedia(manifest.Attributes.ExtendedAssetUrls.EnhancedHls, true)
			if err != nil {
				return fmt.Errorf("failed to extract quality from manifest: %w", err)
			}
		}
	}
	stringsToJoin := []string{}
	if track.Attributes.IsAppleDigitalMaster {
		if Config.AppleMasterChoice != "" {
			stringsToJoin = append(stringsToJoin, Config.AppleMasterChoice)
		}
	}
	if track.Attributes.ContentRating == "explicit" {
		if Config.ExplicitChoice != "" {
			stringsToJoin = append(stringsToJoin, Config.ExplicitChoice)
		}
	}
	if track.Attributes.ContentRating == "clean" {
		if Config.CleanChoice != "" {
			stringsToJoin = append(stringsToJoin, Config.CleanChoice)
		}
	}
	Tag_string := strings.Join(stringsToJoin, " ")

	trackNum := -1
	trackTotal := len(meta.Data[0].Relationships.Tracks.Data)
	for i, t := range meta.Data[0].Relationships.Tracks.Data {
		if t.ID == track.ID {
			trackNum = i + 1
			break
		}
	}
	if trackNum == -1 {
		return errors.New("track not found in metadata")
	}
	var singerFoldername, albumFoldername string
	if Config.ArtistFolderFormat != "" {
		if strings.Contains(albumId, "pl.") {
			singerFoldername = strings.NewReplacer(
				"{ArtistName}", "Apple Music", "{ArtistId}", "", "{UrlArtistName}", "Apple Music",
			).Replace(Config.ArtistFolderFormat)
		} else if len(meta.Data[0].Relationships.Artists.Data) > 0 {
			singerFoldername = strings.NewReplacer(
				"{UrlArtistName}", LimitString(meta.Data[0].Attributes.ArtistName),
				"{ArtistName}", LimitString(meta.Data[0].Attributes.ArtistName),
				"{ArtistId}", meta.Data[0].Relationships.Artists.Data[0].ID,
			).Replace(Config.ArtistFolderFormat)
		} else {
			singerFoldername = strings.NewReplacer(
				"{UrlArtistName}", LimitString(meta.Data[0].Attributes.ArtistName),
				"{ArtistName}", LimitString(meta.Data[0].Attributes.ArtistName),
				"{ArtistId}", "",
			).Replace(Config.ArtistFolderFormat)
		}
	}

	if strings.Contains(albumId, "pl.") {
		albumFoldername = strings.NewReplacer(
			"{PlaylistName}", LimitString(meta.Data[0].Attributes.Name),
			"{PlaylistId}", albumId, "{Quality}", Quality, "{Codec}", Codec, "{Tag}", Tag_string,
		).Replace(Config.PlaylistFolderFormat)
	} else {
		albumFoldername = strings.NewReplacer(
			"{ReleaseDate}", meta.Data[0].Attributes.ReleaseDate, "{ReleaseYear}", meta.Data[0].Attributes.ReleaseDate[:4],
			"{ArtistName}", LimitString(meta.Data[0].Attributes.ArtistName), "{AlbumName}", LimitString(meta.Data[0].Attributes.Name),
			"{UPC}", meta.Data[0].Attributes.Upc, "{RecordLabel}", meta.Data[0].Attributes.RecordLabel,
			"{Copyright}", meta.Data[0].Attributes.Copyright, "{AlbumId}", albumId,
			"{Quality}", Quality, "{Codec}", Codec, "{Tag}", Tag_string,
		).Replace(Config.AlbumFolderFormat)
	}

	songName := strings.NewReplacer(
		"{SongId}", track.ID,
		"{SongNumer}", fmt.Sprintf("%02d", trackNum),
		"{SongName}", LimitString(track.Attributes.Name),
		"{DiscNumber}", fmt.Sprintf("%0d", track.Attributes.DiscNumber),
		"{TrackNumber}", fmt.Sprintf("%0d", track.Attributes.TrackNumber),
		"{Quality}", Quality,
		"{Tag}", Tag_string,
		"{Codec}", Codec,
	).Replace(Config.SongFileFormat)

	sanitizedSingerFolder := forbiddenNames.ReplaceAllString(singerFoldername, "_")
	sanitizedAlbumFolder := forbiddenNames.ReplaceAllString(albumFoldername, "_")
	sanitizedSongName := forbiddenNames.ReplaceAllString(songName, "_")
	filenameWithExt := fmt.Sprintf("%s.m4a", sanitizedSongName)

	finalArtistDir, finalAlbumDir, finalFilename := ensureSafePath(baseSaveFolder, sanitizedSingerFolder, sanitizedAlbumFolder, filenameWithExt)

	var finalSingerFolder string
	if finalArtistDir != "" {
		finalSingerFolder = filepath.Join(baseSaveFolder, finalArtistDir)
	} else {
		finalSingerFolder = baseSaveFolder
	}
	finalAlbumFolder := filepath.Join(finalSingerFolder, finalAlbumDir)
	os.MkdirAll(finalAlbumFolder, os.ModePerm)
	trackPath := filepath.Join(finalAlbumFolder, finalFilename)

	lrcFilename := fmt.Sprintf("%s.%s", strings.TrimSuffix(finalFilename, ".m4a"), Config.LrcFormat)

	var lrc string = ""
	if Config.EmbedLrc || Config.SaveLrcFile {
		lrcStr, err := lyrics.Get(storefront, track.ID, Config.LrcType, Config.Language, Config.LrcFormat, developerToken, account.MediaUserToken)
		if err == nil {
			if Config.SaveLrcFile {
				_, _, finalLrcFilename := ensureSafePath(baseSaveFolder, finalArtistDir, finalAlbumDir, lrcFilename)
				_ = writeLyrics(finalAlbumFolder, finalLrcFilename, lrcStr)
			}
			if Config.EmbedLrc {
				lrc = lrcStr
			}
		}
	}

	exists, err := fileExists(trackPath)
	if err != nil {
		return errors.New("failed to check if track exists")
	}
	if exists {
		okDict[albumId] = append(okDict[albumId], trackNum)
		return nil
	}
	if needDlAacLc {
		if len(account.MediaUserToken) <= 50 {
			return errors.New("invalid media-user-token")
		}
		_, err := runv3.Run(track.ID, trackPath, developerToken, account.MediaUserToken, false)
		if err != nil {
			return fmt.Errorf("failed to dl aac-lc: %w", err)
		}
	} else {
		trackM3u8Url, _, _, err := extractMedia(manifest.Attributes.ExtendedAssetUrls.EnhancedHls, false)
		if err != nil {
			return fmt.Errorf("failed to extract info from manifest: %w", err)
		}
		err = runv10.Run(track.ID, trackM3u8Url, trackPath, account, Config, progressChan)
		if err != nil {
			return fmt.Errorf("failed to run v10 with account %s: %w", account.Name, err)
		}
	}
	tags := []string{
		"tool=",
		fmt.Sprintf("artist=%s", meta.Data[0].Attributes.ArtistName),
	}
	var trackCovPath string
	if Config.EmbedCover {
		if strings.Contains(albumId, "pl.") && Config.DlAlbumcoverForPlaylist {
			_, _, safeCoverFilename := ensureSafePath(baseSaveFolder, finalArtistDir, finalAlbumDir, track.ID+".jpg")
			trackCovPath, err = writeCover(finalAlbumFolder, strings.TrimSuffix(safeCoverFilename, ".jpg"), track.Attributes.Artwork.URL)
			if err != nil {
			}
			tags = append(tags, fmt.Sprintf("cover=%s", trackCovPath))
		} else {
			tags = append(tags, fmt.Sprintf("cover=%s", covPath))
		}
	}
	tagsString := strings.Join(tags, ":")
	cmd := exec.Command("MP4Box", "-itags", tagsString, trackPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("embed failed: %w", err)
	}
	if strings.Contains(albumId, "pl.") && Config.DlAlbumcoverForPlaylist && trackCovPath != "" {
		if err := os.Remove(trackCovPath); err != nil {
		}
	}
	err = writeMP4Tags(trackPath, lrc, meta, trackNum, trackTotal)
	if err != nil {
		return fmt.Errorf("failed to write tags in media: %w", err)
	}

	okDict[albumId] = append(okDict[albumId], trackNum)
	return nil
}

func rip(albumId string, storefront string, urlArg_i string, urlRaw string) error {
	mainAccount, err := getAccountForStorefront(storefront)
	if err != nil {
		return err
	}

	meta, err := getMeta(albumId, mainAccount, storefront)
	if err != nil {
		return err
	}

	if debug_mode {
		fmt.Println(meta.Data[0].Attributes.ArtistName)
		fmt.Println(meta.Data[0].Attributes.Name)

		for trackNum, track := range meta.Data[0].Relationships.Tracks.Data {
			trackNum++
			fmt.Printf("\nTrack %d of %d:\n", trackNum, len(meta.Data[0].Relationships.Tracks.Data))
			fmt.Printf("%02d. %s\n", trackNum, track.Attributes.Name)

			manifest, err := getInfoFromAdam(track.ID, mainAccount, storefront)
			if err != nil {
				fmt.Printf("Failed to get manifest for track %d: %v\n", trackNum, err)
				continue
			}

			var m3u8Url string
			if manifest.Attributes.ExtendedAssetUrls.EnhancedHls != "" {
				m3u8Url = manifest.Attributes.ExtendedAssetUrls.EnhancedHls
			}
			needCheck := false
			if Config.GetM3u8Mode == "all" {
				needCheck = true
			} else if Config.GetM3u8Mode == "hires" && contains(track.Attributes.AudioTraits, "hi-res-lossless") {
				needCheck = true
			}
			if needCheck {
				fullM3u8Url, err := checkM3u8(track.ID, "song", mainAccount)
				if err == nil && strings.HasSuffix(fullM3u8Url, ".m3u8") {
					m3u8Url = fullM3u8Url
				}
			}

			_, _, _, err = extractMedia(m3u8Url, true)
			if err != nil {
				fmt.Printf("Failed to extract quality info for track %d: %v\n", trackNum, err)
				continue
			}
		}
		return nil
	}

	var Codec string
	if dl_atmos {
		Codec = "ATMOS"
	} else if dl_aac {
		Codec = "AAC"
	} else {
		Codec = "ALAC"
	}

	var baseSaveFolder string
	if dl_atmos {
		baseSaveFolder = Config.AtmosSaveFolder
	} else {
		baseSaveFolder = Config.AlacSaveFolder
	}

	var singerFoldername, albumFoldername string
	if Config.ArtistFolderFormat != "" {
		if strings.Contains(albumId, "pl.") {
			singerFoldername = strings.NewReplacer(
				"{ArtistName}", "Apple Music", "{ArtistId}", "", "{UrlArtistName}", "Apple Music",
			).Replace(Config.ArtistFolderFormat)
		} else if len(meta.Data[0].Relationships.Artists.Data) > 0 {
			singerFoldername = strings.NewReplacer(
				"{UrlArtistName}", LimitString(meta.Data[0].Attributes.ArtistName),
				"{ArtistName}", LimitString(meta.Data[0].Attributes.ArtistName),
				"{ArtistId}", meta.Data[0].Relationships.Artists.Data[0].ID,
			).Replace(Config.ArtistFolderFormat)
		} else {
			singerFoldername = strings.NewReplacer(
				"{UrlArtistName}", LimitString(meta.Data[0].Attributes.ArtistName),
				"{ArtistName}", LimitString(meta.Data[0].Attributes.ArtistName),
				"{ArtistId}", "",
			).Replace(Config.ArtistFolderFormat)
		}
	}

	var Quality string
	if strings.Contains(Config.AlbumFolderFormat, "Quality") {
		if dl_atmos {
			Quality = fmt.Sprintf("%dkbps", Config.AtmosMax-2000)
		} else if dl_aac && Config.AacType == "aac-lc" {
			Quality = "256kbps"
		} else {
			manifest1, err := getInfoFromAdam(meta.Data[0].Relationships.Tracks.Data[0].ID, mainAccount, storefront)
			if err != nil {
			} else {
				if manifest1.Attributes.ExtendedAssetUrls.EnhancedHls == "" {
					Codec = "AAC"
					Quality = "256kbps"
				} else {
					needCheck := false
					if Config.GetM3u8Mode == "all" {
						needCheck = true
					} else if Config.GetM3u8Mode == "hires" && contains(meta.Data[0].Relationships.Tracks.Data[0].Attributes.AudioTraits, "hi-res-lossless") {
						needCheck = true
					}
					var EnhancedHls_m3u8 string
					if needCheck {
						EnhancedHls_m3u8, _ = checkM3u8(meta.Data[0].Relationships.Tracks.Data[0].ID, "album", mainAccount)
						if strings.HasSuffix(EnhancedHls_m3u8, ".m3u8") {
							manifest1.Attributes.ExtendedAssetUrls.EnhancedHls = EnhancedHls_m3u8
						}
					}
					_, Quality, _, err = extractMedia(manifest1.Attributes.ExtendedAssetUrls.EnhancedHls, true)
					if err != nil {
					}
				}
			}
		}
	}
	stringsToJoin := []string{}
	if meta.Data[0].Attributes.IsAppleDigitalMaster || meta.Data[0].Attributes.IsMasteredForItunes {
		if Config.AppleMasterChoice != "" {
			stringsToJoin = append(stringsToJoin, Config.AppleMasterChoice)
		}
	}
	if meta.Data[0].Attributes.ContentRating == "explicit" {
		if Config.ExplicitChoice != "" {
			stringsToJoin = append(stringsToJoin, Config.ExplicitChoice)
		}
	}
	if meta.Data[0].Attributes.ContentRating == "clean" {
		if Config.CleanChoice != "" {
			stringsToJoin = append(stringsToJoin, Config.CleanChoice)
		}
	}
	Tag_string := strings.Join(stringsToJoin, " ")

	if strings.Contains(albumId, "pl.") {
		albumFoldername = strings.NewReplacer(
			"{PlaylistName}", LimitString(meta.Data[0].Attributes.Name),
			"{PlaylistId}", albumId, "{Quality}", Quality, "{Codec}", Codec, "{Tag}", Tag_string,
		).Replace(Config.PlaylistFolderFormat)
	} else {
		albumFoldername = strings.NewReplacer(
			"{ReleaseDate}", meta.Data[0].Attributes.ReleaseDate, "{ReleaseYear}", meta.Data[0].Attributes.ReleaseDate[:4],
			"{ArtistName}", LimitString(meta.Data[0].Attributes.ArtistName), "{AlbumName}", LimitString(meta.Data[0].Attributes.Name),
			"{UPC}", meta.Data[0].Attributes.Upc, "{RecordLabel}", meta.Data[0].Attributes.RecordLabel,
			"{Copyright}", meta.Data[0].Attributes.Copyright, "{AlbumId}", albumId,
			"{Quality}", Quality, "{Codec}", Codec, "{Tag}", Tag_string,
		).Replace(Config.AlbumFolderFormat)
	}

	sanitizedSingerFolder := forbiddenNames.ReplaceAllString(singerFoldername, "_")
	sanitizedAlbumFolder := forbiddenNames.ReplaceAllString(albumFoldername, "_")

	var longestFilename string
	for i := range meta.Data[0].Relationships.Tracks.Data {
		if len(meta.Data[0].Relationships.Tracks.Data[i].Attributes.Name) > len(longestFilename) {
			longestFilename = meta.Data[0].Relationships.Tracks.Data[i].Attributes.Name
		}
	}
	longestFilename = strings.NewReplacer(
		"{SongName}", longestFilename,
		"{SongNumer}", "99",
		"{Quality}", "24B-192.0kHz",
		"{Tag}", Config.AppleMasterChoice+" "+Config.ExplicitChoice,
		"{Codec}", "ATMOS",
	).Replace(Config.SongFileFormat) + ".m4a"

	finalArtistDir, finalAlbumDir, _ := ensureSafePath(baseSaveFolder, sanitizedSingerFolder, sanitizedAlbumFolder, longestFilename)

	var finalSingerFolder string
	if finalArtistDir != "" {
		finalSingerFolder = filepath.Join(baseSaveFolder, finalArtistDir)
	} else {
		finalSingerFolder = baseSaveFolder
	}
	finalAlbumFolder := filepath.Join(finalSingerFolder, finalAlbumDir)
	os.MkdirAll(finalAlbumFolder, os.ModePerm)

	fmt.Printf("歌手: %s\n", meta.Data[0].Attributes.ArtistName)
	fmt.Printf("专辑: %s\n", meta.Data[0].Attributes.Name)

	if Config.SaveArtistCover && !(strings.Contains(albumId, "pl.")) {
		if len(meta.Data[0].Relationships.Artists.Data) > 0 {
			_, err = writeCover(finalSingerFolder, "folder", meta.Data[0].Relationships.Artists.Data[0].Attributes.Artwork.Url)
			if err != nil {
			}
		}
	}
	covPath, err := writeCover(finalAlbumFolder, "cover", meta.Data[0].Attributes.Artwork.URL)
	if err != nil {
	}
	if Config.SaveAnimatedArtwork && meta.Data[0].Attributes.EditorialVideo.MotionDetailSquare.Video != "" {
		motionvideoUrlSquare, err := extractVideo(meta.Data[0].Attributes.EditorialVideo.MotionDetailSquare.Video)
		if err == nil {
			exists, _ := fileExists(filepath.Join(finalAlbumFolder, "square_animated_artwork.mp4"))
			if !exists {
				cmd := exec.Command("ffmpeg", "-loglevel", "quiet", "-y", "-i", motionvideoUrlSquare, "-c", "copy", filepath.Join(finalAlbumFolder, "square_animated_artwork.mp4"))
				_ = cmd.Run()
			}
		}

		if Config.EmbyAnimatedArtwork {
			cmd3 := exec.Command("ffmpeg", "-i", filepath.Join(finalAlbumFolder, "square_animated_artwork.mp4"), "-vf", "scale=440:-1", "-r", "24", "-f", "gif", filepath.Join(finalAlbumFolder, "folder.jpg"))
			_ = cmd3.Run()
		}

		motionvideoUrlTall, err := extractVideo(meta.Data[0].Attributes.EditorialVideo.MotionDetailTall.Video)
		if err == nil {
			exists, _ := fileExists(filepath.Join(finalAlbumFolder, "tall_animated_artwork.mp4"))
			if !exists {
				cmd := exec.Command("ffmpeg", "-loglevel", "quiet", "-y", "-i", motionvideoUrlTall, "-c", "copy", filepath.Join(finalAlbumFolder, "tall_animated_artwork.mp4"))
				_ = cmd.Run()
			}
		}
	}

	trackTotal := len(meta.Data[0].Relationships.Tracks.Data)
	arr := make([]int, trackTotal)
	for i := 0; i < trackTotal; i++ {
		arr[i] = i + 1
	}
	selected := []int{}

	if dl_song {
		found := false
		for i, track := range meta.Data[0].Relationships.Tracks.Data {
			if urlArg_i == track.ID {
				selected = append(selected, i+1)
				found = true
				break
			}
		}
		if !found {
			return errors.New("指定的单曲ID未在专辑中找到")
		}
	} else if !dl_select {
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
		fmt.Println("Please select from the track options above (multiple options separated by commas, ranges supported, or type 'all' to select all)")
		cyanColor := color.New(color.FgCyan)
		cyanColor.Print("select: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
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

	fmt.Println("正在进行版权预检，请稍候...")
	var workingAccounts []structs.Account
	firstTrackId := meta.Data[0].Relationships.Tracks.Data[0].ID
	for _, acc := range Config.Accounts {
		_, err := getInfoFromAdam(firstTrackId, &acc, acc.Storefront)
		if err == nil {
			workingAccounts = append(workingAccounts, acc)
		} else {
			fmt.Printf("账户 [%s] 无法访问此专辑 (可能无版权)，本次任务将跳过该账户。\n", acc.Name)
		}
	}

	if len(workingAccounts) == 0 {
		return errors.New("所有账户均无法访问此专辑，任务中止")
	}

	albumQualityType := "AAC"
	albumQualityString := "AAC"
	isHires := false
	isLossless := false

	for _, trackIndex := range selected {
		track := meta.Data[0].Relationships.Tracks.Data[trackIndex-1]
		if contains(track.Attributes.AudioTraits, "hi-res-lossless") {
			isHires = true
			break
		}
		if contains(track.Attributes.AudioTraits, "lossless") {
			isLossless = true
		}
	}

	if isHires {
		albumQualityType = "Hi-Res Lossless"
		albumQualityString = "Hi-Res Lossless"
	} else if isLossless {
		albumQualityType = "Lossless"
		albumQualityString = "Lossless"
	}

	var numThreads int
	switch albumQualityType {
	case "Hi-Res Lossless":
		numThreads = Config.HiresDownloadThreads
	case "Lossless":
		numThreads = Config.LosslessDownloadThreads
	default: // "AAC"
		numThreads = Config.AacDownloadThreads
	}

	if numThreads < 1 {
		numThreads = 1
	}

	regionSet := make(map[string]bool)
	for _, acc := range workingAccounts {
		if acc.Storefront != "" {
			regionSet[strings.ToUpper(acc.Storefront)] = true
		}
	}
	var regionNames []string
	for r := range regionSet {
		regionNames = append(regionNames, r)
	}
	sort.Strings(regionNames)
	regionsStr := strings.Join(regionNames, " / ")

	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("%s %s | %s | %s | %s\n",
	    green("音源:"),
		green(albumQualityString),
		green(fmt.Sprintf("%d 线程", numThreads)),
		yellow(regionsStr),
		green(fmt.Sprintf("%d 个账户并行下载", len(workingAccounts))),
	)
	fmt.Println(strings.Repeat("-", 50))

	trackStatuses = make([]TrackStatus, len(selected))
	for i, trackNum := range selected {
		track := meta.Data[0].Relationships.Tracks.Data[trackNum-1]
		manifest, err := getInfoFromAdam(track.ID, mainAccount, storefront)
		quality := "N/A"
		if err == nil && manifest.Attributes.ExtendedAssetUrls.EnhancedHls != "" {
			_, _, quality, err = extractMedia(manifest.Attributes.ExtendedAssetUrls.EnhancedHls, false)
			if err != nil {
				quality = "获取失败"
			}
		} else {
			quality = "AAC 256kbps"
		}

		trackStatuses[i] = TrackStatus{
			Index:       i,
			TrackNum:    trackNum,
			TrackTotal:  trackTotal,
			TrackName:   track.Attributes.Name,
			Quality:     fmt.Sprintf("(%s)", quality),
			Status:      "等待中",
			StatusColor: color.New(color.FgWhite).SprintFunc(),
		}
	}
	var firstPrint = true
	printUI := func() {
		if !firstPrint {
			fmt.Printf("\033[%dA", len(trackStatuses))
		}
		firstPrint = false

		terminalWidth := 120
		colorRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)

		for _, ts := range trackStatuses {
			displayName := ts.TrackName
			prefixStr := fmt.Sprintf("Track %d of %d: ", ts.TrackNum, ts.TrackTotal)
			qualityStr := ts.Quality
			statusStrWithColor := ts.StatusColor(ts.Status)

			plainStatusStr := colorRegex.ReplaceAllString(statusStrWithColor, "")
			prefixRunes := len([]rune(prefixStr))
			suffixRunes := len([]rune(qualityStr)) + len([]rune(" - ")) + len([]rune(plainStatusStr))

			availableRunesForName := terminalWidth - prefixRunes - suffixRunes
			if availableRunesForName < 15 {
				availableRunesForName = 15
			}

			displayNameRunes := []rune(displayName)
			if len(displayNameRunes) > availableRunesForName {
				displayName = string(displayNameRunes[:availableRunesForName-3]) + "..."
			}

			fmt.Printf("\r\033[K%s%s %s - %s\n", prefixStr, displayName, qualityStr, statusStrWithColor)
		}
	}

	updateStatus := func(index int, status string, sColor func(a ...interface{}) string) {
		uiMutex.Lock()
		defer uiMutex.Unlock()
		trackStatuses[index].Status = status
		trackStatuses[index].StatusColor = sColor
		printUI()
	}
	printUI()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, numThreads)

	for i, trackNum := range selected {
		wg.Add(1)
		go func(trackIndexInMeta int, statusIndex int) {
			semaphore <- struct{}{}
			defer func() {
				<-semaphore
				wg.Done()
			}()

			trackData := meta.Data[0].Relationships.Tracks.Data[trackIndexInMeta-1]

			sharedLock.Lock()
			isDone := isInArray(okDict[albumId], trackIndexInMeta)
			sharedLock.Unlock()

			if isDone {
				updateStatus(statusIndex, "已存在", color.New(color.FgCyan).SprintFunc())
				sharedLock.Lock()
				counter.Total++
				counter.Success++
				sharedLock.Unlock()
				return
			}

			progressChan := make(chan runv10.ProgressUpdate, 10)
			go func() {
				for p := range progressChan {
					speedStr := formatSpeed(p.SpeedBPS)
					var status string
					account := &workingAccounts[statusIndex%len(workingAccounts)]
					accountInfo := fmt.Sprintf("%s 账号", strings.ToUpper(account.Storefront))

					if p.Stage == "decrypt" {
						status = fmt.Sprintf("%s 解密中 %d%% (%s)", yellow(accountInfo), p.Percentage, speedStr)
					} else { // "download"
						status = fmt.Sprintf("%s 下载中 %d%% (%s)", yellow(accountInfo), p.Percentage, speedStr)
					}
					updateStatus(statusIndex, status, color.New(color.FgYellow).SprintFunc())
				}
			}()

			err := downloadTrackWithFallback(trackData, meta, albumId, storefront, baseSaveFolder, Codec, covPath, workingAccounts, statusIndex, statusIndex, updateStatus, progressChan)
			close(progressChan)

			sharedLock.Lock()
			counter.Total++
			if err != nil {
				updateStatus(statusIndex, fmt.Sprintf("下载失败: %v", err), color.New(color.FgRed).SprintFunc())
				counter.Error++
			} else {
				updateStatus(statusIndex, "下载完成", color.New(color.FgGreen).SprintFunc())
				counter.Success++
			}
			sharedLock.Unlock()

		}(trackNum, i)
	}

	wg.Wait()
	fmt.Println(strings.Repeat("-", 50))
	return nil
}

func writeMP4Tags(trackPath, lrc string, meta *structs.AutoGenerated, trackNum, trackTotal int) error {
	index := trackNum - 1

	t := &mp4tag.MP4Tags{
		Title:      meta.Data[0].Relationships.Tracks.Data[index].Attributes.Name,
		TitleSort:  meta.Data[0].Relationships.Tracks.Data[index].Attributes.Name,
		Artist:     meta.Data[0].Relationships.Tracks.Data[index].Attributes.ArtistName,
		ArtistSort: meta.Data[0].Relationships.Tracks.Data[index].Attributes.ArtistName,
		Custom: map[string]string{
			"PERFORMER":   meta.Data[0].Relationships.Tracks.Data[index].Attributes.ArtistName,
			"RELEASETIME": meta.Data[0].Relationships.Tracks.Data[index].Attributes.ReleaseDate,
			"ISRC":        meta.Data[0].Relationships.Tracks.Data[index].Attributes.Isrc,
			"LABEL":       meta.Data[0].Attributes.RecordLabel,
			"UPC":         meta.Data[0].Attributes.Upc,
		},
		Composer:     meta.Data[0].Relationships.Tracks.Data[index].Attributes.ComposerName,
		ComposerSort: meta.Data[0].Relationships.Tracks.Data[index].Attributes.ComposerName,
		Date:         meta.Data[0].Attributes.ReleaseDate,
		CustomGenre:  meta.Data[0].Relationships.Tracks.Data[index].Attributes.GenreNames[0],
		Copyright:    meta.Data[0].Attributes.Copyright,
		Publisher:    meta.Data[0].Attributes.RecordLabel,
		Lyrics:       lrc,
	}

	if !strings.Contains(meta.Data[0].ID, "pl.") {
		albumID, err := strconv.ParseUint(meta.Data[0].ID, 10, 32)
		if err == nil {
			t.ItunesAlbumID = int32(albumID)
		}
	}

	if len(meta.Data[0].Relationships.Artists.Data) > 0 {
		if len(meta.Data[0].Relationships.Tracks.Data[index].Relationships.Artists.Data) > 0 {
			artistID, err := strconv.ParseUint(meta.Data[0].Relationships.Tracks.Data[index].Relationships.Artists.Data[0].ID, 10, 32)
			if err == nil {
				t.ItunesArtistID = int32(artistID)
			}
		}
	}

	if strings.Contains(meta.Data[0].ID, "pl.") && !Config.UseSongInfoForPlaylist {
		t.DiscNumber = 1
		t.DiscTotal = 1
		t.TrackNumber = int16(trackNum)
		t.TrackTotal = int16(trackTotal)
		t.Album = meta.Data[0].Attributes.Name
		t.AlbumSort = meta.Data[0].Attributes.Name
		t.AlbumArtist = meta.Data[0].Attributes.ArtistName
		t.AlbumArtistSort = meta.Data[0].Attributes.ArtistName
	} else if strings.Contains(meta.Data[0].ID, "pl.") && Config.UseSongInfoForPlaylist {
		t.DiscNumber = int16(meta.Data[0].Relationships.Tracks.Data[index].Attributes.DiscNumber)
		t.DiscTotal = int16(meta.Data[0].Relationships.Tracks.Data[trackTotal-1].Attributes.DiscNumber)
		t.TrackNumber = int16(meta.Data[0].Relationships.Tracks.Data[index].Attributes.TrackNumber)
		t.TrackTotal = int16(trackTotal)
		t.Album = meta.Data[0].Relationships.Tracks.Data[index].Attributes.AlbumName
		t.AlbumSort = meta.Data[0].Relationships.Tracks.Data[index].Attributes.AlbumName
		t.AlbumArtist = meta.Data[0].Relationships.Tracks.Data[index].Relationships.Albums.Data[0].Attributes.ArtistName
		t.AlbumArtistSort = meta.Data[0].Relationships.Tracks.Data[index].Relationships.Albums.Data[0].Attributes.ArtistName
	} else {
		t.DiscNumber = int16(meta.Data[0].Relationships.Tracks.Data[index].Attributes.DiscNumber)
		t.DiscTotal = int16(meta.Data[0].Relationships.Tracks.Data[trackTotal-1].Attributes.DiscNumber)
		t.TrackNumber = int16(meta.Data[0].Relationships.Tracks.Data[index].Attributes.TrackNumber)
		t.TrackTotal = int16(trackTotal)
		t.Album = meta.Data[0].Relationships.Tracks.Data[index].Attributes.AlbumName
		t.AlbumSort = meta.Data[0].Relationships.Tracks.Data[index].Attributes.AlbumName
		t.AlbumArtist = meta.Data[0].Attributes.ArtistName
		t.AlbumArtistSort = meta.Data[0].Attributes.ArtistName
	}

	if meta.Data[0].Relationships.Tracks.Data[index].Attributes.ContentRating == "explicit" {
		t.ItunesAdvisory = mp4tag.ItunesAdvisoryExplicit
	} else if meta.Data[0].Relationships.Tracks.Data[index].Attributes.ContentRating == "clean" {
		t.ItunesAdvisory = mp4tag.ItunesAdvisoryClean
	} else {
		t.ItunesAdvisory = mp4tag.ItunesAdvisoryNone
	}

	mp4, err := mp4tag.Open(trackPath)
	if err != nil {
		return err
	}
	defer mp4.Close()
	err = mp4.Write(t, []string{})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	pflag.StringVar(&configPath, "config", "", "指定要使用的配置文件路径 (例如: configs/cn.yaml)")
	pflag.StringVar(&outputPath, "output", "", "指定本次任务的唯一输出目录")

	pflag.BoolVar(&dl_atmos, "atmos", false, "Enable atmos download mode")
	pflag.BoolVar(&dl_aac, "aac", false, "Enable adm-aac download mode")
	pflag.BoolVar(&dl_select, "select", false, "Enable selective download")
	pflag.BoolVar(&dl_song, "song", false, "Enable single song download mode")
	pflag.BoolVar(&artist_select, "all-album", false, "Download all artist albums")
	pflag.BoolVar(&debug_mode, "debug", false, "Enable debug mode to show audio quality information")
	alac_max = pflag.Int("alac-max", Config.AlacMax, "Specify the max quality for download alac")
	atmos_max = pflag.Int("atmos-max", Config.AtmosMax, "Specify the max quality for download atmos")
	aac_type = pflag.String("aac-type", Config.AacType, "Select AAC type, aac aac-binaural aac-downmix")
	mv_audio_type = pflag.String("mv-audio-type", Config.MVAudioType, "Select MV audio type, atmos ac3 aac")
	mv_max = pflag.Int("mv-max", Config.MVMax, "Specify the max quality for download MV")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] url1 url2 ...\n", "[main | main.exe | go run main.go]")
		fmt.Println("Options:")
		pflag.PrintDefaults()
	}

	pflag.Parse()

	if configPath == "" {
		configPath = "config.yaml"
	}

	err := loadConfig(configPath)
	if err != nil {
		if os.IsNotExist(err) && configPath == "config.yaml" {
			fmt.Println("错误: 默认配置文件 config.yaml 未找到，请在程序同目录下创建它，或通过 --config 参数指定一个有效的配置文件。")
			pflag.Usage()
			return
		}
		fmt.Printf("加载配置文件 %s 失败: %v\n", configPath, err)
		return
	}

	if outputPath != "" {
		Config.AlacSaveFolder = outputPath
		Config.AtmosSaveFolder = outputPath
	}

	token, err := getToken()
	if err != nil {
		if len(Config.Accounts) > 0 && Config.Accounts[0].AuthorizationToken != "" && Config.Accounts[0].AuthorizationToken != "your-authorization-token" {
			token = strings.Replace(Config.Accounts[0].AuthorizationToken, "Bearer ", "", -1)
		} else {
			fmt.Println("Failed to get developer token.")
			return
		}
	}
	developerToken = token

	args := pflag.Args()
	if len(args) == 0 {
		fmt.Println("No URLs provided. Please provide at least one URL.")
		pflag.Usage()
		return
	}
	os.Args = args
	if strings.Contains(os.Args[0], "/artist/") {
		artistAccount := &Config.Accounts[0]
		urlArtistName, urlArtistID, err := getUrlArtistName(os.Args[0], artistAccount)
		if err != nil {
			fmt.Println("Failed to get artistname.")
			return
		}
		Config.ArtistFolderFormat = strings.NewReplacer(
			"{UrlArtistName}", LimitString(urlArtistName),
			"{ArtistId}", urlArtistID,
		).Replace(Config.ArtistFolderFormat)
		albumArgs, err := checkArtist(os.Args[0], artistAccount, "albums")
		if err != nil {
			fmt.Println("Failed to get artist albums.")
			return
		}
		mvArgs, err := checkArtist(os.Args[0], artistAccount, "music-videos")
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
				if debug_mode {
					continue
				}
				storefront, albumId = checkUrlMv(urlRaw)
				accountForMV, err := getAccountForStorefront(storefront)
				if err != nil {
					fmt.Printf("MV 下载失败: %v\n", err)
					sharedLock.Lock()
					counter.Error++
					sharedLock.Unlock()
					continue
				}

				sharedLock.Lock()
				counter.Total++
				sharedLock.Unlock()
				if len(accountForMV.MediaUserToken) <= 50 {
					sharedLock.Lock()
					counter.Error++
					sharedLock.Unlock()
					continue
				}
				if _, err := exec.LookPath("mp4decrypt"); err != nil {
					sharedLock.Lock()
					counter.Error++
					sharedLock.Unlock()
					continue
				}
				mvSaveDir := strings.NewReplacer(
					"{ArtistName}", "",
					"{UrlArtistName}", "",
					"{ArtistId}", "",
				).Replace(Config.ArtistFolderFormat)
				if mvSaveDir != "" {
					mvSaveDir = filepath.Join(Config.AlacSaveFolder, forbiddenNames.ReplaceAllString(mvSaveDir, "_"))
				} else {
					mvSaveDir = Config.AlacSaveFolder
				}
				err = mvDownloader(albumId, mvSaveDir, storefront, nil, accountForMV)
				if err != nil {
					sharedLock.Lock()
					counter.Error++
					sharedLock.Unlock()
					continue
				}
				sharedLock.Lock()
				counter.Success++
				sharedLock.Unlock()
				continue
			}

			if strings.Contains(urlRaw, "/song/") {
				tempStorefront, _ := checkUrlSong(urlRaw)
				accountForSong, err := getAccountForStorefront(tempStorefront)
				if err != nil {
					fmt.Printf("获取歌曲信息失败: %v\n", err)
					continue
				}
				urlRaw, err = getUrlSong(urlRaw, accountForSong)
				if err != nil {
					fmt.Println("Failed to get Song info.")
					continue
				}
				dl_song = true
			}

			if strings.Contains(urlRaw, "/playlist/") {
				storefront, albumId = checkUrlPlaylist(urlRaw)
			} else {
				storefront, albumId = checkUrl(urlRaw)
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
			err = rip(albumId, storefront, urlArg_i, urlRaw)
			if err != nil {
				fmt.Println("Album failed:", err)
			}
		}
		fmt.Printf("=======  [\u2714 ] Completed: %d/%d  |  [\u26A0 ] Warnings: %d  |  [\u2716 ] Errors: %d  =======\n", counter.Success, counter.Total, counter.Unavailable+counter.NotSong, counter.Error)
		if counter.Error == 0 {
			break
		}
		fmt.Println("Error detected, press Enter to try again...")
		fmt.Scanln()
		fmt.Println("Start trying again...")
		sharedLock.Lock()
		counter = structs.Counter{}
		sharedLock.Unlock()
	}
}

func mvDownloader(adamID string, saveDir string, storefront string, meta *structs.AutoGenerated, account *structs.Account) error {
	MVInfo, err := getMVInfoFromAdam(adamID, account, storefront)
	if err != nil {
		return err
	}

	var trackTotal int
	var trackNum int
	var index int
	if meta != nil {
		trackTotal = len(meta.Data[0].Relationships.Tracks.Data)
		for i, track := range meta.Data[0].Relationships.Tracks.Data {
			if adamID == track.ID {
				index = i
				trackNum = i + 1
			}
		}
	}

	if strings.HasSuffix(saveDir, ".") {
		saveDir = strings.ReplaceAll(saveDir, ".", "")
	}
	saveDir = strings.TrimSpace(saveDir)

	vidPath := filepath.Join(saveDir, fmt.Sprintf("%s_vid.mp4", adamID))
	audPath := filepath.Join(saveDir, fmt.Sprintf("%s_aud.mp4", adamID))
	mvSaveName := fmt.Sprintf("%s (%s)", MVInfo.Data[0].Attributes.Name, adamID)
	if meta != nil {
		mvSaveName = fmt.Sprintf("%02d. %s", trackNum, MVInfo.Data[0].Attributes.Name)
	}

	sanitizedMvSaveName := forbiddenNames.ReplaceAllString(mvSaveName, "_")
	filenameWithExt := fmt.Sprintf("%s.mp4", sanitizedMvSaveName)

	basePath := filepath.Dir(saveDir)
	dirName := filepath.Base(saveDir)
	_, finalDir, finalFilename := ensureSafePath(basePath, "", dirName, filenameWithExt)
	finalSaveDir := filepath.Join(basePath, finalDir)
	mvOutPath := filepath.Join(finalSaveDir, finalFilename)
	os.MkdirAll(finalSaveDir, os.ModePerm)

	exists, _ := fileExists(mvOutPath)
	if exists {
		return nil
	}

	mvm3u8url, _, _ := runv3.GetWebplayback(adamID, developerToken, account.MediaUserToken, true)
	if mvm3u8url == "" {
		return errors.New("media-user-token may wrong or expired")
	}

	os.MkdirAll(saveDir, os.ModePerm)
	videom3u8url, _ := extractVideo(mvm3u8url)
	videokeyAndUrls, _ := runv3.Run(adamID, videom3u8url, developerToken, account.MediaUserToken, true)
	_ = runv3.ExtMvData(videokeyAndUrls, vidPath)
	audiom3u8url, _ := extractMvAudio(mvm3u8url)
	audiokeyAndUrls, _ := runv3.Run(adamID, audiom3u8url, developerToken, account.MediaUserToken, true)
	_ = runv3.ExtMvData(audiokeyAndUrls, audPath)

	tags := []string{
		"tool=",
		fmt.Sprintf("artist=%s", MVInfo.Data[0].Attributes.ArtistName),
		fmt.Sprintf("title=%s", MVInfo.Data[0].Attributes.Name),
		fmt.Sprintf("genre=%s", MVInfo.Data[0].Attributes.GenreNames[0]),
		fmt.Sprintf("created=%s", MVInfo.Data[0].Attributes.ReleaseDate),
		fmt.Sprintf("ISRC=%s", MVInfo.Data[0].Attributes.Isrc),
	}

	if MVInfo.Data[0].Attributes.ContentRating == "explicit" {
		tags = append(tags, "rating=1")
	} else if MVInfo.Data[0].Attributes.ContentRating == "clean" {
		tags = append(tags, "rating=2")
	} else {
		tags = append(tags, "rating=0")
	}

	if meta != nil {
		if meta.Data[0].Type == "playlists" && !Config.UseSongInfoForPlaylist {
			tags = append(tags, "disk=1/1", fmt.Sprintf("album=%s", meta.Data[0].Attributes.Name), fmt.Sprintf("track=%d", trackNum), fmt.Sprintf("tracknum=%d/%d", trackNum, trackTotal), fmt.Sprintf("album_artist=%s", meta.Data[0].Attributes.ArtistName), fmt.Sprintf("performer=%s", meta.Data[0].Relationships.Tracks.Data[index].Attributes.ArtistName), fmt.Sprintf("copyright=%s", meta.Data[0].Attributes.Copyright), fmt.Sprintf("UPC=%s", meta.Data[0].Attributes.Upc))
		} else {
			tags = append(tags, fmt.Sprintf("album=%s", meta.Data[0].Relationships.Tracks.Data[index].Attributes.AlbumName), fmt.Sprintf("disk=%d/%d", meta.Data[0].Relationships.Tracks.Data[index].Attributes.DiscNumber, meta.Data[0].Relationships.Tracks.Data[trackTotal-1].Attributes.DiscNumber), fmt.Sprintf("track=%d", meta.Data[0].Relationships.Tracks.Data[index].Attributes.TrackNumber), fmt.Sprintf("tracknum=%d/%d", meta.Data[0].Relationships.Tracks.Data[index].Attributes.TrackNumber, meta.Data[0].Attributes.TrackCount), fmt.Sprintf("album_artist=%s", meta.Data[0].Attributes.ArtistName), fmt.Sprintf("performer=%s", meta.Data[0].Relationships.Tracks.Data[index].Attributes.ArtistName), fmt.Sprintf("copyright=%s", meta.Data[0].Attributes.Copyright), fmt.Sprintf("UPC=%s", meta.Data[0].Attributes.Upc))
		}
	} else {
		tags = append(tags, fmt.Sprintf("album=%s", MVInfo.Data[0].Attributes.AlbumName), fmt.Sprintf("disk=%d", MVInfo.Data[0].Attributes.DiscNumber), fmt.Sprintf("track=%d", MVInfo.Data[0].Attributes.TrackNumber), fmt.Sprintf("tracknum=%d", MVInfo.Data[0].Attributes.TrackNumber), fmt.Sprintf("performer=%s", MVInfo.Data[0].Attributes.ArtistName))
	}

	var covPath string
	if true {
		thumbURL := MVInfo.Data[0].Attributes.Artwork.URL
		baseThumbName := forbiddenNames.ReplaceAllString(mvSaveName, "_") + "_thumbnail"
		covPath, err = writeCover(finalSaveDir, baseThumbName, thumbURL)
		if err == nil {
			tags = append(tags, fmt.Sprintf("cover=%s", covPath))
		}
	}

	tagsString := strings.Join(tags, ":")
	muxCmd := exec.Command("MP4Box", "-itags", tagsString, "-quiet", "-add", vidPath, "-add", audPath, "-keep-utc", "-new", mvOutPath)
	if err := muxCmd.Run(); err != nil {
		return err
	}
	defer os.Remove(vidPath)
	defer os.Remove(audPath)
	if covPath != "" {
		defer os.Remove(covPath)
	}
	return nil
}

func extractMvAudio(c string) (string, error) {
	MediaUrl, err := url.Parse(c)
	if err != nil {
		return "", err
	}
	resp, err := http.Get(c)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	audioString := string(body)
	from, listType, err := m3u8.DecodeFrom(strings.NewReader(audioString), true)
	if err != nil || listType != m3u8.MASTER {
		return "", errors.New("m3u8 not of media type")
	}
	audio := from.(*m3u8.MasterPlaylist)

	var audioPriority = []string{"audio-atmos", "audio-ac3", "audio-stereo-256"}
	if Config.MVAudioType == "ac3" {
		audioPriority = []string{"audio-ac3", "audio-stereo-256"}
	} else if Config.MVAudioType == "aac" {
		audioPriority = []string{"audio-stereo-256"}
	}

	re := regexp.MustCompile(`_gr(\d+)_`)

	type AudioStream struct {
		URL     string
		Rank    int
		GroupID string
	}
	var audioStreams []AudioStream

	for _, variant := range audio.Variants {
		for _, audiov := range variant.Alternatives {
			if audiov.URI != "" {
				for _, priority := range audioPriority {
					if audiov.GroupId == priority {
						matches := re.FindStringSubmatch(audiov.URI)
						if len(matches) == 2 {
							var rank int
							fmt.Sscanf(matches[1], "%d", &rank)
							streamUrl, _ := MediaUrl.Parse(audiov.URI)
							audioStreams = append(audioStreams, AudioStream{
								URL:     streamUrl.String(),
								Rank:    rank,
								GroupID: audiov.GroupId,
							})
						}
					}
				}
			}
		}
	}
	if len(audioStreams) == 0 {
		return "", errors.New("no suitable audio stream found")
	}
	sort.Slice(audioStreams, func(i, j int) bool {
		return audioStreams[i].Rank > audioStreams[j].Rank
	})
	return audioStreams[0].URL, nil
}

func checkM3u8(b string, f string, account *structs.Account) (string, error) {
	var EnhancedHls string
	if Config.GetM3u8FromDevice {
		adamID := b
		conn, err := net.Dial("tcp", account.GetM3u8Port)
		if err != nil {
			return "none", err
		}
		defer conn.Close()

		adamIDBuffer := []byte(adamID)
		lengthBuffer := []byte{byte(len(adamIDBuffer))}

		_, err = conn.Write(lengthBuffer)
		if err != nil {
			return "none", err
		}
		_, err = conn.Write(adamIDBuffer)
		if err != nil {
			return "none", err
		}
		response, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			return "none", err
		}
		response = bytes.TrimSpace(response)
		if len(response) > 0 {
			EnhancedHls = string(response)
		}
	}
	return EnhancedHls, nil
}

func formatAvailability(available bool, quality string) string {
	if !available {
		return "Not Available"
	}
	return quality
}

func extractMedia(b string, more_mode bool) (string, string, string, error) {
	masterUrl, err := url.Parse(b)
	if err != nil {
		return "", "", "", err
	}
	resp, err := http.Get(b)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", "", errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}
	masterString := string(body)
	from, listType, err := m3u8.DecodeFrom(strings.NewReader(masterString), true)
	if err != nil || listType != m3u8.MASTER {
		return "", "", "", errors.New("m3u8 not of master type")
	}
	master := from.(*m3u8.MasterPlaylist)
	var streamUrl *url.URL
	sort.Slice(master.Variants, func(i, j int) bool {
		return master.Variants[i].AverageBandwidth > master.Variants[j].AverageBandwidth
	})

	var hasAAC, hasLossless, hasHiRes, hasAtmos, hasDolbyAudio bool
	var aacQuality, losslessQuality, hiResQuality, atmosQuality, dolbyAudioQuality string

	for _, variant := range master.Variants {
		if variant.Codecs == "mp4a.40.2" {
			hasAAC = true
			split := strings.Split(variant.Audio, "-")
			if len(split) >= 3 {
				bitrate, _ := strconv.Atoi(split[2])
				currentBitrate := 0
				if aacQuality != "" {
					fmt.Sscanf(aacQuality, "%d kbps", &currentBitrate)
				}
				if bitrate > currentBitrate {
					aacQuality = fmt.Sprintf("%d kbps", bitrate)
				}
			}
		} else if variant.Codecs == "ec-3" && strings.Contains(variant.Audio, "atmos") {
			hasAtmos = true
			split := strings.Split(variant.Audio, "-")
			if len(split) > 0 {
				bitrateStr := split[len(split)-1]
				if len(bitrateStr) == 4 && bitrateStr[0] == '2' {
					bitrateStr = bitrateStr[1:]
				}
				bitrate, _ := strconv.Atoi(bitrateStr)
				currentBitrate := 0
				if atmosQuality != "" {
					fmt.Sscanf(atmosQuality, "%d kbps", &currentBitrate)
				}
				if bitrate > currentBitrate {
					atmosQuality = fmt.Sprintf("%d kbps", bitrate)
				}
			}
		} else if variant.Codecs == "alac" {
			split := strings.Split(variant.Audio, "-")
			if len(split) >= 3 {
				bitDepth := split[len(split)-1]
				sampleRate := split[len(split)-2]
				sampleRateInt, _ := strconv.Atoi(sampleRate)
				if sampleRateInt > 48000 {
					hasHiRes = true
					hiResQuality = fmt.Sprintf("%sbit/%.1fkHz", bitDepth, float64(sampleRateInt)/1000.0)
				} else {
					hasLossless = true
					losslessQuality = fmt.Sprintf("%sbit/%.1fkHz", bitDepth, float64(sampleRateInt)/1000.0)
				}
			}
		} else if variant.Codecs == "ac-3" {
			hasDolbyAudio = true
			split := strings.Split(variant.Audio, "-")
			if len(split) > 0 {
				bitrate, _ := strconv.Atoi(split[len(split)-1])
				dolbyAudioQuality = fmt.Sprintf("%d kbps", bitrate)
			}
		}
	}

	var qualityForDisplay string
	if hasHiRes {
		qualityForDisplay = hiResQuality
	} else if hasLossless {
		qualityForDisplay = losslessQuality
	} else if hasAtmos {
		qualityForDisplay = "Dolby Atmos"
	} else if hasDolbyAudio {
		qualityForDisplay = "Dolby Audio"
	} else if hasAAC {
		qualityForDisplay = "AAC"
	}

	if debug_mode && more_mode {
		fmt.Println("\nDebug: All Available Variants:")
		var data [][]string
		for _, variant := range master.Variants {
			data = append(data, []string{variant.Codecs, variant.Audio, fmt.Sprint(variant.Bandwidth)})
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Codec", "Audio", "Bandwidth"})
		table.SetRowLine(true)
		table.AppendBulk(data)
		table.Render()

		fmt.Println("Available Audio Formats:")
		fmt.Println("------------------------")
		fmt.Printf("AAC             : %s\n", formatAvailability(hasAAC, aacQuality))
		fmt.Printf("Lossless        : %s\n", formatAvailability(hasLossless, losslessQuality))
		fmt.Printf("Hi-Res Lossless : %s\n", formatAvailability(hasHiRes, hiResQuality))
		fmt.Printf("Dolby Atmos     : %s\n", formatAvailability(hasAtmos, atmosQuality))
		fmt.Printf("Dolby Audio     : %s\n", formatAvailability(hasDolbyAudio, dolbyAudioQuality))
		fmt.Println("------------------------")

		return "", "", "", nil
	}
	var qualityForFilename string
	for _, variant := range master.Variants {
		if dl_atmos {
			if variant.Codecs == "ec-3" && strings.Contains(variant.Audio, "atmos") {
				split := strings.Split(variant.Audio, "-")
				length_int, err := strconv.Atoi(split[len(split)-1])
				if err == nil && length_int <= Config.AtmosMax {
					streamUrl, _ = masterUrl.Parse(variant.URI)
					qualityForFilename = fmt.Sprintf("%s kbps", split[len(split)-1])
					break
				}
			} else if variant.Codecs == "ac-3" {
				streamUrl, _ = masterUrl.Parse(variant.URI)
				split := strings.Split(variant.Audio, "-")
				qualityForFilename = fmt.Sprintf("%s kbps", split[len(split)-1])
				break
			}
		} else if dl_aac {
			if variant.Codecs == "mp4a.40.2" {
				aacregex := regexp.MustCompile(`audio-stereo-\d+`)
				replaced := aacregex.ReplaceAllString(variant.Audio, "aac")
				if replaced == Config.AacType {
					streamUrl, _ = masterUrl.Parse(variant.URI)
					split := strings.Split(variant.Audio, "-")
					qualityForFilename = fmt.Sprintf("%s kbps", split[2])
					break
				}
			}
		} else {
			if variant.Codecs == "alac" {
				split := strings.Split(variant.Audio, "-")
				length_int, err := strconv.Atoi(split[len(split)-2])
				if err == nil && length_int <= Config.AlacMax {
					streamUrl, _ = masterUrl.Parse(variant.URI)
					KHZ := float64(length_int) / 1000.0
					qualityForFilename = fmt.Sprintf("%sB-%.1fkHz", split[len(split)-1], KHZ)
					break
				}
			}
		}
	}
	if streamUrl == nil {
		if len(master.Variants) > 0 {
			streamUrl, _ = masterUrl.Parse(master.Variants[0].URI)
		} else {
			return "", "", qualityForDisplay, errors.New("no variants found in playlist")
		}
	}
	return streamUrl.String(), qualityForFilename, qualityForDisplay, nil
}
func extractVideo(c string) (string, error) {
	MediaUrl, err := url.Parse(c)
	if err != nil {
		return "", err
	}
	resp, err := http.Get(c)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	videoString := string(body)

	from, listType, err := m3u8.DecodeFrom(strings.NewReader(videoString), true)
	if err != nil || listType != m3u8.MASTER {
		return "", errors.New("m3u8 not of media type")
	}
	video := from.(*m3u8.MasterPlaylist)

	var streamUrl *url.URL
	sort.Slice(video.Variants, func(i, j int) bool {
		return video.Variants[i].AverageBandwidth > video.Variants[j].AverageBandwidth
	})

	maxHeight := Config.MVMax
	for _, variant := range video.Variants {
		re := regexp.MustCompile(`_(\d+)x(\d+)`)
		matches := re.FindStringSubmatch(variant.URI)
		if len(matches) == 3 {
			height, _ := strconv.Atoi(matches[2])
			if height <= maxHeight {
				streamUrl, _ = MediaUrl.Parse(variant.URI)
				break
			}
		}
	}

	if streamUrl == nil {
		if len(video.Variants) > 0 {
			streamUrl, _ = MediaUrl.Parse(video.Variants[0].URI)
		} else {
			return "", errors.New("no suitable video stream found")
		}
	}
	return streamUrl.String(), nil
}

func getInfoFromAdam(adamId string, account *structs.Account, storefront string) (*structs.SongData, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("https://amp-api.music.apple.com/v1/catalog/%s/songs/%s", storefront, adamId), nil)
	if err != nil {
		return nil, err
	}
	query := url.Values{}
	query.Set("extend", "extendedAssetUrls")
	query.Set("include", "albums")
	query.Set("l", Config.Language)
	request.URL.RawQuery = query.Encode()

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", developerToken))
	request.Header.Set("User-Agent", "iTunes/12.11.3 (Windows; Microsoft Windows 10 x64 Professional Edition (Build 19041); x64) AppleWebKit/7611.1022.4001.1 (dt:2)")
	request.Header.Set("Origin", "https://music.apple.com")

	do, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return nil, errors.New(do.Status)
	}

	obj := new(structs.ApiResult)
	err = json.NewDecoder(do.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}

	for _, d := range obj.Data {
		if d.ID == adamId {
			return &d, nil
		}
	}
	return nil, nil
}

func getMVInfoFromAdam(adamId string, account *structs.Account, storefront string) (*structs.AutoGeneratedMusicVideo, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("https://amp-api.music.apple.com/v1/catalog/%s/music-videos/%s", storefront, adamId), nil)
	if err != nil {
		return nil, err
	}
	query := url.Values{}
	query.Set("l", Config.Language)
	request.URL.RawQuery = query.Encode()
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", developerToken))
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	request.Header.Set("Origin", "https://music.apple.com")

	do, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK {
		return nil, errors.New(do.Status)
	}

	obj := new(structs.AutoGeneratedMusicVideo)
	err = json.NewDecoder(do.Body).Decode(&obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func getToken() (string, error) {
	req, err := http.NewRequest("GET", "https://beta.music.apple.com", nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	regex := regexp.MustCompile(`/assets/index-legacy-[^/]+\.js`)
	indexJsUri := regex.FindString(string(body))
	if indexJsUri == "" {
		return "", errors.New("could not find JS asset URL in HTML")
	}
	req, err = http.NewRequest("GET", "https://beta.music.apple.com"+indexJsUri, nil)
	if err != nil {
		return "", err
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	regex = regexp.MustCompile(`eyJh([^"]*)`)
	token := regex.FindString(string(body))
	if token == "" {
		return "", errors.New("could not find developer token in JS file")
	}
	return token, nil
}
