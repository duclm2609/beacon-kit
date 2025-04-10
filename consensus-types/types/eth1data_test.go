// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types_test

import (
	"io"
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/encoding/ssz"
	karalabessz "github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

func TestEth1Data_Serialization(t *testing.T) {
	t.Parallel()
	original := types.NewEth1Data(common.Root{})
	data, err := original.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	unmarshalled := new(types.Eth1Data)
	err = ssz.Unmarshal(data, unmarshalled)
	require.NoError(t, err)
	require.Equal(t, original, unmarshalled)

	var buf []byte
	buf, err = original.MarshalSSZTo(buf)
	require.NoError(t, err)

	// The two byte slices should be equal
	require.Equal(t, data, buf)
}

func TestEth1Data_UnmarshalError(t *testing.T) {
	t.Parallel()

	var unmarshalled types.Eth1Data
	err := ssz.Unmarshal([]byte{}, &unmarshalled)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestEth1Data_SizeSSZ(t *testing.T) {
	t.Parallel()
	eth1Data := types.NewEth1Data(common.Root{})
	size := karalabessz.Size(eth1Data)
	require.Equal(t, uint32(72), size)
}

func TestEth1Data_HashTreeRoot(t *testing.T) {
	t.Parallel()
	eth1Data := types.NewEth1Data(common.Root{})

	require.NotPanics(t, func() {
		_ = eth1Data.HashTreeRoot()
	})
}

func TestEth1Data_GetTree(t *testing.T) {
	t.Parallel()
	eth1Data := types.NewEth1Data(common.Root{})
	tree, err := eth1Data.GetTree()

	require.NoError(t, err)
	require.NotNil(t, tree)
}

func TestEth1Data_GetDepositCount(t *testing.T) {
	t.Parallel()
	eth1Data := &types.Eth1Data{
		DepositRoot:  common.Root{},
		DepositCount: 10,
		BlockHash:    common.ExecutionHash{},
	}

	count := eth1Data.GetDepositCount()

	require.Equal(t, uint64(10), count.Unwrap())
}
