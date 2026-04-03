package cli

// func new_issue_cmd(session ssh.Session, app *app.App, _ **tea.Program) *cobra.Command {
// 	projectCmd := &cobra.Command{
// 		Use: "new ",
// 	}
//
// 	newCmd := &cobra.Command{
// 		Use:  "new [prefix] [title]",
// 		Args: cobra.MinimumNArgs(2),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			userId, ok := app.SessionIdToUserIds[session.Context().SessionID()]
// 			if !ok {
// 				return errors.New("your userid was not found in the session map. might be an auth issue.")
// 			}
// 			return nil
// 		},
// 	}
//
// 	return newCmd
// }
