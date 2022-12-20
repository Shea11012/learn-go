package composite

import "fmt"

type folder struct {
	components []component
	name string
}

func (f *folder) search(keyword string) {
	fmt.Printf("search keyword %s in %s\n",keyword,f.name)
	for _,composite := range f.components {
		composite.search(keyword)
	}
}

func (f *folder) add(c component) {
	f.components = append(f.components,c)
}