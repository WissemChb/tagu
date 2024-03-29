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
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func execute(t *testing.T, c *cobra.Command, args ...string) (string, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	c.SetOut(buf)
	c.SetErr(buf)
	c.SetArgs(args)

	err := c.Execute()
	s := strings.TrimSpace(buf.String())
	return s, err
}

func TestRootCmdSuite(t *testing.T) {
	assert := assert.New(t)

	// Get the current project
	absConfig, _ := filepath.Abs("../")

	root := &cobra.Command{Use: "tagu", RunE: rootCmdRunE}
	initRootFlags(root)

	fixtures := []struct {
		name     string
		args     []string
		env      string
		expected string
		err      error
	}{
		{
			name:     "Test root config no option specified",
			args:     []string{},
			env:      "",
			expected: "Error: Config File \"config\" Not Found in \"[]\"\nUsage:\n  tagu [flags]\n\nFlags:\n  -c, --config string   config file (default is $HOME/.tagu.yaml)\n  -h, --help            help for tagu",
			err:      errors.New("Config File \"config\" Not Found in \"[]\""),
		},
		{
			name:     "Test root config from CONFIG",
			args:     []string{},
			env:      absConfig + "/examples/input-tags.yaml",
			expected: "Load configuration file " + absConfig + "/examples/input-tags.yaml",
			err:      nil,
		},
		{
			name: "Test root config flag",
			args: []string{
				"-c",
				absConfig + "/examples/aws-tags.yaml",
			},
			expected: "Load configuration file " + absConfig + "/examples/aws-tags.yaml",
			err:      nil,
		},
	}

	for _, fixture := range fixtures {
		t.Run(fixture.name, func(t *testing.T) {
			if fixture.env != "" {
				os.Setenv("CONFIG", fixture.env)
				defer os.Unsetenv("CONFIG")
			}
			defer viper.Reset()
			res, err := execute(t, root, fixture.args...)
			assert.Equal(res, fixture.expected)
			assert.Equal(err, fixture.err)
		})
	}
}
