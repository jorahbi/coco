package main

import (
	"flag"
	"fmt"
	query "tools/gorm-query/gen"
)

var dsn = flag.String("dsn", "", "the mysql dsn")
var genPath = flag.String("path", "", "the gen path")
var tables = flag.String("tables", "", "the tables ")

type Config struct {
	Db struct {
		Dsn string
	}
}

func main() {
	flag.Parse()
	if *dsn == "" || *genPath == "" {
		fmt.Println("miss params")
		return
	}
	query.GenQuery(*dsn, *genPath, *tables)
}

// func build(dsn string) {
// 	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
// 		Logger: logger.Default.LogMode(logger.Info),
// 	})

// 	q := query.Use(db)
// 	// q.Address
// 	address := q.WithContext(context.Background()).DcomAssetInventoryPlan
// 	condition := []query.DcomAssetInventoryPlanOption{q.DcomAssetInventoryPlan.WithIds([]int32{1, 2})}
// 	mAddress := model.DcomAssetInventoryPlan{}
// 	sql := address.WithOptions(address, condition...).UnderlyingDB().ToSQL(func(tx *gorm.DB) *gorm.DB {
// 		return address.WithOptions(address, condition...).UnderlyingDB().Find(&mAddress)
// 	})
// 	fmt.Println(sql)
// }