package main

import (
	"os"
	"testing"
)

var envVarOptions = []struct {
	title          string
	envName        string
	envConfigVal   string
	envExpectedVal string
}{
	{
		title:          "envvar correctly set",
		envName:        "maintainers_file",
		envConfigVal:   "DEREK",
		envExpectedVal: "DEREK",
	},
	{
		title:          "Misspelt envVar Name",
		envName:        "maintainers_fill",
		envConfigVal:   "DEREK",
		envExpectedVal: "MAINTAINERS",
	},
	{
		title:          "envVar doesnt exist",
		envName:        "",
		envConfigVal:   "",
		envExpectedVal: "MAINTAINERS",
	},
}

func Test_getEnv(t *testing.T) {

	for _, test := range envVarOptions {
		t.Run(test.title, func(t *testing.T) {

			os.Setenv(test.envName, test.envConfigVal)

			envvar := getEnv("maintainers_file", "MAINTAINERS")

			if envvar != test.envExpectedVal {
				t.Errorf("Maintainers File - wanted: %s, found %s", test.envExpectedVal, envvar)
			}
			os.Unsetenv(test.envName)
		})
	}
}

var actionOptions = []struct {
	title          string
	body           string
	expectedAction string
}{
	{
		title:          "Correct reopen command",
		body:           "Derek reopen",
		expectedAction: "reopen",
	},
	{ //this case replaces Test_Parsing_Close
		title:          "Correct close command",
		body:           "Derek close",
		expectedAction: "close",
	},
	{
		title:          "invalid command",
		body:           "Derek dance",
		expectedAction: "",
	},
	{
		title:          "Longer reopen command",
		body:           "Derek reopen: ",
		expectedAction: "reopen",
	},
	{
		title:          "Longer close command",
		body:           "Derek close: ",
		expectedAction: "close",
	},
}

func Test_Parsing_OpenClose(t *testing.T) {

	for _, test := range actionOptions {
		t.Run(test.title, func(t *testing.T) {

			action := parse(test.body)

			if action.Type != test.expectedAction {
				t.Errorf("Action - want: %s, got %s", test.expectedAction, action.Type)
			}

		})
	}
}

var labelOptions = []struct {
	title        string
	body         string
	expectedType string
	expectedVal  string
}{
	{ //this case replaces Test_Parsing_AddLabel
		title:        "Add label of demo",
		body:         "Derek add label: demo",
		expectedType: "AddLabel",
		expectedVal:  "demo",
	},
	{
		title:        "Remove label of demo",
		body:         "Derek remove label: demo",
		expectedType: "RemoveLabel",
		expectedVal:  "demo",
	},
	{
		title:        "Invalid label action",
		body:         "Derek peel label: demo",
		expectedType: "",
		expectedVal:  "",
	},
}

func Test_Parsing_Labels(t *testing.T) {

	for _, test := range labelOptions {
		t.Run(test.title, func(t *testing.T) {

			action := parse(test.body)
			if action.Type != test.expectedType || action.Value != test.expectedVal {
				t.Errorf("Action - wanted: %s, got %s\nLabel - wanted: %s, got %s", test.expectedType, action.Type, test.expectedVal, action.Value)
			}
		})
	}
}

var assignmentOptions = []struct {
	title        string
	body         string
	expectedType string
	expectedVal  string
}{
	{
		title:        "Assign to burt",
		body:         "Derek assign: burt",
		expectedType: "Assign",
		expectedVal:  "burt",
	},
	{
		title:        "Unassign burt",
		body:         "Derek unassign: burt",
		expectedType: "Unassign",
		expectedVal:  "burt",
	},
	{
		title:        "Assign to me",
		body:         "Derek assign: me",
		expectedType: "Assign",
		expectedVal:  "me",
	},
	{
		title:        "Unassign me",
		body:         "Derek unassign: me",
		expectedType: "Unassign",
		expectedVal:  "me",
	},
	{
		title:        "Invalid assignment action",
		body:         "Derek consign: burt",
		expectedType: "",
		expectedVal:  "",
	},
	{
		title:        "Unassign blank",
		body:         "Derek unassign: ",
		expectedType: "",
		expectedVal:  "",
	},
}

func Test_Parsing_Assignments(t *testing.T) {

	for _, test := range assignmentOptions {
		t.Run(test.title, func(t *testing.T) {

			action := parse(test.body)
			if action.Type != test.expectedType || action.Value != test.expectedVal {
				t.Errorf("Action - wanted: %s, got %s\nMaintainer - wanted: %s, got %s", test.expectedType, action.Type, test.expectedVal, action.Value)
			}
		})
	}
}

var titleOptions = []struct {
	title        string
	body         string
	expectedType string
	expectedVal  string
}{
	{
		title:        "Set Title",
		body:         "Derek set title: This is a really great Title!",
		expectedType: "SetTitle",
		expectedVal:  "This is a really great Title!",
	},
	{
		title:        "Mis-spelling of title",
		body:         "Derek set titel: This is a really great Title!",
		expectedType: "",
		expectedVal:  "",
	},
	{
		title:        "Empty Title",
		body:         "Derek set title: ",
		expectedType: "",
		expectedVal:  "",
	},
	{
		title:        "Empty Title (Double Space)",
		body:         "Derek set title:  ",
		expectedType: "SetTitle",
		expectedVal:  "",
	},
}

func Test_Parsing_Titles(t *testing.T) {

	for _, test := range titleOptions {
		t.Run(test.title, func(t *testing.T) {

			action := parse(test.body)
			if action.Type != test.expectedType || action.Value != test.expectedVal {
				t.Errorf("Action - wanted: %s, got %s\nLabel - wanted: %s, got %s", test.expectedType, action.Type, test.expectedVal, action.Value)
			}
		})
	}
}

var validCommandsOptions = []struct {
	title        string
	body         string
	trigger      string
	expectedBool bool
}{
	{
		title:        "Valid Add Label",
		body:         "Derek add label: labelName",
		trigger:      "Derek add label: ",
		expectedBool: true,
	},
	{
		title:        "Valid Remove Label",
		body:         "Derek remove label: labelName",
		trigger:      "Derek remove label: ",
		expectedBool: true,
	},
	{
		title:        "Valid Assign",
		body:         "Derek assign: Burt",
		trigger:      "Derek assign: ",
		expectedBool: true,
	},
	{
		title:        "Valid Unassign",
		body:         "Derek unassign: Burt",
		trigger:      "Derek unassign: ",
		expectedBool: true,
	},
	{
		title:        "Valid close",
		body:         "Derek close",
		trigger:      "Derek close",
		expectedBool: true,
	},
	{
		title:        "Valid reopen",
		body:         "Derek reopen",
		trigger:      "Derek reopen",
		expectedBool: true,
	},
	{
		title:        "Valid Unassign",
		body:         "Derek set title: This is golden",
		trigger:      "Derek set title: ",
		expectedBool: true,
	},
	{
		title:        "Body length > Trigger length with substring match",
		body:         "This body is greater than the legth of the trigger",
		trigger:      "This body is greater",
		expectedBool: true,
	},
	{
		title:        "Body length < Trigger length with substring match",
		body:         "This body is not",
		trigger:      "This body is not greater than the legth of the trigger",
		expectedBool: false,
	},
	{
		title:        "Body length = Trigger length with substring match with colon",
		body:         "This body is the same: ",
		trigger:      "This body is the same: ",
		expectedBool: false,
	},
	{
		title:        "Body the same as Trigger without colon",
		body:         "This body is the same as trigger",
		trigger:      "This body is the same as trigger",
		expectedBool: true,
	},
	{
		title:        "Body not the same as Trigger without colon",
		body:         "This body is the not same as trigger",
		trigger:      "This body is the same as trigger",
		expectedBool: false,
	},
}

func Test_isValidCommand(t *testing.T) {

	for _, test := range validCommandsOptions {
		t.Run(test.title, func(t *testing.T) {

			isValid := isValidCommand(test.body, test.trigger)
			if isValid != test.expectedBool {
				t.Errorf("IsValid(%s,%s) yielded: %t wanted: %t", test.body, test.trigger, isValid, test.expectedBool)
			}
		})
	}
}
