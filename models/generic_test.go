package models

import (
	"context"
	"testing"

	"github.com/idalmasso/clubmngserver/common"
	"github.com/idalmasso/clubmngserver/database"
	"github.com/idalmasso/clubmngserver/database/memdb"
)
var databasesToTest =[]database.ClubDb{ &memdb.MemoryDB{} }
func TestInitDB(t *testing.T){
	t.Log("[TestInitDB] - Start")
	for _,database:=range(databasesToTest){
		InitDB(&database)
	}
}

func TestInitDBAndUser(t *testing.T){
	t.Log("[TestInitDBAndAddedUsers] - Start")
	for _,db:=range(databasesToTest){
		InitDB(&db)
		roles, err:=GetAllRoles(context.Background())
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
				t.Error("[TestInitDBAndAddedUsers] - Cannot find role "+k+" in added roles")
				return
			}
			if !unorderedEqualPrivileges(v, rolesToBeAdded[k]){
				t.Error("[TestInitDBAndAddedUsers] - Cannot find privileges for role "+k+" in added roles")
				return
			}

		}
		users,err := db.GetAllUsers(context.Background())
		if err!=nil{
			t.Error("[TestInitDBAndAddedUsers] - Cannot find users "+err.Error())
			return
		}
		for username, passRole:=range(usersToBeAdded){
			found:=false
			for _,u:=range(users){
				if u.Username==username{
					found=true
					if u.Role!=passRole.role{
						t.Error("[TestInitDBAndAddedUsers] - User " +username+" has wrong role")
						return
					}
				}
				
			}
			if !found{
				t.Error("[TestInitDBAndAddedUsers] - User " +username+" not found")
			}
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
