package main

import (
	"bytes"
	"code.google.com/p/goauth2/oauth"
	"errors"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"net/url"
	"runtime"
	"strings"
)

const VERSION = "0.1.0"

var copyCommands = []string{"pbcopy", "xclip", "xsel", "putclip"}
var pasteCommands = []string{"pbpaste", "xclip -o", "xsel -o", "getclip"}

var private bool
var showHelp bool
var showVersion bool
var paste bool
var copyURL bool
var openGist bool
var description string

func init() {
	flag.BoolVar(&private, "p", false, "private gist")
	flag.BoolVar(&showHelp, "h", false, "show help")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&paste, "P", false, "paste from the clipboard")
	flag.BoolVar(&copyURL, "c", false, "copy gist URL to the clipboard")
	flag.BoolVar(&openGist, "o", false, "open the gist in a browser")
	flag.StringVar(&description, "d", "", "description of the gist")
}

func main() {
	flag.Parse()

	if showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if showVersion {
		fmt.Printf("go-gist v%s\n", VERSION)
		os.Exit(0)
	}

	client := client()
	public := !private
	gist := &github.Gist{Public: &public, Files: gistFiles(), Description: &description}
	gist, _, err := client.Gists.Create(gist)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*gist.HTMLURL)
	if copyURL {
		copyToClipboard(*gist.HTMLURL)
	}
	if openGist {
		openURLInBrowser(*gist.HTMLURL)
	}
}

func client() *github.Client {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	}

	client := github.NewClient(t.Client())
	if githubURL := os.Getenv("GITHUB_URL"); len(githubURL) > 0 {
		if !strings.HasSuffix(githubURL, "/") {
			githubURL = githubURL + "/"
		}
		githubURL += "api/v3/"
		url, err := url.Parse(githubURL)
		if err != nil {
			log.Fatal(err)
		}
		client.BaseURL = url
	}
	return client
}

func gistFiles() map[github.GistFilename]github.GistFile {
	files := map[github.GistFilename]github.GistFile{}
	args := flag.Args()
	if paste {
		name := "-"
		content, err := pasteFromClipboard()
		if err != nil {
			log.Fatal(err)
		}
		file := github.GistFile{Filename: &name, Content: &content}
		files[github.GistFilename(name)] = file
	} else if len(args) == 0 {
		name := "-"
		bytes, _ := ioutil.ReadAll(os.Stdin)
		content := string(bytes)
		file := github.GistFile{Filename: &name, Content: &content}
		files[github.GistFilename(name)] = file
	} else {
		for _, arg := range args {
			bytes, err := ioutil.ReadFile(arg)
			if err != nil {
				log.Fatal(err)
			}
			content := string(bytes)
			basename := path.Base(arg)
			file := github.GistFile{Filename: &basename, Content: &content}
			files[github.GistFilename(arg)] = file
		}
	}

	return files
}

func copyToClipboard(content string) {
	for _, command := range copyCommands {
		var c *exec.Cmd
		c = exec.Command("which", command)
		if err := c.Run(); err != nil {
			continue
		}
		c = exec.Command(command)
		c.Stdin = strings.NewReader(content)
		if err := c.Run(); err != nil {
			log.Fatal("Failed to copy to the clipboard")
		}
		return
	}
	log.Fatal("Failed to copy to the clipboard")
}

func pasteFromClipboard() (string, error) {
	for _, command := range pasteCommands {
		var c *exec.Cmd
		c = exec.Command("which", command)
		if err := c.Run(); err != nil {
			continue
		}
		var stdoutBuffer bytes.Buffer
		c = exec.Command(command)
		c.Stdout = &stdoutBuffer
		if err := c.Run(); err != nil {
			return "", errors.New("Failed to paste from the clipboard. (pasting command failed)")
		}
		return stdoutBuffer.String(), nil
	}
	return "", errors.New("Failed to paste from the clipboard. (pasting command is not found)")
}

func openURLInBrowser(url string) {
	var command string
	if envBrowser := os.Getenv("BROWSER"); len(envBrowser) > 0 {
		command = envBrowser
	} else if runtime.GOOS == "darwin" {
		command = "open"
	} else {
		log.Fatal("Cannot open the gist in the browser")
	}
	c := exec.Command(command, url)
	if err := c.Run(); err != nil {
		log.Fatal("Failed to open the gist in the browser")
	}
}
