package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

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

// Tags is stuct that dedfines the AWS tags input and filter
// params:
// 		account: string
// 		roleName: string
// 		region: string
type Tags struct {
	Account  string
	Region   string
	RoleName string
	Output   []RessourceTagResult
	length   int
}

// GetResourcesTagsPager is the interface that defines the pagination logic
// Used to make the test API simple and easy to mock
type GetResourcesTagsPager interface {
	HasMorePages() bool
	NextPage(ctx context.Context, optFns ...func(*resourcegroupstaggingapi.Options)) (*resourcegroupstaggingapi.GetResourcesOutput, error)
}

// STSAssumeRoleAPI defines the interface for the AssumeRole function.
// We use this interface to test the function using a mocked service.
type STSAssumeRoleAPI interface {
	AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error)
}

// setupCredentials setup aws.Credentials from  STS if is enabled otherwise get the default credentials
// Args:
// 		ctx: context.Context
// 		cfg: aws.Config
// 		stsApi: STSAssumeRoleAPI
// Returns:
// 		aws.CredentialsProvider: The AWS Credential interface
func (t Tags) setupCredentials(ctx context.Context, cfg aws.Config, stsAPI STSAssumeRoleAPI) (creds aws.CredentialsProvider, err error) {
	if t.RoleName != "" {
		input := &sts.AssumeRoleInput{
			RoleArn:         aws.String("arn:aws:iam::" + t.Account + ":role/" + t.RoleName),
			RoleSessionName: aws.String("session-" + t.Account),
		}
		result, err := stsAPI.AssumeRole(ctx, input)
		if err != nil {
			return creds, err
		}
		creds = credentials.NewStaticCredentialsProvider(*result.Credentials.AccessKeyId, *result.Credentials.SecretAccessKey, *result.Credentials.SessionToken)
		return creds, nil
	}
	return cfg.Credentials, nil
}

// setupRegion aims to get the default AWS region if not specified in the params
// Args:
// 		cfg: aws.Config
// Return:
// 		string: the specified region
func (t Tags) setupRegion(cfg aws.Config) string {
	if t.Region != "" {
		return t.Region
	}
	return cfg.Region
}

// GetResourcesTags returns the tags for the given resources from AWS And fatterns the results in a single struct
// Args:
// 		ctx: context.Context
// 		creds: aws.CredentialsProvider
// 		cfg: aws.Config
// Returns:
// 		[]RessourceTagResult: if successful, the tags for the given resources
// 		error: if an error occurred
func (t *Tags) getResourcesTags(ctx context.Context, cfg aws.Config, paginator GetResourcesTagsPager, creds aws.CredentialsProvider, region string) error {
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx, func(opts *resourcegroupstaggingapi.Options) {
			opts.Credentials = creds
			opts.Region = region
		})

		if err != nil {
			return err
		}

		for _, item := range output.ResourceTagMappingList {
			for _, tag := range item.Tags {
				infos := strings.Split(*item.ResourceARN, ":")
				t.Output = append(
					t.Output,
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
		t.length = len(t.Output)
	}
	return nil
}

// Run executes the tagging logic
// It represent the entrypoint for AWS tags modules
func (t *Tags) Run() error {
	var creds aws.CredentialsProvider

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	// Using the Config value, create the ResourceGroupsTagging client
	rsclient := resourcegroupstaggingapi.NewFromConfig(cfg)

	params := &resourcegroupstaggingapi.GetResourcesInput{}

	paginator := resourcegroupstaggingapi.NewGetResourcesPaginator(rsclient, params, func(o *resourcegroupstaggingapi.GetResourcesPaginatorOptions) {
		o.Limit = 50
	})
	stsclient := sts.NewFromConfig(cfg)
	creds, err = t.setupCredentials(ctx, cfg, stsclient)
	if err != nil {
		return err
	}
	err = t.getResourcesTags(ctx, cfg, paginator, creds, t.setupRegion(cfg))
	if err != nil {
		return err
	}
	return nil
}
