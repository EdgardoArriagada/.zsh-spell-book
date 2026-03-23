package main

import (
	"fmt"

	"example.com/workspace/lib/git"
)

func main() {
	repo_name, err := git.GetRepoName()
	if err != nil {
		return
	}

	fmt.Println(repo_name)
}
