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

package agollo

import (
	"github.com/stack-labs/stack-rpc/plugin/config/source/apollo/agollo/agcache"
	"github.com/stack-labs/stack-rpc/plugin/config/source/apollo/agollo/cluster"
	"github.com/stack-labs/stack-rpc/plugin/config/source/apollo/agollo/env/file"
	"github.com/stack-labs/stack-rpc/plugin/config/source/apollo/agollo/extension"
	"github.com/stack-labs/stack-rpc/plugin/config/source/apollo/agollo/protocol/auth"
)

//SetSignature 设置自定义 http 授权控件
func SetSignature(auth auth.HTTPAuth) {
	if auth != nil {
		extension.SetHTTPAuth(auth)
	}
}

//SetBackupFileHandler 设置自定义备份文件处理组件
func SetBackupFileHandler(file file.FileHandler) {
	if file != nil {
		extension.SetFileHandler(file)
	}
}

//SetLoadBalance 设置自定义负载均衡组件
func SetLoadBalance(loadBalance cluster.LoadBalance) {
	if loadBalance != nil {
		extension.SetLoadBalance(loadBalance)
	}
}

//SetCache 设置自定义cache组件
func SetCache(cacheFactory agcache.CacheFactory) {
	if cacheFactory != nil {
		extension.SetCacheFactory(cacheFactory)
	}
}
