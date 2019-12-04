package main

import (
	"fmt"
	transporthttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/kit/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	//"log"
	"net/http"
	"os"
)

func main(){
	logger := kitlog.NewLogfmtLogger(os.Stderr)

	fieldKeys:=[]string{"method","error"}
	requestCount:=kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace:"my_group",
		Subsystem:"string_service",
		Name:"request_count",
		Help:"Number of requests received,",
	},fieldKeys)
	requestLatency:=kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace:   "my_group",
		Subsystem:   "string_service",
		Name:        "request_latency_microseconds",
		Help:        "Total duration of requests in microseconds",
	},fieldKeys)
	countRequest:=kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace:   "my_group",
		Subsystem:   "string_service",
		Name:        "count_result",
		Help:        "The result of each count method.",
	},[]string{})// no field here


	mid := loggingMiddleware(kitlog.With(logger, "method", "uppercase"))
	var svc StringService
	svc=stringService{}
	svc=loggingMiddleware{logger,svc}
	svc=instrumentingMiddleware{requestCount,requestLatency,countRequest,svc}



	uppercaseHandler:=transporthttp.NewServer(
		mid(makeUppercaseEndpoint(svc)),
		decodeUpperCaseRequest,
		encodeResponse,
		)
	countHandler:=transporthttp.NewServer(
		mid(makeCountEndpoint(svc)),
		decodeCountCaseRequest,
		encodeResponse,
		)
	http.Handle("/uppercase",uppercaseHandler)
	http.Handle("/count",countHandler)
	err:=http.ListenAndServe(":8082",nil)
	if err!=nil{
		fmt.Println(" listen server failed :",err.Error())
	}
	fmt.Println("service has started...")
}
