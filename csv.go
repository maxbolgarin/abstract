package abstract

import (
	"encoding/csv"
	"fmt"
	"io"
	"maps"
	"os"
	"sort"
	"strings"
	"sync"
)

// CSVTable represents a table of data from a CSV file where the first column is used as the ID
// for each row, and the remaining columns are stored with row order preserved.
type CSVTable struct {
	// Store ordered row IDs (first column values)
	ids []string
	// Map for fast lookup by ID
	idIndex map[string]int
	// Headers (column names)
	headers []string
	// Map for fast header lookup
	headerIndex map[string]int
	// Store rows data in a slice for each row, preserving order
	rows [][]string
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
// If the records are empty or if there are not enough headers (< 2), returns an empty table.
func NewCSVTable(records [][]string) *CSVTable {
	table := &CSVTable{
		ids:         make([]string, 0),
		idIndex:     make(map[string]int),
		headerIndex: make(map[string]int),
		rows:        make([][]string, 0),
	}

	if len(records) == 0 {
		return table
	}

	headers := records[0]

	if len(headers) < 2 {
		return table
	}

	// Set up headers and header index
	table.headers = make([]string, len(headers))
	copy(table.headers, headers)

	for i, header := range headers {
		table.headerIndex[header] = i
	}

	// Process data rows
	for i := 1; i < len(records); i++ {
		row := records[i]

		if len(row) == 0 || row[0] == "" {
			continue
		}

		rowID := row[0]
		// Store the row index
		table.idIndex[rowID] = len(table.ids)
		// Add ID to ordered list
		table.ids = append(table.ids, rowID)

		// Store row values in the same order as headers
		rowValues := make([]string, len(headers))
		for j := 0; j < len(headers) && j < len(row); j++ {
			rowValues[j] = row[j]
		}
		table.rows = append(table.rows, rowValues)
	}

	return table
}

// AddRow adds a new row to the table with the given ID and data.
// If the row has no data, it will not be added.
func (t *CSVTable) AddRow(id string, row map[string]string) {
	if len(row) == 0 {
		return
	}

	// Create a new row with all values initialized to empty strings
	newRow := make([]string, len(t.headers))
	newRow[0] = id // Set ID as first column

	// Fill in values from the provided map
	for colName, value := range row {
		if colIndex, exists := t.headerIndex[colName]; exists {
			newRow[colIndex] = value
		}
	}

	// If this ID already exists, update the existing row
	if index, exists := t.idIndex[id]; exists {
		t.rows[index] = newRow
	} else {
		// Otherwise add as a new row
		t.idIndex[id] = len(t.ids)
		t.ids = append(t.ids, id)
		t.rows = append(t.rows, newRow)
	}
}

// AppendColumn adds a new column to the table with the given name and values.
// Values are assigned to rows in order. If there are more rows than values,
// the remaining rows will not have a value for this column.
func (t *CSVTable) AppendColumn(column string, values []string) {
	// Add column to headers
	colIndex := len(t.headers)
	t.headers = append(t.headers, column)
	t.headerIndex[column] = colIndex

	// Expand each row to accommodate the new column
	for i := range t.rows {
		t.rows[i] = append(t.rows[i], "")
	}

	// Assign values to rows in order
	for i := 0; i < len(t.rows) && i < len(values); i++ {
		t.rows[i][colIndex] = values[i]
	}
}

// Row returns the data for the row with the given ID.
// If no row with that ID exists, returns an empty map.
func (t *CSVTable) Row(slug string) map[string]string {
	rowIndex, ok := t.idIndex[slug]
	if !ok {
		return make(map[string]string)
	}

	result := make(map[string]string, len(t.headers)-1)
	rowData := t.rows[rowIndex]

	// Skip the first column (ID) when creating the map
	for j := 1; j < len(t.headers) && j < len(rowData); j++ {
		result[t.headers[j]] = rowData[j]
	}

	return result
}

// LookupRow returns the data for the row with the given ID and a boolean indicating
// if the row exists.
func (t *CSVTable) LookupRow(slug string) (map[string]string, bool) {
	rowIndex, ok := t.idIndex[slug]
	if !ok {
		return nil, false
	}

	result := make(map[string]string, len(t.headers)-1)
	rowData := t.rows[rowIndex]

	// Skip the first column (ID) when creating the map
	for j := 1; j < len(t.headers) && j < len(rowData); j++ {
		result[t.headers[j]] = rowData[j]
	}

	return result, true
}

// RowSorted returns a map of ID to row data in the original sorted order.
func (t *CSVTable) RowSorted(id string) []string {
	index, ok := t.idIndex[id]
	if !ok {
		return nil
	}
	if index < 0 || index >= len(t.rows) {
		return nil
	}
	return t.rows[index]
}

// RowSorted returns a map of ID to row data in the original sorted order.
func (t *CSVTable) LookupRowSorted(id string) ([]string, bool) {
	index, ok := t.idIndex[id]
	if !ok {
		return nil, false
	}
	if index < 0 || index >= len(t.rows) {
		return nil, false
	}
	return t.rows[index], true
}

// All returns all rows in the table as a map of ID to row data.
func (t *CSVTable) All() map[string]map[string]string {
	result := make(map[string]map[string]string, len(t.ids))

	for i, id := range t.ids {
		rowMap := make(map[string]string, len(t.headers)-1)
		rowData := t.rows[i]

		// Skip the first column (ID) when creating each map
		for j := 1; j < len(t.headers) && j < len(rowData); j++ {
			rowMap[t.headers[j]] = rowData[j]
		}

		result[id] = rowMap
	}

	return result
}

// AllRows returns all rows in the table as a slice of row data maps.
func (t *CSVTable) AllRows() []map[string]string {
	rows := make([]map[string]string, len(t.rows))

	for i, rowData := range t.rows {
		rowMap := make(map[string]string, len(t.headers)-1)

		// Skip the first column (ID) when creating each map
		for j := 1; j < len(t.headers) && j < len(rowData); j++ {
			rowMap[t.headers[j]] = rowData[j]
		}

		rows[i] = rowMap
	}

	return rows
}

// AllSorted returns all rows in the table as a slice of maps, preserving the original order.
func (t *CSVTable) AllSorted() [][]string {
	result := make([][]string, len(t.rows))

	for i, row := range t.rows {
		rowCopy := make([]string, len(row))
		copy(rowCopy, row)
		result[i] = rowCopy
	}

	return result
}

// Copy creates a deep copy of the CSVTable.
// This is useful if you need to modify the data without affecting the original.
func (t *CSVTable) Copy() *CSVTable {
	table := &CSVTable{
		ids:         make([]string, len(t.ids)),
		idIndex:     make(map[string]int, len(t.idIndex)),
		headers:     make([]string, len(t.headers)),
		headerIndex: make(map[string]int, len(t.headerIndex)),
		rows:        make([][]string, len(t.rows)),
	}

	// Copy IDs and idIndex
	copy(table.ids, t.ids)
	maps.Copy(table.idIndex, t.idIndex)

	// Copy headers and headerIndex
	copy(table.headers, t.headers)
	maps.Copy(table.headerIndex, t.headerIndex)

	// Copy rows (deep copy)
	for i, row := range t.rows {
		table.rows[i] = make([]string, len(row))
		copy(table.rows[i], row)
	}

	return table
}

// AllIDs returns a slice of all row IDs in the table.
func (t *CSVTable) AllIDs() []string {
	ids := make([]string, len(t.ids))
	copy(ids, t.ids)
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
	rowIndex, ok := t.idIndex[slug]
	if !ok {
		return ""
	}

	colIndex, ok := t.headerIndex[key]
	if !ok {
		return ""
	}

	if colIndex < len(t.rows[rowIndex]) {
		return t.rows[rowIndex][colIndex]
	}

	return ""
}

// Has returns true if a row with the given ID exists in the table.
func (t *CSVTable) Has(slug string) bool {
	_, ok := t.idIndex[slug]
	return ok
}

// Bytes returns the table as a CSV-formatted byte slice.
func (t *CSVTable) Bytes() []byte {
	var buf strings.Builder

	// Write headers
	for i, header := range t.headers {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("\"" + header + "\"")
	}
	buf.WriteString("\n")

	// Write rows
	for _, rowData := range t.rows {
		for i, value := range rowData {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString("\"" + strings.ReplaceAll(value, "\"", "\"\"") + "\"")
		}
		buf.WriteString("\n")
	}

	return []byte(buf.String())
}

// DeleteColumns removes the specified columns from the table.
// This affects both the headers and the data in each row.
func (t *CSVTable) DeleteColumns(columns ...string) {
	// Identify columns to delete
	colIndicesToDelete := make(map[int]bool)
	for _, col := range columns {
		if colIndex, exists := t.headerIndex[col]; exists {
			colIndicesToDelete[colIndex] = true
			delete(t.headerIndex, col)
		}
	}

	if len(colIndicesToDelete) == 0 {
		return
	}

	// Create new headers without deleted columns
	newHeaders := make([]string, 0, len(t.headers)-len(colIndicesToDelete))
	for i, header := range t.headers {
		if !colIndicesToDelete[i] {
			newHeaders = append(newHeaders, header)
		}
	}

	// Update rows: remove deleted columns
	for i, row := range t.rows {
		newRow := make([]string, 0, len(row)-len(colIndicesToDelete))
		for j, val := range row {
			if !colIndicesToDelete[j] {
				newRow = append(newRow, val)
			}
		}
		t.rows[i] = newRow
	}

	// Update headers
	t.headers = newHeaders

	// Rebuild header index
	t.headerIndex = make(map[string]int, len(t.headers))
	for i, header := range t.headers {
		t.headerIndex[header] = i
	}
}

// SortDirection represents the sorting direction (ascending or descending)
type SortDirection int

const (
	// ASCSort sorts in ascending order
	ASCSort SortDirection = iota
	// DESCSort sorts in descending order
	DESCSort
)

// Sort reorders the table rows based on the values in the specified column.
// If the column does not exist, no sorting is performed.
// The direction parameter determines whether sorting is done in ascending or descending order.
func (t *CSVTable) Sort(column string, direction SortDirection) *CSVTable {
	colIndex, exists := t.headerIndex[column]
	if !exists {
		return t
	}

	// Create a stable sort to preserve the original order when values are equal
	sort.SliceStable(t.rows, func(i, j int) bool {
		if direction == ASCSort {
			return t.rows[i][colIndex] < t.rows[j][colIndex]
		}
		return t.rows[i][colIndex] > t.rows[j][colIndex]
	})

	// Update the IDs to match the new row order
	for i, row := range t.rows {
		t.ids[i] = row[0]
	}

	// Rebuild the idIndex map to reflect the new ordering
	for i, id := range t.ids {
		t.idIndex[id] = i
	}

	return t
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

// Row returns a copy of the row with the given ID.
func (t *CSVTableSafe) Row(slug string) map[string]string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.table.Row(slug)
}

// LookupRow returns a copy of the row with the given ID and whether it exists.
func (t *CSVTableSafe) LookupRow(slug string) (map[string]string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.table.LookupRow(slug)
}

// All returns a copy of all rows in the table.
func (t *CSVTableSafe) All() map[string]map[string]string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.table.All()
}

// AllRows returns a copy of all rows as a slice of maps.
func (t *CSVTableSafe) AllRows() []map[string]string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.table.AllRows()
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

// Sort reorders the table rows in a thread-safe manner based on the values in the specified column.
func (t *CSVTableSafe) Sort(column string, direction SortDirection) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.table.Sort(column, direction)
}

// Unwrap returns the underlying CSVTable.
// WARNING: This breaks thread safety. Only use when you're sure no other
// goroutines are accessing the table.
func (t *CSVTableSafe) Unwrap() *CSVTable {
	return t.table
}

// AllSorted returns all rows in the table as a slice of maps, preserving the original order.
func (t *CSVTableSafe) AllSorted() [][]string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.table.AllSorted()
}

// RowSorted returns a map of ID to row data in the original sorted order.
func (t *CSVTableSafe) RowSorted(id string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.table.RowSorted(id)
}

// LookupRowSorted returns a map of ID to row data in the original sorted order.
func (t *CSVTableSafe) LookupRowSorted(id string) ([]string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.table.LookupRowSorted(id)
}
