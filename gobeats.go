package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"runtime"

	// "fmt"
	// "os"
)

var (
	RootNode   *Node         = &Node{}
	Conf       Configuration = Configuration{}
	MusicFiles []string
)

const (
	confPath = "gobeats.conf"
)

type Configuration struct {
	MusicPath  string
	Extensions []string
}

func validExtension(path string) bool {
	ext := filepath.Ext(path)

	for _, cExt := range Conf.Extensions {
		if ext == cExt {
			return true
		}
	}

	return false
}

func getCurrentDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func readConf() error {
	data, err := ioutil.ReadFile(path.Join(getCurrentDir(), confPath))
	if err != nil {
		return err
	}

	// TODO: error if required config fields are unset
	// TODO: add '/' to the end of MusicPath if not present

	return json.Unmarshal(data, &Conf)
}

func main() {

	err := readConf()
	if err != nil {
		log.Fatal(err)
	}

	if err := RootNode.Traverse(Conf.MusicPath, &MusicFiles, nil); err != nil {
		log.Fatal("trave: ", err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/stream/", handlerStream)
	http.HandleFunc("/dir/", handlerDir)
	http.HandleFunc("/", handlerMain)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
