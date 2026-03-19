package app

import (
	"context"
	"encoding/base64"

	"charm.land/log/v2"
	"charm.land/wish/v2"
	"github.com/charmbracelet/ssh"
	"github.com/jackc/pgx/v5"
)

func (a *App) AuthMiddleware(next ssh.Handler) ssh.Handler {
	// TODO: handle case where session has no key
	return func(s ssh.Session) {
		// sum := sha256.Sum256(s.PublicKey().Marshal())
		log.Info(s.PublicKey())
		sum := []byte{}
		fingerprint := base64.StdEncoding.EncodeToString(sum[:])

		ctx := context.Background()
		userId, err := a.DB.GetUserIdFromSSHFingerprint(ctx, fingerprint)

		switch err {
		case nil:
			// This SSH key is associated with an account, perfect! just save it to the map and go on.
			a.sessionIdToUserIds[s.Context().SessionID()] = userId
			next(s)
		case pgx.ErrNoRows:
			// This SSH key is not associated with any account!
			// Does the username exist then?
			isTaken, err := a.DB.IsUsernameTaken(ctx, s.User())
			log.Info("login attempt with no key found in db", "isTaken", isTaken, "err", err)
			// Username does not exist --> Create new account with that username and register this SSH key
			// Username exists --> For now, error. In the future, ask other device to register you.
		default:
			wish.Println(s, "There has been an error in authenticating you.")
		}

	}
}
