package core

import (
	"errors"
	"fmt"
	"main/internal/logger"
	"main/utils/structs"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

var (
	ForbiddenNames   = regexp.MustCompile(`[/\\<>:"|?*]`)
	Dl_atmos         bool
	Dl_aac           bool
	Dl_select        bool
	Dl_song          bool
	Artist_select    bool
	Debug_mode       bool
	DisableDynamicUI bool // 禁用动态UI的标志，启用后使用纯日志输出
	ForceDownload    bool // 强制下载模式，覆盖已存在的文件
	Alac_max         *int
	Atmos_max        *int
	Mv_max           *int
	Mv_audio_type    *string
	Aac_type         *string
	StartFrom        int // 从第几个链接开始下载（从1开始计数）
	Config           structs.ConfigSet
	Counter          structs.Counter
	OkDict           = make(map[string][]int)
	ConfigPath       string
	OutputPath       string
	SharedLock       sync.Mutex
	DeveloperToken   string
	MaxPathLength    int
)

type TrackStatus struct {
	Index        int
	TrackNum     int
	TrackTotal   int
	TrackName    string
	Quality      string
	Status       string
	StatusColor  func(a ...interface{}) string
	LastUpdateNs int64 // 最后更新时间（纳秒），用于防抖
}

var UiMutex sync.Mutex

var RipLock sync.Mutex

var TrackStatuses []TrackStatus

func InitCounter() structs.Counter {
	return structs.Counter{}
}

func InitFlags() {
	pflag.StringVar(&ConfigPath, "config", "", "指定要使用的配置文件路径 (例如: configs/cn.yaml)")
	pflag.StringVar(&OutputPath, "output", "", "指定本次任务的唯一输出目录")

	pflag.BoolVar(&Dl_atmos, "atmos", false, "启用杜比全景声下载模式")
	pflag.BoolVar(&Dl_aac, "aac", false, "启用 AAC 下载模式")
	pflag.BoolVar(&Dl_select, "select", false, "启用选择性下载模式（可选择要下载的曲目）")
	pflag.BoolVar(&Dl_song, "song", false, "启用单曲下载模式")
	pflag.BoolVar(&Artist_select, "all-album", false, "下载歌手的所有专辑")
	pflag.BoolVar(&Debug_mode, "debug", false, "启用调试模式，显示音频质量信息")
	pflag.BoolVar(&DisableDynamicUI, "no-ui", false, "禁用动态终端UI，回退到纯日志输出模式（用于CI/调试或兼容性）")
	pflag.BoolVar(&ForceDownload, "cx", false, "强制下载模式，覆盖已存在的文件")
	pflag.IntVar(&StartFrom, "start", 0, "从 TXT 文件的第几个链接开始下载（从 1 开始计数，例如：--start 44）")
	Alac_max = pflag.Int("alac-max", 0, "指定 ALAC 下载的最大音质（如：192000, 96000, 48000）")
	Atmos_max = pflag.Int("atmos-max", 0, "指定 Dolby Atmos 下载的最大音质（如：2768, 2448）")
	Aac_type = pflag.String("aac-type", "aac", "选择 AAC 类型（可选：aac, aac-binaural, aac-downmix）")
	Mv_audio_type = pflag.String("mv-audio-type", "atmos", "选择 MV 音轨类型（可选：atmos, ac3, aac）")
	Mv_max = pflag.Int("mv-max", 1080, "指定 MV 下载的最大分辨率（如：2160, 1080, 720）")
}

func LoadConfig(configPath string) error {
	if configPath == "" {
		ConfigPath = "config.yaml"
	} else {
		ConfigPath = configPath
	}

	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return err
	}

	// 替换环境变量引用（支持 ${VAR_NAME} 格式）
	configContent := string(data)
	configContent = os.ExpandEnv(configContent)

	err = yaml.Unmarshal([]byte(configContent), &Config)
	if err != nil {
		return err
	}

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	if len(Config.Accounts) == 0 {
		return errors.New(red("配置错误: 'accounts' 列表为空，请在 config.yaml 中至少配置一个账户"))
	}

	if Config.BufferSizeKB <= 0 {
		Config.BufferSizeKB = 4096
		logger.Info("📌 配置文件中未设置 'BufferSizeKB'，自动设为默认值 4096KB (4MB)")
	}

	if Config.NetworkReadBufferKB <= 0 {
		Config.NetworkReadBufferKB = 4096
		logger.Info("📌 配置文件中未设置 'NetworkReadBufferKB'，自动设为默认值 4096KB (4MB)")
	}

	useAutoDetect := true
	if Config.MaxPathLength > 0 {
		MaxPathLength = Config.MaxPathLength
		useAutoDetect = false
		logger.Info("%s%s",
			green("📌 从配置文件强制使用最大路径长度限制: "),
			red(fmt.Sprintf("%d", MaxPathLength)),
		)
	}

	if useAutoDetect {
		if runtime.GOOS == "windows" {
			MaxPathLength = 255
			logger.Info("%s%d",
				green("📌 检测到 Windows 系统, 已自动设置最大路径长度限制为: "),
				MaxPathLength,
			)
		} else {
			MaxPathLength = 4096
			logger.Info("%s%s%s%d",
				green("📌 检测到 "),
				red(runtime.GOOS),
				green(" 系统, 已自动设置最大路径长度限制为: "),
				MaxPathLength,
			)
		}
	}

	if *Alac_max == 0 {
		Alac_max = &Config.AlacMax
	}
	if *Atmos_max == 0 {
		Atmos_max = &Config.AtmosMax
	}
	// 如果命令行中没有指定aac-type（使用默认值），则使用配置文件的值
	if *Aac_type == "aac" {
		Aac_type = &Config.AacType
	}
	// 如果命令行中指定了aac-type，则更新Config.AacType以保持一致性
	if *Aac_type != "aac" {
		Config.AacType = *Aac_type
	}
	if *Mv_audio_type == "atmos" {
		Mv_audio_type = &Config.MVAudioType
	}
	if *Mv_max == 1080 {
		Mv_max = &Config.MVMax
	}

	// 设置缓存文件夹默认值
	if Config.CacheFolder == "" {
		Config.CacheFolder = "./Cache"
	}

	// 如果启用缓存，显示缓存配置信息
	if Config.EnableCache {
		logger.Info("%s%s",
			green("📌 缓存中转机制已启用，缓存路径: "),
			red(Config.CacheFolder),
		)
	}

	// 设置分批下载默认值
	if Config.BatchSize == 0 {
		Config.BatchSize = 20
		logger.Info("📌 配置文件中未设置 'batch-size'，自动设为默认值 20（分批处理模式）")
	} else if Config.BatchSize < 0 {
		Config.BatchSize = 0
		logger.Info("📌 'batch-size' 设置为负数，已调整为 0（禁用分批，一次性处理）")
	}

	// 设置工作-休息循环默认值
	if Config.WorkRestEnabled {
		if Config.WorkDurationMinutes <= 0 {
			Config.WorkDurationMinutes = 5
			logger.Info("📌 配置文件中未设置 'work-duration-minutes'，自动设为默认值 5 分钟")
		}
		if Config.RestDurationMinutes <= 0 {
			Config.RestDurationMinutes = 1
			logger.Info("📌 配置文件中未设置 'rest-duration-minutes'，自动设为默认值 1 分钟")
		}
	}

	return nil
}

func GetAccountForStorefront(storefront string) (*structs.Account, error) {
	if len(Config.Accounts) == 0 {
		return nil, errors.New("无可用账户")
	}

	for i := range Config.Accounts {
		acc := &Config.Accounts[i]
		if strings.EqualFold(acc.Storefront, storefront) {
			return acc, nil
		}
	}

	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	logger.Warn("%s 未找到与 %s 匹配的账户,将尝试使用 %s 等区域进行下载",
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
