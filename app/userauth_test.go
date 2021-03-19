package app

import (
	"context"
	"errors"
	"testing"

	"github.com/idalmasso/clubmngserver/common"
)

func TestGetAllUsers(t *testing.T){
	var testApp App
		testApp.InitDB(&testDatabase)
		users, err:=testApp.GetUsersList(context.Background())
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
			_,err=testApp.AddUser(context.Background(), common.UserData{Username: "user01"}, "Abcd")
			t.Cleanup(func(){
				testApp.db.RemoveUser(context.Background(),common.UserData{Username: "user01"})
			})
			if err!=nil{
				t.Errorf("Cannot add one user after initialization")
			}
			users, err=testApp.GetUsersList(context.Background())
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
func TestFindSingleUser(t *testing.T){
	var testApp App
		testApp.InitDB(&testDatabase)
		t.Run("Get all users initialization one by one", func(t *testing.T){
			for userName, passRole:=range(usersToBeAdded){
				user, err:=testApp.FindUser(context.Background(), userName)
				if err!=nil{
					t.Fatalf("Cannot get user %s", user)
				}
				if user==nil{
					t.Fatalf("User %s nil", userName)
				}
				if user.Role!=passRole.role{
					t.Fatalf("User %s should not have role %s but %s", userName, user.Role, passRole.role)
				}
				
				if err=testApp.checkUserPassword(context.Background(), user.Username, passRole.password); err!=nil{
					t.Fatalf("User cannot check old password or user %s", user.Username)
				}
			}
		})
		t.Run("Get one user not existent", func(t *testing.T){
			user, err:=testApp.FindUser(context.Background(), "NotExistentUser")
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
			_,err:=testApp.AddUser(context.Background(), common.UserData{Username: "user01"}, "Abcd" )
	
			t.Cleanup(func(){
				testApp.db.RemoveUser(context.Background(), common.UserData{Username: "user01"})
			})
			if err!=nil{
				t.Fatal("Cannot run test, cannot add user user01")
			}
			user, err:=testApp.FindUser(context.Background(), "user01")
			if err!=nil{
				t.Fatalf("Error in FindUser: %v", err)
			}
			if user.Username!="user01"{
				t.Fatalf("Role name is '%s', not 'user01'", user.Username)
			}
			if err=testApp.checkUserPassword(context.Background(), user.Username, "Abcd"); err!=nil{
				t.Fatalf("User cannot check old password or user %s", user.Username)
			}
		})
	
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
	var testApp App
	testApp.InitDB(&testDatabase)
	for _, test:=range(addUserTestSet){
		t.Run(test.testName, func (t *testing.T){
			user, err:=testApp.AddUser(context.Background(), common.UserData{Username:test.userToAdd}, "PasswordTest")
			if  test.resultError{
				if err==nil {
					if user!=nil{
						t.Cleanup(func(){
							testApp.db.RemoveUser(context.Background(), *user)
						})
					}
					t.Fatal("Expected test to give error but not")
				} 
				if !errors.As(err, &test.errorType){
					t.Fatalf("Expected test to give %t type of error, got %t",  test.errorType, err)
				}
				return
			}
			t.Cleanup(func(){
				testApp.db.RemoveUser(context.Background(), *user)
			})
			if err!=nil{
				t.Fatalf("Eexpected test to not give error but it has given")
			} else if user.Username!=test.userToAdd {
				t.Fatalf("Wrong inserted username, %s instead of %s", user.Username, test.userToAdd)
			} 
			if err=testApp.checkUserPassword(context.Background(), user.Username, "PasswordTest");err!=nil{
				t.Fatalf("Cannot check password after inserted")
			}
		
		})
	}
}

func TestTryAuthenticate(t *testing.T){
	var testApp App
		testApp.InitDB(&testDatabase)
		user,err:=testApp.AddUser(context.Background(), common.UserData{Username:"user01"}, "Abcd")
		if err!=nil{
			t.Fatalf("Cannot create test user")
		}
		t.Cleanup(func(){
			testApp.db.RemoveUser(context.Background(), *user)
		})
		t.Run("Try user authentication with correct password",func(t *testing.T){
			auth, author, err:=testApp.TryAuthenticate(context.Background(), *user,"Abcd")
			if err!=nil{
				t.Fatalf("Error in try authentication,not expected")
			}
			if auth=="" || author==""{
				t.Fatalf("Cannot get both tokens")
			}
		})
		t.Run("Try user authentication with incorrect password",func(t *testing.T){
			auth, author, err:=testApp.TryAuthenticate(context.Background(), *user,"Popopo")
			if err==nil{
				t.Fatalf("No error in try authentication,expected one")
			}
			if auth!="" || author!=""{
				t.Fatalf("Both token should be empty")
			}
		})
	
}

func TestCheckTokenAndGetAuthenticationAuthorization(t *testing.T){
	var testApp App
	testApp.InitDB(&testDatabase)
	user,err:=testApp.AddUser(context.Background(), common.UserData{Username:"user01"}, "Abcd")
	if err!=nil{
		t.Fatalf("Cannot create test user")
	}
	t.Cleanup(func(){
		testApp.db.RemoveUser(context.Background(), *user)
	})
	auth, author, err:=testApp.TryAuthenticate(context.Background(), *user,"Abcd")
	if err!=nil{
		t.Fatalf("Error in try authentication,not expected")
	}
	t.Run("Try check token authentication with correct authentication token",func(t *testing.T){
		_, ok:=testApp.CheckToken(auth, false)
		if !ok{
			t.Fatalf("Error in check token,not expected")
		}
	})
	t.Run("Try check token authentication with incorrect authentication token(use authorization token)",func(t *testing.T){
		_, ok:=testApp.CheckToken(author, false)
		if ok{
			t.Fatalf("No error in check token, expected one")
		}
	})
	t.Run("Try check token authentication with correct authorization token",func(t *testing.T){
		_, ok:=testApp.CheckToken(author, true)
		if !ok{
			t.Fatalf("Error in check token,not expected")
		}
	})
	t.Run("Try check token authentication with incorrect authorization token(use authentication token)",func(t *testing.T){
		_, ok:=testApp.CheckToken(auth, true)
		if ok{
			t.Fatalf("No error in check token, expected one")
		}
	})

}
//TODO: Refactor this test, and also the function
func TestGetNewAuthorizationToken(t *testing.T){
	var testApp App
	testApp.InitDB(&testDatabase)
	user,err:=testApp.AddUser(context.Background(), common.UserData{Username:"user01"}, "Abcd")
	if err!=nil{
		t.Fatalf("Cannot create test user")
	}
	t.Cleanup(func(){
		testApp.db.RemoveUser(context.Background(), *user)
	})
	auth, author, err:=testApp.TryAuthenticate(context.Background(), *user,"Abcd")
	if err!=nil{
		t.Fatalf("Error in try authentication,not expected")
	}
	t.Run("Try get token authorization with correct authentication token",func(t *testing.T){
		author, err:=testApp.GetAuthorizationTokenForAuthenticationToken(context.Background(), auth, *user)
		if err!=nil{
			t.Fatalf("Error in GetAuthorizationTokenForAuthenticationToken,not expected")
		}
		_, ok:=testApp.CheckToken(author, true)
		if !ok{
			t.Fatalf("Error in check token,not expected")
		}
	})
	t.Run("Try get token authorization with incorrect authentication token",func(t *testing.T){
		_, err:=testApp.GetAuthorizationTokenForAuthenticationToken(context.Background(), author, *user)
		if err!=nil{
			t.Fatalf("No error in GetAuthorizationTokenForAuthenticationToken, expected")
		}
	})
}
type testIsValidToken struct{
	name string
	userToTest string
	useTokenAuth bool
	testAgainstTokenAuth bool
	expectOk bool
	expectError bool
	errorType error
} 
var testIsValidTokenSet = []testIsValidToken{
	{
		name:"Test valid token for not existing user on auth vs auth token",
		userToTest: "user03",
		useTokenAuth: true,
		testAgainstTokenAuth: true,
		expectOk: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name:"Test valid token for not existing user on auth vs author token",
		userToTest: "user03",
		useTokenAuth: true,
		testAgainstTokenAuth: false,
		expectOk: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name:"Test valid token for not existing user on author vs auth token",
		userToTest: "user03",
		useTokenAuth: false,
		testAgainstTokenAuth: true,
		expectOk: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name:"Test valid token for not existing user on author vs auth token",
		userToTest: "user03",
		useTokenAuth: false,
		testAgainstTokenAuth: false,
		expectOk: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name:"Test valid token for existing user on auth vs auth token",
		userToTest: "user01",
		useTokenAuth: true,
		testAgainstTokenAuth: true,
		expectOk: true,
		expectError: false,
		errorType: nil,
	},
	{
		name:"Test invalid token for  existing user on auth vs author token",
		userToTest: "user01",
		useTokenAuth: true,
		testAgainstTokenAuth: false,
		expectOk: false,
		expectError: false,
		errorType: nil,
	},
	{
		name:"Test invalid token for  existing user on author vs auth token",
		userToTest: "user01",
		useTokenAuth: false,
		testAgainstTokenAuth: true,
		expectOk: false,
		expectError: false,
		errorType: nil,
	},
	{
		name:"Test valid token for existing user on author vs author token",
		userToTest: "user01",
		useTokenAuth: false,
		testAgainstTokenAuth: false,
		expectOk: true,
		expectError: false,
		errorType: nil,
	},
	{
		name:"Test valid token for existing WRONG user on auth vs auth token",
		userToTest: "user02",
		useTokenAuth: true,
		testAgainstTokenAuth: true,
		expectOk: false,
		expectError: false,
		errorType: nil,
	},
	{
		name:"Test invalid token for  existing WRONG user on auth vs author token",
		userToTest: "user02",
		useTokenAuth: true,
		testAgainstTokenAuth: false,
		expectOk: false,
		expectError: false,
		errorType: nil,
	},
	{
		name:"Test invalid token for  existing WRONG user on author vs auth token",
		userToTest: "user02",
		useTokenAuth: false,
		testAgainstTokenAuth: true,
		expectOk: false,
		expectError: false,
		errorType: nil,
	},
	{
		name:"Test valid token for existing WRONG user on author vs auth token",
		userToTest: "user02",
		useTokenAuth: false,
		testAgainstTokenAuth: false,
		expectOk: false,
		expectError: false,
		errorType: nil,
	},
}
func TestIsValidTokenForUser(t *testing.T){
	var testApp App
	testApp.InitDB(&testDatabase)
	user,err:=testApp.AddUser(context.Background(), common.UserData{Username:"user01"}, "Abcd")
	if err!=nil{
		t.Fatalf("Cannot create user01 user")
	}
	user2,err:=testApp.AddUser(context.Background(), common.UserData{Username:"user02"}, "Abcd")
	if err!=nil{
		t.Fatalf("Cannot create user02 user")
	}
	t.Cleanup(func(){
		testApp.db.RemoveUser(context.Background(), *user)
		testApp.db.RemoveUser(context.Background(), *user2)
	})
	auth, author, err:=testApp.TryAuthenticate(context.Background(), *user,"Abcd")
	if err!=nil{
		t.Fatalf("Error in try authentication,not expected")
	}
	for _,test:=range(testIsValidTokenSet){
		t.Run(test.name, func(t *testing.T){
			tokenToUse:=auth
			if !test.useTokenAuth{
				tokenToUse=author
			}

			ok, err:=testApp.IsValidTokenForUser(context.Background(),test.userToTest,tokenToUse, !test.testAgainstTokenAuth)
			if test.expectError{
				if err==nil{
					t.Fatal("Expected error, but not got one")
				}
				if !errors.As(err, &test.errorType){
					t.Fatalf("Wrong error type, expected %T error, but got %T", test.errorType, err)
				}
				if ok!=test.expectOk{
					t.Fatalf("Ok, expected %v but got %v", test.expectOk, ok)
				}
				
				return
			}
			if err!=nil{
				t.Fatalf("Got error, but expected not: %s", err.Error())
			}
			if ok!=test.expectOk{
				t.Fatalf("Ok, expected %v but got %v", test.expectOk, ok)
			}
		})
	}

	ok, err:=testApp.IsValidTokenForUser(context.Background(),"user01",author, true)
	if !ok || err!=nil{
		t.Fatalf("Error in check authorization token,not expected")
	}
	ok, err=testApp.IsValidTokenForUser(context.Background(),"user01",auth, false)
	if !ok || err!=nil{
		t.Fatalf("Error in check authentication token,not expected")
	}
}
func TestRemoveUserAuth(t *testing.T){
	var testApp App
	testApp.InitDB(&testDatabase)
	user,err:=testApp.AddUser(context.Background(), common.UserData{Username:"user01"}, "Abcd")
	if err!=nil{
		t.Fatalf("Cannot create test user")
	}
	t.Cleanup(func(){
		testApp.db.RemoveUser(context.Background(), *user)
	})
	auth, author, err:=testApp.TryAuthenticate(context.Background(), *user,"Abcd")
	if err!=nil{
		t.Fatalf("Error in try authentication,not expected")
	}
	ok, err:=testApp.IsValidTokenForUser(context.Background(),"user01",author, true)
	if !ok || err!=nil{
		t.Fatalf("Error in check authorization token,not expected")
	}
	ok, err=testApp.IsValidTokenForUser(context.Background(),"user01",auth, false)
	if !ok || err!=nil{
		t.Fatalf("Error in check authentication token,not expected")
	}
	t.Run("Remove auth of a not existing user", func(t *testing.T){
		err=testApp.RemoveUserAuthentication(context.Background(), "user02", auth)
		if err==nil{
			t.Fatal("Expected error, but no")
		}
		if !errors.As(err, &common.NotFoundError{}){
			t.Fatalf("Expected not found error, but no: %s", err.Error())
		}
	})
	t.Run("Remove auth of the existing user", func(t *testing.T){
		err=testApp.RemoveUserAuthentication(context.Background(), "user01", auth)
		if err!=nil{
			t.Fatal("Not expected error but got one")
		}
		ok, _:=testApp.IsValidTokenForUser(context.Background(),"user01",author, true)

		if ok {
			t.Fatalf("Error in check authorization token,not expected")
		}
		ok, err=testApp.IsValidTokenForUser(context.Background(),"user01",auth, false)
		if ok {
			t.Fatalf("Error in check authentication token,not expected")
		}
	})
}

func TestChangePassword(t *testing.T){
	var testApp App
	testApp.InitDB(&testDatabase)
	user,err:=testApp.AddUser(context.Background(), common.UserData{Username:"user01"}, "Abcd")
	if err!=nil{
		t.Fatalf("Cannot create test user")
	}
	t.Cleanup(func(){
		testApp.db.RemoveUser(context.Background(), *user)
	})
	t.Run("Test change password of the user", func(t *testing.T){
		_, err:=testApp.ChangePassword(context.Background(), "user01", "NewPassword")
		if err!=nil{
			t.Fatalf("Got error, should not: %s", err.Error())
		}
		//Check old password
		err=testApp.checkUserPassword(context.Background(), "user01", "Abcd")
		if err==nil{
			t.Fatalf("Got no error, should have because old password")
		}
		//Check new password
		err=testApp.checkUserPassword(context.Background(), "user01", "NewPassword")
		if err!=nil{
			t.Fatalf("Got error, should not: %s", err.Error())
		}
	})
	t.Run("Test change password of another user", func(t *testing.T){
		_, err:=testApp.ChangePassword(context.Background(), "user02", "NewPassword")
		if err==nil{
			t.Fatalf("Got no error, should have")
		}
		if !errors.As(err, &common.NotFoundError{}){
			t.Fatalf("Got wrong error,got %T should have %T", err, common.NotFoundError{})
		}
	})
}
