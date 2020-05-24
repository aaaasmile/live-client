package idl

import (
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
	ObjectID    string
	Name        string
	VersionList string
	Checksum    string
	Filename    string
	FileModTime time.Time
	FileSize    int
}
