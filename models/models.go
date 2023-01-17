package model

import (
	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Id       int    `orm:"column(id);pk"`
	Name     string `orm:"column(name)"`
	Email    string `orm:"column(email)"`
	Password string `orm:"column(password)"`
}

func init() {
	orm.RegisterModel(new(User))
}

func (u *User) UserTable() string {
	return "users"
}
