package util

import (
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/ungerik/go-dry"
)

func StructFields(s interface{}) []reflect.StructField {
	var fields []reflect.StructField

	e := reflect.ValueOf(s).Elem()
	mainStruct := e.Type()

	for i := 0; i < e.NumField(); i++ {
		mf := e.Field(i)
		deepStruct := mf.Type()

		// lvl1 embedded structs
		if mf.Kind() == reflect.Struct {
			for j := 0; j < mf.NumField(); j++ {
				f := deepStruct.Field(j)
				fields = append(fields, f)
			}
		}

		f := mainStruct.Field(i)
		fields = append(fields, f)
	}

	return fields
}

func ChangeStructFieldValue(cStruct interface{}, cField string, cValue string) {
	//e := reflect.ValueOf(&cStruct).Elem()
	//mainStruct := e.Type()
	//log.Println("wat", e.Type().Field(0))

	//s := structs.New(cStruct)
	//err := s.Field(cField).Set(cValue)

	/*for i := 0; i < e.NumField(); i++ {
			mf := e.Field(i)
			deepStruct := mf.Type()

			// lvl1 embedded structs
			if mf.Kind() == reflect.Struct {
				f, ok := deepStruct.FieldByName(j)
				if f.IsValid() && f.CanSet() {
					f.Set(cValue)
				}
				log.Println(ok)
			}

			f, ok := mainStruct.FieldByName(cField)
	    mainStruct.Elem()
			if f.IsValid() && f.CanSet() {
				f.Set(cValue)
			}
			log.Println(ok)
		}*/
}

func PrepareForSQL(s interface{}, exclude []string) []string {
	var columns []string
	var find func(e reflect.Value)

	e := reflect.ValueOf(s).Elem()
	find = func(e reflect.Value) {
		mainStruct := e.Type()

		for i := 0; i < e.NumField(); i++ {
			mf := e.Field(i)

			// embedded structs
			// excludes time.Time
			if mf.Kind() == reflect.Struct && mf.Kind() != reflect.TypeOf(time.Time{}).Kind() {
				find(mf)
			}

			f := mainStruct.Field(i)

			name := strings.ToLower(f.Name)
			tag := f.Tag.Get("sql")

			// ignores columns with tag: sql:"-"
			if tag == "-" || dry.StringListContains(exclude, name) {
				continue
			} else if tag != "" {
				name = tag
			}

			log.Println(name, tag)
			columns = append(columns, name)
		}
	}

	find(e)

	/*for i := 0; i < e.NumField(); i++ {
		mf := e.Field(i)
		deepStruct := mf.Type()

		// lvl1 embedded structs
		if mf.Kind() == reflect.Struct {
			for j := 0; j < mf.NumField(); j++ {
				f := deepStruct.Field(j)

				name := strings.ToLower(f.Name)
				tag := f.Tag.Get("sql")
				log.Println(name, tag)

				// ignores columns with tag: sql:"-"
				if tag == "-" {
					continue
				} else if tag != "" {
					name = tag
				}

				columns = append(columns, name)
			}
		}

		f := mainStruct.Field(i)

		name := strings.ToLower(f.Name)
		tag := f.Tag.Get("sql")
		log.Println(name, tag)

		// ignores columns with tag: sql:"-"
		if tag == "-" {
			continue
		} else if tag != "" {
			name = tag
		}

		columns = append(columns, name)
	}*/

	return columns
}
