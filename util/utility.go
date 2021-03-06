package util

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/kardianos/osext"
)

var rootPath string
var UseRelativeRoot = true

func GetEncryptedPasswordFromPrompt() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter the secret to encode")
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	buf := []byte(response)
	return base64.StdEncoding.EncodeToString(buf)
}

func GetUserLogFile(serviceName string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fn := filepath.Join(usr.HomeDir, "Appdata", "Local", "Temp", fmt.Sprintf("%s.log", serviceName))
	log.Println("Logfile in use is ", fn)
	return fn
}

func GetFullPath(relPath string) string {
	if UseRelativeRoot {
		return relPath
	}

	if rootPath == "" {
		var err error
		rootPath, err = osext.ExecutableFolder()
		if err != nil {
			log.Fatalf("ExecutableFolder failed: %v", err)
		}
		log.Println("Executable folder (rootdir) is ", rootPath)
		//rootPath, _ = filepath.Split(os.Args[0]) // os.Args[0] can be "faked". (https://github.com/kardianos/osext)
	}
	r := filepath.Join(rootPath, relPath)
	return r
}

func WrapJsonToString(res string) ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(res)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func GenerateGUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	guid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return guid
}

func GenerateGUID2() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	guid := fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return guid
}


