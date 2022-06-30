package main

import (
	"encoding/json"
	"log"
	"tagu/aws"
)

func main() {
	tags := aws.Tags{
		Account: "123456789012",
		Region:  "us-east-1",
		// RoleName: "role-name",
	}
	err := aws.Run(&tags)
	if err != nil {
		log.Fatalf("An error occured %v", err)
	}
	json, err := json.Marshal(tags.Output)
	if err != nil {
		log.Fatalf("failed to return json data, %v", err)
	}
	log.Println(string(json))
}
