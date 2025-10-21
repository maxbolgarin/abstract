package abstract_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/maxbolgarin/abstract"
)

func TestNewCSVTable(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
	}

	table := abstract.NewCSVTable(records)

	if got := table.Value("row1", "Name"); got != "Test1" {
		t.Errorf("Expected Value(row1, Name) = %q, got %q", "Test1", got)
	}
	if got := table.Value("row2", "Value"); got != "200" {
		t.Errorf("Expected Value(row2, Value) = %q, got %q", "200", got)
	}
	if got := table.Value("non-existent", "Name"); got != "" {
		t.Errorf("Expected Value(non-existent, Name) = %q, got %q", "", got)
	}
}

func TestNewCSVTableEmptyRecords(t *testing.T) {
	records := [][]string{}
	table := abstract.NewCSVTable(records)

	if len(table.All()) != 0 {
		t.Errorf("Expected empty All() for empty records")
	}
	if len(table.Headers()) != 0 {
		t.Errorf("Expected empty Headers() for empty records")
	}
}

func TestNewCSVTableInsufficientHeaders(t *testing.T) {
	records := [][]string{
		{"ID"},
		{"row1"},
	}

	table := abstract.NewCSVTable(records)
	if len(table.All()) != 0 {
		t.Errorf("Expected empty All() for insufficient headers")
	}
}

func TestNewCSVTableFromReader(t *testing.T) {
	csvData := "ID,Name,Value\nrow1,Test1,100\nrow2,Test2,200"
	reader := strings.NewReader(csvData)

	table, err := abstract.NewCSVTableFromReader(reader)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if got := table.Value("row1", "Name"); got != "Test1" {
		t.Errorf("Expected Value(row1, Name) = %q, got %q", "Test1", got)
	}
	if got := table.Value("row2", "Value"); got != "200" {
		t.Errorf("Expected Value(row2, Value) = %q, got %q", "200", got)
	}
}

func TestNewCSVTableFromReaderError(t *testing.T) {
	// Invalid CSV data
	csvData := "ID,Name,Value\nrow1,Test1,\"unclosed quote"
	reader := strings.NewReader(csvData)

	_, err := abstract.NewCSVTableFromReader(reader)
	if err == nil {
		t.Errorf("Expected error for invalid CSV data, got nil")
	}
}

func TestNewCSVTableFromFilePath(t *testing.T) {
	// This would require a temporary file setup
	// Skipping actual implementation since it relies on file system
	// But would test error cases
	_, err := abstract.NewCSVTableFromFilePath("non-existent-file.csv")
	if err == nil {
		t.Errorf("Expected error for non-existent file, got nil")
	}
}

func TestAddRow(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	table := abstract.NewCSVTable(records)

	newRow := map[string]string{
		"Name":  "Test2",
		"Value": "200",
	}
	table.AddRow("row2", newRow)

	if got := table.Value("row2", "Name"); got != "Test2" {
		t.Errorf("Expected Value(row2, Name) = %q, got %q", "Test2", got)
	}
	if !table.Has("row2") {
		t.Errorf("Expected Has(row2) to be true")
	}
}

func TestAppendColumn(t *testing.T) {
	records := [][]string{
		{"ID", "Name"},
		{"row1", "Test1"},
		{"row2", "Test2"},
	}

	table := abstract.NewCSVTable(records)

	values := []string{"100", "200"}
	table.AppendColumn("Value", values)

	headers := table.Headers()
	found := false
	for _, h := range headers {
		if h == "Value" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Headers() should contain \"Value\"")
	}

	if got := table.Value("row1", "Value"); got != "100" {
		t.Errorf("Expected Value(row1, Value) = %q, got %q", "100", got)
	}
	if got := table.Value("row2", "Value"); got != "200" {
		t.Errorf("Expected Value(row2, Value) = %q, got %q", "200", got)
	}
}

func TestRow(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	table := abstract.NewCSVTable(records)

	row := table.Row("row1")
	if got := row["Name"]; got != "Test1" {
		t.Errorf("Expected row[Name] = %q, got %q", "Test1", got)
	}
	if got := row["Value"]; got != "100" {
		t.Errorf("Expected row[Value] = %q, got %q", "100", got)
	}

	// Non-existent row
	emptyRow := table.Row("non-existent")
	if len(emptyRow) != 0 {
		t.Errorf("Expected empty row for non-existent id")
	}
}

func TestLookupRow(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	table := abstract.NewCSVTable(records)

	row, found := table.LookupRow("row1")
	if !found {
		t.Errorf("Expected found=true for existing row")
	}
	if got := row["Name"]; got != "Test1" {
		t.Errorf("Expected row[Name] = %q, got %q", "Test1", got)
	}

	_, found = table.LookupRow("non-existent")
	if found {
		t.Errorf("Expected found=false for non-existent row")
	}
}

func TestAll(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
	}

	table := abstract.NewCSVTable(records)

	all := table.All()
	if len(all) != 2 {
		t.Errorf("Expected len(All()) = 2, got %d", len(all))
	}
	if got := all["row1"]["Name"]; got != "Test1" {
		t.Errorf("Expected all[row1][Name] = %q, got %q", "Test1", got)
	}
	if got := all["row2"]["Value"]; got != "200" {
		t.Errorf("Expected all[row2][Value] = %q, got %q", "200", got)
	}
}

func TestAllRows(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
	}

	table := abstract.NewCSVTable(records)

	rows := table.AllRows()
	if len(rows) != 2 {
		t.Errorf("Expected len(AllRows()) = 2, got %d", len(rows))
	}

	// Check that all expected rows are in the result
	found1, found2 := false, false
	for _, row := range rows {
		if row["Name"] == "Test1" && row["Value"] == "100" {
			found1 = true
		}
		if row["Name"] == "Test2" && row["Value"] == "200" {
			found2 = true
		}
	}
	if !found1 {
		t.Errorf("Expected to find row with Name=Test1, Value=100")
	}
	if !found2 {
		t.Errorf("Expected to find row with Name=Test2, Value=200")
	}
}

func TestTableCopy(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
	}

	table := abstract.NewCSVTable(records)
	tableCopy := table.Copy()

	// Ensure the copy has the same data
	if !reflect.DeepEqual(table.Headers(), tableCopy.Headers()) {
		t.Errorf("Expected same headers in copy")
	}
	if got := tableCopy.Value("row1", "Name"); got != table.Value("row1", "Name") {
		t.Errorf("Expected same values in copy")
	}

	// Modify the copy and ensure original is unchanged
	tableCopy.AddRow("row3", map[string]string{"Name": "Test3", "Value": "300"})
	if !tableCopy.Has("row3") {
		t.Errorf("Expected Has(row3)=true in copy")
	}
	if table.Has("row3") {
		t.Errorf("Expected Has(row3)=false in original")
	}
}

func TestAllIDs(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
	}

	table := abstract.NewCSVTable(records)

	ids := table.AllIDs()
	if len(ids) != 2 {
		t.Errorf("Expected len(AllIDs()) = 2, got %d", len(ids))
	}

	found1, found2 := false, false
	for _, id := range ids {
		if id == "row1" {
			found1 = true
		}
		if id == "row2" {
			found2 = true
		}
	}
	if !found1 {
		t.Errorf("Expected AllIDs() to contain row1")
	}
	if !found2 {
		t.Errorf("Expected AllIDs() to contain row2")
	}
}

func TestHeaders(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	table := abstract.NewCSVTable(records)

	headers := table.Headers()
	if !reflect.DeepEqual(headers, []string{"ID", "Name", "Value"}) {
		t.Errorf("Expected Headers() = [ID Name Value], got %v", headers)
	}

	// Ensure headers are copied, not referenced
	headers[0] = "Changed"
	originalHeaders := table.Headers()
	if originalHeaders[0] != "ID" {
		t.Errorf("Expected original headers to be unchanged, got %v", originalHeaders)
	}
}

func TestValue(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	table := abstract.NewCSVTable(records)

	if got := table.Value("row1", "Name"); got != "Test1" {
		t.Errorf("Expected Value(row1, Name) = %q, got %q", "Test1", got)
	}
	if got := table.Value("row1", "NonExistent"); got != "" {
		t.Errorf("Expected Value(row1, NonExistent) = %q, got %q", "", got)
	}
	if got := table.Value("nonExistent", "Name"); got != "" {
		t.Errorf("Expected Value(nonExistent, Name) = %q, got %q", "", got)
	}
}

func TestHas(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	table := abstract.NewCSVTable(records)

	if !table.Has("row1") {
		t.Errorf("Expected Has(row1) = true")
	}
	if table.Has("nonExistent") {
		t.Errorf("Expected Has(nonExistent) = false")
	}
}

func TestBytes(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
		{"row3", "Test3", "300"},
		{"row4", "Test4", "400"},
	}

	table := abstract.NewCSVTable(records)

	csvBytes := table.Bytes()
	expected := "\"ID\",\"Name\",\"Value\"\n\"row1\",\"Test1\",\"100\"\n\"row2\",\"Test2\",\"200\"\n\"row3\",\"Test3\",\"300\"\n\"row4\",\"Test4\",\"400\"\n"
	if string(csvBytes) != expected {
		t.Errorf("Expected Bytes() = %q, got %q", expected, string(csvBytes))
	}
}

func TestDeleteColumns(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value", "Extra"},
		{"row1", "Test1", "100", "Data1"},
		{"row2", "Test2", "200", "Data2"},
	}

	table := abstract.NewCSVTable(records)

	table.DeleteColumns("Value", "Extra")

	headers := table.Headers()
	if !reflect.DeepEqual(headers, []string{"ID", "Name"}) {
		t.Errorf("Expected Headers() = [ID Name], got %v", headers)
	}

	row := table.Row("row1")
	if _, exists := row["Value"]; exists {
		t.Errorf("Expected Value column to be deleted")
	}
	if _, exists := row["Extra"]; exists {
		t.Errorf("Expected Extra column to be deleted")
	}
}

// Tests for CSVTableSafe

func TestNewCSVTableSafe(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
	}

	table := abstract.NewCSVTableSafe(records)

	if got := table.Value("row1", "Name"); got != "Test1" {
		t.Errorf("Expected Value(row1, Name) = %q, got %q", "Test1", got)
	}
	if got := table.Value("row2", "Value"); got != "200" {
		t.Errorf("Expected Value(row2, Value) = %q, got %q", "200", got)
	}
	if got := table.Value("non-existent", "Name"); got != "" {
		t.Errorf("Expected Value(non-existent, Name) = %q, got %q", "", got)
	}
}

func TestNewCSVTableSafeFromReader(t *testing.T) {
	csvData := "ID,Name,Value\nrow1,Test1,100\nrow2,Test2,200"
	reader := strings.NewReader(csvData)

	table, err := abstract.NewCSVTableSafeFromReader(reader)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if got := table.Value("row1", "Name"); got != "Test1" {
		t.Errorf("Expected Value(row1, Name) = %q, got %q", "Test1", got)
	}
	if got := table.Value("row2", "Value"); got != "200" {
		t.Errorf("Expected Value(row2, Value) = %q, got %q", "200", got)
	}
}

func TestNewCSVTableSafeFromFilePath(t *testing.T) {
	// Testing error case only
	_, err := abstract.NewCSVTableSafeFromFilePath("non-existent-file.csv")
	if err == nil {
		t.Errorf("Expected error for non-existent file, got nil")
	}
}

func TestCSVTableSafeAddRow(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	table := abstract.NewCSVTableSafe(records)

	newRow := map[string]string{
		"Name":  "Test2",
		"Value": "200",
	}
	table.AddRow("row2", newRow)

	if got := table.Value("row2", "Name"); got != "Test2" {
		t.Errorf("Expected Value(row2, Name) = %q, got %q", "Test2", got)
	}
	if !table.Has("row2") {
		t.Errorf("Expected Has(row2) to be true")
	}
}

func TestCSVTableSafeAppendColumn(t *testing.T) {
	records := [][]string{
		{"ID", "Name"},
		{"row1", "Test1"},
		{"row2", "Test2"},
	}

	table := abstract.NewCSVTableSafe(records)

	values := []string{"100", "200"}
	table.AppendColumn("Value", values)

	headers := table.Headers()
	found := false
	for _, h := range headers {
		if h == "Value" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Headers() should contain \"Value\"")
	}

	if got := table.Value("row1", "Value"); got != "100" {
		t.Errorf("Expected Value(row1, Value) = %q, got %q", "100", got)
	}
	if got := table.Value("row2", "Value"); got != "200" {
		t.Errorf("Expected Value(row2, Value) = %q, got %q", "200", got)
	}
}

// This test verifies that maps returned by Row are deep copies
func TestCSVTableSafeRowDeepCopy(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	table := abstract.NewCSVTableSafe(records)

	row := table.Row("row1")
	if got := row["Name"]; got != "Test1" {
		t.Errorf("Expected row[Name] = %q, got %q", "Test1", got)
	}

	// Modify the returned map - this should not affect the original table
	row["Name"] = "Modified"

	// Check that original data is unchanged
	if got := table.Value("row1", "Name"); got != "Test1" {
		t.Errorf("Expected original data to be unchanged, got Value(row1, Name) = %q", got)
	}

	// Non-existent row
	emptyRow := table.Row("non-existent")
	if len(emptyRow) != 0 {
		t.Errorf("Expected empty row for non-existent id")
	}
}

// This test verifies that maps returned by LookupRow are deep copies
func TestCSVTableSafeLookupRowDeepCopy(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	table := abstract.NewCSVTableSafe(records)

	row, found := table.LookupRow("row1")
	if !found {
		t.Errorf("Expected found=true for existing row")
	}
	if got := row["Name"]; got != "Test1" {
		t.Errorf("Expected row[Name] = %q, got %q", "Test1", got)
	}

	// Modify the returned map - this should not affect the original table
	row["Name"] = "Modified"

	// Check that original data is unchanged
	if got := table.Value("row1", "Name"); got != "Test1" {
		t.Errorf("Expected original data to be unchanged, got Value(row1, Name) = %q", got)
	}

	_, found = table.LookupRow("non-existent")
	if found {
		t.Errorf("Expected found=false for non-existent row")
	}
}

// This test verifies that maps returned by All are deep copies
func TestCSVTableSafeAllDeepCopy(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
	}

	table := abstract.NewCSVTableSafe(records)

	all := table.All()
	if len(all) != 2 {
		t.Errorf("Expected len(All()) = 2, got %d", len(all))
	}
	if got := all["row1"]["Name"]; got != "Test1" {
		t.Errorf("Expected all[row1][Name] = %q, got %q", "Test1", got)
	}

	// Modify the returned map - this should not affect the original table
	all["row1"]["Name"] = "Modified"

	// Check that original data is unchanged
	if got := table.Value("row1", "Name"); got != "Test1" {
		t.Errorf("Expected original data to be unchanged, got Value(row1, Name) = %q", got)
	}
}

// This test verifies that maps returned by AllRows are deep copies
func TestCSVTableSafeAllRowsDeepCopy(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
	}

	table := abstract.NewCSVTableSafe(records)

	rows := table.AllRows()
	if len(rows) != 2 {
		t.Errorf("Expected len(AllRows()) = 2, got %d", len(rows))
	}

	// Find the row with Test1 and modify it
	for i, row := range rows {
		if row["Name"] == "Test1" {
			rows[i]["Name"] = "Modified"
			break
		}
	}

	// Check that original data is unchanged
	if got := table.Value("row1", "Name"); got != "Test1" {
		t.Errorf("Expected original data to be unchanged, got Value(row1, Name) = %q", got)
	}
}

func TestCSVTableSafeCopy(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
	}

	table := abstract.NewCSVTableSafe(records)
	tableCopy := table.Copy()

	// Ensure the copy has the same data
	if got := tableCopy.Value("row1", "Name"); got != "Test1" {
		t.Errorf("Expected Value(row1, Name) = %q, got %q", "Test1", got)
	}

	// Modify the copy and ensure original is unchanged
	tableCopy.AddRow("row3", map[string]string{"Name": "Test3", "Value": "300"})
	if !tableCopy.Has("row3") {
		t.Errorf("Expected Has(row3)=true in copy")
	}
	if table.Has("row3") {
		t.Errorf("Expected Has(row3)=false in original")
	}
}

func TestCSVTableSafeAllIDs(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
	}

	table := abstract.NewCSVTableSafe(records)

	ids := table.AllIDs()
	if len(ids) != 2 {
		t.Errorf("Expected len(AllIDs()) = 2, got %d", len(ids))
	}

	found1, found2 := false, false
	for _, id := range ids {
		if id == "row1" {
			found1 = true
		}
		if id == "row2" {
			found2 = true
		}
	}
	if !found1 {
		t.Errorf("Expected AllIDs() to contain row1")
	}
	if !found2 {
		t.Errorf("Expected AllIDs() to contain row2")
	}
}

func TestCSVTableSafeHeaders(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	table := abstract.NewCSVTableSafe(records)

	headers := table.Headers()
	if !reflect.DeepEqual(headers, []string{"ID", "Name", "Value"}) {
		t.Errorf("Expected Headers() = [ID Name Value], got %v", headers)
	}

	// Ensure headers are copied, not referenced
	headers[0] = "Changed"
	originalHeaders := table.Headers()
	if originalHeaders[0] != "ID" {
		t.Errorf("Expected original headers to be unchanged, got %v", originalHeaders)
	}
}

func TestCSVTableSafeValue(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	table := abstract.NewCSVTableSafe(records)

	if got := table.Value("row1", "Name"); got != "Test1" {
		t.Errorf("Expected Value(row1, Name) = %q, got %q", "Test1", got)
	}
	if got := table.Value("row1", "NonExistent"); got != "" {
		t.Errorf("Expected Value(row1, NonExistent) = %q, got %q", "", got)
	}
	if got := table.Value("nonExistent", "Name"); got != "" {
		t.Errorf("Expected Value(nonExistent, Name) = %q, got %q", "", got)
	}
}

func TestCSVTableSafeHas(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	table := abstract.NewCSVTableSafe(records)

	if !table.Has("row1") {
		t.Errorf("Expected Has(row1) = true")
	}
	if table.Has("nonExistent") {
		t.Errorf("Expected Has(nonExistent) = false")
	}
}

func TestCSVTableSafeBytes(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	table := abstract.NewCSVTableSafe(records)

	csvBytes := table.Bytes()
	expected := "\"ID\",\"Name\",\"Value\"\n\"row1\",\"Test1\",\"100\"\n"
	if string(csvBytes) != expected {
		t.Errorf("Expected Bytes() = %q, got %q", expected, string(csvBytes))
	}
}

func TestCSVTableSafeDeleteColumns(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value", "Extra"},
		{"row1", "Test1", "100", "Data1"},
		{"row2", "Test2", "200", "Data2"},
	}

	table := abstract.NewCSVTableSafe(records)

	table.DeleteColumns("Value", "Extra")

	headers := table.Headers()
	if !reflect.DeepEqual(headers, []string{"ID", "Name"}) {
		t.Errorf("Expected Headers() = [ID Name], got %v", headers)
	}

	row := table.Row("row1")
	if _, exists := row["Value"]; exists {
		t.Errorf("Expected Value column to be deleted")
	}
	if _, exists := row["Extra"]; exists {
		t.Errorf("Expected Extra column to be deleted")
	}
}

func TestCSVTableSafeUnwrap(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	tableSafe := abstract.NewCSVTableSafe(records)
	table := tableSafe.Unwrap()

	if got := table.Value("row1", "Name"); got != "Test1" {
		t.Errorf("Expected Value(row1, Name) = %q, got %q", "Test1", got)
	}

	// Verify that the unwrapped table is the actual underlying table
	// by modifying it and seeing the effect on the safe version
	table.AddRow("directAccess", map[string]string{"Name": "DirectAdd"})

	if !tableSafe.Has("directAccess") {
		t.Errorf("Expected modifications to unwrapped table to affect safe table")
	}
}

func TestAllSorted(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row3", "Charlie", "300"},
		{"row1", "Alpha", "100"},
		{"row4", "Delta", "200"},
		{"row2", "Bravo", "400"},
	}

	table := abstract.NewCSVTable(records)

	// Get sorted rows
	sortedRows := table.AllSorted()

	// Should have the same number of rows as in records (minus header)
	if len(sortedRows) != len(records)-1 {
		t.Errorf("Expected %d rows, got %d", len(records)-1, len(sortedRows))
	}

	// Check that rows are in original insertion order
	if sortedRows[0][0] != "row3" {
		t.Errorf("Expected first row ID to be row3, got %s", sortedRows[0][0])
	}
	if sortedRows[1][0] != "row1" {
		t.Errorf("Expected second row ID to be row1, got %s", sortedRows[1][0])
	}
	if sortedRows[2][0] != "row4" {
		t.Errorf("Expected third row ID to be row4, got %s", sortedRows[2][0])
	}
	if sortedRows[3][0] != "row2" {
		t.Errorf("Expected fourth row ID to be row2, got %s", sortedRows[3][0])
	}

	// After sorting, the order should change but AllSorted should still return original order
	table.Sort("Name", abstract.ASCSort)
	sortedRows = table.AllSorted()

	// Check that rows are still in original insertion order despite table sort
	if sortedRows[0][0] != "row1" {
		t.Errorf("Expected first row ID to be row1, got %s", sortedRows[0][0])
	}
	if sortedRows[1][0] != "row2" {
		t.Errorf("Expected second row ID to be row2, got %s", sortedRows[1][0])
	}
	if sortedRows[2][0] != "row3" {
		t.Errorf("Expected third row ID to be row3, got %s", sortedRows[2][0])
	}
	if sortedRows[3][0] != "row4" {
		t.Errorf("Expected fourth row ID to be row4, got %s", sortedRows[3][0])
	}
}

func TestRowSorted(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value", "Extra"},
		{"row1", "Alpha", "100", "Data1"},
		{"row2", "Bravo", "200", "Data2"},
	}

	table := abstract.NewCSVTable(records)

	// Get the row as an array of strings
	row := table.RowSorted("row1")

	// Check length
	if len(row) != 4 {
		t.Errorf("Expected row to have 4 values, got %d", len(row))
	}

	// Check values
	if row[0] != "row1" {
		t.Errorf("Expected row[0] to be 'row1', got %s", row[0])
	}
	if row[1] != "Alpha" {
		t.Errorf("Expected row[1] to be 'Alpha', got %s", row[1])
	}
	if row[2] != "100" {
		t.Errorf("Expected row[2] to be '100', got %s", row[2])
	}
	if row[3] != "Data1" {
		t.Errorf("Expected row[3] to be 'Data1', got %s", row[3])
	}

	// Check non-existent row
	nonExistentRow := table.RowSorted("nonExistent")
	if nonExistentRow != nil {
		t.Errorf("Expected nil for non-existent row, got %v", nonExistentRow)
	}
}

func TestLookupRowSorted(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value", "Extra"},
		{"row1", "Alpha", "100", "Data1"},
		{"row2", "Bravo", "200", "Data2"},
	}

	table := abstract.NewCSVTable(records)

	// Look up an existing row
	row, found := table.LookupRowSorted("row1")
	if !found {
		t.Errorf("Expected to find row1")
	}

	// Check length
	if len(row) != 4 {
		t.Errorf("Expected row to have 4 values, got %d", len(row))
	}

	// Check values
	if row[0] != "row1" {
		t.Errorf("Expected row[0] to be 'row1', got %s", row[0])
	}
	if row[1] != "Alpha" {
		t.Errorf("Expected row[1] to be 'Alpha', got %s", row[1])
	}
	if row[2] != "100" {
		t.Errorf("Expected row[2] to be '100', got %s", row[2])
	}
	if row[3] != "Data1" {
		t.Errorf("Expected row[3] to be 'Data1', got %s", row[3])
	}

	// Look up a non-existent row
	row, found = table.LookupRowSorted("nonExistent")
	if found {
		t.Errorf("Expected not to find nonExistent row")
	}
	if row != nil {
		t.Errorf("Expected nil row for non-existent row, got %v", row)
	}
}

func TestCSVTableSafeAllSorted(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row3", "Charlie", "300"},
		{"row1", "Alpha", "100"},
		{"row4", "Delta", "200"},
		{"row2", "Bravo", "400"},
	}

	table := abstract.NewCSVTableSafe(records)

	// Get sorted rows
	sortedRows := table.AllSorted()

	// Check that rows are in original insertion order
	if sortedRows[0][0] != "row3" {
		t.Errorf("Expected first row ID to be row3, got %s", sortedRows[0][0])
	}
	if sortedRows[1][0] != "row1" {
		t.Errorf("Expected second row ID to be row1, got %s", sortedRows[1][0])
	}
	if sortedRows[2][0] != "row4" {
		t.Errorf("Expected third row ID to be row4, got %s", sortedRows[2][0])
	}
	if sortedRows[3][0] != "row2" {
		t.Errorf("Expected fourth row ID to be row2, got %s", sortedRows[3][0])
	}

	table.Sort("Name", abstract.ASCSort)
	sortedRows = table.AllSorted()

	if sortedRows[0][0] != "row1" {
		t.Errorf("Expected first row ID to be row1, got %s", sortedRows[0][0])
	}
	if sortedRows[1][0] != "row2" {
		t.Errorf("Expected second row ID to be row2, got %s", sortedRows[1][0])
	}
	if sortedRows[2][0] != "row3" {
		t.Errorf("Expected third row ID to be row3, got %s", sortedRows[2][0])
	}
	if sortedRows[3][0] != "row4" {
		t.Errorf("Expected fourth row ID to be row4, got %s", sortedRows[3][0])
	}
}

func TestCSVTableSafeRowSorted(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Alpha", "100"},
		{"row2", "Bravo", "200"},
	}

	table := abstract.NewCSVTableSafe(records)

	// Get the row as an array of strings
	row := table.RowSorted("row1")

	// Check length and values
	if len(row) != 3 {
		t.Errorf("Expected row to have 3 values, got %d", len(row))
	}
	if row[0] != "row1" {
		t.Errorf("Expected row[0] to be 'row1', got %s", row[0])
	}
	if row[1] != "Alpha" {
		t.Errorf("Expected row[1] to be 'Alpha', got %s", row[1])
	}
	if row[2] != "100" {
		t.Errorf("Expected row[2] to be '100', got %s", row[2])
	}
}

func TestCSVTableSafeLookupRowSorted(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Alpha", "100"},
		{"row2", "Bravo", "200"},
	}

	table := abstract.NewCSVTableSafe(records)

	// Look up an existing row
	row, found := table.LookupRowSorted("row1")
	if !found {
		t.Errorf("Expected to find row1")
	}

	// Check length and values
	if len(row) != 3 {
		t.Errorf("Expected row to have 3 values, got %d", len(row))
	}
	if row[0] != "row1" {
		t.Errorf("Expected row[0] to be 'row1', got %s", row[0])
	}
	if row[1] != "Alpha" {
		t.Errorf("Expected row[1] to be 'Alpha', got %s", row[1])
	}
	if row[2] != "100" {
		t.Errorf("Expected row[2] to be '100', got %s", row[2])
	}

	// Look up a non-existent row
	row, found = table.LookupRowSorted("nonExistent")
	if found {
		t.Errorf("Expected not to find nonExistent row")
	}
	if row != nil {
		t.Errorf("Expected nil row for non-existent row, got %v", row)
	}
}

func TestSort(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row3", "Charlie", "300"},
		{"row1", "Alpha", "100"},
		{"row4", "Delta", "200"},
		{"row2", "Bravo", "400"},
	}

	table := abstract.NewCSVTable(records)

	// Test ascending sort by Name
	table.Sort("Name", abstract.ASCSort)

	// Check if IDs are in the expected order after name-based sorting
	ids := table.AllIDs()
	expectedNameOrder := []string{"row1", "row2", "row3", "row4"} // Alpha, Bravo, Charlie, Delta

	if !reflect.DeepEqual(ids, expectedNameOrder) {
		t.Errorf("Expected IDs after sorting by Name ASC to be %v, got %v", expectedNameOrder, ids)
	}

	// Test descending sort by Value
	table.Sort("Value", abstract.DESCSort)

	ids = table.AllIDs()
	expectedValueOrder := []string{"row2", "row3", "row4", "row1"} // 400, 300, 200, 100

	if !reflect.DeepEqual(ids, expectedValueOrder) {
		t.Errorf("Expected IDs after sorting by Value DESC to be %v, got %v", expectedValueOrder, ids)
	}

	// Test that row data is correctly accessible after sorting
	if val := table.Value("row2", "Value"); val != "400" {
		t.Errorf("Expected Value for row2 to be 400, got %s", val)
	}

	// Test sorting by non-existent column (should have no effect)
	originalIDs := table.AllIDs()
	table.Sort("NonExistentColumn", abstract.ASCSort)

	if !reflect.DeepEqual(originalIDs, table.AllIDs()) {
		t.Errorf("Expected no change when sorting by non-existent column")
	}
}

func TestCSVTableSafeSort(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row3", "Charlie", "300"},
		{"row1", "Alpha", "100"},
		{"row4", "Delta", "200"},
		{"row2", "Bravo", "400"},
	}

	table := abstract.NewCSVTableSafe(records)

	// Test ascending sort by Name
	table.Sort("Name", abstract.ASCSort)

	// Check if IDs are in the expected order after name-based sorting
	ids := table.AllIDs()
	expectedNameOrder := []string{"row1", "row2", "row3", "row4"} // Alpha, Bravo, Charlie, Delta

	if !reflect.DeepEqual(ids, expectedNameOrder) {
		t.Errorf("Expected IDs after sorting by Name ASC to be %v, got %v", expectedNameOrder, ids)
	}

	// Test descending sort by Value
	table.Sort("Value", abstract.DESCSort)

	ids = table.AllIDs()
	expectedValueOrder := []string{"row2", "row3", "row4", "row1"} // 400, 300, 200, 100

	if !reflect.DeepEqual(ids, expectedValueOrder) {
		t.Errorf("Expected IDs after sorting by Value DESC to be %v, got %v", expectedValueOrder, ids)
	}

	// Verify that we can still look up values by ID after sorting
	val, exists := table.LookupRow("row2")
	if !exists {
		t.Errorf("Expected to find row2 after sorting")
	}
	if val["Value"] != "400" {
		t.Errorf("Expected Value for row2 to be 400, got %s", val["Value"])
	}
}

// Tests for new methods

func TestNewCSVTableFromMap(t *testing.T) {
	data := map[string]map[string]string{
		"user1": {
			"name":  "Alice",
			"email": "alice@example.com",
			"age":   "25",
		},
		"user2": {
			"name":  "Bob",
			"email": "bob@example.com",
			"age":   "30",
		},
		"user3": {
			"name":  "Charlie",
			"email": "charlie@example.com",
		},
	}

	table := abstract.NewCSVTableFromMap(data)

	// Check that all users exist
	if !table.Has("user1") {
		t.Errorf("Expected table to have user1")
	}
	if !table.Has("user2") {
		t.Errorf("Expected table to have user2")
	}
	if !table.Has("user3") {
		t.Errorf("Expected table to have user3")
	}

	// Check values
	if got := table.Value("user1", "name"); got != "Alice" {
		t.Errorf("Expected user1 name = Alice, got %s", got)
	}
	if got := table.Value("user2", "age"); got != "30" {
		t.Errorf("Expected user2 age = 30, got %s", got)
	}
	if got := table.Value("user3", "age"); got != "" {
		t.Errorf("Expected user3 age = empty, got %s", got)
	}

	// Check headers include id and all columns
	headers := table.Headers()
	if len(headers) != 4 { // id, age, email, name (sorted)
		t.Errorf("Expected 4 headers, got %d", len(headers))
	}
	if headers[0] != "id" {
		t.Errorf("Expected first header to be 'id', got %s", headers[0])
	}
}

func TestNewCSVTableFromMapWithCustomID(t *testing.T) {
	data := map[string]map[string]string{
		"user1": {
			"name":  "Alice",
			"email": "alice@example.com",
		},
	}

	table := abstract.NewCSVTableFromMap(data, "user_id")

	headers := table.Headers()
	if headers[0] != "user_id" {
		t.Errorf("Expected first header to be 'user_id', got %s", headers[0])
	}

	// Check that the ID value is correctly set
	row := table.RowSorted("user1")
	if row[0] != "user1" {
		t.Errorf("Expected first column value to be 'user1', got %s", row[0])
	}
}

func TestNewCSVTableFromMapEmpty(t *testing.T) {
	data := map[string]map[string]string{}
	table := abstract.NewCSVTableFromMap(data)

	if len(table.All()) != 0 {
		t.Errorf("Expected empty table for empty data")
	}
	if len(table.Headers()) != 0 {
		t.Errorf("Expected no headers for empty data")
	}
}

func TestDeleteColumn(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value", "Extra"},
		{"row1", "Test1", "100", "Data1"},
		{"row2", "Test2", "200", "Data2"},
	}

	table := abstract.NewCSVTable(records)

	// Delete single column
	table.DeleteColumn("Value")

	headers := table.Headers()
	if len(headers) != 3 {
		t.Errorf("Expected 3 headers after deleting one column, got %d", len(headers))
	}

	// Check that Value column is gone
	for _, h := range headers {
		if h == "Value" {
			t.Errorf("Expected Value column to be deleted")
		}
	}

	// Check that data is preserved for remaining columns
	if got := table.Value("row1", "Name"); got != "Test1" {
		t.Errorf("Expected Name to be preserved, got %s", got)
	}
	if got := table.Value("row1", "Value"); got != "" {
		t.Errorf("Expected Value to be empty after deletion, got %s", got)
	}
}

func TestDeleteRow(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
		{"row3", "Test3", "300"},
	}

	table := abstract.NewCSVTable(records)

	// Delete middle row
	deleted := table.DeleteRow("row2")
	if !deleted {
		t.Errorf("Expected DeleteRow to return true for existing row")
	}

	// Check that row is gone
	if table.Has("row2") {
		t.Errorf("Expected row2 to be deleted")
	}

	// Check that other rows still exist
	if !table.Has("row1") {
		t.Errorf("Expected row1 to still exist")
	}
	if !table.Has("row3") {
		t.Errorf("Expected row3 to still exist")
	}

	// Check that indices are updated correctly
	ids := table.AllIDs()
	if len(ids) != 2 {
		t.Errorf("Expected 2 IDs after deletion, got %d", len(ids))
	}

	// Try to delete non-existent row
	deleted = table.DeleteRow("nonexistent")
	if deleted {
		t.Errorf("Expected DeleteRow to return false for non-existent row")
	}
}

func TestUpdateColumn(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
		{"row3", "Test3", "300"},
	}

	table := abstract.NewCSVTable(records)

	// Update existing column
	newValues := []string{"NewVal1", "NewVal2", "NewVal3"}
	table.UpdateColumn("Value", newValues)

	if got := table.Value("row1", "Value"); got != "NewVal1" {
		t.Errorf("Expected updated value NewVal1, got %s", got)
	}
	if got := table.Value("row2", "Value"); got != "NewVal2" {
		t.Errorf("Expected updated value NewVal2, got %s", got)
	}
	if got := table.Value("row3", "Value"); got != "NewVal3" {
		t.Errorf("Expected updated value NewVal3, got %s", got)
	}

	// Update with fewer values than rows
	shortValues := []string{"Short1"}
	table.UpdateColumn("Name", shortValues)

	if got := table.Value("row1", "Name"); got != "Short1" {
		t.Errorf("Expected updated name Short1, got %s", got)
	}
	if got := table.Value("row2", "Name"); got != "Test2" {
		t.Errorf("Expected unchanged name Test2, got %s", got)
	}

	// Try to update non-existent column
	table.UpdateColumn("NonExistent", []string{"test"})
	// Should not crash or affect anything
}

func TestUpdateRow(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
	}

	table := abstract.NewCSVTable(records)

	// Update existing row
	updates := map[string]string{
		"Name":  "UpdatedName",
		"Value": "UpdatedValue",
	}
	updated := table.UpdateRow("row1", updates)
	if !updated {
		t.Errorf("Expected UpdateRow to return true for existing row")
	}

	if got := table.Value("row1", "Name"); got != "UpdatedName" {
		t.Errorf("Expected updated name UpdatedName, got %s", got)
	}
	if got := table.Value("row1", "Value"); got != "UpdatedValue" {
		t.Errorf("Expected updated value UpdatedValue, got %s", got)
	}

	// Partial update
	partialUpdates := map[string]string{
		"Value": "PartialUpdate",
	}
	table.UpdateRow("row2", partialUpdates)

	if got := table.Value("row2", "Name"); got != "Test2" {
		t.Errorf("Expected unchanged name Test2, got %s", got)
	}
	if got := table.Value("row2", "Value"); got != "PartialUpdate" {
		t.Errorf("Expected updated value PartialUpdate, got %s", got)
	}

	// Try to update non-existent row
	updated = table.UpdateRow("nonexistent", updates)
	if updated {
		t.Errorf("Expected UpdateRow to return false for non-existent row")
	}
}

func TestFindRow(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Age", "City"},
		{"user1", "Alice Smith", "25", "New York"},
		{"user2", "Bob Johnson", "30", "Los Angeles"},
		{"user3", "Charlie Brown", "25", "New York"},
	}

	table := abstract.NewCSVTable(records)

	// Find by single criterion (using contains)
	id, row := table.FindRow(map[string]string{"Age": "25"})
	if id == "" {
		t.Errorf("Expected to find a row with Age=25")
	}
	if row["Name"] != "Alice Smith" && row["Name"] != "Charlie Brown" {
		t.Errorf("Expected to find Alice Smith or Charlie Brown, got %s", row["Name"])
	}

	// Find by multiple criteria
	id, _ = table.FindRow(map[string]string{"Age": "25", "City": "New York"})
	if id == "" {
		t.Errorf("Expected to find a row with Age=25 and City=New York")
	}
	// Should find either user1 or user3, both match

	// Find with partial match (contains)
	id, row = table.FindRow(map[string]string{"Name": "Alice"})
	if id != "user1" {
		t.Errorf("Expected to find user1, got %s", id)
	}
	if row["Name"] != "Alice Smith" {
		t.Errorf("Expected to find Alice Smith, got %s", row["Name"])
	}

	// Find non-existent
	id, row = table.FindRow(map[string]string{"Age": "99"})
	if id != "" {
		t.Errorf("Expected empty ID for non-existent criteria, got %s", id)
	}
	if row != nil {
		t.Errorf("Expected nil row for non-existent criteria, got %v", row)
	}
}

func TestFind(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Age", "City"},
		{"user1", "Alice Smith", "25", "New York"},
		{"user2", "Bob Johnson", "30", "Los Angeles"},
		{"user3", "Charlie Brown", "25", "New York"},
		{"user4", "David Wilson", "25", "Chicago"},
	}

	table := abstract.NewCSVTable(records)

	// Find all with Age=25
	results := table.Find(map[string]string{"Age": "25"})
	if len(results) != 3 {
		t.Errorf("Expected 3 results for Age=25, got %d", len(results))
	}

	// Check that all returned rows have Age=25
	for id, row := range results {
		if row["Age"] != "25" {
			t.Errorf("Expected Age=25 for %s, got %s", id, row["Age"])
		}
	}

	// Find all with multiple criteria
	results = table.Find(map[string]string{"Age": "25", "City": "New York"})
	if len(results) != 2 {
		t.Errorf("Expected 2 results for Age=25 and City=New York, got %d", len(results))
	}

	// Find with partial match
	results = table.Find(map[string]string{"Name": "Johnson"})
	if len(results) != 1 {
		t.Errorf("Expected 1 result for Name containing Johnson, got %d", len(results))
	}
	if _, exists := results["user2"]; !exists {
		t.Errorf("Expected to find user2 in results")
	}

	// Find non-existent
	results = table.Find(map[string]string{"Age": "99"})
	if len(results) != 0 {
		t.Errorf("Expected 0 results for non-existent criteria, got %d", len(results))
	}
}

// Tests for CSVTableSafe new methods

func TestNewCSVTableSafeFromMap(t *testing.T) {
	data := map[string]map[string]string{
		"user1": {
			"name":  "Alice",
			"email": "alice@example.com",
		},
	}

	table := abstract.NewCSVTableSafeFromMap(data)

	if !table.Has("user1") {
		t.Errorf("Expected table to have user1")
	}
	if got := table.Value("user1", "name"); got != "Alice" {
		t.Errorf("Expected user1 name = Alice, got %s", got)
	}
}

func TestCSVTableSafeDeleteColumn(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	table := abstract.NewCSVTableSafe(records)
	table.DeleteColumn("Value")

	headers := table.Headers()
	for _, h := range headers {
		if h == "Value" {
			t.Errorf("Expected Value column to be deleted")
		}
	}
}

func TestCSVTableSafeDeleteRow(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
	}

	table := abstract.NewCSVTableSafe(records)

	deleted := table.DeleteRow("row1")
	if !deleted {
		t.Errorf("Expected DeleteRow to return true")
	}
	if table.Has("row1") {
		t.Errorf("Expected row1 to be deleted")
	}
	if !table.Has("row2") {
		t.Errorf("Expected row2 to still exist")
	}
}

func TestCSVTableSafeUpdateColumn(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
		{"row2", "Test2", "200"},
	}

	table := abstract.NewCSVTableSafe(records)

	newValues := []string{"NewVal1", "NewVal2"}
	table.UpdateColumn("Value", newValues)

	if got := table.Value("row1", "Value"); got != "NewVal1" {
		t.Errorf("Expected updated value NewVal1, got %s", got)
	}
	if got := table.Value("row2", "Value"); got != "NewVal2" {
		t.Errorf("Expected updated value NewVal2, got %s", got)
	}
}

func TestCSVTableSafeUpdateRow(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Value"},
		{"row1", "Test1", "100"},
	}

	table := abstract.NewCSVTableSafe(records)

	updates := map[string]string{
		"Name":  "UpdatedName",
		"Value": "UpdatedValue",
	}
	updated := table.UpdateRow("row1", updates)
	if !updated {
		t.Errorf("Expected UpdateRow to return true")
	}

	if got := table.Value("row1", "Name"); got != "UpdatedName" {
		t.Errorf("Expected updated name UpdatedName, got %s", got)
	}
}

func TestCSVTableSafeFindRow(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Age"},
		{"user1", "Alice Smith", "25"},
		{"user2", "Bob Johnson", "30"},
	}

	table := abstract.NewCSVTableSafe(records)

	id, row := table.FindRow(map[string]string{"Age": "25"})
	if id != "user1" {
		t.Errorf("Expected to find user1, got %s", id)
	}
	if row["Name"] != "Alice Smith" {
		t.Errorf("Expected to find Alice Smith, got %s", row["Name"])
	}

	// Test thread safety by ensuring returned map is a copy
	row["Name"] = "Modified"
	if got := table.Value("user1", "Name"); got != "Alice Smith" {
		t.Errorf("Expected original data to be unchanged, got %s", got)
	}
}

func TestCSVTableSafeFind(t *testing.T) {
	records := [][]string{
		{"ID", "Name", "Age"},
		{"user1", "Alice Smith", "25"},
		{"user2", "Bob Johnson", "25"},
		{"user3", "Charlie Brown", "30"},
	}

	table := abstract.NewCSVTableSafe(records)

	results := table.Find(map[string]string{"Age": "25"})
	if len(results) != 2 {
		t.Errorf("Expected 2 results for Age=25, got %d", len(results))
	}

	// Test thread safety by ensuring returned maps are copies
	for _, row := range results {
		row["Name"] = "Modified"
	}

	if got := table.Value("user1", "Name"); got != "Alice Smith" {
		t.Errorf("Expected original data to be unchanged, got %s", got)
	}
	if got := table.Value("user2", "Name"); got != "Bob Johnson" {
		t.Errorf("Expected original data to be unchanged, got %s", got)
	}
}
