package aws

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGetResourceTagPager struct {
	mock.Mock

	PageNumber int
	Pages      []*resourcegroupstaggingapi.GetResourcesOutput
}

func (m *mockGetResourceTagPager) HasMorePages() bool {
	args := m.Called()
	return args.Bool(0)
	// return m.PageNumber < len(m.Pages)
}

func (m *mockGetResourceTagPager) NextPage(ctx context.Context, optFns ...func(*resourcegroupstaggingapi.Options)) (output *resourcegroupstaggingapi.GetResourcesOutput, err error) {

	if m.PageNumber >= len(m.Pages) {
		return nil, fmt.Errorf("no more pages")
	}
	output = m.Pages[m.PageNumber]
	m.PageNumber++
	return output, nil
}

func TestGetResourceTagPagerFunc(t *testing.T) {
	pager := &mockGetResourceTagPager{
		PageNumber: 0,
		Pages: []*resourcegroupstaggingapi.GetResourcesOutput{
			{
				ResourceTagMappingList: []types.ResourceTagMapping{
					{
						ResourceARN: aws.String("arn:aws:ec2:us-east-1:123456789012:instance/i-12345678"),
						Tags: []types.Tag{
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
						Tags: []types.Tag{
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
				ResourceTagMappingList: []types.ResourceTagMapping{
					{
						ResourceARN: aws.String("arn:aws:ec2:us-east-2:123456789012:instance/i-123456100"),
						Tags: []types.Tag{
							{
								Key:   aws.String("Name"),
								Value: aws.String("test-instance-p2"),
							},
						},
					},
					{
						ResourceARN: aws.String("arn:aws:ec2:us-east-2:123456789012:instance/i-12345689"),
						Tags: []types.Tag{
							{
								Key:   aws.String("ENV"),
								Value: aws.String("test-env-p2"),
							},
						},
					},
				},
			},
		},
	}
	assert := assert.New(t)

	t.Run("TestGetResourceTagPagerFunc OK", func(t *testing.T) {
		pager.On("HasMorePages").Return(true)
		result, err := GetResourcesTags(context.TODO(), pager)
		assert.Errorf(err, "no more page")
		assert.EqualValues([]RessourceTagResult{
			{Account: "123456789012", Region: "us-east-1", Service: "ec2", Resource: "instance/i-12345678", Key: "Name", Value: "test-instance"},
			{Account: "123456789012", Region: "us-east-1", Service: "ec2", Resource: "instance/i-12345678", Key: "Owner", Value: "test-owner"},
			{Account: "123456789012", Region: "us-east-1", Service: "ec2", Resource: "instance/i-12345675", Key: "ENV", Value: "test-env"},
			{Account: "123456789012", Region: "us-east-1", Service: "ec2", Resource: "instance/i-12345675", Key: "Name", Value: "test-instance2"},
			{Account: "123456789012", Region: "us-east-2", Service: "ec2", Resource: "instance/i-123456100", Key: "Name", Value: "test-instance-p2"},
			{Account: "123456789012", Region: "us-east-2", Service: "ec2", Resource: "instance/i-12345689", Key: "ENV", Value: "test-env-p2"},
		}, result)
	})

	t.Run("TestGetResourceTagPagerFunc KO", func(t *testing.T) {
		pager.On("HasMorePages").Return(true)
		pager.On("NextPage", mock.Anything, mock.Anything).Return(nil, errors.New("no more pages"))
		result, err := GetResourcesTags(context.TODO(), pager)
		assert.Nil(result)
		assert.Error(err)
		assert.Equal(err.Error(), "no more pages")
	})
}
