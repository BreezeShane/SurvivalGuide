package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sv/email"
	"sync"
)

type Assignee struct {
	Id    int    `json:"id"`
	Name  string `json:"login"`
	Email string `json:"email"`
}

type User struct {
	Assignee
}

type Issue struct {
	Title     string     `json:"title"`
	Assignees []Assignee `json:"assignees"`
	Body      string     `json:"body"`
	User      User       `json:"user"`
}

func main() {
	flag.Parse()

	request, _ := http.NewRequest("GET", "https://api.github.com/repos/sanyuankexie/SurvivalGuide/issues/5", nil)
	request.Header.Add("Authorization", "token "+*email.GithubToken)
	respose, err := email.Client.Do(request)
	if err != nil {
		panic(err.Error())
	}
	defer respose.Body.Close()

	body, err := ioutil.ReadAll(respose.Body)
	if err != nil {
		panic(err.Error())
	}
	jsonData := string(body)
	issue := &Issue{}
	json.Unmarshal([]byte(jsonData), &issue)

	var wg sync.WaitGroup
	var toers = make([]string, 0, 5)
	for _, assignee := range issue.Assignees {
		wg.Add(1)

		go func(assignee Assignee) {
			defer wg.Add(-1)

			var email, err = email.GetEmail(assignee.Name)
			if err != nil {
				log.Println(err.Error())
			}

			if email != "" {
				log.Printf("found [%s] email:[%s]\n", assignee.Name, email)
				toers = append(toers, email)
			}
		}(assignee)
	}
	wg.Wait()
	email.Send(issue.Title, fmt.Sprintf("@%s: \n%s", issue.User.Name, issue.Body), toers)
}
