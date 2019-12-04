package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"time"
)

type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V string `json:"v"`
	Err string `json:"err,omitempty"`
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}


//=======
func makeUppercaseEndpoint(svc StringService)endpoint.Endpoint{
	return func(_ context.Context, request interface{}) ( interface{},  error){
		req:=request.(uppercaseRequest)
		v,err:=svc.Uppercase(req.S)
		if err!=nil{
			return uppercaseResponse{
				V:  v,
				Err: err.Error(),
			},err
		}
		return uppercaseResponse{
			V:  v,
			Err: "",
		},nil
	}
}

func makeCountEndpoint(svc StringService)endpoint.Endpoint{
	return func(_ context.Context, request interface{}) ( interface{},  error){
		req:=request.(countRequest)
		v:=svc.Count(req.S)
		return countResponse{V:v},nil
	}
}

//-----中间件

//--日志---
func loggingMiddleware(logger kitlog.Logger)endpoint.Middleware{
	return func(next endpoint.Endpoint)endpoint.Endpoint{
		return func(ctx context.Context,request interface{})(interface{},error){
			//写任何在中途想做的事，这里是作为日志中间件
			logger.Log("msg","callin endpoint")
			defer logger.Log("msg","called endpoint")
			return next(ctx,request)
		}
	}
}

//--instrumentation
type instrumentingMiddleware struct {
	requestCount metrics.Counter
	requestLatency metrics.Histogram
	countResult metrics.Histogram
	next StringService
}

func (mw instrumentingMiddleware)Uppercase(s string)(output string,err error){
	defer func(begin time.Time){
		lvs:=[]string{"method","uppercase","error",fmt.Sprint(err!=nil)	}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output,err=mw.next.Uppercase(s)
	return
}



func (mw instrumentingMiddleware)Count(s string)(n int){
	defer func(begin time.Time) {
		lvs := []string{"method", "count", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.countResult.Observe(float64(n))
	}(time.Now())

	n = mw.next.Count(s)//TODO------------
	return
}