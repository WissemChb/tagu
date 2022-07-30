/*
Copyright © 2022 Wissem BEN CHAABANE<benchaaben.wissem@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"errors"
	"os"
	"testing"

	"tagu/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// func execute(t *testing.T, c *cobra.Command, args ...string) (string, error) {
// 	t.Helper()

// 	buf := new(bytes.Buffer)
// 	c.SetOut(buf)
// 	c.SetErr(buf)
// 	c.SetArgs(args)

// 	err := c.Execute()
// 	s := strings.TrimSpace(buf.String())
// 	return s, err
// }

func TestAwsCmdSuite(t *testing.T) {
	assert := assert.New(t)

	aws := &cobra.Command{Use: "aws", RunE: awsCmdRunE}
	initAwsFlags(aws)

	fixtures := []struct {
		name     string
		args     []string
		env      string
		expected string
		err      error
	}{
		{
			name:     "Test AWS config no option specified",
			args:     []string{},
			expected: "Error: Config File \"config\" Not Found in \"[]\"\nUsage:\n  aws [flags]\n\nFlags:\n  -h, --help                help for aws\n  -i, --input-file string   the input file",
			err:      errors.New("Config File \"config\" Not Found in \"[]\""),
		},
		{
			name:     "Test AWS config from AWS_CONFIG",
			args:     []string{},
			env:      "/workspaces/tagu/examples/input-tags.yaml",
			expected: "Load configuration file /workspaces/tagu/examples/input-tags.yaml",
			err:      nil,
		},
		{
			name: "Test AWS config flag",
			args: []string{
				"-i",
				"/workspaces/tagu/examples/aws-tags.yaml",
			},
			expected: "Load configuration file /workspaces/tagu/examples/aws-tags.yaml",
			err:      nil,
		},
	}

	for _, fixture := range fixtures {
		t.Run(fixture.name, func(t *testing.T) {
			if fixture.env != "" {
				os.Setenv("AWS_CONFIG", fixture.env)
				defer os.Unsetenv("AWS_CONFIG")
			}
			defer viper.Reset()
			res, err := execute(t, aws, fixture.args...)
			assert.Equal(res, fixture.expected)
			assert.Equal(err, fixture.err)
		})
	}
}

func TestLoadAwsConfigSuite(t *testing.T) {
	assert := assert.New(t)

	fixtures := []struct {
		name     string
		input    string
		expected *models.Spec
		err      error
	}{
		{
			name:  "Test Load config OK",
			input: "/workspaces/tagu/examples/aws-tags.yaml",
			expected: &models.Spec{
				RoleName: "test-role",
				FilterInput: []models.InputTag{
					{
						Account: "236534879095",
						Regions: []string{
							"us-west-1",
							"eu-west-1",
						},
						FilterResources: []string{
							"ec2:instance",
						},
						FilterTags: []models.Tags{
							{
								Key: "env",
								Values: []string{
									"prod",
								},
							},
						},
					},
					{
						Account: "636568979095",
						Regions: []string{
							"us-east-1",
							"eu-east-2",
						},
						FilterResources: []string{
							"rds",
							"s3",
							"lambda",
						},
						FilterTags: []models.Tags{
							{
								Key: "Schedule",
							},
						},
					},
					{
						Account: "436567879095",
						Regions: []string{
							"us-east-3",
							"eu-west-6",
						},
						FilterResources: []string{
							"ec2",
							"rds",
						},
						FilterTags: []models.Tags{
							{
								Values: []string{
									"dev",
								},
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			name:     "Load AWS config file not found",
			input:    "examples/aws-tags.yaml",
			expected: &models.Spec{},
			err:      errors.New("open examples/aws-tags.yaml: no such file or directory"),
		},
		{
			name:  "Build AWS config file OK",
			input: "/workspaces/tagu/examples/aws-general.yaml",
			expected: &models.Spec{
				RoleName: "test-role",
				FilterInput: []models.InputTag{
					{
						Account: "236534879095",
						Regions: []string{
							"us-west-1",
							"eu-west-1",
							"eu-west-6",
						},
						FilterResources: []string{
							"ec2:instance",
							"rds",
							"s3",
							"lambda",
						},
						FilterTags: []models.Tags{
							{
								Key: "env",
								Values: []string{
									"PROD",
								},
							},
						},
					},
					{
						Account: "636568979095",
						Regions: []string{
							"us-west-1",
							"eu-west-1",
							"eu-west-6",
						},
						FilterResources: []string{
							"ec2:instance",
							"rds",
							"s3",
							"lambda",
						},
						FilterTags: []models.Tags{
							{
								Key: "env",
								Values: []string{
									"PROD",
								},
							},
						},
					},
					{
						Account: "436567879095",
						Regions: []string{
							"us-west-1",
							"eu-west-1",
							"eu-west-6",
						},
						FilterResources: []string{
							"ec2:instance",
							"rds",
							"s3",
							"lambda",
						},
						FilterTags: []models.Tags{
							{
								Key: "env",
								Values: []string{
									"PROD",
								},
							},
						},
					},
				},
			},
			err: nil,
		},
	}

	for _, fixture := range fixtures {
		t.Run(fixture.name, func(t *testing.T) {
			defer viper.Reset()
			out, err := awsloadConfig(fixture.input)
			assert.Equal(&out, fixture.expected)
			assert.Equal(err, fixture.err)
		})
	}
}
