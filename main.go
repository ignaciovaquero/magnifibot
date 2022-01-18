package main

import (
	"fmt"
	"log"
	"time"

	"github.com/igvaquero18/magnifibot/archimadrid"
)

func main() {
	client := archimadrid.NewClient()
	gospel, err := client.GetGospel(time.Now())
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(gospel.Day)
	fmt.Println(gospel.Title)
	fmt.Println(gospel.Reference)
	fmt.Println(gospel.Content)
}
