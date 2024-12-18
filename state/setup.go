package state

import (
	"fmt"
	"github.com/peter-mount/go-script/errors"
	"github.com/peter-mount/go-script/packages"
	"github.com/peter-mount/go-script/script"
	"path"
	"strings"
)

func (s *state) setup() error {
	for _, i := range s.script.Import {
		for _, p := range i.Packages {
			if err := s.importPackage(p); err != nil {
				return err
			}
		}
	}

	for _, f := range s.script.FunDec {
		if err := s.declareFunction(f); err != nil {
			return err
		}
	}

	return nil
}

func (s *state) declareFunction(f *script.FuncDec) error {
	// If function name starts with _ then it's local so prefix with the filename
	// so that it's not accessible from outside its own file.
	//
	// The format of the actual name stored in the map is in a format that a
	// script cannot use for a function call, protecting the local method.
	//
	// This is only used here and in GetFunction().
	name := f.Name
	if strings.HasPrefix(f.Name, "_") {
		name = "!" + f.Pos.Filename + "!" + f.Name
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if e, exists := s.functions[name]; exists {
		return fmt.Errorf("%s function %q already defined at %s", f.Pos.String(), f.Name, e.Pos.String())
	}
	s.functions[name] = f
	return nil
}

func (s *state) importPackage(p *script.ImportPackage) error {
	pkg, exists := packages.Lookup(p.Name)
	if !exists {
		return errors.Errorf(p.Pos, "package %q is not available", p.Name)
	}

	if p.As == "" {
		p.As = path.Base(p.Name)
		if strings.ContainsAny(p.As, " /.") {
			return errors.Errorf(p.Pos, "package %q contains invalid characters", p.As)
		}
	}

	key := p.Pos.Filename + "!" + p.As
	if _, ok := s.packages[key]; ok {
		return errors.Errorf(p.Pos, "package %q %q has already been imported", p.As, p.Name)
	}

	s.packages[key] = pkg
	return nil
}
