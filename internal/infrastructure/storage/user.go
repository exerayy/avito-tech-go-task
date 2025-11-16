package storage

import (
	"avito-tech-go-task/internal/domain"
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type UserRepo struct {
	db DB
}

type User struct {
	id       string `db:"id"`
	name     string `db:"name"`
	teamName string `db:"team_name"`
	isActive bool   `db:"is_active"`
}

type UserStat struct {
	UserID        string    `db:"user_id"`
	TotalReviews  int64     `db:"total_reviews" `
	ActiveReviews int64     `db:"active_reviews"`
	MergedReviews int64     `db:"merged_reviews"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func NewUserRepo(db DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u User) toDomain() domain.User {
	return *domain.NewUser(u.id, u.name, u.teamName, u.isActive)
}

func (u UserStat) toDomain() domain.UserStat {
	return *domain.NewUserStat(u.UserID, u.TotalReviews, u.ActiveReviews, u.MergedReviews, u.UpdatedAt)
}

func (r *UserRepo) SetIsActive(ctx context.Context, userID string, isActive bool) (err error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("db.Begin: %w", err)
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
			if err != nil {
				err = fmt.Errorf("tx.Commit: %w", err)
			}
		}
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				err = fmt.Errorf("%w tx.Rollback: %s", err, rbErr)
			}
		}
	}()

	_, err = tx.ExecContext(ctx,
		`UPDATE users
		SET is_active = $1
		WHERE id = $2`,
		isActive,
		userID,
	)
	if err != nil {
		return fmt.Errorf("update is_active user tx.ExecContext: %w", err)
	}

	if isActive == true {
		return nil
	}

	// Удаляем неактивного ревьюера со всех PR со статусом OPEN
	_, err = tx.ExecContext(
		ctx,
		`UPDATE pull_requests
		SET reviewers_ids = array_remove(reviewers_ids, $1)
		WHERE $1 = ANY(reviewers_ids) AND status = $2`,
		userID,
		domain.PRStatusOpen,
	)
	if err != nil {
		return fmt.Errorf("remove not active reviewer tx.ExecContext: %w", err)
	}

	return nil
}

func (r *UserRepo) FindByID(ctx context.Context, userID string) (domain.User, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, team_name, is_active FROM users WHERE id = $1", userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("FindByID db.Query: %w", err)
	}
	defer rows.Close()

	domainUser := domain.User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.id,
			&user.name,
			&user.teamName,
			&user.isActive,
		); err != nil {
			return domain.User{}, fmt.Errorf("FindByID rows.Next: %w", err)
		}
		domainUser = user.toDomain()
	}

	if domainUser.ID == "" {
		return domain.User{}, domain.ErrUserNotExist
	}

	return domainUser, nil
}

func (r *UserRepo) FindTeamByUserID(ctx context.Context, userID string) (string, error) {
	rows, err := r.db.Query(ctx, "SELECT team_name FROM users WHERE id = $1", userID)
	if err != nil {
		return "", fmt.Errorf("FindTeamByUserID db.Query: %w", err)
	}
	defer rows.Close()

	var teamName string
	for rows.Next() {
		if err := rows.Scan(
			&teamName,
		); err != nil {
			return "", fmt.Errorf("FindTeamByUserID rows.Next: %w", err)
		}
	}

	if teamName == "" {
		return "", domain.ErrUserNotExist
	}

	return teamName, nil
}

func (r *UserRepo) FindActiveUserIDsByTeam(ctx context.Context, team string) ([]string, error) {
	rows, err := r.db.Query(ctx, "SELECT id FROM users WHERE team_name = $1", team)
	if err != nil {
		return nil, fmt.Errorf("FindActiveUserIDsByTeam db.Query: %w", err)
	}
	defer rows.Close()

	var userID string
	userIDs := make([]string, 0, 15)
	for rows.Next() {
		if err := rows.Scan(
			&userID,
		); err != nil {
			return nil, fmt.Errorf("FindActiveUserIDsByTeam rows.Next: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

func (r *UserRepo) FindActiveUserIDsByTeamExcludeAuthor(ctx context.Context, team, excludeAuthorID string, reviewersCount int64) ([]string, error) {
	rows, err := r.db.Query(ctx, "SELECT id FROM users WHERE team_name = $1 AND id != $2 LIMIT $3", team, excludeAuthorID, reviewersCount)
	if err != nil {
		return nil, fmt.Errorf("FindActiveUserIDsByTeamExcludeAuthor db.Query: %w", err)
	}
	defer rows.Close()

	var userID string
	userIDs := make([]string, 0, 15)
	for rows.Next() {
		if err := rows.Scan(
			&userID,
		); err != nil {
			return nil, fmt.Errorf("FindActiveUserIDsByTeamExcludeAuthor rows.Next: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

func (r *UserRepo) GetStats(ctx context.Context, limit uint64) ([]domain.UserStat, error) {
	builder := sq.Select("user_id", "total_reviews", "active_reviews", "merged_reviews", "updated_at ").
		From("user_review_stats").
		Limit(limit).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetStats team builder.ToSql: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetStats team db.Query: %w", err)
	}
	defer rows.Close()

	users := make([]domain.UserStat, 0, 20)
	for rows.Next() {
		var userStat UserStat
		if err := rows.Scan(
			&userStat.UserID,
			&userStat.TotalReviews,
			&userStat.ActiveReviews,
			&userStat.MergedReviews,
			&userStat.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("GetStats team rows.Next: %w", err)
		}
		users = append(users, userStat.toDomain())
	}

	return users, nil
}
