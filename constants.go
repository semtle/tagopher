package tagopher

import (
	"fmt"
	"path/filepath"
)

var (
	TAG                   = "tags"
	TAG_NAME              = "tagopher"
	TAG_DIR               = fmt.Sprintf(".%s", TAG_NAME)
	TAG_DATABASE          = filepath.Join(TAG_DIR, "tags.db")
	TAG_DEFAULT_DIR_PERM  = 0755
	TAG_DEFAULT_FILE_PERM = 0644
)
