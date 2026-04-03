package action

import (
	"context"
	"strings"

	"github.com/Tesohh/isshues/app"
	db "github.com/Tesohh/isshues/db/generated"
)

func BulkInsertLabels(app *app.App, projectID int64, names []string) ([]db.Label, error) {
	ctx := context.Background()
	labels := []db.Label{}

	for _, name := range names {
		label, err := app.DB.InsertLabelBasic(ctx, db.InsertLabelBasicParams{
			Name:      strings.ToLower(name),
			ProjectID: projectID,
		})
		if err != nil {
			return nil, err
		}
		labels = append(labels, label)
	}

	return labels, nil
}
