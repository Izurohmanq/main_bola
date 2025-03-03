package middlewares

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"user-service/common/response"
	"user-service/config"
	"user-service/constants"
	errContstant "user-service/constants/error"
	servicesUser "user-service/services/user"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// jadi handle panic ini berguna untuk handle error panic
// error panic ini akan men shutdown atau mematikan secara paksa
// jadi ketika ada error maka aplikasi akan mati secara otomatis
// untuk handlenya buat sebuah recover, ketika ada error panic terus kita buat recovery ini
// aplikasinya tidak akan mati dan masih bisa dijalankan
func HandlePanic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Recovered from panic: %w", r)
				ctx.JSON(http.StatusInternalServerError, response.Response{
					Status:  constants.Error,
					Message: errContstant.ErrInternalServerError.Error(),
				})
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}

// fungsi sebagai batasan request yang masuk ke sistem
func RateLimitter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, ctx.Writer, ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusTooManyRequests, response.Response{
				Status:  constants.Error,
				Message: errContstant.ErrTooManyRequest.Error(),
			})
			ctx.Abort()
		}
		ctx.Next()
	}
}

func extractBearerToken(token string) string {
	arrayToken := strings.Split(token, " ")
	if len(arrayToken) == 2 {
		return arrayToken[1]
	}
	return ""
}

func responseUnauthorized(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, response.Response{
		Status:  constants.Error,
		Message: message,
	})
	ctx.Abort()
}

func validateAPIKey(ctx *gin.Context) error {
	apiKey := ctx.GetHeader(constants.XApiKey)
	requestAt := ctx.GetHeader(constants.XRequestAt)
	serviceName := ctx.GetHeader(constants.XServiceName)
	signatureKey := config.Config.SignatureKey

	validateKey := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)

	hash := sha256.New()
	hash.Write([]byte(validateKey))
	resultHash := hex.EncodeToString(hash.Sum(nil))

	if apiKey != resultHash {
		return errContstant.ErrUnauthorized
	}

	return nil
}

func validateBearerToken(ctx *gin.Context, token string) error {
	if !strings.Contains(token, "Bearer") {
		return errContstant.ErrUnauthorized
	}

	tokenString := extractBearerToken(token)
	if token == "" {
		return errContstant.ErrUnauthorized
	}

	claims := &servicesUser.Claims{}
	tokenJwt, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errContstant.ErrInvalidToken
		}

		jwtSecret := []byte(config.Config.JwtSecretKey)
		return jwtSecret, nil
	})

	if err != nil || !tokenJwt.Valid {
		return errContstant.ErrUnauthorized
	}

	userLogin := ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), constants.UserLogin, claims.User))

	ctx.Request = userLogin
	ctx.Set(constants.Token, token)
	return nil
}

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		token := ctx.GetHeader(constants.Authorization)
		if token == "" {
			responseUnauthorized(ctx, errContstant.ErrUnauthorized.Error())
			return
		}

		err = validateBearerToken(ctx, token)
		if err != nil {
			responseUnauthorized(ctx, err.Error())
			return
		}

		err = validateAPIKey(ctx)
		if err != nil {
			responseUnauthorized(ctx, err.Error())
			return
		}

		ctx.Next()
	}
}
