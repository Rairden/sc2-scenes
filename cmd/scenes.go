package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Rairden/sc2-scenes/scene"
	"github.com/google/go-cmp/cmp"
)

var (
	e, _      = os.Executable()
	currDir   = path.Dir(e)
	logDir    = filepath.Join(currDir, "log")
	logPath   = filepath.Join(logDir, "errors-scene.log")
	scenePath = filepath.Join(currDir, "scene.txt")
)

func main() {
	setup()
	prod()
	// debug()
}

func setup() {
	os.Mkdir(logDir, os.ModePerm)
	logs, _ := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(logs)
}

func prod() {
	scn := scene.Scene{}
	lastUI := scn.UI.ActiveScreens
	lastGame := scn.Game

	for {
		time.Sleep(100 * time.Millisecond)
		_ = scn.SetMenu()

		// check if the scene has changed
		if cmp.Equal(scn.UI.ActiveScreens, lastUI) && cmp.Equal(scn.Game, lastGame) {
			continue
		}

		writeData(scenePath, scn.Menu.String())
		lastUI = scn.UI.ActiveScreens
		lastGame = scn.Game
	}
}

func debug() {
	scn := scene.Scene{}
	lastUI := scn.UI.ActiveScreens
	lastGame := scn.Game
	now := time.Now()
	globalTime := now
	var sb strings.Builder

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	go func() {
		<-c
		fmt.Println(sb.String())
		os.Exit(21)
	}()

	for {
		time.Sleep(100 * time.Millisecond)
		now = time.Now()
		_ = scn.SetMenu()
		wait := time.Since(now).Round(time.Millisecond).String()

		// only print to stdout if the scene is different
		if cmp.Equal(scn.UI.ActiveScreens, lastUI) && cmp.Equal(scn.Game, lastGame) {
			continue
		}

		menu := fmt.Sprintf("#################  %s, %s (prev) #################\n", scn.Menu.String(), scn.Prev.String())
		sb.WriteString(menu)
		waitTime := fmt.Sprintf("%10s  %s\n", wait, time.Since(globalTime).Round(time.Millisecond).String())
		sb.WriteString(waitTime)
		sb.WriteString(scene.PrettyPrint(scn.UI.ActiveScreens) + "\n")
		sb.WriteString(scene.PrettyPrint(scn.Game) + "\n")
		writeData(scenePath, scn.Menu.String())

		lastUI = scn.UI.ActiveScreens
		lastGame = scn.Game
	}
}

func writeData(fullPath, data string) {
	file, err := os.Create(fullPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(data)
	file.Sync()
}
