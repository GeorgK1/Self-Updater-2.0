package main

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func checkError(err error) {

	if err != nil {
		panic(err)
	}

}

func runCommand(command string, argument string, path string) {
	args := strings.Split(argument, " ")

	cmd := exec.Command(command, args...)
	cmd.Dir = path

	stdout, err := cmd.StdoutPipe()
	checkError(err)
	err = cmd.Start()
	checkError(err)
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		t := scanner.Text()
		fmt.Println(t)
	}
}

func server(parsedPath string) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello, %q", r.URL.Path)
		if err != nil {
			return
		}

		if r.Method == "POST" {
			fmt.Println("POST")

			runCommand("git", "pull", parsedPath)

			fmt.Println("Restarting Docker containers")

			runCommand("docker-compose", "restart", parsedPath)

		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func main() {
	fmt.Println("Starting on localhost:8080")
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	checkError(err)

	auth := githttp.BasicAuth{
		Username: viper.GetString("GIT_USERNAME"),
		Password: viper.GetString("GIT_PASSWORD"),
	}
	url := viper.GetString("GIT_URL")

	parsedUrl := strings.Split(url, "/")

	path := filepath.Join("./", parsedUrl[4])

	parsedPath := strings.Split(path, ".git")[0]

	_, err = git.PlainClone(parsedPath, false, &git.CloneOptions{
		URL:      strings.TrimSpace(url),
		Progress: os.Stdout,
		Auth:     &auth,
	})

	if err != nil {
		fmt.Println("Repository already copied")

	}

	server(parsedPath)
}
