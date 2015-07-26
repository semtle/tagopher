package tagopher

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/dgryski/go-farm"
	"github.com/sony/sonyflake"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	os.FileInfo
}

func (f *File) Tags() (tags []string) {
	ref, _ := anchor(f.Name())
	if ref == "" {
		return
	}
	db, err := opendb(true)
	if err != nil {
		return
	}
	defer db.Close()
	db.View(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(TAG))
		if err != nil {
			return
		}
		current := make([]string, 0)
		value := bucket.Get([]byte(ref))
		if value != nil {
			dec := gob.NewDecoder(bytes.NewBuffer(value))
			err = dec.Decode(&current)
			if err != nil {
				return
			}
		}
		tags = current[:]
		return
	})
	return
}

func (f *File) String() string {
	filename := f.Name()
	end := strings.LastIndex(filename, "[")
	if end > -1 {
		extension := filepath.Ext(filename)
		filename = fmt.Sprintf("%s%s", filename[:end], extension)
	}
	return filename
}

func opendb(readonly bool) (db *bolt.DB, err error) {
	options := bolt.DefaultOptions
	options.ReadOnly = readonly
	db, err = bolt.Open(TAG_DATABASE, os.FileMode(TAG_DEFAULT_FILE_PERM), options)
	return
}

func Init() (err error) {
	path, _ := filepath.Abs(TAG_DIR)
	if _, err = os.Stat(TAG_DIR); os.IsNotExist(err) {
		os.Mkdir(TAG_DIR, os.FileMode(TAG_DEFAULT_DIR_PERM))
		fmt.Printf("Initialized empty %s repository in %s\n", TAG_NAME, path)
	} else {
		fmt.Printf("Reinitialized existing %s repository in %s\n", TAG_NAME, path)
	}
	db, err := opendb(false)
	if err != nil {
		return
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) (err error) {
		_, err = tx.CreateBucketIfNotExists([]byte(TAG))
		if err != nil {
			return
		}
		return
	})
	return
}

func applyTags(base []string, tags []string, remove bool) (result []string) {
	found := make(map[string]bool)
	for _, tag := range base {
		found[tag] = true
	}
	for _, tag := range tags {
		if remove {
			delete(found, tag)
		} else {
			found[tag] = true
		}
	}
	for tag := range found {
		result = append(result, tag)
	}
	return
}

func updateTags(db *bolt.DB, ref string, tags []string, remove bool) error {
	return db.Update(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(TAG))
		if err != nil {
			return
		}
		current := make([]string, 0)
		value := bucket.Get([]byte(ref))
		if value != nil {
			dec := gob.NewDecoder(bytes.NewBuffer(value))
			err = dec.Decode(&current)
			if err != nil {
				return
			}
		}
		current = applyTags(current, tags, remove)
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		err = enc.Encode(current)
		if err != nil {
			return
		}
		bucket.Put([]byte(ref), buffer.Bytes())
		return
	})

}

func AddTags(path string, tags ...string) (err error) {
	if _, err = Root(); err != nil {
		return
	}
	ref, path := anchor(path)
	if ref == "" {
		err = addAnchor(path, "")
		if err != nil {
			return
		}
	}
	db, err := opendb(false)
	if err != nil {
		return
	}
	defer db.Close()
	err = updateTags(db, ref, tags, false)
	if err != nil {
		return
	}
	return
}

func RemoveTags(path string, tags ...string) (err error) {
	if _, err = Root(); err != nil {
		return
	}
	ref, _ := anchor(path)
	if ref == "" {
		err = fmt.Errorf("No anchor found: %s", path)
		return
	}
	db, err := opendb(false)
	if err != nil {
		return
	}
	defer db.Close()
	err = updateTags(db, ref, tags, true)
	if err != nil {
		return
	}
	return
}

func RenameTags(path string, tags ...string) (err error) {
	if len(tags)%2 != 0 {
		err = fmt.Errorf("Odd argument count")
		return
	}
	if _, err = Root(); err != nil {
		return
	}
	ref, _ := anchor(path)
	if ref == "" {
		err = fmt.Errorf("No anchor found: %s", path)
		return
	}
	db, err := opendb(false)
	if err != nil {
		return
	}
	oldtags := make([]string, 0)
	newtags := make([]string, 0)
	for i := 0; i < len(tags); i++ {
		if i%2 == 0 {
			oldtags = append(oldtags, tags[i])
		} else {
			newtags = append(newtags, tags[i])
		}
	}
	err = updateTags(db, ref, oldtags, true)
	if err != nil {
		return
	}
	err = updateTags(db, ref, newtags, false)
	if err != nil {
		return
	}
	defer db.Close()
	return
}

func Anchor(path string) (value string, err error) {
	if _, err = Root(); err != nil {
		return
	}
	value, _ = anchor(path)
	return
}

func anchor(path string) (value string, newpath string) {
	path = filepath.Clean(path)
	newpath = path
        extension := filepath.Ext(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		end := strings.LastIndex(path, extension)
		if end == -1  {
			end = len(path)
		}
		pattern := fmt.Sprintf("%s*\\[*\\]%s", path[:end], extension)
		matches, _ := filepath.Glob(pattern)
		if matches == nil {
			return
		}
		if len(matches) == 0 {
			return
		}
		path = matches[0]
		newpath = path
	}
	base := filepath.Base(path)
	filename := base[:strings.LastIndex(base, extension)]
	if !strings.HasSuffix(filename, "]") {
		return
	}
	start := strings.LastIndex(filename, "[")
	value = filename[start:]
	value = value[1 : len(value)-1]
	return
}

func AddAnchor(path string, value string) (err error) {
	if _, err = Root(); err != nil {
		return
	}
	current, path := anchor(path)
	if current != "" {
		err = fmt.Errorf("File anchor already set: %s", current)
		return
	}
	err = addAnchor(path, value)
	return
}

func generateAnchor() string {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	newref, _ := flake.NextID()
	return fmt.Sprintf("%x", farm.Hash32([]byte(fmt.Sprintf("%d", newref))))
}

func addAnchor(path string, value string) (err error) {
	path = filepath.Clean(path)
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return
	}
	if value == "" {
		value = generateAnchor()
	}
	base := filepath.Base(path)
	extension := filepath.Ext(base)
	filename := base[:strings.LastIndex(base, extension)]
	dir := filepath.Dir(path)
	err = os.Rename(path, filepath.Join(dir, fmt.Sprintf("%s[%s]%s", filename, value, extension)))
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
		err = fmt.Errorf("Not a %s repository (or any of the parent directories): %s", TAG_NAME, TAG_DIR)
	}
	return
}

func Get(path string) (file File, err error) {
	root, err := Root()
	if err != nil {
		return
	}
	file, err = get(path, root)
	return
}

func get(path string, root string) (file File, err error) {
	abspath, _ := filepath.Abs(path)
	relpath, err := filepath.Rel(root, abspath)
	if err != nil {
		err = fmt.Errorf("Not a valid path under %s: %s", root, path)
	}
	if strings.Contains(relpath, "..") {
		err = fmt.Errorf("Not a valid path under %s: %s", root, path)
	}
	fi, err := os.Stat(abspath)
	if err != nil {
		return
	}
	file = File{fi}
	return
}

func List(path string) (files []File, err error) {
	root, err := Root()
	if err != nil {
		return
	}
	files = make([]File, 0)
	fi, err := get(path, root)
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
