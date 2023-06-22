package dataencoder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/util/walk"
	"github.com/peter-mount/go-script/tools/dataencoder/jenkinsfile"
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

func (b *Build) AddLibProvider(p LibProvider) {
	b.libProviders = append(b.libProviders, p)
}

func (b *Build) Run() error {
	if *b.Dest != "" {
		arch, err := b.getDist()
		if err != nil {
			return err
		}

		tools, err := b.getTools()
		if err != nil {
			return err
		}

		err = b.generate(tools, arch)
		if err != nil {
			return err
		}

		err = b.platformIndex(arch)
		if err != nil {
			return err
		}

		return b.jenkinsfile(arch)
	}
	return nil
}

func (b *Build) getDist() ([]Arch, error) {
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

func (b *Build) getTools() ([]string, error) {
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

	var a []string
	a = append(a,
		"# Generated Makefile "+time.Now().Format(time.RFC3339),
		"",
		"include Makefile.include",
		"include Go.include",
		"",
	)

	var archListTargets []string

	// Generate all target with either all or subset of platforms
	if *s.Platforms != "" {
		plats := strings.Split(*s.Platforms, " ")
		for _, arch := range arches {
			for _, plat := range plats {
				if strings.TrimSpace(plat) == arch.Platform() {
					archListTargets = append(archListTargets, arch.Target())
				}
			}
		}
	} else if len(archListTargets) == 0 {
		for _, arch := range arches {
			archListTargets = append(archListTargets, arch.Target())
		}
	}
	a = append(a, "all: "+strings.Join(archListTargets, " "), "")

	var archList, toolList []string
	libList := make(map[string][]string)

	los := ""
	var losdep []string
	for _, arch := range arches {
		if los != arch.GOOS {
			if len(losdep) > 0 {
				a = append(a, los+": "+strings.Join(losdep, " "), "")
			}
			los = arch.GOOS
			losdep = nil
		}
		losdep = append(losdep, arch.Target())
	}
	a = append(a, los+": "+strings.Join(losdep, " "), "")

	for _, arch := range arches {
		archList = append(archList,
			"",
			"# "+arch.Platform(),
		)

		archListTargets = nil
		for _, tool := range tools {
			archListTargets = append(archListTargets, arch.Tool(*s.Encoder.Dest, tool))
		}

		// Now rules for each tool
		for _, tool := range tools {
			dest := arch.Tool(*s.Encoder.Dest, tool)
			toolList = append(toolList,
				"",
				dest+":",
				fmt.Sprintf(
					"\t$(call GO-BUILD,%s,%s,%s)",
					arch.Platform(),
					dest,
					filepath.Join("tools", tool, "bin/main.go")),
			)
		}

		// Run LibProvider's
		localLib := make(map[string][]string)
		for _, p := range s.libProviders {
			s.build(arch, localLib, p)
		}

		// Add localLib to targets & global libList
		for k, v := range localLib {
			libList[k] = append(libList[k], v...)
			archListTargets = append(archListTargets, k)
		}

		// Tar/Zip
		archive := filepath.Join(*s.Dist, fmt.Sprintf("%s-%s_%s%s.tgz", *s.Prefix, arch.GOOS, arch.GOARCH, arch.GOARM))
		toolList = append(toolList,
			"",
			archive+":",
			"\t@mkdir -p "+*s.Dist,
			fmt.Sprintf(
				"\t$(call cmd,\"TAR\",%s);tar -P --transform \"s|^%s|%s|\" -czpf %s %s",
				archive,
				arch.BaseDir(*s.Encoder.Dest),
				*s.PackageName,
				archive,
				arch.BaseDir(*s.Encoder.Dest),
			),
		)
		archListTargets = append(archListTargets, archive)

		// Do archList last
		archList = append(archList, arch.Target()+": "+strings.Join(archListTargets, " "))
	}

	a = append(a, archList...)
	a = append(a, toolList...)

	var keys []string
	for k, _ := range libList {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool { return keys[i] < keys[j] })

	for _, k := range keys {
		a = append(a, "", k+":")
		a = append(a, libList[k]...)
	}

	if err := os.MkdirAll(filepath.Dir(*s.Dest), 0755); err != nil {
		return err
	}
	return os.WriteFile(*s.Dest, []byte(strings.Join(a, "\n")), 0644)
}

func (s *Build) build(arch Arch, libList map[string][]string, f LibProvider) {
	dest, args := f(arch.BaseDir(*s.Encoder.Dest))
	libList[dest] = append(libList[dest],
		fmt.Sprintf(
			"\t$(call cmd,\"GENERATE\",\"%s\");%s -d %s %s",
			strings.Join(strings.Split(dest, "/")[1:], " "),
			filepath.Join(*s.Encoder.Dest, "dataencoder"),
			dest,
			strings.Join(args, " "),
		),
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

func (b *Build) jenkinsfile(arches []Arch) error {

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
