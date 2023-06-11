package api

import (
	"context"
	"economic_api/user_web/forms"
	"economic_api/user_web/global"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func GenerateSMSCode(width int) string {
	numberRange := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numberRange)
	rand.Seed(time.Now().UnixNano())

	var stringBuilder strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&stringBuilder, "%d", numberRange[rand.Intn(r)])
	}
	return stringBuilder.String()
}

func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func SendSMS(ctx *gin.Context) {
	smsForm := forms.SMSForm{}
	if err := ctx.ShouldBindJSON(&smsForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "params error:" + err.Error(),
		})
		return
	}

	if len(smsForm.Mobile) < 11 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "please input right phone",
		})
		return
	}

	// send sms by alibaba sms service
	client, err := CreateClient(tea.String(global.ServerConfig.AliSMSInfo.AccessKeyId), tea.String(global.ServerConfig.AliSMSInfo.AccessKeySecret))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "send sms error: " + err.Error(),
		})
		return
	}

	smsCodeStr := GenerateSMSCode(6)
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(smsForm.Mobile),
		SignName:      tea.String(global.ServerConfig.AliSMSInfo.SignName),     // 短信签名名称
		TemplateCode:  tea.String(global.ServerConfig.AliSMSInfo.TemplateCode), // 短信模板code
		TemplateParam: tea.String(smsCodeStr),
	}
	runtime := &util.RuntimeOptions{}
	_, err = client.SendSmsWithOptions(sendSmsRequest, runtime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "send sms error: " + err.Error(),
		})
		return
	}
	// save sms code. pair mobile to smscode by redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: global.ServerConfig.RedisInfo.Password, // no password set
		DB:       0,                                      // use default DB
	})
	err = rdb.Set(context.Background(), smsForm.Mobile, smsCodeStr, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "redis set error: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "has already sended sms!",
	})
}
