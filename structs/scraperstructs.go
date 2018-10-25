package ScraperData

type ScraperData struct {
	Teachers  map[string]Teacher                                          `json:"Teachers"`
	Semesters map[string]map[string]map[string][]map[string]CourseSection `json:"Semesters"`
}

type Coursection struct {
	Open                 bool              `json:"Open"`
	AcademicLevel        string            `json:"AcademicLevel"`
	CourseCode           string            `json:"CourseCode"`
	CourseDescription    string            `json:"CourseDescription"`
	CourseName           string            `json:"CourseName"`
	DateStart            string            `json:"DateStart"`
	DateEnd              string            `json:"DateEnd"`
	Location             string            `json:"Location"`
	MeetingInformation   string            `json:"MeetingInformation"`
	Supplies             string            `json:"Supplies"`
	Credits              float64           `json:"Credits"`
	SlotsAvailable       int64             `json:"SlotsAvailable"`
	SlotsCapacity        int64             `json:"SlotsCapacity"`
	SlotsWaitlist        int64             `json:"SlotsWaitlist"`
	TimeEnd              int64             `json:"TimeEnd"`
	TimeStart            int64             `json:"TimeStart"`
	ProfessorEmails      []string          `json:"ProfessorEmails"`
	PrereqNonCourse      string            `json:"PrereqNonCourse"`
	RecConcurrentCourses []string          `json:"RecConcurrentCourses"`
	ReqConcurrentCourses []string          `json:"ReqConcurrentCourses"`
	PrereqCourses        PrereqCourses     `json:"PrereqCourses"`
	InstructionalMethods map[string]string `json:"InstructionalMethods"`
}

type PrereqCourses struct {
	And []string      `json:"and"`
	Or  []interface{} `json:"or"`
}

type Teacher struct {
	Email string `json:"Email"`
	Name  string `json:"Name"`
	Phone string `json:"Phone"`
}
