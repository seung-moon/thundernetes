package controllers

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DynamicStandby controller tests", func() {
	Context("testing dynamic standby", func() {
		ctx := context.Background()

		// sample test
		It("this is a sample test", func() {
			buildName, buildID := getNewBuildNameAndID()
			gsb := testGenerateGameServerBuild(buildName, testnamespace, buildID, 2, 4, false)
			Expect(testk8sClient.Create(ctx, &gsb)).Should(Succeed())
		})
	})
})
