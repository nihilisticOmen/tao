package dao

import (
	"project-user/internal/database"
	"project-user/internal/database/gorms"
)

type TransactionImpl struct {
	conn database.DbConn
}

func NewTransaction() *TransactionImpl {
	return &TransactionImpl{
		conn: gorms.NewTran(),
	}
}
func (t TransactionImpl) Action(f func(conn database.DbConn) error) error {
	err := f(t.conn)
	if err != nil {
		t.conn.RoolBack()
		return err
	}
	t.conn.Commit()
	return nil
}
