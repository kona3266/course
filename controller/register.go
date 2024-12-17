package controller

import (
	"sync"
)

type CentralSystem struct {
	mu            sync.Mutex
	course_member map[string]*Course
}

/*
s1 s1 s1 s3 s2
*/
var GlobalSystem = CentralSystem{course_member: make(map[string]*Course)}
var stuLock = NewMultipleLock()

func (s *CentralSystem) RegisterCourse(c *Course) error {
	s.mu.Lock()
	s.course_member[c.code] = c
	s.mu.Unlock()
	return nil
}

func (s *CentralSystem) StuEnroll(i *Student, courses []string) []string {
	stuLock.Lock(i.Id)
	defer stuLock.Unlock(i.Id)
	var errs []string
	var cg sync.WaitGroup
	var ret sync.Mutex
	for _, code := range courses {
		course := s.course_member[code]
		cg.Add(1)
		go course.AddStudent(i, &cg, &errs, &ret)
	}
	cg.Wait()

	return errs
}

func (s *CentralSystem) StuDisenroll(i *Student, courses []string) []string {
	stuLock.Lock(i.Id)
	defer stuLock.Unlock(i.Id)
	var cg sync.WaitGroup
	var errs []string
	var ret sync.Mutex
	for _, code := range courses {
		course := s.course_member[code]
		cg.Add(1)
		go course.DelStudent(i, &cg, &errs, &ret)
	}
	cg.Wait()

	return errs
}

func (s *CentralSystem) ShowAll() {
	for _, course := range s.course_member {
		course.Show()
	}
}
