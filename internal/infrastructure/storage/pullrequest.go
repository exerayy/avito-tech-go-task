package storage

import (
	"avito-tech-go-task/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

type PRRepo struct {
	db DB
}

type PullRequest struct {
	id           string         `db:"id"`
	name         string         `db:"name"`
	authorID     string         `db:"author_id"`
	status       string         `db:"status"`
	reviewersIDs pq.StringArray `db:"reviewers_ids"`
	mergedAt     time.Time      `db:"merged_at"`
}

func NewPRRepo(db DB) *PRRepo {
	return &PRRepo{db: db}
}

func (pr PullRequest) toDomain() domain.PullRequest {
	return domain.NewPullRequestFromStorage(pr.id, pr.name, pr.authorID, domain.PRStatus(pr.status), pr.reviewersIDs, pr.mergedAt)
}

func updateReviewStats(ctx context.Context, tx *sql.Tx, status domain.PRStatus, reviewerIDs ...string) error {
	builder := sq.Update("user_review_stats").
		Set("updated_at", time.Now())

	switch status {
	case domain.PRStatusOpen:
		builder = builder.Set("total_reviews", sq.Expr("total_reviews + 1")).
			Set("active_reviews", sq.Expr("active_reviews + 1"))

	case domain.PRStatusMerged:
		builder = builder.Set("merged_reviews", sq.Expr("merged_reviews + 1")).
			Set("active_reviews", sq.Expr("active_reviews - 1"))

	default:
		return errors.New("invalid pull request status")
	}

	builder = builder.Where(sq.Eq{"user_id": reviewerIDs}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("updateReviewStats builder.ToSql: %w", err)
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("updateReviewStats tx.ExecContext: %w", err)
	}

	return nil
}

func (r *PRRepo) CreatePR(ctx context.Context, pr domain.PullRequest) error {
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

	builder := sq.Insert("pull_requests").
		Columns("id", "name", "author_id", "status", "reviewers_ids", "merged_at").
		Values(pr.ID, pr.Name, pr.AuthorID, pr.Status, pq.StringArray(pr.ReviewersIDs), pr.MergedAt).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("CreatePR builder.ToSql: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("CreatePR db.Exec: %w", err)
	}

	err = updateReviewStats(ctx, tx, domain.PRStatusOpen, pr.ReviewersIDs...)
	if err != nil {
		return fmt.Errorf("UpdateReviewStats: %w", err)
	}

	return nil
}

func (r *PRRepo) MergePR(ctx context.Context, pr domain.PullRequest) error {
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

	builder := sq.Update("pull_requests").
		Set("status", domain.PRStatusMerged.String()).
		Set("merged_at", pr.MergedAt).
		Where(sq.Eq{"id": pr.ID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("MergePR builder.ToSql: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("MergePR db.Exec: %w", err)
	}

	err = updateReviewStats(ctx, tx, domain.PRStatusMerged, pr.ReviewersIDs...)
	if err != nil {
		return fmt.Errorf("UpdateReviewStats: %w", err)
	}

	return nil
}

func (r *PRRepo) ReassignPR(ctx context.Context, pr domain.PullRequest, oldReviewer, newReviewer string) error {
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

	builder := sq.Update("pull_requests").
		Set("reviewers_ids", pq.StringArray(pr.ReviewersIDs)).
		Where(sq.Eq{"id": pr.ID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("ReassignPR builder.ToSql: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ReassignPR db.Exec: %w", err)
	}

	err = updateReviewStats(ctx, tx, domain.PRStatusOpen, oldReviewer, newReviewer)
	if err != nil {
		return fmt.Errorf("UpdateReviewStats: %w", err)
	}

	return nil
}

func (r *PRRepo) FindByID(ctx context.Context, prID string) (domain.PullRequest, error) {
	builder := sq.Select("id", "name", "author_id", "status", "reviewers_ids", "merged_at").
		From("pull_requests").
		Where(sq.Eq{"id": prID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("FindByID PR builder.ToSql: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("FindByID PR db.Query: %w", err)
	}
	defer rows.Close()

	domainPR := domain.PullRequest{}
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
			return domain.PullRequest{}, fmt.Errorf("FindByID PR rows.Next: %w", err)
		}
		domainPR = pullRequest.toDomain()
	}

	if domainPR.ID == "" {
		return domain.PullRequest{}, domain.ErrPRNotFound
	}

	return domainPR, nil
}

func (r *PRRepo) FindByReviewerID(ctx context.Context, reviewerID string) ([]domain.PullRequest, error) {
	queryString := `SELECT id, name, author_id, status, reviewers_ids, merged_at
		FROM pull_requests
		WHERE $1 = ANY(reviewers_ids)`

	rows, err := r.db.Query(ctx, queryString, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("FindByReviewerID r.db.Query: %w", err)
	}
	defer rows.Close()

	prs := make([]domain.PullRequest, 0, 10)
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
			return nil, fmt.Errorf("FindByReviewerID rows.Next: %w", err)
		}

		prs = append(prs, pullRequest.toDomain())
	}

	return prs, nil
}
