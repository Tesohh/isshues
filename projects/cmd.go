package projects

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	db "github.com/Tesohh/isshues/db/generated"
)

type UpdateProjectsMsg struct {
	Projects []db.Project
}

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
