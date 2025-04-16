// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package codeowners

import (
	"fmt"
	"os"
	"strings"

	"github.com/hmarr/codeowners"
)

var (
	ErrNoOwners = fmt.Errorf("no owners defined")
)

type Ruleset struct {
	codeowners.Ruleset

	// prefix is removed from all 'Match()' calls, as the CODEOWNERS files
	// typically define relative paths within a repository, but not the
	// full repository path.
	prefix string

	exclude map[string]struct{}
}

// Load code owners from 'paths', ignoring
func Load(paths []string, ignorePrefix string) (*Ruleset, error) {
	var allOwners codeowners.Ruleset

	if len(paths) == 0 {
		owners, err := codeowners.LoadFileFromStandardLocation()
		if err != nil {
			return nil, fmt.Errorf("while loading: %w", err)
		}
		allOwners = owners
	}

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

	return &Ruleset{
		allOwners,
		ignorePrefix,
		make(map[string]struct{}),
	}, nil
}

func (r *Ruleset) WithExcludedOwners(excludedOwners []string) *Ruleset {
	excluded := make(map[string]struct{})
	for _, o := range excludedOwners {
		excluded[o] = struct{}{}
	}
	r.exclude = excluded
	return r
}

func (r *Ruleset) Match(relPath string) ([]string, error) {
	path := strings.TrimPrefix(relPath, r.prefix) + "/"
	rule, err := r.Ruleset.Match(path)
	if err == nil && (rule == nil || rule.Owners == nil) {
		err = ErrNoOwners
	}
	if err != nil {
		return nil, err
	}
	owners := make([]string, 0, len(rule.Owners))
	for _, o := range rule.Owners {
		if _, ok := r.exclude[o.String()]; ok {
			continue
		}
		owners = append(owners, o.String())
	}
	return owners, nil
}
