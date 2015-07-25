package tagopher

import (
	"fmt"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	os.FileInfo
}

func (f *File) Tags() (tags []string) {
	return
}

func (f *File) TagsLine() string {
	return strings.Join(f.Tags(), ", ")
}

func (f *File) String() string {
	return f.Name()
}

func Init() (err error) {
	path, _ := filepath.Abs(TAG_DIR)
	if _, err = os.Stat(TAG_DIR); os.IsNotExist(err) {
		os.Mkdir(TAG_DIR, os.FileMode(TAG_DEFAULT_DIR_PERM))
		fmt.Printf("Initialized empty %s repository in %s\n", TAG_NAME, path)
	} else {
		fmt.Printf("Reinitialized existing %s repository in %s\n", TAG_NAME, path)
	}
	db, err := bolt.Open(TAG_DATABASE, os.FileMode(TAG_DEFAULT_FILE_PERM), nil)
	if err != nil {
		return
	}
	db.Close()
	return
}

func AddTag(path string, name string) (err error) {
	return
}

func RemoveTag(path string, name string) (err error) {
	return
}

func RenameTag(path string, name string) (err error) {
	return
}

// Returns repository's root directory.
func Root() (path string, err error) {
	const MIN_ELEMENTS = 2
	var separator = fmt.Sprintf("%c", filepath.Separator)
	abspath, _ := filepath.Abs(".")
	elements := strings.Split(abspath, separator)
	for i := len(elements); i >= MIN_ELEMENTS; i-- {
		testpath := append(elements[:i], TAG_DIR)
		if testpath[0] == "" {
			testpath[0] = separator
		}
		tagdir := filepath.Join(testpath...)
		if _, err = os.Stat(tagdir); os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return
		}
		path = filepath.Join(testpath[:len(testpath)-1]...)
		break
	}
	if path == "" {
		err = fmt.Errorf("fatal: Not a %s repository (or any of the parent directories): %s", TAG_NAME, TAG_DIR)
	}
	return
}

func Get(path string) (file File, err error) {
	root, err := Root()
	if err != nil {
		return
	}
	abspath, _ := filepath.Abs(path)
	relpath, err := filepath.Rel(root, abspath)
	if err != nil {
		err = fmt.Errorf("fatal: Not a valid path under %s: %s", root, path)
	}
	if strings.Contains(relpath, "..") {
		err = fmt.Errorf("fatal: Not a valid path under %s: %s", root, path)
	}
	fi, err := os.Stat(abspath)
	if err != nil {
		return
	}
	file = File{fi}
	return
}

func List(path string) (files []File, err error) {
	if _, err = Root(); err != nil {
		return
	}
	files = make([]File, 0)
	fi, err := Get(path)
	if err != nil {
		return
	}
	if fi.IsDir() {
		entries, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, entry := range entries {
			files = append(files, File{entry})
		}
	} else {
		files = append(files, fi)
	}
	return
}
