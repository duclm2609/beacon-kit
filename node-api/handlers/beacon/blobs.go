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

package beacon

import (
	"github.com/berachain/beacon-kit/node-api/handlers"
	apitypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/math"
)

// GetBlobSidecars provides an implementation for the
// "/eth/v1/beacon/blob_sidecars/:block_id" API endpoint.
func (h *Handler) GetBlobSidecars(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[apitypes.GetBlobSidecarsRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}

	// Grab the requested slot.
	slot, err := utils.SlotFromBlockID(req.BlockID, h.backend)
	if err != nil {
		return nil, err
	}

	// Convert indices to uint64.
	indices := make([]uint64, len(req.Indices))
	for i, idxS := range req.Indices {
		var idx math.U64
		idx, err = math.U64FromString(idxS)
		if err != nil {
			return nil, err
		}
		indices[i] = idx.Unwrap()
	}

	// Grab the blob sidecars from the backend.
	blobSidecars, err := h.backend.BlobSidecarsByIndices(slot, indices)
	if err != nil {
		return nil, err
	}

	return apitypes.SidecarsResponse{
		Data: blobSidecars,
	}, nil
}
