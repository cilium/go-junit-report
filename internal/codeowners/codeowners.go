// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package codeowners

import (
	"fmt"
	"os"

	"github.com/hmarr/codeowners"
)

var (
	ErrNoOwners = fmt.Errorf("no owners defined")
)

type Ruleset struct {
	codeowners.Ruleset
}

func Load(paths []string) (*Ruleset, error) {
	if len(paths) == 0 {
		owners, err := codeowners.LoadFileFromStandardLocation()
		if err != nil {
			return nil, fmt.Errorf("while loading: %w", err)
		}
		return &Ruleset{owners}, nil
	}

	var allOwners codeowners.Ruleset

	for _, f := range paths {
		coFile, err := os.Open(f)
		if err != nil {
			return nil, fmt.Errorf("while opening %s: %w", f, err)
		}
		defer coFile.Close()

		owners, err := codeowners.ParseFile(coFile)
		if err != nil {
			return nil, fmt.Errorf("while parsing %s: %w", f, err)
		}

		allOwners = append(allOwners, owners...)
	}

	return &Ruleset{allOwners}, nil
}

func (r *Ruleset) Match(relPath string) ([]string, error) {
	rule, err := r.Ruleset.Match(relPath)
	if err == nil && (rule == nil || rule.Owners == nil) {
		err = ErrNoOwners
	}
	if err != nil {
		return nil, err
	}
	owners := make([]string, 0, len(rule.Owners))
	for _, o := range rule.Owners {
		owners = append(owners, o.String())
	}
	return owners, nil
}
