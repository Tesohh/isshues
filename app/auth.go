package app

import (
	"context"
	"crypto/sha256"
	"encoding/base64"

	"charm.land/log/v2"
	"charm.land/wish/v2"
	db "github.com/Tesohh/isshues/db/generated"
	"github.com/charmbracelet/ssh"
	"github.com/jackc/pgx/v5"
)

const DefaultTheme = "git_hub_dark"

func (a *App) AuthMiddleware(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		sum := sha256.Sum256(s.PublicKey().Marshal())
		fingerprint := base64.StdEncoding.EncodeToString(sum[:])

		ctx := context.Background()
		userId, err := a.DB.GetUserIdFromAuth(ctx, db.GetUserIdFromAuthParams{
			Fingerprint: fingerprint,
			Username:    s.User(),
		})

		switch err {
		case nil:
			// This SSH key is associated with an account, perfect! just save it to the map and go on.
			a.SessionIdToUserIds[s.Context().SessionID()] = userId
			next(s)
		case pgx.ErrNoRows:
			// This SSH key is not associated with any account!
			// Does the username exist then?
			isTaken, err := a.DB.IsUsernameTaken(ctx, s.User())
			if err != nil && err != pgx.ErrNoRows {
				log.Error("username check error", "username", s.User(), "err", err)
				return
			}

			if isTaken {
				// Username exists --> For now, error. In the future, ask other device to register you.
				log.Warn("failed login attempt", "user", s.User(), "isTaken", isTaken, "err", err)
				wish.Println(s, "This username is already registered. For now, you can't register any other public key.")
			}

			// Username does not exist --> Create new account with that username and register this SSH key
			userId, err = a.DB.InsertUser(ctx, s.User())
			if err != nil {
				log.Error("user creation error", "username", s.User(), "err", err)
				wish.Println(s, "There has been an error in creating your user.")
				return
			}

			_, err = a.DB.InsertUserSettings(ctx, db.InsertUserSettingsParams{
				UserID: userId,
				Theme:  DefaultTheme,
			})
			if err != nil {
				log.Error("user settings creation err", "username", s.User(), "err", err)
				wish.Println(s, "There has been an error in creating your user.")
				return
			}

			log.Info("new user registered", "id", userId, "username", s.User())

			err = a.DB.RegisterSSHFingerprintToUser(ctx, db.RegisterSSHFingerprintToUserParams{
				UserID:      userId,
				Fingerprint: fingerprint,
			})
			if err != nil {
				log.Error("ssh key registering error", "username", s.User(), "err", err)
				wish.Println(s, "There has been an error in registering your ssh key. This username is now locked. Please contact your admin.")
				return
			}

			log.Info("added ssh key to user", "id", userId, "username", s.User())

			a.SessionIdToUserIds[s.Context().SessionID()] = userId
			next(s)
		default:
			wish.Println(s, "There has been an error in authenticating you.")
		}

	}
}
