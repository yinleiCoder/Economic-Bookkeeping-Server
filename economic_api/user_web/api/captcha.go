package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(ctx *gin.Context) {
	driver := base64Captcha.DefaultDriverDigit
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := captcha.Generate()
	if err != nil {
		zap.S().Errorf("generate captcha error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "generate captcha error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picture":   b64s,
	})
}
