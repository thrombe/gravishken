package types

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
	ExeNotFound Varient = iota
	UserLogin
	LoadRoute
	ReloadUi
	Err
	Unknown // NOTE: keep this as the last constant here.
)

func (self Varient) TSName() string {
	switch self {
	case ExeNotFound:
		return "ExeNotFound"
	case UserLogin:
		return "UserLogin"
	case LoadRoute:
		return "LoadRoute"
	case ReloadUi:
		return "ReloadUi"
	case Err:
		return "Err"
	default:
		return "Unknown"
	}
}
func varientFromName(typ string) Varient {
	switch typ {
	case "ExeNotFound":
		return ExeNotFound
	case "UserLogin":
		return UserLogin
	case "ReloadUi":
		return ReloadUi
	case "LoadRoute":
		return LoadRoute
	case "Err":
		return Err
	default:
		return Unknown
	}
}

type Message struct {
	Typ Varient
	Val string
}

type TExeNotFound struct {
	Name   string
	ErrMsg string
}

type TUserLogin struct {
	Username string
	Password string
	TestCode string
}

type TLoadRoute struct {
	Route string
}

type TReloadUi struct{}

// only for unexpected errors / for errors that we can't do much about, other than telling the user about it
type TErr struct {
	Message string
}

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
	for i := range Unknown + 1 {
		allVarients[i] = i
	}

	converter := typescriptify.New().
		WithInterface(true).
		WithBackupDir("").
		Add(Message{}).
		Add(TExeNotFound{}).
		Add(TUserLogin{}).
		Add(TLoadRoute{}).
		Add(TReloadUi{}).
		Add(TErr{}).
		AddEnum(allVarients)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err.Error())
	}
	err = converter.ConvertToFile(filepath.Join(dir, "types.ts"))
	if err != nil {
		panic(err.Error())
	}
}
