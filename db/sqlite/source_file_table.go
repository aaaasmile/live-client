package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aaaasmile/live-client/web/idl"
)

type SourceFileTable struct {
	Fields       idl.SourceFile
	DebugVerbose bool
	prjID        int
	connDb       *sql.DB
}

func NewSourceFileTable(prjID int, connDb *sql.DB) *SourceFileTable {
	res := SourceFileTable{
		prjID:  prjID,
		connDb: connDb,
	}
	return &res
}

func (tb *SourceFileTable) ReadIDByFileInfo(fileInfo os.FileInfo) (bool, error) {
	q := `SELECT id,FileModTime,FileSize,Type,ObjectID  FROM SourceFile WHERE  Filename = ?  AND ProjectID = ?;`
	found := false
	// fileModTime is stored as Unix timestamp
	var fileModTime, fileSize string
	var id, typeObj, objectID int
	err := tb.connDb.QueryRow(q, fileInfo.Name(), tb.prjID).Scan(&id, &fileModTime, &fileSize, &typeObj, &objectID)

	if err != nil {
		if err == sql.ErrNoRows {
			if tb.DebugVerbose {
				log.Println("No source file found")
			}
			return found, nil
		}
	} else {
		mti, err := strconv.ParseInt(fileModTime, 10, 64)
		if err != nil {
			return false, err
		}
		mt := time.Unix(mti, 0)
		size, err := strconv.Atoi(fileSize)
		if err != nil {
			return false, err
		}
		tb.Fields.DbLiteID = id
		tb.Fields.FileModTime = mt
		tb.Fields.FileSize = size
		tb.Fields.Filename = fileInfo.Name()
		tb.Fields.Type = typeObj
		tb.Fields.ObjectID = objectID
		found = true
	}
	//fmt.Println("*** fields", tb.Fields)
	return found, err

}

func (tb *SourceFileTable) ReadOtherFieldsByID() error {
	if tb.Fields.DbLiteID == 0 {
		return fmt.Errorf("Record id was not fetched. This is a second pass function")
	}
	q := `SELECT Name,Modified,VersionList,Checksum,Timestamp  FROM SourceFile WHERE  id = ?;`
	// fileModTime is stored as Unix timestamp
	var name, versionList, checksum string
	var modified int
	var timestamp int64
	err := tb.connDb.QueryRow(q, tb.Fields.DbLiteID).Scan(&name, &modified, &versionList, &checksum, &timestamp)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Record %d not found", tb.Fields.DbLiteID)
		}
	} else {
		mt := time.Unix(timestamp, 0)
		tb.Fields.Name = name
		tb.Fields.Modified = modified
		tb.Fields.VersionList = versionList
		tb.Fields.Checksum = checksum
		tb.Fields.Timestamp = mt
	}
	//fmt.Println("*** fields in ReadFieldsByID", tb.Fields)
	return err

}
