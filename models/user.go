package models

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/idalmasso/clubmngserver/common"
)


func GetUsersList(ctx context.Context) ([]common.UserData, error){
	users, err:=db.GetAllUsers(ctx)
	return users,err
}
func  GetAuthenticatonAuthorizationToken(username string) (string,string,  error){
	return getAuthenticationAuthorizationTokens(username)
}

func  TryAuthenticate(context context.Context, user common.UserData,  password string) (string,string,  error){
	if err:=checkUserPassword(context,user.Username, password); err!=nil{
		return "","", err
	}
	authentication, authorization, err:= getAuthenticationAuthorizationTokens(user.Username)
	if err!=nil{
		return "","", err
	}
	var val struct {}
	user.AuthenticationTokens[authentication]=authorization
	user.AuthorizationTokens[authorization] = val
	db.UpdateUser(context, user)
	return authentication, authorization, nil
}

func GetUserTokenForAuthenticationToken(context context.Context, authenticationToken string, user common.UserData ) (string, error){
	authorization, err:= getAuthorizationToken(user.Username)
	if err!= nil{
		return "", err
	}
	var val struct{}
	user.AuthenticationTokens[authenticationToken] = authorization
	user.AuthorizationTokens[authorization] = val
	db.UpdateUser(context, user)
	return authorization, nil
}

func  RemoveUserAuthentication(user common.UserData,authenticationToken string) {
	if author, ok:=user.AuthenticationTokens[authenticationToken]; ok{
		delete(user.AuthorizationTokens, author)
	}
	delete(user.AuthenticationTokens, authenticationToken)
}

func FindUser(ctx context.Context, username string) (common.UserData,error) {
	u,err:=db.FindUser(ctx, username)
	if u==nil{
		return common.UserData{}, err
	}
	return *u, nil
}

func AddUser(ctx context.Context,user common.UserData, password string) (common.UserData,error){
	if _,err:=db.FindUser(ctx, user.Username);	 err==nil{
		return common.UserData{}, fmt.Errorf("Username already exists: %v",err)
	}
	user, err:=db.AddUser(ctx, user); 
	if err!=nil{
		return common.UserData{}, err
	} 
	if err:=user.SetPassword(ctx, password); err!=nil{
		return common.UserData{}, err
	}
	user, err=db.UpdateUser(ctx, user)
	if err!=nil{
		return common.UserData{}, err
	}
	return user, nil
}

func ChangePassword(ctx context.Context, user common.UserData, newPassword string)(common.UserData, error){
if err:=user.SetPassword(ctx, newPassword); err!=nil{
		return common.UserData{}, err
	}
	user, err:=db.UpdateUser(ctx, user)
	if err!=nil{
		return common.UserData{}, err
	}
	return user, nil
}

func checkUserPassword(context context.Context, username string, password string) error{
	u,err:=FindUser(context,username)
	if err!=nil{
		return fmt.Errorf("Not found")
	}
	return u.CheckPassword(password)
}

func CheckToken (tokenString string, authorizationToken bool) (*jwt.Token, bool) {
	 token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
     //Make sure that the token method conform to "SigningMethodHMAC"
     if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		 }
		 if authorizationToken{
		 	return []byte(getAuthorizationSecret()), nil
		 }else{
			 return []byte(getAuthenticationSecret()), nil
		 }
	})
	if err!=nil{
		return nil, false
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
     return nil, false
  }
	return token, true
}

func createToken(username string, secret string, duration time.Duration) (string, error) {
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

func getAuthenticationAuthorizationTokens(username string) (string, string, error){
	authenticationToken, err:=createToken(username, getAuthenticationSecret(), time.Hour*24*15)
		if err!=nil{
		return "","", err
	}
	authorizationToken, err:=getAuthorizationToken(username)
	if err!=nil{
		return "","", err
	}
	return authenticationToken, authorizationToken, nil
}

func getAuthorizationToken(username string) (string, error){
	authorizationToken, err:=createToken(username, getAuthorizationSecret(), time.Minute*5)
	if err!=nil{
		return "", err
	}
	return authorizationToken, nil
}

func getAuthorizationSecret() string{
	secret:=os.Getenv("ACCESS_SECRET_AUTHORIZATION")
	if secret==""{
		//That's surely a big secret this way...
		secret="sdmalncnjsdsmfAuthorization"
	}
	return secret
}

func getAuthenticationSecret() string{
	secret:=os.Getenv("ACCESS_SECRET_AUTHENTICATION")
	if secret==""{
		//That's surely a big secret this way...
		secret="sdmalncnjsdsmf"
	}
	return secret
}

