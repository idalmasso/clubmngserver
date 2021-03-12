package memdb

import (
	"github.com/idalmasso/clubmngserver/common"
	"github.com/idalmasso/clubmngserver/common/userprops"
)

//MemoryDB is the in memory database for testing
type MemoryDB struct{
	users map[string]*common.UserData
	roles map[string]*common.SecurityRole
	userProps map[string]userprops.UserPropertyDefinition
}

//Init iinitialize the memory db
func (db *MemoryDB) Init(){

	if db.roles==nil{
		db.roles=make(map[string]*common.SecurityRole)
	
	}
	if db.users==nil{
		db.users=make(map[string]*common.UserData)
	}
}
