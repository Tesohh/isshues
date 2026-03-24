package projects

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"github.com/Tesohh/isshues/action"
	db "github.com/Tesohh/isshues/db/generated"
)

// Ask the projects view to refetch projects
type RefreshProjectsMsg struct{}

// Projects have been fetched, here is what the query got.
type UpdateProjectsMsg struct {
	Projects []db.Project
}

type InitHasCreatePermissionMsg struct{}

func (m ProjectsView) FetchProjectsCmd() tea.Msg {
	log.Info("Fetching projects")

	ctx := context.Background()
	projects, err := m.app.GetDB().GetUserProjectMemberships(ctx, m.userId)
	if err != nil {
		log.Error("GetUserProjectMemberships error", "err", err, "userId", m.userId)

		// TODO: return the error as a message
		return nil
	}
	return UpdateProjectsMsg{Projects: projects}
}

func (m ProjectsView) MakeCreateProjectCmd(title, prefix string) func() tea.Msg {
	return func() tea.Msg {
		_ = action.CreateProject(m.app, m.userId, title, prefix)
		// TODO: handle error

		m.app.Broadcast(RefreshProjectsMsg{})
		return nil
	}
}

func (m ProjectsView) HasCreatePermissionCmd() tea.Msg {
	ctx := context.Background()
	hasPermission, _ := m.app.DB.UserHasGlobalPermission(ctx, db.UserHasGlobalPermissionParams{
		UserID:             m.userId,
		GlobalPermissionID: "create-projects",
	})

	if hasPermission {
		return InitHasCreatePermissionMsg{}
	} else {
		return nil
	}
}
