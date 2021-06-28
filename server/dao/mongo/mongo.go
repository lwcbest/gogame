package mongo

import (
	"fmt"
	"gameserver-997/server/base/logger"
	"gameserver-997/server/util"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	config   *util.MongoConfig
	dialInfo *mgo.DialInfo
)

func init() {

}

// crud接口
type Men interface {
	GetDb()
	Insert(collection string, docs ...interface{})
	Del(collection string, query interface{})
	Update(collection string, selector interface{}, change interface{})
	Find(collection string, query interface{}, obj interface{})
	FindById(collection string, objId string, obj interface{})
	FindOne(collection string, query interface{}, obj interface{})
	FindAll(collection string, obj interface{})
}

// myMogo类
type MyMongo struct {
	S  *mgo.Session
	DB *mgo.Database
}

// 初始化
func InitMongo() *MyMongo {
	config = util.MongoConf
	fmt.Printf("mongo config: %+v", config)
	dialInfo = &mgo.DialInfo{
		Addrs:     []string{config.Url},          // 数据库地址
		Timeout:   time.Duration(config.Timeout), // 连接超时时间
		Source:    config.Authdb,                 // 设置权限的数据库
		Username:  config.Authuser,               // 设置的用户名
		Password:  config.Authpass,               // 设置的密码
		PoolLimit: config.Poollimit,              // 连接池的数量
	}

	M := &MyMongo{}
	if config.IsAuth == true {
		defer fmt.Println("[auth init]")
		session, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			panic(err)
		}
		session.SetMode(mgo.Monotonic, true)
		M.S = session
		M.DB = session.DB(config.Db)
	} else {
		defer fmt.Println("[init] M", M)
		session, err := mgo.Dial(config.Url)
		if err != nil {
			panic(err)
		}
		session.SetMode(mgo.Monotonic, true)
		M.S = session
		M.DB = session.DB(config.Db)
	}
	return M
}

// 连接池
func (this *MyMongo) getSession(collection string) (*mgo.Session, *mgo.Collection) {
	if this.S == nil {
		var session *mgo.Session
		var err error
		if config.IsAuth == false {
			session, err = mgo.Dial(config.Url)
		} else {
			session, err = mgo.DialWithInfo(dialInfo)
		}
		if err != nil {
			panic(err)
		}
		session.SetMode(mgo.Monotonic, true)
		this.S = session
		this.DB = session.DB(config.Db)
	}
	ms := this.S.Clone()
	ms.SetMode(mgo.Monotonic, true)
	c := ms.DB(config.Db).C(collection)
	fmt.Println("[getSession c]", c)
	return ms, c
}

//获取db
func (this *MyMongo) GetDb() *mgo.Database {
	return this.DB
}

//插入
func (this *MyMongo) Insert(collection string, docs ...interface{}) error {
	s, c := this.getSession(collection)
	defer s.Close()
	err := c.Insert(docs[:]...)
	return err
}

//更新或者插入
func (this *MyMongo) Upsert(collection string, sel interface{}, update interface{}) error {
	s, c := this.getSession(collection)
	defer s.Close()
	info, err := c.Upsert(sel, update)
	logger.Info("Upsert %s result %+v", collection, info)
	return err
}

// 删除
func (this *MyMongo) Del(collection string, query interface{}) error {
	s, c := this.getSession(collection)
	defer s.Close()
	err := c.Remove(&query)
	return err
}

// update
func (this *MyMongo) Update(collection string, selector interface{}, change interface{}) error {
	s, c := this.getSession(collection)
	defer s.Close()
	err := c.Update(&selector, &change)
	return err
}

//查询所有
func (this *MyMongo) FindAll(collection string, obj interface{}) error {
	s, c := this.getSession(collection)
	defer s.Close()
	err := c.Find(nil).All(obj)
	return err
}

//其他条件查询
func (this *MyMongo) Find(collection string, query interface{}, obj interface{}) error {
	s, c := this.getSession(collection)
	defer s.Close()
	err := c.Find(query).All(obj)
	return err
}

//查询byId
func (this *MyMongo) FindById(collection string, objId string, obj interface{}) error {
	s, c := this.getSession(collection)
	defer s.Close()
	err := c.FindId(bson.ObjectIdHex(objId)).One(obj)
	return err
}

//查询1个
func (this *MyMongo) FindOne(collection string, query interface{}, obj interface{}) error {
	s, c := this.getSession(collection)
	defer s.Close()
	err := c.Find(query).One(obj)
	return err
}

//其他条件查询
func (this *MyMongo) FindAndSort(collection string, query interface{}, sort string, obj interface{}) error {
	s, c := this.getSession(collection)
	defer s.Close()
	err := c.Find(query).Sort(sort).All(obj)
	return err
}

//分页查询
func (this *MyMongo) FindInPage(collection string, query interface{}, sort []string, skip int, limit int, obj interface{}) error {
	s, c := this.getSession(collection)
	defer s.Close()
	logger.Info("find in page %+v, %+v, %d %d %+v\n", query, sort, skip, limit, obj)
	err := c.Find(query).Sort(sort...).Skip(skip).Limit(limit).All(obj)
	return err
}

//查询数量
func (this *MyMongo) Count(collection string, query interface{}) (int, error) {
	s, c := this.getSession(collection)
	defer s.Close()
	n, err := c.Find(query).Count()
	return n, err
}

//创建一个批量操作
func (this *MyMongo) CreateBulk(collection string) (*MgoBulk, error) {
	s, c := this.getSession(collection)
	bulk := c.Bulk()
	mBulk := &MgoBulk{
		session: s,
		bulk:    bulk,
	}
	return mBulk, nil
}

type MgoBulk struct {
	session *mgo.Session
	bulk    *mgo.Bulk
}

func (this *MgoBulk) Insert(obj ...interface{}) {
	this.bulk.Insert(obj...)
}

func (this *MgoBulk) Run() error {
	defer this.session.Close()
	ret, err := this.bulk.Run()
	logger.Info("mongoBulk ret: %+v err: %v", ret, err)
	return err
}
