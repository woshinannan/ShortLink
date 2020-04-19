/*
@Time       : 2020/4/1 23:39
@Author     : stevinpan
@File       : Main
@Software   : GoLand
@Description: <>
*/
package main

import (
	"ShortLink/private/service"
	"fmt"
)

func main() {
	var longLink = "xxcvcvffgtghh"
	shortLink, err := service.Process(longLink, "")
	if err != nil {
		fmt.Printf("failed, err=%+v\n", err)
	} else {
		fmt.Println(shortLink)
	}
}
