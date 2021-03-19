package app

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/idalmasso/clubmngserver/common"
	"github.com/idalmasso/clubmngserver/database"
	"github.com/idalmasso/clubmngserver/database/memdb"
)

var testDatabase database.ClubDb
func TestMain(m *testing.M){
	var code int
	
	for _,clubDB:=range(databasesToTest){
		testDatabase=clubDB
		
		log.Printf("Using database %T\n", testDatabase)
		code=m.Run()
		log.Printf("Returned code %d\n", code)
		clubDB.TearDown()
		clubDB=nil
		testDatabase=nil
	}
	os.Exit(code)
}
var databasesToTest =[]database.ClubDbTest{ &memdb.MemoryDBTest{} }
func TestInitDB(t *testing.T){
	var testApp App
	testApp.InitDB(&testDatabase)
}

func TestInitDBAndUser(t *testing.T){
	var testApp App
	testApp.InitDB(&testDatabase)
	roles, err:=testApp.GetAllRoles(context.Background())
	if err!=nil{
		t.Error(err)
		return
	}
	addedRoles:=make(map[string][]common.SecurityPrivilege)
	for _,role:=range(roles){
		addedRoles [role.Name]=role.Privileges
	}
	for k,v:=range(addedRoles){
		if _, ok:=rolesToBeAdded[k]; !ok{
			t.Fatalf("Cannot find role %s in added roles", k)
		}
		if !unorderedEqualPrivileges(v, rolesToBeAdded[k]){
			t.Fatalf("Cannot find privileges for role %s in added roles",k)
		}

	}
	users,err := testApp.GetUsersList(context.Background())
	if err!=nil{
		t.Fatalf("Cannot find users %v",err.Error())
	}
	for username, passRole:=range(usersToBeAdded){
		found:=false
		for _,u:=range(users){
			if u.Username==username{
				found=true
				if u.Role!=passRole.role{
					t.Fatalf("User %s has wrong role", username)
					return
				}
			}
			
		}
		if !found{
			t.Fatalf("User %s not found",username)
		}
	}
}


func unorderedEqualPrivileges(first, second []common.SecurityPrivilege) bool {
    if len(first) != len(second) {
        return false
    }
    exists := make(map[common.SecurityPrivilege]bool)
    for _, value := range first {
        exists[value] = true
    }
    for _, value := range second {
        if !exists[value] {
            return false
        }
    }
    return true
}
