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

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		state = tstate.New()
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {

		_ = kubernetes.DeleteNamespace(kubernetes.KubeClient, state.Namespace)
		return ctx, nil
	})
}

func anIngressResourceWithRewriteAnnotation(spec *godog.DocString) error {
	ns, err := kubernetes.NewNamespace(kubernetes.KubeClient)
	if err != nil {
		return err
	}

	state.Namespace = ns

	ingress, err := kubernetes.IngressFromManifest(state.Namespace, spec.Content)
	if err != nil {
		return err
	}

	err = kubernetes.DeploymentsFromIngress(kubernetes.KubeClient, ingress)
	if err != nil {
		return err
	}

	err = kubernetes.NewIngress(kubernetes.KubeClient, state.Namespace, ingress)
	if err != nil {
		return err
	}

	state.IngressName = ingress.GetName()

	return nil
}

func theIngressStatusShowsTheIPAddressOrFQDNWhereItIsExposed() error {
	ingress, err := kubernetes.WaitForIngressAddress(kubernetes.KubeClient, state.Namespace, state.IngressName)
	if err != nil {
		return err
	}

	state.IPOrFQDN = ingress

	time.Sleep(3 * time.Second)

	return err
}

func iSendARequestTo(method string, rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	return state.CaptureRoundTrip(method, "http", u.Host, u.Path, u.Query(), nil, false)

}

func theResponseStatusCodeMustBe(arg1 int) error {
	return state.AssertStatusCode(statusCode)
}
//遍历看最终结果
func theValueOfTheFieldInTheHeaderMustBe(key string, value string) error {
	return state.AssertRequestHeader(key,value)  
	// AssertRequestHeader返回一个错误，如果捕获的请求头不包含预期的headerKey，
//或者匹配的请求头值与期望的headerValue不匹配。
//如果headerValue字符串等于' * '，头值检查将被忽略。
}

//去找那个被删了的
func headersFieldAndItsValueMustBeDeleted(key string) error {
    err=state.AssertRequestHeader(key,value)  
	if err != nil {
		return fmt.Errorf("the key=%s not be deleted",key)
	}
	return nil
}