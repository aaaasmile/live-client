package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/aaaasmile/live-client/web/idl"
)

type ProjectTable struct {
	Fields       idl.Project
	DebugVerbose bool
	prjID        int
	connDb       *sql.DB
}

func NewProjectTable(prjID int, connDb *sql.DB, debug bool) (*ProjectTable, error) {
	res := ProjectTable{
		prjID:        prjID,
		DebugVerbose: debug,
		connDb:       connDb,
	}
	err := res.ReadFields()
	return &res, err
}

func (tb *ProjectTable) ReadFields() error {
	q := fmt.Sprintf(`SELECT Name,Description,SourceDir,SQLDatabase  FROM Project WHERE id = %d;`, tb.prjID)
	if tb.DebugVerbose {
		log.Println(q)
	}

	var name, description, sourceDir, sqlDatabase string
	err := tb.connDb.QueryRow(q).Scan(&name, &description, &sourceDir, &sqlDatabase)

	if err != nil {
		if err == sql.ErrNoRows {
			if tb.DebugVerbose {
				log.Println("No project source file found")
			}
			return err
		}
	} else {
		tb.Fields.DbLiteID = tb.prjID
		tb.Fields.Name = name
		tb.Fields.Description = description
		tb.Fields.SourceDir = sourceDir
		tb.Fields.SQLDatabase = sqlDatabase
	}
	//fmt.Println("*** project fields", tb.Fields)
	return err

}

func (tb *ProjectTable) Update() error {
	q := `UPDATE Project SET Name=?,Description=?,SourceDir=?,SQLDatabase=? WHERE id=?;`
	if tb.DebugVerbose {
		log.Println(q)
	}
	updateMore, err := tb.connDb.Prepare(q)
	if err != nil {
		return err
	}

	_, err = updateMore.Exec(tb.Fields.Name, tb.Fields.Description, tb.Fields.SourceDir, tb.Fields.SQLDatabase, tb.prjID)
	if err != nil {
		return err
	}
	if tb.DebugVerbose {
		log.Println("Project update ok")
	}
	return nil
}
