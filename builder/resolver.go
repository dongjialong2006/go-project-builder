package builder

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"go-project-builder/pkg/file"

	"github.com/BurntSushi/toml"
)

type Resolver struct {
	ctx      context.Context
	cancel   context.CancelFunc
	paths    []string
	builders []*Builder
}

func NewResolver(ctx context.Context, path string) (*Resolver, error) {
	if "" == path {
		return nil, fmt.Errorf("config file path is empty.")
	}
	path = strings.Trim(path, " ")

	var paths []string = nil
	if strings.HasSuffix(path, ".toml") {
		if err := file.FileExist(path); nil != err {
			return nil, err
		}
		paths = append(paths, path)
	} else {
		tmp, err := file.FileExistInPath(path, ".toml")
		if nil != err {
			return nil, err
		}
		paths = append(paths, tmp...)
	}

	ctx, cancel := context.WithCancel(ctx)

	return &Resolver{
		ctx:    ctx,
		paths:  paths,
		cancel: cancel,
	}, nil
}

func (r *Resolver) init() error {
	for _, path := range r.paths {
		if "" == path {
			continue
		}

		builder := NewBuilder()
		_, err := toml.DecodeFile(path, builder)
		if nil != err {
			return err
		}
		r.builders = append(r.builders, builder)
	}

	if 0 == len(r.builders) {
		return fmt.Errorf("invalide toml config file.")
	}

	return nil
}

func (r *Resolver) Start() error {
	err := r.init()
	if nil != err {
		return err
	}

	var wg sync.WaitGroup
	for _, builder := range r.builders {
		if nil == builder {
			continue
		}
		wg.Add(1)
		go builder.Start(r.ctx, &wg)
	}
	wg.Wait()
	return nil
}

func (r *Resolver) Stop() {
	r.cancel()
	return
}
