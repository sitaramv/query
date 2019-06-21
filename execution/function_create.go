//  Copyright (c) 2019 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package execution

import (
	"encoding/json"

	"github.com/couchbase/query/functions"
	"github.com/couchbase/query/plan"
	"github.com/couchbase/query/value"
)

type CreateFunction struct {
	base
	plan *plan.CreateFunction
}

func NewCreateFunction(plan *plan.CreateFunction, context *Context) *CreateFunction {
	rv := &CreateFunction{
		plan: plan,
	}

	newRedirectBase(&rv.base)
	rv.output = rv
	return rv
}

func (this *CreateFunction) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitCreateFunction(this)
}

func (this *CreateFunction) Copy() Operator {
	rv := &CreateFunction{plan: this.plan}
	this.base.copy(&rv.base)
	return rv
}

func (this *CreateFunction) RunOnce(context *Context, parent value.Value) {
	this.once.Do(func() {
		defer context.Recover(&this.base) // Recover from any panic
		if !this.active() {
			return
		}
		defer this.close(context)
		this.switchPhase(_EXECTIME)
		defer this.switchPhase(_NOTIME)
		defer this.notify() // Notify that I have stopped

		if context.Readonly() {
			return
		}

		// Actually create function
		this.switchPhase(_SERVTIME)
		err := functions.AddFunction(this.plan.Name(), this.plan.Body())
		this.switchPhase(_EXECTIME)
		if err != nil {
			context.Error(err)
		}
	})
}

func (this *CreateFunction) MarshalJSON() ([]byte, error) {
	r := this.plan.MarshalBase(func(r map[string]interface{}) {
		this.marshalTimes(r)
	})
	return json.Marshal(r)
}