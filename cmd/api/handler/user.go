package handler

import (
	"github.com/G-Villarinho/social-network/domain"
	"github.com/G-Villarinho/social-network/pkg"
	"github.com/labstack/echo"
)

type userHandler struct {
	di          *pkg.Di
	userService domain.UserService
}

func NewUserHandler(di *pkg.Di) (domain.UserHandler, error) {
	userService, err := pkg.Invoke[domain.UserService](di)
	if err != nil {
		return nil, err
	}

	return &userHandler{
		di:          di,
		userService: userService,
	}, nil
}

func (u *userHandler) CreateUser(ctx echo.Context) error {
	panic("unimplemented")
}
