// Copyright 2019 TiKV Project Authors.

//

// Licensed under the Apache License, Version 2.0 (the "License");

// you may not use this file except in compliance with the License.

// You may obtain a copy of the License at

//

//     http://www.apache.org/licenses/LICENSE-2.0

//

// Unless required by applicable law or agreed to in writing, software

// distributed under the License is distributed on an "AS IS" BASIS,

// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.

// See the License for the specific language governing permissions and

// limitations under the License.

package pd

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tikv/pd/pkg/testutil"
	"github.com/tikv/pd/pkg/tsoutil"
	"github.com/tikv/pd/tests"
)

const testBaseClientURL = "tmp://test.url:5255"

func TestNewBaseClient(t *testing.T) {
	re := require.New(t)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()
	urls := []string{testBaseClientURL}
	bc := newBaseClient(ctx, urls, SecurityOption{})
	re.Equal(urls, bc.GetURLs())
	re.Equal(defaultPDTimeout, bc.option.timeout)

}

func TestNewBaseClientValues(t *testing.T) {
	re := require.New(t)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()
	urls := []string{testBaseClientURL}
	bc := newBaseClient(ctx, urls, SecurityOption{})
	err := bc.initClusterID()
	re.EqualError(err, "[pd] failed to get cluster id")
	members, err := bc.getMembers(ctx, testBaseClientURL, defaultPDTimeout)
	re.NotEqual(err, nil)
	re.Nil(members)
	err = bc.updateMember()
	re.EqualError(err, fmt.Sprintf("[PD:client:ErrClientGetLeader]get leader from %v error", urls))

	re.Equal(uint64(0), bc.GetClusterID(ctx))
	re.Equal([]string{}, bc.GetFollowerAddrs())
	re.Equal(map[string]string{}, bc.GetAllocatorLeaderURLs())
	addr, ok := bc.getAllocatorLeaderAddrByDCLocation("dc-1")
	re.False(ok)
	re.Empty(addr)

	re.Equal("", bc.GetLeaderAddr())
	err = bc.switchLeader([]string{testBaseClientURL})
	re.Nil(err)
	re.Equal(testBaseClientURL, bc.GetLeaderAddr())
	bc.leader.Swap("")
	re.Equal("", bc.GetLeaderAddr())
	re.Equal([]string{}, bc.GetFollowerAddrs())
	bc.followers.Store([]string{testBaseClientURL})
	re.Equal([]string{testBaseClientURL}, bc.GetFollowerAddrs())

}

func TestNewBaseClientSchedule(t *testing.T) {
	re := require.New(t)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()
	urls := []string{testBaseClientURL}
	bc := newBaseClient(ctx, urls, SecurityOption{})
	bc.ScheduleCheckLeader()
	checkLeader := <-bc.checkLeaderCh
	re.Equal(struct{}{}, checkLeader)
	bc.scheduleCheckTSODispatcher()
	checkTSODispatcher := <-bc.checkTSODispatcherCh
	re.Equal(struct{}{}, checkTSODispatcher)
	bc.scheduleUpdateConnectionCtxs()
	updateConnectionCtxs := <-bc.updateConnectionCtxsCh
	re.Equal(struct{}{}, updateConnectionCtxs)
}

func TestInitRetry(t *testing.T) {
	re := require.New(t)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()
	urls := []string{testBaseClientURL}
	bc := newBaseClient(ctx, urls, SecurityOption{})
	re.Nil(bc.initRetry(func() error { return nil }))
	re.EqualError(bc.init(), "[pd] failed to get cluster id")

}

func runServer(re *require.Assertions, cluster *tests.TestCluster) []string {
	err := cluster.RunInitialServers()
	re.NoError(err)
	cluster.WaitLeader()
	leaderServer := cluster.GetServer(cluster.GetLeader())
	re.NoError(leaderServer.BootstrapCluster())

	testServers := cluster.GetServers()
	endpoints := make([]string, 0, len(testServers))
	for _, s := range testServers {
		endpoints = append(endpoints, s.GetConfig().AdvertiseClientUrls)
	}
	return endpoints
}

func TestClientGetLeader(t *testing.T) {
	re := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cluster, err := tests.NewTestCluster(ctx, 3)
	re.NoError(err)
	defer cluster.Destroy()

	endpoints := runServer(re, cluster)
	cli, err := NewClientWithContext(ctx, endpoints, SecurityOption{})
	re.NoError(err)
	var ts1, ts2 uint64
	testutil.Eventually(re, func() bool {
		p1, l1, err := cli.GetTS(context.TODO())
		if err == nil {
			ts1 = tsoutil.ComposeTS(p1, l1)
			return true
		}
		t.Log(err)
		return false
	})
	re.True(cluster.CheckTSOUnique(ts1))
	leader := cluster.GetLeader()

	err = cluster.GetServer(leader).Stop()
	re.NoError(err)
	leader = cluster.WaitLeader()
	re.NotEmpty(leader)

	// Check TS won't fall back after leader changed.
	testutil.Eventually(re, func() bool {
		p2, l2, err := cli.GetTS(context.TODO())
		if err == nil {
			ts2 = tsoutil.ComposeTS(p2, l2)
			return true
		}
		t.Log(err)
		return false
	})
	re.True(cluster.CheckTSOUnique(ts2))
	re.Less(ts1, ts2)

	// Check URL list.
	cli.Close()
}
