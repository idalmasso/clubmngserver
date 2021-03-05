package models

import "github.com/idalmasso/clubmngserver/database"


var db database.ClubDb

func InitDB(database *database.ClubDb){
	db=*database
	db.Init()
}
