// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/M1chlCZ
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"context"
)

// GetSnapshotSlot returns the highest slot that the node has a snapshot for.
//
// Deprecated: use GetHighestSnapshotSlot instead.
func (cl *Client) GetSnapshotSlot(ctx context.Context) (out uint64, err error) {
	err = cl.rpcClient.CallForInto(ctx, &out, "getSnapshotSlot", nil)
	return
}
