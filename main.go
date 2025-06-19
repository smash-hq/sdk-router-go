package main

import (
	"context"
	"encoding/json"
	"github.com/scrapeless-ai/sdk-go/scrapeless/actor"
	"github.com/scrapeless-ai/sdk-go/scrapeless/log"
	"strings"
	"time"
)

var req = &RequestParams{}

type RequestParams struct {
	RunID  string `json:"run_id"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   string `json:"body"`
	Header string `json:"header"`
}

func main() {
	ctx := context.Background()
	a := actor.New()
	if err := a.Input(req); err != nil {
		log.Warnf("input error: %v", err)
	}
	var header = &map[string]string{}
	err := json.Unmarshal([]byte(req.Header), header)
	if err != nil {
		log.Errorf("json.Unmarshal error: %v", err)
	}
	request, err := a.Router.Request(req.RunID, req.Method, req.Path, strings.NewReader(req.Body), *header)
	if err != nil {
		log.Errorf("request error: %v", err)
		return
	}
	inputSave(a, err, ctx, request)
	duration := 10 * time.Second
	log.Infof("sleep %v ......", duration)
	time.Sleep(duration)
	log.Infof("wake up")
}

func inputSave(a *actor.Actor, err error, ctx context.Context, data []byte) {
	items, err := a.AddItems(ctx, []map[string]interface{}{
		{"run_id": req.RunID, "method": req.Method, "path": req.Path, "body": req.Body, "header": req.Header, "data": string(data)},
	})
	if err != nil {
		log.Warnf("input--> save input failed: %v", err)
	}
	if !items {
		log.Infof("input--> save input failed, isErr: %v", items)
	}
}
