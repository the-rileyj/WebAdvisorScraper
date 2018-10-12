package structs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// CourseInfomation is the data structure representing a given course
type CourseInfomation struct {
	Open                 bool                `json:"open"`
	AcademicLevel        string              `json:"academicLevel"`
	CourseCode           string              `json:"courseCode"`
	CourseDescription    string              `json:"courseDescription"`
	CourseName           string              `json:"courseName"`
	DateStart            string              `json:"dateStart"`
	DateEnd              string              `json:"dateEnd"`
	Location             string              `json:"location"`
	MeetingInformation   string              `json:"meetingInformation"`
	PrereqNonCourse      string              `json:"prereqNonCourse"`
	Supplies             string              `json:"supplies"`
	Credits              int                 `json:"credits"`
	SlotsAvailable       int                 `json:"slotsAvailable"`
	SlotsCapacity        int                 `json:"slotsCapacity"`
	SlotsWaitlist        int                 `json:"slotsWaitlist"`
	TimeEnd              int                 `json:"timeEnd"`
	TimeStart            int                 `json:"timeStart"`
	ProfessorEmails      []string            `json:"professorEmails"`
	RecConcurrentCourses []string            `json:"recConcurrentCourses"`
	ReqConcurrentCourses []string            `json:"reqConcurrentCourses"`
	PrereqCourses        map[string][]string `json:"prereqCourses"`
	InstructionalMethods string              `json:"instructionalMethods"`
}

// RJSignal struct for indicating that shell process finished properly and the main process can stop listening for CTRL+C so it can kill that process */
type RJSignal struct{}

// String is for satisfying the os.Signal interface
func (rjs RJSignal) String() string { return "RJ" }

// Signal is for satisfying the os.Signal interface
func (rjs RJSignal) Signal() {}

// Struct for holding SDBOR login credentials
type info struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetCredentials reads in the json at 'path' and returns an info struct
func GetCredentials(path string) (info, error) {
	credentialsFile, err := os.Open(path)

	if err != nil {
		return info{}, fmt.Errorf("could not open credentials file at %s", path)
	}

	var credentials info

	err = json.NewDecoder(credentialsFile).Decode(&credentials)

	if err != nil {
		return info{}, errors.New("could not unmarshal file into struct")
	}

	return credentials, nil
}

// TeacherInformation is the data structure representing a DSU teacher
type TeacherInformation struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}
