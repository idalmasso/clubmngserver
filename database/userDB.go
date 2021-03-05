package database

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)


type UserData struct {
	Username string
	Email string
	passwordHash []byte
	AuthenticationTokens map[string] string
	AuthorizationTokens map[string] struct{}
	Role string
}

func (user *UserData) SetPassword(ctx context.Context, password string, db ClubDb) error {
	passwordHash, err:=bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err!=nil{
		return  err
	}
	user.passwordHash=passwordHash
	*user, err=db.UpdateUser(ctx, *user)
	if err!=nil{
		return  err
	} 
	return nil
}
func (user *UserData) CheckPassword(password string) error{
		return bcrypt.CompareHashAndPassword(user.passwordHash, []byte(password))
}

func (user *UserData) AddRole(ctx context.Context,roleName string,db ClubDb) error{
	_, err:=db.FindRole(ctx, roleName)
	if err!=nil{
		return err
	}
	user.Role=roleName
	if _, err=db.UpdateUser(ctx, *user); err!=nil{
		return err
	}
	return nil
}
