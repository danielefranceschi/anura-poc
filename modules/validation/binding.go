// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package validation

import (
	"fmt"
	"regexp"
	"strings"

	"gitea.com/go-chi/binding"
	"github.com/gobwas/glob"
)

const (
	// ErrGlobPattern is returned when glob pattern is invalid
	ErrGlobPattern = "GlobPattern"
	// ErrRegexPattern is returned when a regex pattern is invalid
	ErrRegexPattern = "RegexPattern"
	// ErrUsername is username error
	ErrUsername = "UsernameError"
)

// AddBindingRules adds additional binding rules
func AddBindingRules() {
	addValidURLBindingRule()
	addValidSiteURLBindingRule()
	addGlobPatternRule()
	addRegexPatternRule()
	addGlobOrRegexPatternRule()
	addUsernamePatternRule()
}

func addValidURLBindingRule() {
	// URL validation rule
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return strings.HasPrefix(rule, "ValidUrl")
		},
		IsValid: func(errs binding.Errors, name string, val any) (bool, binding.Errors) {
			str := fmt.Sprintf("%v", val)
			if len(str) != 0 && !IsValidURL(str) {
				errs.Add([]string{name}, binding.ERR_URL, "Url")
				return false, errs
			}

			return true, errs
		},
	})
}

func addValidSiteURLBindingRule() {
	// URL validation rule
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return strings.HasPrefix(rule, "ValidSiteUrl")
		},
		IsValid: func(errs binding.Errors, name string, val any) (bool, binding.Errors) {
			str := fmt.Sprintf("%v", val)
			if len(str) != 0 && !IsValidSiteURL(str) {
				errs.Add([]string{name}, binding.ERR_URL, "Url")
				return false, errs
			}

			return true, errs
		},
	})
}

func addGlobPatternRule() {
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "GlobPattern"
		},
		IsValid: globPatternValidator,
	})
}

func globPatternValidator(errs binding.Errors, name string, val any) (bool, binding.Errors) {
	str := fmt.Sprintf("%v", val)

	if len(str) != 0 {
		if _, err := glob.Compile(str); err != nil {
			errs.Add([]string{name}, ErrGlobPattern, err.Error())
			return false, errs
		}
	}

	return true, errs
}

func addRegexPatternRule() {
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "RegexPattern"
		},
		IsValid: regexPatternValidator,
	})
}

func regexPatternValidator(errs binding.Errors, name string, val any) (bool, binding.Errors) {
	str := fmt.Sprintf("%v", val)

	if _, err := regexp.Compile(str); err != nil {
		errs.Add([]string{name}, ErrRegexPattern, err.Error())
		return false, errs
	}

	return true, errs
}

func addGlobOrRegexPatternRule() {
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "GlobOrRegexPattern"
		},
		IsValid: func(errs binding.Errors, name string, val any) (bool, binding.Errors) {
			str := strings.TrimSpace(fmt.Sprintf("%v", val))

			if len(str) >= 2 && strings.HasPrefix(str, "/") && strings.HasSuffix(str, "/") {
				return regexPatternValidator(errs, name, str[1:len(str)-1])
			}
			return globPatternValidator(errs, name, val)
		},
	})
}

func addUsernamePatternRule() {
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "Username"
		},
		IsValid: func(errs binding.Errors, name string, val any) (bool, binding.Errors) {
			str := fmt.Sprintf("%v", val)
			if !IsValidUsername(str) {
				errs.Add([]string{name}, ErrUsername, "invalid username")
				return false, errs
			}
			return true, errs
		},
	})
}

func portOnly(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return ""
	}
	if i := strings.Index(hostport, "]:"); i != -1 {
		return hostport[i+len("]:"):]
	}
	if strings.Contains(hostport, "]") {
		return ""
	}
	return hostport[colon+len(":"):]
}

func validPort(p string) bool {
	for _, r := range []byte(p) {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
