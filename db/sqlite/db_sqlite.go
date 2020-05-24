package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aaaasmile/live-client/util"
	"github.com/aaaasmile/live-client/web/idl"
	_ "github.com/mattn/go-sqlite3"
)

type LiteDB struct {
	connDb       *sql.DB
	DebugSQL     bool
	SqliteDBPath string
}

func (ld *LiteDB) GetConnDB() *sql.DB {
	return ld.connDb
}

func (ld *LiteDB) OpenSqliteDatabase() error {
	var err error
	// Source control should be only an empty navrepo.db.
	dbname := util.GetFullPath(ld.SqliteDBPath)
	log.Println("Using the sqlite file: ", dbname)
	ld.connDb, err = sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	return nil
}
func (ld *LiteDB) GetNewTransaction() (*sql.Tx, error) {
	tx, err := ld.connDb.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

//  Interface ObjProvider - start

func (ld *LiteDB) DoReadAllObj(ot idl.ObjTypeInProv) ([]*idl.ObjectInfo, error) {
	switch ot {
	case idl.OTPSourceFile:
		return ld.readAllFile("SourceFile")
	case idl.OTPServerFile:
		return ld.readAllFile("ServerFile")
	}
	return nil, fmt.Errorf("Type povider not recgonized ", ot)
}

func (ld *LiteDB) DoInsertObject(tx *sql.Tx, obj *idl.ObjectInfo, ot idl.ObjTypeInProv) error {
	switch ot {
	case idl.OTPServerFile:
		return ld.insertFile("ServerFile", tx, &obj.SourceFile)
	case idl.OTPSourceFile:
		return ld.insertFile("SourceFile", tx, &obj.SourceFile)
	}
	return fmt.Errorf("Type povider not recgonized %v", ot)
}

func (ld *LiteDB) DoUpdateObject(tx *sql.Tx, obj *idl.ObjectInfo, ot idl.ObjTypeInProv) error {
	switch ot {
	case idl.OTPServerFile:
		return ld.updateFile("ServerFile", tx, &obj.SourceFile)
	case idl.OTPSourceFile:
		return ld.updateFile("SourceFile", tx, &obj.SourceFile)
	}
	return fmt.Errorf("Type povider not recgonized %v", ot)
}

func (ld *LiteDB) DoDeleteObject(tx *sql.Tx, obj *idl.ObjectInfo, ot idl.ObjTypeInProv) error {
	switch ot {
	case idl.OTPSourceFile:
		return ld.deleteSourceFile(tx, obj.SourceFile.DbLiteID)
	}
	return fmt.Errorf("Type povider not recgonized %v", ot)
}

//  Interface ObjProvider - end

func (ld *LiteDB) deleteSourceFile(tx *sql.Tx, recID int) error {
	q := fmt.Sprintf(`DELETE FROM SourceFile WHERE id=%d;`, recID)
	if ld.DebugSQL {
		log.Println("Query is", q)
	}

	stmt, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(stmt).Exec()
	return err
}

func (ld *LiteDB) updateFile(tableName string, tx *sql.Tx, srcItem *idl.SourceFile) error {
	//fmt.Println("** oi updateFile ", srcItem.FileModTime, srcItem.FileModTime.Local().Unix(), srcItem)
	var q string
	recID := srcItem.DbLiteID
	if recID == 0 {
		panic("Recid could not be null")
	}

	q = fmt.Sprintf(`UPDATE %s SET VersionList='%s',Filename='%s',Checksum='%s',FileModTime=%d,FileSize=%d WHERE id=%d;`,
		tableName,
		srcItem.VersionList, srcItem.Filename, srcItem.Checksum, srcItem.FileModTime.Local().Unix(), srcItem.FileSize,
		recID)

	if ld.DebugSQL {
		log.Println("Query is", q)
	}
	updateMore, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	res, err := tx.Stmt(updateMore).Exec()
	if err != nil {
		log.Println("Error in updateFile")
		return err
	} else {
		if ld.DebugSQL {
			log.Println("Update result", res)
		}
	}

	if ld.DebugSQL {
		log.Println("SourceFile updated OK: ", srcItem.Name)
	}
	return nil
}

func (ld *LiteDB) readAllFile(tableName string) ([]*idl.ObjectInfo, error) {
	q := `SELECT id,Name,VersionList,Checksum,ObjectID,FileModTime,FileSize,Filename FROM %s;`
	q = fmt.Sprintf(q, tableName)
	if ld.DebugSQL {
		log.Println("Query is", q)
	}
	rows, err := ld.connDb.Query(q)
	if err != nil {
		return nil, err
	}
	result := make([]*idl.ObjectInfo, 0)
	defer rows.Close()
	for rows.Next() {
		item := idl.SourceFile{}
		var fileModTime, fileSize string
		if err := rows.Scan(&item.DbLiteID, &item.Name,
			&item.VersionList, &item.Checksum, &item.ObjectID, &fileModTime, &fileSize, &item.Filename); err != nil {
			log.Println("Error in scan lite src ", err)
			return nil, err
		}
		mti, err := strconv.ParseInt(fileModTime, 10, 64)
		if err != nil {
			return nil, err
		}
		mt := time.Unix(mti, 0)
		size, err := strconv.Atoi(fileSize)
		if err != nil {
			return nil, err
		}
		item.FileModTime = mt
		item.FileSize = size

		result = append(result, idl.NewObjectInfoFromSF(item))
	}
	return result, nil
}

func (ld *LiteDB) insertFile(tableName string, tx *sql.Tx, srcItem *idl.SourceFile) error {
	q := `INSERT INTO %s(Name,ObjectID,VersionList,Checksum,Filename,FileModTime,Filesize) 
	VALUES (?,?,?,?,?,?,?);`
	q = fmt.Sprintf(q, tableName)
	if ld.DebugSQL {
		log.Println("Query is", q)
	}

	insertMore, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(insertMore).Exec(srcItem.Name, srcItem.ObjectID, srcItem.VersionList, srcItem.Checksum,
		srcItem.Filename, srcItem.FileModTime.Local().Unix(), srcItem.FileSize)
	if err != nil {
		return err
	}
	if ld.DebugSQL {
		log.Println("SourceFile added OK: ", srcItem.Name)
	}
	return nil
}
