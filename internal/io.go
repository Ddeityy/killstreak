package internal

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/ddeityy/steamlocate-go"
)

// get user's steamID using tf2 appmanifest: LastUser

func GetUserSteamId() string {
	s := steamlocate.SteamDir{}
	s.Locate()
	file, err := os.ReadDir(path.Join(s.Path, "userdata"))
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("[U:1:%s]", file[0].Name())
}

func ParseDemo(demoPath string) string {
	parserPath, err := filepath.Abs("parse_demo")
	if err != nil {
		log.Println(err)
	}
	command := exec.Command(parserPath, demoPath)
	var out bytes.Buffer

	command.Stdout = &out
	err = command.Run()
	if err != nil {
		log.Println(err)
	}

	return out.String()
}
