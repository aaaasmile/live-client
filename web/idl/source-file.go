package idl

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aaaasmile/live-client/util"
)

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

func (sf *SourceFile) CreateNewFile(dirTo, name string) error {
	initialVersion := "1"
	intName, id := sf.getInternalFilename(name, initialVersion)
	fsrc, err := sf.touchNewFile(dirTo, intName)
	if err != nil {
		return err
	}
	sf.ObjectID = id
	sf.Name = name
	sf.VersionList = initialVersion

	log.Println("Created an empty file ", intName, id, fsrc)
	return nil
}

func (sf *SourceFile) FillFromFname(fname string) error {
	var name, id, version string
	//fmt.Println("*** FillFromFname ", fname)
	arr := strings.Split(fname, "-")
	if len(arr) != 3 {
		return fmt.Errorf("File name structure unrecognized")
	}
	name = arr[0]
	id = arr[1]
	version = arr[2]
	//fmt.Println("*** FillFromFname res ", name, id, version)
	if version == "" {
		return fmt.Errorf("Version is empty")
	}
	sf.VersionList = version
	sf.Name = name
	if name == "" {
		return fmt.Errorf("Name is empty")
	}
	sf.ObjectID = id
	if id == "" {
		return fmt.Errorf("Id is empty")
	}
	return nil
}

func (sf *SourceFile) getInternalFilename(name string, version string) (string, string) {
	id := util.GenerateGUID2()
	res := fmt.Sprintf("%s-%s-%s", name, id, version)
	return res, id
}

func (sf *SourceFile) touchNewFile(destDir string, baseName string) (string, error) {
	fname := path.Join(destDir, baseName)
	_, err := os.Stat(fname)
	if os.IsNotExist(err) {
		file, err := os.Create(fname)
		if err != nil {
			return "", err
		}
		defer file.Close()
	}
	return fname, nil
}
