package models

type Tags struct {
	Key    string   `mapstructure:"key,omitempty"`
	Values []string `mapstructure:"values,omitempty"`
}

type InputTag struct {
	Account         string   `mapstructure:"account"`
	Regions         []string `mapstructure:"regions"`
	FilterResources []string `mapstructure:"resources,omitempty"`
	FilterTags      []Tags   `mapstructure:"filter-tags,omitempty"`
}

// Spec is the input data passed to
// the AWS resourcegroupstaggingapi to filter fetched resources
type Spec struct {
	RoleName    string     `mapstructure:"role-name"`
	FilterInput []InputTag `mapstructure:"filter-input"`
}

type GeneralSpec struct {
	RoleName        string   `mapstructure:"role-name"`
	Accounts        []string `mapstructure:"accounts"`
	Regions         []string `mapstructure:"regions"`
	FilterResources []string `mapstructure:"resources,omitempty"`
	FilterTags      []Tags   `mapstructure:"filter-tags,omitempty"`
}

// UniformConfig is function that take general input specification
// and convert it to the Detailed Spec
func (spec *Spec) UniformConfig(gSpec GeneralSpec) {
	spec.RoleName = gSpec.RoleName
	for _, acc := range gSpec.Accounts {
		spec.FilterInput = append(spec.FilterInput, InputTag{
			Account:         acc,
			Regions:         gSpec.Regions,
			FilterResources: gSpec.FilterResources,
			FilterTags:      gSpec.FilterTags,
		})
	}
}
