// Copyright 2024 CloudWeGo Authors
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

package service

import (
	"context"
	"testing"

	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"
)

// GO_ENV=dev go test -run TestEmptyCart_Run
func TestEmptyCart_Run(t *testing.T) {
	ctx := context.Background()
	s := NewEmptyCartService(ctx)

	// init req and assert value
	req := &cart.EmptyCartReq{
		UserId: 1,
	}
	resp, err := s.Run(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Errorf("unexpected nil response")
	}
}
