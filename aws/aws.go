package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
)

// GetResourcesTagsPager is the interface that defines the pagination logic
// Used to make the test API simple and easy to mock
type GetResourcesTagsPager interface {
	HasMorePages() bool
	NextPage(ctx context.Context, optFns ...func(*resourcegroupstaggingapi.Options)) (*resourcegroupstaggingapi.GetResourcesOutput, error)
}

// RessourceTagResult is the output of a GetResourcesTags call
// It contains the account, region, service, resource, key and value of the tag
// flattened in a single struct
type RessourceTagResult struct {
	Account  string
	Region   string
	Service  string
	Resource string
	Key      string
	Value    string
}

// GetResourcesTags returns the tags for the given resources from AWS
// And fatterns the results in a single struct
func GetResourcesTags(ctx context.Context, paginator GetResourcesTagsPager) (result []RessourceTagResult, err error) {

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return result, err
		}

		for _, item := range output.ResourceTagMappingList {
			for _, tag := range item.Tags {
				infos := strings.Split(*item.ResourceARN, ":")
				result = append(
					result,
					RessourceTagResult{
						Account:  infos[4],
						Region:   infos[3],
						Service:  infos[2],
						Resource: infos[5],
						Key:      *tag.Key,
						Value:    *tag.Value,
					},
				)
			}
		}
	}
	return result, nil
}
