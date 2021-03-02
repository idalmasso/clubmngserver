package database

//UserInterface is an interface that returns the data of a user
type UserInterface interface{
	GetUsername() string
	GetEmail() string
	
}
