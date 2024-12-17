package controller

import (
	"github.com/google/uuid"
)

type Student struct {
	Id   uuid.UUID
	Name string
}

func NewStudent(name string) *Student {
	return &Student{
		Id:   uuid.New(),
		Name: name}
}
