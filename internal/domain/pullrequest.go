package domain

import (
	"avito-tech-go-task/internal/infrastructure/http/model"
	"errors"
	"math/rand"
	"time"
)

const (
	PRStatusOpen      PRStatus = "OPEN"
	PRStatusMerged    PRStatus = "MERGED"
	ReviewersMaxCount int64    = 2
)

var (
	ErrPRMerged            = errors.New("cannot modify ReviewersIDs for merged PR")
	ErrNoCandidate         = errors.New("no active replacement candidate in team")
	ErrReviewerNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrPRExists            = errors.New("PR already exists")
	ErrPRNotFound          = errors.New("PR not found")
)

type PRStatus string

func (s PRStatus) String() string {
	return string(s)
}

type PullRequest struct {
	ID           string
	Name         string
	AuthorID     string
	Status       PRStatus
	ReviewersIDs []string
	MergedAt     time.Time
}

func NewPullRequest(prID, name, authorID string, reviewersIDs []string) (*PullRequest, error) {
	if authorID == "" {
		return nil, errors.New("authorID ID is required")
	}

	return &PullRequest{
		ID:           prID,
		Name:         name,
		AuthorID:     authorID,
		Status:       PRStatusOpen,
		ReviewersIDs: reviewersIDs,
	}, nil
}

func NewPullRequestFromStorage(prID, name, authorID string, status PRStatus, reviewersIDs []string, mergedAt time.Time) PullRequest {
	return PullRequest{
		ID:           prID,
		Name:         name,
		AuthorID:     authorID,
		Status:       status,
		ReviewersIDs: reviewersIDs,
		MergedAt:     mergedAt,
	}
}

func (pr *PullRequest) IsOpen() bool {
	return pr.Status == PRStatusOpen
}

func (pr *PullRequest) IsMerged() bool {
	return pr.Status == PRStatusMerged
}

func (pr *PullRequest) SetMergedStatus() {
	pr.Status = PRStatusMerged
	pr.MergedAt = time.Now()
}

func (pr *PullRequest) GetReviewerIndex(reviewerID string) (index int64, exist bool) {
	for i, id := range pr.ReviewersIDs {
		if id == reviewerID {
			return int64(i), true
		}
	}
	return -1, false
}

func (pr *PullRequest) ReassignReviewer(oldReviewerIndex int64, candidatesForReview []string) (string, error) {
	if pr.IsMerged() {
		return "", ErrPRMerged
	}

	if len(candidatesForReview) == 0 {
		return "", ErrNoCandidate
	}

	randomInt := rand.Intn(len(candidatesForReview))
	randomReviewerID := candidatesForReview[randomInt]

	pr.ReviewersIDs[oldReviewerIndex] = randomReviewerID

	return randomReviewerID, nil
}

func (pr *PullRequest) ToJSON() model.PullRequest {
	return model.PullRequest{
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status.String(),
		AssignedReviewers: pr.ReviewersIDs,
		MergedAt:          pr.MergedAt,
	}
}

func (pr *PullRequest) ToJSONShort() model.PullRequestShort {
	return model.PullRequestShort{
		PullRequestID:   pr.ID,
		PullRequestName: pr.Name,
		AuthorID:        pr.AuthorID,
		Status:          pr.Status.String(),
	}
}
