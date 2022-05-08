package main

import (
	"context"
	"encoding/json"
	"log"
	"tagu/aws"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the ResourceGroupsTagging client
	client := resourcegroupstaggingapi.NewFromConfig(cfg)

	params := &resourcegroupstaggingapi.GetResourcesInput{}

	paginator := resourcegroupstaggingapi.NewGetResourcesPaginator(client, params, func(o *resourcegroupstaggingapi.GetResourcesPaginatorOptions) {
		o.Limit = 50
	})
	res, err := aws.GetResourcesTags(context.TODO(), paginator)
	if err != nil {
		log.Fatalf("unable to get resources tags, %v", err)
	}
	json, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("failed to return json data, %v", err)
	}
	log.Println(string(json))
}
