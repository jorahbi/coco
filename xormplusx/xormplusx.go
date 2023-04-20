package xormplusx

import (
	"fmt"
	"log"

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

type EngineGroup struct {
	group *xorm.EngineGroup
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

type SqlArgs map[string]any

var single = &singleflight.Group{}

func OrmGroup(o Options) *EngineGroup {
	group, err, _ := single.Do("engine_group", func() (interface{}, error) {
		return newOrm(o)
	})
	if err != nil {
		log.Panicln(err.Error())
	}

	return group.(*EngineGroup)
}

func newOrm(o Options) (*EngineGroup, error) {
	var (
		slaves        []*xorm.Engine
		master, slave *xorm.Engine
		group         *xorm.EngineGroup
		err           error
	)
	orm := &EngineGroup{}
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
	group.SetMaxOpenConns(100)
	// 连接池中最大空闲连接数
	group.SetMaxIdleConns(10)
	// 单个连接最大存活时间(单位:秒)
	group.SetConnMaxLifetime(10000)
	orm.group = group

	return orm, nil

}

func (orm *EngineGroup) NewSqlArgs(queryId string) SqlArgs {
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

func (orm *EngineGroup) ReadConn() *xorm.Engine {
	return orm.group.Main()
}

func (orm *EngineGroup) WriteConn() *xorm.Engine {
	return orm.group.Subordinate()
}
