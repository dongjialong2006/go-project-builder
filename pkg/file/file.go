package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func FileExist(name string) error {
	if "" == name {
		return fmt.Errorf("file name is empty.")
	}
	file, err := os.Open(name)
	if nil != err {
		if os.IsNotExist(err) {
			return fmt.Errorf("file name:%s is not exist.", name)
		}
		return err
	}
	file.Close()
	return nil
}

func PathExist(path string) ([]os.FileInfo, error) {
	if "" == path {
		return nil, fmt.Errorf("config dir is empty.")
	}

	return ioutil.ReadDir(path)
}

func FileExistInPath(dir string, suffix string) ([]string, error) {
	files, err := PathExist(dir)
	if nil != err {
		return nil, err
	}

	var paths []string = nil
	for _, file := range files {
		if nil == file || file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), suffix) {
			continue
		}
		paths = append(paths, path.Join(dir, file.Name()))
	}

	if 0 == len(paths) {
		return nil, fmt.Errorf("path:%s is not exist files with suffix:%s.", dir, suffix)
	}
	return paths, nil
}

func ReadFile(name string) ([]byte, error) {
	if "" == name {
		return nil, fmt.Errorf("name is empty.")
	}

	fp, err := os.Open(name)
	if nil != err {
		return nil, err
	}
	defer fp.Close()

	return ioutil.ReadAll(fp)
}

func CreatePath(path string) error {
	if "" == path {
		return fmt.Errorf("path is empty.")
	}

	_, err := os.Stat(path)
	if nil != err {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
		}
	}

	return err
}

func CreateFile(path string) (*os.File, error) {
	if "" == path {
		return nil, fmt.Errorf("path is empty.")
	}

	var dir string = "./"
	pos := strings.LastIndex(path, "/")
	if -1 != pos {
		dir = path[:pos]
	} else {
		path = fmt.Sprintf("./%s", path)
	}
	_, err := os.Stat(dir)
	if nil != err {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
		}
	}

	return os.Create(path)
}
