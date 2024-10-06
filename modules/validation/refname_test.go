// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package validation

import (
	"testing"

	"gitea.com/go-chi/binding"
)

var gitRefNameValidationTestCases = []validationTestCase{
	{
		description: "Reference name contains only characters",
		data: TestForm{
			BranchName: "test",
		},
		expectedErrors: binding.Errors{},
	},
	{
		description: "Reference name contains single slash",
		data: TestForm{
			BranchName: "feature/test",
		},
		expectedErrors: binding.Errors{},
	},
	{
		description: "Reference name has allowed special characters",
		data: TestForm{
			BranchName: "debian/1%1.6.0-2",
		},
		expectedErrors: binding.Errors{},
	},
}

func Test_GitRefNameValidation(t *testing.T) {
	AddBindingRules()

	for _, testCase := range gitRefNameValidationTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			performValidationTest(t, testCase)
		})
	}
}
