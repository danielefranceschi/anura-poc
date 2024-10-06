// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

// Package private contains all internal routes. The package name "internal" isn't usable because Golang reserves it for disabling cross-package usage.
package private

import (
	"net/http"
	"strings"

	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/private"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/services/context"

	"gitea.com/go-chi/binding"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
)

// CheckInternalToken check internal token is set
func CheckInternalToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		tokens := req.Header.Get("Authorization")
		fields := strings.SplitN(tokens, " ", 2)
		if setting.InternalToken == "" {
			log.Warn(`The INTERNAL_TOKEN setting is missing from the configuration file: %q, internal API can't work.`, setting.CustomConf)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		if len(fields) != 2 || fields[0] != "Bearer" || fields[1] != setting.InternalToken {
			log.Debug("Forbidden attempt to access internal url: Authorization header: %s", tokens)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		} else {
			next.ServeHTTP(w, req)
		}
	})
}

// bind binding an obj to a handler
func bind[T any](_ T) any {
	return func(ctx *context.PrivateContext) {
		theObj := new(T) // create a new form obj for every request but not use obj directly
		binding.Bind(ctx.Req, theObj)
		web.SetForm(ctx, theObj)
	}
}

// Routes registers all internal APIs routes to web application.
// These APIs will be invoked by internal commands for example `gitea serv` and etc.
func Routes() *web.Router {
	r := web.NewRouter()
	r.Use(context.PrivateContexter())
	r.Use(CheckInternalToken)
	// Log the real ip address of the request from SSH is really helpful for diagnosing sometimes.
	// Since internal API will be sent only from Gitea sub commands and it's under control (checked by InternalToken), we can trust the headers.
	r.Use(chi_middleware.RealIP)

	r.Post("/manager/shutdown", Shutdown)
	r.Post("/manager/restart", Restart)
	r.Post("/manager/reload-templates", ReloadTemplates)
	r.Post("/manager/flush-queues", bind(private.FlushOptions{}), FlushQueues)
	r.Post("/manager/pause-logging", PauseLogging)
	r.Post("/manager/resume-logging", ResumeLogging)
	r.Post("/manager/release-and-reopen-logging", ReleaseReopenLogging)
	r.Post("/manager/set-log-sql", SetLogSQL)
	r.Post("/manager/add-logger", bind(private.LoggerOptions{}), AddLogger)
	r.Post("/manager/remove-logger/{logger}/{writer}", RemoveLogger)
	r.Get("/manager/processes", Processes)

	return r
}
