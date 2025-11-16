package domain

import (
	"avito-tech-go-task/internal/infrastructure/http/model"
	"errors"
	"time"
)

var ErrUserNotExist = errors.New("user not exist")

type User struct {
	ID       string
	Name     string
	TeamName string
	IsActive bool
}

type UserStat struct {
	UserID        string
	TotalReviews  int64
	ActiveReviews int64
	MergedReviews int64
	UpdatedAt     time.Time
}

func NewUser(id, name, teamName string, isActive bool) *User {
	return &User{
		ID:       id,
		Name:     name,
		TeamName: teamName,
		IsActive: isActive,
	}
}

func NewUserStat(userID string, totalReviews, activeReviews, mergedReviews int64, updatedAt time.Time) *UserStat {
	return &UserStat{
		UserID:        userID,
		TotalReviews:  totalReviews,
		ActiveReviews: activeReviews,
		MergedReviews: mergedReviews,
		UpdatedAt:     updatedAt,
	}
}

func (u *User) SetIsActive(isActive bool) {
	u.IsActive = isActive
}

func (u *User) GetTeamID() string {
	return u.TeamName
}

func (u *User) ToJSON() model.User {
	return model.User{
		UserID:   u.ID,
		Username: u.Name,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func (u *User) ToJSONTeamMember() model.TeamMember {
	return model.TeamMember{
		UserID:   u.ID,
		Username: u.Name,
		IsActive: u.IsActive,
	}
}

func (u *UserStat) ToJSON() model.UserStat {
	return model.UserStat{
		UserID:        u.UserID,
		TotalReviews:  u.TotalReviews,
		ActiveReviews: u.ActiveReviews,
		MergedReviews: u.MergedReviews,
		UpdatedAt:     u.UpdatedAt,
	}
}
