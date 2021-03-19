package app

import (
	"context"
	"fmt"

	"github.com/idalmasso/clubmngserver/common"
)

//AddRole tries to add a role with name roleName and privileges passed
func (app *App)AddRole(ctx context.Context, roleName string, privileges ...common.SecurityPrivilege) ( *common.SecurityRole,error) {
	role, err:=app.db.FindRole(ctx, roleName)
	if err!=nil{
		return nil,err
	}
	if role!=nil{
		return role,  common.AlreadyExistsError{ID:roleName}
	}
	newRole:= common.SecurityRole{Name: roleName , Privileges: privileges}
	theNewRole, err := app.db.AddRole(ctx, newRole); 
	if err!=nil{
		return nil,err
	}
	return theNewRole, nil
}
//DeleteRole deletes a role with name roleName
func (app *App)DeleteRole(ctx context.Context, roleName string) error {
	role, err:=app.db.FindRole(ctx, roleName)
	if err!=nil{
		return err
	}
	if role==nil{
		return common.NotFoundError{ID:roleName}
	}
	users, err:= app.db.GetAllUsersWithRole(ctx, roleName)
	if err!=nil{
		return err
	}
	if len(users)>0{
		return fmt.Errorf("Users have the role %v cannot delete", roleName)
	}
	return app.db.RemoveRole(ctx, roleName)
}

//UpdateRole updates the privileges inside the role with name roleName
func (app *App)UpdateRole(ctx context.Context, roleName string, privileges ...common.SecurityPrivilege) (*common.SecurityRole,error) {
	role, err:=app.db.FindRole(ctx, roleName)
	if err!=nil{
		return nil,err
	}
	if role==nil{
		return nil,common.NotFoundError{ID:roleName}
	}
	role.Privileges = privileges
	role, err=app.db.UpdateRole(ctx, *role)
	if err!=nil{
		return nil, err
	}
	return role, err
}

//GetRole returns the role called roleName
func (app *App)GetRole(ctx context.Context, roleName string) (*common.SecurityRole, error){
	role, err:=app.db.FindRole(ctx, roleName)
	if role==nil{
		return nil, common.NotFoundError{ID: roleName}
	}
	return role, err
}
//GetAllRoles returns the list of all securityRoles
func (app *App)GetAllRoles(ctx context.Context) ([]common.SecurityRole, error){
	return app.db.GetAllRoles(ctx)
}
//AddRoleToUser adds the role with name roleName to the user
func (app *App)AddRoleToUser(ctx context.Context, username string ,roleName string) (*common.UserData,error){
	role, err:=app.db.FindRole(ctx, roleName)
	if err!=nil{
		return nil,err
	}
	if role==nil{
		return nil,common.NotFoundError{ID: roleName}
	}
	user, err:=app.db.FindUser(ctx, username)
	if err!=nil{
		return nil, err
	}
	if user==nil{
		return nil,common.NotFoundError{ID:username}
	}
	user.AddRole(ctx, *role)
	updpatedUser, err:= app.db.UpdateUser(ctx, *user);
	if err!=nil{
		return nil, err
	} 
	return updpatedUser, nil
}
//IsRoleAuthorized returns true if a role with roleName is enabled for a privilege
func (app *App)IsRoleAuthorized(ctx context.Context, roleName string, privilege common.SecurityPrivilege) (bool,error){
	role, err:=app.db.FindRole(ctx, roleName)
	if err!=nil{
		return false, err
	}
	if role==nil{
		return false, common.NotFoundError{ID: roleName}
	}
	if role.HasPrivilege(privilege) || role.HasPrivilege(common.SecurityAdmin) {
		return true, nil
	}
	return false, nil
}
//GetUserRole returns the role object for a user
func(app *App) GetUserRole(ctx context.Context, user common.UserData) (*common.SecurityRole, error){
	role, err:=app.db.FindRole(ctx, user.Role)
	if err!=nil{
		return nil, err
	}
	if role==nil{
		return nil, common.NotFoundError{ID: user.Role}
	}
	return role, nil
}

