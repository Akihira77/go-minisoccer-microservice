package response

import (
	"net/http"
	"user-service/constants"
	errConstant "user-service/constants/custom-error"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Data    any     `json:"data"`
	Token   *string `json:"token,omitempty"`
}

type ParamHTTPResp struct {
	Code    int
	Err     error
	Message *string
	Gin     *gin.Context
	Data    any
	Token   *string
}

func HttpResponse(param ParamHTTPResp) {
	if param.Err == nil {
		param.Gin.JSON(param.Code, Response{
			Status:  constants.Success,
			Message: http.StatusText(http.StatusOK),
			Data:    param.Data,
			Token:   param.Token,
		})

		return
	}

	message := errConstant.ErrInternalServer.Error()
	if param.Message != nil {
		message = *param.Message
	} else if param.Err != nil && errConstant.ErrMapping(param.Err) {
		message = param.Err.Error()
	}

	param.Gin.JSON(param.Code, Response{
		Status:  constants.Error,
		Message: message,
		Data:    param.Data,
	})
}
