// Copyright 2014 The Cayley Authors. All rights reserved.
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

package iterator

// Define the general iterator interface.

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/cayleygraph/cayley/graph/values"
)

var nextIteratorID uint64

func init() {
	atomic.StoreUint64(&nextIteratorID, 1)
}

func NextUID() uint64 {
	return atomic.AddUint64(&nextIteratorID, 1) - 1
}

var (
	_ Iterator = &Null{}
	_ Iterator = &Error{}
)

type Morphism func(Iterator) Iterator

// Here we define the simplest iterator -- the Null iterator. It contains nothing.
// It is the empty set. Often times, queries that contain one of these match nothing,
// so it's important to give it a special iterator.
type Null struct {
	uid uint64
}

// Fairly useless New function.
func NewNull() *Null {
	return &Null{uid: NextUID()}
}

func (it *Null) UID() uint64 {
	return it.uid
}

// Fill the map based on the tags assigned to this iterator.
func (it *Null) TagResults(dst map[string]values.Value) {}

func (it *Null) Contains(ctx context.Context, v values.Value) bool {
	return false
}

func (it *Null) String() string {
	return "Null"
}

func (it *Null) Next(ctx context.Context) bool {
	return false
}

func (it *Null) Err() error {
	return nil
}

func (it *Null) Result() values.Value {
	return nil
}

func (it *Null) SubIterators() []Generic {
	return nil
}

func (it *Null) NextPath(ctx context.Context) bool {
	return false
}

func (it *Null) Size() (int64, bool) {
	return 0, true
}

func (it *Null) Reset() {}

func (it *Null) Close() error {
	return nil
}

// A null iterator costs nothing. Use it!
func (it *Null) Stats() IteratorStats {
	return IteratorStats{}
}

// Error iterator always returns a single error with no other results.
type Error struct {
	uid uint64
	err error
}

func NewError(err error) *Error {
	return &Error{uid: NextUID(), err: err}
}

func (it *Error) UID() uint64 {
	return it.uid
}

// Fill the map based on the tags assigned to this iterator.
func (it *Error) TagResults(dst map[string]values.Value) {}

func (it *Error) Contains(ctx context.Context, v values.Value) bool {
	return false
}

func (it *Error) String() string {
	return fmt.Sprintf("Error(%v)", it.err)
}

func (it *Error) Next(ctx context.Context) bool {
	return false
}

func (it *Error) Err() error {
	return it.err
}

func (it *Error) Result() values.Value {
	return nil
}

func (it *Error) SubIterators() []Generic {
	return nil
}

func (it *Error) NextPath(ctx context.Context) bool {
	return false
}

func (it *Error) Size() (int64, bool) {
	return 0, true
}

func (it *Error) Reset() {}

func (it *Error) Close() error {
	return it.err
}

func (it *Error) Stats() IteratorStats {
	return IteratorStats{}
}
