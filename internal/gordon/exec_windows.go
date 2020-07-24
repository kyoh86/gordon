//+build windows

package gordon

import (
	"os"
	"strings"
)

func isExecutable(fi os.FileInfo) bool {
	return strings.HasSuffix(fi.Name(), ".exe")
}
