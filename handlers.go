package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func handlerMain(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile(path.Join(getCurrentDir(), "views", "index.html"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(body))
}

func handlerDir(w http.ResponseWriter, r *http.Request) {
	p := len("/dir/")

	dir, err := url.QueryUnescape(r.URL.Path[p:])
	if err != nil {
		http.Error(w, fmt.Sprintf("File not found: %s", dir), 404)
		return
	}

	node, err := RootNode.GetNode(dir)
	if err != nil {
		http.Error(w, fmt.Sprintf("File not found: %s", dir), 404)
		return
	}

	dirJsonInfo, err := node.Json()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating json for: %s", dir), 500)
		return
	}

	fmt.Fprint(w, string(dirJsonInfo))
}

func handlerStream(w http.ResponseWriter, r *http.Request) {
	var (
		out io.Reader
		err error
	)

	p := len("/stream/")
	urlFilename, err := url.QueryUnescape(r.URL.Path[p:])
	if err != nil {
		log.Fatal(err)
	}
	filename := path.Join(Conf.MusicPath, urlFilename)

	if _, err = os.Stat(filename); os.IsNotExist(err) {
		http.Error(w, fmt.Sprintf("File not found: %s", urlFilename), 404)
		return
	}

	ext := filepath.Ext(filename)
	switch ext {
	case ".flac":
		c1 := exec.Command("flac", "-dc", filename)
		c2 := exec.Command("lame", "-V", "-2", "-")

		p1, err := c1.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		out, err = c2.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		c2.Stdin = p1

		if err = c2.Start(); err != nil {
			log.Fatal(err)
		}

		if err = c1.Start(); err != nil {
			log.Fatal(err)
		}

	case ".mp3":
		out, err = os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}

	default:
		http.Error(w, fmt.Sprintf("Invalid file: %s", urlFilename), 404)
		return
	}

	w.Header().Set("Content-type", "audio/mpeg")
	io.Copy(w, out)
}
