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

//nolint:dupl // separate file for ease of future implementation
package merkle

import (
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/merkle"
)

// ProveExecutionFeeRecipientInBlock generates a proof for the fee recipient in
// the latest execution payload header in the beacon block. The proof is then
// verified against the beacon block root as a sanity check. Returns the proof
// along with the beacon block root. It uses the fastssz library to generate the
// proof.
func ProveExecutionFeeRecipientInBlock[
	BeaconStateMarshallableT types.BeaconStateMarshallable,
](
	bbh *ctypes.BeaconBlockHeader,
	bs types.BeaconState[BeaconStateMarshallableT],
) ([]common.Root, common.Root, error) {
	// Get the proof of the execution fee recipient in the beacon state.
	feeRecipientInStateProof, leaf, err := ProveExecutionFeeRecipientInState(bs)
	if err != nil {
		return nil, common.Root{}, err
	}

	// Then get the proof of the beacon state in the beacon block.
	stateInBlockProof, err := ProveBeaconStateInBlock(bbh, false)
	if err != nil {
		return nil, common.Root{}, err
	}

	// Sanity check that the combined proof verifies against our beacon root.
	//
	//nolint:gocritic // ok.
	combinedProof := append(feeRecipientInStateProof, stateInBlockProof...)
	beaconRoot, err := verifyExecutionFeeRecipientInBlock(
		bbh, combinedProof, leaf,
	)
	if err != nil {
		return nil, common.Root{}, err
	}

	return combinedProof, beaconRoot, nil
}

// ProveExecutionFeeRecipientInState generates a proof for the execution fee
// recipient in the beacon state. It uses the fastssz library to generate the
// proof.
func ProveExecutionFeeRecipientInState[
	BeaconStateMarshallableT types.BeaconStateMarshallable,
](
	bs types.BeaconState[BeaconStateMarshallableT],
) ([]common.Root, common.Root, error) {
	bsm, err := bs.GetMarshallable()
	if err != nil {
		return nil, common.Root{}, err
	}
	stateProofTree, err := bsm.GetTree()
	if err != nil {
		return nil, common.Root{}, err
	}

	feeRecipientInStateProof, err := stateProofTree.Prove(
		ExecutionFeeRecipientGIndexDenebState,
	)
	if err != nil {
		return nil, common.Root{}, err
	}

	proof := make([]common.Root, len(feeRecipientInStateProof.Hashes))
	for i, hash := range feeRecipientInStateProof.Hashes {
		proof[i] = common.NewRootFromBytes(hash)
	}
	return proof, common.NewRootFromBytes(feeRecipientInStateProof.Leaf), nil
}

// verifyExecutionFeeRecipientInBlock verifies the execution fee recipient in
// the
// beacon block, returning the beacon block root used to verify against.
//
// TODO: verifying the proof is not absolutely necessary.
func verifyExecutionFeeRecipientInBlock(
	bbh *ctypes.BeaconBlockHeader,
	proof []common.Root,
	leaf common.Root,
) (common.Root, error) {
	beaconRoot := bbh.HashTreeRoot()
	if !merkle.VerifyProof(beaconRoot, leaf, ExecutionFeeRecipientGIndexDenebBlock, proof) {
		return common.Root{}, errors.Wrapf(
			errors.New("execution fee recipient proof failed to verify against beacon root"),
			"beacon root: 0x%s", beaconRoot,
		)
	}

	return beaconRoot, nil
}
