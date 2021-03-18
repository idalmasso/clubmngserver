package database

type ClubDbTest interface{
	ClubDb
	TearDown()
}
