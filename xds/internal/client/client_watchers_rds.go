/*
 *
 * Copyright 2020 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package client

import (
	"fmt"
	"time"
)

type rdsUpdate struct {
	clusterName string
}
type rdsCallbackFunc func(rdsUpdate, error)

// watchRDS starts a listener watcher for the service..
//
// Note that during race (e.g. an xDS response is received while the user is
// calling cancel()), there's a small window where the callback can be called
// after the watcher is canceled. The caller needs to handle this case.
func (c *Client) watchRDS(routeName string, cb rdsCallbackFunc) (cancel func()) {
	wi := &watchInfo{
		typeURL:     rdsURL,
		target:      routeName,
		rdsCallback: cb,
	}

	wi.expiryTimer = time.AfterFunc(defaultWatchExpiryTimeout, func() {
		c.scheduleCallback(wi, rdsUpdate{}, fmt.Errorf("xds: RDS target %s not found, watcher timeout", routeName))
	})
	return c.watch(wi)
}
