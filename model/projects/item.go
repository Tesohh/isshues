package projects

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	db "github.com/Tesohh/isshues/db/generated"
)

type item struct {
	title, desc string
}

func itemFromProject(p db.Project) item {
	return item{
		title: fmt.Sprintf("[#%s] %s", p.Prefix, p.Title),
		desc:  "TODO!",
	}
}
func itemsFromProjects(ps []db.Project) []list.Item {
	items := []list.Item{}
	for _, p := range ps {
		items = append(items, itemFromProject(p))
	}
	return items
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
