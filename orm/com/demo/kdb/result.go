package kdb

import (
	"database/sql"
	"fmt"
	"reflect"
)


type Row struct {
	rs *Rows
	lastError error
}

func (r *Row) ToArray() (result []string,err error) {
	items,err :=r.rs.ToArray()
	if err != nil{
		r.lastError = err
		return nil,err
	}
	return items[0],nil
}

func (r *Row) ToMap()(result map[string]string,err error){

	items,err := r.rs.ToMap()
	if err !=nil{
		r.lastError = err
		return nil,err
	}
	return items[0],nil
}

func (r *Row) ToStruct(st interface{})error{
	 stType:=reflect.TypeOf(st)
	 stVal :=reflect.ValueOf(st)

	 if stType.Kind() !=reflect.Ptr{
	 	return fmt.Errorf("this variavle type is %v,not a pointer",stType.Kind())
	 }
	 stTypeInd := stType.Elem()

	 if r.rs.rs == nil{
	 	return r.lastError
	 }

	 defer r.rs.rs.Close()

	 v:=reflect.New(stTypeInd)
	 tagList,err :=extractTagInfo(v)
	 if err != nil{
	 	return err
	 }

	 fields,err := r.rs.rs.Columns()

	 if err != nil{
	 	r.rs.lastError = err
	 	return err
	 }

	 refs :=make([]interface{},len(fields))
	 for i,field := range fields{
	 	if f,ok :=tagList[field];ok{
	 		refs[i] = f.Addr().Interface()
		}else{
			refs[i]=new(interface{})
		}
	 }

	 if !r.rs.rs.Next(){
	 	return sql.ErrNoRows
	 }

	 if err :=r.rs.rs.Scan(refs...);err!=nil{
	 	return err
	 }

	 stVal.Elem().Set(v.Elem())
	 return nil
}

type Rows struct {
	rs *sql.Rows
	lastError error
}
func (r *Rows) ToArray()(data [][]string,err error){
	if r.rs==nil{
		return nil,r.lastError
	}
	defer r.rs.Close()

	fields,err := r.rs.Columns()
	if err !=nil{
		r.lastError = err
		return nil,err
	}

	data = make([][]string,0)

	num :=len(fields)

	refs :=make([]interface{},num)
	for i:=0;i<num;i++{
		var ref interface{}
		refs[i]=&ref
	}

	for r.rs.Next(){
		result :=make([]string,len(fields))
		if err:=r.rs.Scan(refs...);err!=nil{
			return nil,err
		}

		for i:=range fields{
			if val,err:=toString(refs[i]);err==nil{
				result[i]=val
			}else{
				return nil,err
			}
		}
		if err !=nil{
			r.lastError = err
			return nil,err
		}
		data = append(data,result)
	}
	return data,nil
}

func (r *Rows) ToMap() (data []map[string]string,err error){
	if r.rs ==nil{
		return nil,r.lastError
	}

	defer r.rs.Close()

	fields,err :=r.rs.Columns()
	if err !=nil{
		r.lastError = err
		return nil,err
	}

	data = make([]map[string]string,0)
	num :=len(fields)

	result :=make(map[string]string)

	refs :=make([]interface{},num)

	for i:=0;i<num;i++{
		var ref interface{}
		refs[i] = &ref
	}

	for r.rs.Next(){
		if err:=r.rs.Scan(refs...);err!=nil{
			return nil,err
		}

		for i,field := range fields{
			if val,err :=toString(refs[i]);err==nil{
				result[field]=val
			}else{
				return nil,err
			}
		}
		data = append(data,result)
	}
	return data,nil
}

func (r *Rows) ToStruct(st interface{}) error {
	stType := reflect.TypeOf(st)

	stVal := reflect.ValueOf(st)
	stValInd := reflect.Indirect(stVal)

	stTypeInd :=stType.Elem()
	if stType.Kind() != reflect.Ptr {
		return fmt.Errorf("this variavle type is %v,not a pointer", stType.Kind())
	}

	if r.rs == nil {
		return r.lastError
	}

	defer r.rs.Close()

	v:=reflect.New(stTypeInd.Elem())

	tagList,err:=extractTagInfo(v)
	if err!=nil{
		return err
	}

	fields,err :=r.rs.Columns()
	if err!=nil{
		r.lastError = err
		return err
	}

	refs:=make([]interface{},len(fields))

	for i,field :=range fields{
		if f,ok :=tagList[field];ok{
			refs[i]=f.Addr().Interface()
		}else{
			refs[i]=new(interface{})
		}
	}

	for r.rs.Next(){
		if err:=r.rs.Scan(refs...);err!=nil{
			return err
		}
		stValInd = reflect.Append(stValInd,v.Elem())
	}
	stVal.Elem().Set(stValInd)
	return nil

}