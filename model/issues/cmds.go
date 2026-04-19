package issues

import (
	"context"
	"errors"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/Tesohh/isshues/model"
)

var (
	NotMemberErr = errors.New("cannot load project, as you are not a member")
)

type UpdateProjectMsg struct {
	Project db.Project
}

func (m Model) LoadProjectCmd() tea.Msg {
	ctx := context.Background()
	// start transaction
	tx, err := m.app.DBPool.Begin(ctx)
	if err != nil {
		log.Error("issues.Model.LoadProjectCmd: cannot start transaction", "err", err)
		return model.InternalErrMsg()
	}
	defer tx.Rollback(ctx)
	query := db.New(tx)

	// check permission
	hasPermission, err := query.UserIsMemberOfProject(ctx, db.UserIsMemberOfProjectParams{UserID: m.userId, ProjectID: m.projectId})
	if err != nil {
		log.Error("issues.Model.LoadProjectCmd: error when checking permission", "err", err)
		return model.InternalErrMsg()
	}

	if !hasPermission {
		return model.ErrMsg{Err: NotMemberErr}
	}

	// get project
	project, err := query.GetProjectById(ctx, m.projectId)
	if err != nil {
		log.Error("issues.Model.LoadProjectCmd: error when querying project", "err", err)
		return model.InternalErrMsg()
	}

	return UpdateProjectMsg{Project: project}
}
