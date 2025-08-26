package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gpm/common/res"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// FileSearchReq 搜索请求参数（无分页）
type FileSearchReq struct {
	FilePath string `form:"filePath" binding:"required"` // 日志文件路径（相对路径）
	LogId    string `form:"logId"`                       // 日志 ID
	Type     string `form:"type"`                        // 日志类型：action, db, system
	Level    string `form:"level"`                       // 日志级别：info, error, debug, warn
	Keyword  string `form:"keyword"`                     // 关键词（msg、sql、request、response 中匹配）
}

// FileSearchView 搜索日志文件中的匹配行，返回所有符合条件的日志（无分页）
func (File) FileSearchView(c *gin.Context) {
	var cr FileSearchReq
	if err := c.ShouldBindQuery(&cr); err != nil {
		res.FailWithError(c, err)
		return
	}

	// 检查文件是否存在
	file, err := os.Open(cr.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			res.FailWithMsg(c, "日志文件不存在")
			return
		} else {
			res.FailWithMsg(c, "打开文件失败")
			return
		}
	}
	defer file.Close()

	// 验证是否为普通文件
	info, err := file.Stat()
	if err != nil || info.IsDir() {
		res.FailWithMsg(c, "无效的日志文件")
		return
	}

	var results []map[string]interface{}

	// 使用 bufio.Scanner 逐行读取
	scanner := bufio.NewScanner(file)
	// 增大缓冲区以支持超长日志行（如大 SQL 或 JSON）
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024) // 最大支持 1MB 的单行

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()

		// 快速关键词过滤（大小写不敏感）
		if cr.Keyword != "" {
			if !containsIgnoreCase(line, cr.Keyword) {
				continue
			}
		}

		// 解析 JSON
		var entry map[string]interface{}
		if err = json.Unmarshal(line, &entry); err != nil {
			logrus.WithContext(c.Request.Context()).Info("第 %d 行解析失败（跳过）: %s", lineNum, truncate(string(line), 200))
			continue
		}

		// 提取过滤字段
		logId, _ := entry["logId"].(string)
		logType, _ := entry["type"].(string)
		level, _ := entry["level"].(string)

		// 过滤条件（空值表示不限制）
		if cr.LogId != "" && logId != cr.LogId {
			continue
		}
		if cr.Type != "" && logType != cr.Type {
			continue
		}
		if cr.Level != "" && level != cr.Level {
			continue
		}

		// 匹配成功，转换为精简结构并加入结果
		results = append(results, entry)
	}

	// 检查扫描错误
	if err = scanner.Err(); err != nil {
		res.FailWithMsg(c, "读取文件时发生错误")
		return
	}

	res.SuccessWithList(c, results, int64(len(results)))
}

// getString 安全获取字符串
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
		return toString(v)
	}
	return ""
}

// toString 通用转字符串（避免 panic）
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%g", val)
	case int, int64:
		return fmt.Sprintf("%v", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", val)
	}
}

// containsIgnoreCase 大小写不敏感匹配（支持中文）
func containsIgnoreCase(data []byte, substr string) bool {
	return strings.Contains(
		strings.ToLower(string(data)),
		strings.ToLower(substr),
	)
}

// isSubPath 检查 path 是否在 root 目录下
func isSubPath(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..")
}

// truncate 字符串截断（避免日志过长）
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
