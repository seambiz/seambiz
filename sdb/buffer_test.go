package sdb_test

import (
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/seambiz/seambiz/sdb"
)

// Helper functions for older Go versions
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func TestAppend(t *testing.T) {
	sql := sdb.NewSQLStatement()
	sql.Append("SELECT")

	got := sql.String()
	if got != "SELECT " {
		t.Errorf("got '%s', want '%s'", got, "SELECT ")
	}

	got = sql.Query()
	if got != "SELECT " {
		t.Errorf("got '%s', want '%s'", got, "SELECT ")
	}
}

func TestAppendRaw(t *testing.T) {
	sql := sdb.NewSQLStatement()
	sql.AppendRaw("SELECT")

	got := sql.String()
	if got != "SELECT" {
		t.Errorf("got '%s', want '%s'", got, "SELECT")
	}

	got = sql.Query()
	if got != "SELECT" {
		t.Errorf("got '%s', want '%s'", got, "SELECT")
	}
}

func TestAppendInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		sql := sdb.NewSQLStatement()
		n := int(rand.Int31())
		sql.AppendInt(n)

		got := sql.Query()
		if got != strconv.Itoa(n) {
			t.Errorf("got '%s', want '%d'", got, n)
		}
	}
}

func TestFieldsSimple(t *testing.T) {
	sql := sdb.NewSQLStatement()
	fields := []string{"id", "test", "third"}

	sql.Fields("", fields)
	got := sql.Query()
	want := "id,test,third "
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func TestFieldsFull(t *testing.T) {
	sql := sdb.NewSQLStatement()
	fields := []string{"id", "test", "third"}

	sql.Fields("abc", fields)
	got := sql.Query()
	want := "abc.id,abc.test,abc.third "
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func TestFieldsCodegen(t *testing.T) {
	sql := sdb.NewSQLStatement()
	fields := []string{"id", "test", "third"}
	fields2 := []string{"id", "test", "third"}

	sql.Fields("a", fields)
	sql.Fields("b", fields2)
	got := sql.Query()
	want := "a.id,a.test,a.third ,b.id,b.test,b.third "
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}

func TestSQLStatement_AppendFields(t *testing.T) {
	type args struct {
		prepend   string
		prefix    string
		separator string
		append    string
		fields    []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empy",
			args: args{},
			want: "",
		},
		{
			name: "simple field list",
			args: args{
				separator: ",",
				fields:    []string{"test", "yes"},
			},
			want: "test,yes",
		},
		{
			name: "insert use case",
			args: args{
				prepend:   "(",
				prefix:    "",
				append:    ")",
				separator: ",",
				fields:    []string{"field1", "field2"},
			},
			want: "(field1,field2)",
		},
		{
			name: "select use case",
			args: args{
				prepend:   "SELECT ",
				prefix:    "A.",
				append:    " FROM table A",
				separator: ",",
				fields:    []string{"field1", "field2"},
			},
			want: "SELECT A.field1,A.field2 FROM table A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := sdb.NewSQLStatement()
			sql.AppendFields(tt.args.prepend, tt.args.prefix, tt.args.separator, tt.args.append, tt.args.fields)
			got := sql.Query()
			if got != tt.want {
				t.Errorf("got '%s', want '%s'", got, tt.want)
			}
		})
	}
}

func TestSQLStatement_AppendFiller(t *testing.T) {
	type args struct {
		prepend   string
		separator string
		append    string
		filler    string
		n         int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{},
			want: "",
		},
		{
			name: "insert use case",
			args: args{
				prepend:   "(",
				append:    ")",
				separator: ",",
				filler:    "?",
				n:         3,
			},
			want: "(?,?,?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := sdb.NewSQLStatement()
			sql.AppendFiller(tt.args.prepend, tt.args.separator, tt.args.append, tt.args.filler, tt.args.n)
			got := sql.Query()
			if got != tt.want {
				t.Errorf("got '%s', want '%s'", got, tt.want)
			}
		})
	}
}

func TestSQLStatement_InInt(t *testing.T) {
	tests := []struct {
		name string
		args []int
		want string
	}{
		{
			name: "empty",
			args: nil,
			want: "",
		},
		{
			name: "single",
			args: []int{99},
			want: "99",
		},
		{
			name: "three",
			args: []int{1, 2, 3},
			want: "1,2,3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := sdb.NewSQLStatement()
			sql.InInt(tt.args)
			got := sql.Query()
			if got != tt.want {
				t.Errorf("got '%s', want '%s'", got, tt.want)
			}
		})
	}
}

// PRIORITY 1 ISSUE TESTS

// TestBytes_DataRace tests that Bytes() returns data that won't be corrupted
// when the buffer is returned to the pool and reused.
// This is a CRITICAL bug - the current implementation returns a reference to
// the internal buffer which gets reset when Release() is called.
func TestBytes_DataRace(t *testing.T) {
	sql := sdb.NewSQLStatement()
	sql.Append("SELECT * FROM users WHERE id =")
	sql.AppendInt(42)

	// Get bytes - this should be a safe copy
	result := sql.Bytes()

	// At this point, the buffer has been returned to the pool and reset
	// Let's reuse the pool to prove the buffer gets corrupted
	sql2 := sdb.NewSQLStatement()
	sql2.Append("DIFFERENT QUERY")
	_ = sql2.Query()

	// The original result should still contain our original data
	// Note: Append() adds trailing space, AppendInt() does not
	expected := "SELECT * FROM users WHERE id = 42"
	got := string(result)

	if got != expected {
		t.Errorf("Bytes() data was corrupted after buffer reuse!\ngot:  '%s'\nwant: '%s'", got, expected)
	}
}


// TestFieldsCalled_PoolReuse verifies that fieldsCalled is properly reset
// when a buffer is returned to the pool. This was the original production bug
// causing SQL statements to start with a comma.
func TestFieldsCalled_PoolReuse(t *testing.T) {
	// First use: call Fields() which sets fieldsCalled = true
	sql1 := sdb.NewSQLStatement()
	sql1.Append("SELECT")
	sql1.Fields("t1", []string{"id", "name"})
	query1 := sql1.Query() // Returns to pool with Reset()

	expected1 := "SELECT t1.id,t1.name "
	if query1 != expected1 {
		t.Errorf("First query incorrect:\ngot:  '%s'\nwant: '%s'", query1, expected1)
	}

	// Second use: get from pool and immediately call Fields()
	// If fieldsCalled wasn't reset, this will start with a comma
	sql2 := sdb.NewSQLStatement()
	sql2.Append("SELECT")
	sql2.Fields("t2", []string{"id", "email"})
	query2 := sql2.Query()

	expected2 := "SELECT t2.id,t2.email "
	if query2 != expected2 {
		t.Errorf("Second query has leading comma bug!\ngot:  '%s'\nwant: '%s'", query2, expected2)
	}

	// Most important: verify no leading comma after the first word
	// This is the production bug - Fields() adding "," when fieldsCalled is true
	if query2 == ",t2.id,t2.email " || query2 == "SELECT ,t2.id,t2.email " {
		t.Error("CRITICAL: Query has leading comma - fieldsCalled was not reset!")
	}
}

// TestAppendUInt_ErrorHandling tests that appendUInt doesn't silently ignore errors
// (though in practice Write() never returns an error, this tests consistency)
func TestAppendUInt_Consistency(t *testing.T) {
	sql := sdb.NewSQLStatement()

	// Use the internal append method that calls appendUInt
	sql.Append(uint(12345))

	got := sql.Query()
	want := "12345 "

	if got != want {
		t.Errorf("appendUInt failed:\ngot:  '%s'\nwant: '%s'", got, want)
	}
}

// TestRelease_ResetsAllState verifies that Release() properly resets
// both the buffer and the fieldsCalled flag
func TestRelease_ResetsAllState(t *testing.T) {
	sql := sdb.NewSQLStatement()
	sql.Append("SELECT * FROM users")
	sql.Fields("", []string{"id", "name"})

	// Manually release instead of using Query()
	sql.Release()

	// Get the same object back from the pool
	sql2 := sdb.NewSQLStatement()

	// Verify buffer is empty
	if len(sql2.String()) != 0 {
		t.Errorf("Buffer not reset after Release(): '%s'", sql2.String())
	}

	// Verify fieldsCalled is reset by checking Fields() doesn't add leading comma
	sql2.Fields("t", []string{"col"})
	result := sql2.Query()

	if result[0] == ',' {
		t.Error("fieldsCalled was not reset by Release()")
	}
}

// TestConcurrentStatementBuilding_RaceCondition tests that concurrent goroutines
// building SQL statements don't share or corrupt each other's buffers.
//
// This reproduces the production bug where:
// - Order query: "SELECT o.id, o.created_at... FROM order o..."
// - Sensor query: "SELECT S.ver, S.geo_lat... FROM sensor S..."
// Were being built concurrently and ended up with interleaved content like:
// "SELECT o.id typ,S.ver,S.geo_lat...S.FROM order o..."
//
// Run with: go test -race -run TestConcurrentStatementBuilding_RaceCondition
func TestConcurrentStatementBuilding_RaceCondition(t *testing.T) {
	const numGoroutines = 100
	const iterations = 50

	var wg sync.WaitGroup
	errors := make(chan string, numGoroutines*iterations)

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for i := 0; i < iterations; i++ {
				// Simulate order query
				if goroutineID%2 == 0 {
					sql := sdb.NewSQLStatement()
					sql.AppendStr("SELECT ")
					sql.AppendStr("o.id, ")
					sql.AppendStr("o.created_at, ")
					sql.AppendStr("o.status ")
					sql.AppendStr("FROM order o ")
					sql.AppendStr("WHERE o.company_id = ?")

					query := sql.Query()

					// Validate: query must start with "SELECT o" (order query)
					// and must NOT contain sensor columns
					if !strings.HasPrefix(query, "SELECT o") {
						errors <- "Order query doesn't start with 'SELECT o': " + query[:minInt(50, len(query))]
					}
					if strings.Contains(query, "S.ver") || strings.Contains(query, "S.geo") {
						errors <- "Order query contains sensor columns (buffer corruption): " + query[:minInt(100, len(query))]
					}
					if query == "" {
						errors <- "Order query is empty (race condition)"
					}
				} else {
					// Simulate sensor query
					sql := sdb.NewSQLStatement()
					sql.AppendStr("SELECT ")
					sql.AppendStr("S.id, ")
					sql.AppendStr("S.ver, ")
					sql.AppendStr("S.geo_lat, ")
					sql.AppendStr("S.geo_long ")
					sql.AppendStr("FROM sensor S ")
					sql.AppendStr("WHERE S.company_id = ?")

					query := sql.Query()

					// Validate: query must start with "SELECT S" (sensor query)
					// and must NOT contain order columns
					if !strings.HasPrefix(query, "SELECT S") {
						errors <- "Sensor query doesn't start with 'SELECT S': " + query[:minInt(50, len(query))]
					}
					if strings.Contains(query, "o.id") || strings.Contains(query, "o.created") {
						errors <- "Sensor query contains order columns (buffer corruption): " + query[:minInt(100, len(query))]
					}
					if query == "" {
						errors <- "Sensor query is empty (race condition)"
					}
				}
			}
		}(g)
	}

	wg.Wait()
	close(errors)

	// Collect all errors
	var errMsgs []string
	for err := range errors {
		errMsgs = append(errMsgs, err)
	}

	if len(errMsgs) > 0 {
		t.Errorf("Found %d buffer corruption/race errors:\n", len(errMsgs))
		// Show first 10 errors to avoid flooding output
		for i, msg := range errMsgs {
			if i >= 10 {
				t.Errorf("... and %d more errors", len(errMsgs)-10)
				break
			}
			t.Errorf("  - %s", msg)
		}
	}
}

// TestConcurrentStatementBuilding_NoInterleaving is a simpler test that
// verifies statements built concurrently don't get mixed together.
func TestConcurrentStatementBuilding_NoInterleaving(t *testing.T) {
	const numGoroutines = 50

	var wg sync.WaitGroup
	results := make([]string, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			sql := sdb.NewSQLStatement()

			// Each goroutine builds a unique query with its ID
			marker := strconv.Itoa(id)
			sql.AppendStr("QUERY_")
			sql.AppendStr(marker)
			sql.AppendStr("_START ")
			// Add some bulk to increase chance of interleaving
			for j := 0; j < 10; j++ {
				sql.AppendStr("FIELD_")
				sql.AppendStr(marker)
				sql.AppendStr("_")
				sql.AppendInt(j)
				sql.AppendStr(" ")
			}
			sql.AppendStr("QUERY_")
			sql.AppendStr(marker)
			sql.AppendStr("_END")

			results[id] = sql.Query()
		}(i)
	}

	wg.Wait()

	// Verify each result contains ONLY its own marker
	for id, result := range results {
		marker := strconv.Itoa(id)
		expectedStart := "QUERY_" + marker + "_START"
		expectedEnd := "QUERY_" + marker + "_END"

		if !strings.HasPrefix(result, expectedStart) {
			t.Errorf("Goroutine %d: result doesn't start correctly\ngot: %s", id, result[:minInt(50, len(result))])
		}
		if !strings.HasSuffix(result, expectedEnd) {
			t.Errorf("Goroutine %d: result doesn't end correctly\ngot: ...%s", id, result[maxInt(0, len(result)-50):])
		}

		// Check for contamination from other goroutines
		for otherID := 0; otherID < numGoroutines; otherID++ {
			if otherID == id {
				continue
			}
			otherMarker := "QUERY_" + strconv.Itoa(otherID) + "_"
			if strings.Contains(result, otherMarker) {
				t.Errorf("Goroutine %d: result contains content from goroutine %d (buffer sharing detected)", id, otherID)
				break
			}
		}
	}
}
