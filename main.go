package main

import (
	"fmt"
	"io"

	"net/http"
	"os"

	//"os/exec"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	git_http "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/viper"
)

func check_error(err error) {

	if err != nil {
		panic(err)
	}

}

func readFilesToLocalFileSystem(dir []os.FileInfo, parsedPath string, fs billy.Filesystem, ran bool, foldername string) {
	 
	for _, file := range dir {
		
		if file.IsDir() {
			//read a new directory using directory name'
			fmt.Println("Found a div")
			folderName := file.Name()
			fmt.Println(folderName)
			newpath := fmt.Sprintf("%s/%s", parsedPath, folderName)
			
			newfolder, err := fs.ReadDir(folderName)
			fmt.Println(newpath)
			check_error(err)
			
			os.Mkdir(newpath, 0755)
			
			check_error(err)
			readFilesToLocalFileSystem(newfolder, newpath, fs, true, folderName)
		} else {
			//create a new destination using the folder name and file name
			dst, err := os.Create(fmt.Sprintf("%s/%s", parsedPath, file.Name()))
			fmt.Printf("%s/%s", parsedPath, file.Name()+ "\n")
			check_error(err)
			
			
			if(ran) {
				f, err := fs.Open(fmt.Sprintf("%s/%s", foldername, file.Name()))
				check_error(err)
				io.Copy(dst, f)
			} else {
				f, err := fs.Open(file.Name())
				check_error(err)
				io.Copy(dst, f)
			}
			
			
		}

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

		fs := memfs.New()
		mr, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
			URL:      strings.TrimSpace(url),
			Progress: os.Stdout,
			Auth:     &auth,
		})

		if err != nil {
			fmt.Println("Repository already copied")
		}

		err = os.Mkdir(parsedPath, 0755)

		check_error(err)

		if r.Method == "POST" {
			fmt.Println("POST")

			mr.Fetch(&git.FetchOptions{
				RemoteName: viper.GetString("GIT_REMOTE"),
				Progress:   os.Stdout,
			})

			dir, _ := fs.ReadDir("./")

			//assemble files from the in memory filesystem and put them to the local filesystem
			//temppath := "./"
			readFilesToLocalFileSystem(dir, parsedPath, fs, false, "./")

			fmt.Println("Fetched")

			fmt.Println("Restarting Docker containers")

			//cmd := exec.Command("docker-compose", "up")
			//cmd.Dir = parsedPath

			//stdout, err := cmd.Output()

			//fmt.Println(string(stdout))
			check_error(err)

		}
	})

	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Starting on localhost:8080")

	server()
}
