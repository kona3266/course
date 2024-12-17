package main

import (
	"course_system/controller"
	"course_system/server"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	c1 := controller.NewCourse("CS61A", "Structure and Interpretation of Computer Programs", 1)
	c2 := controller.NewCourse("MATH101", "Linear algebra", 3)
	c3 := controller.NewCourse("MATH102", "Calculus", 3)

	controller.GlobalSystem.RegisterCourse(c1)
	controller.GlobalSystem.RegisterCourse(c2)
	controller.GlobalSystem.RegisterCourse(c3)

	port := flag.Int("port", 8000, "listen port")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	flag.Parse()
	fmt.Printf("listen on %d \n", *port)
	address := fmt.Sprintf(":%d", *port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	go func() {
		<-sigs
		controller.GlobalSystem.ShowAll()
		os.Exit(0)
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go server.Handle(conn)
	}

}
