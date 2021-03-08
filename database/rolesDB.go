package database

type SecurityRole struct{
	Name string
	privileges []SecurityPrivilege
}

func (role SecurityRole)HasPrivilege(privilege SecurityPrivilege) bool{
	for _, secPrivilege := range(role.privileges){
		if privilege==secPrivilege{
			return true
		}
	}
	return false
}
