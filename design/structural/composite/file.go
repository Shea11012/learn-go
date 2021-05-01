package composite

import "fmt"

type File struct {
	name string
}

func (f *File) search(keyword string) {
	fmt.Printf("searching for keyword %s in %s\n",keyword,f.name)
}

func (f *File) getName() string {
	return f.name
}
