package server

import (
	"course_system/controller"
	"encoding/gob"
	"errors"
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
	Stub["StuDisenroll"] = controller.GlobalSystem.StuDisenroll
	Stub["StuEnroll"] = controller.GlobalSystem.StuEnroll
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
	f := reflect.ValueOf(Stub[name])
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
