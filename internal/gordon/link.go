package gordon

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/saracen/walker"
)

const executable = 0400

func Link(ev Env, version Version) error {
	bins := map[string]string{}
	mans := map[string]string{}
	// unlink old links
	if err := Unlink(ev, version.App); err != nil {
		return err
	}

	// link all executables and mans
	if err := walker.Walk(VersionPath(ev, version), func(path string, fi os.FileInfo) error {
		switch {
		case (fi.Mode() & executable) == executable:
			// executable file
			bins[path] = fi.Name()
		case strings.HasSuffix(path, ".1"):
			// man file
			mans[path] = fi.Name()
		}
		return nil
	}); err != nil {
		return err
	}

	for binPath, binName := range bins {
		if err := os.Link(binPath, filepath.Join(ev.Bin(), binName)); err != nil {
			return err
		}

		manPath := binPath + ".1"
		manName, ok := mans[manPath]
		if !ok {
			continue
		}

		if err := os.Link(manPath, filepath.Join(ev.Man(), manName)); err != nil {
			return err
		}
	}
	return nil
}
