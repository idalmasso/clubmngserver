package models

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

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


type UserDetails struct {
	Username string
	Email string
	PasswordHash []byte
	AuthenticationTokens map[string] string
	AuthorizationTokens map[string] struct{}
}
func (user *UserDetails) GetAuthenticatonAuthorizationToken() (string,string,  error){
	return getAuthenticationAuthorizationTokens(user.Username)
}

func (user *UserDetails) TryAuthenticate(context context.Context, password string) (string,string,  error){
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
	return authentication, authorization, nil
}

func (user *UserDetails) GetUserTokenForAuthenticationToken(authenticationToken string ) (string, error){
	authorization, err:= getAuthorizationToken(user.Username)
	if err!= nil{
		return "", err
	}
	var val struct{}
	user.AuthenticationTokens[authenticationToken] = authorization
	user.AuthorizationTokens[authorization] = val
	return authorization, nil
}

func (user *UserDetails) RemoveUserAuthentication(authenticationToken string) {
	if author, ok:=user.AuthenticationTokens[authenticationToken]; ok{
		delete(user.AuthorizationTokens, author)
	}
	delete(user.AuthenticationTokens, authenticationToken)
}

var users map[string]UserDetails=make(map[string]UserDetails, 0) 


func FindUser(ctx context.Context, username string) *UserDetails {
	u,ok:=users[username]
	if !ok{
		return nil
	}
	return &u
}

func AddUser(user *UserDetails, password string) (UserDetails,error){
	_, ok:=users[user.Username]
	if ok{
		return UserDetails{}, fmt.Errorf("Username already exists")
	}
	passwordHash, err:=bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err!=nil{
		return UserDetails{}, err
	}
	user.PasswordHash=passwordHash
	users[user.Username] = *user
	return users[user.Username], nil
}

func checkUserPassword(context context.Context, username string, password string) error{
	u:=FindUser(context,username)
	if u==nil{
		return fmt.Errorf("Not found")
	}
	return bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password))
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
