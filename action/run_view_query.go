package action

import (
	"context"

	"github.com/Tesohh/isshues/app"
	db "github.com/Tesohh/isshues/db/generated"
)

// NOTE: changing the shape of Issue will require changing this too.
func RunViewQuery(app *app.App, tx db.DBTX, query string, binds []any) ([]db.Issue, error) {
	ctx := context.Background()
	rows, err := tx.Query(ctx, query, binds...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []db.Issue
	for rows.Next() {
		var i db.Issue
		if err := rows.Scan(
			&i.ID, &i.Title, &i.Code,
			&i.Description, &i.Status, &i.Priority,
			&i.ProjectID, &i.RecruiterUserID,
			&i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
