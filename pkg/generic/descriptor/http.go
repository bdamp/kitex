/*
 * Copyright 2021 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package descriptor

import (
	"net/http"
	"net/url"
)

// Cookies ...
type Cookies map[string]string

// MIMEType ...
type MIMEType string

const (
	MIMEApplicationJson     = "application/json"
	MIMEApplicationForm     = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf = "application/x-protobuf"
)

// HTTPRequest ...
type HTTPRequest struct {
	Header      http.Header
	Query       url.Values
	Cookies     Cookies
	Method      string
	Host        string
	Path        string
	Params      *Params // path params
	RawBody     []byte
	Body        interface{}
	ContentType MIMEType
}

// HTTPResponse ...
type HTTPResponse struct {
	Header      http.Header
	StatusCode  int32
	Body        interface{}
	ContentType MIMEType
	Renderer    Renderer
}

// NewHTTPJsonResponse ...
func NewHTTPJsonResponse() *HTTPResponse {
	return &HTTPResponse{
		Header:      http.Header{},
		ContentType: MIMEApplicationJson,
		Body:        map[string]interface{}{},
		Renderer:    JsonRenderer{},
	}
}

func NewHTTPPbResponse(initBody interface{}) *HTTPResponse {
	return &HTTPResponse{
		Header:      http.Header{},
		ContentType: MIMEApplicationProtobuf,
		Body:        initBody,
		Renderer:    PbRenderer{},
	}
}

// NewHTTPResponse init response with given MIMEType and body
func NewHTTPResponse(contentType MIMEType, initBody interface{}, renderer Renderer) *HTTPResponse {
	return &HTTPResponse{
		Header:      http.Header{},
		ContentType: contentType,
		Body:        initBody,
		Renderer:    renderer,
	}
}

// Write to ResponseWriter
func (resp *HTTPResponse) Write(w http.ResponseWriter) error {
	w.WriteHeader(int(resp.StatusCode))
	for k := range resp.Header {
		w.Header().Set(k, resp.Header.Get(k))
	}

	resp.Renderer.WriteContentType(w)
	return resp.Renderer.Render(w, resp.Body)
}

// Param in request path
type Param struct {
	Key   string
	Value string
}

// Params and recyclable
type Params struct {
	params   []Param
	recycle  func(*Params)
	recycled bool
}

// Recycle the Params
func (ps *Params) Recycle() {
	if ps.recycled {
		return
	}
	ps.recycled = true
	ps.recycle(ps)
}

// ByName search Param by given name
func (ps *Params) ByName(name string) string {
	for _, p := range ps.params {
		if p.Key == name {
			return p.Value
		}
	}
	return ""
}
