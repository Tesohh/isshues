package shorthand

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"charm.land/log/v2"
	"github.com/Tesohh/isshues/app"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	WarningMentionFailed                = errors.New("no user (lenient) or group (exact) found")
	WarningIssueNotFound                = errors.New("no issue found")
	WarningLabelNotFoundAndNoPermission = errors.New("no label found, and you lack the `create-labels` permission to create a new one")
	WarningInternalError                = errors.New("internal error")
)

type ShorthandResults struct {
	Text string

	UserMentions  []db.User
	GroupMentions []db.Group
	Nobody        bool

	Dependencies []db.Issue

	Labels   []db.Label
	Priority int

	// anything that is problematic but doesn't break the issue creation
	Warnings []error
}

// Processes raw results from shorthand parser, giving out all info required to create a new issue
func Process(captures parserCaptures, app *app.App, projectId int64, userId int64) (ShorthandResults, error) {
	result := ShorthandResults{}
	ctx := context.Background()

	// merge raws into text
	result.Text = strings.Join(captures.Raws, " ")

	// figure out which mentions are a. Users b. Groups c. Users and groups that don't exist and thus must be discarded
	if slices.Contains(captures.Mentions, "nobody") {
		result.Nobody = true
	} else {
		for _, mention := range captures.Mentions {
			// first fetch users leniently (more likely)
			// TODO: consider doing a single id in ("...") query for faster querying

			user, err := app.DB.GetUserByUsernameLenient(ctx, pgtype.Text{String: mention, Valid: true})
			if err != nil && err != pgx.ErrNoRows {
				return result, err
			} else if err == nil {
				result.UserMentions = append(result.UserMentions, user)
				continue
			}

			group, err := app.DB.GetGroupByName(ctx, db.GetGroupByNameParams{
				Name:      pgtype.Text{String: mention, Valid: true},
				ProjectID: projectId,
			})
			if err != nil && err != pgx.ErrNoRows {
				return result, err
			} else if err == pgx.ErrNoRows {
				// no user or group was found, warn
				result.Warnings = append(result.Warnings, fmt.Errorf("%w with name: %s", WarningMentionFailed, mention))
			} else if err == nil {
				result.GroupMentions = append(result.GroupMentions, group)
			}
		}
	}

	// fetch the dependencies and warn if they don't exist
	// TODO: if the user lacks the `view-unassigned` permission, do not add the dependency to avoid leaking details, and warn
	for _, code := range captures.Dependencies {
		codeInt, _ := strconv.Atoi(code) // ignore the error, it would not be regex parsed anyway
		code := int64(codeInt)

		issue, err := app.DB.GetIssueFromCode(ctx, db.GetIssueFromCodeParams{
			Code:      code,
			ProjectID: projectId,
		})

		if err != nil && err != pgx.ErrNoRows {
			return result, err
		} else if err == pgx.ErrNoRows {
			result.Warnings = append(result.Warnings, fmt.Errorf("%w with code: %d", WarningIssueNotFound, code))
		} else if err == nil {
			result.Dependencies = append(result.Dependencies, issue)
		}
	}

	// fetch labels and create new labels if they don't exist and user has the "create-label" permission, otherwise warn
	hasCreateLabelPermission, err := app.DB.UserHasProjectPermission(ctx, db.UserHasProjectPermissionParams{
		UserID:              userId,
		ProjectPermissionID: "create-labels",
		ProjectID:           projectId,
	})
	if err != nil {
		log.Warn("cannot check if user has create-labels permission while processing shorthand", "userId", userId, "err", err)
	}

	for _, labelName := range captures.Labels {
		label, err := app.DB.GetLabelFromName(ctx, db.GetLabelFromNameParams{
			Name:      labelName,
			ProjectID: projectId,
		})

		if err != nil && err != pgx.ErrNoRows {
			return result, err
		} else if err == pgx.ErrNoRows {
			// if user has `create-label` permission, then create the label
			if hasCreateLabelPermission {
				// create the label
				label, err = app.DB.InsertLabelBasic(ctx, db.InsertLabelBasicParams{
					Name:      labelName,
					ProjectID: projectId,
				})
				if err != nil {
					log.Error("error while creating new label when processing shorthand", "labelName", labelName, "projectId", projectId, "err", err)
					result.Warnings = append(result.Warnings, WarningInternalError)
				}

				result.Labels = append(result.Labels, label)
			} else {
				result.Warnings = append(result.Warnings, fmt.Errorf("%w. label: %s", WarningLabelNotFoundAndNoPermission, labelName))
			}
		} else if err == nil {
			result.Labels = append(result.Labels, label)
		}
	}

	// TODO: check viper for Priority. if there is any problem, set the priority to 1 and warn.
	return result, nil
}
