/*
Copyright 2021 The BFE Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package header

import (
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v16"

	"github.com/bfenetworks/ingress-bfe/test/e2e/pkg/kubernetes"
	tstate "github.com/bfenetworks/ingress-bfe/test/e2e/pkg/state"
)

var (
	state *tstate.Scenario
)

// IMPORTANT: Steps definitions are generated and should not be modified
// by hand but rather through make codegen. DO NOT EDIT.

// InitializeScenario configures the Feature to test
func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^an Ingress resource with rewrite annotation$`, anIngressResourceWithRewriteAnnotation)
	ctx.Step(`^The Ingress status shows the IP address or FQDN where it is exposed$`, theIngressStatusShowsTheIPAddressOrFQDNWhereItIsExposed)
	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)"$`, iSendARequestTo)
	ctx.Step(`^the response status code must be (\d+)$`, theResponseStatusCodeMustBe)
	ctx.Step(`^the value of the "([^"]*)" field in the header must be "([^"]*)"$`, theValueOfTheFieldInTheHeaderMustBe)
	ctx.Step(`^header\'s "([^"]*)" field and its value must be deleted$`, headersFieldAndItsValueMustBeDeleted)

	ctx.BeforeScenario(func(*godog.Scenario) {
		state = tstate.New()
	})

	ctx.AfterScenario(func(*messages.Pickle, error) {
		// delete namespace an all the content
		_ = kubernetes.DeleteNamespace(kubernetes.KubeClient, state.Namespace)
	})
}

func anIngressResourceWithRewriteAnnotation(arg1 *godog.DocString) error {
	return godog.ErrPending
}

func theIngressStatusShowsTheIPAddressOrFQDNWhereItIsExposed() error {
	return godog.ErrPending
}

func iSendARequestTo(arg1 string, arg2 string) error {
	return godog.ErrPending
}

func theResponseStatusCodeMustBe(arg1 int) error {
	return godog.ErrPending
}

func theValueOfTheFieldInTheHeaderMustBe(arg1 string, arg2 string) error {
	return godog.ErrPending
}

func headersFieldAndItsValueMustBeDeleted(arg1 string) error {
	return godog.ErrPending
}

