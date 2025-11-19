package service

import (
	"avito-tech-go-task/internal/domain"
	"context"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go

type TeamRepository interface {
	Save(ctx context.Context, team domain.Team, teamMembers []domain.User) (err error)
	FindByName(ctx context.Context, teamName string) ([]domain.User, error)
	DeactivateTeam(ctx context.Context, teamName string) ([]domain.PullRequest, error)
}

type PullRequestRepository interface {
	CreatePR(ctx context.Context, pr domain.PullRequest) error
	MergePR(ctx context.Context, pr domain.PullRequest) error
	ReassignPR(ctx context.Context, pr domain.PullRequest, oldReviewer, newReviewer string) error
	FindByID(ctx context.Context, prID string) (domain.PullRequest, error)
	FindByReviewerID(ctx context.Context, reviewerID string) ([]domain.PullRequest, error)
}

type UserRepository interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (err error)
	FindByID(ctx context.Context, userID string) (domain.User, error)
	FindTeamByUserID(ctx context.Context, userID string) (string, error)
	FindActiveUserIDsByTeam(ctx context.Context, team string) ([]string, error)
	FindActiveUserIDsByTeamExcludeAuthor(ctx context.Context, team, excludeAuthorID string, reviewersCount int64) ([]string, error)
	GetStats(ctx context.Context, limit uint64) ([]domain.UserStat, error)
}
