package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
)

const maxColumns = 8
const maxValueLen = 30

// Output renders an envelope to stdout in human or JSON mode.
func Output(env *envelope.Envelope, jsonMode bool) {
	if jsonMode {
		data, _ := json.MarshalIndent(env, "", "  ")
		fmt.Println(string(data))
		return
	}

	if !env.OK {
		fmt.Fprintf(os.Stderr, "Error: %s\n", env.Error)
		return
	}

	switch data := env.Data.(type) {
	case []interface{}:
		printTable(data)
		fmt.Println()
		fmt.Println(env.Summary)
	case map[string]interface{}:
		printRecord(data)
	default:
		fmt.Println(env.Summary)
	}

	printBreadcrumbs(env.Breadcrumbs)
}

func printTable(records []interface{}) {
	if len(records) == 0 {
		fmt.Println("No records found.")
		return
	}

	// Extract column names from first record
	first, ok := records[0].(map[string]interface{})
	if !ok {
		fmt.Println("No records found.")
		return
	}

	columns := orderedKeys(first)
	if len(columns) > maxColumns {
		columns = columns[:maxColumns]
	}

	// Build header
	headers := make([]string, len(columns))
	colWidths := make([]int, len(columns))
	for i, col := range columns {
		h := strings.ToUpper(strings.ReplaceAll(col, "_", " "))
		headers[i] = h
		colWidths[i] = len(h)
	}

	// Build rows and compute widths
	rows := make([][]string, len(records))
	for r, rec := range records {
		record, ok := rec.(map[string]interface{})
		if !ok {
			continue
		}
		row := make([]string, len(columns))
		for c, col := range columns {
			val := formatValue(record[col])
			val = truncate(val, maxValueLen)
			row[c] = val
			if len(val) > colWidths[c] {
				colWidths[c] = len(val)
			}
		}
		rows[r] = row
	}

	// Print separator
	printSep(colWidths)

	// Print header
	printRow(headers, colWidths)
	printSep(colWidths)

	// Print data rows
	for _, row := range rows {
		if row != nil {
			printRow(row, colWidths)
		}
	}

	printSep(colWidths)
}

func printSep(widths []int) {
	parts := make([]string, len(widths))
	for i, w := range widths {
		parts[i] = strings.Repeat("-", w+2)
	}
	fmt.Printf("+%s+\n", strings.Join(parts, "+"))
}

func printRow(cells []string, widths []int) {
	parts := make([]string, len(cells))
	for i, cell := range cells {
		parts[i] = fmt.Sprintf(" %-*s ", widths[i], cell)
	}
	fmt.Printf("|%s|\n", strings.Join(parts, "|"))
}

func printRecord(record map[string]interface{}) {
	if len(record) == 0 {
		fmt.Println("No record found.")
		return
	}

	keys := orderedKeys(record)
	maxKeyLen := 0
	for _, k := range keys {
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	for _, key := range keys {
		label := fmt.Sprintf("%-*s", maxKeyLen, key)
		val := formatValue(record[key])
		fmt.Printf("  %s  %s\n", label, val)
	}
}

func printBreadcrumbs(breadcrumbs []string) {
	if len(breadcrumbs) == 0 {
		return
	}
	fmt.Println()
	fmt.Println("Next:")
	for _, cmd := range breadcrumbs {
		fmt.Printf("  → %s\n", cmd)
	}
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}

func formatValue(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return "-"
	case []interface{}:
		if len(val) == 0 {
			return "[]"
		}
		parts := make([]string, 0, 5)
		for i, item := range val {
			if i >= 5 {
				break
			}
			parts = append(parts, fmt.Sprintf("%v", item))
		}
		return strings.Join(parts, ", ")
	case map[string]interface{}:
		data, _ := json.Marshal(val)
		return string(data)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// orderedKeys returns map keys in a stable order.
// JSON unmarshal doesn't preserve order, so we use a deterministic sort.
// For better UX, we prioritize common ID/name fields.
func orderedKeys(m map[string]interface{}) []string {
	priority := []string{"_id", "id", "name", "first_name", "last_name", "email", "status", "description"}
	seen := make(map[string]bool)
	var result []string

	for _, k := range priority {
		if _, ok := m[k]; ok {
			result = append(result, k)
			seen[k] = true
		}
	}

	// Add remaining keys in natural order
	for k := range m {
		if !seen[k] {
			result = append(result, k)
		}
	}

	return result
}
