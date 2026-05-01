package db_complex

import (
	"fmt"
	"strings"

	db "github.com/Tesohh/isshues/db/generated"
)

// [x] ProjectID              int64
// [x] Title                  pgtype.Text
// [x] Status                 NullStatus
// [x] Priority               pgtype.Int4
// [x] PriorityMode           ViewPriorityMode
// [x] LabelsMode             ViewManyMode
// [x] AssigneesMode          ViewManyMode
// [x] AssigneesIncludeViewer bool
// [ ] AssigneeGroupsMode     ViewManyMode TODO
// [x] SortBy                 ViewSortBy
// [x] SortOrder              ViewSortOrder

type ViewQueryParams struct {
	ViewerUserID int64
}

// Generate a VERY large query for getting all issues related to a view.
func GenerateViewQuery(view db.View, params ViewQueryParams) (string, []any) {
	var b strings.Builder
	binds := []any{}

	// TODO: allow user to see only issues they have permission to view

	binds = append(binds, view.ProjectID)
	fmt.Fprintf(&b, "SELECT * FROM issues WHERE (issues.project_id = $%d)", len(binds))

	if view.Title.Valid {
		sanitized := strings.ReplaceAll(view.Title.String, `%`, `\%`)
		sanitized = strings.ReplaceAll(sanitized, `_`, `\_`)
		binds = append(binds, "%"+sanitized+"%")
		fmt.Fprintf(&b, " AND (issues.title ILIKE $%d ESCAPE '\\')", len(binds))
	}

	if len(view.Statuses) > 0 {
		binds = append(binds, view.Statuses)
		fmt.Fprintf(&b, " AND (issues.status = any($%d::status[]))", len(binds))
	}

	if view.Priority.Valid {
		symbol := map[db.ViewPriorityMode]string{
			db.ViewPriorityModeLt: "<",
			db.ViewPriorityModeLe: "<=",
			db.ViewPriorityModeEq: "=",
			db.ViewPriorityModeGe: ">=",
			db.ViewPriorityModeGt: ">",
		}[view.PriorityMode]

		binds = append(binds, view.Priority.Int32)
		fmt.Fprintf(&b, " AND (issues.priority %s $%d)", symbol, len(binds))
	}

	appendManyFilter(&b, &binds, view.ID, manyFilter{
		mode:       view.LabelsMode,
		viewTable:  "view_labels",
		viewIDCol:  "label_id",
		issueTable: "issue_labels",
		issueIDCol: "label_id",
	})

	var additionalAssigneeIDs []int64
	if view.AssigneesIncludeViewer {
		additionalAssigneeIDs = []int64{params.ViewerUserID}
	}

	appendManyFilter(&b, &binds, view.ID, manyFilter{
		mode:          view.AssigneesMode,
		viewTable:     "view_assignees",
		viewIDCol:     "assignee_id",
		issueTable:    "issue_assignees",
		issueIDCol:    "assignee_id",
		additionalIDs: additionalAssigneeIDs,
	})

	sortBy := map[db.ViewSortBy]string{
		db.ViewSortByCode:     "issues.code",
		db.ViewSortByEditDate: "issues.updated_at",
		db.ViewSortByPriority: "issues.priority",
	}[view.SortBy]

	sortOrder := map[db.ViewSortOrder]string{
		db.ViewSortOrderAscending:  "ASC",
		db.ViewSortOrderDescending: "DESC",
	}[view.SortOrder]
	fmt.Fprintf(&b, " ORDER BY %s %s", sortBy, sortOrder)

	return b.String(), binds
}
