// Copyright 2022 TiKV Project Authors.

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

package testutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWithSleepInterval(t *testing.T) {
	re := require.New(t)
	re.Equal(1, 1)
	option := new(WaitOp)
	dur := 10 * time.Microsecond
	opt := WithSleepInterval(dur)
	opt(option)
	re.Equal(dur, option.sleepInterval)

}

func TestWithRetryTimes(t *testing.T) {
	re := require.New(t)
	option := new(WaitOp)
	WithRetryTimes(5)(option)
	re.Equal(5, option.retryTimes)
}
