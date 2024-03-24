package main

import (
	"fmt"
	"strings"
	"time"
)

type IWebServiceResult interface {
	HasError() bool
	ToString() string
}

type WebServiceResult[T interface{}] struct {
	Result           T `json:"result"`
	Error            string
	Message          string
	ExceptionMessage string
}

type PersonCenter struct {
	CardID        string           `json:"-"`
	UserID        string           `json:"userid"`
	Username      string           `json:"username"`
	UserPhoto     string           `json:"userphoto"`
	Classroom     string           `json:"classroom"`
	ClassName     string           `json:"classname"`
	CourseList    WeekCourse       `json:"courselist"`
	MsgCount      int              `json:"msgcount"`
	Error         string           `json:"error"`
	ExtendOperate []OperateSetting `json:"extendoperate"`
}

func (pc PersonCenter) ToString() {
	fmt.Printf("User ID: %s\n", pc.UserID)
	fmt.Printf("Username: %s\n", pc.Username)
	fmt.Printf("User Photo: %s\n", pc.UserPhoto)
	fmt.Printf("Classroom: %s\n", pc.Classroom)
	fmt.Printf("Class Name: %s\n", pc.ClassName)
	fmt.Printf("Message Count: %d\n", pc.MsgCount)
	fmt.Printf("Error: %s\n", pc.Error)
	for i := 0; i < len(pc.ExtendOperate); i++ {
		fmt.Printf("Extend Operate Settings #%d:\n", i)
		pc.ExtendOperate[i].ToString()
	}
}

type OperateSetting struct {
	BackgroundColor string `json:"BackgroundColor"`
	Index           int    `json:"Index"`
	IsEnable        bool   `json:"IsEnable"`
	IsUpdate        bool   `json:"IsUpdate"`
	OperateName     string `json:"OperateName"`
	OperIcon        string `json:"OperateIcon"`
	OperUrl         string `json:"OperateUrl"`
	FontIcon        string `json:"FontIcon"`
	SystemID        int    `json:"SystemID"`
}

func (os OperateSetting) ToString() {
	fmt.Printf("Background Color: %s\n", os.BackgroundColor)
	fmt.Printf("Index: %d\n", os.Index)
	fmt.Printf("Is Enable: %v\n", os.IsEnable)
	fmt.Printf("Is Update: %v\n", os.IsUpdate)
	fmt.Printf("Operate Name: %s\n", os.OperateName)
	fmt.Printf("Operate Icon: %s\n", os.OperIcon)
	fmt.Printf("Operate URL: %s\n", os.OperUrl)
	fmt.Printf("Font Icon: %s\n", os.FontIcon)
	fmt.Printf("System ID: %d\n", os.SystemID)
}

type WeekCourse struct {
	ShowSaturday bool
	ShowSunday   bool
	ShowTeacher  bool
	Title        string `json:"title"`
	SubTitle     string `json:"subtitle"`
	StartDate    string
	EndDate      string
	Days         int       `json:"days,omitempty"`
	Sections     []Section `json:"sections"`
}

type Section struct {
	Name      string   `json:"name"`
	Index     int      `json:"index"`
	Courses   []Course `json:"courses"`
	StartTime string
	EndTime   string
}

type Course struct {
	Name        string `json:"name"`
	Week        int    `json:"week"`
	Teacher     string `json:"teacher"`
	TypeID      int
	Show        bool          `json:"show"`
	ClassRoom   string        `json:"classroom"`
	StartTime   time.Duration `json:"StartTime,omitempty"`
	EndTime     time.Duration `json:"EndTime,omitempty"`
	KQStartTime time.Duration `json:"KQStartTime,omitempty"`
	KQEndTime   time.Duration `json:"KQEndTime,omitempty"`
}

func (result WebServiceResult[T]) HasError() bool {
	return len(result.Error) != 0 || len(result.Message) != 0
}

func (result WebServiceResult[T]) ToString() string {
	s := []string{
		"Error->",
		result.Error,
		";Message->",
		result.Message,
		";ExceptionMessage->",
		result.ExceptionMessage,
	}
	return strings.Join(s, "")
}
