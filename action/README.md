This package contains common DB actions.
Common themes are:
- actions must assume that all input is valid
- actions must assume that the users have the necessary permissions
- actions must be "ui-agnostic", ie. they are suitable for CLI and TUI usage
- actions must take a `*db.Queries` instead of using the global app db, so that they can be used in a transaction.
