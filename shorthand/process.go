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
	"github.com/Tesohh/isshues/config"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/spf13/viper"
)

var (
	WarningMentionFailed                = errors.New("no user (lenient) or group (exact) found")
	WarningIssueNotFound                = errors.New("no issue found")
	WarningLabelNotFoundAndNoPermission = errors.New("no label found, and you lack the `create-labels` permission to create a new one")
	WarningInternalError                = errors.New("internal error. tell your admin!")
	WarningInternalErrorDefaulting      = errors.New("internal error. tell your admin! defaulting to priority 1")
	WarningInvalidPriority              = errors.New("invalid priority")
)

type ShorthandResults struct {
	Title       string
	Description string

	UserMentions  []db.User
	GroupMentions []db.Group

	Dependencies []db.Issue

	Labels        []db.Label
	PendingLabels []string
	Priority      int

	// anything that is problematic but doesn't break the issue creation
	Warnings []error
}

// Processes raw results from shorthand parser, giving out all info required to create a new issue
func Process(captures parserCaptures, app *app.App, projectId int64, userId int64) (ShorthandResults, error) {
	result := ShorthandResults{}
	ctx := context.Background()

	result.Title = captures.Text

	result.Description = strings.Join(captures.Descriptions, "\n")

	// figure out which mentions are a. Users b. Groups c. Users and groups that don't exist and thus must be discarded
	if !slices.Contains(captures.Mentions, "NOBODY") {
		if len(captures.Mentions) == 0 {
			user, err := app.DB.GetUserByID(ctx, userId)
			if err != nil {
				return result, err
			}
			result.UserMentions = append(result.UserMentions, user)
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

	// fetch labels
	// hasCreateLabelPermission, err := app.DB.UserHasProjectPermission(ctx, db.UserHasProjectPermissionParams{
	// 	UserID:              userId,
	// 	ProjectPermissionID: "create-labels",
	// 	ProjectID:           projectId,
	// })
	// if err != nil {
	// 	log.Warn("cannot check if user has create-labels permission while processing shorthand", "userId", userId, "err", err)
	// }

	for _, labelName := range captures.Labels {
		label, err := app.DB.GetLabelFromName(ctx, db.GetLabelFromNameParams{
			Name:      strings.ToLower(labelName),
			ProjectID: projectId,
		})

		if err != nil && err != pgx.ErrNoRows {
			return result, err
		} else if err == pgx.ErrNoRows {
			// if user has `create-label` permission, then create the label
			result.PendingLabels = append(result.PendingLabels, strings.ToLower(labelName))

			// if hasCreateLabelPermission {
			// 	// create the label
			// 	label, err = app.DB.InsertLabelBasic(ctx, db.InsertLabelBasicParams{
			// 		Name:      labelName,
			// 		ProjectID: projectId,
			// 	})
			// 	if err != nil {
			// 		log.Error("error while creating new label when processing shorthand", "labelName", labelName, "projectId", projectId, "err", err)
			// 		result.Warnings = append(result.Warnings, WarningInternalError)
			// 	}
			// 	// TODO: broadcast RefreshLabels msg
			//
			// 	result.Labels = append(result.Labels, label)
			// } else {
			// 	result.Warnings = append(result.Warnings, fmt.Errorf("%w. label: %s", WarningLabelNotFoundAndNoPermission, labelName))
			// }
		} else if err == nil {
			result.Labels = append(result.Labels, label)
		}
	}

	// check viper for Priority. if there is any problem, set the priority to 1 and warn.
	var err error
	result.Priority, err = parsePriorityWithViper(captures.Priorities, app.Viper)
	if err != nil {
		result.Warnings = append(result.Warnings, err)
	}

	return result, nil
}

func parsePriorityWithViper(captures []string, viper *viper.Viper) (int, error) {
	if len(captures) == 0 {
		var defaultPriority config.Priority
		err := viper.UnmarshalKey("priorities.default", &defaultPriority)
		if err != nil {
			log.Warn("config error, in 'priorities', with default path. (probably priorities.default is not defined). defaulting to 1", "err", err)
			return 1, WarningInternalErrorDefaulting
		}

		return defaultPriority.Value, nil
	} else {
		var priorities config.Priorities
		err := viper.UnmarshalKey("priorities", &priorities)
		if err != nil {
			log.Warn("config error, in 'priorities', with captures path. defaulting to 1", "err", err)
			return 1, WarningInternalErrorDefaulting
		}

		// in case the actual string was requested
		target := captures[len(captures)-1]
		if priority, ok := priorities[target]; ok {
			return priority.Value, nil
		}

		// in case an integer was (hopefully) requested
		if value, err := strconv.Atoi(target); err == nil {
			return value, nil
		}

		// the priority was invalid, default to 1 and give a useful warning
		keys := make([]string, 0, len(priorities))
		for k := range priorities {
			keys = append(keys, k)
		}

		options := strings.Join(keys, ", ")

		return 1, fmt.Errorf("%w. must be an integer, or one of (%s)", WarningInvalidPriority, options)
	}
}
