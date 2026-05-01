package db_complex

import (
	"fmt"
	"strings"

	db "github.com/Tesohh/isshues/db/generated"
)

type manyFilter struct {
	mode db.ViewManyMode

	viewTable string // the view 2 something join table e.g. "view_labels"
	viewIDCol string // the id of the other something in the join table e.g. "label_id"

	issueTable string // the issue 2 something join table e.g. "issue_labels"
	issueIDCol string // the id of the other something in the join table e.g. "label_id"

	additionalIDs []int64 // additional ids to check in `subquery`, eg. the viewer user id
}

// auxiliary function for GenerateViewQuery
func appendManyFilter(b *strings.Builder, binds *[]any, viewID int64, f manyFilter) {
	if f.mode == db.ViewManyModeIgnore {
		return
	}

	*binds = append(*binds, viewID)
	view_bind_idx := len(*binds)

	// selects the `<thing>_id` from the specified view_<thing> table
	subquery := fmt.Sprintf("SELECT %s FROM %s WHERE view_id = $%d", f.viewIDCol, f.viewTable, view_bind_idx)

	if len(f.additionalIDs) > 0 {
		*binds = append(*binds, f.additionalIDs)
		additional_ids_bind_idx := len(*binds)
		subquery += fmt.Sprintf(" UNION SELECT unnest($%d::bigint[])", additional_ids_bind_idx)
	}

	switch f.mode {
	case db.ViewManyModeAny:
		// checks if any one of the <thing> related to the view are also in relation to the issue
		fmt.Fprintf(b, `
			AND EXISTS (
				SELECT 1 FROM %s
				WHERE issue_id = issues.id
				AND %s IN (%s)
			)`, f.issueTable, f.issueIDCol, subquery)
	case db.ViewManyModeAll:
		// checks if there is no <thing> related to the view that is not also associated with the issue
		fmt.Fprintf(b, `
			AND NOT EXISTS (
				SELECT 1 FROM (%s) AS wanted(id)
				WHERE NOT EXISTS (
					SELECT 1 FROM %s
					WHERE issue_id = issues.id AND %s = wanted.id
				)
			)`, subquery, f.issueTable, f.issueIDCol)
	case db.ViewManyModeExact:
		// if there are 0 of <thing> related to view: check if issue also has 0 of that <thing>.
		// otherwise: check if all <thing> are satistfied (similar to all mode)
		// 			  AND also check if there is any <thing> that is related to the isseu but not the view
		fmt.Fprintf(b, `
			AND CASE
				WHEN (SELECT COUNT(*) FROM %s WHERE view_id = $%d) = 0
				THEN NOT EXISTS (
					SELECT 1 FROM %s WHERE issue_id = issues.id
				)
				ELSE (
					NOT EXISTS (
						SELECT 1 FROM (%s) AS wanted(id)
						WHERE NOT EXISTS (
							SELECT 1 FROM %s
							WHERE issue_id = issues.id AND %s = wanted.id
						)
					)
					AND NOT EXISTS (
						SELECT 1 FROM %s
						WHERE issue_id = issues.id
						AND %s NOT IN (%s)
					)
				)
			END`, f.viewTable, view_bind_idx, f.issueTable, subquery, f.issueTable, f.issueIDCol, f.issueTable, f.issueIDCol, subquery)
	}
}
