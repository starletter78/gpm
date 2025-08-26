package log

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LogByGin(c *gin.Context, fields logrus.Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}
