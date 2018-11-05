package models

import (
	_ "github.com/go-sql-driver/mysql"
	"time"
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id int `orm:"pk;auto"`
	Name string
	Pwd string

	Article []*Article `orm:"rel(m2m)"`
}

type Article struct {
	Id int  `orm:"pk;auto"`  //主键，自动增长
	Artiname string `orm:"size(20)"`   //ArtiName长度为20
	Atime time.Time   `orm:"auto_now"`//
	Acount  int			`orm:"default(0);null"`
	Acontent string
	Aimg string
	ArticleType *ArticleType `orm:"rel(fk)"`

	User []*User `orm:"reverse(many)"`
}

type ArticleType struct {
	Id int `orm:"pk;auto"`
	Tname string
	Articles []*Article   `orm:"reverse(many)"`
}



func init()  {
	//orm.RegisterDriver("mysql", orm.MySQL)
	orm.RegisterDataBase("default","mysql","root:root@/test1?charset=utf8")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default",false,true)
}