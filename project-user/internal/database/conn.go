package database

type DbConn interface {
	RoolBack()
	Commit()
}
