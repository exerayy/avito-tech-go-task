package tests

import (
	"avito-tech-go-task/internal/domain"
	"avito-tech-go-task/internal/infrastructure/http/model"
	"context"
)

func (s *TestSuite) TestAddTeam() {
	tests := []struct {
		name    string
		request *model.AddTeamRequest
		wantErr bool
		setup   func()
	}{
		{
			name: "success - add team with members",
			request: &model.AddTeamRequest{
				TeamName: "payments",
				Members: []model.TeamMember{
					{
						UserID:   "u1",
						Username: "Alice",
						IsActive: true,
					},
					{
						UserID:   "u2",
						Username: "Bot",
						IsActive: true,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success - add team with members",
			request: &model.AddTeamRequest{
				TeamName: "backend",
				Members: []model.TeamMember{
					{
						UserID:   "u3",
						Username: "Paul",
						IsActive: true,
					},
					{
						UserID:   "u4",
						Username: "Anna",
						IsActive: true,
					},
					{
						UserID:   "u5",
						Username: "Nina",
						IsActive: true,
					},
					{
						UserID:   "u6",
						Username: "Jack",
						IsActive: false,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success - update team members",
			request: &model.AddTeamRequest{
				TeamName: "payments",
				Members: []model.TeamMember{
					{
						UserID:   "u2",
						Username: "Lili",
						IsActive: true,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "fail - not valid team members",
			request: &model.AddTeamRequest{
				TeamName: "payments",
				Members: []model.TeamMember{
					{
						UserID:   "u2",
						Username: "",
						IsActive: true,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()

			result, err := s.ApiService.AddTeam(ctx, tt.request)

			if tt.wantErr {
				s.Error(err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.NotNil(result)
				s.Equal(tt.request.TeamName, result.TeamName)
				s.ElementsMatch(tt.request.Members, result.Members)
			}
		})
	}
}

func (s *TestSuite) TestCreatePullRequest() {
	tests := []struct {
		name                   string
		request                *model.CreatePullRequestRequest
		responseReviewersCount int
		wantErr                bool
		setup                  func()
	}{
		{
			name: "success - create PR from payments team",
			request: &model.CreatePullRequestRequest{
				PullRequestID:   "pr-100",
				PullRequestName: "add some features",
				AuthorID:        "u1",
			},
			responseReviewersCount: 1,
			wantErr:                false,
		},
		{
			name: "success - create PR from backend team",
			request: &model.CreatePullRequestRequest{
				PullRequestID:   "pr-102",
				PullRequestName: "add some validations",
				AuthorID:        "u3",
			},
			responseReviewersCount: 2,
			wantErr:                false,
		},
		{
			name: "fail - PR with this ID already exists",
			request: &model.CreatePullRequestRequest{
				PullRequestID:   "pr-100",
				PullRequestName: "add some features",
				AuthorID:        "u1",
			},
			wantErr: true,
		},
		{
			name: "fail - author not exist",
			request: &model.CreatePullRequestRequest{
				PullRequestID:   "pr-159",
				PullRequestName: "add some validations",
				AuthorID:        "u404",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()

			result, err := s.ApiService.CreatePullRequest(ctx, tt.request)

			if tt.wantErr {
				s.Error(err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.NotNil(result)
				s.Equal(tt.request.PullRequestID, result.PR.PullRequestID)
				s.Equal(tt.request.PullRequestName, result.PR.PullRequestName)
				s.Equal(tt.request.AuthorID, result.PR.AuthorID)
				s.Equal(domain.PRStatusOpen.String(), result.PR.Status)
			}
		})
	}
}

func (s *TestSuite) TestMergePullRequest() {
	tests := []struct {
		name    string
		request *model.MergePullRequestRequest
		wantErr bool
		setup   func()
	}{
		{
			name: "success - merge PR from payments team",
			request: &model.MergePullRequestRequest{
				PullRequestID: "pr-100",
			},
			wantErr: false,
		},
		{
			name: "success - repeated merge PR from payments team (without changes)",
			request: &model.MergePullRequestRequest{
				PullRequestID: "pr-100",
			},
			wantErr: false,
		},
		{
			name: "fail - PR not exist",
			request: &model.MergePullRequestRequest{
				PullRequestID: "pr-404",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()

			result, err := s.ApiService.MergePullRequest(ctx, tt.request)

			if tt.wantErr {
				s.Error(err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.NotNil(result)
				s.Equal(tt.request.PullRequestID, result.PR.PullRequestID)
				s.Equal(domain.PRStatusMerged.String(), result.PR.Status)
			}
		})
	}
}

func (s *TestSuite) TestSetIsActiveUser() {
	tests := []struct {
		name    string
		request *model.SetIsActiveUserRequest
		wantErr bool
		setup   func()
	}{
		{
			name: "success - activate user",
			request: &model.SetIsActiveUserRequest{
				UserID:   "u6",
				IsActive: true,
			},
			wantErr: false,
		},
		{
			name: "success - deactivate user",
			request: &model.SetIsActiveUserRequest{
				UserID:   "u1",
				IsActive: false,
			},
			wantErr: false,
		},
		{
			name: "fail - user not exist",
			request: &model.SetIsActiveUserRequest{
				UserID:   "u404",
				IsActive: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()

			result, err := s.ApiService.SetIsActiveUser(ctx, tt.request)

			if tt.wantErr {
				s.Error(err)
				s.Nil(result)
			} else {
				s.NoError(err)
				s.NotNil(result)
				s.Equal(tt.request.UserID, result.User.UserID)
				s.Equal(tt.request.IsActive, result.User.IsActive)
			}
		})
	}
}
