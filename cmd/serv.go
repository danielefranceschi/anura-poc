// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2016 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"unicode"

	"code.gitea.io/gitea/models/perm"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/private"
	"code.gitea.io/gitea/modules/setting"

	"github.com/urfave/cli/v2"
)

// CmdServ represents the available serv sub-command.
var CmdServ = &cli.Command{
	Name:        "serv",
	Usage:       "(internal) Should only be called by SSH shell",
	Description: "Serv provides access auth for repositories",
	Before:      PrepareConsoleLoggerLevel(log.FATAL),
	Action:      runServ,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "enable-pprof",
		},
		&cli.BoolFlag{
			Name: "debug",
		},
	},
}

func setup(ctx context.Context, debug bool) {
	if debug {
		setupConsoleLogger(log.TRACE, false, os.Stderr)
	} else {
		setupConsoleLogger(log.FATAL, false, os.Stderr)
	}
	setting.MustInstalled()
	if _, err := os.Stat(setting.RepoRootPath); err != nil {
		_ = fail(ctx, "Unable to access repository path", "Unable to access repository path %q, err: %v", setting.RepoRootPath, err)
		return
	}
}

var (
	allowedCommands     = map[string]perm.AccessMode{}
	alphaDashDotPattern = regexp.MustCompile(`[^\w-\.]`)
)

// fail prints message to stdout, it's mainly used for git serv and git hook commands.
// The output will be passed to git client and shown to user.
func fail(ctx context.Context, userMessage, logMsgFmt string, args ...any) error {
	if userMessage == "" {
		userMessage = "Internal Server Error (no specific error)"
	}

	// There appears to be a chance to cause a zombie process and failure to read the Exit status
	// if nothing is outputted on stdout.
	_, _ = fmt.Fprintln(os.Stdout, "")
	_, _ = fmt.Fprintln(os.Stderr, "Anura:", userMessage)

	if logMsgFmt != "" {
		logMsg := fmt.Sprintf(logMsgFmt, args...)
		if !setting.IsProd {
			_, _ = fmt.Fprintln(os.Stderr, "Anura:", logMsg)
		}
		if userMessage != "" {
			if unicode.IsPunct(rune(userMessage[len(userMessage)-1])) {
				logMsg = userMessage + " " + logMsg
			} else {
				logMsg = userMessage + ". " + logMsg
			}
		}
		_, _ = fmt.Fprintln(os.Stderr, "User:", logMsg)
	}
	return cli.Exit("", 1)
}

// handleCliResponseExtra handles the extra response from the cli sub-commands
// If there is a user message it will be printed to stdout
// If the command failed it will return an error (the error will be printed by cli framework)
func handleCliResponseExtra(extra private.ResponseExtra) error {
	if extra.UserMsg != "" {
		_, _ = fmt.Fprintln(os.Stdout, extra.UserMsg)
	}
	if extra.HasError() {
		return cli.Exit(extra.Error, 1)
	}
	return nil
}

func runServ(c *cli.Context) error {
	ctx, cancel := installSignals()
	defer cancel()

	// FIXME: This needs to internationalised
	setup(ctx, c.Bool("debug"))

	println("Anura: SSH has been disabled")
	return nil
}
