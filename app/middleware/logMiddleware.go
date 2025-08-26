package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gpm/app/model"
	"gpm/global"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func LogMiddleware(c *gin.Context) {
	// 记录请求开始的时间
	startTime := time.Now()
	logId := uuid.New().String()
	c.Set("logId", logId)
	// 读取请求体（需在 c.Next() 前读取，因为请求体只能读一次）
	// 使用 context.WithValue 创建一个新的 context，并附加 logId
	ctx := context.WithValue(c.Request.Context(), "logId", logId)
	// 将新的 context 替换回请求中
	c.Request = c.Request.WithContext(ctx)
	var requestBody []byte
	if c.Request.Body != nil {
		requestBody, _ = c.GetRawData()
		c.Request.Body = io.NopCloser(bytes.NewReader(requestBody)) // 恢复 Body 供后续使用
	}
	// 创建一个 ResponseWriter 包装器以捕获响应体
	w := &responseBodyWriter{ResponseWriter: c.Writer, body: new(bytes.Buffer)}
	c.Writer = w

	// 执行后续处理器
	c.Next()
	requestBodyStr := string(requestBody)
	responseBodyStr := string(w.body.Bytes())
	responseBodyStrDB := responseBodyStr
	if len(requestBodyStr) > 1024*128 {
		responseBodyStrDB = "数据过大，请查看文件日志"
	}
	header, err := json.Marshal(c.Request.Header)
	if err != nil {
		return
	}
	headerStr := string(header)
	duration := time.Since(startTime).Seconds() * 1000
	// 构造操作日志结构体
	actionLog := model.ActionLog{
		LogID:        logId,
		UserID:       c.GetString("userID"),
		IP:           c.ClientIP(),
		UA:           c.Request.UserAgent(),
		Action:       c.GetString("action"),
		Path:         c.Request.URL.Path,
		Method:       c.Request.Method,
		Tenant:       c.GetString("tenant"),
		Header:       &headerStr,
		RequestBody:  &requestBodyStr,
		ResponseBody: &responseBodyStrDB,
		Status:       c.Writer.Status(),
		Duration:     fmt.Sprintf("%.5f", duration),
	}
	// 写入数据库
	if err = global.DB.WithContext(c.Request.Context()).Create(&actionLog).Error; err != nil {
		logrus.WithError(err).WithField("logId", logId).Error("Failed to save action log")
	}
	duration = time.Since(startTime).Seconds() * 1000
	// 使用 logrus 输出结构化日志
	logrus.WithContext(c.Request.Context()).WithFields(logrus.Fields{
		"type":     "action",
		"userId":   c.GetInt("action"),
		"ip":       c.ClientIP(),
		"method":   c.Request.Method,
		"path":     c.Request.URL.Path,
		"status":   c.Writer.Status(),
		"duration": fmt.Sprintf("%.5f", duration),
		"request":  string(requestBody),
		"response": responseBodyStr,
	}).Info("Request processed")
}

// responseBodyWriter 用于捕获响应体
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
