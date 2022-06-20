// Copyright 2022 TiKV Project Authors.

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

// NOTE: The code in this file is based on code from the

// etcd project, licensed under the Apache License v2.0

//

// https://github.com/etcd-io/etcd/blob/release-3.3/pkg/transport/listener.go

//

// Copyright 2015 The etcd Authors

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

package tlsutil

import (
	"crypto/tls"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTLSInfoClientConfig(t *testing.T) {
	re := require.New(t)
	info := new(TLSInfo)
	cfg, err := info.ClientConfig()
	re.Nil(err)
	re.Equal(tls.NoClientCert, cfg.ClientAuth)
	isEmpty := info.Empty()
	re.Equal(true, isEmpty)
	re.Equal("", cfg.ServerName)
	info.CertFile = "../../tests/client/cert/pd-server.pem"
	re.Equal(false, info.Empty())
	cfg, err = info.ClientConfig()
	re.Equal(fmt.Errorf("KeyFile and CertFile must both be present[key: %v, cert: %v]", info.KeyFile, info.CertFile), err)
	re.Nil(cfg)
}
