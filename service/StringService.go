package main

import (
	"context"
	"errors"
	"strings"
)

type StringService interface {
	Uppercase(string)(string ,error)//字符串大写话
	Count(string )int//计算字符串长度
}

type stringService struct {
}

//
func (stringService)Uppercase(_ context.Context,s string)(string,error){
	if s==""{
		return "",errors.New("empty string")
	}
	return strings.ToUpper(s),nil
}

func (stringService)Count(s string)int{
	return len(s)
}


