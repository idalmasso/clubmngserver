package models

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/idalmasso/clubmngserver/common"
)
type addRoleTest struct{
	testName string
	testNumber int
	roleToAdd string
	privileges []common.SecurityPrivilege
	resultError bool
	errorType error
	
}
var addRoleTestSet= []addRoleTest{
	{testName:"Test add admin user again",testNumber:1, roleToAdd: "Admin", privileges:[]common.SecurityPrivilege{common.SecurityAdmin} , resultError: true, errorType: common.AlreadyExistsError{} },
	{testName:"Test normal user add",testNumber:2, roleToAdd: "NormalUserRole", privileges:[]common.SecurityPrivilege{common.SecuritySelfUserView, common.SecuritySelfUserUpdate, common.SecuritySelfPaymentsView} , resultError: false, errorType: nil },
	{testName:"Test add new role with other privilege",testNumber:3, roleToAdd: "TestRole2", privileges:[]common.SecurityPrivilege{common.SecuritySelfUserView, common.SecuritySelfUserUpdate, common.SecuritySelfPaymentsView} , resultError: false, errorType: nil },
}


func TestAddRole(t *testing.T){
	t.Logf("[TestAddRole] - start")
	for _,db:=range(databasesToTest){
		InitDB(&db)
		for _, test:=range(addRoleTestSet){
			t.Run(test.testName, func (t *testing.T){
				role, err:=AddRole(context.Background(), test.roleToAdd, test.privileges...)
				if  test.resultError{
					if err==nil {
						t.Errorf("[TestAddRole] - %d - expected test to give error but not", test.testNumber)
						if role!=nil{
							t.Cleanup(func(){
								db.RemoveRole(context.Background(), test.roleToAdd)
							})
						}
					} else if !errors.As(err, &test.errorType){
						t.Errorf("[TestAddRole] - %d - expected test to give %t type of error, got %t", test.testNumber, test.errorType, err)
					}
				}else{
					if err!=nil{
						t.Errorf("[TestAddRole] - %d - expected test to not give error but it has given", test.testNumber)
						return
					} else if role.Name!=test.roleToAdd {
						t.Errorf("[TestAddRole] - %d - wrong inserted role name, %s instead of %s", test.testNumber, role.Name, test.roleToAdd)
					} else if !unorderedEqualPrivileges(role.Privileges, test.privileges){
						t.Errorf("[TestAddRole] - %d - wrong inserted privileges", test.testNumber)
					}
					t.Cleanup(func(){
						db.RemoveRole(context.Background(), test.roleToAdd)
					})
				}
			})
		}
	}
}
type deleteRoleTest struct{
	testName string
	testNumber int
	rolesToAdd []string
	privileges map[string][]common.SecurityPrivilege
	roleToDelete string
	resultError bool
	errorType error
	
}
var deleteRoleTestSet= []deleteRoleTest{
	{testNumber:1, 
		rolesToAdd: []string{"role1"}, 
		privileges:map[string][]common.SecurityPrivilege{ "role1":{common.SecurityAllUsersDelete}} , 
		roleToDelete: "role1",
		resultError: false, 
		errorType: nil },
	{testNumber:2, 
		rolesToAdd: []string{"role2"}, 
		privileges:map[string][]common.SecurityPrivilege{ "role2":{common.SecurityAllUsersDelete}} , 
		roleToDelete: "role1",
		resultError: true, 
		errorType: common.NotFoundError{} },
	{testNumber:3, 
		rolesToAdd: []string{"role3"}, 
		privileges:map[string][]common.SecurityPrivilege{ "role1":{common.SecurityAllUsersDelete}} , 
		roleToDelete: "role3",
		resultError: false, 
		errorType: nil },
}
func TestDeleteRole(t *testing.T){
	t.Logf("[TestDeleteRole] - start")
	for _,db:=range(databasesToTest){
		
		InitDB(&db)
		for _, test:=range(deleteRoleTestSet){
			t.Run(test.testName, func (t *testing.T){
				for _, testRole:=range(test.rolesToAdd){
					_, err:=AddRole(context.Background(), testRole, test.privileges[testRole]...)
					if err!=nil {
						t.Errorf("[TestDeleteRole] - %d - cannot add role", test.testNumber)
						continue
					}
					
				}
				t.Cleanup(func(){
					for _, testRole:=range(test.rolesToAdd){
						db.RemoveRole(context.Background(),testRole)
					}
				})
				err:=DeleteRole(context.Background(), test.roleToDelete)
				if test.resultError{
					if err==nil{
						t.Errorf("[TestToDeleteRole] - %d expected test to give error, does not", test.testNumber)
					} else if !errors.As(err, &test.errorType){
						t.Errorf("[TestToDeleteRole] - %d expected test to give error %v, does give %v", test.testNumber,test.errorType, err )
					}
				} else {
					if err!=nil{
						t.Errorf("[TestDeleteRole] - %d - expected test to not give error but it has given", test.testNumber)
					} else  {
						role, err:=GetRole(context.Background(), test.roleToDelete)
						var errorNotFound common.NotFoundError
						if err==nil || !errors.As(err, &errorNotFound){
							t.Errorf("[TestDeleteRole] - %d - expected find after delete to give not found error", test.testNumber)
						} 
						if role!=nil{
							t.Errorf("[TestDeleteRole] - %d - expected find to return null user", test.testNumber)
						}

					}
				}
			})
		}
	
	}
}

func TestDeleteRoleExistsUser(t *testing.T){
	t.Logf("[TestDeleteRoleExistsUser] - start")
	for _,db:=range(databasesToTest){
		
		InitDB(&db)

		_, err:=AddRole(context.Background(), "roleToTest", []common.SecurityPrivilege{common.SecurityAddUserEntrances}...)
		if err!=nil {
			t.Errorf("[TestDeleteRoleExistsUser] -  cannot add role: %v", err)
			return
		}
	
		user,err:=AddUser(context.Background(), common.UserData{Username: "Pippo"}, "Abcd")
		if err!=nil {
			t.Errorf("[TestDeleteRoleExistsUser] -  cannot add user: %v", err)
			t.Cleanup(func(){
				db.RemoveRole(context.Background(),"roleToTest")
			})
			return
		}
		t.Cleanup(func(){
			db.RemoveRole(context.Background(),"roleToTest")
			db.RemoveUser(context.Background(), *user)
		})
		user, err=AddRoleToUser(context.Background(),user.Username,"roleToTest")
		if err!=nil {
			t.Errorf("[TestDeleteRoleExistsUser] -  cannot add role to user: %v",err)
			return
		}
		err=DeleteRole(context.Background(), "roleToTest")
		
		if err==nil || !strings.Contains(err.Error(), "Users have the role"){
			t.Errorf("[TestDeleteRoleExistsUser] - expected test to give error with 'users have the role', does not")
		}
	}
}

func TestGetAllRoles(t *testing.T){
	t.Logf("[TestGetAllRoles] - start")
	for _,db:=range(databasesToTest){
		
		InitDB(&db)
		roles, err:=GetAllRoles(context.Background())
		if err!=nil{
			t.Errorf("[TestGetAllRoles] - cannot get all roles after initialization")
		} else if len(roles)!=len(rolesToBeAdded){
			t.Errorf("[TestGetAllRoles] - len roles =%d, expected %d role after initialization : %v", len(roles), len(rolesToBeAdded),roles)
		} else{
			for _,role:=range(roles){
				_, ok:=rolesToBeAdded[role.Name]
				if !ok{
					t.Fatalf("[TestGetAllRoles] - find role %s in db roles but not in initialized ones", role.Name)
				}
				if !unorderedEqualPrivileges( role.Privileges, rolesToBeAdded[role.Name]){
					t.Fatalf("[TestGetAllRoles] - privileges expected different then in role %s", role.Name)
				}
			}
			_,err=AddRole(context.Background(), "role1", common.SecurityAllUsersPaymentsView)
			t.Cleanup(func(){
				db.RemoveRole(context.Background(), "role1")
			})
			if err!=nil{
				t.Errorf("[TestGetAllRoles] - cannot add one role after initialization")
			}
			roles, err:=GetAllRoles(context.Background())
			if err!=nil{
				t.Errorf("[TestGetAllRoles] - cannot get all roles after initialization and added 1")
			} else if len(roles)!=len(rolesToBeAdded)+1{
				t.Errorf("[TestGetAllRoles] - len roles =%d, expected %d role after initialization and added 1", len(roles), len(rolesToBeAdded)+1)
			} else{
				for _,role:=range(roles){
					_, ok:=rolesToBeAdded[role.Name]
					if !ok && role.Name!="role1"{
						t.Errorf("[TestGetAllRoles] - find role %s in db roles but not in initialized ones", role.Name)
						return
					}
					if role.Name=="role1"{
						if !unorderedEqualPrivileges( role.Privileges, []common.SecurityPrivilege{common.SecurityAllUsersPaymentsView}){
							t.Errorf("[TestGetAllRoles] - privileges expected different then in role %s", role.Name)
							return
						}
					}
				}
			}
		}
	}
}

func TestGetSingleRole(t *testing.T){
	t.Logf("[TestGetSingleRole] - start")
	for _,db:=range(databasesToTest){
		
		InitDB(&db)
		t.Run("Get all roles initialization one by one", func(t *testing.T){
			for roleName, privileges:=range(rolesToBeAdded){
				role, err:=GetRole(context.Background(), roleName)
				if err!=nil{
					t.Fatalf("Cannot get role %s", roleName)
				}
				if role==nil{
					t.Fatalf("Role %s nil", roleName)
				}
				if !unorderedEqualPrivileges(privileges, role.Privileges){
					t.Fatalf("Privileges on role %s are not equal to the ones expected, %v - %v", roleName, role.Privileges,privileges)
				}
			}
		})
		t.Run("Get one role not existent", func(t *testing.T){
			role, err:=GetRole(context.Background(), "NotExistsRole")
			var notFoundError common.NotFoundError
			if err==nil{
				t.Errorf("Error get Role nil")
			} else if !errors.As(err, &notFoundError){
				t.Errorf("Error get Role is not of correct type")
			} else if role != nil{
				t.Errorf("Role is not nil")
			}
		})
		t.Run("Get one role after having added it", func(t *testing.T){
			var privs = []common.SecurityPrivilege{common.SecurityAllUsersDelete, common.SecurityLinkedUsersPaymentsView, common.SecurityRolesAdd}
			_,err:=AddRole(context.Background(), "ExistsRole",privs... )
			if err!=nil{
				t.Fatal("Cannot run test, cannot add role ExistsRole")
			}
			t.Cleanup(func(){
				db.RemoveRole(context.Background(), "ExistsRole")
			})
			role, err:=GetRole(context.Background(), "ExistsRole")
			if err!=nil{
				t.Fatalf("Error in GetRole: %v", err)
			}
			if role.Name!="ExistsRole"{
				t.Fatalf("Role name is '%s', not 'ExistsRole'", role.Name)
			}
			if !unorderedEqualPrivileges(role.Privileges, privs){
				t.Fatalf("Privileges on role %s are not equal to the ones expected, %v - %v", role.Name,  role.Privileges,privs)
			}
		})
	}
}

type updateRoleTest struct{
	testName string
	roleToAdd string
	roleToUpdate string
	privilegesBefore []common.SecurityPrivilege
	privilegesAfter []common.SecurityPrivilege
	resultError bool
	errorType error
}
var updateRoleTestSet= []updateRoleTest{
	{testName:"Test update role not exists",
		roleToAdd: "role01", 
		roleToUpdate:"role02",
		privilegesBefore:[]common.SecurityPrivilege{common.SecurityAdmin} , 
		privilegesAfter: []common.SecurityPrivilege{common.SecurityAdmin} , 
		resultError: true, 
		errorType: common.NotFoundError{} },
	{testName:"Test update role exists, with that privilege and others ",
		roleToAdd: "role01", 
		roleToUpdate:"role01",
		privilegesBefore:[]common.SecurityPrivilege{common.SecurityAdmin} , 
		privilegesAfter: []common.SecurityPrivilege{common.SecurityAdmin, common.SecurityAddUserEntrances, common.SecuritySelfUserUpdate, common.SecurityRolesAdd} , 
		resultError: false, 
		errorType: nil},
	{testName:"Test update user exists, without that privilege  ",
		roleToAdd: "role01", 
		roleToUpdate:"role01",
		privilegesBefore:[]common.SecurityPrivilege{common.SecurityAdmin} , 
		privilegesAfter: []common.SecurityPrivilege{common.SecurityAddUserEntrances, common.SecuritySelfUserUpdate, common.SecurityRolesAdd} , 
		resultError: false, 
		errorType: nil},
	{testName:"Test update user exists, with no privileges",
		roleToAdd: "role01", 
		roleToUpdate:"role01",
		privilegesBefore:[]common.SecurityPrivilege{common.SecurityAdmin} , 
		privilegesAfter: []common.SecurityPrivilege{} , 
		resultError: false, 
		errorType: nil},
	{testName:"Test update user exists, with nil list",
		roleToAdd: "role01", 
		roleToUpdate:"role01",
		privilegesBefore:[]common.SecurityPrivilege{common.SecurityAdmin} , 
		privilegesAfter: nil , 
		resultError: false, 
		errorType: nil},
	}
func TestUpdateRoles(t *testing.T){
		t.Logf("[TestUpdateRoles] - start")
	for _,db:=range(databasesToTest){
		InitDB(&db)
		for _, test:=range(updateRoleTestSet){
			t.Run(test.testName, func (t *testing.T){
				_, err:=AddRole(context.Background(), test.roleToAdd, test.privilegesBefore...)
				t.Cleanup(func(){
					db.RemoveRole(context.Background(), test.roleToAdd)
				})
				if err!=nil{
					t.Fatal("Cannot add the role")
				}
				role, err:=UpdateRole(context.Background(), test.roleToUpdate, test.privilegesAfter...)
				if test.resultError{
					if err==nil{
						t.Fatal("Expected error, not given")
					}
					if !errors.As(err, &test.errorType){
						t.Fatalf("Wrong expected errortype, got %v, expected %v", err, test.errorType)
					}
					return
				}
				if !unorderedEqualPrivileges(role.Privileges, test.privilegesAfter){
					t.Fatalf("Error: expected role to have privileges %v, it has %v", test.privilegesAfter, role.Privileges)
				}
				role, err=GetRole(context.Background(), test.roleToUpdate)
				if err!=nil{
					t.Fatalf("Cannot find anymore role %s", test.roleToUpdate)
				}
				if !unorderedEqualPrivileges(role.Privileges, test.privilegesAfter){
					t.Fatalf("Error: expected role to have privileges %v, it has %v", test.privilegesAfter, role.Privileges)
				}
			})
		}
	}
}

type addRoleToUserTest struct{
	testName string
	userToCreate string
	userToUse string
	roleToCreate string
	roleToAdd string
	privileges []common.SecurityPrivilege
	resultError bool
	errorType error
}
var addRoleToUserTestSet= []addRoleToUserTest{
	{testName:"Add role to not exists user",
		roleToCreate:"role01",
		roleToAdd: "role01", 
		userToCreate:"user01",
		userToUse: "user02",
		privileges:[]common.SecurityPrivilege{common.SecurityAdmin} , 
		resultError: true, 
		errorType: common.NotFoundError{} },
	{testName:"Add not exists role to not exists user",
		roleToCreate:"role01",
		roleToAdd: "role02", 
		userToCreate:"user01",
		userToUse: "user02",
		privileges:[]common.SecurityPrivilege{common.SecurityAdmin} , 
		resultError: true, 
		errorType: common.NotFoundError{} },
	{testName:"Add not exists role to exists user",
		roleToCreate:"role01",
		roleToAdd: "role02", 
		userToCreate:"user01",
		userToUse: "user01",
		privileges:[]common.SecurityPrivilege{common.SecurityAdmin} , 
		resultError: true, 
		errorType: common.NotFoundError{} },
	{testName:"Add exists role to exists user",
		roleToCreate:"role01",
		roleToAdd: "role01", 
		userToCreate:"user01",
		userToUse: "user01",
		privileges:[]common.SecurityPrivilege{common.SecurityAdmin} , 
		resultError: false, 
		errorType: nil },
	}
func TestAddRolesToUser(t *testing.T){
	t.Logf("[TestAddRolesToUser] - start")
	for _,db:=range(databasesToTest){
		InitDB(&db)
		for _, test:=range(addRoleToUserTestSet){
			t.Run(test.testName, func (t *testing.T){
				_, err:=AddRole(context.Background(), test.roleToCreate, test.privileges...)
				_, err2 := AddUser(context.Background(), common.UserData{Username: test.userToCreate}, "Abcd")

				t.Cleanup(func(){
					db.RemoveUser(context.Background(), common.UserData{Username:test.userToCreate})
					db.RemoveRole(context.Background(), test.roleToCreate)
				})
				if err!=nil{
					t.Fatal("Cannot add the role")
				}
				if err2!=nil{
					t.Fatal("Cannot add the user")
				}
				
				user, err:=AddRoleToUser(context.Background(), test.userToUse, test.roleToAdd)
				if test.resultError{
					if err==nil{
						t.Fatal("Expected error, not given")
					}
					if !errors.As(err, &test.errorType){
						t.Fatalf("Wrong expected errortype, got %v, expected %v", err, test.errorType)
					}
					return
				}
				
				if user.Role!=test.roleToAdd{
					t.Fatalf("Error: expected role on user to be %v, it has %v", test.roleToAdd, user.Role)
				}
				user, err = FindUser(context.Background(), test.userToUse)
				if err!=nil{
					t.Fatal("Cannot find user at end, strange error")
				}
				if user==nil{
					t.Fatalf("Cannot find anymore user %s", test.userToUse)
				}
					if user.Role!=test.roleToAdd{
					t.Fatalf("Error after find: expected role on user to be %v, it has %v", test.roleToAdd, user.Role)
				}
			})
		}
	}
}


type isRoleAuthorizedTest struct{
	testName string
	roleToUse string
	roleToAdd string
	rolePrivileges []common.SecurityPrivilege
	testPrivilege common.SecurityPrivilege
	resultValue bool
	resultError bool
	errorType error
	
}
var isRoleAuthorizedTestSet= []isRoleAuthorizedTest{
	{testName:"Role has not the privilege",
		roleToUse:"role01",
		roleToAdd:"role01",
		rolePrivileges:[]common.SecurityPrivilege{common.SecurityAllUsersDelete, common.SecurityAllUsersView, common.SecurityLinkedUsersDelete} , 
		testPrivilege: common.SecurityRolesAdd,
		resultValue: false,
		resultError: false, 
		errorType: nil },
	{testName:"Role has the privilege 1",
		roleToUse:"role01",
		roleToAdd:"role01",
		rolePrivileges:[]common.SecurityPrivilege{common.SecurityAllUsersDelete, common.SecurityAllUsersView, common.SecurityLinkedUsersDelete} , 
		testPrivilege: common.SecurityAllUsersDelete,
		resultValue: true,
		resultError: false, 
		errorType: nil },
	{testName:"Role has the privilege 2",
		roleToUse:"role01",
		roleToAdd:"role01",
		rolePrivileges:[]common.SecurityPrivilege{common.SecurityAllUsersDelete, common.SecurityAllUsersView, common.SecurityLinkedUsersDelete} , 
		testPrivilege: common.SecurityAllUsersView,
		resultValue: true,
		resultError: false, 
		errorType: nil },
	{testName:"Role has the privilege 3",
		roleToUse:"role01",
		roleToAdd:"role01",
		rolePrivileges:[]common.SecurityPrivilege{common.SecurityAllUsersDelete, common.SecurityAllUsersView, common.SecurityLinkedUsersDelete} , 
		testPrivilege: common.SecurityLinkedUsersDelete,
		resultValue: true,
		resultError: false, 
		errorType: nil },
	{testName:"Role has the privilege SysAdmin",
		roleToUse:"role01",
		roleToAdd:"role01",
		rolePrivileges:[]common.SecurityPrivilege{common.SecurityAdmin, common.SecurityAllUsersView, common.SecurityLinkedUsersDelete} , 
		testPrivilege: common.SecurityViewParameters,
		resultValue: true,
		resultError: false, 
		errorType: nil },
	{testName:"Role not exist",
		roleToUse:"role02",
		roleToAdd:"role01",
		rolePrivileges:[]common.SecurityPrivilege{common.SecurityAdmin, common.SecurityAllUsersView, common.SecurityLinkedUsersDelete} , 
		testPrivilege: common.SecurityViewParameters,
		resultValue: false,
		resultError: true, 
		errorType: common.NotFoundError{} },
}
func TestIsRoleAuthorized(t *testing.T){
	t.Logf("[TestIsRoleAuthorized] - start")
	for _,db:=range(databasesToTest){
		InitDB(&db)
		for _, test:=range(isRoleAuthorizedTestSet){
			t.Run(test.testName, func (t *testing.T){
				_, err:=AddRole(context.Background(), test.roleToAdd, test.rolePrivileges...)
				t.Cleanup(func(){
					db.RemoveRole(context.Background(), test.roleToAdd)
				})
				if err!=nil{
					t.Fatal("Cannot add the role")
				}
				auth, err:=IsRoleAuthorized(context.Background(), test.roleToUse, test.testPrivilege)
				if test.resultError{
					if auth!=test.resultValue{
						t.Errorf("Expected result %v, got %v", test.resultValue, auth)
					}
					if err==nil{
						t.Fatal("Expected error, not got one")
					}
					if !errors.As(err, &test.errorType){
						t.Fatalf("Wrong expected errortype, got %v, expected %v", err, test.errorType)
					}
				} else {
					if err!=nil{
						t.Fatalf("Not expected error, got one: %v", err)
					}
					if auth!=test.resultValue{
						t.Fatalf("Expected result %v, got %v", test.resultValue, auth)
					}
				}
			})
		}
	}
}

func TestGetUserRole(t *testing.T){
	t.Logf("[TestGetUserRole] - start")
	for _,db:=range(databasesToTest){
		InitDB(&db)
		t.Run("Get user role exists role", func (t *testing.T){
			_, err:=AddRole(context.Background(), "role01", common.SecurityAdmin)

			t.Cleanup(func(){
				db.RemoveRole(context.Background(), "role01")
			})
			if err!=nil{
				t.Fatal("Cannot add the role")
			}
			role, err:=GetUserRole(context.Background(), common.UserData{Username: "user", Role: "role01"})
			if err!=nil{
				t.Fatalf("Error gotten not expected")
			}
			if role==nil{
				t.Fatalf("Role returned is nil")
			}
			if role.Name!="role01" {
				t.Fatalf("Expected role01, but got %s", role.Name)
			}
			if len(role.Privileges)!=1 || role.Privileges[0]!=common.SecurityAdmin{
				t.Fatalf("Expected only securityAdmin role, got %v", role.Privileges)
			}
		})
		t.Run("Get user role not exists role", func (t *testing.T){
		
			role, err:=GetUserRole(context.Background(), common.UserData{Username: "user", Role: "role01"})
			if err==nil{
				t.Fatalf("Error not gotten but expected")
			}
			if role!=nil{
				t.Fatalf("Role returned is not nil")
			}
			if !errors.As(err, &common.NotFoundError{}){
				t.Fatalf("Error  gotten but wrong type: %v but expected %v", err, common.NotFoundError{})
			}
		})
	}
}
