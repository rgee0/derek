package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/types"
)

const maintainersFileEnv = "maintainers_file"
const defaultMaintFile = ".DEREK.yml"

func getEnv(envVar string, assumed string) string {
	if value, exists := os.LookupEnv(envVar); exists {
		return value
	}
	return assumed
}

func enabledFeature(attemptedFeature string, controls *types.DerekConfig) bool {

	featureEnabled := false

	for _, availableFeature := range controls.Features {
		if strings.EqualFold(attemptedFeature, availableFeature) {
			featureEnabled = true
			break
		}
	}
	return featureEnabled
}

func permittedUserFeature(attemptedFeature string, controls *types.DerekConfig, user string) bool {

	permitted := false
	featureEnabled := false

	for _, availableFeature := range controls.Features {
		if strings.EqualFold(attemptedFeature, availableFeature) {
			featureEnabled = true
			break
		}
	}

	if featureEnabled {
		for _, maintainer := range controls.Maintainers {
			if strings.EqualFold(user, maintainer) {
				permitted = true
				break
			}
		}
	}

	return permitted
}

func getControls(owner string, repository string) (*types.DerekConfig, error) {

	maintainersFile := getEnv(maintainersFileEnv, defaultMaintFile)
	maintainersFile = fmt.Sprintf("https://github.com/%s/%s/raw/master/%s", owner, repository, strings.Trim(maintainersFile, "/"))

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	req, _ := http.NewRequest(http.MethodGet, maintainersFile, nil)

	res, resErr := client.Do(req)
	if resErr != nil {
		log.Fatalln(resErr)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalln(fmt.Sprintf("HTTP Status code: %d while fetching maintainers list (%s)", res.StatusCode, maintainersFile))
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	bytesOut, _ := ioutil.ReadAll(res.Body)
	var controls types.DerekConfig

	err := yaml.Unmarshal(bytesOut, &controls)
	if err != nil {
		return nil, err
	}

	return &controls, nil
}
