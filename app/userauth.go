package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/idalmasso/clubmngserver/common"
)

//GetUsersList return a list of all users
func (app *App)GetUsersList(ctx context.Context) ([]common.UserData, error){
	users, err:=app.db.GetAllUsers(ctx)
	return users,err
}

//FindUser finds an actual user
func (app *App)FindUser(ctx context.Context, username string) (*common.UserData,error) {
	u, err:=app.db.FindUser(ctx, username)
	if err!=nil{
		return nil, err
	}
	if u==nil{
		return nil, common.NotFoundError{ID: username}
	}
	return u, nil
}
//AddUser add a user to the database and set its password. Returns the user created
func (app *App)AddUser(ctx context.Context,user common.UserData, password string) (*common.UserData,error){
	if userOld,err:=app.db.FindUser(ctx, user.Username);	 err==nil && userOld!=nil{
		return nil,  common.AlreadyExistsError{ID:user.Username}
	}
	addedUser, err:=app.db.AddUser(ctx, user); 
	if err!=nil{
		return nil, err
	} 
	if err:=addedUser.SetPassword(ctx, password); err!=nil{
		return nil, err
	}
	userUpdated, err:=app.db.UpdateUser(ctx, *addedUser)
	if err!=nil{
		return nil, err
	}
	return userUpdated, nil
}
//ChangePassword change the password to a user. Returns the user updated
func (app *App)ChangePassword(ctx context.Context, username string, newPassword string)(*common.UserData, error){
user,err:=app.FindUser(ctx,username)
if err!=nil{
	return nil, common.NotFoundError{ID: username}
}	
if err:=user.SetPassword(ctx, newPassword); err!=nil{
		return nil, err
	}
	return app.db.UpdateUser(ctx, *user)
}

func (app *App)checkUserPassword(context context.Context, username string, password string) error{
	u,err:=app.FindUser(context,username)
	if err!=nil{
		return common.NotFoundError{ID: username}
	}
	return u.CheckPassword(password)
}
//CheckToken check if a token string is ok for the authentication token or, if authorizationToken is true, for the authorization token
//
//only check if the token is valid, not if it is ok in the list of the user (the logoff removes it)
func (app *App)CheckToken (tokenString string, authorizationToken bool) (*jwt.Token, bool) {
	  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if authorizationToken{
		return []byte(app.getAuthorizationSecret()), nil
		}
		return []byte(app.getAuthenticationSecret()), nil
		 
	})
	if err!=nil{
		return nil, false
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
     return nil, false
  }
	return token, true
}

func (app *App)createToken(username string, secret string, duration time.Duration) (string, error) {
  var err error
  //Creating Access Token
  
  atClaims := jwt.MapClaims{}
  atClaims["authorized"] = true
  atClaims["username"] = username
  atClaims["exp"] = time.Now().Add(duration).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	
  token, err := at.SignedString([]byte(secret))
  if err != nil {
     return "", err
  }
  return token, nil
}

func (app *App)getAuthenticationAuthorizationTokens(username string) (string, string, error){
	authenticationToken, err:=app.createToken(username, app.getAuthenticationSecret(), time.Hour*24*15)
		if err!=nil{
		return "","", err
	}
	authorizationToken, err:=app.getAuthorizationToken(username)
	if err!=nil{
		return "","", err
	}
	return authenticationToken, authorizationToken, nil
}

func (app *App)getAuthorizationToken(username string) (string, error){
	authorizationToken, err:=app.createToken(username, app.getAuthorizationSecret(), time.Minute*5)
	if err!=nil{
		return "", err
	}
	return authorizationToken, nil
}

func (app *App)getAuthorizationSecret() string{
	secret:=os.Getenv("ACCESS_SECRET_AUTHORIZATION")
	if secret==""{
		//That's surely a big secret this way...
		secret="sdmalncnjsdsmfAuthorization"
	}
	return secret
}

func (app *App)getAuthenticationSecret() string{
	secret:=os.Getenv("ACCESS_SECRET_AUTHENTICATION")
	if secret==""{
		//That's surely a big secret this way...
		secret="sdmalncnjsdsmf"
	}
	return secret
}

//TryAuthenticate tries to authenticate the user with a password passed. 
//
//If ok, it sets and returns also the two tokens; returns authentication/authorization/error
//
//Else returns empty strings and error
func (app *App) TryAuthenticate(context context.Context, user common.UserData,  password string) (string,string,  error){
	if err:=app.checkUserPassword(context,user.Username, password); err!=nil{
		return "","", err
	}
	authentication, authorization, err:= app.getAuthenticationAuthorizationTokens(user.Username)
	if err!=nil{
		return "","", err
	}
	var val struct {}
	user.AuthenticationTokens[authentication]=authorization
	user.AuthorizationTokens[authorization] = val
	app.db.UpdateUser(context, user)
	return authentication, authorization, nil
}
//GetAuthorizationTokenForAuthenticationToken get a new Authorizationtoken and set it in the user for an existing authorization token
func (app *App)GetAuthorizationTokenForAuthenticationToken(context context.Context, authenticationToken string, user common.UserData ) (string, error){
	authorization, err:= app.getAuthorizationToken(user.Username)
	if err!= nil{
		return "", err
	}
	var val struct{}
	user.AuthenticationTokens[authenticationToken] = authorization
	user.AuthorizationTokens[authorization] = val
	app.db.UpdateUser(context, user)
	return authorization, nil
}
//RemoveUserAuthentication removes a token from a user/invalidate an authentication user token so it cannot be used anymore
func (app *App) RemoveUserAuthentication(ctx context.Context,username string,authenticationToken string) error {
	user,err:=app.FindUser(ctx, username)
	if err!=nil{
		return err
	}
	if user==nil{
		return common.NotFoundError{ID:username}
	}
	if author, ok:=user.AuthenticationTokens[authenticationToken]; ok{
		delete(user.AuthorizationTokens, author)
	}
	delete(user.AuthenticationTokens, authenticationToken)
	_, err= app.db.UpdateUser(ctx, *user)
	return err
}
//IsValidTokenForUser check if a token is valid and if it is ok for a user, so if it has been created and not revoked
func (app *App)IsValidTokenForUser(ctx context.Context, username string, token string, isAuthorization bool) (bool,error){
	user, err := app.db.FindUser(ctx, username)
	if err!=nil{
		return false, err
	}
	if user==nil{
		return false, common.NotFoundError{ID:username}
	}
	_, ok:=app.CheckToken(token, isAuthorization)
	if !ok{
		return false, nil
	}
	if isAuthorization{
		_, ok=user.AuthorizationTokens[token]
	}else{
		_,ok=user.AuthenticationTokens[token]
	}
	return ok, nil
}
