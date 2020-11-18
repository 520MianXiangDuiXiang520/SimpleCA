package utils

import (
	"errors"
	"github.com/520MianXiangDuiXiang520/GinTools/utils"
	"reflect"
)

type DaoLogic interface{}

// 使用数据库事务执行方法 def, 返回执行结果, 如果 def 执行过程中发生错误导致
// panic 或 主动返回一个 error 时数据库会回滚保证一致性.
//   - def:  要执行的方法, 该方法的最后一个返回值应该是 error 类型
//   - args: 方法参数, 第一个参数需要是 *DB 类型, 在真实执行时,会被替换掉
func UseTransaction(def DaoLogic, args []interface{}) ([]reflect.Value, error) {
	var err error
	tx := GetDB().Begin()
	tx.LogMode(true)
	defer func() {
		// def 抛出 panic, 回滚
		if pan := recover(); pan != nil {
			utils.ExceptionLog(pan.(error), "Transaction execution failed and has been rolled back！")
			tx.Rollback()
		}
		// def 返回了一个 err, 回滚
		if err != nil {
			utils.ExceptionLog(err, "Transaction return false and has been rolled back！")
			tx.Rollback()
		}
		tx.Commit()
	}()
	value := reflect.ValueOf(def)
	if value.Kind() != reflect.Func {
		return nil, errors.New("TypeError: def is not a Func type")
	}
	if reflect.TypeOf(args[0]) != reflect.TypeOf(tx) {
		return nil, errors.New("TypeError: the first parameter must be of type *DB")
	}
	argsVal := make([]reflect.Value, len(args))
	for i, arg := range args {
		argsVal[i] = reflect.ValueOf(arg)
	}

	argsVal[0] = reflect.ValueOf(tx)
	res := value.Call(argsVal)

	// 如果 def 主动抛出异常，回滚
	errVal := res[len(res)-1]

	if errVal.Interface() != nil {
		if _, ok := errVal.Interface().(error); ok {
			err = errVal.Interface().(error)
		} else {
			err = errors.New("return error")
		}
		return res, err
	}
	return res, nil
}
