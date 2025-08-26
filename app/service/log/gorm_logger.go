package log

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"time"
)

// GormLogger 实现 logger.Interface 接口
type GormLogger struct {
	SlowThreshold time.Duration
}

// NewGormLogger 创建一个新的 GormLogger 实例
func NewGormLogger() *GormLogger {
	return &GormLogger{
		// 设置慢查询阈值，例如 200ms
		SlowThreshold: 200 * time.Millisecond,
	}
}

// LogMode 设置日志模式
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	// 这里可以根据需要调整 logrus 的级别
	// 例如，如果是 Silent 模式，可以设置 logrus.SetLevel(logrus.PanicLevel)
	// 为了简单，我们这里不改变 logrus 的级别
	return &newLogger
}

// Info 打印信息日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logrus.WithContext(ctx).Infof(msg, data...)
}

// Warn 打印警告日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logrus.WithContext(ctx).Warnf(msg, data...)
}

// Error 打印错误日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logrus.WithContext(ctx).Errorf(msg, data...)
}

// Trace 打印SQL语句和耗时
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// 获取执行的 SQL 和影响的行数
	sql, rows := fc()
	// 计算耗时
	elapsed := time.Since(begin)
	elapsedMs := float64(elapsed.Nanoseconds()) / 1e6

	// 获取包含 logId 的 logrus.Fields
	fields := logrus.Fields{}
	fields["rows"] = rows
	fields["duration"] = fmt.Sprintf("%.5f", elapsedMs) // 保留3位小数
	fields["sql"] = sql
	fields["type"] = "db"
	if len(sql) > 1*1024 {
		fields["sql"] = ""
	}
	switch {
	case err != nil && !errors.Is(err, logger.ErrRecordNotFound):
		// 发生错误
		logrus.WithContext(ctx).WithFields(fields).Error("Database Error")
	case elapsed > l.SlowThreshold:
		// 慢查询
		logrus.WithContext(ctx).WithFields(fields).Warn("Database Slow Query")
	default:
		// 正常查询
		logrus.WithContext(ctx).WithFields(fields).Debug("Database Query")
	}
}
