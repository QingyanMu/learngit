// Copyright (c) 2021 The BFE Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//Package header is the module of rewrite url.
//This file defines header & cache's struct, also implements update ingress method.
package header

import (
	netv1 "k8s.io/api/networking/v1"

	"github.com/bfenetworks/ingress-bfe/internal/bfeConfig/annotations"
	"github.com/bfenetworks/ingress-bfe/internal/bfeConfig/configs/cache"
	"github.com/bfenetworks/ingress-bfe/internal/bfeConfig/util"
	"github.com/bfenetworks/bfe/bfe_modules/mod_header"
)

type headerRule struct {
	*cache.BaseRule

	// action is the header action
	actions *mod_header.ActionFileList

	// if true, stop processing the next rule
	last    bool
}

type headerRuleCache struct {
	*cache.BaseCache
}

func newHeaderRuleCache(version string) *headerRuleCache {
	return &headerRuleCache{
		BaseCache: cache.NewBaseCache(version),
	}
}

func (c headerRuleCache) UpdateByIngress(ingress *netv1.Ingress) error {
	headerActions, err := annotations.GetHeaderAction(ingress.Annotations)
	if err != nil {
		return err
	}

	if headerActions == nil {
		return nil
	}

	return c.BaseCache.UpdateByIngressFramework(
		ingress,
		func(ingress *netv1.Ingress, host, path string, _ netv1.HTTPIngressPath) (cache.Rule, error) {
			// preCheck. Duplicate with mod_header.ActionFileListCheck?
			if err := checkAction(headerActions); err != nil {
				return nil, err
			}

			//acts := []mod_header.ActionFile{}
			acts := make(mod_header.ActionFileList, 0)
			for _, action := range *headerActions {
				actionFile := mod_header.ActionFile{
					Cmd: &action.Cmd,
					Params: action.Params,
				}
				acts = append(acts, actionFile)
			}
			actions := &acts
			if err = mod_header.ActionFileListCheck(actions); err != nil {
				return nil, err
			}
			return &headerRule{
				BaseRule: cache.NewBaseRule(
					util.NamespacedName(ingress.Namespace, ingress.Name),
					host,
					path,
					ingress.Annotations,
					ingress.CreationTimestamp.Time,
				),
				actions:    actions,
				last:       true,				// Last = true as default??
			}, nil
		},
		nil,
		nil,
	)

	return nil
}
