// Copyright 2021 github.com/gagliardetto
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

package token2022

import (
	"bytes"
	"strconv"
	"testing"

	solana "github.com/M1chlCZ/solana-go"
	ag_gofuzz "github.com/gagliardetto/gofuzz"
	ag_require "github.com/stretchr/testify/require"
)

func TestEncodeDecode_InitializeMint(t *testing.T) {
	fu := ag_gofuzz.New().NilChance(0)
	for i := 0; i < 1; i++ {
		t.Run("InitializeMint"+strconv.Itoa(i), func(t *testing.T) {
			{
				params := new(InitializeMint)
				fu.Fuzz(params)
				params.AccountMetaSlice = nil
				buf := new(bytes.Buffer)
				err := encodeT(*params, buf)
				ag_require.NoError(t, err)
				//
				got := new(InitializeMint)
				err = decodeT(got, buf.Bytes())
				got.AccountMetaSlice = nil
				ag_require.NoError(t, err)
				ag_require.Equal(t, params, got)
			}
		})
	}
}

func TestDecodeInstruction_InitializeMint(t *testing.T) {
	mintAuthority := solana.NewWallet().PublicKey()
	freezeAuthority := solana.NewWallet().PublicKey()
	mint := solana.NewWallet().PublicKey()

	builder := NewInitializeMintInstructionBuilder().
		SetDecimals(9).
		SetMintAuthority(mintAuthority).
		SetFreezeAuthority(freezeAuthority).
		SetMintAccount(mint).
		SetSysVarRentPubkeyAccount(solana.SysVarRentPubkey)

	inst, err := builder.ValidateAndBuild()
	ag_require.NoError(t, err)

	data, err := inst.Data()
	ag_require.NoError(t, err)

	decoded, err := DecodeInstruction(nil, data)
	ag_require.NoError(t, err)
	ag_require.IsType(t, &InitializeMint{}, decoded.Impl)

	got := decoded.Impl.(*InitializeMint)
	ag_require.Equal(t, uint8(9), *got.Decimals)
	ag_require.Equal(t, mintAuthority, *got.MintAuthority)
	ag_require.Equal(t, freezeAuthority, *got.FreezeAuthority)
}
