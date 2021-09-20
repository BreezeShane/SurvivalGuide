package email

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Event struct {
	Type    string       `json:"type"`
	Payload EventPayload `json:"payload"`
}

type EventPayload struct {
	Commits []Commit `json:"commits"`
}

type Commit struct {
	Author CommitAuthor `json:""`
}

type CommitAuthor struct {
	Email string `json:"email"`
}

var staticData map[string]string
var Client *http.Client
var GithubToken = flag.String("githubToken", "", "github token")

func init() {
	staticData = map[string]string{
		"xiaoguokf":   "panzguo@qq.com",
		"ZLzzzzzzz":   "2290885422@qq.com",
		"Lei1900":     "2674312206@qq.com",
		"Volerde":     "volerde@qq.com",
		"Onlytonight": "1198076162@qq.com",
	}
	Client = &http.Client{}
}

func GetEmail(name string) (string, error) {
	request, _ := http.NewRequest("GET", "https://api.github.com/users/"+name+"/events/public", nil)
	request.Header.Add("Authorization", "token "+*GithubToken)

	respose, err := Client.Do(request)
	if err != nil {
		panic(err.Error())
	}
	defer respose.Body.Close()

	body, err := ioutil.ReadAll(respose.Body)
	if err != nil {
		panic(err.Error())
	}

	jsonData := string(body)
	var events []Event
	json.Unmarshal([]byte(jsonData), &events)

	for _, event := range events {
		if event.Type == "PushEvent" {
			for _, commit := range event.Payload.Commits {
				if commit.Author.Email != "" {
					return commit.Author.Email, nil
				}
			}
		}
	}

	switch {
	case staticData[name] != "":
		{
			return staticData[name], nil
		}
	default:
		{
			return "", fmt.Errorf("not found [%s] email", name)
		}
	}
}
