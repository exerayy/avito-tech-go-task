package controller

import (
	"avito-tech-go-task/internal/infrastructure/http/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddTeamHandler godoc
//
//	@Summary		Создать команду с участниками (создаёт/обновляет пользователей)
//	@Description
//	@Tags			Teams
//	@Accept			json
//	@Produce		json
//	@Param			request    body		model.AddTeamRequest	true	"team"
//	@Success		200	{object}	model.Team
//	@Failure		400	{object}	model.ErrorResponse
//	@Failure		404	{object}	model.ErrorResponse
//	@Failure		500	{object}	model.ErrorResponse
//	@Router			/teams/add [post]
func (s *ApiService) AddTeamHandler(ctx *gin.Context) {
	var req model.AddTeamRequest

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

	res, err := s.AddTeam(ctx, &req)
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

// GetTeamHandler godoc
//
//	@Summary		Получить команду с участниками
//	@Description
//	@Tags			Teams
//	@Accept			json
//	@Produce		json
//	@Param			team_name	query		string	true	"team_name"
//	@Success		200	{object}	model.Team
//	@Failure		400	{object}	model.ErrorResponse
//	@Failure		404	{object}	model.ErrorResponse
//	@Failure		500	{object}	model.ErrorResponse
//	@Router			/teams/get [get]
func (s *ApiService) GetTeamHandler(ctx *gin.Context) {
	teamName := ctx.Query("team_name")
	if teamName == "" {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: &model.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "team_name can't be empty",
			},
		})
		return
	}

	res, err := s.GetTeam(ctx, teamName)
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

// DeactivateTeamHandler godoc
//
//	@Summary		Деактивировать команду
//	@Description
//	@Tags			Teams
//	@Accept			json
//	@Produce		json
//	@Param			team_name	query		string	true	"team_name"
//	@Success		200	{object}	model.DeactivateTeamResponse
//	@Failure		400	{object}	model.ErrorResponse
//	@Failure		404	{object}	model.ErrorResponse
//	@Failure		500	{object}	model.ErrorResponse
//	@Router			/teams/deactivate [patch]
func (s *ApiService) DeactivateTeamHandler(ctx *gin.Context) {
	teamName := ctx.Query("team_name")
	if teamName == "" {
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: &model.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "team_name can't be empty",
			},
		})
		return
	}

	res, err := s.DeactivateTeam(ctx, teamName)
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
