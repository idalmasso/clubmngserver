package database

//SecurityRole is the type that is linked to a user to see if he can do actions
type SecurityRole struct{
	Name string
	Privileges []SecurityPrivilege
}
//HasPrivilege returns true if a user has a security privilege
func (role SecurityRole)HasPrivilege(privilege SecurityPrivilege) bool{
	for _, secPrivilege := range(role.Privileges){
		if privilege==secPrivilege{
			return true
		}
	}
	return false
}

