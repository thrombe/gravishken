package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

// this is a custom error type for use throughout the app
type Error struct {
	message string
}

func NewError(msg string) Error {
	return Error{message: msg}
}

func (self Error) Error() string {
	return fmt.Sprintf("Error: %s", self.message)
}

type Varient int

const (
	Err Varient = iota
	Notification
	ExeNotFound
	Quit
	UserLoginRequest
	WarnUser
	LoadRoute
	ReloadUi
	StartTest
	TestFinished
	CheckSystem
	OpenApp
	QuitApp
	Unknown // NOTE: keep this as the last constant here.
)

func (self Varient) TSName() string {
	switch self {
	case Err:
		return "Err"
	case Notification:
		return "Notification"
	case ExeNotFound:
		return "ExeNotFound"
	case Quit:
		return "Quit"
	case UserLoginRequest:
		return "UserLoginRequest"
	case WarnUser:
		return "WarnUser"
	case LoadRoute:
		return "LoadRoute"
	case ReloadUi:
		return "ReloadUi"
	case StartTest:
		return "StartTest"
	case TestFinished:
		return "TestFinished"
	case CheckSystem:
		return "CheckSystem"
	case OpenApp:
		return "OpenApp"
	case QuitApp:
		return "QuitApp"
	default:
		return "Unknown"
	}
}
func varientFromName(typ string) Varient {
	switch typ {
	case "Err":
		return Err
	case "Notification":
		return Notification
	case "ExeNotFound":
		return ExeNotFound
	case "Quit":
		return Quit
	case "UserLoginRequest":
		return UserLoginRequest
	case "WarnUser":
		return WarnUser
	case "ReloadUi":
		return ReloadUi
	case "LoadRoute":
		return LoadRoute
	case "StartTest":
		return StartTest
	case "TestFinished":
		return TestFinished
	case "CheckSystem":
		return CheckSystem
	case "OpenApp":
		return OpenApp
	case "QuitApp":
		return QuitApp
	default:
		return Unknown
	}
}

// only for unexpected errors / for errors that we can't do much about, other than telling the user about it
type TErr struct {
	Message string
}

type TNotification struct {
	Message string
	Typ     string
}

type Message struct {
	Typ Varient
	Val string
}

type TCheckSystem struct{}

type TExeNotFound struct {
	Name   string
	ErrMsg string
}

type TQuit struct{}

type TUserLoginRequest struct {
	Username string
	Password string
}
type TWarnUser struct {
	Message string
}

type TLoadRoute struct {
	Route string
}

type TReloadUi struct{}

type TStartTest struct{}

// TODO: Send this to every client from server when time ends and submit whatever stuff the client can submit
type TTestFinished struct{}

type AppType int

const (
	TXT AppType = iota
	DOCX
	XLSX
	PPTX
)

func (self AppType) TSName() string {
	switch self {
	case TXT:
		return "TXT"
	case DOCX:
		return "DOCX"
	case XLSX:
		return "XLSX"
	case PPTX:
		return "PPTX"
	default:
		return "Unknown"
	}
}

type TOpenApp struct {
	Typ    AppType
	TestId ID `ts_type:"string"`
}

type TQuitApp struct{}

func NewMessage(typ interface{}) Message {
	name := reflect.TypeOf(typ).Name()[1:]
	varient := varientFromName(name)
	json, err := json.Marshal(typ)
	if err != nil {
		panic(err)
	}
	return Message{
		Typ: varient,
		Val: string(json),
	}
}

func Get[T any](msg Message) (*T, error) {
	var val T

	name := reflect.TypeOf(val).Name()[1:]
	if name != msg.Typ.TSName() {
		err_msg := fmt.Sprintf("message of type '%s' but asked to be decoded as '%s'", msg.Typ.TSName(), name)
		return nil, NewError(err_msg)
	}

	err := json.Unmarshal([]byte(msg.Val), &val)
	return &val, err
}

// - [tkrajina/tkypescriptify-golang-structs](https://github.com/tkrajina/typescriptify-golang-structs)
func DumpTypes(dir string) {
	allVarients := make([]Varient, Unknown+1)

	for i := 0; i <= int(Unknown); i++ {
		allVarients[i] = Varient(i)
	}

	converter := typescriptify.New().
		WithInterface(true).
		WithBackupDir("").
		Add(TErr{}).
		Add(TNotification{}).
		Add(Message{}).
		Add(TExeNotFound{}).
		Add(TQuit{}).
		Add(TUserLoginRequest{}).
		Add(TWarnUser{}).
		Add(TLoadRoute{}).
		Add(TReloadUi{}).
		Add(TStartTest{}).
		Add(TTestFinished{}).
		Add(TCheckSystem{}).
		Add(TOpenApp{}).
		Add(TQuitApp{}).
		AddEnum([]AppType{TXT, DOCX, XLSX, PPTX}).
		AddEnum(allVarients)

	converter = converter.
		Add(User{}).
		Add(TestSubmission{}).
		Add(TypingTestInfo{}).
		Add(AppTestInfo{}).
		Add(McqTestInfo{}).
		// Add(UserBatchRequestData{}).
		Add(Test{}).
		Add(Admin{}).
		// Add(AdminRequest{}).
		Add(Batch{}).
		AddEnum([]TestType{TypingTest, DocxTest, ExcelTest, PptTest, MCQTest})

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err.Error())
	}
	err = converter.ConvertToFile(filepath.Join(dir, "types.ts"))
	if err != nil {
		panic(err.Error())
	}
}
