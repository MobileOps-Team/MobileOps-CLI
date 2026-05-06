---
name: mobileops
description: Manage maritime operations — vessels, crew, jobs, maintenance, inspections, compliance, and more — via the MobileOps CLI.
---

# MobileOps CLI Skill

MobileOps is a maritime operations platform for managing vessels, crew, jobs, maintenance, inspections, and compliance. This CLI lets you interact with the MobileOps REST API from the command line.

## Install the Skill

```bash
npx skills add MobileOps-Team/mobileops-cli
```

## Setup

```bash
# Install
curl -fsSL https://www.mobileops.at/install-cli | bash

# Authenticate (get keys from MobileOps Settings > REST API Keys)
mobileops auth login
```

## Terminology Aliases

Users may use different terms for the same concepts. Always map to the correct CLI command:

| User says | CLI command | Notes |
|-----------|-------------|-------|
| asset, boat, ship, barge, tug, fleet | `vessels` | "Vessel" is the CLI term for any maritime asset |
| user, employee, team member, worker, staff, people | `crew` | All personnel are managed under `crew` |
| task, work order, maintenance job | `jobs` | Scheduled work assignments |
| equipment, machinery, system | `components` | Equipment installed on a vessel |
| spare, supply, inventory | `parts` | Parts belong to components |
| issue, problem, ticket | `work-requests` | Reported issues needing resolution |
| finding, audit item | `deficiencies` | Issues found during inspections |
| NCR, non-conformance | `nonconformities` | Compliance violations |
| observation, audit observation, SIRE finding | `observations` | Observations recorded against an audit (e.g. SIRE Inspection) |
| certificate, cert, license, doc | `vessel-documents` or `personnel-documents` | Depends on whether it belongs to a vessel or person |
| PM, preventive maintenance, routine | `routine-templates` / `routine-calculations` | Templates define schedules; calculations show due dates |
| PO | `purchase-orders` | Purchase orders |
| vendor | `suppliers` | Parts/service providers |
| manufacturer, brand | `makes` | Equipment manufacturers |
| position, role, title | `employee-positions` | Crew job titles |
| department, fleet group | `divisions` | Organizational groupings |

When a user asks about "assets" or "boats", use `mobileops vessels`. When they ask about "employees", use `mobileops crew`. Always translate to the correct CLI resource name.

## Global Flags

| Flag | Description |
|------|-------------|
| `--json` | Output as structured JSON with envelope (ok, data, summary, breadcrumbs) |
| `--help --agent` | Machine-readable help for AI agents |
| `--env <name>` | Environment (default: production) |

## Always Use --json

When using this CLI programmatically, always add `--json` for structured output:

```bash
mobileops vessels list --json
```

Response envelope:

```json
{
  "ok": true,
  "data": [ ... ],
  "summary": "5 vessels found",
  "breadcrumbs": [
    "mobileops vessels get <id>",
    "mobileops components list --vessel-id <id>"
  ]
}
```

## Commands Reference

### Authentication

```bash
mobileops auth login                          # Authenticate (production)
mobileops auth status                         # Check current auth status
mobileops auth logout                         # Remove stored credentials
```

---

### Vessels (Assets)

```bash
# List all vessels
mobileops vessels list --json
mobileops vessels list --division-id <division_id> --json

# Get a specific vessel
mobileops vessels get <id> --json

# Create a vessel
mobileops vessels create --name "MV Atlantic" --division-id <id> --vessel-type-id <id> --json

# Update a vessel
mobileops vessels update <id> --name "MV Atlantic II" --status "active" --json

# List specs for a vessel
mobileops vessels specs <vessel_id> --json

# Filters: --division-id
```

**API: GET /api/assets, GET /api/assets/:id, POST /api/assets, PUT /api/assets/:id**

Writable fields: `--name`, `--division-id`, `--status`, `--category`, `--vessel-type-id`, `--customer-id`, `--customer-name`, `--imo-number`, `--uscg-number`, `--mmsi-number`, `--call-sign`, `--color`, `--length`, `--height`, `--width`, `--dimension-unit`, `--active` (bool), `--voyage-enabled` (bool), `--external` (bool)

---

### Jobs

```bash
# List jobs
mobileops jobs list --json
mobileops jobs list --vessel-id <vessel_id> --json
mobileops jobs list --from 2026-01-01 --to 2026-03-31 --json

# Get a specific job
mobileops jobs get <id> --json

# Create a job
mobileops jobs create --from 2026-04-01 --to 2026-04-15 --notes "Engine overhaul" --vessels <id1>,<id2> --json

# Update a job
mobileops jobs update <id> --notes "Updated notes" --cancel --json

# Filters: --vessel-id, --from, --to, --active-only
```

**API: GET /api/jobs, GET /api/jobs/:id, POST /api/jobs, PUT /api/jobs/:id**

Writable fields: `--from` (date_beginning), `--to` (date_ending), `--notes`, `--user-id`, `--division-id`, `--po-number`, `--ref-number`, `--tbd` (bool), `--cancel` (bool), `--vessels` (comma-separated), `--customers` (comma-separated), `--work-types` (comma-separated), `--locations` (comma-separated)

---

### Crew (Users)

```bash
# List crew members
mobileops crew list --json

# Get a specific crew member
mobileops crew get <id> --json

# Find by employee number
mobileops crew find-by-employee-number <number> --json

# Create a crew member
mobileops crew create --first-name "John" --last-name "Smith" --email "john@example.com" --json

# Update a crew member
mobileops crew update <id> --phone "+1234567890" --time-zone "US/Eastern" --json

# Delete a crew member
mobileops crew delete <id> --json

# Filters: --include-ultra-fields
```

**API: GET /api/users, GET /api/users/:id, GET /api/users/by_employee_number, POST /api/users, PUT /api/users/:id, DELETE /api/users/:id**

Writable fields: `--first-name`, `--last-name`, `--email`, `--password`, `--phone`, `--time-zone`, `--address`, `--birthday`, `--passport-number`, `--mmc-number`, `--employee-number`, `--receive-notifications` (bool), `--login-disabled` (bool), `--archived` (bool), `--division-ids` (comma-separated), `--employee-positions` (comma-separated), `--primary-assets` (comma-separated)

---

### Components

```bash
# List components for a vessel
mobileops components list --vessel-id <vessel_id> --json

# Get a specific component
mobileops components get <id> --json

# Create a component
mobileops components create --name "Main Engine" --vessel-id <id> --has-hours --json

# Update a component
mobileops components update <id> --hours "1500" --critical --json

# Filters: --vessel-id, --component-id, --part-id
```

**API: GET /api/components, GET /api/components/:id, POST /api/components, PUT /api/components/:id**

Writable fields: `--name`, `--vessel-id`, `--model-id`, `--model-number`, `--ultra-field-template-id`, `--hours`, `--lifetime-hours`, `--has-hours` (bool), `--critical` (bool), `--reset-hours` (bool)

---

### Parts

```bash
# List parts for a component
mobileops parts list --vessel-id <vessel_id> --component-id <component_id> --json

# Get a specific part
mobileops parts get <id> --json

# Create a part
mobileops parts create --name "Oil Filter" --vessel-id <id> --component-id <id> --json

# Update a part
mobileops parts update <id> --hours "500" --critical --json

# Filters: --vessel-id, --component-id, --part-id
```

**API: GET /api/parts, GET /api/parts/:id, POST /api/parts, PUT /api/parts/:id**

Writable fields: `--name`, `--vessel-id`, `--component-id`, `--model-id`, `--model-number`, `--code-id`, `--ultra-field-template-id`, `--hours`, `--lifetime-hours`, `--has-hours` (bool), `--critical` (bool), `--reset-hours` (bool)

---

### Work Requests

```bash
# List work requests
mobileops work-requests list --json
mobileops work-requests list --vessel-id <vessel_id> --json

# Get a specific work request
mobileops work-requests get <id> --json

# Create a work request
mobileops work-requests create --vessel-id <id> --description "Pump leaking" --priority "high" --json

# Update a work request
mobileops work-requests update <id> --status "resolved" --resolution-date "2026-03-25" --json

# Filters: --vessel-id, --component-id, --part-id
```

**API: GET /api/work-requests, GET /api/work-requests/:id, POST /api/work-requests, PUT /api/work-requests/:id**

Writable fields: `--user-id`, `--date`, `--resolution-date`, `--vessel-id`, `--component-id`, `--part-id`, `--description`, `--status`, `--plan-of-action`, `--priority`, `--notification-group-id`, `--location`, `--manager-notes`, `--code-id`, `--reference-number`, `--self-resolve` (bool), `--modify-resolvers` (bool), `--tags` (comma-separated), `--resolvers` (comma-separated)

---

### Deficiencies

```bash
# List deficiencies
mobileops deficiencies list --vessel-id <vessel_id> --json

# Get a specific deficiency
mobileops deficiencies get <id> --json

# Create a deficiency
mobileops deficiencies create --vessel-id <id> --description "Hull corrosion" --priority "high" --json

# Update a deficiency
mobileops deficiencies update <id> --status "resolved" --corrective-actions "Repainted hull" --json

# Filters: --vessel-id, --component-id, --part-id
```

**API: GET /api/deficiencies, GET /api/deficiencies/:id, POST /api/deficiencies, PUT /api/deficiencies/:id**

Writable fields: `--user-id`, `--date`, `--resolution-date`, `--vessel-id`, `--component-id`, `--part-id`, `--description`, `--status`, `--priority`, `--notification-group-id`, `--root-cause`, `--corrective-actions`, `--location`, `--manager-notes`, `--code-id`, `--reference-number`, `--self-resolve` (bool), `--modify-resolvers` (bool), `--tags` (comma-separated), `--resolvers` (comma-separated)

---

### Nonconformities

```bash
# List nonconformities
mobileops nonconformities list --vessel-id <vessel_id> --json

# Get a specific nonconformity
mobileops nonconformities get <id> --json

# Create a nonconformity
mobileops nonconformities create --vessel-id <id> --description "Safety protocol breach" --major --json

# Update a nonconformity
mobileops nonconformities update <id> --status "resolved" --corrective-actions "Retrained crew" --json

# Filters: --vessel-id, --component-id, --part-id
```

**API: GET /api/nonconformities, GET /api/nonconformities/:id, POST /api/nonconformities, PUT /api/nonconformities/:id**

Writable fields: `--user-id`, `--date`, `--vessel-id`, `--plan-of-action`, `--description`, `--status`, `--priority`, `--notification-group-id`, `--root-cause`, `--corrective-actions`, `--location`, `--manager-notes`, `--resolution-date`, `--major` (bool), `--external` (bool), `--shoreside` (bool), `--self-resolve` (bool), `--modify-resolvers` (bool), `--tags` (comma-separated), `--resolvers` (comma-separated)

---

### Observations

```bash
# List observations
mobileops observations list --vessel-id <vessel_id> --json

# Pull all observations connected to SIRE Inspection audits (the "SIRE report")
mobileops observations list --audit-type "SIRE Inspection" --json

# All observations linked to a specific audit
mobileops observations list --audit-id <audit_id> --json

# Get a specific observation
mobileops observations get <id> --json

# Create an observation (optionally tied to an audit)
mobileops observations create --vessel-id <id> --user-id <id> --description "Hose leak" --audit-id <audit_id> --json

# Update an observation
mobileops observations update <id> --status "Resolved" --corrective-actions "Replaced hose" --json

# Filters: --audit-type, --audit-id, --vessel-id, --component-id, --part-id, --start-date, --end-date, --status
```

**API: GET /api/observations, GET /api/observations/:id, POST /api/observations, PATCH /api/observations/:id**

Each observation in the response includes an embedded `audit` object (id, type, custom_type_name, name, status, date, completion_date, external_auditor_name, vessel_id, vessel_name) when `audit_id` is set, so a SIRE/audit report does not require a second call.

`--audit-type` accepts the literal Audit type — built-in values: `SIRE Inspection`, `External Audit`, `External Survey`, `External Dry Dock Survey`, `Internal Audit`, `Internal Survey`, `Internal Dry Dock Survey`, `Internal Inspection` — or any per-company custom category name. Match is exact and case-sensitive.

Writable fields: `--user-id`, `--date`, `--vessel-id`, `--description`, `--corrective-actions`, `--plan-of-action`, `--status`, `--notification-group-id`, `--resolution-date`, `--manager-notes`, `--reference-number`, `--assigned-to-id`, `--assigned-to`, `--audit-id`, `--routine-id`, `--self-resolve` (bool), `--modify-resolvers` (bool), `--tags` (comma-separated), `--resolvers` (comma-separated)

---

### Maintenance Reports

```bash
# List maintenance reports
mobileops maintenance-reports list --vessel-id <vessel_id> --json

# Get a specific maintenance report
mobileops maintenance-reports get <id> --json

# Create a maintenance report
mobileops maintenance-reports create --vessel-id <id> --component-id <id> --description "Oil change" --date "2026-03-25" --json

# Update a maintenance report
mobileops maintenance-reports update <id> --description "Oil change completed" --component-hours "1500" --json

# Filters: --vessel-id, --component-id, --part-id
```

**API: GET /api/maintenance-reports, GET /api/maintenance-reports/:id, POST /api/maintenance-reports, PUT /api/maintenance-reports/:id**

Writable fields: `--description`, `--date`, `--vessel-id`, `--component-id`, `--component-hours`, `--part-id`, `--user-id`, `--users` (comma-separated)

---

### Vessel Documents

```bash
# List vessel documents
mobileops vessel-documents list --vessel-id <vessel_id> --json

# Get a specific vessel document
mobileops vessel-documents get <id> --json

# Get attachments
mobileops vessel-documents attachments <id> --json

# Create a vessel document
mobileops vessel-documents create --vessel-id <id> --name "Safety Certificate" --issue-date "2026-01-01" --expire-date "2027-01-01" --json

# Update a vessel document
mobileops vessel-documents update <id> --expire-date "2027-06-01" --notes "Extended" --json

# Delete a vessel document
mobileops vessel-documents delete <id> --json

# Filters: --vessel-id
```

**API: GET /api/vessel-documents, GET /api/vessel-documents/:id, POST /api/vessel-documents, PUT /api/vessel-documents/:id, DELETE /api/vessel-documents/:id, GET /api/vessel-documents/:id/attachments**

Writable fields: `--issue-date`, `--expire-date`, `--endorsement-date`, `--name`, `--notification-group-id`, `--reminder-interval`, `--user-id`, `--vessel-id`, `--notes`, `--template-id`, `--expire-notification` (bool), `--captain-edit` (bool), `--confidential` (bool), `--tags` (comma-separated)

---

### Personnel Documents

```bash
# List personnel documents
mobileops personnel-documents list --user-id <user_id> --json

# Get a specific personnel document
mobileops personnel-documents get <id> --json

# Get attachments
mobileops personnel-documents attachments <id> --json

# Create a personnel document
mobileops personnel-documents create --user-id <id> --template-id <id> --issue-date "2026-01-01" --expire-date "2027-01-01" --json

# Update a personnel document
mobileops personnel-documents update <id> --expire-date "2027-06-01" --notes "Renewed" --json

# Delete a personnel document
mobileops personnel-documents delete <id> --json

# Filters: --user-id, --employee-number, --template-id, --include-ultra-fields
```

**API: GET /api/personnel-documents, GET /api/personnel-documents/:id, POST /api/personnel-documents, PUT /api/personnel-documents/:id, DELETE /api/personnel-documents/:id, GET /api/personnel-documents/:id/attachments**

Writable fields: `--user-id`, `--template-id`, `--issue-date`, `--expire-date`, `--notes`, `--expire-notification` (bool)

---

### Vessel Specs

```bash
# List vessel specs
mobileops vessel-specs list --vessel-id <vessel_id> --json

# Get a specific vessel spec
mobileops vessel-specs get <id> --json

# Create a vessel spec
mobileops vessel-specs create --vessel-id <id> --template-id <id> --name "Engine Specs" --json

# Update a vessel spec
mobileops vessel-specs update <id> --name "Updated Engine Specs" --json

# Filters: --vessel-id
```

**API: GET /api/vessel-specs, GET /api/vessel-specs/:id, POST /api/vessel-specs, PUT /api/vessel-specs/:id**

Writable fields: `--name`, `--vessel-id`, `--template-id`

---

### Vessel Spec Templates

```bash
# List vessel spec templates
mobileops vessel-spec-templates list --json

# Get a specific template
mobileops vessel-spec-templates get <id> --json
```

**API: GET /api/vessel-spec-templates, GET /api/vessel-spec-templates/:id (read-only)**

---

### Customers

```bash
# List customers
mobileops customers list --json

# Get a specific customer
mobileops customers get <id> --json

# Create a customer
mobileops customers create --name "Acme Shipping" --email-address "info@acme.com" --json

# Update a customer
mobileops customers update <id> --phone "+1234567890" --add-vessels <id1>,<id2> --json
```

**API: GET /api/customers, GET /api/customers/:id, POST /api/customers, PUT /api/customers/:id**

Writable fields: `--name`, `--address`, `--city`, `--postal-code`, `--state-province`, `--country`, `--email-address`, `--phone`, `--hex-color`, `--third-party` (bool), `--add-vessels` (comma-separated)

---

### Suppliers

```bash
# List suppliers
mobileops suppliers list --json

# Get a specific supplier
mobileops suppliers get <id> --json

# Create a supplier
mobileops suppliers create --name "Marine Parts Co" --email-address "sales@marineparts.com" --json

# Update a supplier
mobileops suppliers update <id> --notes "Preferred vendor" --phone "+1234567890" --json
```

**API: GET /api/suppliers, GET /api/suppliers/:id, POST /api/suppliers, PUT /api/suppliers/:id**

Writable fields: `--name`, `--notes`, `--address`, `--city`, `--postal-code`, `--country`, `--state-province`, `--email-address`, `--phone`

---

### Makes

```bash
# List makes (manufacturers)
mobileops makes list --json

# Get a specific make
mobileops makes get <id> --json

# Create a make
mobileops makes create --name "Caterpillar" --json

# Update a make
mobileops makes update <id> --phone "+1234567890" --json
```

**API: GET /api/makes, GET /api/makes/:id, POST /api/makes, PUT /api/makes/:id**

Writable fields: `--name`, `--address`, `--city`, `--postal-code`, `--country`, `--state-province`, `--email-address`, `--phone`

---

### Models

```bash
# List models for a make
mobileops models list --make-id <make_id> --json

# Get a specific model
mobileops models get <id> --make-id <make_id> --json

# Create a model
mobileops models create --make-id <id> --name "C32 ACERT" --part-number "CAT-C32" --json

# Update a model
mobileops models update <id> --unit-cost "15000" --notes "Updated pricing" --json
```

**API: GET /api/models, GET /api/models/:id, POST /api/models, PUT /api/models/:id**

Writable fields: `--make-id`, `--serial-number`, `--name`, `--part-number`, `--notes`, `--vessel-id`, `--unit-cost`, `--expected-lifetime-hours`, `--critical` (bool)

---

### Divisions

```bash
# List divisions
mobileops divisions list --json

# Get a specific division
mobileops divisions get <id> --json

# Create a division
mobileops divisions create --name "West Coast" --hex-color "#FF5733" --json

# Update a division
mobileops divisions update <id> --description "West Coast operations" --active --json
```

**API: GET /api/divisions, GET /api/divisions/:id, POST /api/divisions, PUT /api/divisions/:id**

Writable fields: `--name`, `--description`, `--reference-number`, `--hex-color`, `--active` (bool)

---

### Employee Positions

```bash
# List employee positions
mobileops employee-positions list --json

# Get a specific position
mobileops employee-positions get <id> --json

# Create an employee position
mobileops employee-positions create --name "Chief Engineer" --json

# Update an employee position
mobileops employee-positions update <id> --description "Senior engineering role" --active --json
```

**API: GET /api/employee-positions, GET /api/employee-positions/:id, POST /api/employee-positions, PUT /api/employee-positions/:id**

Writable fields: `--name`, `--description`, `--reference-number`, `--hex-color`, `--active` (bool)

---

### Vessel Types

```bash
# List vessel types
mobileops vessel-types list --json

# Get a specific vessel type
mobileops vessel-types get <id> --json

# Create a vessel type
mobileops vessel-types create --name "Tugboat" --json

# Update a vessel type
mobileops vessel-types update <id> --name "Ocean Tugboat" --json
```

**API: GET /api/vessel-types, GET /api/vessel-types/:id, POST /api/vessel-types, PUT /api/vessel-types/:id**

Writable fields: `--name`

---

### Work Types

```bash
# List work types
mobileops work-types list --json

# Get a specific work type
mobileops work-types get <id> --json

# Create a work type
mobileops work-types create --name "Dry Dock" --json

# Update a work type
mobileops work-types update <id> --name "Dry Dock Maintenance" --json
```

**API: GET /api/work-types, GET /api/work-types/:id, POST /api/work-types, PUT /api/work-types/:id**

Writable fields: `--name`

---

### Locations

```bash
# List locations
mobileops locations list --json

# Get a specific location
mobileops locations get <id> --json

# Create a location
mobileops locations create --name "Port of Houston" --type "port" --latitude "29.7604" --longitude "-95.3698" --json

# Update a location
mobileops locations update <id> --name "Port of Houston (Main)" --json
```

**API: GET /api/locations, GET /api/locations/:id, POST /api/locations, PUT /api/locations/:id**

Writable fields: `--name`, `--type`, `--latitude`, `--longitude`

---

### Invoice Statements

```bash
# List invoice statements
mobileops invoice-statements list --json
mobileops invoice-statements list --status <status> --from 2026-01-01 --to 2026-03-31 --json
mobileops invoice-statements list --include-jobs --json

# Filters: --status, --integration-status, --from, --to, --include-jobs
```

**API: GET /api/invoice-statements (read-only)**

---

### Form Instances

```bash
# List form instances
mobileops form-instances list --json
mobileops form-instances list --template-id <form_template_id> --json
mobileops form-instances list --job-id <job_id> --json

# Filters: --template-id, --job-id
```

**API: GET /api/form-instances (read-only)**

---

### Forms (Templates)

```bash
# List form templates
mobileops forms list --json

# Filters: --name
```

**API: GET /api/forms (read-only)**

---

### Routine Templates

```bash
# List routine templates
mobileops routine-templates list --json
mobileops routine-templates list --vessel-id <vessel_id> --json

# Filters: --vessel-id, --component-id, --part-id, --category, --type, --frequency-type
```

**API: GET /api/routine-templates (read-only)**

---

### Routine Calculations

```bash
# List routine calculations
mobileops routine-calculations list --json
mobileops routine-calculations list --vessel-id <vessel_id> --json
mobileops routine-calculations list --due-date-lte 2026-04-01 --json

# Filters: --vessel-id, --component-id, --part-id, --division-id, --frequency-type, --value-lte, --value-gte, --due-date-lte, --due-date-gte, --routine-template-id, --master-template-id
```

**API: GET /api/routine-calculations (read-only)**

---

### Position Reports

```bash
# List position reports
mobileops position-reports list --vessel-id <vessel_id> --json

# Create a position report
mobileops position-reports create --asset-id <vessel_id> --latitude "29.7604" --longitude "-95.3698" --heading "180" --json

# Filters: --vessel-id, --include-ultra-fields
```

**API: GET /api/position-reports, POST /api/position-reports**

Writable fields: `--asset-id` (required), `--latitude` (required, -90 to 90), `--longitude` (required, -180 to 180), `--heading` (optional)

---

### Component Hour Logs

```bash
# List component hour logs
mobileops component-logs list --vessel-id <vessel_id> --json
mobileops component-logs list --vessel-id <vessel_id> --sum --json

# Filters: --vessel-id, --created-after, --created-before, --updated-after, --updated-before, --sum
```

**API: GET /api/component-logs (read-only)**

When `--sum` is used, returns aggregated hours and fuel consumption by component.

---

### Part Requests

```bash
# List part requests
mobileops part-requests list --vessel-id <vessel_id> --json

# Get a specific part request
mobileops part-requests get <id> --json

# Filters: --vessel-id, --component-id, --part-id
```

**API: GET /api/part-requests, GET /api/part-requests/:id (read-only in practice)**

---

### Terms

```bash
# List terms for a supplier
mobileops terms list --supplier-id <supplier_id> --json

# Get a specific term
mobileops terms get <id> --supplier-id <supplier_id> --json

# Create a term
mobileops terms create --supplier-id <id> --model-id <id> --price "1500.00" --date-beginning "2026-01-01" --date-ending "2026-12-31" --json

# Update a term
mobileops terms update <id> --price "1600.00" --notes "Price increase" --json
```

**API: GET /api/terms, GET /api/terms/:id, POST /api/terms, PUT /api/terms/:id**

Writable fields: `--make-id`, `--model-id`, `--model-title`, `--supplier-id`, `--price`, `--notes`, `--date-beginning`, `--date-ending`

---

### Purchase Orders

```bash
# List purchase orders
mobileops purchase-orders list --json
mobileops purchase-orders list --status <status> --from 2026-01-01 --to 2026-03-31 --json

# Filters: --status, --from, --to
```

**API: GET /api/purchase-orders (read-only)**

---

### Wheelhouse Logs

```bash
# List wheelhouse logs
mobileops wheelhouse-logs list --vessel-id <vessel_id> --json
mobileops wheelhouse-logs list --vessel-id <vessel_id> --from 2026-01-01 --to 2026-03-31 --json

# Filters: --vessel-id, --from, --to
```

**API: GET /api/wheelhouse-logs (read-only)**

---

### Cargo Types

```bash
# List cargo types
mobileops cargo-types list --json

# Get a specific cargo type
mobileops cargo-types get <id> --json

# Create a cargo type
mobileops cargo-types create --name "Crude Oil" --color "#8B4513" --measurement-ids <id1>,<id2> --json

# Update a cargo type
mobileops cargo-types update <id> --reference-number "CT-001" --json
```

**API: GET /api/cargo-types, GET /api/cargo-types/:id, POST /api/cargo-types, PUT /api/cargo-types/:id**

Writable fields: `--name`, `--reference-number`, `--color`, `--measurement-ids` (comma-separated)

---

## Write Operations Summary

| Resource | Create | Update | Delete |
|----------|--------|--------|--------|
| vessels | ✅ | ✅ | — |
| jobs | ✅ | ✅ | — |
| crew | ✅ | ✅ | ✅ |
| components | ✅ | ✅ | — |
| parts | ✅ | ✅ | — |
| work-requests | ✅ | ✅ | — |
| deficiencies | ✅ | ✅ | — |
| nonconformities | ✅ | ✅ | — |
| observations | ✅ | ✅ | — |
| maintenance-reports | ✅ | ✅ | — |
| vessel-documents | ✅ | ✅ | ✅ |
| personnel-documents | ✅ | ✅ | ✅ |
| vessel-specs | ✅ | ✅ | — |
| customers | ✅ | ✅ | — |
| suppliers | ✅ | ✅ | — |
| makes | ✅ | ✅ | — |
| models | ✅ | ✅ | — |
| divisions | ✅ | ✅ | — |
| employee-positions | ✅ | ✅ | — |
| vessel-types | ✅ | ✅ | — |
| work-types | ✅ | ✅ | — |
| locations | ✅ | ✅ | — |
| cargo-types | ✅ | ✅ | — |
| terms | ✅ | ✅ | — |
| position-reports | ✅ | — | — |

Read-only resources: vessel-spec-templates, forms, form-instances, routine-templates, routine-calculations, component-logs, part-requests, purchase-orders, wheelhouse-logs, invoice-statements

## Common Workflows

### Find overdue maintenance for a vessel

```bash
# 1. Find the vessel
mobileops vessels list --json | jq '.data[] | select(.name | test("Atlantic"; "i"))'

# 2. Get routine calculations with overdue items
mobileops routine-calculations list --vessel-id <vessel_id> --due-date-lte $(date +%Y-%m-%d) --json

# 3. Check open work requests
mobileops work-requests list --vessel-id <vessel_id> --json
```

### Check crew credentials expiring soon

```bash
# 1. List crew
mobileops crew list --json

# 2. Check personnel documents for a crew member
mobileops personnel-documents list --user-id <user_id> --json
```

### Review vessel compliance

```bash
# 1. Get vessel documents
mobileops vessel-documents list --vessel-id <vessel_id> --json

# 2. Check nonconformities
mobileops nonconformities list --vessel-id <vessel_id> --json

# 3. Check deficiencies
mobileops deficiencies list --vessel-id <vessel_id> --json
```

### Get vessel equipment hierarchy

```bash
# 1. List components for a vessel
mobileops components list --vessel-id <vessel_id> --json

# 2. List parts for a component
mobileops parts list --vessel-id <vessel_id> --component-id <component_id> --json

# 3. Get component hour logs
mobileops component-logs list --vessel-id <vessel_id> --sum --json
```

### Job management

```bash
# 1. List active jobs
mobileops jobs list --active-only --json

# 2. Get job details with forms
mobileops jobs get <job_id> --json
mobileops form-instances list --job-id <job_id> --json
```

### Create a work request from a deficiency

```bash
# 1. Find the deficiency
mobileops deficiencies list --vessel-id <vessel_id> --json

# 2. Create a work request to address it
mobileops work-requests create --vessel-id <vessel_id> --component-id <component_id> \
  --description "Fix hull corrosion from deficiency" --priority "high" --json
```

### Onboard a new crew member

```bash
# 1. Create the crew member
mobileops crew create --first-name "Jane" --last-name "Doe" --email "jane@example.com" \
  --employee-number "EMP-100" --division-ids <div_id> --primary-assets <vessel_id> --json

# 2. Create their personnel documents
mobileops personnel-documents create --user-id <new_user_id> --template-id <template_id> \
  --issue-date "2026-03-25" --expire-date "2027-03-25" --json
```

### Set up a new vessel

```bash
# 1. Create the vessel
mobileops vessels create --name "MV Pacific" --division-id <id> --vessel-type-id <id> --json

# 2. Add components
mobileops components create --name "Main Engine" --vessel-id <new_vessel_id> --has-hours --json

# 3. Add parts to component
mobileops parts create --name "Oil Filter" --vessel-id <new_vessel_id> --component-id <new_component_id> --json

# 4. Create vessel documents
mobileops vessel-documents create --vessel-id <new_vessel_id> --name "Safety Certificate" \
  --issue-date "2026-01-01" --expire-date "2027-01-01" --expire-notification --json
```

## Utility Commands

```bash
mobileops version                             # Check CLI version
mobileops tree                                # Print full command tree
```

## Pagination

All list commands return paginated results. Use `--page` and `--limit`:

```bash
mobileops vessels list --page 1 --limit 50 --json
mobileops vessels list --page 2 --limit 50 --json
```

Default: page 1, limit 10. Maximum limit: 100.

## Error Handling

Errors return structured JSON with `--json`:

```json
{
  "ok": false,
  "error": "auth_error",
  "error_message": "Invalid or Disabled API Key"
}
```

Common errors:
- `auth_error` - Invalid API keys. Run `mobileops auth login`.
- `permission_error` - API key lacks write/delete permission.
- `not_found` - Resource ID does not exist.
- `api_error` - Server error. Retry or check MobileOps status.
