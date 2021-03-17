package userprops

import (
	"encoding/json"
)

//StringUserProperty is a userProperty definition
type StringUserProperty struct{
	Name string `json:"name"`
	Mandatory bool `json:"mandatory"`
	System bool `json:"system"`
}
//MarshalJSON needed for json-marshalling
func(property *StringUserProperty)	MarshalJSON()([]byte, error){
	return json.Marshal(*property)
}
//GetName returns the name of the property
func(property *StringUserProperty)	GetName()string{
	return property.Name
}
//SetName sets the name of the property
func(property *StringUserProperty)	SetName(value string){
	property.Name=value
}

//GetTypeString returns the type in a string formatted fashion
func(property *StringUserProperty)	GetTypeString()UserPropertyType{
	return UserTypeString
}
//SetMandatory sets the fact of a value being mandatory
func(property *StringUserProperty) SetMandatory(mandatory bool){
	property.Mandatory=mandatory
}
//IsMandatory should be used ONLY on frontend to decide if a field should be mandatory. no logic on backend here, because if added after would be hard
func(property *StringUserProperty)	IsMandatory()bool {
	return property.Mandatory
}

//SetIsSystem set value if a prop is a system property (that should not be deleted or updated)
func (property *StringUserProperty)	SetIsSystem(isSystem bool) {
	property.System=isSystem
}
//IsSystem will returns if a prop is a system property
func(property *StringUserProperty)	IsSystem()bool{
	return property.System
}
