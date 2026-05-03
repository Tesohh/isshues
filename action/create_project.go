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

var ErrDuplicatePrefix = fmt.Errorf("prefix already exists")

// note: prefix must be 4 chars and uppercase already
func CreateProject(app *app.App, query *db.Queries, userId int64, title string, prefix string) error {
	ctx := context.Background()

	// create the project
	projectId, err := query.InsertProject(ctx, db.InsertProjectParams{
		Title:  title,
		Prefix: prefix,
	})
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) && pgerr.Code == "23505" {
		return ErrDuplicatePrefix
	} else if err != nil {
		return fmt.Errorf("project insertion error: %w", err)
	}

	// add default groups
	var defaultGroups []config.DefaultGroup
	err = app.Viper.UnmarshalKey("default_groups", &defaultGroups)
	if err != nil {
		return fmt.Errorf("default_groups config error: %w", err)
	}

	groupMap := make(map[string]int64, len(defaultGroups))
	for _, group := range defaultGroups {
		params := db.InsertGroupParams{
			Name:        pgtype.Text{String: group.Name, Valid: true},
			ColorKey:    pgtype.Text{String: group.Color, Valid: group.Color != ""},
			Mentionable: group.Mentionable,
			ProjectID:   projectId,
		}
		groupId, err := query.InsertGroup(ctx, params)
		if err != nil {
			return fmt.Errorf("group insertion error: %w", err)
		}
		groupMap[group.Name] = groupId

		for _, permission := range group.Permissions {
			err = query.GrantPermissionToGroup(ctx, db.GrantPermissionToGroupParams{
				GroupID:             groupId,
				ProjectPermissionID: permission,
			})
			if err != nil {
				return fmt.Errorf("group grant permission error: %w", err)
			}
		}

		// add creator to the admins group
		if group.AddCreator {
			err = query.AddUserToGroup(ctx, db.AddUserToGroupParams{
				GroupID: groupId,
				UserID:  userId,
			})
			if err != nil {
				return fmt.Errorf("cannot add creator to grouop: %w", err)
			}
		}
	}

	// add default labels
	var defaultLabels []config.DefaultLabel
	err = app.Viper.UnmarshalKey("default_labels", &defaultLabels)
	if err != nil {
		return fmt.Errorf("default_labels config error: %w", err)
	}

	labelsMap := make(map[string]int64, len(defaultLabels))
	for _, label := range defaultLabels {
		params := db.InsertLabelParams{
			Name:      label.Name,
			ColorKey:  pgtype.Text{String: label.Color, Valid: label.Color != ""},
			Symbol:    pgtype.Text{String: label.Symbol, Valid: label.Symbol != ""},
			ProjectID: projectId,
		}

		label, err := query.InsertLabel(ctx, params)
		if err != nil {
			return fmt.Errorf("label insertion error: %w", err)
		}

		labelsMap[label.Name] = label.ID
	}

	// add default views
	var defaultViews []config.DefaultView
	err = app.Viper.UnmarshalKey("default_views", &defaultViews)
	if err != nil {
		return fmt.Errorf("default_views config error: %w", err)
	}

	for _, view := range defaultViews {
		priority := 0
		if view.Priority != nil {
			priority = *view.Priority
		}

		priorityMode := db.ViewPriorityModeEq
		if view.PriorityMode != "" {
			priorityMode = view.PriorityMode
		}

		labelsMode := db.ViewManyModeIgnore
		if view.LabelsMode != "" {
			labelsMode = view.LabelsMode
		}

		assigneesMode := db.ViewManyModeIgnore
		if view.LabelsMode != "" {
			labelsMode = view.LabelsMode
		}

		assigneeGroupMode := db.ViewManyModeIgnore
		if view.AssigneeGroupsMode != "" {
			assigneeGroupMode = view.AssigneeGroupsMode
		}

		sortBy := db.ViewSortByCode
		if view.SortBy != "" {
			sortBy = view.SortBy
		}

		sortOrder := db.ViewSortOrderAscending
		if view.SortOrder != "" {
			sortOrder = view.SortOrder
		}

		style := db.ViewStylePanels
		if view.Style != "" {
			style = view.Style
		}

		// insert view itself
		viewId, err := query.InsertView(ctx, db.InsertViewParams{
			ProjectID:              projectId,
			Name:                   view.Name,
			Title:                  pgtype.Text{String: view.Title, Valid: view.Title != ""},
			Statuses:               view.Statuses,
			Priority:               pgtype.Int4{Int32: int32(priority), Valid: view.Priority != nil},
			PriorityMode:           priorityMode,
			LabelsMode:             labelsMode,
			AssigneesMode:          assigneesMode,
			AssigneesIncludeViewer: view.AssigneesIncludeViewer,
			AssigneeGroupsMode:     assigneeGroupMode,
			SortBy:                 sortBy,
			SortOrder:              sortOrder,
			Style:                  style,
		})
		if err != nil {
			return fmt.Errorf("cannot insert default view: %w", err)
		}

		// get users in bulk and insert them
		users, err := query.GetUsersByUsernameBulk(ctx, view.Assignees)
		if err != nil {
			return fmt.Errorf("cannot fetch users for view: %w", err)
		}

		userArgs := make([]db.BulkInsertViewAssigneesParams, 0, len(users))
		for _, user := range users {
			userArgs = append(userArgs, db.BulkInsertViewAssigneesParams{ViewID: viewId, UserID: user.ID})
		}

		_, err = query.BulkInsertViewAssignees(ctx, userArgs)
		if err != nil {
			return fmt.Errorf("cannot insert view assignees: %w", err)
		}

		// insert groups
		groupArgs := make([]db.BulkInsertViewGroupAssigneesParams, 0)
		for _, name := range view.AssigneeGroups {
			if groupId, ok := groupMap[name]; ok {
				groupArgs = append(groupArgs, db.BulkInsertViewGroupAssigneesParams{ViewID: viewId, GroupID: groupId})
			}
		}

		_, err = query.BulkInsertViewGroupAssignees(ctx, groupArgs)
		if err != nil {
			return fmt.Errorf("cannot insert view group assignees: %w", err)
		}

		// insert labels
		labelArgs := make([]db.BulkInsertViewLabelsParams, 0)
		for _, name := range view.Labels {
			if labelId, ok := labelsMap[name]; ok {
				labelArgs = append(labelArgs, db.BulkInsertViewLabelsParams{ViewID: viewId, LabelID: labelId})
			}
		}

		_, err = query.BulkInsertViewLabels(ctx, labelArgs)
		if err != nil {
			return fmt.Errorf("cannot insert view labels: %w", err)
		}
	}

	return nil
}
