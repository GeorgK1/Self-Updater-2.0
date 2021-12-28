package main

import (
	
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"github.com/go-git/go-git/v5"
	git_http "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func check_error(err error) {
	if err != nil {
		panic(err)
	}
}



func server(url string) {
	fmt.Println(os.Args[2])
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", r.URL.Path)	
		
		auth:=git_http.BasicAuth{
			Username: os.Args[2],
			Password: os.Args[3],
		}
		
		parsedUrl := strings.Split(url, "/")
		path := fmt.Sprintf("../%s", parsedUrl[4])
		
		fmt.Println(auth)

		// Clone the given repository to the given directory

		repo, err := git.PlainClone(strings.Split(path, ".git")[0], false , &git.CloneOptions{ 
			URL: url,
			Progress: os.Stdout,
			Auth: &auth,
		})
		check_error(err)
		
		if(r.Method == "GET") {
			fmt.Println("POST")		
			
			wt, err := repo.Worktree()
			
			check_error(err)
			
			wt.Pull(&git.PullOptions{
				RemoteName: "main",
				Progress: os.Stdout,
			})
		
			fmt.Println("Pulled")

			fmt.Println("Restarting Docker containers")

			cmd := exec.Command("docker-compose", "restart")
			stdout, err := cmd.Output()

			check_error(err)

			fmt.Println(string(stdout))
		}
	})

	

	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Starting")

	server(os.Args[1])
}