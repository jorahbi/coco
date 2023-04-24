package xormplusx

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
	"golang.org/x/sync/singleflight"
)

const (
	UNDELETED = 0
	DELETED   = 1
	INSERT    = "insert"
	UPDATE    = "update"
)

type engineGroup struct {
	group   *xorm.EngineGroup
	options Options
}
type Options struct {
	Sources struct {
		Master string
		Slave  []string
	}
	ShowSQL     bool
	StplPath    string
	MaxOpen     int
	MaxIdle     int
	MaxLifetime int
}

type Single func() *engineGroup
type SqlArgs map[string]any

var sf = &singleflight.Group{}
var eg *engineGroup

func Orm(o Options) Single {
	if eg != nil {
		return single
	}
	group, err, _ := sf.Do("engine_group", func() (interface{}, error) {
		return newOrm(o)
	})
	if err != nil {
		log.Panicln(err.Error())
	}
	eg = group.(*engineGroup)

	return single
}

func single() *engineGroup {
	return eg
}

func newOrm(o Options) (*engineGroup, error) {
	var (
		slaves        []*xorm.Engine
		master, slave *xorm.Engine
		group         *xorm.EngineGroup
		err           error
	)
	orm := &engineGroup{options: o}
	if master, err = xorm.NewEngine("mysql", o.Sources.Master); err != nil {
		return orm, err
	}
	if err != nil {
		log.Panicln(err.Error())
	}
	for _, source := range o.Sources.Slave {
		if slave, err = xorm.NewEngine("mysql", source); err != nil {
			return orm, err
		}
		slaves = append(slaves, slave)
	}

	if group, err = xorm.NewEngineGroup(master, slaves); err != nil {
		return orm, err
	}
	if err = group.Ping(); err != nil {
		return orm, err
	}
	if err = group.RegisterSqlTemplate(xorm.Pongo2(o.StplPath, ".stpl")); err != nil {
		return orm, err
	}
	group.StartFSWatcher()
	group.ShowSQL(o.ShowSQL)

	// 连接池中最大连接数
	group.SetMaxOpenConns(o.MaxOpen)
	// 连接池中最大空闲连接数
	group.SetMaxIdleConns(o.MaxIdle)
	// 单个连接最大存活时间(单位:秒)
	group.SetConnMaxLifetime(time.Duration(o.MaxLifetime))
	orm.group = group

	return orm, nil

}

func (orm *engineGroup) NewSqlArgs(queryId string) SqlArgs {
	args := SqlArgs{}
	args.Set("queryId", queryId)
	args.Set("isDelete", UNDELETED)
	return args
}

func (args SqlArgs) Set(key string, value any) SqlArgs {
	args[key] = value
	return args
}

func (args SqlArgs) Done() map[string]any {
	return map[string]any(args)
}

func (args SqlArgs) String() string {
	return fmt.Sprintln("sql args", map[string]any(args))
}

func (orm *engineGroup) ReadConn() *xorm.Engine {
	return orm.group.Main()
}

func (orm *engineGroup) WriteConn() *xorm.Engine {
	return orm.group.Subordinate()
}
