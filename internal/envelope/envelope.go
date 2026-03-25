package envelope

import "fmt"

// Envelope is the standard response wrapper.
type Envelope struct {
	OK          bool                   `json:"ok"`
	Data        interface{}            `json:"data,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Summary     string                 `json:"summary,omitempty"`
	Breadcrumbs []string               `json:"breadcrumbs,omitempty"`
}

// Wrap creates a success envelope.
func Wrap(data interface{}, summary string, breadcrumbs []string) *Envelope {
	return &Envelope{
		OK:          true,
		Data:        data,
		Summary:     summary,
		Breadcrumbs: breadcrumbs,
	}
}

// WrapError creates an error envelope.
func WrapError(message string, breadcrumbs []string) *Envelope {
	return &Envelope{
		OK:          false,
		Error:       message,
		Breadcrumbs: breadcrumbs,
	}
}

// WrapCollection wraps an index/list API response.
// The API typically returns { "data": [...], "total": N }.
func WrapCollection(response map[string]interface{}, resourceName string, breadcrumbs []string) *Envelope {
	data := extractData(response)
	total := extractTotal(response)

	var summary string
	records, ok := data.([]interface{})
	if ok {
		if total >= 0 {
			summary = fmt.Sprintf("Showing %d of %d %s", len(records), total, resourceName)
		} else {
			summary = fmt.Sprintf("Found %d %s", len(records), resourceName)
		}
	} else {
		summary = fmt.Sprintf("Found %s", resourceName)
	}

	return Wrap(data, summary, breadcrumbs)
}

// WrapRecord wraps a single record API response.
func WrapRecord(record map[string]interface{}, resourceName string, breadcrumbs []string) *Envelope {
	id := ""
	if v, ok := record["_id"]; ok {
		id = fmt.Sprintf("%v", v)
	} else if v, ok := record["id"]; ok {
		id = fmt.Sprintf("%v", v)
	}

	summary := fmt.Sprintf("%s %s", resourceName, id)
	return Wrap(record, summary, breadcrumbs)
}

// ExtractID gets the _id or id field from a record.
func ExtractID(record map[string]interface{}) string {
	if v, ok := record["_id"]; ok {
		return fmt.Sprintf("%v", v)
	}
	if v, ok := record["id"]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func extractData(response map[string]interface{}) interface{} {
	if d, ok := response["data"]; ok {
		return d
	}
	return response
}

func extractTotal(response map[string]interface{}) int {
	if t, ok := response["total"]; ok {
		switch v := t.(type) {
		case float64:
			return int(v)
		case int:
			return v
		}
	}
	return -1
}
