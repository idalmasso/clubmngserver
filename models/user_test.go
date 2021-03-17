package models

import (
	"context"
	"errors"
	"testing"

	"github.com/idalmasso/clubmngserver/common"
)

func TestGetAllUsers(t *testing.T){
	for _,db:=range(databasesToTest){
		
		InitDB(&db)
		users, err:=GetUsersList(context.Background())
		if err!=nil{
			t.Errorf("Cannot get all users after initialization")
		} else if len(users)!=len(usersToBeAdded){
			t.Errorf("Len users =%d, expected %d users after initialization : %v", len(users), len(usersToBeAdded),users)
		} else{
			for _,user:=range(users){
				_, ok:=usersToBeAdded[user.Username]
				if !ok{
					t.Fatalf("Find user %s in db users but not in initialized ones: %v", user.Username, usersToBeAdded)
				}
			}
			_,err=AddUser(context.Background(), common.UserData{Username: "user01"}, "Abcd")
			t.Cleanup(func(){
				db.RemoveUser(context.Background(),common.UserData{Username: "user01"})
			})
			if err!=nil{
				t.Errorf("Cannot add one user after initialization")
			}
			users, err=GetUsersList(context.Background())
			if err!=nil{
				t.Errorf("Cannot get all users after initialization and added 1")
			} else if len(users)!=len(usersToBeAdded)+1{
				t.Errorf("Len users =%d, expected %d users after initialization and added 1", len(users), len(usersToBeAdded)+1)
			} else{
				for _,user:=range(users){
					_, ok:=usersToBeAdded[user.Username]
					if !ok && user.Username!="user01"{
						t.Fatalf("Find user %s in db users but not in initialized ones", user.Username)
					}
				}
			}
		}
	}
}
func TestFindSingleUser(t *testing.T){
	for _,db:=range(databasesToTest){
		
		InitDB(&db)
		t.Run("Get all users initialization one by one", func(t *testing.T){
			for userName, passRole:=range(usersToBeAdded){
				user, err:=FindUser(context.Background(), userName)
				if err!=nil{
					t.Fatalf("Cannot get user %s", user)
				}
				if user==nil{
					t.Fatalf("User %s nil", userName)
				}
				if user.Role!=passRole.role{
					t.Fatalf("User %s should not have role %s but %s", userName, user.Role, passRole.role)
				}
				
				if err=checkUserPassword(context.Background(), user.Username, passRole.password); err!=nil{
					t.Fatalf("User cannot check old password or user %s", user.Username)
				}
			}
		})
		t.Run("Get one user not existent", func(t *testing.T){
			user, err:=FindUser(context.Background(), "NotExistentUser")
			var notFoundError common.NotFoundError
			if err==nil{
				t.Errorf("Error get User nil")
			} else if !errors.As(err, &notFoundError){
				t.Errorf("Error get User is not of correct type")
			} else if user != nil{
				t.Errorf("User is not nil")
			}
		})
		t.Run("Get one user after having added it", func(t *testing.T){
			_,err:=AddUser(context.Background(), common.UserData{Username: "user01"}, "Abcd" )
	
			t.Cleanup(func(){
				db.RemoveUser(context.Background(), common.UserData{Username: "user01"})
			})
			if err!=nil{
				t.Fatal("Cannot run test, cannot add user user01")
			}
			user, err:=FindUser(context.Background(), "user01")
			if err!=nil{
				t.Fatalf("Error in FindUser: %v", err)
			}
			if user.Username!="user01"{
				t.Fatalf("Role name is '%s', not 'user01'", user.Username)
			}
			if err=checkUserPassword(context.Background(), user.Username, "Abcd"); err!=nil{
				t.Fatalf("User cannot check old password or user %s", user.Username)
			}
		})
	}
}

type addUserTest struct{
	testName string
	userToAdd string
	resultError bool
	errorType error
}
var addUserTestSet= []addUserTest{
	{testName:"Test add admin user again", userToAdd: "admin" , resultError: true, errorType: common.AlreadyExistsError{} },
	{testName:"Test normal user add", userToAdd: "user01",  resultError: false, errorType: nil },
}


func TestAddUser(t *testing.T){
	for _,db:=range(databasesToTest){
		InitDB(&db)
		for _, test:=range(addUserTestSet){
			t.Run(test.testName, func (t *testing.T){
				user, err:=AddUser(context.Background(), common.UserData{Username:test.userToAdd}, "PasswordTest")
				if  test.resultError{
					if err==nil {
						t.Errorf("Expected test to give error but not")
						if user!=nil{
							t.Cleanup(func(){
								db.RemoveUser(context.Background(), *user)
							})
						}
					} else if !errors.As(err, &test.errorType){
						t.Fatalf("Expected test to give %t type of error, got %t",  test.errorType, err)
					}
				}else{
					t.Cleanup(func(){
						db.RemoveUser(context.Background(), *user)
					})
					if err!=nil{
						t.Fatalf("Eexpected test to not give error but it has given")
					} else if user.Username!=test.userToAdd {
						t.Fatalf("Wrong inserted username, %s instead of %s", user.Username, test.userToAdd)
					} 
					if err=checkUserPassword(context.Background(), user.Username, "PasswordTest");err!=nil{
						t.Fatalf("Cannot check password after inserted")
					}

				}
			})
		}
	}
}

func TestTryAuthenticate(t *testing.T){
	for _,db:=range(databasesToTest){
		InitDB(&db)
		user,err:=AddUser(context.Background(), common.UserData{Username:"user01"}, "Abcd")
		if err!=nil{
			t.Fatalf("Cannot create test user")
		}
		t.Cleanup(func(){
			db.RemoveUser(context.Background(), *user)
		})
		t.Run("Try user authentication with correct password",func(t *testing.T){
			auth, autor, err:=TryAuthenticate(context.Background(), *user,"Abcd")
			if err!=nil{
				t.Fatalf("Error in try authentication,not expected")
			}
			if auth=="" || autor==""{
				t.Fatalf("Cannot get both tokens")
			}
		})
		t.Run("Try user authentication with incorrect password",func(t *testing.T){
			auth, autor, err:=TryAuthenticate(context.Background(), *user,"Popopo")
			if err==nil{
				t.Fatalf("No error in try authentication,expected one")
			}
			if auth!="" || autor!=""{
				t.Fatalf("Both token should be empty")
			}
		})
	}
}

func TestCheckTokenAndGetAuthenticationAuthorization(t *testing.T){
	for _,db:=range(databasesToTest){
		InitDB(&db)
		user,err:=AddUser(context.Background(), common.UserData{Username:"user01"}, "Abcd")
		if err!=nil{
			t.Fatalf("Cannot create test user")
		}
		t.Cleanup(func(){
			db.RemoveUser(context.Background(), *user)
		})
		auth, autor, err:=TryAuthenticate(context.Background(), *user,"Abcd")
		if err!=nil{
			t.Fatalf("Error in try authentication,not expected")
		}
		t.Run("Try check token authentication with correct authentication token",func(t *testing.T){
			_, ok:=CheckToken(auth, false)
			if !ok{
				t.Fatalf("Error in check token,not expected")
			}
		})
		t.Run("Try check token authentication with incorrect authentication token(use authorization token)",func(t *testing.T){
			_, ok:=CheckToken(autor, false)
			if ok{
				t.Fatalf("No error in check token, expected one")
			}
		})
		t.Run("Try check token authentication with correct authorization token",func(t *testing.T){
			_, ok:=CheckToken(autor, true)
			if !ok{
				t.Fatalf("Error in check token,not expected")
			}
		})
		t.Run("Try check token authentication with incorrect authorization token(use authentication token)",func(t *testing.T){
			_, ok:=CheckToken(auth, true)
			if ok{
				t.Fatalf("No error in check token, expected one")
			}
		})
	}
}
//TODO: Refactor this test, and also the function
func TestGetNewAuthorizationToken(t *testing.T){
	for _,db:=range(databasesToTest){
		InitDB(&db)
		user,err:=AddUser(context.Background(), common.UserData{Username:"user01"}, "Abcd")
		if err!=nil{
			t.Fatalf("Cannot create test user")
		}
		t.Cleanup(func(){
			db.RemoveUser(context.Background(), *user)
		})
		auth, autor, err:=TryAuthenticate(context.Background(), *user,"Abcd")
		if err!=nil{
			t.Fatalf("Error in try authentication,not expected")
		}
		t.Run("Try get token authorization with correct authentication token",func(t *testing.T){
			autor, err:=GetAuthorizationTokenForAuthenticationToken(context.Background(), auth, *user)
			if err!=nil{
				t.Fatalf("Error in GetAuthorizationTokenForAuthenticationToken,not expected")
			}
			_, ok:=CheckToken(autor, true)
			if !ok{
				t.Fatalf("Error in check token,not expected")
			}
		})
		t.Run("Try get token authorization with incorrect authentication token",func(t *testing.T){
			_, err:=GetAuthorizationTokenForAuthenticationToken(context.Background(), autor, *user)
			if err!=nil{
				t.Fatalf("No error in GetAuthorizationTokenForAuthenticationToken, expected")
			}
			
		})
		
	}
}
