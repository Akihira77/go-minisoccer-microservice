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
	"user-service/constants/custom-error"
	"user-service/services/user"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Recovered from panic: %v", r)
				c.JSON(http.StatusInternalServerError, response.Response{
					Status:  constants.Error,
					Message: customerror.ErrInternalServer.Error(),
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}

func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, response.Response{
				Status:  constants.Error,
				Message: customerror.ErrTooManyRequest.Error(),
			})

			c.Abort()
		}

		c.Next()
	}
}

func extractBearerToken(token string) string {
	arrToken := strings.Split(token, " ")
	if len(arrToken) == 2 {
		return arrToken[1]
	}

	return ""
}

func responseUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, response.Response{
		Status:  constants.Error,
		Message: message,
	})
	c.Abort()
}

func validateAPIKey(c *gin.Context) error {
	apiKey := c.GetHeader(constants.XApiKey)
	requestAt := c.GetHeader(constants.XRequestAt)
	serviceName := c.GetHeader(constants.XServiceName)
	signatureKey := config.Config.SignatureKey

	validateKey := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)
	hash := sha256.New()
	hash.Write([]byte(validateKey))
	resultHash := hex.EncodeToString(hash.Sum(nil))

	if apiKey != resultHash {
		logrus.Infof("ResultHash APIKey: %v", resultHash)
		return customerror.ErrUnauthorized
	}

	return nil
}

func validateBearerToken(c *gin.Context, token string) error {
	if !strings.Contains(token, "Bearer") {
		logrus.Errorf("Token is invalid")
		return customerror.ErrUnauthorized
	}

	tokenStr := extractBearerToken(token)
	if tokenStr == "" {
		logrus.Errorf("Token is empty")
		return customerror.ErrUnauthorized
	}

	claims := &user.Claims{}
	tokenJwt, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			logrus.Errorf("Token is invalid JWT")
			return nil, customerror.ErrInvalidToken
		}

		return []byte(config.Config.JwtSecretKey), nil
	})
	if err != nil || !tokenJwt.Valid {
		logrus.Errorf("Parsing token error: %v", err)
		return customerror.ErrUnauthorized
	}

	userLogin := c.Request.WithContext(context.WithValue(c.Request.Context(), constants.UserLogin, claims.User))
	c.Request = userLogin
	c.Set(constants.Token, token)

	return nil
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		token := c.GetHeader(constants.Authorization)
		if token == "" {
			logrus.Errorf("Token is empty inside Authorization header")
			responseUnauthorized(c, customerror.ErrUnauthorized.Error())
			return
		}

		err = validateBearerToken(c, token)
		if err != nil {
			logrus.Errorf("Token is invalid bearer token: %v", err)
			responseUnauthorized(c, err.Error())
			return
		}

		err = validateAPIKey(c)
		if err != nil {
			logrus.Errorf("Validating API Key invalid: %v", err)
			responseUnauthorized(c, err.Error())
			return
		}

		c.Next()
	}
}
