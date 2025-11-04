package config

import (
	"fmt"
	"main/internal/constants"
	"main/internal/logger"
	"main/utils/structs"
	"os"
	"path/filepath"
	"strings"
)

// ValidationError 表示配置验证错误
type ValidationError struct {
	Field   string // 配置字段名
	Message string // 错误描述
}

// ValidationResult 表示验证结果
type ValidationResult struct {
	Errors   []ValidationError // 错误列表
	Warnings []ValidationError // 警告列表
}

// IsValid 返回验证是否通过（没有错误）
func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

// HasWarnings 返回是否有警告
func (r *ValidationResult) HasWarnings() bool {
	return len(r.Warnings) > 0
}

// Print 输出验证结果
func (r *ValidationResult) Print() {
	if len(r.Errors) > 0 {
		logger.Error("❌ 配置文件验证失败，发现 %d 个错误:", len(r.Errors))
		for _, err := range r.Errors {
			logger.Error("  - %s: %s", err.Field, err.Message)
		}
	}
	
	if len(r.Warnings) > 0 {
		logger.Warn("⚠️  配置文件有 %d 个警告:", len(r.Warnings))
		for _, warn := range r.Warnings {
			logger.Warn("  - %s: %s", warn.Field, warn.Message)
		}
	}
	
	if r.IsValid() && !r.HasWarnings() {
		logger.Info("✅ 配置文件验证通过")
	}
}

// ValidateConfig 验证配置文件的完整性和有效性
func ValidateConfig(cfg *structs.ConfigSet) *ValidationResult {
	result := &ValidationResult{
		Errors:   []ValidationError{},
		Warnings: []ValidationError{},
	}
	
	// 1. 验证账号配置
	validateAccounts(cfg, result)
	
	// 2. 验证保存路径
	validateSavePaths(cfg, result)
	
	// 3. 验证缓存配置
	validateCache(cfg, result)
	
	// 4. 验证下载性能配置
	validateDownloadPerformance(cfg, result)
	
	// 5. 验证工作-休息循环配置
	validateWorkRest(cfg, result)
	
	// 6. 验证音质配置
	validateAudioQuality(cfg, result)
	
	// 7. 验证 MV 配置
	validateMVConfig(cfg, result)
	
	// 8. 验证路径限制
	validatePathLimits(cfg, result)
	
	// 9. 验证日志配置
	validateLogging(cfg, result)
	
	return result
}

// validateAccounts 验证账号配置
func validateAccounts(cfg *structs.ConfigSet, result *ValidationResult) {
	if len(cfg.Accounts) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "accounts",
			Message: "至少需要配置一个 Apple Music 账号",
		})
		return
	}
	
	storefronts := make(map[string]int)
	for i, account := range cfg.Accounts {
		fieldPrefix := fmt.Sprintf("accounts[%d]", i)
		
		// 验证 storefront
		if account.Storefront == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fieldPrefix + ".storefront",
				Message: "storefront 不能为空",
			})
		} else {
			storefronts[account.Storefront]++
		}
		
		// 验证 media_user_token
		if account.MediaUserToken == "" {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   fieldPrefix + ".media_user_token",
				Message: "media_user_token 为空（部分功能可能无法使用）",
			})
		} else if len(account.MediaUserToken) < constants.MinTokenLength {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   fieldPrefix + ".media_user_token",
				Message: fmt.Sprintf("token 长度过短（%d < %d）", len(account.MediaUserToken), constants.MinTokenLength),
			})
		}
	}
	
	// 检查重复的 storefront
	for sf, count := range storefronts {
		if count > 1 {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   "accounts",
				Message: fmt.Sprintf("storefront '%s' 配置了 %d 次（可能是重复配置）", sf, count),
			})
		}
	}
}

// validateSavePaths 验证保存路径
func validateSavePaths(cfg *structs.ConfigSet, result *ValidationResult) {
	paths := map[string]string{
		"alac_save_folder":  cfg.AlacSaveFolder,
		"atmos_save_folder": cfg.AtmosSaveFolder,
		"aac_save_folder":   cfg.AacSaveFolder,
	}
	
	for field, path := range paths {
		if path == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   field,
				Message: "保存路径不能为空",
			})
			continue
		}
		
		// 检查路径是否可访问（尝试创建目录）
		absPath, err := filepath.Abs(path)
		if err != nil {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   field,
				Message: fmt.Sprintf("路径解析失败: %v", err),
			})
			continue
		}
		
		// 检查父目录是否存在或可创建
		if err := os.MkdirAll(absPath, 0755); err != nil {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   field,
				Message: fmt.Sprintf("无法创建目录: %v", err),
			})
		}
	}
}

// validateCache 验证缓存配置
func validateCache(cfg *structs.ConfigSet, result *ValidationResult) {
	if cfg.CacheFolder == "" {
		result.Warnings = append(result.Warnings, ValidationError{
			Field:   "cache-folder",
			Message: "缓存目录未设置（将使用默认值）",
		})
	}
	
	if cfg.BatchSize < 0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "batch-size",
			Message: "batch-size 不能为负数",
		})
	}
}

// validateDownloadPerformance 验证下载性能配置
func validateDownloadPerformance(cfg *structs.ConfigSet, result *ValidationResult) {
	threads := map[string]int{
		"chunk_downloadthreads":    cfg.ChunkDownloadThreads,
		"lossless_downloadthreads": cfg.LosslessDownloadThreads,
		"aac_downloadthreads":      cfg.AacDownloadThreads,
		"hires_downloadthreads":    cfg.HiresDownloadThreads,
	}
	
	for field, value := range threads {
		if value <= 0 {
			result.Errors = append(result.Errors, ValidationError{
				Field:   field,
				Message: "线程数必须大于 0",
			})
		} else if value > 50 {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   field,
				Message: fmt.Sprintf("线程数过高（%d）可能导致性能问题或被限流", value),
			})
		}
	}
	
	// 验证缓冲区大小
	if cfg.BufferSizeKB <= 0 {
		result.Warnings = append(result.Warnings, ValidationError{
			Field:   "BufferSizeKB",
			Message: "缓冲区大小未设置或无效",
		})
	}
	
	if cfg.NetworkReadBufferKB <= 0 {
		result.Warnings = append(result.Warnings, ValidationError{
			Field:   "NetworkReadBufferKB",
			Message: "网络读取缓冲区大小未设置或无效",
		})
	}
}

// validateWorkRest 验证工作-休息循环配置
func validateWorkRest(cfg *structs.ConfigSet, result *ValidationResult) {
	if !cfg.WorkRestEnabled {
		return
	}
	
	if cfg.WorkDurationMinutes < 1 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "work-duration-minutes",
			Message: "工作时长至少为 1 分钟",
		})
	}
	
	if cfg.RestDurationMinutes < 1 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "rest-duration-minutes",
			Message: "休息时长至少为 1 分钟",
		})
	}
	
	if cfg.WorkDurationMinutes > 180 {
		result.Warnings = append(result.Warnings, ValidationError{
			Field:   "work-duration-minutes",
			Message: fmt.Sprintf("工作时长过长（%d 分钟 > 3 小时）", cfg.WorkDurationMinutes),
		})
	}
}

// validateAudioQuality 验证音质配置
func validateAudioQuality(cfg *structs.ConfigSet, result *ValidationResult) {
	// 验证 AAC 类型
	if cfg.AacType != "" {
		valid := false
		for _, aacType := range constants.ValidAacTypes {
			if cfg.AacType == aacType {
				valid = true
				break
			}
		}
		if !valid {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   "aac-type",
				Message: fmt.Sprintf("不支持的 AAC 类型 '%s'（有效值: %v）", cfg.AacType, constants.ValidAacTypes),
			})
		}
	}
	
	// 验证封面尺寸格式
	if cfg.CoverSize != "" {
		// 简单验证格式是否为 WxH
		if !strings.Contains(cfg.CoverSize, "x") {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   "cover-size",
				Message: fmt.Sprintf("封面尺寸格式可能不正确: '%s'（应为 WxH 格式，如 3000x3000）", cfg.CoverSize),
			})
		}
	}
}

// validateMVConfig 验证 MV 配置
func validateMVConfig(cfg *structs.ConfigSet, result *ValidationResult) {
	// 验证 MV 音频类型
	if cfg.MVAudioType != "" {
		validTypes := []string{"aac", "ac3", "ec3"}
		valid := false
		for _, t := range validTypes {
			if cfg.MVAudioType == t {
				valid = true
				break
			}
		}
		if !valid {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   "mv-audio-type",
				Message: fmt.Sprintf("MV 音频类型 '%s' 可能不支持（建议值: %v）", cfg.MVAudioType, validTypes),
			})
		}
	}
	
	// 验证 MV 最大分辨率
	if cfg.MVMax != 0 {
		valid := false
		for _, res := range constants.ValidMVResolutions {
			if cfg.MVMax == res {
				valid = true
				break
			}
		}
		if !valid {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   "mv-max",
				Message: fmt.Sprintf("MV 分辨率上限 %d 不在标准值中（有效值: %v）", cfg.MVMax, constants.ValidMVResolutions),
			})
		}
	}
	
	// 验证 MV 最小分辨率
	if cfg.MVMin != 0 {
		valid := false
		for _, res := range constants.ValidMVResolutions {
			if cfg.MVMin == res {
				valid = true
				break
			}
		}
		if !valid {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   "mv-min",
				Message: fmt.Sprintf("MV 分辨率下限 %d 不在标准值中（有效值: %v）", cfg.MVMin, constants.ValidMVResolutions),
			})
		}
	}
	
	// 验证上限和下限的逻辑关系
	if cfg.MVMax > 0 && cfg.MVMin > 0 && cfg.MVMin > cfg.MVMax {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "mv-min/mv-max",
			Message: fmt.Sprintf("MV 分辨率下限 (%d) 不能大于上限 (%d)", cfg.MVMin, cfg.MVMax),
		})
	}
}

// validatePathLimits 验证路径限制
func validatePathLimits(cfg *structs.ConfigSet, result *ValidationResult) {
	if cfg.LimitMax > 0 {
		if cfg.LimitMax < 50 {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   "limit-max",
				Message: fmt.Sprintf("路径长度限制过小（%d），可能导致文件名被过度截断", cfg.LimitMax),
			})
		}
	}
	
	if cfg.MaxPathLength > 0 {
		if cfg.MaxPathLength < 100 {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   "max-path-length",
				Message: fmt.Sprintf("最大路径长度过小（%d）", cfg.MaxPathLength),
			})
		}
	}
}

// validateLogging 验证日志配置
func validateLogging(cfg *structs.ConfigSet, result *ValidationResult) {
	if cfg.Logging.Level != "" {
		validLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
		levelValid := false
		for _, level := range validLevels {
			if strings.EqualFold(cfg.Logging.Level, level) {
				levelValid = true
				break
			}
		}
		
		if !levelValid {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   "logging.level",
				Message: fmt.Sprintf("无效的日志级别 '%s'（有效值: %v）", cfg.Logging.Level, validLevels),
			})
		}
	}
	
	// 验证日志输出配置
	if cfg.Logging.Output != "" {
		validOutputs := []string{"stdout", "stderr", "file"}
		outputValid := false
		for _, output := range validOutputs {
			if strings.EqualFold(cfg.Logging.Output, output) || strings.HasSuffix(cfg.Logging.Output, ".log") {
				outputValid = true
				break
			}
		}
		
		if !outputValid {
			result.Warnings = append(result.Warnings, ValidationError{
				Field:   "logging.output",
				Message: fmt.Sprintf("日志输出配置可能不正确: '%s'", cfg.Logging.Output),
			})
		}
	}
}

