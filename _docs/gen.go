// Copyright 2014 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build ignore

package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/internal"
)

const (
	url = "https://hajimehoshi.github.io/ebiten/"
)

var (
	examplesDir = filepath.Join("public", "examples")
)

func execute(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	msg, err := ioutil.ReadAll(stderr)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%v: %s", err, string(msg))
	}
	return nil
}

var (
	copyright     = ""
	stableVersion = ""
	devVersion    = ""
)

func init() {
	year, err := internal.LicenseYear()
	if err != nil {
		panic(err)
	}
	copyright = fmt.Sprintf("© %d Hajime Hoshi", year)
}

func init() {
	b, err := exec.Command("git", "tag").Output()
	if err != nil {
		panic(err)
	}
	lastStableVersion := ""
	lastCommitTime := 0
	for _, tag := range strings.Split(string(b), "\n") {
		m := regexp.MustCompile(`^v(\d.+)$`).FindStringSubmatch(tag)
		if m == nil {
			continue
		}
		t, err := exec.Command("git", "log", tag, "-1", "--format=%ct").Output()
		if err != nil {
			panic(err)
		}
		tt, err := strconv.Atoi(strings.TrimSpace(string(t)))
		if err != nil {
			panic(err)
		}
		if lastCommitTime >= tt {
			continue
		}
		lastCommitTime = tt
		lastStableVersion = m[1]
	}
	// See the HEAD commit time
	stableVersion = lastStableVersion
}

func init() {
	b, err := exec.Command("git", "show", "master:version.txt").Output()
	if err != nil {
		panic(err)
	}
	devVersion = strings.TrimSpace(string(b))
}

func comment(text string) template.HTML {
	// http://www.w3.org/TR/html-markup/syntax.html#comments
	// The text part of comments has the following restrictions:
	// * must not start with a ">" character
	// * must not start with the string "->"
	// * must not contain the string "--"
	// * must not end with a "-" character

	for strings.HasPrefix(text, ">") {
		text = text[1:]
	}
	for strings.HasPrefix(text, "->") {
		text = text[2:]
	}
	text = strings.Replace(text, "--", "", -1)
	for strings.HasSuffix(text, "-") {
		text = text[:len(text)-1]
	}
	return template.HTML("<!--\n" + text + "\n-->")
}

func safeHTML(text string) template.HTML {
	return template.HTML(text)
}

type example struct {
	Name         string
	ThumbWidth   int
	ThumbHeight  int
	ScreenWidth  int
	ScreenHeight int
}

func (e *example) Width() int {
	if e.ScreenWidth == 0 {
		return e.ThumbWidth * 2
	}
	return e.ScreenWidth
}

func (e *example) Height() int {
	if e.ScreenHeight == 0 {
		return e.ThumbHeight * 2
	}
	return e.ScreenHeight
}

func (e *example) Source() string {
	const (
		commentFor2048   = `// Please read examples/2048/main.go and examples/2048/2048/*.go`
		commentForBlocks = `// Please read examples/blocks/main.go and examples/blocks/blocks/*.go
// NOTE: If Gamepad API is available in your browswer, you can use gamepads. Try it out!`
	)
	if e.Name == "2048" {
		return commentFor2048
	}
	if e.Name == "blocks" {
		return commentForBlocks
	}

	path := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "hajimehoshi", "ebiten", "examples", e.Name, "main.go")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	str := regexp.MustCompile("(?s)^.*?\n\n").ReplaceAllString(string(b), "")
	str = strings.Replace(str, "\t", "        ", -1)
	return str
}

func versions() string {
	return fmt.Sprintf("v%s (dev: v%s)", stableVersion, devVersion)
}

var (
	gamesExamples = []example{
		{Name: "2048", ThumbWidth: 210, ThumbHeight: 300},
		{Name: "blocks", ThumbWidth: 256, ThumbHeight: 240},
	}
	graphicsExamples = []example{
		{Name: "alphablending", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "flood", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "font", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "highdpi", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "hsv", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "hue", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "infinitescroll", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "life", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "mandelbrot", ThumbWidth: 320, ThumbHeight: 320, ScreenWidth: 640, ScreenHeight: 640},
		{Name: "masking", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "mosaic", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "noise", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "paint", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "perspective", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "rotate", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "sprites", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "tiles", ThumbWidth: 240, ThumbHeight: 240},
	}
	inputExamples = []example{
		{Name: "gamepad", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "keyboard", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "typewriter", ThumbWidth: 320, ThumbHeight: 240},
	}
	audioExamples = []example{
		{Name: "audio", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "piano", ThumbWidth: 320, ThumbHeight: 240},
		{Name: "sinewave", ThumbWidth: 320, ThumbHeight: 240},
	}
)

func clear() error {
	if err := filepath.Walk("public", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if m, _ := regexp.MatchString("~$", path); m {
			return nil
		}
		// Remove auto-generated html files.
		m, err := regexp.MatchString(".html$", path)
		if err != nil {
			return err
		}
		if m {
			return os.Remove(path)
		}
		// Remove example resources that are copied.
		m, err = regexp.MatchString("^public/examples/_resources/images$", path)
		if err != nil {
			return err
		}
		if m {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
			return filepath.SkipDir
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func outputMain() error {
	f, err := os.Create("public/index.html")
	if err != nil {
		return err
	}
	defer f.Close()

	funcs := template.FuncMap{
		"comment":  comment,
		"safeHTML": safeHTML,
	}
	const templatePath = "index.tmpl.html"
	name := filepath.Base(templatePath)
	t, err := template.New(name).Funcs(funcs).ParseFiles(templatePath)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"URL":              url,
		"Copyright":        copyright,
		"StableVersion":    stableVersion,
		"DevVersion":       devVersion,
		"GraphicsExamples": graphicsExamples,
		"InputExamples":    inputExamples,
		"AudioExamples":    audioExamples,
		"GamesExamples":    gamesExamples,
	}
	return t.Funcs(funcs).Execute(f, data)
}

func createExamplesDir() error {
	if err := os.RemoveAll(examplesDir); err != nil {
		return err
	}
	if err := os.MkdirAll(examplesDir, 0755); err != nil {
		return err
	}
	return nil
}

func outputExampleResources() error {
	// TODO: Using cp command might not be portable.
	// Use io.Copy instead.
	if err := execute("cp", "-R", filepath.Join("..", "examples", "_resources"), filepath.Join(examplesDir, "_resources")); err != nil {
		return err
	}
	return nil
}

func outputExampleContent(e *example) error {
	f, err := os.Create(filepath.Join(examplesDir, e.Name+".content.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	funcs := template.FuncMap{
		"comment":  comment,
		"safeHTML": safeHTML,
	}
	const templatePath = "examplecontent.tmpl.html"
	name := filepath.Base(templatePath)
	t, err := template.New(name).Funcs(funcs).ParseFiles(templatePath)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"Copyright": copyright,
		"Example":   e,
	}
	if err := t.Funcs(funcs).Execute(f, data); err != nil {
		return err
	}

	out := filepath.Join(examplesDir, e.Name+".js")
	path := "github.com/hajimehoshi/ebiten/examples/" + e.Name
	if err := execute("gopherjs", "build", "--tags", "example", "-m", "-o", out, path); err != nil {
		return err
	}

	return nil
}

func outputExample(e *example) error {
	f, err := os.Create(filepath.Join(examplesDir, e.Name+".html"))
	if err != nil {
		return err
	}
	defer f.Close()

	funcs := template.FuncMap{
		"comment":  comment,
		"safeHTML": safeHTML,
	}
	const templatePath = "example.tmpl.html"
	name := filepath.Base(templatePath)
	t, err := template.New(name).Funcs(funcs).ParseFiles(templatePath)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"URL":       url,
		"Copyright": copyright,
		"Example":   e,
	}
	return t.Funcs(funcs).Execute(f, data)
}

func main() {
	if err := clear(); err != nil {
		log.Fatal(err)
	}
	if err := outputMain(); err != nil {
		log.Fatal(err)
	}
	if err := createExamplesDir(); err != nil {
		log.Fatal(err)
	}
	if err := outputExampleResources(); err != nil {
		log.Fatal(err)
	}

	examples := []example{}
	examples = append(examples, graphicsExamples...)
	examples = append(examples, inputExamples...)
	examples = append(examples, audioExamples...)
	examples = append(examples, gamesExamples...)
	for _, e := range examples {
		if err := outputExampleContent(&e); err != nil {
			log.Fatal(err)
		}
		if err := outputExample(&e); err != nil {
			log.Fatal(err)
		}
	}
}
