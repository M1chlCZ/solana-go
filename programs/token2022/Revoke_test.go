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

func TestEncodeDecode_Revoke(t *testing.T) {
	fu := ag_gofuzz.New().NilChance(0)
	for i := 0; i < 1; i++ {
		t.Run("Revoke"+strconv.Itoa(i), func(t *testing.T) {
			{
				params := new(Revoke)
				fu.Fuzz(params)
				params.Accounts = nil
				params.Signers = nil
				buf := new(bytes.Buffer)
				err := encodeT(*params, buf)
				ag_require.NoError(t, err)
				//
				got := new(Revoke)
				err = decodeT(got, buf.Bytes())
				got.Accounts = nil
				got.Signers = nil
				ag_require.NoError(t, err)
				ag_require.Equal(t, params, got)
			}
		})
	}
}

func TestRevoke_SetOwnerAccountDefaultsToSigner(t *testing.T) {
	builder := NewRevokeInstructionBuilder()
	owner := solana.NewWallet().PublicKey()

	builder.SetOwnerAccount(owner)

	ag_require.True(t, builder.GetOwnerAccount().IsSigner)
	ag_require.Empty(t, builder.Signers)
}

func TestRevoke_ValidateRequiresSigner(t *testing.T) {
	builder := NewRevokeInstructionBuilder().
		SetSourceAccount(solana.NewWallet().PublicKey())

	builder.Accounts[1] = solana.Meta(solana.NewWallet().PublicKey())

	_, err := builder.ValidateAndBuild()
	ag_require.EqualError(t, err, "accounts.Signers is not set")
}

func TestRevoke_ValidateRejectsTooManySigners(t *testing.T) {
	builder := NewRevokeInstructionBuilder().
		SetSourceAccount(solana.NewWallet().PublicKey()).
		SetOwnerAccount(solana.NewWallet().PublicKey(), solana.NewWallet().PublicKey())

	for i := 0; i < MAX_SIGNERS; i++ {
		builder.Signers = append(builder.Signers, solana.Meta(solana.NewWallet().PublicKey()).SIGNER())
	}

	_, err := builder.ValidateAndBuild()
	ag_require.EqualError(t, err, "too many signers; got 12, but max is 11")
}

func TestDecodeInstruction_Revoke(t *testing.T) {
	source := solana.NewWallet().PublicKey()
	owner := solana.NewWallet().PublicKey()
	signer := solana.NewWallet().PublicKey()

	builder := NewRevokeInstructionBuilder().
		SetSourceAccount(source).
		SetOwnerAccount(owner, signer)

	inst, err := builder.ValidateAndBuild()
	ag_require.NoError(t, err)
	ag_require.Len(t, inst.Accounts(), 3)

	data, err := inst.Data()
	ag_require.NoError(t, err)

	decoded, err := DecodeInstruction(inst.Accounts(), data)
	ag_require.NoError(t, err)
	ag_require.IsType(t, &Revoke{}, decoded.Impl)

	got := decoded.Impl.(*Revoke)
	ag_require.Len(t, got.Accounts, 2)
	ag_require.Len(t, got.Signers, 1)
}
