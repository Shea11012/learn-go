package composite

import "testing"

func TestComposite(t *testing.T) {
	file1 := &File{name: "file1"}
	file2 := &File{name: "file2"}
	file3 := &File{name: "file3"}
	file4 := &File{name: "file4"}

	folder1 := &folder{name: "folder1"}

	folder1.add(file1)

	folder2 := &folder{name: "folder2"}
	folder2.add(file2)
	folder2.add(file3)
	folder2.add(file4)
	folder2.add(folder1)

	folder2.search("rose")
}
