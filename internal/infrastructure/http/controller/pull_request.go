package controller

import (
	"avito-tech-go-task/internal/infrastructure/http/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreatePullRequestHandler godoc
//
//	@Summary		Создать PR и автоматически назначить до 2 ревьюверов из команды автора
//	@Description
//	@Tags			PullRequests
//	@Accept			json
//	@Produce		json
//	@Param			request  body		model.CreatePullRequestRequest	true	"pull_request"
//	@Success		200	{object}	model.CreatePullRequestResponse
//	@Failure		400	{object}	model.ErrorResponse
//	@Failure		404	{object}	model.ErrorResponse
//	@Failure		500	{object}	model.ErrorResponse
//	@Router			/pullRequests/create [post]
func (s *ApiService) CreatePullRequestHandler(ctx *gin.Context) {
	var req model.CreatePullRequestRequest

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

	res, err := s.CreatePullRequest(ctx, &req)
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

// MergePullRequestHandler godoc
//
//	@Summary		Пометить PR как MERGED (идемпотентная операция)
//	@Description
//	@Tags			PullRequests
//	@Accept			json
//	@Produce		json
//	@Param			request  body		model.MergePullRequestRequest	true	"pull_request_id"
//	@Success		200	{object}	model.MergePullRequestResponse
//	@Failure		400	{object}	model.ErrorResponse
//	@Failure		404	{object}	model.ErrorResponse
//	@Failure		500	{object}	model.ErrorResponse
//	@Router			/pullRequests/merge [post]
func (s *ApiService) MergePullRequestHandler(ctx *gin.Context) {
	var req model.MergePullRequestRequest

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

	res, err := s.MergePullRequest(ctx, &req)
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

// ReassignPullRequestHandler godoc
//
//	@Summary		Переназначить конкретного ревьювера на другого из его команды
//	@Description
//	@Tags			PullRequests
//	@Accept			json
//	@Produce		json
//	@Param			request  body		model.ReassignPullRequestRequest	true	"pull_request"
//	@Success		200	{object}	model.ReassignPullRequestResponse
//	@Failure		400	{object}	model.ErrorResponse
//	@Failure		404	{object}	model.ErrorResponse
//	@Failure		500	{object}	model.ErrorResponse
//	@Router			/pullRequests/reassign [post]
func (s *ApiService) ReassignPullRequestHandler(ctx *gin.Context) {
	var req model.ReassignPullRequestRequest

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

	res, err := s.ReassignPullRequest(ctx, &req)
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
