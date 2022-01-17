package main

import (
	"log"
	dynamo "nst-go-course-deliverable/aws"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
)

var (
	UUID string
	CFG  aws.Config
)

func init() {
	UUID = uuid.New().String()
	CFG = dynamo.LoadConfig()
}

func main() {
	if dynamo.Insert(CFG, UUID) {
		if dynamo.Read(CFG, UUID) {
			log.Println("Success!")
		} else {
			log.Println("!STATUS could not read")
		}
	} else {
		log.Println("!STATUS could not insert")
	}
}
