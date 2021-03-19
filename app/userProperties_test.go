package app

import (
	"context"
	"errors"
	"testing"

	"github.com/idalmasso/clubmngserver/common"
	"github.com/idalmasso/clubmngserver/common/userprops"
)
func checkPropertyValues(prop userprops.UserPropertyDefinition, name string, propType userprops.UserPropertyType, mandatory bool , isSystem bool ) bool{
	return prop.IsMandatory()==mandatory && prop.GetName()==name && prop.GetTypeString()==propType && prop.IsSystem()==isSystem;
}
type addUserPropertyTest struct{
	name string
	userPropertyName string
	userPropertyType userprops.UserPropertyType
	mandatory bool
	isSystem bool
	expectError bool
	errorType error
}
var addUserPropertyTestSet = []addUserPropertyTest{
	{
		name: "Test user property with int value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with int value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with int value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with int value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},

	{
		name: "Test user property with int value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
	{
		name: "Test user property with int value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
	{
		name: "Test user property with int value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
	{
		name: "Test user property with int value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
	//Float
		{
		name: "Test user property with float value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with float value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with float value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with float value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},

	{
		name: "Test user property with float value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
	{
		name: "Test user property with float value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
	{
		name: "Test user property with float value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
	{
		name: "Test user property with float value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
	//bool
		{
		name: "Test user property with bool value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with bool value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with bool value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeBool,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with bool value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},

	{
		name: "Test user property with bool value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
	{
		name: "Test user property with bool value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
	{
		name: "Test user property with bool value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeBool,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
	{
		name: "Test user property with bool value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
	//string
		{
		name: "Test user property with string value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with string value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with string value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeString,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with string value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},

	{
		name: "Test user property with string value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
	{
		name: "Test user property with string value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
	{
		name: "Test user property with string value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeString,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
	{
		name: "Test user property with string value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.AlreadyExistsError{},
	},
}
func TestAddUserProperty(t *testing.T){
	var testApp App
	testApp.InitDB(&testDatabase)
	_, err:=testApp.AddUserProperty(context.Background(), "test01", false,false,userprops.UserPropertyType("int64"))
	if err!=nil{
		t.Fatalf("Insertion of first property failed, skipping other tests")
	}
	t.Cleanup(func(){
		testApp.db.RemoveUserPropertyDefinition(context.Background(), "test01")
	})
	for _, test:=range(addUserPropertyTestSet){
		t.Run(test.name, func(t *testing.T){
			prop, err:=testApp.AddUserProperty(context.Background(), test.userPropertyName, test.mandatory, test.isSystem, test.userPropertyType)
			
			if test.expectError{
				if err==nil{
					t.Fatal("Should have an error, not got")
				}
				if !errors.As(err, &test.errorType) {
					t.Fatalf("Should have an error of type %T, got %T", test.errorType, err)
				}
				if prop!=nil{
					t.Fatal("Should not get a property result with the error")
				}
				return
			}
			t.Cleanup(func(){
				testApp.db.RemoveUserPropertyDefinition(context.Background(),test.userPropertyName)
			})
			if err!=nil{
				t.Fatalf("Should not have an error, got %s", err.Error())
			}
			if prop==nil{
				t.Fatal("Property should not be nil")
			}
			if !checkPropertyValues(prop, test.userPropertyName, test.userPropertyType, test.mandatory, test.isSystem){
				t.Fatal("Property values not correct returned from addProperty")
			}
			prop, err=testApp.GetUserProperty(context.Background(), test.userPropertyName)
			if err!=nil{
				t.Fatalf("Cannot find added property: %s", err.Error())
			}
			if prop==nil{
				t.Fatal("Property found should not be nil")
			}
			if !checkPropertyValues(prop, test.userPropertyName, test.userPropertyType, test.mandatory, test.isSystem){
				t.Fatal("Property values for found prop not correct returned from addProperty")
			}
		})
	}
}

type removeUserPropertyTest struct{
	name string
	userPropertyToDeleteName string
	userPropertyName string
	userPropertyType userprops.UserPropertyType
	mandatory bool
	isSystem bool
	expectError bool
	errorType error
}
var removeUserPropertyTestSet = []removeUserPropertyTest{
	{
		name: "Test user property with int value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with int value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.BadRequestParametersError{},
	},
	{
		name: "Test user property with int value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with int value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.BadRequestParametersError{},
	},

	//Float
		{
		name: "Test user property with float value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with float value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.BadRequestParametersError{},
	},
	{
		name: "Test user property with float value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with float value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.BadRequestParametersError{},
	},

	//bool
		{
		name: "Test user property with bool value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with bool value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.BadRequestParametersError{},
	},
	{
		name: "Test user property with bool value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeBool,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with bool value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.BadRequestParametersError{},
	},

	//string
		{
		name: "Test user property with string value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with string value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.BadRequestParametersError{},
	},
	{
		name: "Test user property with string value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeString,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with string value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.BadRequestParametersError{},
	},
	//Now not existing ones
	{
		name: "Test user property with int value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with int value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with int value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with int value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},

	//Float
		{
		name: "Test user property with float value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with float value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with float value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with float value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},

	//bool
		{
		name: "Test user property with bool value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with bool value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with bool value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeBool,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with bool value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},

	//string
		{
		name: "Test user property with string value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with string value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with string value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeString,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with string value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToDeleteName: "TestProperty1",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
}
func TestRemoveUserProperty(t *testing.T){
	var testApp App
	testApp.InitDB(&testDatabase)
	for _, test:=range(removeUserPropertyTestSet){
		t.Run(test.name, func(t *testing.T){
			_, err:=testApp.AddUserProperty(context.Background(), test.userPropertyName, test.mandatory, test.isSystem, test.userPropertyType)
			t.Cleanup(func(){
				testApp.db.RemoveUserPropertyDefinition(context.Background(),test.userPropertyName)
			})
			if err!=nil{
				t.Fatalf("Should not have an error on add, got %s", err.Error())
			}
			err=testApp.DeleteUserProperty(context.Background(), test.userPropertyToDeleteName)
			if test.expectError{
				if err==nil{
					t.Fatalf("Should get an error, does not")
				}
				if !errors.As(err, &test.errorType){
					t.Fatalf("Got error %T, should have %T", err, test.errorType)
				}
				return
			}
			if err!=nil{
				t.Fatalf("Got an error, should not: %s", err.Error())
			}
		})
	}
}

func TestFindAllUserProperties(t *testing.T){
	var testApp App
	testApp.InitDB(&testDatabase)

	properties,err:=testApp.GetUserPropertiesList(context.Background())
	if err!=nil{
		t.Fatalf("Error on get: %s", err.Error())
	}
	lengthAlready:=len(buildUserPropertiesToBeAdded())
	if len(properties)!=lengthAlready{
		t.Fatalf("Length of properties gotten from call: %d, old one %d", len(properties), lengthAlready)
	}
	testApp.AddUserProperty(context.Background(), "test01", false, false, userprops.UserTypeBool)
	t.Cleanup(func(){
		testApp.db.RemoveUserPropertyDefinition(context.Background(), "test01")
	})
	properties,err=testApp.GetUserPropertiesList(context.Background())
	if err!=nil{
		t.Fatalf("Error on get: %s", err.Error())
	}
	if len(properties)!=lengthAlready+1{
		t.Fatalf("Length of properties gotten from call: %d, old one %d", len(properties), lengthAlready+1)
	}
}
type findUserPropertyTest struct{
	name string
	userPropertyToFind string
	userPropertyName string
	userPropertyType userprops.UserPropertyType
	mandatory bool
	isSystem bool
	expectError bool
	errorType error
}
var findUserProperyTestSet = []findUserPropertyTest{
	{
		name: "Test user property with int value, not existing, not mandatory, not system",
		userPropertyToFind:"test01",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with int value, not existing, not mandatory, system",
		userPropertyToFind:"test01",
		userPropertyName: "TestProperty",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with int value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with int value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},

	{
		name: "Test user property with int value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with int value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with int value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with int value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	//Float
		{
		name: "Test user property with float value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with float value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with float value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with float value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},

	{
		name: "Test user property with float value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with float value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with float value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with float value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	//bool
		{
		name: "Test user property with bool value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with bool value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with bool value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeBool,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with bool value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},

	{
		name: "Test user property with bool value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with bool value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with bool value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeBool,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with bool value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	//string
		{
		name: "Test user property with string value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with string value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with string value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeString,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with string value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},

	{
		name: "Test user property with string value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with string value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with string value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeString,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with string value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyToFind:"test01",
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
}

func TestFindSingleUserProperty(t *testing.T){
	var testApp App
	testApp.InitDB(&testDatabase)
	for _,test:=range(findUserProperyTestSet){
		t.Run(test.name, func(t *testing.T){
			_, err:=testApp.AddUserProperty(context.Background(), test.userPropertyName, test.mandatory, test.isSystem, test.userPropertyType )
			t.Cleanup(func(){
				testApp.db.RemoveUserPropertyDefinition(context.Background(), test.userPropertyName)
			})
			if err!=nil{
				t.Fatalf("Cannot insert the user property")
			}
			prop, err:=testApp.GetUserProperty(context.Background(), test.userPropertyToFind)
			if test.expectError{
				if err==nil{
					t.Fatalf("Expected error but got nil")
				}
				if !errors.As(err, &test.errorType){
					t.Fatalf("Wrong error type, got %T, expexted %T", err, test.errorType)
				}
				if prop!=nil{
					t.Fatalf("Returned property should be null")
				}
				return
			}
			if err!=nil{
				t.Fatalf("Expected nnot error but got one")
			}
			if prop==nil{
				t.Fatalf("No prop returned")
			}
			if !checkPropertyValues(prop, test.userPropertyName, test.userPropertyType, test.mandatory, test.isSystem){
				t.Fatalf("Prop values not ok")
			}
		})
	}
}


type updateUserPropertyTest struct{
	name string
	userPropertyToUpdate string
	mandatoryToUpdate bool
	userPropertyName string
	userPropertyType userprops.UserPropertyType
	mandatory bool
	isSystem bool
	expectError bool
	errorType error
}
var updateUserProperyTestSet = []updateUserPropertyTest{
	{
		name: "Test user property with int value, not existing, not mandatory, not system",
		userPropertyToUpdate:"test01",
		userPropertyName: "TestProperty",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with int value, not existing, not mandatory, system",
		userPropertyToUpdate:"test01",
		userPropertyName: "TestProperty",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with int value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with int value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},

	{
		name: "Test user property with int value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with int value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with int value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with int value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	//Float
		{
		name: "Test user property with float value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with float value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with float value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with float value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},

	{
		name: "Test user property with float value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with float value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with float value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with float value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	//bool
		{
		name: "Test user property with bool value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with bool value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with bool value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeBool,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with bool value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},

	{
		name: "Test user property with bool value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with bool value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with bool value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeBool,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with bool value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	//string
		{
		name: "Test user property with string value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with string value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with string value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeString,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with string value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},

	{
		name: "Test user property with string value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with string value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with string value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeString,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with string value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: true,
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	//HERE Additional tests for mandatoryToUpdate: false

		{
		name: "Test user property with int value, not existing, not mandatory, not system",
		userPropertyToUpdate:"test01",
		userPropertyName: "TestProperty",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with int value, not existing, not mandatory, system",
		userPropertyToUpdate:"test01",
		userPropertyName: "TestProperty",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with int value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with int value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},

	{
		name: "Test user property with int value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with int value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with int value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with int value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeInt64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	//Float
		{
		name: "Test user property with float value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with float value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with float value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with float value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},

	{
		name: "Test user property with float value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with float value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with float value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with float value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeFloat64,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	//bool
		{
		name: "Test user property with bool value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with bool value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with bool value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeBool,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with bool value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},

	{
		name: "Test user property with bool value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with bool value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with bool value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeBool,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with bool value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeBool,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	//string
		{
		name: "Test user property with string value, not existing, not mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with string value, not existing, not mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with string value, not existing, mandatory, not system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeString,
		mandatory: true,
		isSystem: false,
		expectError: true,
		errorType: common.NotFoundError{},
	},
	{
		name: "Test user property with string value, not existing, mandatory, system",
		userPropertyName: "TestProperty",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: true,
		errorType: common.NotFoundError{},
	},

	{
		name: "Test user property with string value, existing, not mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with string value, existing, not mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with string value, existing, mandatory, not system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeString,
		mandatory: true,
		isSystem: false,
		expectError: false,
		errorType: nil,
	},
	{
		name: "Test user property with string value, existing, mandatory, system",
		userPropertyName: "test01",
		userPropertyToUpdate:"test01",
		mandatoryToUpdate: false,
		userPropertyType: userprops.UserTypeString,
		mandatory: false,
		isSystem: true,
		expectError: false,
		errorType: nil,
	},
}

func TestUpdateUserProperty(t *testing.T){
	var testApp App
	testApp.InitDB(&testDatabase)
	for _,test:=range(updateUserProperyTestSet){
		t.Run(test.name, func(t *testing.T){
			_, err:=testApp.AddUserProperty(context.Background(), test.userPropertyName, test.mandatory, test.isSystem, test.userPropertyType )
			t.Cleanup(func(){
				testApp.db.RemoveUserPropertyDefinition(context.Background(), test.userPropertyName)
			})
			if err!=nil{
				t.Fatalf("Cannot insert the user property")
			}
			prop, err:=testApp.UpdateUserProperty(context.Background(), test.userPropertyToUpdate, test.mandatoryToUpdate)
			if test.expectError{
				if err==nil{
					t.Fatalf("Expected error but got nil")
				}
				if !errors.As(err, &test.errorType){
					t.Fatalf("Wrong error type, got %T, expexted %T", err, test.errorType)
				}
				if prop!=nil{
					t.Fatalf("Returned property should be null")
				}
				return
			}
			if err!=nil{
				t.Fatalf("Expected nnot error but got one")
			}
			if prop==nil{
				t.Fatalf("No prop returned")
			}
			if !checkPropertyValues(prop, test.userPropertyName, test.userPropertyType, test.mandatoryToUpdate, test.isSystem){
				t.Fatalf("Prop values not ok")
			}
			prop, err=testApp.GetUserProperty(context.Background(), test.userPropertyName)
			if err!=nil{
				t.Fatalf("Cannot get property for check")
			}
			if !checkPropertyValues(prop, test.userPropertyName, test.userPropertyType, test.mandatoryToUpdate, test.isSystem){
				t.Fatalf("Prop values are not ok after GET")
			}
		})
	}
}
