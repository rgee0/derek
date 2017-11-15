package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/alexellis/derek/auth"
	"github.com/alexellis/derek/types"
)

const dcoCheck = "dco_check"
const comments = "comments"
const maintainersFileEnv = "maintainers_file"
const defaultMaintFile = ".DEREK.yml"

func hmacValidation() bool {
	val := os.Getenv("validate_hmac")
	return len(val) > 0 && (val == "1" || val == "true")
}

func getEnv(envVar string, assumed string) string {
	if value, exists := os.LookupEnv(envVar); exists {
		return value
	}
	return assumed
}

func enabledFeature(attemptedFeature string, reqRep types.Repository) bool {

	derekControls, err := getControls(reqRep.Owner.Login, reqRep.Name)

	if err != nil {
		log.Fatalf("Unable to verify access maintainers file: %s/%s", reqRep.Owner.Login, reqRep.Name)
	}

	featureEnabled := false

	for _, availableFeature := range derekControls.Features {
		if attemptedFeature == availableFeature {
			featureEnabled = true
			break
		}
	}
	return featureEnabled
}

func permittedUserFeature(attemptedFeature string, reqRep types.Repository, user string) bool {

	derekControls, err := getControls(reqRep.Owner.Login, reqRep.Name)

	if err != nil {
		log.Fatalf("Unable to verify access maintainers file: %s/%s", reqRep.Owner.Login, reqRep.Name)
	}

	permitted := false
	featureEnabled := false

	for _, availableFeature := range derekControls.Features {
		if attemptedFeature == availableFeature {
			featureEnabled = true
			break
		}
	}

	if featureEnabled {
		for _, maintainer := range derekControls.Maintainers {
			if user == maintainer {
				permitted = true
				break
			}
		}
	}

	return permitted
}

func getControls(owner string, repository string) (*types.DerekControl, error) {

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	maintainersFile := getEnv(maintainersFileEnv, defaultMaintFile)
	maintainersFile = fmt.Sprintf("https://github.com/%s/%s/raw/master/%s", owner, repository, strings.Trim(maintainersFile, "/"))

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
	var controls types.DerekControl

	err := yaml.Unmarshal(bytesOut, &controls)
	if err != nil {
		return nil, err
	}

	return &controls, nil
}

func main() {
	bytesIn, _ := ioutil.ReadAll(os.Stdin)

	xHubSignature := os.Getenv("Http_X_Hub_Signature")
	if hmacValidation() && len(xHubSignature) == 0 {
		log.Fatal("must provide X_Hub_Signature")
		return
	}

	if len(xHubSignature) > 0 {
		secretKey := os.Getenv("secret_key")

		err := auth.ValidateHMAC(bytesIn, xHubSignature, secretKey)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	}

	// HMAC Validated or not turned on.
	eventType := os.Getenv("Http_X_Github_Event")

	switch eventType {
	case "pull_request":
		req := types.PullRequestOuter{}
		if err := json.Unmarshal(bytesIn, &req); err != nil {
			log.Fatalf("Cannot parse input %s", err.Error())
		}

		customer, err := auth.IsCustomer(req.Repository)
		if err != nil {
			log.Fatalf("Unable to verify customer: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		} else if !customer {
			log.Fatalf("No customer found for: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		}

		if enabledFeature(dcoCheck, req.Repository) {
			handlePullRequest(req)
		}
		break

	case "issue_comment":
		req := types.IssueCommentOuter{}
		if err := json.Unmarshal(bytesIn, &req); err != nil {
			log.Fatalf("Cannot parse input %s", err.Error())
		}

		customer, err := auth.IsCustomer(req.Repository)
		if err != nil {
			log.Fatalf("Unable to verify customer: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		} else if !customer {
			log.Fatalf("No customer found for: %s/%s", req.Repository.Owner.Login, req.Repository.Name)
		}

		if permittedUserFeature(comments, req.Repository, req.Comment.User.Login) {
			handleComment(req)
		}
		break
	default:
		log.Fatalln("X_Github_Event want: ['pull_request', 'issue_comment'], got: " + eventType)
	}
}
