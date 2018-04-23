package iterator

import (
	"context"

	"github.com/cayleygraph/cayley/graph/values"
)

var _ Iterator = &Not{}

// Not iterator acts like a complement for the primary iterator.
// It will return all the vertices which are not part of the primary iterator.
type Not struct {
	uid       uint64
	primaryIt Iterator
	allIt     Iterator
	result    values.Ref
	runstats  IteratorStats
	err       error
}

func NewNot(primaryIt, allIt Iterator) *Not {
	return &Not{
		uid:       NextUID(),
		primaryIt: primaryIt,
		allIt:     allIt,
	}
}

func (it *Not) UID() uint64 {
	return it.uid
}

// Reset resets the internal iterators and the iterator itself.
func (it *Not) Reset() {
	it.result = nil
	it.primaryIt.Reset()
	it.allIt.Reset()
}

func (it *Not) TagResults(dst map[string]values.Ref) {
	if it.primaryIt != nil {
		it.primaryIt.TagResults(dst)
	}
}

// SubIterators returns a slice of the sub iterators.
// The first iterator is the primary iterator, for which the complement
// is generated.
func (it *Not) SubIterators() []Generic {
	return []Generic{it.primaryIt, it.allIt}
}

// Next advances the Not iterator. It returns whether there is another valid
// new value. It fetches the next value of the all iterator which is not
// contained by the primary iterator.
func (it *Not) Next(ctx context.Context) bool {
	it.runstats.Next += 1

	for it.allIt.Next(ctx) {
		if curr := it.allIt.Result(); !it.primaryIt.Contains(ctx, curr) {
			it.result = curr
			it.runstats.ContainsNext += 1
			return true
		}
	}
	it.err = it.allIt.Err()
	return false
}

func (it *Not) Err() error {
	return it.err
}

func (it *Not) Result() values.Ref {
	return it.result
}

// Contains checks whether the passed value is part of the primary iterator's
// complement. For a valid value, it updates the Result returned by the iterator
// to the value itself.
func (it *Not) Contains(ctx context.Context, val values.Ref) bool {
	it.runstats.Contains += 1

	if it.primaryIt.Contains(ctx, val) {
		return false
	}

	it.err = it.primaryIt.Err()
	if it.err != nil {
		// Explicitly return 'false', since an error occurred.
		return false
	}

	it.result = val
	return true
}

// NextPath checks whether there is another path. Not applicable, hence it will
// return false.
func (it *Not) NextPath(ctx context.Context) bool {
	return false
}

// Close closes the primary and all iterators.  It closes all subiterators
// it can, but returns the first error it encounters.
func (it *Not) Close() error {
	err := it.primaryIt.Close()

	_err := it.allIt.Close()
	if _err != nil && err == nil {
		err = _err
	}

	return err
}

func (it *Not) Stats() IteratorStats {
	primaryStats := it.primaryIt.Stats()
	allStats := it.allIt.Stats()
	return IteratorStats{
		NextCost:     allStats.NextCost + primaryStats.ContainsCost,
		ContainsCost: primaryStats.ContainsCost,
		Size:         allStats.Size - primaryStats.Size,
		ExactSize:    false,
		Next:         it.runstats.Next,
		Contains:     it.runstats.Contains,
		ContainsNext: it.runstats.ContainsNext,
	}
}

func (it *Not) Size() (int64, bool) {
	st := it.Stats()
	return st.Size, st.ExactSize
}

func (it *Not) String() string {
	return "Not"
}
