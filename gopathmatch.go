package gopathmatch

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	FlgFileOnly = 0x01
)

func MatchPathes(path string, flg int) []string {
	if os.PathSeparator != '/' {
		path = strings.Replace(path, "\\", "/", -1)
	}
	if strings.Index(path, "./") == 0 {
		path = path[2:]
	}

	basepath := "."
	if filepath.IsAbs(path) || strings.Index(path, "/") == 0 {
		var pos int
		basepath, pos = _fetchToken(path)
		path = path[pos:]
	}

	return _findMatchPaths(basepath, path, flg)
}

func _findMatchPaths(basepath, path string, flg int) []string {
	rets := []string{}

	if path == "" {
		return make([]string, 0, 0)
	}

	t, i := _fetchToken(path)
	nextpath := path[i:]

	if t[len(t)-1:] == "/" {
		t = t[0 : len(t)-1]
	}

	if t == "**" {
		return _findMatchPathsRecursive(basepath, nextpath, flg)
	} else if t == "." || t == ".." {
		return _findMatchPaths(basepath+"/"+t, nextpath, flg)
	}

	regexptn := strings.Replace(regexp.QuoteMeta(t), "\\*", ".*", -1)
	r := regexp.MustCompile(regexptn)

	files, _ := ioutil.ReadDir(basepath)
	for _, file := range files {
		childfile := file.Name()
		if r.MatchString(childfile) {
			if nextpath == "" {
				if flg&FlgFileOnly > 0 && file.IsDir() {
					continue
				}
				rets = append(rets, basepath+"/"+file.Name())
			} else if file.IsDir() {
				rets = append(rets, _findMatchPaths(basepath+"/"+file.Name(), nextpath, flg)...)
			}

		}
	}
	return rets
}

func _findMatchPathsRecursive(basepath, path string, flg int) []string {
	if path == "" {
		return []string{}
	}
	if path == "**" {
		_, i := _fetchToken(path)
		return _findMatchPathsRecursive(basepath, path[i:], flg)
	}

	rets := []string{}

	t, i := _fetchToken(path)
	regexptn := strings.Replace(regexp.QuoteMeta(t), "\\*", ".*", -1)
	r := regexp.MustCompile(regexptn)

	files, _ := ioutil.ReadDir(basepath)
	for _, file := range files {
		childfile := file.Name()
		if r.MatchString(childfile) {
			if file.IsDir() {
				if flg&FlgFileOnly == 0 {
					rets = append(rets, basepath+"/"+file.Name())
				}
				rets = append(rets, _findMatchPaths(basepath+"/"+file.Name(), path[i:], flg)...)
			} else {
				rets = append(rets, basepath+"/"+file.Name())
			}
		}
		if file.IsDir() {
			rets = append(rets, _findMatchPathsRecursive(basepath+"/"+childfile, path, flg)...)
		}
	}

	return rets
}

func _fetchToken(path string) (string, int) {
	pos := _min(strings.Index(path, "/")+1, len(path))
	if pos <= 0 {
		pos = len(path)
	}
	return path[0:pos], pos
}

func _min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
