package domain

import "errors"

var (
	ErrTeamMembersNotFound  = errors.New("team members not found")
	ErrTeamMemberIsNotValid = errors.New("team member is not valid")
)

type Team struct {
	Name string
}

func NewTeam(name string) *Team {
	return &Team{
		Name: name,
	}
}
