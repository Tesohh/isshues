package shorthand

import (
	"strings"

	"github.com/Tesohh/isshues/app"
)

type ShorthandResults struct {
	Text            string
	UserMentionIDs  []int64
	GroupMentionIDs []int64
	DependencyIDs   []int64
	LabelIds        []int64
	Priority        int

	// anything that is problematic but doesn't break the issue creation
	Warnings []error
}

// Processes raw results from shorthand parser, giving out all info required to create a new issue
func Process(captures parserCaptures, app *app.App, projectId int64, userId int64) (ShorthandResults, error) {
	result := ShorthandResults{}
	// merge raws into text
	result.Text = strings.Join(captures.Raws, " ")

	// figure out which mentions are a. Users b. Groups c. Users and groups that don't exist and thus must be discarded
	// TODO: add GetUserFromUsername query
	// TODO: add GetGroupFromName query

	// fetch the dependencies and warn if they don't exist
	// TODO: add GetIssueFromCode query

	// fetch labels and create new labels if they don't exist and user has the "create-label" permission, otherwise warn
	// TODO: add GetLabelFromName query

	// check viper for Priority. if there is any problem, set the priority to 1 and warn.
	return result, nil
}
