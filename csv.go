package abstract

import (
	"encoding/csv"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strings"
	"sync"
)

// CSVTable represents a table of data from a CSV file where the first column is used as the ID
// for each row, and the remaining columns are stored as key-value pairs.
type CSVTable struct {
	data    map[string]map[string]string
	headers []string
}

// NewCSVTableFromFilePath creates a new CSVTable from a file at the given path.
// Returns an error if the file cannot be opened or parsed.
func NewCSVTableFromFilePath(path string) (*CSVTable, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	return NewCSVTableFromReader(file)
}

// NewCSVTableFromReader creates a new CSVTable from any io.Reader that contains CSV data.
// Returns an error if the CSV data cannot be parsed.
func NewCSVTableFromReader(reader io.Reader) (*CSVTable, error) {
	records, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return NewCSVTable(records), nil
}

// NewCSVTable creates a new CSVTable from the given records.
// The first row is considered the header row, and the first column is used as the ID for each row.
// Each row's data is stored as a map of column name to value.
// If the records are empty or if there are not enough headers (< 2), returns an empty table.
func NewCSVTable(records [][]string) *CSVTable {
	table := &CSVTable{
		data: make(map[string]map[string]string),
	}

	if len(records) == 0 {
		return table
	}

	headers := records[0]

	if len(headers) < 2 {
		return table
	}

	table.headers = headers

	for i := 1; i < len(records); i++ {
		row := records[i]

		if len(row) == 0 || row[0] == "" {
			continue
		}

		rowID := row[0]
		rowData := make(map[string]string, len(headers)-1)

		for j := 1; j < len(headers) && j < len(row); j++ {
			rowData[headers[j]] = row[j]
		}
		table.data[rowID] = rowData
	}

	return table
}

// AddRow adds a new row to the table with the given ID and data.
// If the row has no data, it will not be added.
func (t *CSVTable) AddRow(id string, row map[string]string) {
	if len(row) == 0 {
		return
	}
	t.data[id] = row
}

// AppendColumn adds a new column to the table with the given name and values.
// Values are assigned to rows in order. If there are more rows than values,
// the remaining rows will not have a value for this column.
func (t *CSVTable) AppendColumn(column string, values []string) {
	t.headers = append(t.headers, column)
	var i int
	for _, row := range t.data {
		if i >= len(values) {
			break
		}
		row[column] = values[i]
		i++
	}
}

// Row returns the data for the row with the given ID.
// If no row with that ID exists, returns an empty map.
// WARNING: The returned map is a direct reference, not a copy.
// Modifying the returned map will modify the internal state of the CSVTable,
// which is generally not recommended. Use Copy() if you need to modify the data.
func (t *CSVTable) Row(slug string) map[string]string {
	row, ok := t.data[slug]
	if !ok {
		return make(map[string]string)
	}
	return row
}

// LookupRow returns the data for the row with the given ID and a boolean indicating
// if the row exists.
// WARNING: The returned map is a direct reference, not a copy.
// Modifying the returned map will modify the internal state of the CSVTable,
// which is generally not recommended. Use Copy() if you need to modify the data.
func (t *CSVTable) LookupRow(slug string) (map[string]string, bool) {
	row, ok := t.data[slug]
	return row, ok
}

// All returns all rows in the table as a map of ID to row data.
// WARNING: The returned maps are direct references, not copies.
// Modifying the returned maps will modify the internal state of the CSVTable,
// which is generally not recommended. Use Copy() if you need to modify the data.
func (t *CSVTable) All() map[string]map[string]string {
	return t.data
}

// AllRows returns all rows in the table as a slice of row data maps.
// WARNING: The returned maps are direct references, not copies.
// Modifying the returned maps will modify the internal state of the CSVTable,
// which is generally not recommended. Use Copy() if you need to modify the data.
func (t *CSVTable) AllRows() []map[string]string {
	rows := make([]map[string]string, 0, len(t.data))
	for _, row := range t.data {
		rows = append(rows, row)
	}
	return rows
}

// Copy creates a deep copy of the CSVTable.
// This is useful if you need to modify the data without affecting the original.
func (t *CSVTable) Copy() *CSVTable {
	table := &CSVTable{
		data:    make(map[string]map[string]string),
		headers: make([]string, len(t.headers)),
	}
	copy(table.headers, t.headers)
	for slug, row := range t.data {
		table.data[slug] = make(map[string]string, len(row))
		maps.Copy(table.data[slug], row)
	}
	return table
}

// AllIDs returns a slice of all row IDs in the table.
func (t *CSVTable) AllIDs() []string {
	ids := make([]string, 0, len(t.data))
	for id := range t.data {
		ids = append(ids, id)
	}
	return ids
}

// Headers returns a copy of the headers for the table.
func (t *CSVTable) Headers() []string {
	headers := make([]string, len(t.headers))
	copy(headers, t.headers)
	return headers
}

// Value returns the value for the given ID and key.
// If no row with that ID exists, or if the key doesn't exist in that row,
// returns an empty string.
func (t *CSVTable) Value(slug, key string) string {
	row := t.Row(slug)
	return row[key]
}

// Has returns true if a row with the given ID exists in the table.
func (t *CSVTable) Has(slug string) bool {
	_, ok := t.data[slug]
	return ok
}

// Bytes returns the table as a CSV-formatted byte slice.
func (t *CSVTable) Bytes() []byte {
	var rows []byte
	rows = append(rows, []byte(strings.Join(t.headers, ","))...)
	rows = append(rows, '\n')

	for slug, row := range t.data {
		values := make([]string, 0, len(row)+1)
		values = append(values, "\""+slug+"\"")
		for _, header := range t.headers[1:] {
			values = append(values, "\""+row[header]+"\"")
		}
		rows = append(rows, []byte(strings.Join(values, ","))...)
		rows = append(rows, '\n')
	}

	return rows
}

// DeleteColumns removes the specified columns from the table.
// This affects both the headers and the data in each row.
func (t *CSVTable) DeleteColumns(columns ...string) {
	for _, row := range t.data {
		for _, col := range columns {
			delete(row, col)
		}
	}
	for _, col := range columns {
		for i, header := range t.headers {
			if header == col {
				t.headers = slices.Delete(t.headers, i, i+1)
				break
			}
		}
	}
}

// CSVTableSafe is a thread-safe wrapper around CSVTable that provides
// synchronized access to the underlying data using a mutex.
type CSVTableSafe struct {
	table *CSVTable
	mu    sync.RWMutex
}

// NewCSVTableSafeFromFilePath creates a new thread-safe CSVTable from a file path.
func NewCSVTableSafeFromFilePath(path string) (*CSVTableSafe, error) {
	table, err := NewCSVTableFromFilePath(path)
	if err != nil {
		return nil, err
	}
	return &CSVTableSafe{table: table}, nil
}

// NewCSVTableSafeFromReader creates a new thread-safe CSVTable from a reader.
func NewCSVTableSafeFromReader(reader io.Reader) (*CSVTableSafe, error) {
	table, err := NewCSVTableFromReader(reader)
	if err != nil {
		return nil, err
	}
	return &CSVTableSafe{table: table}, nil
}

// NewCSVTableSafe creates a new thread-safe CSVTable from records.
func NewCSVTableSafe(records [][]string) *CSVTableSafe {
	return &CSVTableSafe{
		table: NewCSVTable(records),
	}
}

// AddRow adds a new row to the table in a thread-safe manner.
func (t *CSVTableSafe) AddRow(id string, row map[string]string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.table.AddRow(id, row)
}

// AppendColumn adds a new column to the table in a thread-safe manner.
func (t *CSVTableSafe) AppendColumn(column string, values []string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.table.AppendColumn(column, values)
}

// Row returns a copy of the row with the given ID to avoid concurrent modification issues.
func (t *CSVTableSafe) Row(slug string) map[string]string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	row := t.table.Row(slug)
	// Create a copy to avoid returning references to internal data
	result := make(map[string]string, len(row))
	maps.Copy(result, row)
	return result
}

// LookupRow returns a copy of the row with the given ID and whether it exists.
func (t *CSVTableSafe) LookupRow(slug string) (map[string]string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	row, exists := t.table.LookupRow(slug)
	if !exists {
		return nil, false
	}
	// Create a copy to avoid returning references to internal data
	result := make(map[string]string, len(row))
	maps.Copy(result, row)
	return result, true
}

// All returns a deep copy of all rows in the table.
func (t *CSVTableSafe) All() map[string]map[string]string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	all := t.table.All()
	// Create a deep copy to avoid returning references to internal data
	result := make(map[string]map[string]string, len(all))
	for id, row := range all {
		rowCopy := make(map[string]string, len(row))
		maps.Copy(rowCopy, row)
		result[id] = rowCopy
	}
	return result
}

// AllRows returns a deep copy of all rows as a slice of maps.
func (t *CSVTableSafe) AllRows() []map[string]string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	rows := t.table.AllRows()
	// Create a deep copy
	result := make([]map[string]string, len(rows))
	for i, row := range rows {
		rowCopy := make(map[string]string, len(row))
		for k, v := range row {
			rowCopy[k] = v
		}
		result[i] = rowCopy
	}
	return result
}

// Copy creates a deep copy of the CSVTableSafe, including its internal table.
func (t *CSVTableSafe) Copy() *CSVTableSafe {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return &CSVTableSafe{
		table: t.table.Copy(),
	}
}

// AllIDs returns a copy of all row IDs in the table.
func (t *CSVTableSafe) AllIDs() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.table.AllIDs()
}

// Headers returns a copy of the headers for the table.
func (t *CSVTableSafe) Headers() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.table.Headers()
}

// Value returns the value for the given ID and key.
func (t *CSVTableSafe) Value(slug, key string) string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.table.Value(slug, key)
}

// Has returns true if a row with the given ID exists in the table.
func (t *CSVTableSafe) Has(slug string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.table.Has(slug)
}

// Bytes returns the table as a CSV-formatted byte slice.
func (t *CSVTableSafe) Bytes() []byte {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.table.Bytes()
}

// DeleteColumns removes the specified columns from the table.
func (t *CSVTableSafe) DeleteColumns(columns ...string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.table.DeleteColumns(columns...)
}

// Unwrap returns the underlying CSVTable.
// WARNING: This breaks thread safety. Only use when you're sure no other
// goroutines are accessing the table.
func (t *CSVTableSafe) Unwrap() *CSVTable {
	return t.table
}
