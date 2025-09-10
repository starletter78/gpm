package core

import (
	"context"
	"fmt"
	"gpm/global"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func initFile(logPath, appName string) {
	fileName := time.Now().Format("2006010215")
	err := os.MkdirAll(fmt.Sprintf("%s", logPath), os.ModePerm)
	if err != nil {
		logrus.Error(err)
		return
	}

	// 按小时创建日志文件
	filename := fmt.Sprintf("%s/%s.%s.log", logPath, appName, fileName)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		logrus.Error(err)
		return
	}
	fileHook := FileDateHook{file, logPath, filename, appName}
	logrus.AddHook(&fileHook)
}

func InitLogrus() {
	//新建一个实例
	logrus.SetReportCaller(false) //开启返回函数名和行号
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   false,
	})
	logrus.SetLevel(logrus.DebugLevel)
	initFile(global.Config.Log.Dir, global.Config.Log.App) //设置最低的Level
	logrus.Infof("日志初始化成功")
}

type FileDateHook struct {
	file     *os.File
	logPath  string
	fileName string // 小时，用于判断是否切换文件
	appName  string
}

func (hook *FileDateHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
func (hook *FileDateHook) Fire(entry *logrus.Entry) error {
	entry.Data["app"] = global.Config.Log.App
	if entry.Data["type"] == "" {
		entry.Data["type"] = "system"
	} else {
		for k, v := range hook.getFiles(entry.Context) {
			entry.Data[k] = v
		}
	}
	currentDate := entry.Time.Format("2006010215")
	line, _ := entry.String()
	if hook.fileName == currentDate {
		_, err := hook.file.Write([]byte(line))
		if err != nil {
			return err
		}
		return nil
	}
	// 时间或小时变化，关闭旧文件，创建新文件
	err := hook.file.Close()
	if err != nil {
		return err
	}
	// 创建新文件（按小时）
	filename := fmt.Sprintf("%s/%s.%s.log", hook.logPath, hook.appName, currentDate)
	hook.file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	_, err = hook.file.Write([]byte(line))
	return err
}

func (hook *FileDateHook) getFiles(ctx context.Context) logrus.Fields {
	fields := logrus.Fields{}
	fields["logId"] = ""
	fields["userId"] = "0"
	if ctx != nil {
		if logId, ok := ctx.Value("logId").(string); ok {
			fields["logId"] = logId
		}
		if userId, ok := ctx.Value("userId").(uint); ok {
			fields["userId"] = strconv.Itoa(int(userId))
		}
	}
	return fields
}
