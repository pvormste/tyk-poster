package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/TykTechnologies/tyk/apidef"
)

type CreateAPIDefinitionModel struct {
	APIDefinition apidef.APIDefinition `json:"api_definition"`
}

func main() {
	fmt.Println(os.Args)
	if len(os.Args) < 3 {
		fmt.Println("Please provide API-Key and number of APIs to create")
		os.Exit(1)
	}

	apiKey := os.Args[1]
	numOfAPIs, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	apiNamePrefix := "Auto-Created"
	if len(os.Args) > 3 {
		apiNamePrefix = os.Args[3]
	}

	httpClient := &http.Client{}
	for i := 1; i <= numOfAPIs; i++ {
		name := fmt.Sprintf("%s-%d", apiNamePrefix, i)
		fmt.Printf("Creating API number %d with name %s", i, name)
		rawDef := apidef.DummyAPI()
		rawDef.Name = name
		sendApiDef(httpClient, apiKey, CreateAPIDefinitionModel{APIDefinition: rawDef})
		fmt.Println(" ......... DONE!")
	}

	fmt.Println("\nALL DONE!")
}


func sendApiDef(httpClient *http.Client, apiKey string, rawDef CreateAPIDefinitionModel) {
	rawDefBytes, err := json.Marshal(rawDef)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:3000/api/apis/", bytes.NewBuffer(rawDefBytes))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		var respBody map[string]interface{}
		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(respBodyBytes, &respBody)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Status: %d, %+v\n", resp.StatusCode, respBody)
		panic(errors.New("status is not 200"))
	}
}