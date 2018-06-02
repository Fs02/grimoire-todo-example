# GO Paranoid
Paranoid is utilities supplemental to go panic, it usefull when you want to make your code panic on unhandled error, it also make your test case simpler by reducing the test branch.

## Usage
```golang
package main

import (
	"errors"
	"github.com/Fs02/go-paranoid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Transaction struct {
	ID uint
}

func retrieve() (Transaction, error) {
	db, err := gorm.Open("mysql", "root@(127.0.0.1:3306)/db?charset=utf8&parseTime=True&loc=Local")
	// panic if error
	paranoid.Panic(err, "Error opening database connection")

	trx := Transaction{}
	query := db.First(&trx, 1000)

	if query.RecordNotFound() {
		return trx, errors.New("not found")
	}

	// it'll panic on unknown or untestable error
	paranoid.Panic(query.Error, "Failed when fetching transaction %+v", trx)
	return trx, nil
}

func main() {
	retrieve()
}
```

There's also `paranoid.PanicFunc` which usefull when you want to run any specific function before panic (ie: rollback).
