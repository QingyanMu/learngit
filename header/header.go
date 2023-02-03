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

//Package header is the module of header url.
//This file implements operate rule cache, generate and reload config file methods.
package header

import (
	"fmt"

	netv1 "k8s.io/api/networking/v1"

	"github.com/bfenetworks/bfe/bfe_modules/mod_header"
	"github.com/bfenetworks/ingress-bfe/internal/bfeConfig/annotations"
	"github.com/bfenetworks/ingress-bfe/internal/bfeConfig/configs"
	"github.com/bfenetworks/ingress-bfe/internal/bfeConfig/util"
)

const (
	ConfigNameHeader  = "mod_header"
	RuleData          = "mod_header/header.data"
)

type ModHeaderConfig struct {
	version          string
	headerRuleCache *headerRuleCache
	headerConfFile  *mod_header.HeaderConfFile
}

func NewHeaderConfig(version string) *ModHeaderConfig {
	return &ModHeaderConfig{
		version:          version,
		headerRuleCache:  newHeaderRuleCache(version),
		headerConfFile:   newHeaderConfFile(version),
	}
}

func newHeaderConfFile(version string) *mod_header.HeaderConfFile {
	ruleFileList := make([]mod_header.HeaderRuleFile, 0)
	productRulesFile := make(mod_header.ProductRulesFile)
	productRulesFile[configs.DefaultProduct] = (*mod_header.RuleFileList)(&ruleFileList)
	return &mod_header.HeaderConfFile{
		Version: &version,
		Config:  &productRulesFile,
	}
}

func (c *ModHeaderConfig) Name() string {
	return ConfigNameHeader
}


func (c *ModHeaderConfig) UpdateIngress(ingress *netv1.Ingress) error {
	// clear cache
	ingressName := util.NamespacedName(ingress.Namespace, ingress.Name)
	if c.headerRuleCache.ContainsIngress(ingressName) {
		c.headerRuleCache.DeleteByIngress(ingressName)
	}
	// nothing to update
	if len(ingress.Spec.Rules) == 0 {
		return nil
	}
	//update the cache
	return c.headerRuleCache.UpdateByIngress(ingress)
}

func (c *ModHeaderConfig) DeleteIngress(namespace, name string) {
	ingressName := util.NamespacedName(namespace, name)
	if !c.headerRuleCache.ContainsIngress(ingressName) {
		return
	}

	c.headerRuleCache.DeleteByIngress(ingressName)
}

//bfe process to reload headerConfFile through bfe monitor port
func (c *ModHeaderConfig) Reload() error {
	if err := c.updateHeaderConf(); err != nil {
		return fmt.Errorf("update %s config error: %v", RuleData, err)
	}

	if *c.headerConfFile.Version != c.version {
		// dump config file
		err := util.DumpBfeConf(RuleData, c.headerConfFile)
		if err != nil {
			return fmt.Errorf("dump %s error: %v", RuleData, err)
		}
		// reload bfe engine
		err = util.ReloadBfe(ConfigNameHeader)
		if err != nil {
			return err
		}
		c.version = *c.headerConfFile.Version
	}

	return nil
}

func (c *ModHeaderConfig) updateHeaderConf() error {

	if *c.headerConfFile.Version == c.headerRuleCache.Version {
		// if the version is the same, no need to update
		return nil
	}

	//try to parse the rule Cond ,Actions,and Last from the ingress annotations
	ruleList := c.headerRuleCache.GetRules()
	headerRuleList := make(mod_header.RuleFileList, 0, len(ruleList))
	for _, rule := range ruleList {
		rule := rule.(*headerRule)
		cond, err := rule.GetCond()
		if err != nil {
			return err
		}
		headerRuleList = append(headerRuleList, mod_header.HeaderRuleFile{
			Cond:    &cond,
			Actions: rule.actions,
			Last:    &(rule.last),
		})
	}

	headerConfFile := newHeaderConfFile(c.headerRuleCache.Version)
	(*headerConfFile.Config)[configs.DefaultProduct] = &headerRuleList

	if err := mod_header.HeaderConfCheck(*headerConfFile); err != nil {
		return err
	}

	c.headerConfFile = headerConfFile
	return nil
}


//checkAction checks the add, set, and del directives for the header and returns an error if the action is invalid
func checkAction(actions *annotations.HeaderActions) error {
	if actions == nil {
		return fmt.Errorf("No cmd for header action")
	}
	for _, action := range *actions {
		switch action.Cmd {
		case annotations.HeaderActionReqHeaderAdd:
			fallthrough
	
		case annotations.HeaderActionReqHeaderSet:
			if len(action.Params) != 2 {
				return fmt.Errorf("Invalid parameters for header action: %s", action.Cmd)
			}
	
		case annotations.HeaderActionReqHeaderDel:
			if len(action.Params) != 1 {
				return fmt.Errorf("Invalid parameters for header action: %s", action.Cmd)
			}
		
		case annotations.HeaderActionRspHeaderAdd:
			fallthrough

		case annotations.HeaderActionRspHeaderSet:
			if len(action.Params) != 2 {
				return fmt.Errorf("Invalid parameters for header action: %s", action.Cmd)
			}

		case annotations.HeaderActionRspHeaderDel:
			if len(action.Params) != 1 {
				return fmt.Errorf("Invalid parameters for header action: %s", action.Cmd)
			}
		
		default:
			return fmt.Errorf("unsupported cmd for header action: %s", action.Cmd)
		}
	}
	return nil
}
