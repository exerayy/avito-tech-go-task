package service

import (
	"avito-tech-go-task/internal/domain"
	"avito-tech-go-task/internal/infrastructure/http/model"
	"avito-tech-go-task/pkg/helper"
	"context"
	"errors"
)

type PRService struct {
	prRepo   PullRequestRepository
	userRepo UserRepository
	teamRepo TeamRepository
}

func NewPRService(prRepo PullRequestRepository, userRepo UserRepository, teamRepo TeamRepository) *PRService {
	return &PRService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *PRService) CreatePR(ctx context.Context, prID, prName, authorID string) (domain.PullRequest, error) {
	_, err := s.prRepo.FindByID(ctx, prID)
	if !errors.Is(err, domain.ErrPRNotFound) {
		return domain.PullRequest{}, domain.ErrPRExists
	}

	teamID, err := s.userRepo.FindTeamByUserID(ctx, authorID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	reviewersIDs, err := s.userRepo.FindActiveUserIDsByTeamExcludeAuthor(ctx, teamID, authorID, domain.ReviewersMaxCount)
	if err != nil {
		return domain.PullRequest{}, err
	}

	pr, err := domain.NewPullRequest(prID, prName, authorID, reviewersIDs)
	if err != nil {
		return domain.PullRequest{}, err
	}

	err = s.prRepo.CreatePR(ctx, *pr)
	if err != nil {
		return domain.PullRequest{}, err
	}

	return *pr, nil
}

func (s *PRService) MergePR(ctx context.Context, prID string) (domain.PullRequest, error) {
	pr, err := s.prRepo.FindByID(ctx, prID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	if pr.IsMerged() {
		return pr, nil
	}

	pr.SetMergedStatus()

	err = s.prRepo.MergePR(ctx, pr)
	if err != nil {
		return domain.PullRequest{}, err
	}

	return pr, nil
}

func (s *PRService) ReassignPR(ctx context.Context, prID, oldReviewerID string) (prVal domain.PullRequest, newReviewerID string, err error) {
	pr, err := s.prRepo.FindByID(ctx, prID)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	oldReviewerIndexInPR, exist := pr.GetReviewerIndex(oldReviewerID)
	if !exist {
		return domain.PullRequest{}, "", domain.ErrReviewerNotAssigned
	}

	oldReviewerTeam, err := s.userRepo.FindTeamByUserID(ctx, oldReviewerID)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	activeCandidatesForReview, err := s.userRepo.FindActiveUserIDsByTeam(ctx, oldReviewerTeam)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	activeCandidatesForReview = helper.RemoveElement(activeCandidatesForReview, pr.AuthorID)
	for _, id := range pr.ReviewersIDs {
		activeCandidatesForReview = helper.RemoveElement(activeCandidatesForReview, id)
	}

	newReviewerID, err = pr.ReassignReviewer(oldReviewerIndexInPR, activeCandidatesForReview)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	err = s.prRepo.ReassignPR(ctx, pr, oldReviewerID, newReviewerID)
	if err != nil {
		return domain.PullRequest{}, "", err
	}

	return pr, newReviewerID, nil
}

func (s *PRService) SetIsActiveUser(ctx context.Context, userID string, isActive bool) (domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if errors.Is(err, domain.ErrUserNotExist) {
		return domain.User{}, err
	}

	err = s.userRepo.SetIsActive(ctx, userID, isActive)
	if err != nil {
		return domain.User{}, err
	}

	user.IsActive = isActive

	return user, nil
}

func (s *PRService) GetReviewUser(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	prs, err := s.prRepo.FindByReviewerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}

func (s *PRService) AddTeam(ctx context.Context, teamName string, members []model.TeamMember) error {
	for _, m := range members {
		if !m.Validate() {
			return domain.ErrTeamMemberIsNotValid
		}
	}

	team := domain.NewTeam(teamName)
	domainMembers := make([]domain.User, 0, len(members))
	for _, m := range members {
		domainMembers = append(domainMembers, *domain.NewUser(
			m.UserID,
			m.Username,
			teamName,
			m.IsActive,
		),
		)
	}

	err := s.teamRepo.Save(ctx, *team, domainMembers)
	if err != nil {
		return err
	}

	return nil
}

func (s *PRService) GetTeam(ctx context.Context, teamName string) ([]domain.User, error) {
	teamMembers, err := s.teamRepo.FindByName(ctx, teamName)
	if err != nil {
		return nil, err
	}

	if len(teamMembers) == 0 {
		return nil, domain.ErrTeamMembersNotFound
	}

	return teamMembers, nil
}

func (s *PRService) GetStats(ctx context.Context, limit uint64) ([]domain.UserStat, error) {
	userStats, err := s.userRepo.GetStats(ctx, limit)
	if err != nil {
		return nil, err
	}

	return userStats, nil
}
