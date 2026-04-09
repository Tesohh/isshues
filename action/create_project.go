package action

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tesohh/isshues/app"
	"github.com/Tesohh/isshues/config"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var DuplicatePrefixErr = fmt.Errorf("prefix already exists")

// note: prefix must be 4 chars and uppercase already
func CreateProject(app *app.App, userId int64, title string, prefix string) error {
	ctx := context.Background()

	// create the project
	projectId, err := app.GetDB().InsertProject(ctx, db.InsertProjectParams{
		Title:  title,
		Prefix: prefix,
	})
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) && pgerr.Code == "23505" {
		return DuplicatePrefixErr
	} else if err != nil {
		return fmt.Errorf("project insertion error: %w", err)
	}

	// add default groups
	var defaultGroups []config.DefaultGroup
	err = app.Viper.UnmarshalKey("default_groups", &defaultGroups)
	if err != nil {
		return fmt.Errorf("default_groups config error: %w", err)
	}

	for _, group := range defaultGroups {
		params := db.InsertGroupParams{
			Name:        pgtype.Text{String: group.Name, Valid: true},
			ColorKey:    pgtype.Text{String: group.Color, Valid: group.Color != ""},
			Mentionable: group.Mentionable,
			ProjectID:   projectId,
		}
		groupId, err := app.GetDB().InsertGroup(ctx, params)
		if err != nil {
			return fmt.Errorf("group insertion error: %w", err)
		}

		for _, permission := range group.Permissions {
			err = app.GetDB().GrantPermissionToGroup(ctx, db.GrantPermissionToGroupParams{
				GroupID:             groupId,
				ProjectPermissionID: permission,
			})
			if err != nil {
				return fmt.Errorf("group grant permission error: %w", err)
			}
		}

		// add creator to the admins group
		if group.AddCreator {
			err = app.GetDB().AddUserToGroup(ctx, db.AddUserToGroupParams{
				GroupID: groupId,
				UserID:  userId,
			})
			if err != nil {
				return fmt.Errorf("cannot add creator to grouop: %w", err)
			}
		}
	}

	// add default groups
	var defaultLabels []config.DefaultLabel
	err = app.Viper.UnmarshalKey("default_labels", &defaultLabels)
	if err != nil {
		return fmt.Errorf("default_labels config error: %w", err)
	}

	for _, label := range defaultLabels {
		params := db.InsertLabelParams{
			Name:      label.Name,
			ColorKey:  pgtype.Text{String: label.Color, Valid: label.Color != ""},
			Symbol:    pgtype.Text{String: label.Symbol, Valid: label.Symbol != ""},
			ProjectID: projectId,
		}

		_, err := app.GetDB().InsertLabel(ctx, params)
		if err != nil {
			return fmt.Errorf("label insertion error: %w", err)
		}
	}

	return nil
}
