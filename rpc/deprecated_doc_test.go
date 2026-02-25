// Copyright 2026 github.com/M1chlCZ
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
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeprecatedRPCMethodsHaveGoDocTag(t *testing.T) {
	expectedDeprecated := map[string]bool{
		"GetConfirmedBlock":                 true,
		"GetConfirmedBlockWithOpts":         true,
		"GetConfirmedBlocks":                true,
		"GetConfirmedBlocksWithLimit":       true,
		"GetConfirmedSignaturesForAddress2": true,
		"GetConfirmedTransaction":           true,
		"GetConfirmedTransactionWithOpts":   true,
		"GetFeeCalculatorForBlockhash":      true,
		"GetFeeRateGovernor":                true,
		"GetFees":                           true,
		"GetRecentBlockhash":                true,
		"GetSnapshotSlot":                   true,
		"GetStakeActivation":                true,
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", nil, parser.ParseComments)
	require.NoError(t, err)

	seen := map[string]bool{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Recv == nil || fn.Name == nil {
					continue
				}
				if len(fn.Recv.List) != 1 {
					continue
				}

				star, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
				if !ok {
					continue
				}
				ident, ok := star.X.(*ast.Ident)
				if !ok || ident.Name != "Client" {
					continue
				}

				if !expectedDeprecated[fn.Name.Name] {
					continue
				}
				seen[fn.Name.Name] = true

				doc := ""
				if fn.Doc != nil {
					doc = fn.Doc.Text()
				}
				require.Containsf(t, doc, "Deprecated:", "%s should have a Go deprecation doc tag", fn.Name.Name)
			}
		}
	}

	for name := range expectedDeprecated {
		require.Truef(t, seen[name], "expected method %s not found while checking deprecations", name)
	}
}
