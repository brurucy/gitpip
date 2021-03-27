package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	githubBaseURL = "https://api.github.com/users/"
)

func NewGistApiRequest(name string) (*[]Gist, error) {

	var gistResponse []Gist

	logrus.Infof("Sending a request to get gists for user %s\n", name)

	resp, err := http.Get(githubBaseURL + name + "/gists")

	if err != nil {

		return nil, fmt.Errorf("Failed to get user %s gists with error %v\n", name, err)

	}

	logrus.Tracef("Got an ok response")

	if err := json.NewDecoder(resp.Body).Decode(&gistResponse); err != nil {

		return nil, fmt.Errorf("Failed to unmarshal %v\n", err)

	}

	return &gistResponse, nil

}

func NewPipedriveActivity(request *PipeAddActivityRequest) error {

	apiCALL := "https://" + os.Getenv("PIPEDRIVE_ORG") + ".pipedrive.com/v1/activities?api_token=" + os.Getenv("PIPEDRIVE_TOKEN")

	var marshalledRequest []byte
	marshalledRequest, err := json.Marshal(request)

	if err != nil {

		return fmt.Errorf("Failed to marshall the activity request %v", err)

	}

	logrus.Info("Sending a POST to Pipedrive's API")

	resp, err := http.Post(apiCALL, "application/json", bytes.NewBuffer(marshalledRequest))

	if err != nil && resp.StatusCode != 201 {

		return fmt.Errorf("Failed to POST %v", err)

	}

	logrus.Tracef("Successfully added new activity")

	resp.Body.Close()

	return nil

}

func GistTextDownloader(url string) (string, error) {

	var gistText string

	logrus.Info("Sending a request to get gist text\n")

	resp, err := http.Get(url)

	if err != nil {

		return "", fmt.Errorf("Failed to GET raw gist %v", err)

	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {

		return "", fmt.Errorf("Failed to read %v\n", err)

	}

	gistText = string(body)

	resp.Body.Close()

	logrus.Tracef("Succesfully downloaded Gist")

	return gistText, nil

}
