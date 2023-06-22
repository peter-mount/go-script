package dataencoder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/util/walk"
	"github.com/peter-mount/go-script/tools/dataencoder/jenkinsfile"
	"github.com/peter-mount/go-script/tools/dataencoder/makefile"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Build struct {
	Encoder      *Encoder `kernel:"inject"`
	Dest         *string  `kernel:"flag,build,generate build files"`
	Platforms    *string  `kernel:"flag,build-platform,platform(s) to build"`
	PackageName  *string  `kernel:"flag,package,package name"`
	Dist         *string  `kernel:"flag,dist,distribution destination"`
	Prefix       *string  `kernel:"flag,prefix,Prefix to archive"`
	libProviders []LibProvider
}

// LibProvider handles calls to generate additional files/directories in a build
// returns destPath and arguments to pass
type LibProvider func(builds string) (string, []string)

func (s *Build) AddLibProvider(p LibProvider) {
	s.libProviders = append(s.libProviders, p)
}

func (s *Build) Run() error {
	if *s.Dest != "" {
		arch, err := s.getDist()
		if err != nil {
			return err
		}

		tools, err := s.getTools()
		if err != nil {
			return err
		}

		err = s.generate(tools, arch)
		if err != nil {
			return err
		}

		err = s.platformIndex(arch)
		if err != nil {
			return err
		}

		return s.jenkinsfile(arch)
	}
	return nil
}

func (s *Build) getDist() ([]Arch, error) {
	var buf bytes.Buffer
	cmd := exec.Command("go", "tool", "dist", "list", "-json")
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var arch []Arch
	if err := json.Unmarshal(buf.Bytes(), &arch); err != nil {
		return nil, err
	}

	sort.SliceStable(arch, func(i, j int) bool {
		a, b := arch[i], arch[j]

		if a.GOOS != b.GOOS {
			return a.GOOS < b.GOOS
		}

		return a.GOARCH < b.GOARCH
	})

	// Filter out blocked platforms
	var a []Arch
	for _, e := range arch {
		if !e.IsBlocked() {
			if e.GOARCH == "arm" {
				// We support arm 6 & 7 for 32bits
				e.GOARM = "6"
				a = append(a, e)

				e.GOARM = "7"
				a = append(a, e)
			} else {
				a = append(a, e)
			}
		}
	}
	return a, nil
}

func (s *Build) getTools() ([]string, error) {
	var tools []string

	if err := walk.NewPathWalker().
		Then(func(path string, info os.FileInfo) error {
			if info.Name() == "main.go" {
				tool := filepath.Base(filepath.Dir(filepath.Dir(path)))
				if tool != "dataencoder" {
					tools = append(tools, tool)
				}
			}
			return nil
		}).
		IsFile().
		Walk("tools"); err != nil {
		return nil, err
	}

	sort.SliceStable(tools, func(i, j int) bool {
		return tools[i] < tools[j]
	})

	return tools, nil
}

func (s *Build) generate(tools []string, arches []Arch) error {

	builder := makefile.New()
	builder.Comment("Generated Makefile %s", time.Now().Format(time.RFC3339)).
		Line("").
		Include("Makefile.include").
		Include("Go.include").
		Line("")

	root := s.allRule(arches, builder)

	targetGroups := s.targetGroups(arches, root)

	for _, arch := range arches {
		target := targetGroups.Get(arch.Target())

		for _, tool := range tools {
			s.goBuild(arch, target, tool)
		}

		for _, p := range s.libProviders {
			s.libProvider(arch, target, p)
		}

		s.tar(arch, target)
	}

	if err := os.MkdirAll(filepath.Dir(*s.Dest), 0755); err != nil {
		return err
	}

	return os.WriteFile(*s.Dest, []byte(builder.Build()), 0644)
}

func (s *Build) allRule(arches []Arch, builder makefile.Builder) makefile.Builder {
	all := builder.Rule("all")

	// Generate all target with either all or subset of platforms
	if *s.Platforms != "" {
		plats := strings.Split(*s.Platforms, " ")
		for _, arch := range arches {
			for _, plat := range plats {
				if strings.TrimSpace(plat) == arch.Platform() {
					all.AddDependency(arch.Target())
				}
			}
		}
	}

	// If all is still empty then return it so the Operating System rules
	// will get added to it automatically
	if all.IsEmptyRule() {
		return all
	}

	// All is not empty so return the original builder
	return builder
}

func (s *Build) targetGroups(arches []Arch, builder makefile.Builder) makefile.Map {
	osGroups := makefile.NewMap(builder)
	targetGroups := makefile.NewMap(builder)

	for _, arch := range arches {
		goos := arch.GOOS
		if !osGroups.Contains(goos) {
			osGroups.Add(goos, func(builder makefile.Builder) makefile.Builder {
				return builder.Block().
					Blank().
					Comment("==================").
					Comment(goos).
					Comment("==================").
					Rule(goos)
			})
		}

		target := arch.Target()
		if !targetGroups.Contains(target) {
			targetGroups.Add(target, func(_ makefile.Builder) makefile.Builder {
				return osGroups.Get(goos).
					Block().
					Blank().
					Comment("------------------").
					Comment("%s %s", arch.GOOS, arch.Arch()).
					Comment("------------------").
					Rule(target)
			})
		}
	}

	return targetGroups
}

// Build a tool in go
func (s *Build) goBuild(arch Arch, target makefile.Builder, tool string) {
	dest := arch.Tool(*s.Encoder.Dest, tool)

	target.Rule(dest).
		Line(`@echo "%-8s %s";\`, "GO-BUILD", dest).
		Line(
			"CGO_ENABLED=0 GOOS=%s GOARCH=%s GOARM=%s go build"+
				` -ldflags="-X '%s.Version=%s (%s %s %s) $(shell id -u -n) $(shell date))'"`+
				" -o %s %s",
			arch.GOOS,
			arch.GOARCH,
			arch.GOARM,
			"PACKAGE_PREFIX",
			filepath.Base(dest),
			"version", arch.GOOS, arch.Arch(),
			dest,
			filepath.Join("tools", tool, "bin/main.go"),
		)
}

// Add rules for a LibProvider
func (s *Build) libProvider(arch Arch, target makefile.Builder, f LibProvider) {
	dest, args := f(arch.BaseDir(*s.Encoder.Dest))
	target.Rule(dest).
		Line(`@echo "%-8s %s";\`,
			"GENERATE",
			strings.Join(strings.Split(dest, "/")[1:], " "),
		).
		Line("%s -d %s %s",
			filepath.Join(*s.Encoder.Dest, "dataencoder"),
			dest,
			strings.Join(args, " "),
		)
}

// Add rule for a tar distribution
func (s *Build) tar(arch Arch, target makefile.Builder) {
	archive := filepath.Join(*s.Dist, fmt.Sprintf("%s-%s_%s%s.tgz", *s.Prefix, arch.GOOS, arch.GOARCH, arch.GOARM))
	target.Rule(archive).
		Line("@mkdir -p %s", *s.Dist).
		Line(`@echo "%-8s %s";\`, "TAR", archive).
		Line(
			"tar -P --transform \"s|^%s|%s|\" -czpf %s %s",
			arch.BaseDir(*s.Encoder.Dest),
			*s.PackageName,
			archive,
			arch.BaseDir(*s.Encoder.Dest),
		)
}
func (s *Build) platformIndex(arches []Arch) error {
	var a []string
	a = append(a,
		"# Supported Platforms",
		"",
		"The following platforms are supported by virtue of how the build system works:",
		"",
		"| Operating System | CPU Architectures |",
		"| ---------------- | ----------------- |",
	)

	larch := ""
	for _, arch := range arches {
		if arch.GOOS != larch {
			larch = arch.GOOS

			var as []string
			as = append(as, "|", larch, "|")
			for _, arch2 := range arches {
				if arch2.GOOS == larch {
					as = append(as, arch2.GOARCH+arch2.GOARM)
				}
			}
			as = append(as, "|")
			a = append(a, strings.Join(as, " "))
		}
	}

	a = append(a, "")
	return os.WriteFile("platforms.md", []byte(strings.Join(a, "\n")), 0644)
}

func (s *Build) jenkinsfile(arches []Arch) error {

	builder := jenkinsfile.New()

	builder.Begin("properties([").
		Array().
		Begin("buildDiscarder(").
		Begin("logRotator(").
		Array().
		Property("artifactDaysToKeepStr", "").
		Property("artifactNumToKeepStr", "").
		Property("daysToKeepStr", "").
		Property("numToKeepStr", 10).
		End().End().
		Simple("disableConcurrentBuilds").
		Simple("disableResume").
		Begin("pipelineTriggers([").
		Simple("cron", `"H H * * *"`)

	node := builder.Node("go")

	node.Stage("Checkout").
		Line("checkout scm")

	node.Stage("Init").
		Sh("make clean init test")

	// Map of stages -> arch -> steps
	stages := make(map[string]*OsStage)
	for _, arch := range arches {
		stage := stages[arch.GOOS]
		if stage == nil {
			stage = &OsStage{
				arch:     arch,
				builder:  node.Stage(arch.GOOS).Parallel(),
				children: make(map[string]*ArchStage),
			}
		}
		stage1 := stage.children[arch.Arch()]
		if stage1 == nil {
			stage1 = &ArchStage{
				arch:    arch,
				builder: stage.builder.Stage(arch.Arch()),
			}
		}
		stage1.builder.Sh("make -f Makefile.gen " + arch.Target())
		stage.children[arch.Arch()] = stage1
		stages[arch.GOOS] = stage
	}

	// Sort stages
	for _, s1 := range stages {
		s1.builder.Sort()
		for _, s2 := range s1.children {
			s2.builder.Sort()
		}
	}

	return os.WriteFile("Jenkinsfile", []byte(builder.Build()), 0644)
}

type OsStage struct {
	arch     Arch
	builder  jenkinsfile.Builder
	children map[string]*ArchStage
}
type ArchStage struct {
	arch    Arch
	builder jenkinsfile.Builder
}
