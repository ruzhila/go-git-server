package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var gitRootPath string
var prefix string

func system(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func HandleGit(w http.ResponseWriter, r *http.Request) {
	pathInfo := strings.TrimPrefix(r.URL.Path, strings.TrimSuffix(prefix, "/"))
	handler := cgi.Handler{
		Path: "/usr/bin/git",
		Root: gitRootPath,
		Env: []string{
			"GIT_PROJECT_ROOT=" + gitRootPath,
			"GIT_HTTP_EXPORT_ALL=",
			"PATH_INFO=" + pathInfo,
			"QUERY_STRING=" + r.URL.RawQuery,
			"REQUEST_METHOD=" + r.Method,
			"CONTENT_TYPE=" + r.Header.Get("Content-Type"),
		},
		Args:       []string{"http-backend"},
		InheritEnv: []string{"PATH"},
		Stderr:     os.Stderr,
	}
	defer log.Println(r.RemoteAddr, r.Method, r.URL.Path)
	handler.ServeHTTP(w, r)
}

func CreateRepository(root, name string) error {
	reposDir := filepath.Join(root, name)
	err := os.MkdirAll(reposDir, 0755)
	if err != nil {
		return err
	}
	cmds := [][]string{
		{"git", "init", reposDir, "--bare", "--shared"},
		{"git", "--git-dir", reposDir, "update-server-info"},
		{"git", "--git-dir", reposDir, "config", "http.receivepack", "true"},
		{"git", "--git-dir", reposDir, "config", "init.defaultBranch", "master"},
	}
	for _, cmd := range cmds {
		if err = system(cmd[0], cmd[1:]...); err != nil {
			log.Println("Fail:", strings.Join(cmd, " "), err)
			return err
		}
	}
	return nil
}

func main() {
	var addr string
	var create bool = false
	flag.StringVar(&addr, "addr", "127.0.0.1:8080", "address to listen on")
	flag.StringVar(&gitRootPath, "root", "/tmp/git", "root repository path")
	flag.StringVar(&prefix, "prefix", "/", "prefix path for git server")
	flag.BoolVar(&create, "create", false, "create repository only without serving")
	flag.Parse()

	if create {
		if flag.NArg() == 0 {
			panic("missing repository name")
		}
		reposName := flag.Arg(0)
		err := CreateRepository(gitRootPath, reposName)
		if err != nil {
			panic(err)
		}
		return
	}
	fullAddr := filepath.Join(addr, prefix)
	welcome := `Simple Git Server started on %s
Quit the server with CONTROL-C.

Clone repository with:
  git clone http://%s/<repository>.git
`
	fmt.Printf(welcome, fullAddr, fullAddr)
	err := http.ListenAndServe(addr, http.HandlerFunc(HandleGit))
	if err != nil {
		panic(err)
	}
}
