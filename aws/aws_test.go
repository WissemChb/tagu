package aws

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	rt "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	st "github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGetResourceTagPager struct {
	mock.Mock

	PageNumber int
	Pages      []*resourcegroupstaggingapi.GetResourcesOutput
}

type mockSTSAssumeRoleAPI struct {
	mock.Mock
}

func (s *mockSTSAssumeRoleAPI) AssumeRole(ctx context.Context, params *sts.AssumeRoleInput, optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error) {
	args := s.Called(ctx, params, optFns)
	if args.Get(0) != nil {
		return args.Get(0).(*sts.AssumeRoleOutput), args.Error(1)
	}
	return nil, args.Error(1)

}

func (m *mockGetResourceTagPager) HasMorePages() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockGetResourceTagPager) NextPage(ctx context.Context, optFns ...func(*resourcegroupstaggingapi.Options)) (output *resourcegroupstaggingapi.GetResourcesOutput, err error) {
	if m.PageNumber >= len(m.Pages) {
		return nil, fmt.Errorf("no more pages")
	}
	output = m.Pages[m.PageNumber]
	m.PageNumber++
	return output, nil
}

var ResourceTagsPagesOutput = []*resourcegroupstaggingapi.GetResourcesOutput{
	{
		ResourceTagMappingList: []rt.ResourceTagMapping{
			{
				ResourceARN: aws.String("arn:aws:ec2:us-east-1:123456789012:instance/i-12345678"),
				Tags: []rt.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-instance"),
					},
					{
						Key:   aws.String("Owner"),
						Value: aws.String("test-owner"),
					},
				},
			},
			{
				ResourceARN: aws.String("arn:aws:ec2:us-east-1:123456789012:instance/i-12345675"),
				Tags: []rt.Tag{
					{
						Key:   aws.String("ENV"),
						Value: aws.String("test-env"),
					},
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-instance2"),
					},
				},
			},
		},
	},
	{
		ResourceTagMappingList: []rt.ResourceTagMapping{
			{
				ResourceARN: aws.String("arn:aws:ec2:us-east-2:123456789012:instance/i-123456100"),
				Tags: []rt.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-instance-p2"),
					},
				},
			},
			{
				ResourceARN: aws.String("arn:aws:ec2:us-east-2:123456789012:instance/i-12345689"),
				Tags: []rt.Tag{
					{
						Key:   aws.String("ENV"),
						Value: aws.String("test-env-p2"),
					},
				},
			},
		},
	},
}

func TestGetResourceTagPagerSuite(t *testing.T) {
	cfg := aws.Config{
		Region: "us-east-2",
	}
	pager := &mockGetResourceTagPager{
		PageNumber: 0,
		Pages:      ResourceTagsPagesOutput,
	}
	assert := assert.New(t)
	creds := credentials.StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID:     "XXXXXXXXXXXXXXXXXXX",
			SecretAccessKey: "xxxxxxxXxxxxXXXXXXXx1455xxxxxxxxxxxxx",
			SessionToken:    "tokenxxxxxxxxxx",
		},
	}

	var fixtures = []struct {
		name     string
		err      error
		expected []RessourceTagResult
		count    int
	}{
		{
			"TestGetResourceTagPagerFunc",
			errors.New("no more pages"),
			[]RessourceTagResult{
				{Account: "123456789012", Region: "us-east-1", Service: "ec2", Resource: "instance/i-12345678", Key: "Name", Value: "test-instance"},
				{Account: "123456789012", Region: "us-east-1", Service: "ec2", Resource: "instance/i-12345678", Key: "Owner", Value: "test-owner"},
				{Account: "123456789012", Region: "us-east-1", Service: "ec2", Resource: "instance/i-12345675", Key: "ENV", Value: "test-env"},
				{Account: "123456789012", Region: "us-east-1", Service: "ec2", Resource: "instance/i-12345675", Key: "Name", Value: "test-instance2"},
				{Account: "123456789012", Region: "us-east-2", Service: "ec2", Resource: "instance/i-123456100", Key: "Name", Value: "test-instance-p2"},
				{Account: "123456789012", Region: "us-east-2", Service: "ec2", Resource: "instance/i-12345689", Key: "ENV", Value: "test-env-p2"},
			},
			6,
		},
	}
	for _, fixture := range fixtures {
		t.Run(fixture.name, func(t *testing.T) {
			tags := Tags{
				Account: "123456789012",
				Region:  "us-east-1",
			}
			var expectedErr error
			if fixture.err != nil {
				expectedErr = fixture.err
			} else {
				expectedErr = fmt.Errorf("no more pages")
			}
			pager.On("HasMorePages").Return(true)
			pager.On("NextPage", mock.Anything, mock.Anything).Return(nil, fixture.err)
			err := tags.getResourcesTags(context.TODO(), cfg, pager, creds, "eu-east-1")
			assert.EqualValues(fixture.expected, tags.Output)
			assert.Equal(tags.length, fixture.count)
			assert.Equal(err, expectedErr)
		})
	}
}

func TestSetupCredentialsSuite(t *testing.T) {
	cfg := aws.Config{
		Region: "us-east-2",
	}

	tags := Tags{
		Account:  "123456789012",
		Region:   "us-east-1",
		RoleName: "role-name",
	}

	var fixtures = []struct {
		name         string
		defaultCreds aws.CredentialsProvider
		input        *sts.AssumeRoleOutput
		expected     aws.CredentialsProvider
		err          error
	}{
		{
			"SetupCredential with STS Assume role OK",
			nil,
			&sts.AssumeRoleOutput{
				Credentials: &st.Credentials{
					AccessKeyId:     aws.String("XXXXXXXXXXXXXXXXXXX"),
					SecretAccessKey: aws.String("xxxxxxxXxxxxXXXXXXXx1455xxxxxxxxxxxxx"),
					SessionToken:    aws.String("tokenxxxxxxxxxx"),
				},
			},
			credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     "XXXXXXXXXXXXXXXXXXX",
					SecretAccessKey: "xxxxxxxXxxxxXXXXXXXx1455xxxxxxxxxxxxx",
					SessionToken:    "tokenxxxxxxxxxx",
				},
			},
			nil,
		},
		{
			"SetupCredential with STS Assume role KO",
			nil,
			nil,
			nil,
			errors.New("AssumeRole error"),
		},
		{
			"SetupCredential with Default Auth",
			credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     "XXXXXXXXXXXXXXXXXXX",
					SecretAccessKey: "xxxxxxxXxxxxXXXXXXXx1455xxxxxxxxxxxxx",
					SessionToken:    "tokenxxxxxxxxxx",
				},
			},
			nil,
			credentials.StaticCredentialsProvider{
				Value: aws.Credentials{
					AccessKeyID:     "XXXXXXXXXXXXXXXXXXX",
					SecretAccessKey: "xxxxxxxXxxxxXXXXXXXx1455xxxxxxxxxxxxx",
					SessionToken:    "tokenxxxxxxxxxx",
				},
			},
			nil,
		},
	}

	assert := assert.New(t)
	for _, fixture := range fixtures {
		t.Run(fixture.name, func(t *testing.T) {
			if fixture.defaultCreds != nil {
				cfg.Credentials = fixture.defaultCreds
				tags.RoleName = ""
			}
			stsMock := mockSTSAssumeRoleAPI{}
			stsMock.On("AssumeRole", mock.Anything, mock.Anything, mock.Anything).Return(fixture.input, fixture.err)
			result, err := tags.setupCredentials(context.TODO(), cfg, &stsMock)
			assert.EqualValues(fixture.expected, result)
			assert.Equal(err, fixture.err)
		})
	}

}

func TestSetupRegionSuite(t *testing.T) {
	cfg := aws.Config{
		Region: "us-east-2",
	}
	tags := Tags{
		Account: "123456789012",
		Region:  "",
	}
	var fixtures = []struct {
		name     string
		input    string
		expected string
	}{
		{
			"SetupRegion with region parameter",
			"us-west-2",
			"us-west-2",
		},
		{
			"SetupRegion with default region",
			"",
			"us-east-2",
		},
	}

	assert := assert.New(t)
	for _, fixture := range fixtures {
		t.Run(fixture.name, func(t *testing.T) {
			tags.Region = fixture.input
			result := tags.setupRegion(cfg)
			assert.EqualValues(fixture.expected, result)
		})
	}
}
