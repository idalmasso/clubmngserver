package database

type SecurityRole struct{
	Name string
	privileges []SecurityPrivilege
}
