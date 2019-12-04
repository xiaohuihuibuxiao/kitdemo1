package main

import (
	"context"
	"encoding/json"
	"net/http"
)

func decodeUpperCaseRequest(_ context.Context,r *http.Request)(interface{},error){
	var request uppercaseRequest
	if err:=json.NewDecoder(r.Body).Decode(&request);err!=nil{
		return nil,err
	}
	return request,nil
}

func decodeCountCaseRequest(_ context.Context,r *http.Request)(interface{},error){
	var request countRequest
	if err:=json.NewDecoder(r.Body).Decode(&request);err!=nil{
		return nil,err
	}
	return request,nil
}

func encodeResponse(_ context.Context,w http.ResponseWriter,response interface{})error{
	return json.NewEncoder(w).Encode(response)
}


