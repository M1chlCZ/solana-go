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

func TestEncodeDecode_InitializeAccount(t *testing.T) {
	fu := ag_gofuzz.New().NilChance(0)
	for i := 0; i < 1; i++ {
		t.Run("InitializeAccount"+strconv.Itoa(i), func(t *testing.T) {
			{
				params := new(InitializeAccount)
				fu.Fuzz(params)
				params.AccountMetaSlice = nil
				buf := new(bytes.Buffer)
				err := encodeT(*params, buf)
				ag_require.NoError(t, err)
				//
				got := new(InitializeAccount)
				err = decodeT(got, buf.Bytes())
				got.AccountMetaSlice = nil
				ag_require.NoError(t, err)
				ag_require.Equal(t, params, got)
			}
		})
	}
}

func TestNewInitializeAccountInstructionBuilder_DefaultsRentSysvar(t *testing.T) {
	builder := NewInitializeAccountInstructionBuilder()
	ag_require.Equal(t, solana.SysVarRentPubkey, builder.GetSysVarRentPubkeyAccount().PublicKey)
}

func TestInitializeAccount_ValidateRequiresAccounts(t *testing.T) {
	builder := NewInitializeAccountInstructionBuilder()
	_, err := builder.ValidateAndBuild()
	ag_require.EqualError(t, err, "accounts.Account is not set")
}

func TestDecodeInstruction_InitializeAccount(t *testing.T) {
	account := solana.NewWallet().PublicKey()
	mint := solana.NewWallet().PublicKey()
	owner := solana.NewWallet().PublicKey()

	builder := NewInitializeAccountInstructionBuilder().
		SetAccount(account).
		SetMintAccount(mint).
		SetOwnerAccount(owner).
		SetSysVarRentPubkeyAccount(solana.SysVarRentPubkey)

	inst, err := builder.ValidateAndBuild()
	ag_require.NoError(t, err)
	ag_require.Len(t, inst.Accounts(), 4)

	data, err := inst.Data()
	ag_require.NoError(t, err)

	decoded, err := DecodeInstruction(nil, data)
	ag_require.NoError(t, err)
	ag_require.IsType(t, &InitializeAccount{}, decoded.Impl)
}
