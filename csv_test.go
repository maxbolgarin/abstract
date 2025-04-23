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
