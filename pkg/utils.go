package pkg

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	githubBaseURL = "https://api.github.com/users/"
)

func NewGistApiRequest(name string) (*[]Gist, error) {

	var gistResponse []Gist

	log.Printf("Sending a request for user %s\n", name)

	resp, err := http.Get(githubBaseURL + name + "/gists")

	if err != nil {

		log.Printf("Failed to get user %s with error %v\n", name, err)

		return nil, err

	}

	log.Printf("Got an ok response")

	if err := json.NewDecoder(resp.Body).Decode(&gistResponse); err != nil {

		log.Printf("Failed to unmarshal %v\n", err)

		return nil, err

	}

	return &gistResponse, nil

}

func NewPipedriveActivity(request *PipeAddActivityRequest) error {

	apiCALL := "https://" + os.Getenv("PIPEDRIVE_ORG") + ".pipedrive.com/v1/activities?api_token=" + os.Getenv("PIPEDRIVE_TOKEN")

	var marshalledRequest []byte
	marshalledRequest, err := json.Marshal(request)

	if err != nil {

		log.Printf("Failed to marshall the activity request %v", err)

		return err

	}

	log.Printf("Sending a POST to Pipedrive's API")

	resp, err := http.Post(apiCALL, "application/json", bytes.NewBuffer(marshalledRequest))

	if err != nil && resp.StatusCode != 201 {

		log.Printf("Failed to POST %v", err)

		return err

	}

	log.Printf("Successfully added new activity")

	resp.Body.Close()

	return nil

}

func GistTextDownloader(url string) (string, error) {

	var gistText string

	log.Print("Sending a request to get gist text\n")

	resp, err := http.Get(url)

	if err != nil {

		log.Print("Failed to get gist\n")

		return "", err

	}

	log.Printf("Got an ok response\n")

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {

		log.Printf("Failed to read %v\n", err)

		return "", err

	}

	gistText = string(body)

	resp.Body.Close()

	return gistText, nil

}
