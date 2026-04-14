package shorthand

import (
	"regexp"
	"slices"
	"strings"
)

type parserCaptures struct {
	Text         string
	Mentions     []string
	Labels       []string
	Priorities   []string
	Dependencies []string
	Descriptions []string
}

var parser = regexp.MustCompile(`(?mi)@(?P<Mention>[a-z_][a-z0-9_-]*)|\+(?P<Label>\w*)|>(?P<DependencyCode>\d*)|!(?P<Priority>\w*)|desc\:(?P<Description>.+):`)

// Parses a message by running the regex, with no additional processing (except for content)
func Parse(msg string) parserCaptures {
	captures := parserCaptures{}
	rawCaptures := []string{}
	ptrs := []*[]string{&rawCaptures, &captures.Mentions, &captures.Labels, &captures.Dependencies, &captures.Priorities, &captures.Descriptions}

	matches := parser.FindAllStringSubmatch(msg, -1)
	for _, match := range matches {
		for i, subexp := range match {
			if subexp != "" {
				*ptrs[i] = append(*ptrs[i], subexp)
			}
		}
	}

	for _, ptr := range ptrs {
		slices.Sort(*ptr)
		*ptr = slices.Compact(*ptr)
	}

	for _, raw := range rawCaptures {
		msg = strings.ReplaceAll(msg, raw, "")
	}
	captures.Text = strings.Join(strings.Fields(msg), " ")

	return captures
}
