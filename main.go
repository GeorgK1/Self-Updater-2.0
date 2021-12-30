package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
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

func server() {
	viper.SetConfigFile(".env")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", r.URL.Path)
		
		
		err := viper.ReadInConfig()

		check_error(err)

		auth := git_http.BasicAuth{
			Username: viper.GetString("GIT_USERNAME"),
			Password: viper.GetString("GIT_PASSWORD"),
		}
		url := viper.GetString("GIT_URL")

		parsedUrl := strings.Split(url, "/")
		path := fmt.Sprintf("./%s", parsedUrl[4])

		parsedPath := strings.Split(path, ".git")[0]

		// Clone the given repository to the given directory

		repo, err := git.PlainClone(parsedPath, false, &git.CloneOptions{
			URL:      strings.TrimSpace(url),
			Progress: os.Stdout,
			Auth:     &auth,
		})

		if(err != nil){
			fmt.Println("Repository already copied")
		}

		if r.Method == "POST" {
			fmt.Println("POST")

			wt, err := repo.Worktree()

			check_error(err)

			wt.Pull(&git.PullOptions{
				RemoteName: viper.GetString("GIT_REMOTE"),
				Progress:   os.Stdout,
			})

			fmt.Println("Pulled")

			fmt.Println("Restarting Docker containers")

			cmd := exec.Command("docker-compose", "up")
			cmd.Dir = parsedPath
		
			stdout, err := cmd.Output()

			fmt.Println(string(stdout))
			check_error(err)

		}
	})

	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Starting on localhost:8080")
	

	server()
}
