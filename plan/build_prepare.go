//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package plan

import (
	"github.com/couchbase/query/algebra"
	"github.com/couchbase/query/value"
)

func (this *builder) VisitPrepare(stmt *algebra.Prepare) (interface{}, error) {
	plan, err := BuildPrepared(stmt.Statement(), this.datastore, this.systemstore, this.namespace, false)
	if err != nil {
		return nil, err
	}

	PreparedCache().AddPrepared(plan)

	json_bytes, err := plan.MarshalJSON()
	if err != nil {
		return nil, err
	}

	value := value.NewValue(json_bytes)

	return NewPrepare(value), nil
}
