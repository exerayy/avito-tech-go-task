package model

import "time"

type User struct {
	UserID   string `json:"user_id" example:"u2"`
	Username string `json:"username" example:"Bob"`
	TeamName string `json:"team_name" example:"backend"`
	IsActive bool   `json:"is_active" example:"false"`
}

type SetIsActiveUserRequest struct {
	UserID   string `json:"user_id" binding:"required" example:"u2"`
	IsActive bool   `json:"is_active" binding:"required" example:"false"`
}

type SetIsActiveUserResponse struct {
	User User `json:"user"`
}

type GetReviewUserResponse struct {
	UserID       string             `json:"user_id" example:"u2"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

type UserStat struct {
	UserID        string    `json:"user_id" example:"u2"`
	TotalReviews  int64     `json:"total_reviews" example:"2"`
	ActiveReviews int64     `json:"active_reviews" example:"1"`
	MergedReviews int64     `json:"merged_reviews" example:"1"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GetStatsResponse struct {
	UserStats []UserStat `json:"user_review_stats"`
}
