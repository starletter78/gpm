package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gpm/common/res"
	"gpm/global"
	"io"

	"github.com/gin-gonic/gin"
)

func ArgsCheckMiddleware(c *gin.Context) {
	tenant := c.GetHeader("tenant")
	timestamp := c.GetHeader("timestamp")
	signature := c.GetHeader("signature")
	if tenant == "" || timestamp == "" || signature == "" {
		res.FailValid(c, "请求信息不全")
		c.Abort()
		return
	}
	//err := global.DB.WithContext(c.Request.Context()).Find(&model.Tenant{}, "id = ?", tenant).Error
	//fmt.Println("err:", err)
	//if err != nil {
	//	res.FailWithError(c, err)
	//	c.Abort()
	//	return
	//}

	prefix := global.Config.ArgsCheck.Prefix
	suffix := global.Config.ArgsCheck.Suffix
	url := c.Request.RequestURI

	var requestBody []byte
	if c.Request.Body != nil {
		requestBody, _ = c.GetRawData()
		c.Request.Body = io.NopCloser(bytes.NewReader(requestBody)) // 恢复 Body 供后续使用
	}
	// --- 核心修改点：使用“解析再封装”方法规范化请求体 ---
	if len(requestBody) != 0 {
		// 1. 定义一个变量来存储解析后的JSON数据
		var jsonData interface{}
		// 2. 使用 json.Unmarshal 解析原始请求体
		// 使用 json.NewDecoder 更适合处理流式数据，但这里 c.GetString() 已经是字符串，Unmarshal 更直接
		err := json.Unmarshal(requestBody, &jsonData)
		if err != nil {
			// 如果JSON解析失败，说明请求体格式不合法，应拒绝请求
			res.FailValid(c, "请求体格式不合法: "+err.Error())
			return
		}
		// 3. 使用 json.Marshal 将解析后的数据重新序列化
		// Marshal 会生成一个紧凑的、无多余空格、键已排序的JSON字符串
		requestBody, err = json.Marshal(jsonData)
		if err != nil {
			res.FailWithMsg(c, "服务器内部错误: 无法序列化JSON")
			return
		}
	}
	signatureStr := prefix + tenant + url + string(requestBody) + timestamp + suffix
	fmt.Println(signatureStr)
	hash := sha256.New()
	hash.Write([]byte(signatureStr))
	hashSum := hash.Sum(nil)
	signatureHash, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		res.FailWithMsg(c, "sign解析失败")
		c.Abort()
		return
	}
	if string(signatureHash) != string(hashSum) {
		res.FailValid(c, "参数校验失败")
		c.Abort()
		return
	}
}
