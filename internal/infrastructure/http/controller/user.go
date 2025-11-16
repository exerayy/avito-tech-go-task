package controller

import (
	"avito-tech-go-task/internal/infrastructure/http/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetIsActiveUserHandler godoc
//
//	@Summary		Установить флаг активности пользователя
//	@Description
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request body		model.SetIsActiveUserRequest	true	"user"
//	@Success		200	{object}	model.SetIsActiveUserResponse
//	@Failure		400	{object}	model.ErrorResponse
//	@Failure		404	{object}	model.ErrorResponse
//	@Failure		500	{object}	model.ErrorResponse
//	@Router			/users/setIsActive [post]
func (s *ApiService) SetIsActiveUserHandler(ctx *gin.Context) {
	var req model.SetIsActiveUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: &model.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	res, err := s.SetIsActiveUser(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: &model.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// GetReviewerUserHandler godoc
//
//	@Summary		Получить PR'ы, где пользователь назначен ревьювером
//	@Description
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			user_id	   query		string	true	"user_id"
//	@Success		200	{object}	model.GetReviewUserResponse
//	@Failure		400	{object}	model.ErrorResponse
//	@Failure		404	{object}	model.ErrorResponse
//	@Failure		500	{object}	model.ErrorResponse
//	@Router			/users/getReview [get]
func (s *ApiService) GetReviewerUserHandler(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: &model.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "user_id can't be empty",
			},
		})
		return
	}

	res, err := s.GetReviewerUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: &model.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// GetStatsHandler godoc
//
//	@Summary		Получить cтатистику пользователей по PR
//	@Description	limit = 0 - убирает лимит
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			limit	   query		string	true	"limit"
//	@Success		200	{object}	model.GetStatsResponse
//	@Failure		400	{object}	model.ErrorResponse
//	@Failure		404	{object}	model.ErrorResponse
//	@Failure		500	{object}	model.ErrorResponse
//	@Router			/users/getStatsHandler [get]
func (s *ApiService) GetStatsHandler(ctx *gin.Context) {
	limit := ctx.Query("limit")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: &model.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "limit can't be string",
			},
		})
	}

	if limitInt < 0 {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: &model.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "limit can't be < 0",
			},
		})
		return
	}

	res, err := s.GetStats(ctx, uint64(limitInt))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: &model.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
