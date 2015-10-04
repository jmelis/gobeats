package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Name string
	Path string
}

type Node struct {
	Name string
	Path string

	Children []*Node
	Files    []File
}

func (node *Node) Traverse(path string, allFiles *[]string, info os.FileInfo) error {
	relPath := path[len(Conf.MusicPath):]

	node.Path = relPath

	if info != nil {
		node.Name = info.Name()
	} else {
		node.Name = "."
	}

	dirInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, fileInfo := range dirInfo {
		newPath := filepath.Join(path, fileInfo.Name())
		newRelPath := filepath.Join(relPath, fileInfo.Name())

		if fileInfo.IsDir() {
			child := &Node{}
			if err := child.Traverse(newPath, allFiles, fileInfo); err != nil {
				return err
			}
			node.Children = append(node.Children, child)
		} else {
			if validExtension(newPath) {
				file := File{fileInfo.Name(), newRelPath}
				*allFiles = append(*allFiles, newRelPath)
				node.Files = append(node.Files, file)
			}
		}
	}

	return nil
}

func (node *Node) GetChild(name string) (*Node, error) {
	for _, child := range node.Children {
		if child.Name == name {
			return child, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Node '%s' has no child '%s'", node.Name, name))
}

func (node *Node) GetNode(path string) (*Node, error) {
	pathComponents := strings.Split(path, "/")
	n := node
	for _, p := range pathComponents {
		if p != "" {
			c, err := n.GetChild(p)
			if err != nil {
				return nil, err
			} else {
				n = c
			}
		}
	}
	return n, nil
}

func (node *Node) Json() ([]byte, error) {
	type DirJson struct {
		Dirs  []string
		Files []string
	}

	dirJson := DirJson{}

	for _, f := range node.Files {
		dirJson.Files = append(dirJson.Files, f.Path)
	}

	for _, c := range node.Children {
		dirJson.Dirs = append(dirJson.Dirs, c.Path)
	}

	j, err := json.Marshal(dirJson)
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}
