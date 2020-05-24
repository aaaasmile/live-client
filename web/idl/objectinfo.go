package idl

import (
	"database/sql"
	"time"
)

type ObjTypeInProv int

func (ot *ObjTypeInProv) String() string {
	switch *ot {
	case OTPSourceFile:
		return "SourceFile"
	}
	return ""
}

const (
	OTPSourceFile ObjTypeInProv = iota
)

type ObjProvider interface {
	DoReadAllObj(ObjTypeInProv) ([]*ObjectInfo, error)
	GetNewTransaction() (*sql.Tx, error)
	DoInsertObject(*sql.Tx, *ObjectInfo, ObjTypeInProv) error
	DoUpdateObject(*sql.Tx, *ObjectInfo, ObjTypeInProv) error
	DoDeleteObject(*sql.Tx, *ObjectInfo, ObjTypeInProv) error
}

type ObjectInfoColl []*ObjectInfo

func (a ObjectInfoColl) Len() int      { return len(a) }
func (a ObjectInfoColl) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ObjectInfoColl) Less(i, j int) bool {
	return a[i].Key < a[j].Key
}

type ObjectInfo struct {
	Key         string
	Name        string
	VersionList string
	Checksum    string
	IDInDb      int
	Timestamp   time.Time
	SourceFile  SourceFile
}

func NewObjectInfoFromSF(sf SourceFile) *ObjectInfo {
	oi := ObjectInfo{
		Key:         sf.ObjectID,
		Name:        sf.Name,
		VersionList: sf.VersionList,
		Checksum:    sf.Checksum,
		Timestamp:   sf.FileModTime,
		SourceFile:  sf,
	}

	return &oi
}

func (oi *ObjectInfo) IsEqual(other *ObjectInfo) bool {
	return (oi.Name == other.Name) &&
		(oi.VersionList == other.VersionList) &&
		((oi.Checksum == "") || (other.Checksum == "") || (oi.Checksum == other.Checksum)) &&
		(oi.Timestamp.Unix() == other.Timestamp.Unix())
}

type ObjOpChangeType int

const (
	OOCTinsert = iota
	OOCTupdate
	OOCTconfirm
	OOCTdelete
)

type ObjInfoChange struct {
	ChangeType ObjOpChangeType
	Obj        *ObjectInfo
}
