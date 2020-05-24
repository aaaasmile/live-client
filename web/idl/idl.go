package idl

import (
	"database/sql"
	"time"
)

var (
	Appname         = "LiveClient"
	Buildnr         = "00.01.20200524-00"
	LocalServiceURL = ""
	WebServiceName  = "LiveClient"
)

const ()

type ServiceHandlerSettings struct {
	ConfigFile    string
	ServiceName   string
	Command       string
	AutoStartPage bool
}

type ResErr struct {
	Err error
}

type SourceFileWithErr struct {
	SourceFile SourceFile
	Err        error
}
type ChanSourceFiles chan SourceFileWithErr

type SourceFile struct {
	DbLiteID    int
	Name        string
	VersionList string
	Checksum    string
	Filename    string
	FileModTime time.Time
	FileSize    int
}

type ObjTypeInProv int

func (ot *ObjTypeInProv) String() string {
	switch *ot {
	case OTPSourceFile:
		return "SourceFile"
	}
	return ""
}

const (
	OTPNavObj ObjTypeInProv = iota
	OTPSourceFile
	OTPIgnorelist
)

type ObjProvider interface {
	DoReadAllObj(ObjTypeInProv) ([]*ObjectInfo, error)
	GetNewTransaction() (*sql.Tx, error)
	DoInsertObject(*sql.Tx, *ObjectInfo, ObjTypeInProv) error
	DoUpdateObject(*sql.Tx, *ObjectInfo, ObjTypeInProv) error
	DoDeleteObject(*sql.Tx, *ObjectInfo, ObjTypeInProv) error
}
