package database

type DbConn interface {
	Begin()
	RoolBack()
	Commit()
}
