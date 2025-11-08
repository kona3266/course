package main

import (
	"course_system/controller"
	"course_system/server"
	"encoding/gob"
	"fmt"
	"net"
	"sync"
)

func main() {

	var c sync.WaitGroup
	wait_list := []string{"CS61A", "MATH101", "MATH102"}
	stu_list := []*controller.Student{}
	for i := 0; i < 5; i++ {
		name := fmt.Sprint("Eminem", i)
		std := controller.NewStudent(name)
		stu_list = append(stu_list, std)
	}

	for i := 0; i < 5; i++ {
		c.Add(1)
		go func(c *sync.WaitGroup) {
			conn, err := net.Dial("tcp", "localhost:8000")
			if err != nil {
				panic(err)
			}
			enc := gob.NewEncoder(conn)

			j := i % 3
			var list []string
			if j == 2 {
				list = wait_list[:2]
			} else if j == 1 {
				list = []string{wait_list[0], wait_list[2]}
			} else {
				list = wait_list[1:]
			}

			k := i % 3
			msg := server.Msg{
				Method: "StuEnroll",
				S:      stu_list[k],
				Course: list,
			}
			fmt.Printf("Student %s enroll list %s\n", stu_list[k].Name, list)
			enc.Encode(msg)
			dec := gob.NewDecoder(conn)
			rsp := server.Rep{}
			dec.Decode(&rsp)
			fmt.Println(rsp)
			c.Done()
		}(&c)

	}

	c.Wait()

	for i := 0; i < 5; i++ {
		c.Add(1)
		go func(c *sync.WaitGroup) {
			conn, err := net.Dial("tcp", "localhost:8000")
			if err != nil {
				panic(err)
			}
			enc := gob.NewEncoder(conn)

			j := i % 3
			var list []string
			if j == 2 {
				list = wait_list[:2]
			} else if j == 1 {
				list = []string{wait_list[0], wait_list[2]}
			} else {
				list = wait_list[1:]
			}

			k := i % 3
			msg := server.Msg{
				Method: "StuDisenroll",
				S:      stu_list[k],
				Course: list,
			}
			fmt.Printf("Student %s disenroll list %s\n", stu_list[k].Name, list)
			enc.Encode(msg)
			dec := gob.NewDecoder(conn)
			rsp := server.Rep{}
			dec.Decode(&rsp)
			fmt.Println(rsp)
			c.Done()
		}(&c)

	}
	c.Wait()

}
