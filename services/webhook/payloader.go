// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	"bytes"
	"fmt"
	"net/http"

	webhook_model "code.gitea.io/gitea/models/webhook"
	"code.gitea.io/gitea/modules/json"
	api "code.gitea.io/gitea/modules/structs"
	webhook_module "code.gitea.io/gitea/modules/webhook"
)

// payloadConvertor defines the interface to convert system payload to webhook payload
type payloadConvertor[T any] interface {
	// Create(*api.CreatePayload) (T, error)
	Package(*api.PackagePayload) (T, error)
}

func convertUnmarshalledJSON[T, P any](convert func(P) (T, error), data []byte) (t T, err error) {
	var p P
	if err = json.Unmarshal(data, &p); err != nil {
		return t, fmt.Errorf("could not unmarshal payload: %w", err)
	}
	return convert(p)
}

func newPayload[T any](rc payloadConvertor[T], data []byte, event webhook_module.HookEventType) (t T, err error) {
	switch event {
	// case webhook_module.HookEventCreate:
	// 	return convertUnmarshalledJSON(rc.Create, data)
	// case webhook_module.HookEventDelete:
	// 	return convertUnmarshalledJSON(rc.Delete, data)
	// case webhook_module.HookEventRepository:
	// 	return convertUnmarshalledJSON(rc.Repository, data)
	case webhook_module.HookEventPackage:
		return convertUnmarshalledJSON(rc.Package, data)
	}
	return t, fmt.Errorf("newPayload unsupported event: %s", event)
}

func newJSONRequest[T any](pc payloadConvertor[T], w *webhook_model.Webhook, t *webhook_model.HookTask, withDefaultHeaders bool) (*http.Request, []byte, error) {
	payload, err := newPayload(pc, []byte(t.PayloadContent), t.EventType)
	if err != nil {
		return nil, nil, err
	}

	body, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	method := w.HTTPMethod
	if method == "" {
		method = http.MethodPost
	}

	req, err := http.NewRequest(method, w.URL, bytes.NewReader(body))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	if withDefaultHeaders {
		return req, body, addDefaultHeaders(req, []byte(w.Secret), t, body)
	}
	return req, body, nil
}
