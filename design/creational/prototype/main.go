package main

import (
	"fmt"
	"strings"
)

type inode interface {
	print(string)
	clone() inode
}

type file struct {
	name string
}

func (f *file) print(indent string) {
	fmt.Println(indent + f.name)
}

func (f *file) clone() inode {
	return &file{name: f.name + "_clone"}
}

type folder struct {
	children []inode
	name string
}

func (f *folder) print(indent string) {
	fmt.Println(indent + f.name)
	for _,i := range f.children {
		i.print(strings.Repeat(indent,2))
	}
}

func (f *folder) clone() inode {
	cloneFolder := &folder{name: f.name + "_clone",children: make([]inode,len(f.children))}
	copy(cloneFolder.children,f.children)
	return cloneFolder
}

func main() {
	file1 := &file{name: "File1"}
	file2 := &file{name: "File2"}
	file3 := &file{name: "File3"}

	folder1 := &folder{name: "Folder1",children: []inode{file1}}

	folder2 := &folder{name: "Folder2",children: []inode{folder1,file2,file3}}

	fmt.Println("printing hierarchy for folder2")
	folder2.print("  ")

	cloneFolder := folder2.clone()
	fmt.Printf("%v\n",cloneFolder)
	fmt.Println("printing hierarchy for clone Folder")
	cloneFolder.print("  ")
}


