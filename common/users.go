package common

import (
	"context"

	"github.com/idalmasso/clubmngserver/common/userprops"
	"golang.org/x/crypto/bcrypt"
)

//UserData is the actual data of a user. TBD: create map[string]UserAttribute, create UserAttribute{type , value {}} for attributes and manage them (so, no email!)
type UserData struct {
	Username string
	passwordHash []byte
	AuthenticationTokens map[string] string
	AuthorizationTokens map[string] struct{}
	Role string
	properties map[string]userprops.UserPropertyValue
}
//SetPassword sets the password to the user 
func (user *UserData) SetPassword(ctx context.Context, password string) error {
	passwordHash, err:=bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err!=nil{
		return  err
	}
	user.passwordHash=passwordHash
	return nil
}
//CheckPassword check if the password passed in the method is ok for the user
func (user *UserData) CheckPassword(password string) error{
		return bcrypt.CompareHashAndPassword(user.passwordHash, []byte(password))
}
//AddRole adds a role to the User
func (user *UserData) AddRole(ctx context.Context,role SecurityRole) {
	user.Role=role.Name
}
