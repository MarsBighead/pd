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

package grpcutil

import (
	"context"
	"reflect"
	"testing"
)

func TestBuildForwardContext(t *testing.T) {
	type args struct {
		ctx  context.Context
		addr string
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildForwardContext(tt.args.ctx, tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildForwardContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
