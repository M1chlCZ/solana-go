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

func TestEncodeDecode_SetAuthority(t *testing.T) {
	fu := ag_gofuzz.New().NilChance(0)
	for i := 0; i < 1; i++ {
		t.Run("SetAuthority"+strconv.Itoa(i), func(t *testing.T) {
			{
				params := new(SetAuthority)
				fu.Fuzz(params)
				params.Accounts = nil
				params.Signers = nil
				buf := new(bytes.Buffer)
				err := encodeT(*params, buf)
				ag_require.NoError(t, err)
				//
				got := new(SetAuthority)
				err = decodeT(got, buf.Bytes())
				got.Accounts = nil
				got.Signers = nil
				ag_require.NoError(t, err)
				ag_require.Equal(t, params, got)
			}
		})
	}
}

func TestSetAuthority_SetAuthorityAccountDefaultsToSigner(t *testing.T) {
	builder := NewSetAuthorityInstructionBuilder()
	authority := solana.NewWallet().PublicKey()

	builder.SetAuthorityAccount(authority)

	ag_require.True(t, builder.GetAuthorityAccount().IsSigner)
	ag_require.Empty(t, builder.Signers)
}

func TestSetAuthority_ValidateRequiresAuthorityType(t *testing.T) {
	builder := NewSetAuthorityInstructionBuilder().
		SetSubjectAccount(solana.NewWallet().PublicKey()).
		SetAuthorityAccount(solana.NewWallet().PublicKey())

	_, err := builder.ValidateAndBuild()
	ag_require.EqualError(t, err, "AuthorityType parameter is not set")
}

func TestSetAuthority_ValidateRequiresSigner(t *testing.T) {
	builder := NewSetAuthorityInstructionBuilder().
		SetAuthorityType(AuthorityMintTokens).
		SetSubjectAccount(solana.NewWallet().PublicKey())

	builder.Accounts[1] = solana.Meta(solana.NewWallet().PublicKey())

	_, err := builder.ValidateAndBuild()
	ag_require.EqualError(t, err, "accounts.Signers is not set")
}

func TestSetAuthority_AllowsNilNewAuthority(t *testing.T) {
	builder := NewSetAuthorityInstructionBuilder().
		SetAuthorityType(AuthorityMintTokens).
		SetSubjectAccount(solana.NewWallet().PublicKey()).
		SetAuthorityAccount(solana.NewWallet().PublicKey())

	inst, err := builder.ValidateAndBuild()
	ag_require.NoError(t, err)

	data, err := inst.Data()
	ag_require.NoError(t, err)

	decoded, err := DecodeInstruction(inst.Accounts(), data)
	ag_require.NoError(t, err)

	got := decoded.Impl.(*SetAuthority)
	ag_require.Nil(t, got.NewAuthority)
}

func TestDecodeInstruction_SetAuthority(t *testing.T) {
	subject := solana.NewWallet().PublicKey()
	authority := solana.NewWallet().PublicKey()
	newAuthority := solana.NewWallet().PublicKey()
	signer := solana.NewWallet().PublicKey()

	builder := NewSetAuthorityInstructionBuilder().
		SetAuthorityType(AuthorityCloseAccount).
		SetNewAuthority(newAuthority).
		SetSubjectAccount(subject).
		SetAuthorityAccount(authority, signer)

	inst, err := builder.ValidateAndBuild()
	ag_require.NoError(t, err)
	ag_require.Len(t, inst.Accounts(), 3)

	data, err := inst.Data()
	ag_require.NoError(t, err)

	decoded, err := DecodeInstruction(inst.Accounts(), data)
	ag_require.NoError(t, err)
	ag_require.IsType(t, &SetAuthority{}, decoded.Impl)

	got := decoded.Impl.(*SetAuthority)
	ag_require.Equal(t, AuthorityCloseAccount, *got.AuthorityType)
	ag_require.Equal(t, newAuthority, *got.NewAuthority)
	ag_require.Len(t, got.Accounts, 2)
	ag_require.Len(t, got.Signers, 1)
}
