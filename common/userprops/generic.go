package userprops

import (
	"encoding/json"
	"fmt"
)

//UserPropertyDefinition is an interface that can contain any input field for a user, and its type.
type UserPropertyDefinition interface{
	
	GetName()string
	SetName(string)
	GetTypeString()UserPropertyType
	SetMandatory(bool)
	//IsMandatory should be used ONLY on frontend
	IsMandatory()bool
}
//UserPropertyValue is the interface of the value of a userproperty
type UserPropertyValue interface {
	json.Marshaler
	UserPropertyDefinition
	GetValueString()string
	SetValueString(string) error
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
//NewUserPropertyDefinition returns a UserProperty of the requested type, or error
func NewUserPropertyDefinition(userPropertyType UserPropertyType) (UserPropertyDefinition,error){
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

//NewUserPropertyValue returns a UserPropertyvalue of the requested type, or error
func NewUserPropertyValue(userPropertyType UserPropertyType) (UserPropertyValue,error){
	switch userPropertyType{
	case UserTypeBool:
		return new(BoolUserPropertyValue), nil
	case UserTypeInt64:
		return new(IntUserPropertyValue), nil
	case UserTypeFloat64:
		return new(FloatUserPropertyValue), nil
	case UserTypeString:
		return new(StringUserPropertyValue), nil
	}
	return nil,fmt.Errorf("Type not admitted")
}
