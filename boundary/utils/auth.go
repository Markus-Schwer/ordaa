package utils

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/boundary/auth"
	"gitlab.com/sfz.aalen/hackwerk/dotinder/entity"
	"gorm.io/gorm"
)

func CurrentUser(c echo.Context, repo entity.Repository, tx *gorm.DB) (*entity.User, error) {
	jwt := c.Get("user").(*jwt.Token)
	claims := jwt.Claims.(*auth.Claims)
	userUuid, err := uuid.FromString(claims.UserUuid)
	if err != nil {
		return nil, fmt.Errorf("could not parse user uuid %s as uuid: %w", userUuid, err)
	}
	return repo.GetUser(tx, &userUuid)
}

func BuildJwt(user *entity.User) *jwt.Token {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &auth.Claims{
		Username: user.Name,
		UserUuid: user.Uuid.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}
