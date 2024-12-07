package interfaces

import (
	"github.com/OxytocinGroup/theca-backend/pkg"
)

type UserUseCase interface {
	Register(email, password, username string) pkg.Response
	VerifyEmail(email, code string) pkg.Response
}
