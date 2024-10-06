// Copyright 2016 The Gogs Authors. All rights reserved.
// Copyright 2016 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import (
	"github.com/urfave/cli/v2"
)

var (
	// CmdAdmin represents the available admin sub-command.
	CmdAdmin = &cli.Command{
		Name:  "admin",
		Usage: "Perform common administrative operations",
		Subcommands: []*cli.Command{
			subcmdUser,
			subcmdAuth,
		},
	}

	subcmdAuth = &cli.Command{
		Name:  "auth",
		Usage: "Modify external auth providers",
		Subcommands: []*cli.Command{
			microcmdAuthAddOauth,
			microcmdAuthUpdateOauth,
			microcmdAuthAddLdapBindDn,
			microcmdAuthUpdateLdapBindDn,
			microcmdAuthAddLdapSimpleAuth,
			microcmdAuthUpdateLdapSimpleAuth,
			microcmdAuthAddSMTP,
			microcmdAuthUpdateSMTP,
			microcmdAuthList,
			microcmdAuthDelete,
		},
	}

	idFlag = &cli.Int64Flag{
		Name:  "id",
		Usage: "ID of authentication source",
	}
)
