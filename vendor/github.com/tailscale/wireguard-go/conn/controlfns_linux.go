/* SPDX-License-Identifier: MIT
 *
 * Copyright (C) 2017-2022 WireGuard LLC. All Rights Reserved.
 */

package conn

import (
	"fmt"
	"syscall"
)

func init() {
	controlFns = append(controlFns,

		// Attempt to set the socket buffer size beyond net.core.{r,w}mem_max by
		// using SO_*BUFFORCE. This requires CAP_NET_ADMIN, and is allowed here to
		// fail silently - the result of failure is lower performance on very fast
		// links or high latency links.
		func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				// Set up to *mem_max
				_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, socketBufferSize)
				_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_SNDBUF, socketBufferSize)
				// Set beyond *mem_max if CAP_NET_ADMIN
				_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUFFORCE, socketBufferSize)
				_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_SNDBUFFORCE, socketBufferSize)
			})
		},

		// Enable receiving of the packet information (IP_PKTINFO for IPv4,
		// IPV6_PKTINFO for IPv6) that is used to implement sticky socket support.
		func(network, address string, c syscall.RawConn) error {
			var err error
			switch network {
			case "udp4":
				c.Control(func(fd uintptr) {
					err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_PKTINFO, 1)
				})
			case "udp6":
				c.Control(func(fd uintptr) {
					err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IPV6, syscall.IPV6_RECVPKTINFO, 1)
					if err != nil {
						return
					}
					err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, 1)
				})
			default:
				err = fmt.Errorf("unhandled network: %s: %w", network, syscall.EINVAL)
			}
			return err
		},
	)
}
