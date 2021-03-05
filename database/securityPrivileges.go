package database

type SecurityPrivilege int;
const  (
	SecurityAdmin SecurityPrivilege=iota
	SecuritySelfUserView
	SecuritySelfUserUpdate
	SecuritySelfUserDelete
	SecurityLinkedUsersView
	SecurityLinkedUsersUpdate
	SecurityLinkedUsersDelete
	SecurityAllUsersUpdate
	SecurityAllUsersView
	SecurityAllUsersDelete
	SecurityRolesView
	SecurityRolesUpdate
	SecurityRolesDelete
	SecurityRolesToUserMaintain
	SecuritySelfPaymentsView
	SecurityLinkedUsersPaymentsView
	SecurityAllUsersPaymentsView
	SecurityAddUserEntrances
	SecurityRemoveUserEntrances
	SecurityUpdateParameters
	SecurityViewParameters
)

func (securityPrivilege SecurityPrivilege) String() string{
	switch securityPrivilege{
	case SecurityAdmin:
		return "SecurityAdmin"
	case SecuritySelfUserView:
		return "SecuritySelfUserView"
	case SecuritySelfUserUpdate:
		return "SecuritySelfUserUpdate"
	case SecuritySelfUserDelete:
		return "SecuritySelfUserDelete"
	case SecurityLinkedUsersView:
		return "SecurityLinkedUsersView"
	case SecurityLinkedUsersUpdate:
		return "SecurityLinkedUsersUpdate"
	case SecurityLinkedUsersDelete:
		return "SecurityLinkedUsersDelete"
	case SecurityAllUsersUpdate:
		return "SecurityAllUsersUpdate"
	case SecurityAllUsersView:
		return "SecurityAllUsersView"
	case SecurityAllUsersDelete:
		return "SecurityAllUsersDelete"
	case SecurityRolesView:
		return "SecurityRolesView"
	case SecurityRolesUpdate:
		return "SecurityRolesUpdate"
	case SecurityRolesDelete:
		return "SecurityRolesDelete"
	case SecurityRolesToUserMaintain:
		return "SecurityRolesToUserMaintain"
	case SecuritySelfPaymentsView:
		return "SecuritySelfPaymentsView"
	case SecurityLinkedUsersPaymentsView:
		return "SecurityLinkedUsersPaymentsView"
	case SecurityAllUsersPaymentsView:
		return "SecurityAllUsersPaymentsView"
	case SecurityAddUserEntrances:
		return "SecurityAddUserEntrances"
	case SecurityRemoveUserEntrances:
		return "SecurityRemoveUserEntrances"
	case SecurityUpdateParameters:
		return "SecurityUpdateParameters"
	case SecurityViewParameters:
		return "SecurityViewParameters"
	}
	return ""
}
