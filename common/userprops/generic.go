package userprops

import (
	"encoding/json"
	"fmt"
)

//UserProperty is an interface that can contain any input field for a user, and its type.
type UserProperty interface{
	json.Marshaler
	GetName()string
	SetName(string)
	GetValueString()string
	SetValueString(string) error
	GetTypeString()UserPropertyType
	//IsMandatory should be used ONLY on frontend
	IsMandatory()bool
}
//AdmittedUserPropertyTypes is a list of the admitted types in string
var AdmittedUserPropertyTypes =[]UserPropertyType{UserTypeInt64, UserTypeString, UserTypeFloat64,UserTypeBool}
//UserPropertyType is the type mapped to string for UserProperty enabled
type UserPropertyType string
const (
	//UserTypeInt64 is the string rapresentation for the type of a userproperty of type in64
	UserTypeInt64 UserPropertyType="int64"
	//UserTypeString is the string rapresentation for the type of a userproperty of type string
	UserTypeString UserPropertyType="string"
	//UserTypeFloat64 is the string rapresentation for the type of a userproperty of type float64
	UserTypeFloat64 UserPropertyType="float64"
	//UserTypeBool is the string rapresentation for the type of a userproperty of type bool
	UserTypeBool UserPropertyType="bool"
)
//NewUserProperty returns a UserProperty of the requested type, or error
func NewUserProperty(userPropertyType UserPropertyType) (UserProperty,error){
	switch userPropertyType{
	case UserTypeBool:
		return new(BoolUserProperty), nil
	case UserTypeInt64:
		return new(IntUserProperty), nil
	case UserTypeFloat64:
		return new(FloatUserProperty), nil
	case UserTypeString:
		return new(StringUserProperty), nil
	}
	return nil,fmt.Errorf("Type not admitted")
}
