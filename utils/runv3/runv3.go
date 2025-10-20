package runv3

import (
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"

	"github.com/sky8282/requests"
	"google.golang.org/protobuf/proto"

	//"log/slog"
	"main/internal/logger"
	cdm "main/utils/runv3/cdm"
	key "main/utils/runv3/key"
	"os"

	"bytes"
	"errors"
	"io"

	"github.com/Eyevinn/mp4ff/mp4"

	//"io/ioutil"
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/grafov/m3u8"
	"github.com/schollz/progressbar/v3"
)

// 自定义 context key 类型，避免使用内置字符串类型
type contextKey string

const (
	psshContextKey   contextKey = "pssh"
	adamIdContextKey contextKey = "adamId"
)

type PlaybackLicense struct {
	ErrorCode  int    `json:"errorCode"`
	License    string `json:"license"`
	RenewAfter int    `json:"renew-after"`
	Status     int    `json:"status"`
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getPSSH(contentId string, kidBase64 string) (string, error) {
	kidBytes, err := base64.StdEncoding.DecodeString(kidBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 KID: %v", err)
	}
	contentIdEncoded := base64.StdEncoding.EncodeToString([]byte(contentId))
	algo := cdm.WidevineCencHeader_AESCTR
	widevineCencHeader := &cdm.WidevineCencHeader{
		KeyId:     [][]byte{kidBytes},
		Algorithm: &algo,
		Provider:  new(string),
		ContentId: []byte(contentIdEncoded),
		Policy:    new(string),
	}
	widevineCenc, err := proto.Marshal(widevineCencHeader)
	if err != nil {
		return "", fmt.Errorf("failed to marshal WidevineCencHeader: %v", err)
	}
	//最前面添加32字节
	widevineCenc = append([]byte("0123456789abcdef0123456789abcdef"), widevineCenc...)
	pssh := base64.StdEncoding.EncodeToString(widevineCenc)
	return pssh, nil
}
func BeforeRequest(cl *requests.Client, preCtx context.Context, method string, href string, options ...requests.RequestOption) (resp *requests.Response, err error) {
	data := options[0].Data

	// 安全地从 context 中获取值
	pssh, ok := preCtx.Value(psshContextKey).(string)
	if !ok {
		return nil, fmt.Errorf("pssh not found in context or invalid type")
	}

	adamId, ok := preCtx.Value(adamIdContextKey).(string)
	if !ok {
		return nil, fmt.Errorf("adamId not found in context or invalid type")
	}

	jsondata := map[string]interface{}{
		"challenge":      base64.StdEncoding.EncodeToString(data.([]byte)),
		"key-system":     "com.widevine.alpha",
		"uri":            "data:;base64," + pssh,
		"adamId":         adamId,
		"isLibrary":      false,
		"user-initiated": true,
	}
	options[0].Data = nil
	options[0].Json = jsondata
	resp, err = cl.Request(preCtx, method, href, options...)
	if err != nil {
		logger.Error("Request failed: %v", err)
	}

	return
}
func AfterRequest(Response *requests.Response) ([]byte, error) {
	var ResponseData PlaybackLicense
	_, err := Response.Json(&ResponseData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	if ResponseData.ErrorCode != 0 || ResponseData.Status != 0 {
		return nil, fmt.Errorf("error code: %d", ResponseData.ErrorCode)
	}
	License, err := base64.StdEncoding.DecodeString(ResponseData.License)
	if err != nil {
		return nil, fmt.Errorf("failed to decode license: %v", err)
	}
	return License, nil
}
func GetWebplayback(adamId string, authtoken string, mutoken string, mvmode bool) (string, string, error) {
	url := "https://play.music.apple.com/WebObjects/MZPlay.woa/wa/webPlayback"
	postData := map[string]string{
		"salableAdamId": adamId,
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		logger.Error("Error encoding JSON: %v", err)
		return "", "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		logger.Error("Error creating request: %v", err)
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://music.apple.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Referer", "https://music.apple.com/")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authtoken))
	req.Header.Set("x-apple-music-user-token", mutoken)
	// 创建 HTTP 客户端
	//client := &http.Client{}
	resp, err := http.DefaultClient.Do(req)
	// 发送请求
	//resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error sending request: %v", err)
		return "", "", err
	}
	defer resp.Body.Close()
	//fmt.Println("Response Status:", resp.Status)
	obj := new(Songlist)
	err = json.NewDecoder(resp.Body).Decode(&obj)
	if err != nil {
		logger.Error("json err: %v", err)
		return "", "", err
	}
	if len(obj.List) > 0 {
		if mvmode {
			return obj.List[0].HlsPlaylistUrl, "", nil
		}
		
		// 调试：打印所有可用的assets
		logger.Debug("🔍 webPlayback返回的Assets:")
		for i, asset := range obj.List[0].Assets {
			logger.Debug("  [%d] Flavor=%s, URL=%s", i, asset.Flavor, asset.URL[:min(80, len(asset.URL))]+"...")
		}
		
		// 遍历 Assets，查找匹配的flavor
		for i := range obj.List[0].Assets {
			if obj.List[0].Assets[i].Flavor == "28:ctrp256" {
				kidBase64, fileurl, err := extractKidBase64(obj.List[0].Assets[i].URL, false)
				if err != nil {
					return "", "", err
				}
				return fileurl, kidBase64, nil
			}
			continue
		}
	}
	return "", "", errors.New("Unavailable")
}

type Songlist struct {
	List []struct {
		Hlsurl         string `json:"hls-key-cert-url"`
		HlsPlaylistUrl string `json:"hls-playlist-url"`
		Assets         []struct {
			Flavor string `json:"flavor"`
			URL    string `json:"URL"`
		} `json:"assets"`
	} `json:"songList"`
	Status int `json:"status"`
}

func extractKidBase64(b string, mvmode bool) (string, string, error) {
	resp, err := http.Get(b)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	masterString := string(body)
	from, listType, err := m3u8.DecodeFrom(strings.NewReader(masterString), true)
	if err != nil {
		return "", "", err
	}
	var kidbase64 string
	var urlBuilder strings.Builder
	if listType == m3u8.MEDIA {
		mediaPlaylist := from.(*m3u8.MediaPlaylist)
		if mediaPlaylist.Key != nil {
			split := strings.Split(mediaPlaylist.Key.URI, ",")
			kidbase64 = split[1]
			lastSlashIndex := strings.LastIndex(b, "/")
			// 截取最后一个斜杠之前的部分
			urlBuilder.WriteString(b[:lastSlashIndex])
			urlBuilder.WriteString("/")
			urlBuilder.WriteString(mediaPlaylist.Map.URI)
			//fileurl = b[:lastSlashIndex] + "/" + mediaPlaylist.Map.URI
			//fmt.Println("Extracted URI:", mediaPlaylist.Map.URI)
			if mvmode {
				for _, segment := range mediaPlaylist.Segments {
					if segment != nil {
						//fmt.Println("Extracted URI:", segment.URI)
						urlBuilder.WriteString(";")
						urlBuilder.WriteString(b[:lastSlashIndex])
						urlBuilder.WriteString("/")
						urlBuilder.WriteString(segment.URI)
						//fileurl = fileurl + ";" + b[:lastSlashIndex] + "/" + segment.URI
					}
				}
			}
		} else {
			logger.Warn("No key information found")
		}
	} else {
		logger.Warn("Not a media playlist")
	}
	return kidbase64, urlBuilder.String(), nil
}
func extsong(b string) bytes.Buffer {
	resp, err := http.Get(b)
	if err != nil {
		// 静默处理错误，不干扰UI
	}
	defer resp.Body.Close()
	var buffer bytes.Buffer

	// 将进度条输出重定向到 io.Discard，避免干扰UI系统
	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetElapsedTime(false),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetDescription("Downloading..."),
		progressbar.OptionSetWriter(io.Discard), // 禁用输出
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "",
			SaucerHead:    "",
			SaucerPadding: "",
			BarStart:      "",
			BarEnd:        "",
		}),
	)
	// 忽略 io.Copy 错误，因为这是下载进度显示，主要内容已经缓冲
	_, _ = io.Copy(io.MultiWriter(&buffer, bar), resp.Body)
	return buffer
}
func Run(adamId string, trackpath string, authtoken string, mutoken string, mvmode bool) (string, error) {
	var keystr string //for mv key
	var fileurl string
	var kidBase64 string
	var err error
	if mvmode {
		kidBase64, fileurl, err = extractKidBase64(trackpath, true)
		if err != nil {
			return "", err
		}
	} else {
		fileurl, kidBase64, err = GetWebplayback(adamId, authtoken, mutoken, false)
		if err != nil {
			return "", err
		}
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, psshContextKey, kidBase64)
	ctx = context.WithValue(ctx, adamIdContextKey, adamId)
	pssh, err := getPSSH("", kidBase64)
	//fmt.Println(pssh)
	if err != nil {
		logger.Error("getPSSH failed: %v", err)
		return "", err
	}
	headers := map[string]interface{}{
		"authorization":            "Bearer " + authtoken,
		"x-apple-music-user-token": mutoken,
	}
	client, _ := requests.NewClient(context.TODO(), requests.ClientOption{
		Headers: headers,
	})
	key := key.Key{
		ReqCli:        client,
		BeforeRequest: BeforeRequest,
		AfterRequest:  AfterRequest,
	}
	key.CdmInit()
	var keybt []byte
	if strings.Contains(adamId, "ra.") {
		keystr, keybt, err = key.GetKey(ctx, "https://play.itunes.apple.com/WebObjects/MZPlay.woa/web/radio/versions/1/license", pssh, nil)
		if err != nil {
			logger.Error("GetKey failed (radio): %v", err)
			return "", err
		}
	} else {
		keystr, keybt, err = key.GetKey(ctx, "https://play.itunes.apple.com/WebObjects/MZPlay.woa/wa/acquireWebPlaybackLicense", pssh, nil)
		if err != nil {
			logger.Error("GetKey failed: %v", err)
			return "", err
		}
	}
	if mvmode {
		keyAndUrls := "1:" + keystr + ";" + fileurl
		return keyAndUrls, nil
	}
	body := extsong(fileurl)
	// 静默下载，不打印以避免干扰UI系统
	//bodyReader := bytes.NewReader(body)
	var buffer bytes.Buffer

	err = DecryptMP4(&body, keybt, &buffer)
	if err != nil {
		// 静默处理解密错误，不干扰UI
		return "", err
	}
	// 解密成功，静默继续
	// create output file
	ofh, err := os.Create(trackpath)
	if err != nil {
		logger.Error("创建文件失败: %v", err)
		return "", err
	}
	defer ofh.Close()

	_, err = ofh.Write(buffer.Bytes())
	if err != nil {
		logger.Error("写入文件失败: %v", err)
		return "", err
	}
	return "", nil
}

// Segment 结构体用于在 Channel 中传递分段数据
type Segment struct {
	Index int
	Data  []byte
}

func downloadSegment(url string, index int, wg *sync.WaitGroup, segmentsChan chan<- Segment, client *http.Client, limiter chan struct{}) {
	// 函数退出时，从 limiter 中接收一个值，释放一个并发槽位
	defer func() {
		<-limiter
		wg.Done()
	}()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("错误(分段 %d): 创建请求失败: %v", index, err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("错误(分段 %d): 下载失败: %v", index, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("错误(分段 %d): 服务器返回状态码 %d", index, resp.StatusCode)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("错误(分段 %d): 读取数据失败: %v", index, err)
		return
	}

	// 将下载好的分段（包含序号和数据）发送到 Channel
	segmentsChan <- Segment{Index: index, Data: data}
}

// fileWriter 从 Channel 接收分段并按顺序写入文件
func fileWriter(wg *sync.WaitGroup, segmentsChan <-chan Segment, outputFile io.Writer, totalSegments int) {
	defer wg.Done()

	// 缓冲区，用于存放乱序到达的分段
	// key 是分段序号，value 是分段数据
	segmentBuffer := make(map[int][]byte)
	nextIndex := 0 // 期望写入的下一个分段的序号

	for segment := range segmentsChan {
		// 检查收到的分段是否是当前期望的
		if segment.Index == nextIndex {
			//fmt.Printf("写入分段 %d\n", segment.Index)
			_, err := outputFile.Write(segment.Data)
			if err != nil {
				logger.Error("错误(分段 %d): 写入文件失败: %v", segment.Index, err)
			}
			nextIndex++

			// 检查缓冲区中是否有下一个连续的分段
			for {
				data, ok := segmentBuffer[nextIndex]
				if !ok {
					break // 缓冲区里没有下一个，跳出循环，等待下一个分段到达
				}

				//fmt.Printf("从缓冲区写入分段 %d\n", nextIndex)
				_, err := outputFile.Write(data)
				if err != nil {
					logger.Error("错误(分段 %d): 从缓冲区写入文件失败: %v", nextIndex, err)
				}
				// 从缓冲区删除已写入的分段，释放内存
				delete(segmentBuffer, nextIndex)
				nextIndex++
			}
		} else {
			// 如果不是期望的分段，先存入缓冲区
			//fmt.Printf("缓冲分段 %d (等待 %d)\n", segment.Index, nextIndex)
			segmentBuffer[segment.Index] = segment.Data
		}
	}

	// 确保所有分段都已写入
	if nextIndex != totalSegments {
		logger.Warn("警告: 写入完成，但似乎有分段丢失。期望 %d 个, 实际写入 %d 个。", totalSegments, nextIndex)
	}
}

// getTotalSize 并发获取所有分片的总大小
func getTotalSize(urls []string, client *http.Client) int64 {
	var totalSize int64
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 限制并发数，避免过多HEAD请求
	semaphore := make(chan struct{}, 10)

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			req, err := http.NewRequest("HEAD", u, nil)
			if err != nil {
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				size := resp.ContentLength
				if size > 0 {
					mu.Lock()
					totalSize += size
					mu.Unlock()
				}
			}
		}(url)
	}

	wg.Wait()
	return totalSize
}

func ExtMvData(keyAndUrls string, savePath string) error {
	return ExtMvDataWithDesc(keyAndUrls, savePath, "")
}

func ExtMvDataWithDesc(keyAndUrls string, savePath string, description string) error {
	segments := strings.Split(keyAndUrls, ";")
	key := segments[0]
	//fmt.Println(key)
	urls := segments[1:]
	tempFile, err := os.CreateTemp("", "enc_mv_data-*.mp4")
	if err != nil {
		logger.Error("创建文件失败：%v", err)
		return err
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	var downloadWg, writerWg sync.WaitGroup
	segmentsChan := make(chan Segment, len(urls))
	// --- 新增代码: 定义最大并发数 ---
	const maxConcurrency = 10
	// --- 新增代码: 创建带缓冲的 Channel 作为信号量 ---
	limiter := make(chan struct{}, maxConcurrency)
	client := &http.Client{}

	// 获取总大小：并发发送 HEAD 请求
	totalSize := getTotalSize(urls, client)

	// 设置描述文本
	desc := description
	if desc == "" {
		desc = "Downloading..."
	}

	// 初始化进度条（简洁模式：只显示文本信息，不显示图形条）
	// 输出重定向到 io.Discard 避免干扰UI系统
	var bar *progressbar.ProgressBar
	if totalSize > 0 {
		bar = progressbar.NewOptions64(
			totalSize,
			progressbar.OptionSetDescription(desc),
			progressbar.OptionShowBytes(true),
			progressbar.OptionShowCount(),
			progressbar.OptionSetWidth(0),           // 0 = 不显示进度条图形
			progressbar.OptionSetWriter(io.Discard), // 禁用输出，避免干扰UI
			progressbar.OptionThrottle(100*time.Millisecond),
		)
	} else {
		// 如果无法获取总大小，使用未知大小模式
		bar = progressbar.NewOptions64(
			-1,
			progressbar.OptionSetDescription(desc),
			progressbar.OptionShowBytes(true),
			progressbar.OptionShowCount(),
			progressbar.OptionSetWriter(io.Discard), // 禁用输出，避免干扰UI
			progressbar.OptionSetWidth(0),           // 0 = 不显示进度条图形
			progressbar.OptionThrottle(100*time.Millisecond),
		)
	}
	barWriter := io.MultiWriter(tempFile, bar)

	// 启动写入 Goroutine
	writerWg.Add(1)
	go fileWriter(&writerWg, segmentsChan, barWriter, len(urls))

	// 启动下载 Goroutines
	for i, url := range urls {
		//fmt.Printf("请求启动任务 %d...\n", i)
		limiter <- struct{}{}
		//fmt.Printf("...任务 %d 已启动\n", i)

		downloadWg.Add(1)
		go downloadSegment(url, i, &downloadWg, segmentsChan, client, limiter)
	}

	downloadWg.Wait()
	close(segmentsChan)

	writerWg.Wait()

	if err := tempFile.Close(); err != nil {

		return err
	}

	cmd1 := exec.Command("mp4decrypt", "--key", key, tempFile.Name(), savePath)
	outlog, err := cmd1.CombinedOutput()
	if err != nil {
		return fmt.Errorf("decrypt failed: %w, output: %s", err, string(outlog))
	}

	// 清理可能产生的 out_ 前缀文件
	outDir := filepath.Dir(savePath)
	outFileName := "out_" + filepath.Base(savePath)
	outFilePath := filepath.Join(outDir, outFileName)
	if _, err := os.Stat(outFilePath); err == nil {
		os.Remove(outFilePath)
	}

	return nil
}

// DecryptMP4 decrypts a fragmented MP4 file with keys from widevice license. Supports CENC and CBCS schemes.
func DecryptMP4(r io.Reader, key []byte, w io.Writer) error {
	// Initialization
	inMp4, err := mp4.DecodeFile(r)
	if err != nil {
		return fmt.Errorf("failed to decode file: %w", err)
	}
	if !inMp4.IsFragmented() {
		return errors.New("file is not fragmented")
	}
	// Handle init segment
	if inMp4.Init == nil {
		return errors.New("no init part of file")
	}
	decryptInfo, err := mp4.DecryptInit(inMp4.Init)
	if err != nil {
		return fmt.Errorf("failed to decrypt init: %w", err)
	}
	if err = inMp4.Init.Encode(w); err != nil {
		return fmt.Errorf("failed to write init: %w", err)
	}
	// Decode segments
	for _, seg := range inMp4.Segments {
		if err = mp4.DecryptSegment(seg, decryptInfo, key); err != nil {
			if err.Error() == "no senc box in traf" {
				// No SENC box, skip decryption for this segment as samples can have
				// unencrypted segments followed by encrypted segments. See:
				// https://github.com/iyear/gowidevine/pull/26#issuecomment-2385960551
				err = nil
			} else {
				return fmt.Errorf("failed to decrypt segment: %w", err)
			}
		}
		if err = seg.Encode(w); err != nil {
			return fmt.Errorf("failed to encode segment: %w", err)
		}
	}
	return nil
}
