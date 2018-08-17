package kube

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/jenkins-x/jx/pkg/gits"
	"github.com/jenkins-x/jx/pkg/tests"
	"github.com/jenkins-x/jx/pkg/util"
	"github.com/stretchr/testify/assert"
)

// TODO refactor to encapsulate
func TestEnvModifyNamespace(t *testing.T) {
	if os.Getenv("RUN_UNENCAPSULATED_TESTS") == "true" {
		tempDir, err := ioutil.TempDir("", "test-env-modify-namespace")
		assert.NoError(t, err)

		testData := path.Join("test_data", "env_modify_namespace")
		_, err = os.Stat(testData)
		assert.NoError(t, err)

		files, err := ioutil.ReadDir(testData)
		assert.NoError(t, err)

		for _, f := range files {
			if !f.IsDir() {
				name := f.Name()
				srcDir := filepath.Join(testData, name)
				testDir := filepath.Join(tempDir, name)
				util.CopyFile(srcDir, testDir)

			}
		}

		git := gits.NewGitCLI()
		err = git.Init(tempDir)
		assert.NoError(t, err)

		testNs := "jx-staging"

		env := NewPermanentEnvironment("jx")
		env.Spec.Namespace = testNs

		err = modifyNamespace(os.Stdout, tempDir, env, git)
		assert.NoError(t, err)

		tests.AssertFileContains(t, filepath.Join(tempDir, "Makefile"), `NAMESPACE := "`+testNs+`"`)
		tests.AssertFileContains(t, filepath.Join(tempDir, "Jenkinsfile"), `DEPLOY_NAMESPACE = "`+testNs+`"`)
	} else {
		t.Skip("skipping TestEnvModifyNamespace; RUN_UNENCAPSULATED_TESTS not set")
	}
}
