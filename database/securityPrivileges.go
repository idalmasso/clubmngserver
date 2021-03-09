package database

import "fmt"

//SecurityPrivilege is the type that is incorporated in a role to view if a user can or not do things
type SecurityPrivilege int;

var securityPrivilegesStringMap  = map[SecurityPrivilege]string{
	SecurityAdmin: "SecurityAdmin",
	SecuritySelfUserView: "SecuritySelfUserView",
	SecuritySelfUserUpdate:"SecuritySelfUserUpdate",
	SecurityLinkedUsersView:"SecurityLinkedUsersView",
	SecurityLinkedUsersUpdate:"SecurityLinkedUsersUpdate",
	SecurityLinkedUsersDelete:"SecurityLinkedUsersDelete",
	SecurityAllUsersUpdate:"SecurityAllUsersUpdate",
	SecurityAllUsersView:"SecurityAllUsersView",
	SecurityAllUsersDelete:"SecurityAllUsersDelete",
	SecurityRolesView:"SecurityRolesView",
	SecurityRolesUpdate:"SecurityRolesUpdate",
	SecurityRolesDelete:"SecurityRolesDelete",
	SecurityRolesToUserMaintain:"SecurityRolesToUserMaintain",
	SecuritySelfPaymentsView:"SecuritySelfPaymentsView",
	SecurityLinkedUsersPaymentsView:"SecurityLinkedUsersPaymentsView",
	SecurityAllUsersPaymentsView:"SecurityAllUsersPaymentsView",
	SecurityAddUserEntrances:"SecurityAddUserEntrances",
	SecurityRemoveUserEntrances:"SecurityRemoveUserEntrances",
	SecurityUpdateParameters:"SecurityUpdateParameters",
	SecurityViewParameters:"SecurityViewParameters",
}
const  (
	//SecurityAdmin is all ok for any function
	SecurityAdmin SecurityPrivilege=iota
	//SecuritySelfUserView only view self user
	SecuritySelfUserView
	//SecuritySelfUserUpdate can update only self
	SecuritySelfUserUpdate
	//SecuritySelfUserDelete can delete only self
	SecuritySelfUserDelete
	//SecurityLinkedUsersView can view linked users
	SecurityLinkedUsersView
	//SecurityLinkedUsersUpdate can update linked users
	SecurityLinkedUsersUpdate
	//SecurityLinkedUsersDelete can delete linked users
	SecurityLinkedUsersDelete
	//SecurityAllUsersUpdate can Update all users
	SecurityAllUsersUpdate
	//SecurityAllUsersView Can View all users
	SecurityAllUsersView
	//SecurityAllUsersDelete Can Delete all users
	SecurityAllUsersDelete
	//SecurityRolesView can view security roles
	SecurityRolesView
	//SecurityRolesUpdate can update security roles, add remove privileges
	SecurityRolesUpdate
	//SecurityRolesDelete can delete a security role, only if no user has it
	SecurityRolesDelete
	//SecurityRolesToUserMaintain can add and remove roles to a user
	SecurityRolesToUserMaintain
	//SecuritySelfPaymentsView can view own payment history
	SecuritySelfPaymentsView
	//SecurityLinkedUsersPaymentsView can view linked users payment history
	SecurityLinkedUsersPaymentsView
	//SecurityAllUsersPaymentsView can view all users payment history
	SecurityAllUsersPaymentsView
	//SecurityAddUserEntrances Can add user entrance for a user
	SecurityAddUserEntrances
	//SecurityRemoveUserEntrances Can remove user entrance for a user 
	SecurityRemoveUserEntrances
	//SecurityUpdateParameters Can update parameters
	SecurityUpdateParameters
	//SecurityViewParameters Can view the parameters
	SecurityViewParameters
)


//String value representing the privilege
func (securityPrivilege SecurityPrivilege) String() (string,error){
	v, ok:= securityPrivilegesStringMap[securityPrivilege]
	if !ok{
		return "",fmt.Errorf("%v is not a valid privilege", securityPrivilege)
	}
	return v, nil
}
//StringToSecurityPrivilege returns a securityPrivilege value from a string that represent it
func StringToSecurityPrivilege(security string) (SecurityPrivilege, error){
	for k, v := range(securityPrivilegesStringMap){
		if v==security{
			return k, nil
		}
	}
	return -1, fmt.Errorf("%v is not a valid privilege", security)
}
//ListPrivileges returns a list of all existing privileges in strings
func ListPrivileges() []string{
	var listP []string
	for _, v:=range(securityPrivilegesStringMap){
		listP=append(listP, v)
	}
	return listP
}
