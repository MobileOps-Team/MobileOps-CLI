package commands

var apiCoverage = map[string][]string{
	"vessels list":   {"GET /assets"},
	"vessels get":    {"GET /assets/{id}"},
	"vessels create": {"POST /assets"},
	"vessels update": {"PUT /assets/{id}"},

	"jobs list":                {"GET /jobs"},
	"jobs get":                 {"GET /jobs/{id}"},
	"jobs create":              {"POST /jobs"},
	"jobs update":              {"PUT /jobs/{id}"},
	"jobs vessel-availability": {"POST /jobs/vessel-availability"},
	"jobs export":              {"POST /jobs/job-export"},

	"crew list":                    {"GET /users"},
	"crew get":                     {"GET /users/{id}"},
	"crew create":                  {"POST /users"},
	"crew update":                  {"PUT /users/{id}"},
	"crew delete":                  {"DELETE /users/{id}"},
	"crew find-by-employee-number": {"GET /users/by_employee_number"},

	"components list":   {"GET /components"},
	"components get":    {"GET /components/{id}"},
	"components create": {"POST /components"},
	"components update": {"PUT /components/{id}"},

	"parts list":   {"GET /parts"},
	"parts get":    {"GET /parts/{id}"},
	"parts create": {"POST /parts"},
	"parts update": {"PUT /parts/{id}"},

	"work-requests list":   {"GET /work-requests"},
	"work-requests get":    {"GET /work-requests/{id}"},
	"work-requests create": {"POST /work-requests"},
	"work-requests update": {"PUT /work-requests/{id}"},

	"deficiencies list":   {"GET /deficiencies"},
	"deficiencies get":    {"GET /deficiencies/{id}"},
	"deficiencies create": {"POST /deficiencies"},
	"deficiencies update": {"PUT /deficiencies/{id}"},

	"nonconformities list":   {"GET /nonconformities"},
	"nonconformities get":    {"GET /nonconformities/{id}"},
	"nonconformities create": {"POST /nonconformities"},
	"nonconformities update": {"PUT /nonconformities/{id}"},

	"observations list":   {"GET /observations"},
	"observations get":    {"GET /observations/{id}"},
	"observations create": {"POST /observations"},
	"observations update": {"PUT /observations/{id}"},

	"maintenance-reports list":   {"GET /maintenance-reports"},
	"maintenance-reports get":    {"GET /maintenance-reports/{id}"},
	"maintenance-reports create": {"POST /maintenance-reports"},
	"maintenance-reports update": {"PUT /maintenance-reports/{id}"},

	"vessel-documents list":        {"GET /vessel-documents"},
	"vessel-documents get":         {"GET /vessel-documents/{id}"},
	"vessel-documents create":      {"POST /vessel-documents"},
	"vessel-documents update":      {"PUT /vessel-documents/{id}"},
	"vessel-documents delete":      {"DELETE /vessel-documents/{id}"},
	"vessel-documents attachments": {"GET /vessel-documents/{id}/attachments"},

	"personnel-documents list":        {"GET /personnel-documents"},
	"personnel-documents get":         {"GET /personnel-documents/{id}"},
	"personnel-documents create":      {"POST /personnel-documents"},
	"personnel-documents update":      {"PUT /personnel-documents/{id}"},
	"personnel-documents delete":      {"DELETE /personnel-documents/{id}"},
	"personnel-documents attachments": {"GET /personnel-documents/{id}/attachments"},

	"vessel-specs list":   {"GET /vessel-specs"},
	"vessel-specs get":    {"GET /vessel-specs/{id}"},
	"vessel-specs create": {"POST /vessel-specs"},
	"vessel-specs update": {"PUT /vessel-specs/{id}"},

	"vessel-spec-templates list": {"GET /vessel-spec-templates"},
	"vessel-spec-templates get":  {"GET /vessel-spec-templates/{id}"},

	"customers list":   {"GET /customers"},
	"customers get":    {"GET /customers/{id}"},
	"customers create": {"POST /customers"},
	"customers update": {"PUT /customers/{id}"},

	"suppliers list": {"GET /suppliers"},
	"suppliers get":  {"GET /suppliers/{id}"},

	"makes list":   {"GET /makes"},
	"makes get":    {"GET /makes/{id}"},
	"makes create": {"POST /makes"},

	"models list":   {"GET /models"},
	"models get":    {"GET /models/{id}"},
	"models create": {"POST /models"},

	"terms list":   {"GET /terms"},
	"terms get":    {"GET /terms/{id}"},
	"terms create": {"POST /terms"},

	"divisions list":   {"GET /divisions"},
	"divisions get":    {"GET /divisions/{id}"},
	"divisions create": {"POST /divisions"},
	"divisions update": {"PUT /divisions/{id}"},

	"employee-positions list":   {"GET /employee-positions"},
	"employee-positions get":    {"GET /employee-positions/{id}"},
	"employee-positions create": {"POST /employee-positions"},
	"employee-positions update": {"PUT /employee-positions/{id}"},

	"vessel-types list":   {"GET /vessel-types"},
	"vessel-types get":    {"GET /vessel-types/{id}"},
	"vessel-types create": {"POST /vessel-types"},
	"vessel-types update": {"PUT /vessel-types/{id}"},

	"work-types list":   {"GET /work-types"},
	"work-types get":    {"GET /work-types/{id}"},
	"work-types create": {"POST /work-types"},
	"work-types update": {"PUT /work-types/{id}"},

	"locations list":   {"GET /locations"},
	"locations get":    {"GET /locations/{id}"},
	"locations create": {"POST /locations"},
	"locations update": {"PUT /locations/{id}"},

	"codes list":   {"GET /codes"},
	"codes get":    {"GET /codes/{id}"},
	"codes create": {"POST /codes"},
	"codes update": {"PUT /codes/{id}"},

	"cargo-types list": {"GET /cargo-types"},
	"cargo-types get":  {"GET /cargo-types/{id}"},

	"invoice-statements list":   {"GET /invoice-statements"},
	"invoice-statements update": {"PUT /invoice-statements/{id}"},

	"forms list":                {"GET /forms"},
	"form-instances list":       {"GET /form-instances"},
	"routine-templates list":    {"GET /routine-templates"},
	"routine-calculations list": {"GET /routine-calculations"},
	"component-logs list":       {"GET /component-logs"},
	"part-requests list":        {"GET /part-requests"},
	"part-requests get":         {"GET /part-requests/{id}"},
	"purchase-orders list":      {"GET /purchase-orders"},
	"wheelhouse-logs list":      {"GET /wheelhouse-logs"},
	"audits list":               {"GET /audits"},
	"work-rests list":           {"GET /work-rests"},
	"events search":             {"POST /events/search"},
	"fuel-level-readings list":  {"GET /fuel-level-readings"},
	"bunker-partitions list":    {"GET /bunker-partitions"},

	"position-reports list":   {"GET /position-reports"},
	"position-reports create": {"POST /position-reports"},
}

// apiIgnored lists spec operations the CLI deliberately does not expose, with
// the reason. The coverage test will not ask for a command for these.
var apiIgnored = map[string]string{
	"POST /position-reports/particle": "webhook for Particle.io IoT trackers, not a user-facing call",
}

// nonAPICommands are CLI commands that do not map to a single REST operation.
var nonAPICommands = map[string]bool{
	"auth login":  true,
	"auth logout": true,
	"auth status": true,
	"version":     true,
	"update":      true,
	"tree":        true,
}

// apiFieldIgnored lists contract fields the CLI deliberately has no flag for,
// keyed "VERB /path field" (body fields as "wrapper.field", query params as
// "?name"), each with the reason. TestAPIFieldCoverage skips these.
var apiFieldIgnored = map[string]string{
	"POST /personnel-documents personnel_document.issue_date_string":      "free-text variant of issue_date that bypasses date validation; use --issue-date",
	"POST /personnel-documents personnel_document.expire_date_string":     "free-text variant of expire_date that bypasses date validation; use --expire-date",
	"PUT /personnel-documents/{id} personnel_document.issue_date_string":  "free-text variant of issue_date that bypasses date validation; use --issue-date",
	"PUT /personnel-documents/{id} personnel_document.expire_date_string": "free-text variant of expire_date that bypasses date validation; use --expire-date",
}
