package model

type TeamMember struct {
	UserID   string `json:"user_id" binding:"required" example:"u1"`
	Username string `json:"username" binding:"required" example:"Alice"`
	IsActive bool   `json:"is_active" binding:"required" example:"true"`
}

type Team struct {
	TeamName string       `json:"team_name" example:"payments"`
	Members  []TeamMember `json:"members"`
}

type AddTeamRequest struct {
	TeamName string       `json:"team_name" binding:"required" example:"payments"`
	Members  []TeamMember `json:"members" binding:"required"`
}

func (tm *TeamMember) Validate() bool {
	if tm == nil {
		return false
	}

	if tm.UserID == "" {
		return false
	}

	if tm.Username == "" {
		return false
	}

	return true
}

type DeactivateTeamResponse struct {
	PullRequests []PullRequest `json:"pull_requests"`
}
