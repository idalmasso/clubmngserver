package database

//ClubDb is the interface that gives all function to be used from the "model"
type ClubDb interface{
	SecurityRoleDBInterface
	UserAuthDBInterface
	UserPropsDB
	Init()
	
}
