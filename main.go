package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/atotto/clipboard"
	"github.com/gabriel-vasile/mimetype"
	"github.com/jessevdk/go-flags"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

var usr = func() *user.User {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr
}()

// These variables are set in build step
var (
	Version = "unset"
)

// Option represents application options
type Option struct {
	Copy    bool `short:"c" long:"copy" description:"Copy to clipboard"`
	Init    bool `short:"i" long:"init" description:"Initialize temple config file"`
	Version bool `short:"v" long:"version" description:"Show temple version"`
}

// Config represents the settings for this application
type Config struct {
	HeadSize        int    `json:"headSize"`
	ItemSize        int    `json:"itemSize"`
	SyntaxHighlight string `json:"syntaxHighlight"`
}

// Template represents a template
type Template struct {
	Name string   `json:"name"`
	Path string   `json:"path"`
	Tags []string `json:"tags"`
}

// CLI represents this application itself
type CLI struct {
	Config    Config
	Templates []Template
}

func copy(srcPath, dstPath string) error {
	var err error

	srcPath, err = filepath.Abs(strings.Replace(srcPath, "~", usr.HomeDir, -1))
	if err != nil {
		return err
	}
	dstPath, err = filepath.Abs(strings.Replace(dstPath, "~", usr.HomeDir, -1))
	if err != nil {
		return err
	}

	if srcPath == dstPath {
		return nil
	}

	src, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	fmt.Printf("%s -> %s\n", srcPath, dstPath)
	return nil
}

func download(url string, path string) error {
	fmt.Printf("%s -> %s\n", url, path)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func isBinary(path string) bool {
	mine, err := mimetype.DetectFile(path)
	if err != nil {
		return true
	}
	isBinary := true
	for m := mine; m != nil; m = m.Parent() {
		if m.Is("text/plain") {
			isBinary = false
		}
	}
	return isBinary
}

// Prompt prompts CLI
func (c CLI) Prompt() (Template, error) {
	min := func(x, y int) int {
		if x < y {
			return x
		}
		return y
	}

	funcMap := promptui.FuncMap
	funcMap["join"] = func(texts []string) string {
		return strings.Join(texts, ", ")
	}
	funcMap["head"] = func(path string) string {
		p := strings.Replace(path, "~", usr.HomeDir, -1)
		if isBinary(p) {
			return "binary"
		}

		lexer := lexers.Match(p)
		if lexer == nil {
			lexer = lexers.Fallback
		}

		style := styles.Get(c.Config.SyntaxHighlight)
		if style == nil {
			style = styles.Fallback
		}

		formatter := formatters.Get("terminal256")
		if formatter == nil {
			formatter = formatters.Fallback
		}

		r, err := ioutil.ReadFile(p)
		if err != nil {
			return err.Error()
		}

		iterator, err := lexer.Tokenise(nil, string(r))
		if err != nil {
			return err.Error()
		}

		w := new(bytes.Buffer)
		if err := formatter.Format(w, style, iterator); err != nil {
			return err.Error()
		}

		s := w.String()
		ss := strings.Split(s, "\n")

		head := ""
		for i := 0; i < min(len(ss), c.Config.HeadSize); i++ {
			head += fmt.Sprintf("\n%4d %s", i+1, ss[i])
		}

		return head
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   promptui.IconSelect + " {{ .Path | cyan }}",
		Inactive: "  {{ .Path | faint }}",
		Selected: promptui.IconGood + " {{ .Path }}",
		Details: `
{{ "Name:" | faint }}	{{ .Name }}
{{ "Path:" | faint }}	{{ .Path }}
{{ "Tags:" | faint }}	{{ .Tags | join }}
{{ "Content:" | faint }}	{{ .Path | head }}
		`,
		FuncMap: funcMap,
	}

	searcher := func(input string, index int) bool {
		t := c.Templates[index]
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		name := strings.Replace(strings.ToLower(t.Name), " ", "", -1)
		path := strings.Replace(strings.ToLower(t.Path), " ", "", -1)
		tags := strings.Replace(strings.ToLower(strings.Join(c.Templates[index].Tags, " ")), " ", "", -1)

		if fuzzy.Match(input, name) || fuzzy.Match(input, path) || fuzzy.Match(input, tags) {
			return true
		}
		return false
	}

	sort.Slice(c.Templates, func(i, j int) bool {
		return c.Templates[i].Path < c.Templates[j].Path
	})

	prompt := promptui.Select{
		Label:             "Select",
		Items:             c.Templates,
		Size:              c.Config.ItemSize,
		HideSelected:      false,
		Templates:         templates,
		Searcher:          searcher,
		StartInSearchMode: true,
	}

	index, _, err := prompt.Run()
	return c.Templates[index], err
}

func run(args []string) int {
	var opt Option
	args, err := flags.ParseArgs(&opt, args)
	if err != nil {
		return 2
	}

	if opt.Init {
		p := path.Join(usr.HomeDir, ".config/temple.json")

		if _, err := os.Stat(p); !os.IsNotExist(err) {
			prompt := promptui.Prompt{
				Label:     "Overwrite",
				IsConfirm: true,
			}

			r, _ := prompt.Run()
			if strings.ToLower(r) != "y" {
				return 0
			}
		}

		if err = download("https://raw.githubusercontent.com/skmatz/temple/master/temple.json", p); err != nil {
			log.Print(err)
			return 1
		}
		return 0
	}

	if opt.Version {
		fmt.Printf("temple v%s\n", Version)
		return 0
	}

	viper.SetConfigName("temple")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.config")

	err = viper.ReadInConfig()
	if err != nil {
		log.Print(err)
		return 1
	}

	var cli CLI
	viper.Unmarshal(&cli)

	t, err := cli.Prompt()
	if err != nil {
		log.Print(err)
		return 1
	}

	if _, err := os.Stat(t.Name); !os.IsNotExist(err) {
		prompt := promptui.Prompt{
			Label:     "Overwrite",
			IsConfirm: true,
		}

		r, _ := prompt.Run()
		if strings.ToLower(r) != "y" {
			return 0
		}
	}

	if opt.Copy {
		p := strings.Replace(t.Path, "~", usr.HomeDir, -1)
		bytes, err := ioutil.ReadFile(p)
		if err != nil {
			log.Print(err)
			return 1
		}

		if err := clipboard.WriteAll(string(bytes)); err != nil {
			log.Print(err)
			return 1
		}

		fmt.Printf("%s -> clipboard\n", p)
		return 0
	}

	if err := copy(t.Path, t.Name); err != nil {
		log.Print(err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
