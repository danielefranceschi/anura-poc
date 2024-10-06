// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package doctor

import (
	"context"
	"fmt"
	"os/exec"

	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/setting"
)

func checkScriptType(ctx context.Context, logger log.Logger, autofix bool) error {
	path, err := exec.LookPath(setting.ScriptType)
	if err != nil {
		logger.Critical("ScriptType \"%q\" is not on the current PATH. Error: %v", setting.ScriptType, err)
		return fmt.Errorf("ScriptType \"%q\" is not on the current PATH. Error: %w", setting.ScriptType, err)
	}
	logger.Info("ScriptType %s is on the current PATH at %s", setting.ScriptType, path)
	return nil
}

func init() {
	Register(&Check{
		Title:     "Check if SCRIPT_TYPE is available",
		Name:      "script-type",
		IsDefault: false,
		Run:       checkScriptType,
		Priority:  5,
	})

}
