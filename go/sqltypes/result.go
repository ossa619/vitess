// Copyright 2015, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqltypes

import (
	"github.com/youtube/vitess/go/vt/binlog/eventtoken"
	querypb "github.com/youtube/vitess/go/vt/proto/query"
)

// Result represents a query result.
type Result struct {
	Fields       []*querypb.Field      `json:"fields"`
	RowsAffected uint64                `json:"rows_affected"`
	InsertID     uint64                `json:"insert_id"`
	Rows         [][]Value             `json:"rows"`
	Extras       *querypb.ResultExtras `json:"extras"`
}

// ResultStream is an interface for receiving Result. It is used for
// RPC interfaces.
type ResultStream interface {
	// Recv returns the next result on the stream.
	// It will return io.EOF if the stream ended.
	Recv() (*Result, error)
}

// Repair fixes the type info in the rows
// to conform to the supplied field types.
func (result *Result) Repair(fields []*querypb.Field) {
	// Usage of j is intentional.
	for j, f := range fields {
		for _, r := range result.Rows {
			if r[j].typ != Null {
				r[j].typ = f.Type
			}
		}
	}
}

// Copy creates a deep copy of Result.
func (result *Result) Copy() *Result {
	out := &Result{
		InsertID:     result.InsertID,
		RowsAffected: result.RowsAffected,
	}
	if result.Fields != nil {
		fieldsp := make([]*querypb.Field, len(result.Fields))
		fields := make([]querypb.Field, len(result.Fields))
		for i, f := range result.Fields {
			fields[i] = *f
			fieldsp[i] = &fields[i]
		}
		out.Fields = fieldsp
	}
	if result.Rows != nil {
		rows := make([][]Value, len(result.Rows))
		for i, r := range result.Rows {
			rows[i] = make([]Value, len(r))
			totalLen := 0
			for _, c := range r {
				totalLen += len(c.val)
			}
			arena := make([]byte, 0, totalLen)
			for j, c := range r {
				start := len(arena)
				arena = append(arena, c.val...)
				rows[i][j] = MakeTrusted(c.typ, arena[start:start+len(c.val)])
			}
		}
		out.Rows = rows
	}
	if result.Extras != nil {
		out.Extras = &querypb.ResultExtras{
			Fresher: result.Extras.Fresher,
		}
		if result.Extras.EventToken != nil {
			out.Extras.EventToken = &querypb.EventToken{
				Timestamp: result.Extras.EventToken.Timestamp,
				Shard:     result.Extras.EventToken.Shard,
				Position:  result.Extras.EventToken.Position,
			}
		}
	}
	return out
}

// MakeRowTrusted converts a *querypb.Row to []Value based on the types
// in fields. It does not sanity check the values against the type.
// Every place this function is called, a comment is needed that explains
// why it's justified.
func MakeRowTrusted(fields []*querypb.Field, row *querypb.Row) []Value {
	sqlRow := make([]Value, len(row.Lengths))
	var offset int64
	for i, length := range row.Lengths {
		if length < 0 {
			continue
		}
		sqlRow[i] = MakeTrusted(fields[i].Type, row.Values[offset:offset+length])
		offset += length
	}
	return sqlRow
}

// StripFieldNames will return a new Result that has the same Rows,
// but the Field objects will have their Name emptied.  Note we don't
// proto.Copy each Field for performance reasons, but we only copy the
// individual fields.
func (result *Result) StripFieldNames() *Result {
	if len(result.Fields) == 0 {
		return result
	}
	r := *result
	r.Fields = make([]*querypb.Field, len(result.Fields))
	newFieldsArray := make([]querypb.Field, len(result.Fields))
	for i, f := range result.Fields {
		r.Fields[i] = &newFieldsArray[i]
		newFieldsArray[i].Type = f.Type
	}
	return &r
}

// AppendResult will combine the Results Objects of one result
// to another result.Note currently it doesn't handle cases like
// if two results have different fields.We will enhance this function.
func (result *Result) AppendResult(qr, innerqr *Result) {
	if innerqr.RowsAffected == 0 && len(innerqr.Fields) == 0 {
		return
	}
	if qr.Fields == nil {
		qr.Fields = innerqr.Fields
	}
	qr.RowsAffected += innerqr.RowsAffected
	if innerqr.InsertID != 0 {
		qr.InsertID = innerqr.InsertID
	}
	if len(qr.Rows) == 0 {
		// we haven't gotten any result yet, just save the new extras.
		qr.Extras = innerqr.Extras
	} else {
		// Merge the EventTokens / Fresher flags within Extras.
		if innerqr.Extras == nil {
			// We didn't get any from innerq. Have to clear any
			// we'd have gotten already.
			if qr.Extras != nil {
				qr.Extras.EventToken = nil
				qr.Extras.Fresher = false
			}
		} else {
			// We may have gotten an EventToken from
			// innerqr.  If we also got one earlier, merge
			// it. If we didn't get one earlier, we
			// discard the new one.
			if qr.Extras != nil {
				// Note if any of the two is nil, we get nil.
				qr.Extras.EventToken = eventtoken.Minimum(qr.Extras.EventToken, innerqr.Extras.EventToken)

				qr.Extras.Fresher = qr.Extras.Fresher && innerqr.Extras.Fresher
			}
		}
	}
	qr.Rows = append(qr.Rows, innerqr.Rows...)
}
