package config

type Priorities map[string]Priority

func (p Priorities) FindClosest(value int) (Priority, string) {
	var (
		best          Priority
		bestK         string
		bestFound     bool
		fallback      Priority
		fallbackFound bool
	)

	for k, pr := range p {
		if !fallbackFound || pr.Value < fallback.Value {
			fallback = pr
			fallbackFound = true
		}

		if pr.Value <= value {
			if !bestFound || pr.Value > best.Value {
				best = pr
				bestK = k
				bestFound = true
			}
		}
	}

	if bestFound {
		return best, bestK
	}
	return fallback, "default"
}

type Priority struct {
	Value    int    `mapstructure:"value"`
	ColorKey string `mapstructure:"color"`
}
