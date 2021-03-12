package userprops

import "encoding/json"

//IntUserProperty is a userProperty definition
type IntUserProperty struct{
	Name string `json:"name"`
	Mandatory bool `json:"mandatory"`
}
//MarshalJSON needed for json-marshalling
func(property *IntUserProperty)	MarshalJSON()([]byte, error){
	return json.Marshal(*property)
}
//GetName returns the name of the property
func(property *IntUserProperty)	GetName()string{
	return property.Name
}
//SetName sets the name of the property
func(property *IntUserProperty)	SetName(value string){
	property.Name=value
}
//GetTypeString returns the type in a string formatted fashion
func(property *IntUserProperty)	GetTypeString()UserPropertyType{
	return UserTypeInt64
}
//SetMandatory sets the fact of a value being mandatory
func(property *IntUserProperty) SetMandatory(mandatory bool){
	property.Mandatory=mandatory
}
//IsMandatory should be used ONLY on frontend to decide if a field should be mandatory. no logic on backend here, because if added after would be hard
func(property *IntUserProperty)	IsMandatory()bool {
	return property.Mandatory
}

