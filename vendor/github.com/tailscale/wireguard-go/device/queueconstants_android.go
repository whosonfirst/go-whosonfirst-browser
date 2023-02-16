/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2017-2022 WireGuard LLC. All Rights Reserved.
 */

package device

import "github.com/tailscale/wireguard-go/conn"

/* Reduce memory consumption for Android */

const (
	QueueStagedSize            = conn.DefaultBatchSize
	QueueOutboundSize          = 1024
	QueueInboundSize           = 1024
	QueueHandshakeSize         = 1024
	MaxSegmentSize             = 2200
	PreallocatedBuffersPerPool = 4096
)
