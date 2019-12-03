package kdb

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func toString(src interface{}) (dst string,err error){
	inf :=reflect.Indirect(reflect.ValueOf(src)).Interface()

	if inf ==nil{
		return "",nil
	}

	switch v:=inf.(type) {
	case string:
		dst =v
		return
	case []byte:
		dst=string(v)
		return
	}

	val :=reflect.ValueOf(inf)
	typ := reflect.TypeOf(inf)

	switch typ.Kind() {
	case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
		dst = strconv.FormatInt(val.Int(),10)
	case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
		dst = strconv.FormatUint(val.Uint(),10)
	case reflect.Float32,reflect.Float64:
		dst = strconv.FormatFloat(val.Float(),'f',-1,64)
	case reflect.Bool:
		dst = strconv.FormatBool(val.Bool())
	case reflect.Complex64,reflect.Complex128:
		dst = fmt.Sprintf("%v",val.Complex())
	case reflect.Struct:
		var timeType time.Time
		if typ.ConvertibleTo(reflect.TypeOf(timeType)){
			dst = val.Convert(reflect.TypeOf(timeType)).Interface().(time.Time).Format(time.RFC3339Nano)
		}else{
			err = fmt.Errorf("unsuported struct type %v",val.Type())
		}
	default:
		err=fmt.Errorf("unsupported struct type %v",val.Type())
	}
	return
}

func extractTagInfo(st reflect.Value)(tagList map[string]reflect.Value,err error) {
	stVal := reflect.Indirect(st)
	if stVal.Kind() != reflect.Struct {
		return nil, fmt.Errorf("this variable type is %v,not a struct", st.Kind())
	}

	tagList = make(map[string]reflect.Value)

	for i := 0; i < stVal.NumField(); i++ {
		v := stVal.Field(i)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				var typ reflect.Type
				if v.Type().Kind() == reflect.Ptr {
					typ = v.Type().Elem()
				} else {
					typ = v.Type()
				}
				vv := reflect.New(typ)
				v.Set(vv)
			}

			if v.Elem().Kind() == reflect.Struct {
				t, err := extractTagInfo(v.Elem())
				if err != nil {
					return nil, err
				}

				for k, ptr := range t {
					if _, ok := tagList[k]; ok {
						return nil, fmt.Errorf("%s:%s is exists", kdb.structTag, k)
					}
					tagList[k] = ptr
				}
			} else if v.Kind() == reflect.Map && v.IsNil() {
				v.Set(reflect.MakeMap(v.Type()))
			} else if v.Kind() == reflect.Struct {
				var ignore bool
				switch v.Interface().(type) {
				case time.Time:
					ignore = true
				case sql.NullTime:
					ignore = true
				case sql.NullString:
					ignore = true
				case sql.NullBool:
					ignore = true
				case sql.NullInt32:
					ignore = true
				case sql.NullInt64:
					ignore = true
				}
				if !ignore {
					t, err := extractTagInfo(v)
					if err != nil {
						return nil, err
					}
					for k, ptr := range t {
						if _, ok := tagList[k]; ok {
							return nil, fmt.Errorf("%s:%s is exists", kdb.structTag, k)
						}
						tagList[k] = ptr
					}
				}
			}
			tagName := stVal.Type().Field(i).Tag.Get(kdb.structTag)
			if tagName != "" {
				attr := strings.Split(tagName, ";")
				column := attr[0]
				if _, ok := tagList[column]; ok {
					return nil, fmt.Errorf("%s:%s is exists", kdb.structTag, tagName)
				}
				tagList[column] = v
			}
		}
	}
	return
}
