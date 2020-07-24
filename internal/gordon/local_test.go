package gordon

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestIsDir(t *testing.T) {
	tdir, err := ioutil.TempDir(os.TempDir(), "gordon-local-test")
	if err != nil {
		t.Fatalf("failed to create tempdir: %+v", err)
	}

	defer func() {
		if err := os.RemoveAll(tdir); err != nil {
			t.Fatalf("failed to cleanup tempdir: %+v", err)
		}
	}()

	tfile, err := ioutil.TempFile(tdir, "dummy")
	if err != nil {
		t.Fatalf("failed to create tempfile: %+v", err)
	}
	defer func() {
		if err := tfile.Close(); err != nil {
			t.Fatalf("failed to close tempfile: %+v", err)
		}
	}()

	for _, testcase := range []struct {
		path        string
		expectIsDir bool
		expectErr   error
	}{{
		tdir, true, nil,
	}, {
		tfile.Name(), false, nil,
	}, {
		tfile.Name() + ".notexist", false, nil,
	}} {
		actualIsDir, actualErr := isDir(testcase.path)
		if actualErr != testcase.expectErr {
			t.Errorf("expect err %s on the path %s, but actually %s", testcase.expectErr, testcase.path, actualErr)
		} else if actualIsDir != testcase.expectIsDir {
			t.Errorf("expect isDir %v on the path %s, but actually %v", testcase.expectIsDir, testcase.path, actualIsDir)
		}
	}
}

func TestIsDirWithChild(t *testing.T) {
	tdir1, err := ioutil.TempDir(os.TempDir(), "gordon-local-test")
	if err != nil {
		t.Fatalf("failed to create tempdir 1: %+v", err)
	}

	defer func() {
		if err := os.RemoveAll(tdir1); err != nil {
			t.Fatalf("failed to cleanup tempdir 1: %+v", err)
		}
	}()

	tdir2, err := ioutil.TempDir(os.TempDir(), "gordon-local-test")
	if err != nil {
		t.Fatalf("failed to create tempdir 2: %+v", err)
	}

	defer func() {
		if err := os.RemoveAll(tdir2); err != nil {
			t.Fatalf("failed to cleanup tempdir 2: %+v", err)
		}
	}()

	tfile, err := ioutil.TempFile(tdir1, "dummy")
	if err != nil {
		t.Fatalf("failed to create tempfile: %+v", err)
	}
	defer func() {
		if err := tfile.Close(); err != nil {
			t.Fatalf("failed to close tempfile: %+v", err)
		}
	}()

	for _, testcase := range []struct {
		path        string
		expectIsDir bool
		expectErr   error
	}{{
		tdir1, true, nil,
	}, {
		tdir2, false, nil,
	}, {
		tfile.Name(), false, nil,
	}, {
		tfile.Name() + ".notexist", false, nil,
	}} {
		actualIsDir, actualErr := isDirWithChild(testcase.path)
		if actualErr != testcase.expectErr {
			t.Errorf("expect err %s on the path %s, but actually %s", testcase.expectErr, testcase.path, actualErr)
		} else if actualIsDir != testcase.expectIsDir {
			t.Errorf("expect isDir %v on the path %s, but actually %v", testcase.expectIsDir, testcase.path, actualIsDir)
		}
	}
}
