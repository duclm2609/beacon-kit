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

package validator

import (
	"context"

	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/crypto"
)

// Service is responsible for building beacon blocks and sidecars.
type Service struct {
	// cfg is the validator config.
	cfg *Config
	// logger is a logger.
	logger log.Logger
	// chainSpec is the chain spec.
	chainSpec ChainSpec
	// signer is used to retrieve the public key of this node.
	signer crypto.BLSSigner
	// blobFactory is used to create blob sidecars for blocks.
	blobFactory BlobFactory
	// sb is the beacon state backend.
	sb StorageBackend
	// stateProcessor is responsible for processing the state.
	stateProcessor StateProcessor
	// localPayloadBuilder represents the local block builder, this builder
	// is connected to this nodes execution client via the EngineAPI.
	// Building blocks are done by submitting forkchoice updates through.
	// The local Builder.
	localPayloadBuilder PayloadBuilder
	// metrics is a metrics collector.
	metrics *validatorMetrics
}

// NewService creates a new validator service.
func NewService(
	cfg *Config,
	logger log.Logger,
	chainSpec ChainSpec,
	sb StorageBackend,
	stateProcessor StateProcessor,
	signer crypto.BLSSigner,
	blobFactory BlobFactory,
	localPayloadBuilder PayloadBuilder,
	ts TelemetrySink,
) *Service {
	return &Service{
		cfg:                 cfg,
		logger:              logger,
		sb:                  sb,
		chainSpec:           chainSpec,
		signer:              signer,
		stateProcessor:      stateProcessor,
		blobFactory:         blobFactory,
		localPayloadBuilder: localPayloadBuilder,
		metrics:             newValidatorMetrics(ts),
	}
}

// Name returns the name of the service.
func (s *Service) Name() string {
	return "validator"
}

func (s *Service) Start(
	_ context.Context,
) error {
	return nil
}

func (s *Service) Stop() error {
	return nil
}
