package tagopher

import (
	"fmt"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type File struct {
	os.FileInfo
}

func (f *File) Tags() (tags []string) {
	return
}

func (f *File) String() string {
	return ""
}

func Init() (err error) {
	path, _ := filepath.Abs(TAG_DIR)
	if _, err = os.Stat(TAG_DIR); os.IsNotExist(err) {
		os.Mkdir(TAG_DIR, os.FileMode(TAG_DEFAULT_DIR_PERM))
		fmt.Printf("Initialized empty %s repository in %s\n", TAG_NAME, path)
	} else {
		fmt.Printf("Reinitialized existing %s repository in %s\n", TAG_NAME, path)
	}
	db, err := bolt.Open("tags.db", os.FileMode(TAG_DEFAULT_FILE_PERM), nil)
	if err != nil {
		return
	}
	db.Close()
	return
}

func AddTag(path string, name string) (err error) {
	return
}

func RemoteTag(path string, name string) (err error) {
	return
}

func RenameTag(path string, name string) (err error) {
	return
}

func Get(path string) (file File, err error) {
	return
}

func List(path string) (files []File, err error) {
	files = make([]File, 0)
	fi, err := os.Stat(path)
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
		files = append(files, File{fi})
	}
	return
}
