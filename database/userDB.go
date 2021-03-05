package database

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)


type UserData struct {
	Username string
	Email string
	PasswordHash []byte
	AuthenticationTokens map[string] string
	AuthorizationTokens map[string] struct{}
	Role string
	
}


func (user *UserData) SetPassword(ctx context.Context, password string, db ClubDb) error {
	passwordHash, err:=bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err!=nil{
		return  err
	}
	user.PasswordHash=passwordHash
	*user, err=db.UpdateUser(ctx, *user)
	if err!=nil{
		return  err
	} 
	return nil
}
