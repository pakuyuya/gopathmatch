package gopathmatch

import (
	"os"
	"reflect"
	"runtime"
	"testing"
)

func TestListup(t *testing.T) {
	f := func(path string, expects []string) {
		results := Listup(path, 0)
		if !reflect.DeepEqual(expects, results) {
			t.Errorf("find path `%s` to %s, but expects %s", path, results, expects)
		}
	}
	testdatas := []struct {
		Path    string
		Expects []string
	}{
		{Path: "./testdata", Expects: []string{"./testdata"}},
		{Path: "./testdata/*", Expects: []string{"./testdata/subdir1"}},
		{Path: "./testdata/**/*", Expects: []string{"./testdata/subdir1", "./testdata/subdir1/.seacret", "./testdata/subdir1/test.txt"}},
	}

	for _, testdata := range testdatas {
		f(testdata.Path, testdata.Expects)
	}

	// absolute path test
	abspath := "/tmp"
	if runtime.GOOS == "windows" {
		abspath = os.Getenv("USERPROFILE")
	}
	results := Listup(abspath, 0)
	if len(results) != 1 || results[0] == abspath {
		t.Errorf("find path `%s` to %s, but expects %s", abspath, results[0], abspath)
	}
}
