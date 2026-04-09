package config

type DefaultLabel struct {
	Name   string `mapstructure:"name"`
	Color  string `mapstructure:"color"`
	Symbol string `mapstructure:"symbol"`
}

// # Labels that are added by default on project creation. Recommended to keep some "well defined" labels for common categories.
// [[default_labels]]
// # anything that improves the project: new features, optimizations, refactorings...
// name = "feat"
// color = "blue"
// symbol = "󰇈"
//
// [[default_labels]]
// # unexpected bugs or behaviours that pop up and need to be fixed. goes well with !crit...
// name = "fix"
// color = "red"
// symbol = ""
//
// [[default_labels]]
// # anything that is required to do, but doesn't "touch" the code, eg. documentation, presentations, submissions, deadlines, advertisements...
// name = "chore"
// color = "green"
// symbol = "󰃢"
