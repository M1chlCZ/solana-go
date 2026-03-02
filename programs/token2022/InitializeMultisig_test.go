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

func TestEncodeDecode_InitializeMultisig(t *testing.T) {
	fu := ag_gofuzz.New().NilChance(0)
	for i := 0; i < 1; i++ {
		t.Run("InitializeMultisig"+strconv.Itoa(i), func(t *testing.T) {
			{
				params := new(InitializeMultisig)
				fu.Fuzz(params)
				params.Accounts = nil
				params.Signers = nil
				buf := new(bytes.Buffer)
				err := encodeT(*params, buf)
				ag_require.NoError(t, err)
				//
				got := new(InitializeMultisig)
				err = decodeT(got, buf.Bytes())
				got.Accounts = nil
				params.Signers = nil
				ag_require.NoError(t, err)
				ag_require.Equal(t, params, got)
			}
		})
	}
}

func TestNewInitializeMultisigInstructionBuilder_DefaultsRentSysvar(t *testing.T) {
	builder := NewInitializeMultisigInstructionBuilder()
	ag_require.Equal(t, solana.SysVarRentPubkey, builder.GetSysVarRentPubkeyAccount().PublicKey)
}

func TestInitializeMultisig_ValidateRequiresM(t *testing.T) {
	builder := NewInitializeMultisigInstructionBuilder()
	_, err := builder.ValidateAndBuild()
	ag_require.EqualError(t, err, "M parameter is not set")
}

func TestInitializeMultisig_ValidateRejectsTooManySigners(t *testing.T) {
	builder := NewInitializeMultisigInstructionBuilder().
		SetM(2).
		SetAccount(solana.NewWallet().PublicKey())

	signers := make([]solana.PublicKey, 0, MAX_SIGNERS+1)
	for i := 0; i < MAX_SIGNERS+1; i++ {
		signers = append(signers, solana.NewWallet().PublicKey())
	}
	builder.AddSigners(signers...)

	_, err := builder.ValidateAndBuild()
	ag_require.EqualError(t, err, "too many signers; got 12, but max is 11")
}

func TestDecodeInstruction_InitializeMultisig(t *testing.T) {
	account := solana.NewWallet().PublicKey()
	signerA := solana.NewWallet().PublicKey()
	signerB := solana.NewWallet().PublicKey()

	builder := NewInitializeMultisigInstructionBuilder().
		SetM(2).
		SetAccount(account).
		SetSysVarRentPubkeyAccount(solana.SysVarRentPubkey).
		AddSigners(signerA, signerB)

	inst, err := builder.ValidateAndBuild()
	ag_require.NoError(t, err)
	ag_require.Len(t, inst.Accounts(), 4)

	data, err := inst.Data()
	ag_require.NoError(t, err)

	decoded, err := DecodeInstruction(inst.Accounts(), data)
	ag_require.NoError(t, err)
	ag_require.IsType(t, &InitializeMultisig{}, decoded.Impl)

	got := decoded.Impl.(*InitializeMultisig)
	ag_require.Len(t, got.Accounts, 2)
	ag_require.Len(t, got.Signers, 2)
}
