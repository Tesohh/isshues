package shorthand

import "regexp"

type parserCaptures struct {
	Raws         []string
	Mentions     []string
	Labels       []string
	Priorities   []string
	Dependencies []string
}

var parser = regexp.MustCompile(`/(?mi)@(?P<Mention>\w+)|\+(?P<Label>\w+)|>(?P<DependencyCode>\d+)|!(?P<Priority>\w+)/gm`)

// Parses a message by running the regex, with no additional processing
func Parse(raw string) parserCaptures {
	captures := parserCaptures{}
	ptrs := []*[]string{&captures.Raws, &captures.Mentions, &captures.Labels, &captures.Dependencies, &captures.Priorities}

	matches := parser.FindAllStringSubmatch(raw, -1)
	for _, match := range matches {
		for i, subexp := range match {
			if subexp != "" {
				*ptrs[i] = append(*ptrs[i], subexp)
			}
		}
	}

	return captures
}
