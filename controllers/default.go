package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"beeBlog2/models"
	"fmt"
	"path"
	"time"
	"math"
)
var OneUser string
type MainController struct {
	beego.Controller
}
func (c *MainController) Get() {
	c.TplName = "register.html"
}

func (c *MainController) Post() {
	var user models.User
	userName := c.GetString("userName")
	passWord := c.GetString("password")
	user.Name =userName
	user.Pwd = passWord
	o :=orm.NewOrm()
	err := o.Read(&user,"Name","Pwd")
	if err==nil {
		c.Redirect("/login",302)
	}
	_,err =o.Insert(&user)
	if err != nil {
		beego.Info("注册成功")
	}
	c.Redirect("/login",302)
}

func (c *MainController) ShowLogin()  {
	userName := c.Ctx.GetCookie("userName")
	if userName != ""{
		c.Data["userName"] = userName
		c.Data["checked"] = "checked"
		}else {
			c.Data["userName"] = " "
		}
	c.TplName="login.html"
}
func (c *MainController) HandleLogin() {
	//拿到数据
	userName := c.GetString("userName")
	pwd := c.GetString("password")
	//判断
	if userName == ""|| pwd == "" {
		beego.Info("输入数据不合法")
		c.TplName="login.html"
		return
	}
	//查询账户和密码是否正确

	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	user.Pwd = pwd
	err := o.Read(&user,"Name","Pwd")
	if err != nil {
		beego.Info("查询失败")
		c.Redirect("/login",302)
		return
	}

	remember := c.GetString("remember")
	if remember=="on"{
		c.Ctx.SetCookie("userName",userName,200)
	}else {
		c.Ctx.SetCookie("userName",userName,-1)
	}
	c.SetSession("userName",userName)
	//c.Layout ="layout.html"
	//c.LayoutSections["contentHead"]="Hindex.html"
	OneUser=userName
	c.Redirect("/article/index",302)
}

func (c *MainController)ShowIndex() {

	o := orm.NewOrm()
	id,_ := c.GetInt("select")
	qs := o.QueryTable("Article")

	var articles []models.Article

	var count int64
	count,err := qs.RelatedSel("ArticleType").Count()
	if id != 0 {
		count,err = qs.RelatedSel("ArticleType").Filter("ArticleType__Id",id).Count()
	}
	if err != nil{
		beego.Info("查询错误")
		return
	}

	pageSize := 2
	pageCount := math.Ceil(float64(count)/float64(pageSize))
	pageIndex,err := c.GetInt("pageIndex")
	if err != nil {
		pageIndex=1
	}
	start:= pageSize*(pageIndex -1)
	if id == 0 {
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)
	}else{
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__Id",id).All(&articles)
	}

	FirstPage := false
	if pageIndex==1 {
		FirstPage=true
	}
	LastPage := false
	if pageIndex==int(pageCount) {
		LastPage=true
	}
	//count,err :=qs.Count()
	var articleType []models.ArticleType

	_,err = o.QueryTable("ArticleType").All(&articleType)
	if err != nil {
		beego.Info("获取类型错误")
		return
	}

	//qs.Limit(pageSize,start).Filter("articleType__Id",id).All(&articles)
	c.Data["oneuser"]=OneUser
	c.Data["articleType"]=articleType
	c.Data["FirstPage"]=FirstPage
	c.Data["LastPage"]=LastPage
	c.Data["pageIndex"]=pageIndex
	c.Data["pageCount"]=pageCount
	c.Data["count"]=count
	c.Data["articles"] = articles
	c.Data["typeid"] = id
	////var user := models.User{}

	c.Layout ="layout.html"
	//c.LayoutSections["contentHead"]="Hindex.html"
	c.TplName = "index.html"
	}
//文章添加
func (c *MainController) ShowAdd()  {
	o := orm.NewOrm()
	var artiTypes  []models.ArticleType
	_,err := o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil {
		beego.Info("获取类型错误")
		return
	}
	beego.Info(artiTypes)
	c.Data["oneuser"]=OneUser
	c.Data["articleType"]=artiTypes
	//c.LayoutSections["titleContent"]="Hadd.html"
	c.Layout ="layout.html"
	c.TplName="add.html"
}
func (c *MainController) HandleAdd()  {
	//c.TplName="add.html"
	//拿到数据
	//artiType :=c.GetString("select")
	artiName :=c.GetString("articleName")
	artiContent :=c.GetString("content")
	//文件上传
	f,h,err := c.GetFile("uploadname")
	defer f.Close()
	//限定格式png，jpg
	fileecxt := path.Ext(h.Filename)
	if fileecxt != ".jpg" && fileecxt != ".png"{
		beego.Info("上传文件格式错误")
		return
	}
	//限制大小
	if h.Size > 400000 {
		beego.Info("上传文件过大")
		return
	}
	//对文件重新命名，防止重复
	filename :=time.Now().Format("2016-01-02") +fileecxt
	if err != nil {
		fmt.Println("getFile err=",err)
	}else {
		c.SaveToFile("uploadname","./static/img"+filename)
	}
	if artiName == ""||artiContent =="" {
		beego.Info("添加文章数据错误")
		return
	}
	//c.Ctx.WriteString("添加文章成功")
	//c.Redirect("/index",302)
	////插入数据库
	//o := orm.NewOrm()
	//arti := models.Artcle{}
	var artiType models.ArticleType
	id,err := c.GetInt("select")
	if err != nil {
		beego.Info("err=",err)
	}
	o := orm.NewOrm()
	artiType = models.ArticleType{Id:id}
	o.Read(&artiType,"id")

	//_=o.Read(&artiType,"Tname")

	arti := models.Article{}
	arti.ArticleType=&artiType
	arti.Artiname = artiName
	arti.Acontent = artiContent
	arti.Aimg = h.Filename
	arti.Atime= time.Now()
	arti.ArticleType=&artiType
	_,err =o.Insert(&arti)
	if err!=nil {
		fmt.Println("插入数据库失败")
		return
	}
	//c.LayoutSections["titleContent"]="Hindex.html"
	c.Layout ="layout.html"
	c.Redirect("/article/index",302)
}

//显示详情
func (c *MainController)ShowContent()  {
	userName := c.GetSession("userName")
	if userName == nil {
		c.Redirect("/login", 302)
		return
	}

	id,err := c.GetInt("id")
	if err != nil {
		beego.Info("获取文章Id错误",err)
		return
	}
	o:=orm.NewOrm()
	arti := models.Article{Id:id}
	err = o.Read(&arti)
	if err != nil {
		beego.Info("查询失败")
		return
	}
	arti.Acount +=1

	beego.Info(userName)
	//
	m2m := o.QueryM2M(&arti,"User")
	user := models.User{Name:userName.(string)}
	beego.Info(user)
	err =o.Read(&user,"Name")
	if err != nil {
		beego.Info("err=",err)
	}
	m2m.Add(&user)

	o.Update(&arti)

	var users  []models.User
	o.QueryTable("User").Filter("Article__Article__id",id).Distinct().All(&users)
	beego.Info(users)
	c.Data["oneuser"]=OneUser
	c.Data["users"]=users
	c.Data["article"]=arti
	//c.LayoutSections["titleContent"]="Hcontent.html"
	c.Layout ="layout.html"
	c.TplName="content.html"

}

func (c *MainController) ShowUpdate()  {
	userName := c.GetSession("userName")
	if userName == nil {
		c.Redirect("/login", 302)
		return
	}

	id ,err := c.GetInt("id")
	if err != nil {
		beego.Info("获取文章Id错误",err)
		return
	}
	o:=orm.NewOrm()
	arti := models.Article{Id:id}
	err = o.Read(&arti)
	if err != nil {
		beego.Info("查询失败",err)
		return
	}
	//传递数据给视图
	c.Data["oneuser"]=OneUser
	c.Data["article"]=arti
	//c.LayoutSections["titleContent"]="Hupdate.html"
	c.Layout ="layout.html"
	c.TplName="update.html"
}

func (c *MainController)HandleUpdate()  {
	userName := c.GetSession("userName")
	if userName == nil {
		c.Redirect("/login", 302)
		return
	}
	id,_:=c.GetInt("id")
	artiname := c.GetString("articleName")
	content := c.GetString("content")


	f,h,err := c.GetFile("uploadname")
	if err !=nil {
		beego.Info("文章上传失败",err)
		return
	}else {
		defer f.Close()
	}
	//限定格式png，jpg
	fileecxt := path.Ext(h.Filename)
	if fileecxt != ".jpg" && fileecxt != ".png"{
		beego.Info("上传文件格式错误")
		return
	}
	//限制大小
	if h.Size > 400000 {
		beego.Info("上传文件过大")
		return
	}
	//对文件重新命名，防止重复
	filename :=time.Now().Format("2016-01-02") +fileecxt
	if err != nil {
		fmt.Println("getFile err=",err)
	}else {
		c.SaveToFile("uploadname","./static/img"+filename)
	}
	if artiname == ""||content =="" {
		beego.Info("添加文章数据错误")
		return
	}
	//更新数据
	o := orm.NewOrm()
	arti := models.Article{Id:id}
	err = o.Read(&arti)
	if err != nil {
		beego.Info("查询失败",err)
		return
	}
	arti.Artiname=artiname
	arti.Acontent=content
	arti.Aimg="./static/img/"+filename

	_,err = o.Update(&arti,"ArtiName","Acontent","Aimg")
	if err != nil {
		beego.Info("更新失败")
		return
	}
	//c.LayoutSections["titleContent"]="Hindex.html"
	c.Layout ="layout.html"
	c.Redirect("/article/index",302)

}

func (c *MainController)HandleDelete()  {
	userName := c.GetSession("userName")
	if userName == nil {
		c.Redirect("/login", 302)
		return
	}
	id,err := c.GetInt("id")
	if err != nil {
		beego.Info("获取Id错误")
		return
	}

	o := orm.NewOrm()
	arti:= models.Article{Id:id}

	err = o.Read(&arti)
	if err!=nil {
		beego.Info("")
		return
	}
	o.Delete(&arti)
	//c.LayoutSections["titleContent"]="Hindex.html"
	c.Layout ="layout.html"
	c.Redirect("/article/index",302)
}
func (c *MainController)ShowAddType()  {
	//找到元素

	o := orm.NewOrm()
	var artiType []models.ArticleType
	_,err := o.QueryTable("ArticleType").All(&artiType)
	if err != nil {
		beego.Info("未找到artiType",err)
		//return
	}
	c.Data["oneuser"]=OneUser
	c.Data["articleType"] = artiType
	//c.LayoutSections["titleContent"]="HaddType.html"
	c.Layout ="layout.html"
	c.TplName="addType.html"
}
func (c *MainController)HandleAddType()  {
	userName := c.GetSession("userName")
	if userName == nil {
		c.Redirect("/login", 302)
		return
	}
	typeName := c.GetString("typeName")
	//判断数据是否合法
	if typeName == ""{
		beego.Info("获取类型错误")
		return
	}
	o := orm.NewOrm()
	artiType := models.ArticleType{}
	artiType.Tname=typeName
	_,err := o.Insert(&artiType)
	if err != nil {
		beego.Info("插入失败")
		return
	}
	//c.LayoutSections["titleContent"]="Hindex.html"
	c.Layout ="layout.html"
	c.Redirect("/article/addType",302)
}

func (c *MainController)LogOut()  {
	c.DelSession("userName")
	c.Redirect("/login",302)
}