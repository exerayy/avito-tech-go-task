package model

import "time"

type PullRequest struct {
	PullRequestID     string    `json:"pull_request_id" example:"pr-1001"`
	PullRequestName   string    `json:"pull_request_name" example:"Add search"`
	AuthorID          string    `json:"author_id" example:"u1"`
	Status            string    `json:"status" example:"OPEN"`
	AssignedReviewers []string  `json:"assigned_reviewers"`
	MergedAt          time.Time `json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id" example:"pr-1001"`
	PullRequestName string `json:"pull_request_name" example:"Add search"`
	AuthorID        string `json:"author_id" example:"u1"`
	Status          string `json:"status" example:"OPEN"`
}

type CreatePullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id" binding:"required" example:"pr-1001"`
	PullRequestName string `json:"pull_request_name" binding:"required" example:"Add search"`
	AuthorID        string `json:"author_id" binding:"required" example:"u1"`
}

type CreatePullRequestResponse struct {
	PR PullRequest `json:"pr"`
}

type MergePullRequestRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required" example:"pr-1001"`
}

type MergePullRequestResponse struct {
	PR PullRequest `json:"pr"`
}

type ReassignPullRequestRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required" example:"pr-1001"`
	OldReviewerID string `json:"old_reviewer_id" binding:"required" example:"u2"`
}

type ReassignPullRequestResponse struct {
	PR         PullRequest `json:"pr"`
	ReplacedBy string      `json:"replaced_by" example:"u5"`
}
