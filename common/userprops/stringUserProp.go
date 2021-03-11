package userprops

import (
	"encoding/json"
)

//StringUserProperty is a userProperty containing a string value
type StringUserProperty struct{
	Name string `json:"name"`
	Value string `json:"value"`
	Mandatory bool `json:"mandatory"`
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
//GetValueString get the value formatted as string
func(property *StringUserProperty)	GetValueString()string{
	return property.Value
}
//SetValueString set the value by a string
func(property *StringUserProperty)	SetValueString(value string) error{
	property.Value=value
	return nil
}

//GetTypeString returns the type in a string formatted fashion
func(property *StringUserProperty)	GetTypeString()UserPropertyType{
	return UserTypeString
}

//IsMandatory should be used ONLY on frontend to decide if a field should be mandatory. no logic on backend here, because if added after would be hard
func(property *StringUserProperty)	IsMandatory()bool {
	return property.Mandatory
}
