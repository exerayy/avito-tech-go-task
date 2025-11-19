package controller

import (
	"avito-tech-go-task/internal/domain"
	"avito-tech-go-task/internal/infrastructure/http/model"
	"context"
)

type PRService interface {
	CreatePR(ctx context.Context, prID, prName, authorID string) (domain.PullRequest, error)
	MergePR(ctx context.Context, prID string) (domain.PullRequest, error)
	ReassignPR(ctx context.Context, prID, oldReviewerID string) (prVal domain.PullRequest, newReviewerID string, err error)
	SetIsActiveUser(ctx context.Context, userID string, isActive bool) (domain.User, error)
	GetReviewUser(ctx context.Context, userID string) ([]domain.PullRequest, error)
	AddTeam(ctx context.Context, teamName string, members []model.TeamMember) error
	GetTeam(ctx context.Context, teamName string) ([]domain.User, error)
	GetStats(ctx context.Context, limit uint64) ([]domain.UserStat, error)
	DeactivateTeam(ctx context.Context, teamName string) ([]domain.PullRequest, error)
}

type ApiService struct {
	prService PRService
}

func NewApiService(prService PRService) *ApiService {
	return &ApiService{prService: prService}
}

func (s *ApiService) CreatePullRequest(ctx context.Context, req *model.CreatePullRequestRequest) (*model.CreatePullRequestResponse, error) {
	pr, err := s.prService.CreatePR(ctx, req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		return nil, err
	}

	res := &model.CreatePullRequestResponse{
		PR: pr.ToJSON(),
	}

	return res, nil
}

func (s *ApiService) MergePullRequest(ctx context.Context, req *model.MergePullRequestRequest) (*model.MergePullRequestResponse, error) {
	pr, err := s.prService.MergePR(ctx, req.PullRequestID)
	if err != nil {
		return nil, err
	}

	res := &model.MergePullRequestResponse{
		PR: pr.ToJSON(),
	}

	return res, nil
}

func (s *ApiService) ReassignPullRequest(ctx context.Context, req *model.ReassignPullRequestRequest) (*model.ReassignPullRequestResponse, error) {
	pr, replacedBy, err := s.prService.ReassignPR(ctx, req.PullRequestID, req.OldReviewerID)
	if err != nil {
		return nil, err
	}

	res := &model.ReassignPullRequestResponse{
		PR:         pr.ToJSON(),
		ReplacedBy: replacedBy,
	}

	return res, nil
}

func (s *ApiService) SetIsActiveUser(ctx context.Context, req *model.SetIsActiveUserRequest) (*model.SetIsActiveUserResponse, error) {
	user, err := s.prService.SetIsActiveUser(ctx, req.UserID, req.IsActive)
	if err != nil {
		return nil, err
	}

	res := &model.SetIsActiveUserResponse{
		User: user.ToJSON(),
	}

	return res, nil
}

func (s *ApiService) GetReviewerUser(ctx context.Context, userID string) (*model.GetReviewUserResponse, error) {
	prs, err := s.prService.GetReviewUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	jsonPRs := make([]model.PullRequestShort, 0, len(prs))
	for _, pr := range prs {
		jsonPRs = append(jsonPRs, pr.ToJSONShort())
	}

	res := &model.GetReviewUserResponse{
		UserID:       userID,
		PullRequests: jsonPRs,
	}

	return res, nil
}

func (s *ApiService) AddTeam(ctx context.Context, req *model.AddTeamRequest) (*model.Team, error) {
	err := s.prService.AddTeam(ctx, req.TeamName, req.Members)
	if err != nil {
		return nil, err
	}

	res := &model.Team{
		TeamName: req.TeamName,
		Members:  req.Members,
	}

	return res, nil
}

func (s *ApiService) GetTeam(ctx context.Context, teamName string) (*model.Team, error) {
	users, err := s.prService.GetTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}

	jsonMembers := make([]model.TeamMember, 0, len(users))
	for _, u := range users {
		jsonMembers = append(jsonMembers, u.ToJSONTeamMember())
	}

	res := &model.Team{
		TeamName: teamName,
		Members:  jsonMembers,
	}

	return res, nil
}

func (s *ApiService) GetStats(ctx context.Context, limit uint64) (*model.GetStatsResponse, error) {
	userStats, err := s.prService.GetStats(ctx, limit)
	if err != nil {
		return nil, err
	}

	jsonUserStats := make([]model.UserStat, 0, len(userStats))
	for _, u := range userStats {
		jsonUserStats = append(jsonUserStats, u.ToJSON())
	}

	res := &model.GetStatsResponse{
		UserStats: jsonUserStats,
	}

	return res, nil
}

func (s *ApiService) DeactivateTeam(ctx context.Context, teamName string) (*model.DeactivateTeamResponse, error) {
	users, err := s.prService.DeactivateTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}

	jsonPRs := make([]model.PullRequest, 0, len(users))
	for _, u := range users {
		jsonPRs = append(jsonPRs, u.ToJSON())
	}

	res := &model.DeactivateTeamResponse{
		PullRequests: jsonPRs,
	}

	return res, nil
}
