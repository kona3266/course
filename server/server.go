package server

import (
	"course_system/controller"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
)

type Msg struct {
	Method string
	S      *controller.Student
	Course []string
}

type Rep struct {
	Code int
	Err  []string
}

type funcMap map[string]interface{}

var Stub = funcMap{}

func init() {

}

func RegisterMethod(o interface{}) {
	v := reflect.ValueOf(o)
	fmt.Printf("%d\n", v.NumMethod())

	for i := 0; i < v.NumMethod(); i++ {
		m := v.Method(i)
		if m.Type().NumOut() != 1 {
			continue
		}
		if m.Type().Out(0).Kind() != reflect.Slice {
			continue
		}
		if m.Type().Out(0).Elem().Kind() != reflect.String {
			continue
		}

		methodName := v.Type().Method(i).Name
		Stub[methodName] = m
		fmt.Println("regist", methodName)
	}
}

func Handle(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)
	var msg Msg
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal(err)
	}
	method := msg.Method
	rsp, _ := Call(method, msg.S, msg.Course)
	e := rsp[0].Interface().([]string)
	if len(e) != 0 {
		rep := Rep{Code: -1, Err: e}
		err = enc.Encode(&rep)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		rep := Rep{Code: 1, Err: nil}
		err = enc.Encode(&rep)
		if err != nil {
			log.Fatal(err)
		}
	}
	conn.Close()
}

func Call(name string, params ...interface{}) (result []reflect.Value, err error) {
	f := Stub[name].(reflect.Value)
	if len(params) != f.Type().NumIn() {
		err = errors.New("the number of params is out of index")
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	result = f.Call(in)

	return
}
