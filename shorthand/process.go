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
	ErrWarningMentionFailed                = errors.New("no user (lenient) or group (exact) found")
	ErrWarningIssueNotFound                = errors.New("no issue found")
	ErrWarningLabelNotFoundAndNoPermission = errors.New("no label found, and you lack the `create-labels` permission to create a new one")
	ErrWarningInternalError                = errors.New("internal error. tell your admin")
	ErrWarningInternalErrorDefaulting      = errors.New("internal error. tell your admin! defaulting to priority 1")
	ErrWarningInvalidPriority              = errors.New("invalid priority")
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
			// add the recruiter as an assignee in case noone was mentioned
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

				// then in case the user doesn't exist, try looking for a group...
				group, err := app.DB.GetGroupByName(ctx, db.GetGroupByNameParams{
					Name:      pgtype.Text{String: mention, Valid: true},
					ProjectID: projectId,
				})
				switch err {
				case pgx.ErrNoRows:
					// no user or group was found, warn
					result.Warnings = append(result.Warnings, fmt.Errorf("%w with name: %s", ErrWarningMentionFailed, mention))
				case nil:
					result.GroupMentions = append(result.GroupMentions, group)
				default:
					return result, err
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

		switch err {
		case pgx.ErrNoRows:
			result.Warnings = append(result.Warnings, fmt.Errorf("%w with code: %d", ErrWarningIssueNotFound, code))
		case nil:
			result.Dependencies = append(result.Dependencies, issue)
		default:
			return result, err
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

		switch err {
		case pgx.ErrNoRows:
			// if user has `create-label` permission, then create the label
			result.PendingLabels = append(result.PendingLabels, strings.ToLower(labelName))
		case nil:
			result.Labels = append(result.Labels, label)
		default:
			return result, err
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
			return 1, ErrWarningInternalErrorDefaulting
		}

		return defaultPriority.Value, nil
	}

	var priorities config.Priorities
	err := viper.UnmarshalKey("priorities", &priorities)
	if err != nil {
		log.Warn("config error, in 'priorities', with captures path. defaulting to 1", "err", err)
		return 1, ErrWarningInternalErrorDefaulting
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

	return 1, fmt.Errorf("%w. must be an integer, or one of (%s)", ErrWarningInvalidPriority, options)
}
