/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	tls "github.com/dubbogo/getty/demo/hello/tls"
)

import (
	"github.com/dubbogo/getty"
	gxsync "github.com/dubbogo/gost/sync"
)

import (
	"github.com/dubbogo/getty/demo/util"
)

var (
	taskPoolMode        = flag.Bool("taskPool", false, "task pool mode")
	taskPoolQueueLength = flag.Int("task_queue_length", 100, "task queue length")
	taskPoolQueueNumber = flag.Int("task_queue_number", 4, "task queue number")
	taskPoolSize        = flag.Int("task_pool_size", 2000, "task poll size")
	pprofPort           = flag.Int("pprof_port", 65432, "pprof http port")
	Sessions            []getty.Session
)

var (
	taskPool *gxsync.TaskPool
)

func main() {
	flag.Parse()

	util.SetLimit()

	util.Profiling(*pprofPort)
	c := &getty.ServerTlsConfigBuilder{
		ServerKeyCertChainPath:        `E:\Projects\openSource\dubbo-samples\java\dubbo-samples-ssl\dubbo-samples-ssl-provider\src\main\resources\certs\server0.pem`,
		ServerPrivateKeyPath:          `E:\Projects\openSource\dubbo-samples\java\dubbo-samples-ssl\dubbo-samples-ssl-provider\src\main\resources\certs\server0.key`,
		ServerTrustCertCollectionPath: `E:\Projects\openSource\dubbo-samples\java\dubbo-samples-ssl\dubbo-samples-ssl-consumer\src\main\resources\certs\ca.pem`,
	}

	options := []getty.ServerOption{getty.WithLocalAddress(":8090"),
		getty.WithServerSslEnabled(true),
		getty.WithServerTlsConfigBuilder(c),
	}

	if *taskPoolMode {
		taskPool = gxsync.NewTaskPool(
			gxsync.WithTaskPoolTaskQueueLength(*taskPoolQueueLength),
			gxsync.WithTaskPoolTaskQueueNumber(*taskPoolQueueNumber),
			gxsync.WithTaskPoolTaskPoolSize(*taskPoolSize),
		)
	}

	server := getty.NewTCPServer(options...)

	go server.RunEventLoop(NewHelloServerSession)
	util.WaitCloseSignals(server)
}

func NewHelloServerSession(session getty.Session) (err error) {
	err = tls.InitialSession(session)
	Sessions = append(Sessions, session)
	if err != nil {
		return
	}
	session.SetTaskPool(taskPool)

	return
}
