//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2022 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

package backup

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExcludeClasses(t *testing.T) {
	tests := []struct {
		in  BackupDescriptor
		xs  []string
		out []string
	}{
		{in: BackupDescriptor{}, xs: []string{}, out: []string{}},
		{in: BackupDescriptor{Classes: []ClassDescriptor{{Name: "a"}}}, xs: []string{}, out: []string{"a"}},
		{in: BackupDescriptor{Classes: []ClassDescriptor{{Name: "a"}}}, xs: []string{"a"}, out: []string{}},
		{in: BackupDescriptor{Classes: []ClassDescriptor{{Name: "1"}, {Name: "2"}, {Name: "3"}, {Name: "4"}}}, xs: []string{"2", "3"}, out: []string{"1", "4"}},
		{in: BackupDescriptor{Classes: []ClassDescriptor{{Name: "1"}, {Name: "2"}, {Name: "3"}}}, xs: []string{"1", "3"}, out: []string{"2"}},

		// {in: []BackupDescriptor{"1", "2", "3", "4"}, xs: []string{"2", "3"}, out: []string{"1", "4"}},
		// {in: []BackupDescriptor{"1", "2", "3"}, xs: []string{"1", "3"}, out: []string{"2"}},
	}
	for _, tc := range tests {
		tc.in.Exclude(tc.xs)
		lst := tc.in.List()
		assert.Equal(t, tc.out, lst)
	}
}

func TestIncludeClasses(t *testing.T) {
	tests := []struct {
		in  BackupDescriptor
		xs  []string
		out []string
	}{
		{in: BackupDescriptor{}, xs: []string{}, out: []string{}},
		{in: BackupDescriptor{Classes: []ClassDescriptor{{Name: "a"}}}, xs: []string{}, out: []string{"a"}},
		{in: BackupDescriptor{Classes: []ClassDescriptor{{Name: "a"}}}, xs: []string{"a"}, out: []string{"a"}},
		{in: BackupDescriptor{Classes: []ClassDescriptor{{Name: "1"}, {Name: "2"}, {Name: "3"}, {Name: "4"}}}, xs: []string{"2", "3"}, out: []string{"2", "3"}},
		{in: BackupDescriptor{Classes: []ClassDescriptor{{Name: "1"}, {Name: "2"}, {Name: "3"}}}, xs: []string{"1", "3"}, out: []string{"1", "3"}},
	}
	for _, tc := range tests {
		tc.in.Include(tc.xs)
		lst := tc.in.List()
		assert.Equal(t, tc.out, lst)
	}
}

func TestAllExist(t *testing.T) {
	x := BackupDescriptor{Classes: []ClassDescriptor{{Name: "a"}}}
	if y := x.AllExists(nil); y != "" {
		t.Errorf("x.AllExists(nil) got=%v want=%v", y, "")
	}
	if y := x.AllExists([]string{"a"}); y != "" {
		t.Errorf("x.AllExists(['a']) got=%v want=%v", y, "")
	}
	if y := x.AllExists([]string{"b"}); y != "b" {
		t.Errorf("x.AllExists(['a']) got=%v want=%v", y, "b")
	}
}

func TestValidateBackup(t *testing.T) {
	timept := time.Now().UTC()
	bytes := []byte("hello")
	tests := []struct {
		desc    BackupDescriptor
		success bool
	}{
		// first level check
		{desc: BackupDescriptor{}},
		{desc: BackupDescriptor{ID: "1"}},
		{desc: BackupDescriptor{ID: "1", Version: "1"}},
		{desc: BackupDescriptor{ID: "1", Version: "1", ServerVersion: "1"}},
		{desc: BackupDescriptor{ID: "1", Version: "1", ServerVersion: "1", StartedAt: timept}, success: true},
		{desc: BackupDescriptor{ID: "1", Version: "1", ServerVersion: "1", StartedAt: timept, Error: "err"}},
		{desc: BackupDescriptor{
			ID: "1", Version: "1", ServerVersion: "1", StartedAt: timept,
			Classes: []ClassDescriptor{{}},
		}},
		{desc: BackupDescriptor{
			ID: "1", Version: "1", ServerVersion: "1", StartedAt: timept,
			Classes: []ClassDescriptor{{Name: "n"}},
		}},
		{desc: BackupDescriptor{
			ID: "1", Version: "1", ServerVersion: "1", StartedAt: timept,
			Classes: []ClassDescriptor{{Name: "n", Schema: bytes}},
		}},
		{desc: BackupDescriptor{
			ID: "1", Version: "1", ServerVersion: "1", StartedAt: timept,
			Classes: []ClassDescriptor{{Name: "n", Schema: bytes, ShardingState: bytes}},
		}, success: true},
		{desc: BackupDescriptor{
			ID: "1", Version: "1", ServerVersion: "1", StartedAt: timept,
			Classes: []ClassDescriptor{{
				Name: "n", Schema: bytes, ShardingState: bytes,
				Shards: []ShardDescriptor{{Name: ""}},
			}},
		}},
		{desc: BackupDescriptor{
			ID: "1", Version: "1", ServerVersion: "1", StartedAt: timept,
			Classes: []ClassDescriptor{{
				Name: "n", Schema: bytes, ShardingState: bytes,
				Shards: []ShardDescriptor{{Name: "n", Node: ""}},
			}},
		}},
		{desc: BackupDescriptor{
			ID: "1", Version: "1", ServerVersion: "1", StartedAt: timept,
			Classes: []ClassDescriptor{{
				Name: "n", Schema: bytes, ShardingState: bytes,
				Shards: []ShardDescriptor{{Name: "n", Node: "n"}},
			}},
		}},
		{desc: BackupDescriptor{
			ID: "1", Version: "1", ServerVersion: "1", StartedAt: timept,
			Classes: []ClassDescriptor{{
				Name: "n", Schema: bytes, ShardingState: bytes,
				Shards: []ShardDescriptor{{
					Name: "n", Node: "n",
					PropLengthTrackerPath: "n", DocIDCounterPath: "n", ShardVersionPath: "n",
				}},
			}},
		}, success: true},
		{desc: BackupDescriptor{
			ID: "1", Version: "1", ServerVersion: "1", StartedAt: timept,
			Classes: []ClassDescriptor{{
				Name: "n", Schema: bytes, ShardingState: bytes,
				Shards: []ShardDescriptor{{
					Name: "n", Node: "n",
					PropLengthTrackerPath: "n", DocIDCounterPath: "n", ShardVersionPath: "n",
					Files: []string{"file"},
				}},
			}},
		}},
		{desc: BackupDescriptor{
			ID: "1", Version: "1", ServerVersion: "1", StartedAt: timept,
			Classes: []ClassDescriptor{{
				Name: "n", Schema: bytes, ShardingState: bytes,
				Shards: []ShardDescriptor{{
					Name: "n", Node: "n",
					PropLengthTrackerPath: "n", DocIDCounterPath: "n", ShardVersionPath: "n",
					DocIDCounter: bytes, Files: []string{"file"},
				}},
			}},
		}},
		{desc: BackupDescriptor{
			ID: "1", Version: "1", ServerVersion: "1", StartedAt: timept,
			Classes: []ClassDescriptor{{
				Name: "n", Schema: bytes, ShardingState: bytes,
				Shards: []ShardDescriptor{{
					Name: "n", Node: "n",
					PropLengthTrackerPath: "n", DocIDCounterPath: "n", ShardVersionPath: "n",
					DocIDCounter: bytes, Version: bytes, PropLengthTracker: bytes, Files: []string{""},
				}},
			}},
		}},
		{desc: BackupDescriptor{
			ID: "1", Version: "1", ServerVersion: "1", StartedAt: timept,
			Classes: []ClassDescriptor{{
				Name: "n", Schema: bytes, ShardingState: bytes,
				Shards: []ShardDescriptor{{
					Name: "n", Node: "n",
					PropLengthTrackerPath: "n", DocIDCounterPath: "n", ShardVersionPath: "n",
					DocIDCounter: bytes, Version: bytes, PropLengthTracker: bytes, Files: []string{"file"},
				}},
			}},
		}, success: true},
	}
	for i, tc := range tests {
		err := tc.desc.Validate()
		if got := err == nil; got != tc.success {
			t.Errorf("%d. validate(%+v): want=%v got=%v err=%v", i, tc.desc, tc.success, got, err)
		}
	}
}