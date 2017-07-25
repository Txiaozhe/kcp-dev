/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/07/15     Tang Xiaoji
 */

package main

import (
	"fmt"
	"io"
	"net"

	"github.com/xtaci/kcp-go"

	"kcp-dev/config"
)

func main() {
	// kcp
	conn, err := kcp.DialWithOptions(config.REMOTE_ADDR, nil, config.DATA_SHARD, config.PARITY_SHARD)
	if err != nil {
		fmt.Println("line 41 connection error: ", err)
		return
	}

	fmt.Println("connect success")

	// 监听tcp
	listenTcp(conn)
}

func listenTcp(kcpConn net.Conn) {
	defer kcpConn.Close()

	listener, err := net.Listen("tcp", ":50000")
	if err != nil {
		fmt.Println("error tcp listen", err.Error())
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error accept", err.Error())
		}

		go func(conn net.Conn) {
			for {
				buf := make([]byte, 512)

				_, err := conn.Read(buf)
				if err == io.EOF {
					fmt.Println("read finish")
					return
				}

				_, err = kcpConn.Write([]byte(string(buf)))
				if err != nil {
					fmt.Println("write error")
					return
				}
			}
		}(conn)
	}
}
