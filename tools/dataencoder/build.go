package dataencoder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/util/walk"
	"github.com/peter-mount/go-script/tools/dataencoder/jenkinsfile"
	"github.com/peter-mount/go-script/tools/dataencoder/makefile"
	"github.com/peter-mount/go-script/tools/dataencoder/meta"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type Build struct {
	Encoder      *Encoder `kernel:"inject"`
	Dest         *string  `kernel:"flag,build,generate build files"`
	Platforms    *string  `kernel:"flag,build-platform,platform(s) to build"`
	Dist         *string  `kernel:"flag,dist,distribution destination"`
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
		meta, err := meta.New()
		if err != nil {
			return err
		}

		arch, err := s.getDist()
		if err != nil {
			return err
		}

		tools, err := s.getTools()
		if err != nil {
			return err
		}

		err = s.generate(tools, arch, meta)
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

func (s *Build) generate(tools []string, arches []Arch, meta *meta.Meta) error {

	builder := makefile.New()
	builder.Comment("Generated Makefile %s", meta.Time).
		SetVar("BUILD", meta.ToolName).
		SetVar("export BUILD_VERSION", "%q", meta.Version).
		SetVar("export BUILD_TIME", "%q", meta.Time).
		SetVar("export BUILD_PACKAGE_NAME", "%q", meta.PackageName).
		SetVar("export BUILD_PACKAGE_PREFIX", "%q", meta.PackagePrefix).
		Phony("all", "clean", "init", "test")

	s.init(builder)
	s.clean(builder)
	s.test(builder)

	root := s.allRule(arches, builder)

	targetGroups := s.targetGroups(arches, root)

	for _, arch := range arches {
		target := targetGroups.Get(arch.Target())

		for _, tool := range tools {
			s.goBuild(arch, target, tool, meta)
		}

		for _, p := range s.libProviders {
			s.libProvider(arch, target, p, meta)
		}

		s.tar(arch, target, meta)
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
					Rule(goos, "init")
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
					Rule(target, "init")
			})
		}
	}

	return targetGroups
}

func (s *Build) init(builder makefile.Builder) {
	builder.Rule("init").
		Mkdir(*s.Encoder.Dest, *s.Dist)
}

func (s *Build) callBuilder(builder makefile.Builder, action, cmd string, args ...string) {
	builder.Line("@$(BUILD) -d %s -%s %s %s", *s.Encoder.Dest, action, cmd, strings.Join(args, " "))
}

func (s *Build) clean(builder makefile.Builder) {
	rule := builder.Rule("clean").
		RM(*s.Encoder.Dest, *s.Dist)
	s.callBuilder(rule, "go", "clean", "--", "-testcache")
}

func (s *Build) test(builder makefile.Builder) {
	out := filepath.Join(*s.Encoder.Dest, "go-text.txt")

	rule := builder.Rule("test", "init").
		Mkdir(filepath.Dir(out))

	s.callBuilder(rule, "go", "test")
}

// Build a tool in go
func (s *Build) goBuild(arch Arch, target makefile.Builder, tool string, meta *meta.Meta) {
	dest := arch.Tool(*s.Encoder.Dest, tool)

	rule := target.Rule(dest).
		Mkdir(filepath.Dir(dest))

	if arch.GOARM == "" {
		rule.Line("@$(BUILD) -d %s -go build %s %s %s", *s.Encoder.Dest, arch.GOOS, arch.GOARCH, tool)
	} else {
		rule.Line("@$(BUILD) -d %s -go build %s %s %s %s", *s.Encoder.Dest, arch.GOOS, arch.GOARCH, arch.GOARM, tool)
	}
}

// Add rules for a LibProvider
func (s *Build) libProvider(arch Arch, target makefile.Builder, f LibProvider, meta *meta.Meta) {
	dest, args := f(arch.BaseDir(*s.Encoder.Dest))
	target.Rule(dest).
		Echo("GENERATE", strings.Join(strings.Split(dest, "/")[1:], " ")).
		Line("$(BUILD) -d %s %s", dest, strings.Join(args, " "))
}

// Add rule for a tar distribution
func (s *Build) tar(arch Arch, target makefile.Builder, meta *meta.Meta) {
	archive := filepath.Join(
		*s.Dist,
		fmt.Sprintf("%s_%s_%s_%s%s.tgz", meta.PackageName, meta.Version, arch.GOOS, arch.GOARCH, arch.GOARM),
	)

	rule := target.Rule(archive).
		Mkdir(*s.Dist)

	s.callBuilder(rule, "tar", archive, arch.BaseDir(*s.Encoder.Dest))
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
