package controllers

import (
	"net/http"
	errValidation "order-service/common/error"
	"order-service/common/response"
	"order-service/domain/dto"
	"order-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OrderController struct {
	service services.IServiceRegistry
}

type IOrderController interface {
	GetAllWithPagination(*gin.Context)
	GetByUUID(*gin.Context)
	GetOrderByUserID(*gin.Context)
	Create(*gin.Context)
}

func NewOrderController(service services.IServiceRegistry) IOrderController {
	return &OrderController{service: service}
}

func (o *OrderController) GetAllWithPagination(c *gin.Context) {
	var params dto.OrderRequestParam
	err := c.ShouldBindQuery(&params)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})

		return
	}

	validate := validator.New()
	if err = validate.Struct(params); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     c,
		})

		return
	}

	result, err := o.service.GetOrder().GetAllWithPagination(c, &params)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})
		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (o *OrderController) GetByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	result, err := o.service.GetOrder().GetByUUID(c, uuid)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})

		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (o *OrderController) GetOrderByUserID(c *gin.Context) {
	result, err := o.service.GetOrder().GetOrderByUserID(c.Request.Context())
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})

		return
	}
	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (o *OrderController) Create(c *gin.Context) {
	var (
		request dto.OrderRequest
		ctx     = c.Request.Context()
	)

	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})

		return
	}

	validate := validator.New()
	if err = validate.Struct(request); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HTTPResponse(response.ParamHTTPResp{
			Err:     err,
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Gin:     c,
		})

		return
	}

	result, err := o.service.GetOrder().Create(ctx, &request)
	if err != nil {
		response.HTTPResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  c,
		})

		return
	}

	response.HTTPResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}
