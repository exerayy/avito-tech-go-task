package storage

import (
	"avito-tech-go-task/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type TeamRepo struct {
	db DB
}

type Team struct {
	name string `db:"name"`
}

func NewTeamRepo(db DB) *TeamRepo {
	return &TeamRepo{db: db}
}

func (t Team) toDomain() domain.Team {
	return *domain.NewTeam(t.name)
}

func (r *TeamRepo) Save(ctx context.Context, team domain.Team, teamMembers []domain.User) (err error) {
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

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO teams (name) VALUES ($1) ON CONFLICT (name) DO NOTHING",
		team.Name,
	)
	if err != nil {
		return fmt.Errorf("save team tx.ExecContext: %w", err)
	}

	builder := sq.Insert("users").
		Columns("id", "name", "team_name", "is_active").
		PlaceholderFormat(sq.Dollar).
		Suffix(`ON CONFLICT (id) DO UPDATE SET 
            name = EXCLUDED.name,
            team_name = EXCLUDED.team_name,
            is_active = EXCLUDED.is_active`)

	for _, member := range teamMembers {
		builder = builder.Values(member.ID, member.Name, team.Name, member.IsActive)
	}
	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("save team builder.ToSql: %w", err)
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("insert users tx.ExecContext: %w", err)
	}

	err = createReviewStats(ctx, tx, teamMembers)
	if err != nil {
		return fmt.Errorf("createStats: %w", err)
	}

	return nil
}

func createReviewStats(ctx context.Context, tx *sql.Tx, users []domain.User) error {
	builder := sq.Insert("user_review_stats").
		Columns("user_id", "updated_at").
		PlaceholderFormat(sq.Dollar).
		Suffix(`ON CONFLICT (user_id) DO NOTHING`)

	for _, user := range users {
		builder = builder.Values(user.ID, time.Now())
	}
	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("create stats builder.ToSql: %w", err)
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("create stats users tx.ExecContext: %w", err)
	}

	return nil
}

func (r *TeamRepo) FindByName(ctx context.Context, teamName string) ([]domain.User, error) {
	builder := sq.Select("id", "name", "team_name", "is_active").
		From("users").
		Where(sq.Eq{"team_name": teamName}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("FindByName team builder.ToSql: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("FindByName team db.Query: %w", err)
	}
	defer rows.Close()

	users := make([]domain.User, 0, 20)
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.id,
			&user.name,
			&user.teamName,
			&user.isActive,
		); err != nil {
			return nil, fmt.Errorf("FindByName team rows.Next: %w", err)
		}
		users = append(users, user.toDomain())
	}

	return users, nil
}

func (r *TeamRepo) DeactivateTeam(ctx context.Context, teamName string) ([]domain.PullRequest, error) {
	rows, err := r.db.Query(ctx,
		`WITH updated_users AS (
			UPDATE users
			SET is_active = FALSE
			WHERE team_name = $1
			  AND is_active = TRUE
			RETURNING id
		),
		open_prs AS (
			SELECT pr.id, pr.reviewers_ids
			FROM pull_requests pr
			WHERE pr.status = 'OPEN'::pr_status
			  AND EXISTS (
				  SELECT 1
				  FROM updated_users uu
				  WHERE uu.id = ANY(pr.reviewers_ids)
			  )
		),
		new_reviewers AS (
			SELECT u.id
			FROM users u
			WHERE u.team_name != $1
			  AND u.is_active = TRUE
			ORDER BY RANDOM()
			LIMIT 1
		)
		UPDATE pull_requests pr
		SET reviewers_ids = ARRAY(SELECT id FROM new_reviewers)
		FROM open_prs op
		WHERE pr.id = op.id
		  AND EXISTS (SELECT 1 FROM new_reviewers)
		RETURNING pr.id, pr.name, pr.author_id, pr.status, pr.reviewers_ids, pr.merged_at;`,
		teamName,
	)
	if err != nil {
		return nil, fmt.Errorf("deactivate team db.Exec: %w", err)
	}

	defer rows.Close()

	prs := make([]domain.PullRequest, 0, 20)
	for rows.Next() {
		var pullRequest PullRequest
		if err := rows.Scan(
			&pullRequest.id,
			&pullRequest.name,
			&pullRequest.authorID,
			&pullRequest.status,
			&pullRequest.reviewersIDs,
			&pullRequest.mergedAt,
		); err != nil {
			return nil, fmt.Errorf("DeactivateTeam rows.Next: %w", err)
		}

		prs = append(prs, pullRequest.toDomain())
	}

	return prs, nil
}
