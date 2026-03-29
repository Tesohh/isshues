package config

type Priorities map[string]Priority

func (p Priorities) FindClosest(value int) Priority {
	var (
		best          Priority
		bestFound     bool
		fallback      Priority
		fallbackFound bool
	)

	for _, pr := range p {
		if !fallbackFound || pr.Value < fallback.Value {
			fallback = pr
			fallbackFound = true
		}

		if pr.Value <= value {
			if !bestFound || pr.Value > best.Value {
				best = pr
				bestFound = true
			}
		}
	}

	if bestFound {
		return best
	}
	return fallback
}

type Priority struct {
	Value    int    `mapstructure:"value"`
	ColorKey string `mapstructure:"color"`
}
