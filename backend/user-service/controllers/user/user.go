package user

import (
	"net/http"
	customerror "user-service/common/custom-error"
	"user-service/common/response"
	"user-service/domain/dto"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	service services.IServiceRegistry
}

type IUserController interface {
	Login(*gin.Context)
	Register(*gin.Context)
	Update(*gin.Context)
	GetUserLogin(*gin.Context)
	GetUserByUUID(*gin.Context)
}

func NewUserController(service services.IServiceRegistry) IUserController {
	return &UserController{
		service: service,
	}
}

func (uc *UserController) GetUserByUUID(c *gin.Context) {
	user, err := uc.service.GetUser().GetUserByUUID(c, c.Param("uuid"))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})

		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  c,
	})
}

func (uc *UserController) GetUserLogin(c *gin.Context) {
	user, err := uc.service.GetUser().GetUserLogin(c.Request.Context())
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})

		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  c,
	})
}

func (uc *UserController) Login(c *gin.Context) {
	req := &dto.LoginRequest{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})

		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(req)
	if err != nil {
		errMsg := http.StatusText(http.StatusUnprocessableEntity)
		errResp := customerror.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Err:     err,
			Message: &errMsg,
			Data:    errResp,
			Gin:     c,
		})

		return
	}

	res, err := uc.service.GetUser().Login(c, req)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})

		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code:  http.StatusOK,
		Data:  res.User,
		Token: &res.Token,
		Gin:   c,
	})
}

func (uc *UserController) Register(c *gin.Context) {
	req := &dto.RegisterRequest{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})

		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(req)
	if err != nil {
		errMsg := http.StatusText(http.StatusUnprocessableEntity)
		errResp := customerror.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Err:     err,
			Message: &errMsg,
			Data:    errResp,
			Gin:     c,
		})

		return
	}

	res, err := uc.service.GetUser().Register(c, req)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})

		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: res.User,
		Gin:  c,
	})
}

func (uc *UserController) Update(c *gin.Context) {
	req := &dto.UpdateRequest{}
	uuid := c.Param("uuid")
	err := c.ShouldBindJSON(req)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})

		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(req)
	if err != nil {
		errMsg := http.StatusText(http.StatusUnprocessableEntity)
		errResp := customerror.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusUnprocessableEntity,
			Err:     err,
			Message: &errMsg,
			Data:    errResp,
			Gin:     c,
		})

		return
	}

	user, err := uc.service.GetUser().Update(c, req, uuid)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})

		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: user,
		Gin:  c,
	})
}
