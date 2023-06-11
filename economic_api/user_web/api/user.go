package api

import (
	"context"
	"economic_api/user_web/forms"
	"economic_api/user_web/global"
	"economic_api/user_web/global/response"
	"economic_api/user_web/middlewares"
	"economic_api/user_web/models"
	"economic_api/user_web/proto"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"time"
)

func HandleGrpcErrorToHttp(err error, ctx *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, gin.H{
					"message": e.Message(),
				})
			case codes.Internal:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "go server internal error",
				})
			case codes.InvalidArgument:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "params error",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "oh no! what happened?",
				})
			}
			return
		}
	}
}

// api: query some users.
func GetUserList(ctx *gin.Context) {
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserServiceInfo.Host,
		global.ServerConfig.UserServiceInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] rpc dial error: ",
			"msg", err.Error())
	}

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomJwtClaims)
	zap.S().Infof("current user: %d", currentUser.ID)

	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	userClientService := proto.NewUserClient(userConn)
	resp, err := userClientService.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] rpc query user list error: ")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range resp.Data {
		user := response.UserResponse{
			Id:       value.Id,
			Mobile:   value.Mobile,
			NickName: value.NickName,
			Birthday: time.Time(time.Unix(int64(value.Birthday), 0)).Format("1998-05-05"),
			Gender:   value.Gender,
		}
		//data := make(map[string]interface{})
		//data["id"] = value.Id
		//data["name"] = value.NickName
		//data["gender"] = value.Gender
		//data["birthday"] = value.Birthday
		//data["mobile"] = value.Mobile
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}

// api: login by password.
func LoginByPassword(ctx *gin.Context) {
	loginForm := forms.LoginPasswordForm{}
	if err := ctx.ShouldBindJSON(&loginForm); err != nil {
		// TODO: form validate before delivery to user service.
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "params error:" + err.Error(),
		})
		return
	}

	// captcha verify
	if !store.Verify(loginForm.CaptchaId, loginForm.Captcha, true) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "captcha error, please input right captcha",
		})
		return
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserServiceInfo.Host,
		global.ServerConfig.UserServiceInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] rpc dial error: ",
			"msg", err.Error())
	}
	userClientService := proto.NewUserClient(userConn)
	if resp, err := userClientService.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: loginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "user don't exists",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "login error: " + err.Error(),
				})
			}
			return
		}
	} else {
		if passwordResp, passwordErr := userClientService.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          loginForm.Password,
			EncryptedPassword: resp.Password,
		}); passwordErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "login error, maybe your network error",
			})
		} else {
			if passwordResp.Success {
				// generate token.
				jwtAuth := middlewares.NewJWT()
				claims := models.CustomJwtClaims{
					ID:       uint(resp.Id),
					NickName: resp.NickName,
					RoleId:   uint(resp.Role),
					RegisteredClaims: jwt.RegisteredClaims{
						Issuer:    "economic_bookkeeping",
						NotBefore: jwt.NewNumericDate(time.Now()),
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
					},
				}
				token, err := jwtAuth.CreateToken(claims)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"message": "generate token error!",
					})
					return
				}
				ctx.JSON(http.StatusOK, gin.H{
					"id":         resp.Id,
					"nick_name":  resp.NickName,
					"token":      token,
					"expired_at": time.Now().Add(7 * 24 * time.Hour).Unix(),
				})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "login failed, please fill right information",
				})
			}

		}

	}
}

// api: user register
func RegisterUser(ctx *gin.Context) {
	registerForm := forms.RegisterForm{}
	if err := ctx.ShouldBindJSON(&registerForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "params error:" + err.Error(),
		})
		return
	}

	// check sms code right?
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: global.ServerConfig.RedisInfo.Password, // no password set
		DB:       0,                                      // use default DB
	})
	redisVal, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "please retrive sms code by your phone",
		})
		return
	} else {
		if redisVal != registerForm.SMSCode {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "sms code don't match your phone",
			})
			return
		}
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserServiceInfo.Host,
		global.ServerConfig.UserServiceInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] rpc dial error: ",
			"msg", err.Error())
	}
	userClientService := proto.NewUserClient(userConn)

	user, err := userClientService.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		Password: registerForm.Password,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "register user error, please concat developer." + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":        user.Id,
		"mobile":    user.Mobile,
		"nick_name": user.NickName,
	})
}
