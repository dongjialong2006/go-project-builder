package builder

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"go-project-builder/pkg/file"
)

type Builder struct {
	Title    string
	DestPath string

	ReadMe      ReadMe
	Version     Version
	Example     Example
	Author      Author
	License     License
	MakeFile    MakeFile
	DockerFile  DockerFile
	JenkinsFile JenkinsFile

	Cmds     map[string]Cmd
	Docs     map[string]Doc
	Pkgs     map[string]Pkg
	Types    map[string]Type
	Configs  map[string]Config
	Catalogs map[string]Catalog
}

func NewBuilder() *Builder {
	return &Builder{
		Cmds:     make(map[string]Cmd),
		Docs:     make(map[string]Doc),
		Pkgs:     make(map[string]Pkg),
		Types:    make(map[string]Type),
		Configs:  make(map[string]Config),
		Catalogs: make(map[string]Catalog),
	}
}

func (b *Builder) Start(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	dir := path.Join(b.DestPath, b.Title)
	err := b.title(dir)
	if nil != err {
		return err
	}

	if err = b.readMe(dir); nil != err {
		return err
	}

	if err = b.license(dir); nil != err {
		return err
	}

	if err = b.authors(dir); nil != err {
		return err
	}

	if err = b.version(dir); nil != err {
		return err
	}

	if err = b.example(dir); nil != err {
		return err
	}

	if err = b.makeFile(dir); nil != err {
		return err
	}

	if err = b.dockerFile(dir); nil != err {
		return err
	}

	if err = b.jenkinsFile(dir); nil != err {
		return err
	}

	if err = b.cmds(dir); nil != err {
		return err
	}

	if err = b.docs(dir); nil != err {
		return err
	}

	if err = b.writePkgs(dir); nil != err {
		return err
	}

	if err = b.types(dir); nil != err {
		return err
	}

	if err = b.configs(dir); nil != err {
		return err
	}

	return b.catalogs(dir)
}

func (b *Builder) title(dir string) error {
	if "" == b.Title {
		return fmt.Errorf("project name is empty.")
	}

	if "" == b.DestPath {
		return fmt.Errorf("project path is empty.")
	}

	return file.CreatePath(dir)
}

func (b *Builder) authors(dir string) error {
	if 0 == len(b.Author.Authors) && "" == b.Author.Description && "" == b.Author.BuildTime {
		return nil
	}

	fp, err := file.CreateFile(path.Join(dir, "AUTHORS"))
	if nil != err {
		return err
	}
	defer fp.Close()

	if "" != b.Author.BuildTime {
		fp.Write([]byte(fmt.Sprintf("# BuildTime: %s\n", b.Author.BuildTime)))
	}
	if "" != b.Author.Description {
		fp.Write([]byte(fmt.Sprintf("# Description: %s\n", b.Author.Description)))
	}

	if len(b.Author.Authors) > 0 {
		fp.Write([]byte(fmt.Sprintf("# Authors:\n%s\n", strings.Join(b.Author.Authors, "\n"))))
	}

	return nil
}

func (b *Builder) readMe(dir string) error {
	if "" == b.ReadMe.Content && "" == b.ReadMe.Description {
		return nil
	}

	fp, err := file.CreateFile(path.Join(dir, "README.md"))
	if nil != err {
		return err
	}
	defer fp.Close()

	fp.Write([]byte(fmt.Sprintf("# Time: %s\n# Authors: %s\n", time.Now().Format("2006-01-02 15:04:05"), strings.Join(b.Author.Authors, " "))))
	if "" != b.ReadMe.Description {
		fp.Write([]byte(fmt.Sprintf("# Description: %s\n", b.ReadMe.Description)))
	}

	fp.Write([]byte(fmt.Sprintf("\n%s\n\n\n", b.Title)))

	if "" != b.ReadMe.Content {
		fp.Write([]byte(fmt.Sprintf("\n%s\n", b.ReadMe.Content)))
	} else {
		fp.Write([]byte("## Installation\n------------\n\n\n"))
		fp.Write([]byte("## Prerequisites\n------------\n\n\n"))
		fp.Write([]byte("## Constraints\n------------\n\n\n"))
		fp.Write([]byte("## Documentation\n------------\n\n\n"))
		fp.Write([]byte("## Performance\n------------\n\n\n"))
		fp.Write([]byte("## Status\n------------\n\n\n"))
		fp.Write([]byte("## Example\n------------\n\n\n"))
		fp.Write([]byte("## FAQ\n------------\n\n\n"))
	}

	return nil
}

func (b *Builder) version(dir string) error {
	if "" == b.Version.Content && "" == b.Version.Description {
		return nil
	}

	fp, err := file.CreateFile(path.Join(dir, "version", "version.go"))
	if nil != err {
		return err
	}
	defer fp.Close()

	fp.Write([]byte(fmt.Sprintf("// Time: %s\n// Authors: %s\n", time.Now().Format("2006-01-02 15:04:05"), strings.Join(b.Author.Authors, " "))))
	if "" != b.Version.Description {
		fp.Write([]byte(fmt.Sprintf("// Description: %s\n\npackage version\n", b.Version.Description)))
	}
	fp.Write([]byte("\npackage version\n"))
	if "" != b.Version.Content {
		fp.Write([]byte(fmt.Sprintf("\n%s\n", b.Version.Content)))
	}

	return nil
}

func (b *Builder) example(dir string) error {
	if "" == b.Example.Content && "" == b.Example.Description {
		return nil
	}

	fp, err := file.CreateFile(path.Join(dir, "example", "example.go"))
	if nil != err {
		return err
	}
	defer fp.Close()

	fp.Write([]byte(fmt.Sprintf("// Time: %s\n// Authors: %s\n", time.Now().Format("2006-01-02 15:04:05"), strings.Join(b.Author.Authors, " "))))
	if "" != b.Example.Description {
		fp.Write([]byte(fmt.Sprintf("// Description: %s\n", b.Example.Description)))
	}
	fp.Write([]byte("\npackage example\n"))
	if "" != b.Example.Content {
		fp.Write([]byte(fmt.Sprintf("\n%s\n", b.Example.Content)))
	}

	return nil
}

func (b *Builder) makeFile(dir string) error {
	if "" == b.MakeFile.Content && "" == b.MakeFile.Description {
		return nil
	}

	fp, err := file.CreateFile(path.Join(dir, "Makefile"))
	if nil != err {
		return err
	}
	defer fp.Close()

	fp.Write([]byte(fmt.Sprintf("# Time: %s\n# Authors: %s\n", time.Now().Format("2006-01-02 15:04:05"), strings.Join(b.Author.Authors, " "))))
	if "" != b.MakeFile.Description {
		fp.Write([]byte(fmt.Sprintf("# Description: %s\n", b.MakeFile.Description)))
	}

	if "" != b.MakeFile.Content {
		fp.Write([]byte(fmt.Sprintf("\n%s\n", b.MakeFile.Content)))
	} else {
		fp.Write([]byte(fmt.Sprintf("\nIMPORT_PATH := %s\n", b.Title)))
		fp.Write([]byte("\nexport GOPATH := $(CURDIR)/.GOPATH\n"))
		fp.Write([]byte("\nunexport GOBIN\n"))

		var cmds []string = nil
		for key, _ := range b.Cmds {
			cmds = append(cmds, key)
		}
		fp.Write([]byte(fmt.Sprintf("\nall: %s\n", strings.Join(cmds, " "))))

		for _, cmd := range cmds {
			fp.Write([]byte(fmt.Sprintf("%s: .GOPATH\n", cmd)))
			if 1 == len(cmds) {
				fp.Write([]byte(fmt.Sprintf("\tgo install -tags netgo $(IMPORT_PATH)/cmd\n")))
				continue
			}
			fp.Write([]byte(fmt.Sprintf("\tgo install -tags netgo $(IMPORT_PATH)/cmd/%s\n", cmd)))
		}

		fp.Write([]byte("\n.GOPATH:\n"))
		fp.Write([]byte("\trm -rf $(CURDIR)/.GOPATH\n\tmkdir -p $(CURDIR)/.GOPATH/src\n\tln -sf $(CURDIR) $(CURDIR)/.GOPATH/src/$(IMPORT_PATH)"))
		fp.Write([]byte("\n\tmkdir -p $(CURDIR)/bin\n\tln -sf $(CURDIR)/bin $(CURDIR)/.GOPATH/bin\n\ttouch $@"))

		fp.Write([]byte("\ninit:\n\tglide init\n"))
		fp.Write([]byte("\nclean:\n\trm -rf bin .GOPATH\n"))
		fp.Write([]byte("\nupdate:\n\tglide up -v\n"))
	}

	return nil
}

func (b *Builder) dockerFile(dir string) error {
	if "" == b.DockerFile.Content && "" == b.DockerFile.Description {
		return nil
	}

	fp, err := file.CreateFile(path.Join(dir, "Dockerfile"))
	if nil != err {
		return err
	}
	defer fp.Close()

	fp.Write([]byte(fmt.Sprintf("# Time: %s\n# Authors: %s\n", time.Now().Format("2006-01-02 15:04:05"), strings.Join(b.Author.Authors, " "))))
	if "" != b.DockerFile.Description {
		fp.Write([]byte(fmt.Sprintf("# Description: %s\n", b.DockerFile.Description)))
	}
	if "" != b.DockerFile.Content {
		fp.Write([]byte(fmt.Sprintf("\n%s\n", b.DockerFile.Content)))
	}

	return nil
}

func (b *Builder) jenkinsFile(dir string) error {
	if "" == b.JenkinsFile.Content && "" == b.JenkinsFile.Description {
		return nil
	}

	fp, err := file.CreateFile(path.Join(dir, "Jenkinsfile"))
	if nil != err {
		return err
	}
	defer fp.Close()

	fp.Write([]byte(fmt.Sprintf("# Time: %s\n# Authors: %s\n", time.Now().Format("2006-01-02 15:04:05"), strings.Join(b.Author.Authors, " "))))
	if "" != b.JenkinsFile.Description {
		fp.Write([]byte(fmt.Sprintf("# Description: %s\n", b.JenkinsFile.Description)))
	}
	if "" != b.JenkinsFile.Content {
		fp.Write([]byte(fmt.Sprintf("\n%s\n", b.JenkinsFile.Content)))
	}

	return nil
}

func (b *Builder) cmds(dir string) error {
	if 0 == len(b.Cmds) {
		return nil
	}

	var err error = nil
	var fp *os.File = nil
	for key, value := range b.Cmds {
		if "" == key {
			continue
		}

		if strings.HasSuffix(key, ".go") {
			key = strings.Trim(key, ".go")
		}
		if 1 == len(b.Cmds) {
			fp, err = file.CreateFile(path.Join(dir, "cmd", fmt.Sprintf("%s.go", key)))
		} else {
			fp, err = file.CreateFile(path.Join(dir, "cmd", key, fmt.Sprintf("%s.go", key)))
		}
		if nil != err {
			return err
		}
		defer fp.Close()

		fp.Write([]byte(fmt.Sprintf("// Time: %s\n// Authors: %s\n", time.Now().Format("2006-01-02 15:04:05"), strings.Join(b.Author.Authors, " "))))
		if "" != value.Description {
			fp.Write([]byte(fmt.Sprintf("// Description: %s\n", value.Description)))
		}

		fp.Write([]byte("\npackage main\n"))
		if "" != value.Content {
			fp.Write([]byte(fmt.Sprintf("\n%s\n", value.Description)))
		} else {
			fp.Write([]byte(fmt.Sprintf("\nfunc main() {\n")))
			fp.Write([]byte("\treturn\n}\n"))
		}
	}

	return nil
}

func (b *Builder) docs(dir string) error {
	if 0 == len(b.Docs) {
		return nil
	}

	for key, value := range b.Docs {
		if "" == key {
			continue
		}

		fp, err := file.CreateFile(path.Join(dir, "doc", key))
		if nil != err {
			return err
		}
		defer fp.Close()

		fp.Write([]byte(fmt.Sprintf("# Time: %s\n# Authors: %s\n", time.Now().Format("2006-01-02 15:04:05"), strings.Join(b.Author.Authors, " "))))
		if "" != value.Description {
			fp.Write([]byte(fmt.Sprintf("# Description: %s\n", value.Description)))
		}

		if "" != value.Content {
			fp.Write([]byte(fmt.Sprintf("\n%s\n", value.Description)))
		}
	}

	return nil
}

func (b *Builder) writePkgs(dir string) error {
	if 0 == len(b.Pkgs) {
		return nil
	}

	var err error = nil
	var fp *os.File = nil
	for key, value := range b.Pkgs {
		if "" == key {
			continue
		}

		if strings.HasSuffix(key, ".go") {
			key = strings.Trim(key, ".go")
		}

		if 1 == len(b.Pkgs) {
			fp, err = file.CreateFile(path.Join(dir, "pkg", "pkg.go"))
		} else {
			fp, err = file.CreateFile(path.Join(dir, "pkg", key, fmt.Sprintf("%s.go", key)))
		}
		if nil != err {
			return err
		}
		defer fp.Close()

		fp.Write([]byte(fmt.Sprintf("// Time: %s\n// Authors: %s\n", time.Now().Format("2006-01-02 15:04:05"), strings.Join(b.Author.Authors, " "))))
		if "" != value.Description {
			fp.Write([]byte(fmt.Sprintf("// Description: %s\n", value.Description)))
		}
		fp.Write([]byte(fmt.Sprintf("\npackage %s\n", key)))
		if "" != value.Content {
			fp.Write([]byte(fmt.Sprintf("\n%s\n", value.Description)))
		}
	}

	return nil
}

func (b *Builder) types(dir string) error {
	if 0 == len(b.Types) {
		return nil
	}

	for key, value := range b.Types {
		if "" == key {
			continue
		}

		if strings.HasSuffix(key, ".go") {
			key = strings.Trim(key, ".go")
		}
		fp, err := file.CreateFile(path.Join(dir, "types", fmt.Sprintf("%s.go", key)))
		if nil != err {
			return err
		}
		defer fp.Close()

		fp.Write([]byte(fmt.Sprintf("// Time: %s\n// Authors: %s\n", time.Now().Format("2006-01-02 15:04:05"), strings.Join(b.Author.Authors, " "))))
		if "" != value.Description {
			fp.Write([]byte(fmt.Sprintf("// Description: %s\n", value.Description)))
		}
		fp.Write([]byte(fmt.Sprintf("\npackage %s\n", key)))
		if "" != value.Content {
			fp.Write([]byte(fmt.Sprintf("\n%s\n", value.Description)))
		}
	}

	return nil
}

func (b *Builder) configs(dir string) error {
	if 0 == len(b.Configs) {
		return nil
	}

	for key, value := range b.Configs {
		if "" == key {
			continue
		}

		fp, err := file.CreateFile(path.Join(dir, "config", key))
		if nil != err {
			return err
		}
		defer fp.Close()

		fp.Write([]byte(fmt.Sprintf("# Time: %s\n# Authors: %s\n", time.Now().Format("2006-01-02 15:04:05"), strings.Join(b.Author.Authors, " "))))
		if "" != value.Description {
			fp.Write([]byte(fmt.Sprintf("# Description: %s\n", value.Description)))
		}

		if "" != value.Content {
			fp.Write([]byte(fmt.Sprintf("\n%s\n", value.Description)))
		}
	}

	return nil
}

func (b *Builder) catalogs(dir string) error {
	if 0 == len(b.Catalogs) {
		return nil
	}

	for key, value := range b.Catalogs {
		if "" == key {
			continue
		}

		if strings.HasSuffix(key, ".go") {
			key = strings.Trim(key, ".go")
		}
		fp, err := file.CreateFile(path.Join(dir, key, fmt.Sprintf("%s.go", key)))
		if nil != err {
			return err
		}
		defer fp.Close()

		fp.Write([]byte(fmt.Sprintf("// Time: %s\n// Authors: %s\n", time.Now().Format("2006-01-02 15:04:05"), strings.Join(b.Author.Authors, " "))))
		if "" != value.Description {
			fp.Write([]byte(fmt.Sprintf("// Description: %s\n", value.Description)))
		}
		fp.Write([]byte(fmt.Sprintf("\npackage %s\n", key)))
		if "" != value.Content {
			fp.Write([]byte(fmt.Sprintf("\n%s\n", value.Description)))
		}
	}

	return nil
}

func (b *Builder) license(dir string) error {
	if "" == b.License.Content && "" == b.License.Description {
		return nil
	}

	fp, err := file.CreateFile(path.Join(dir, "LICENSE"))
	if nil != err {
		return err
	}
	defer fp.Close()

	fp.Write([]byte(fmt.Sprintf("// Time: %s\n// Authors: %s\n", time.Now().Format("2006-01-02 15:04:05"), strings.Join(b.Author.Authors, " "))))
	if "" != b.License.Description {
		fp.Write([]byte(fmt.Sprintf("// Description: %s\n", b.License.Description)))
	}

	fp.Write([]byte(fmt.Sprintf("\nproject: %s\n", b.Title)))

	if "" != b.License.Content {
		fp.Write([]byte(fmt.Sprintf("\n%s\n", b.License.Content)))
	}

	return nil
}
