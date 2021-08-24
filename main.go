// To execute Go code, please declare a func main() in a package "main"

package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//https://golang.org/pkg/sort/

//enum
const (
	Spring = "spring"
	Winter = "winter"
	Summer = "summer"
	Fall   = "fall"
	S      = "s"
	Su     = "su"
	W      = "w"
	F      = "f"
)

//models
type Course struct {
	Department   string `json:"department"`
	CourseNumber string `json:"course_number"`
	Year         int    `json:"year"`
	Semester     string `json:"semester"`
}

type Semester string

//main
func main() {
	inputData := []string{
		"CS111ss 2016 Fall",
		"CS-111 Fall 2016",
		"CS~111 Fall 2016",
		"MATH 123 2015 Spring",
		"MATH cs 2015 Spring",
		"CS-222 Fall2017",
		"CS-222 F2018",
		"CS333 Su2019",
		"CS 333 S2020",
		"CS 444 Winter 21",
		"CS 444 Spring2022",
		"CS 444 Spring10000",
		"CS 444 Sprin2022",
		"CS 444 Invalid dd 2022",
		"qww1111124 3434 sd",
		"111CS Fall 2016",
		"CS 444 Winter 99",
		"Chem 444 Winter su88",
		"Physics1234 202S2",
		"Physics1234 2022s22",
		"Physics1234 202su2",
		"Physics1234 Su 2021",
	}
	var departmentCourse map[string]Course
	var keys []string
	departmentCourse = make(map[string]Course, len(inputData))
	for _, input := range inputData {
		course, err := parseCourse(input)
		if err == nil {
			// println(course.String())
			key := fmt.Sprintf("%s-%s", course.Department, course.CourseNumber)

			departmentCourse[key] = course
		}
	}

	for data, _ := range departmentCourse {
		keys = append(keys, data)
	}
	sort.Strings(keys)

	for _, key := range keys {
		//  println(key)
		println(departmentCourse[key].String())

	}

}

// main method
func parseCourse(input string) (Course, error) {
	var course Course
	cols := strings.Split(input, " ")
	if err := validateColumns(input, cols); err != nil {
		return Course{}, err
	}
	var courseNumberFound, semesterFound, yearFound bool
	for index, col := range cols {
		if index == 0 {
			if err := segmentZero(input, col, &course, &courseNumberFound); err != nil {
				return Course{}, err
			}
		} else if index == 1 && !courseNumberFound {
			if err := segmentOne(input, col, &course, &courseNumberFound); err != nil {
				return Course{}, err
			}
		} else if index == 1 || index == 2 || index == 3 {
			if err := segmentBeyondOne(input, col, &course, &semesterFound, &yearFound); err != nil {
				return Course{}, err
			}
		}
	}
	return course, nil
}

//utilities
func (c Course) String() string {
	return fmt.Sprintf("\n[\n Department : %s \n Course Number : %s \n Year : %d \n Semester : %s \n] \n", c.Department, c.CourseNumber, c.Year, c.Semester)
}

func (sem Semester) IsValid() (string, error) {
	switch sem {
	case Su, Summer:
		return Summer, nil
	case S, Spring:
		return Spring, nil
	case W, Winter:
		return Winter, nil
	case F, Fall:
		return Fall, nil
	}
	return "", errors.New("invalid Semester")
}

func IsValidYear(year *int) error {
	if (*year < 0 || *year > 99) && (*year < 1900 || *year > 2100) {
		return errors.New("invalid year")
	}
	if *year > 0 && *year <= 40 {
		*year += 2000
	} else if *year > 40 && *year <= 99 {
		*year += 1900
	}
	return nil
}

func segregateDepartmentCourse(s string) (letters string, numbers string, err error) {
	var l, n []rune
	isNumberFound := false
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z':
			if !isNumberFound {
				l = append(l, r)
			} else {
				return "", "", errors.New("invalid department/course number")
			}
		case r >= '0' && r <= '9':
			isNumberFound = true
			n = append(n, r)
		case r == ' ' || r == ':' || r == '-':
		default:
			return "", "", errors.New("invalid department/course number")
		}
	}
	return string(l), string(n), nil
}

func segregateSemYear(s string, semesterFound *bool, yearFound *bool) (letters string, numbers string, err error) {
	var l, n []rune
	var isNumber bool
	var isAlpha bool
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z':
			if !*semesterFound && !isNumber {
				l = append(l, r)
				isAlpha = true
			} else {
				return "", "", errors.New("invalid alphabets")
			}
		case r >= '0' && r <= '9':
			if !*yearFound {
				n = append(n, r)
				isNumber = true
			} else {
				return "", "", errors.New("invalid numbers")
			}

		}
	}
	if isNumber && isAlpha && (len(l)+len(n) != len(s)) {
		return "", "", errors.New("invalid Sem/Year")
	}
	return string(l), string(n), nil
}

func populateSemester(semester string, line string, course *Course, semesterFound *bool) error {
	sem := Semester(strings.ToLower(semester))
	s, err := sem.IsValid()
	if err != nil {
		return errors.New(fmt.Sprintf("%s %s", "Skipping-Invalid semester -->", line))
	}
	course.Semester = strings.Title(s)
	*semesterFound = true
	return nil
}

func populateYear(year string, line string, course *Course, yearFound *bool) error {
	yearInt, _ := strconv.Atoi(year)
	if err := IsValidYear(&yearInt); err != nil {
		return errors.New(fmt.Sprintf("%s %s", "Skipping-Invalid Year -->", line))
	}
	course.Year = yearInt
	*yearFound = true
	return nil
}

func segmentOne(input string, col string, course *Course, courseNumberFound *bool) error {
	letters, numbers, err := segregateDepartmentCourse(col)
	if err != nil {
		return errors.New(fmt.Sprintf("%s %s", "Skipping-Invalid Department or Course number --> ", input))
	}
	if letters == "" {
		course.CourseNumber = numbers
		*courseNumberFound = true
	} else {
		return errors.New(fmt.Sprintf("%s %s", "Skipping-course number can not have alphabets -->", input))

	}
	return nil
}

func segmentZero(input string, col string, course *Course, courseNumberFound *bool) error {
	department, courseNumber, err := segregateDepartmentCourse(col)
	if err != nil {
		return errors.New(fmt.Sprintf("%s %s", "Skipping-Invalid Department or Course number --> ", input))
	}
	if courseNumber != "" && department != "" {
		course.Department = department
		course.CourseNumber = courseNumber
		*courseNumberFound = true
	} else {
		course.Department = department
	}
	return nil
}

func segmentBeyondOne(input string, col string, course *Course, semesterFound *bool, yearFound *bool) error {
	semester, year, err := segregateSemYear(col, semesterFound, yearFound)
	if err != nil {
		return errors.New(fmt.Sprintf("%s %s", "Skipping-Invalid Semester or Year --> ", input))
	}
	if semester != "" && year != "" {
		if err := populateSemester(semester, input, course, semesterFound); err != nil {
			return err
		}
		if err := populateYear(year, input, course, yearFound); err != nil {
			return err
		}
	} else if semester != "" {
		if err := populateSemester(semester, input, course, semesterFound); err != nil {
			return err
		}

	} else if year != "" {
		if err := populateYear(year, input, course, yearFound); err != nil {
			return err
		}
	}
	return nil
}

func validateColumns(input string, cols []string) error {
	if len(cols) > 4 {
		return errors.New(fmt.Sprintf("%s %s", "Skipping-rows can not exceed more than 4 columns -->", input))
	}
	return nil
}
