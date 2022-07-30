package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildInputFilerSuite(t *testing.T) {
	assert := assert.New(t)

	fixtures := []struct {
		name     string
		input    GeneralSpec
		expected Spec
		err      error
	}{
		{
			name: "Test input Filter Builder OK",
			input: GeneralSpec{
				RoleName: "test-role",
				Accounts: []string{
					"236534879095",
				},
				Regions: []string{
					"eu-west-1",
				},
				FilterResources: []string{
					"ec2:instance",
				},
				FilterTags: []Tags{
					{
						Key: "env",
						Values: []string{
							"PROD",
						},
					},
				},
			},
			expected: Spec{
				RoleName: "test-role",
				FilterInput: []InputTag{
					{
						Account: "236534879095",
						Regions: []string{
							"eu-west-1",
						},
						FilterResources: []string{
							"ec2:instance",
						},
						FilterTags: []Tags{
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

	for _, fx := range fixtures {
		t.Run(fx.name, func(t *testing.T) {
			spec := new(Spec)
			spec.UniformConfig(fx.input)
			assert.Equal(&fx.expected, spec)
		})
	}
}
