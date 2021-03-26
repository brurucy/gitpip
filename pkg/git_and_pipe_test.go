package src

import (
	"fmt"
	"testing"
)

func TestApiUnmarshal(t *testing.T) {

	resp, err := NewGistApiRequest("brurucy")

	if err != nil {

		t.Fatalf("%v", err)

	}

	fmt.Println(resp)

}

func TestPipeApi(t *testing.T) {

	req := &PipeAddActivityRequest{
		Subject: "Testing-again",
		Note:    "Gist",
		Done:    false,
	}

	err := NewPipedriveActivity(req)

	if err != nil {

		t.Fatalf("%v", err)

	}

}

func TestGistDownloader(t *testing.T) {

	gistUrl := "https://gist.githubusercontent.com/robotsnowfall/e823db4efe48513088bc74f08670d78c/raw/614f4fe4713a7fdef59b469508e6960286b94d92/gistfile1.txt"

	_, err := GistTextDownloader(gistUrl)

	if err != nil {

		t.Fatalf("%v", err)

	}

}
