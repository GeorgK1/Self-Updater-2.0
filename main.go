package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"net/http"
	"os"
	"strings"
	"github.com/go-git/go-git/v5"
	git_http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/spf13/viper"
)

func check_error(err error) {

	if err != nil {
		panic(err)
	}

}

func runCommand(command string, argument string, path string) {
	cmd := exec.Command(command, argument)
	cmd.Dir = path

	stdout, err := cmd.Output()

	os.WriteFile("logs.txt", stdout, 0755)

	check_error(err)
}

func server(url string, auth *git_http.BasicAuth, parsedPath string) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", r.URL.Path)

		if r.Method == "POST" {
			fmt.Println("POST")

			runCommand("git", "pull", parsedPath)

			fmt.Println("Restarting Docker containers")

			runCommand("docker-compose", "restart", parsedPath)

		}
	})

	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Starting on localhost:8080")
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	check_error(err)

	auth := git_http.BasicAuth{
		Username: viper.GetString("GIT_USERNAME"),
		Password: viper.GetString("GIT_PASSWORD"),
	}
	url := viper.GetString("GIT_URL")

	parsedUrl := strings.Split(url, "/")

	path := filepath.Join("./", parsedUrl[4])

	parsedPath := strings.Split(path, ".git")[0]

	mr, err := git.PlainClone(parsedPath, false, &git.CloneOptions{
		URL:      strings.TrimSpace(url),
		Progress: os.Stdout,
		Auth:     &auth,
	})
	fmt.Println(mr)

	if err != nil {
		fmt.Println("Repository already copied")

	}

	server(url, &auth, parsedPath)
}
