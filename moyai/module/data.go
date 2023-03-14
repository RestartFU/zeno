package module

import (
	"github.com/jmoiron/sqlx"
	"github.com/moyai-studio/practice-revamp/moyai/user"
)

type Data struct {
	Remover interface{ RemoveUser(u *user.User) }
	DB      *sqlx.DB
	NopModule
}

func (d *Data) HandleQuit(u *user.User) {
	u.UpdateData(d.DB)
	d.Remover.RemoveUser(u)
}
