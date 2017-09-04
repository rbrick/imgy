package db

import "reflect"

type Invite struct {
	ID    uint   `gorm:"primary_key"`
	Code  string `sql:"not null;unique"`
	Valid bool   `sql:"not null"`
}

func (i *Invite) validate() bool {
	return !reflect.DeepEqual(i, &Invite{})
}

func (i *Invite) Save() {
	database.Save(i)
}

func GetInviteByCode(code string) *Invite {
	var invite Invite
	database.Model(&invite).Where("code = ?", code).Scan(&invite)
	if !invite.validate() {
		return nil
	}
	return &invite
}
