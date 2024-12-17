package controller

import (
	"fmt"
	"slices"
	"sync"
)

type Course struct {
	mu       sync.Mutex
	name     string
	code     string
	capacity int
	students []*Student
}

func NewCourse(code string, name string, cap int) *Course {
	return &Course{name: name,
		code:     code,
		capacity: cap,
		students: make([]*Student, 0, cap)}
}

func (c *Course) AddStudent(std *Student, wg *sync.WaitGroup, errs *[]string, ret *sync.Mutex) error {
	defer wg.Done()
	c.mu.Lock()
	defer c.mu.Unlock()

	if contains(c.students, std) {
		ret.Lock()
		defer ret.Unlock()
		*errs = append(*errs, fmt.Sprintf("%s already in the class %s", std.Name, c.name))
		return fmt.Errorf("%s already in the class", std.Name)
	}

	if len(c.students) >= c.capacity {
		ret.Lock()
		defer ret.Unlock()
		*errs = append(*errs, fmt.Sprintf("%s max capacity reached", c.code))
		return fmt.Errorf("max capacity reached")
	}
	fmt.Printf("course %s add %s\n", c.code, std.Name)
	c.students = append(c.students, std)
	return nil
}

func contains(stu_list []*Student, std *Student) bool {
	for _, v := range stu_list {
		if v.Id == std.Id {
			return true
		}
	}
	return false
}

func (c *Course) DelStudent(std *Student, cg *sync.WaitGroup, errs *[]string, ret *sync.Mutex) error {
	defer cg.Done()
	c.mu.Lock()
	defer c.mu.Unlock()
	if !contains(c.students, std) {
		ret.Lock()
		defer ret.Unlock()
		*errs = append(*errs, fmt.Sprintf("students [%s] does not enrolls the course", std.Name))
		return fmt.Errorf("students [%s] does not enrolls the course", std.Id)
	}
	c.students = slices.DeleteFunc(c.students, func(s *Student) bool {
		return s.Id == std.Id
	})

	return nil
}

func (c *Course) Show() {
	fmt.Printf("name: %s code: %s\n", c.name, c.code)
	fmt.Printf("member:")
	for _, u := range c.students {
		fmt.Printf(" %s ", u.Name)
	}
	fmt.Println()
}
