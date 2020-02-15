package main

import (
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

	"github.com/jessevdk/go-flags"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

// Option represents application options
type Option struct {
	Init bool `short:"i" long:"init" description:"Initialize temple config file"`
}

// Config represents the settings for this application
type Config struct {
	HeadSize int `json:"headSize"`
	ItemSize int `json:"itemSize"`
}

// Templates represents all templates
type Templates []struct {
	Path string   `json:"path"`
	Tags []string `json:"tags"`
}

// CLI represents this application itself
type CLI struct {
	Config    Config
	Templates Templates
}

func copy(srcPath, dstPath string) error {
	var err error

	usr, err := user.Current()
	if err != nil {
		return err
	}

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

// Prompt prompts CLI
func (c CLI) Prompt() (string, error) {
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
		usr, err := user.Current()
		if err != nil {
			return err.Error()
		}

		p := strings.Replace(path, "~", usr.HomeDir, -1)
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return err.Error()
		}

		s := string(b)
		ss := strings.Split(s, "\n")
		indent := fmt.Sprintf("\n%s", "  ")

		return indent + strings.Join(ss[:min(len(ss), c.Config.HeadSize)], indent)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   promptui.IconSelect + " {{ .Path | cyan }}",
		Inactive: "  {{ .Path | faint }}",
		Selected: promptui.IconGood + " {{ .Path }}",
		Details: `
{{ "Path:" | faint }}	{{ .Path }}
{{ "Tags:" | faint }}	{{ .Tags | join }}
{{ "Content:" | faint }}	{{ .Path | head }}
		`,
		FuncMap: funcMap,
	}

	searcher := func(input string, index int) bool {
		var path string
		path = c.Templates[index].Path
		path = strings.Replace(strings.ToLower(path), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(path, input)
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

	return c.Templates[index].Path, err
}

func run(args []string) int {
	var opt Option
	args, err := flags.ParseArgs(&opt, args)
	if err != nil {
		return 2
	}

	usr, err := user.Current()
	if err != nil {
		return 1
	}

	if opt.Init {
		if err = download("https://raw.githubusercontent.com/skmatz/temple/master/temple.json", path.Join(usr.HomeDir, ".config/temple.json")); err != nil {
			log.Print(err)
			return 1
		}
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

	p, err := cli.Prompt()
	if err != nil {
		log.Print(err)
		return 1
	}

	ps := strings.Split(p, "/")
	name := ps[len(ps)-1]
	copy(p, name)

	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
