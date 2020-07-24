//+build !windows

package gordon

import "os"

const executable = 0001

func isExecutable(fi os.FileInfo) bool {
	return (fi.Mode() & executable) == executable
}
