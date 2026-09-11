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
| vendor | `suppliers` | Parts/service providers (read-only) |
| audit, inspection, SIRE, survey | `audits` | Audits embed their observations, deficiencies and nonconformities |
| hours of rest, work/rest, fatigue log | `work-rests` | One record per user per day |
| event, timeline entry, wheelhouse event, job log | `events search` | Read-only search; events are created in the app |
| code, accounting code, billing code, job code | `codes` | Company-wide or scoped to one vessel |
| fuel reading, tank level, sounding | `fuel-level-readings` | Tank levels over time |
| fuel lift, bunker, fuel consumption | `bunker-partitions` | One lift plus the consumption drawn from it |
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

# Get a specific vessel (by numeric ID, or by its code)
mobileops vessels get <id> --json
mobileops vessels get <code> --by-code --json

# Create a vessel
mobileops vessels create --name "MV Atlantic" --division-id <id> --vessel-type-id <id> --json

# Update a vessel
mobileops vessels update <id> --name "MV Atlantic II" --status "active" --json

# List specs for a vessel
mobileops vessels specs <vessel_id> --json

# Filters: --division-id
```

**API: GET /api/assets, GET /api/assets/:id, POST /api/assets, PUT /api/assets/:id**

Writable fields: `--name`, `--division-id`, `--status` (In Service, Out of Service, In Repair, Tied Up, or a company status), `--category`, `--vessel-type-id`, `--vessel-type-subtype-id`, `--customer-id`, `--customer-name`, `--imo-number`, `--uscg-number`, `--mmsi-number`, `--call-sign`, `--color`, `--length`, `--height`, `--width`, `--dimension-unit`, `--activation-date`, `--specifications`, `--active` (bool), `--voyage-enabled` (bool), `--external` (bool), `--notify-sync` (bool), `--tow-diagram-enabled` (bool)

---

### Jobs

```bash
# List jobs
mobileops jobs list --json
mobileops jobs list --vessel-id <vessel_id> --json
mobileops jobs list --from 2026-01-01 --to 2026-03-31 --json
mobileops jobs list --status LAUNCHED,IN_PROGRESS --json

# Get a specific job
mobileops jobs get <id> --json

# Create a job
mobileops jobs create --from 2026-04-01 --to 2026-04-15 --notes "Engine overhaul" --vessels <id1>,<id2> --json

# Update a job
mobileops jobs update <id> --notes "Updated notes" --cancel --json

# Is a vessel free? Returns the jobs that overlap the window (empty = free)
mobileops jobs vessel-availability --vessel-id <vessel_id> --from 2026-04-01 --to 2026-04-15 --json

# Job report rows from the reporting engine (detailed by default)
mobileops jobs export --from 2026-04-01 --to 2026-04-30 --json
mobileops jobs export --from 2026-04-01 --to 2026-04-30 --simple --json
mobileops jobs export --from 2026-04-01 --to 2026-04-30 --actuals --json   # planned vs actual segments
mobileops jobs export --from 2026-04-01 --to 2026-04-30 --dispatch-segment-template-id <id> --json

# Filters: --vessel-id, --from, --to, --active-only, --status
```

**API: GET /api/jobs, GET /api/jobs/:id, POST /api/jobs, PUT /api/jobs/:id, POST /api/jobs/vessel-availability, POST /api/jobs/job-export**

`--status` takes one or more of `NOT_READY`, `READY`, `LAUNCHED`, `IN_PROGRESS`, `COMPLETE`, `JOB_CANCELED`, `JOB_RESCHEDULED` (comma-separated). `--active-only` is the shorthand for everything except NOT_READY, JOB_CANCELED and JOB_RESCHEDULED. Deleting a launched job is not possible; set `--cancel` instead.

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

Writable fields: `--first-name`, `--last-name`, `--email`, `--password`, `--password-confirmation` (required on create, 8-128 chars), `--phone`, `--time-zone`, `--address`, `--birthday`, `--passport-number`, `--mmc-number`, `--employee-number`, `--employee-code`, `--archived-date`, `--receive-notifications` (bool), `--login-disabled` (bool), `--archived` (bool), `--exclude-tr` (bool), `--division-ids` (comma-separated), `--employee-positions` (comma-separated), `--primary-assets` (comma-separated), `--roles` (comma-separated; replaces roles wholesale, admin roles are stripped)

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

# Filters: --vessel-id
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

# Filters: --vessel-id, --component-id
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

# Filters: --vessel-id
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

# Filters: --audit-type, --audit-id, --vessel-id, --start-date, --end-date, --status
```

**API: GET /api/observations, GET /api/observations/:id, POST /api/observations, PUT /api/observations/:id**

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

# Create a vessel spec (one per vessel; spec-fields keys must exist on the template)
mobileops vessel-specs create --vessel-id <id> --template-id <id> --spec-fields '{"length":"120 ft","engine":"CAT C32"}' --json

# Update a vessel spec
mobileops vessel-specs update <id> --template-id <id> --spec-fields '{"engine":"CAT C32 ACERT"}' --json

# Filters: --vessel-id
```

**API: GET /api/vessel-specs, GET /api/vessel-specs/:id, POST /api/vessel-specs, PUT /api/vessel-specs/:id**

Writable fields: `--vessel-id`, `--template-id` (required), `--spec-fields` (JSON object keyed by template field). `spec-fields` is filtered against the stored template, so change the template in one call and send its values in a second.

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

# Create a customer together with its own vessels
mobileops customers create --name "Acme Shipping" --add-vessels "Barge 1,Barge 2" --json

# Update a customer; add vessels by name, remove them by the key from the customer's vessels map
mobileops customers update <id> --phone "+1234567890" --add-vessels "Barge 3" --json
mobileops customers update <id> --remove-vessels <vessel_key> --json
```

**API: GET /api/customers, GET /api/customers/:id, POST /api/customers, PUT /api/customers/:id**

Writable fields: `--name`, `--address`, `--city`, `--postal-code`, `--state-province`, `--country`, `--email-address`, `--phone`, `--hex-color`, `--third-party` (bool), `--add-vessels` (comma-separated vessel *names*), `--remove-vessels` (update only; comma-separated vessel *keys*). Customer vessels are the customer-owned boats a job can be booked against; each is stored under a generated key that appears in the customer's `vessels` map.

---

### Suppliers

```bash
# List suppliers
mobileops suppliers list --json

# Get a specific supplier
mobileops suppliers get <id> --json
```

**API: GET /api/suppliers, GET /api/suppliers/:id (read-only; suppliers are managed in the app)**

---

### Makes

```bash
# List makes (manufacturers)
mobileops makes list --json

# Get a specific make
mobileops makes get <id> --json

# Create a make (create the make first, then the models filed under it)
mobileops makes create --name "Caterpillar" --json
```

**API: GET /api/makes, GET /api/makes/:id, POST /api/makes (no update through the API)**

Writable fields: `--name`, `--address`, `--city`, `--postal-code`, `--country`, `--state-province`, `--email-address`, `--phone`

---

### Models

```bash
# List models for a make
mobileops models list --make-id <make_id> --json

# Get a specific model
mobileops models get <id> --make-id <make_id> --json

# Create a model (starts with quantity 0)
mobileops models create --make-id <id> --name "C32 ACERT" --part-number "CAT-C32" --json
```

**API: GET /api/models, GET /api/models/:id, POST /api/models (no update through the API)**

Writable fields: `--make-id`, `--serial-number`, `--name`, `--part-number`, `--notes`, `--vessel-id`, `--unit-cost`, `--expected-lifetime-hours`, `--critical-spares-threshold`, `--critical` (bool)

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

# Mark a statement as pushed to an external accounting system
mobileops invoice-statements update <id> --integration-status "synced" --json

# Filters: --status, --integration-status, --from, --to, --include-jobs
```

**API: GET /api/invoice-statements, PUT /api/invoice-statements/:id**

`integration_status` is the only writable field: a free-form marker for where the statement stands in an external system. Poll with `--integration-status` to find statements not yet pushed.

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

# Filters: --id (single template, other filters ignored), --vessel-id, --component-id, --part-id, --category, --type, --frequency-type, --master (true/false), --master-id, --universal (true/false)
```

**API: GET /api/routine-templates (read-only)**

---

### Routine Calculations

```bash
# List routine calculations
mobileops routine-calculations list --json
mobileops routine-calculations list --vessel-id <vessel_id> --json
mobileops routine-calculations list --due-date-lte 2026-04-01 --json

# Filters: --vessel-id, --component-id, --part-id, --division-id, --frequency-type, --value-lte, --value-gte, --due-date-lte, --due-date-gte, --routine-template-id, --master-template-id, --component-risk-score (>=, 1-10), --part-risk-score (>=, 1-10)

# `value` is the distance to due in the frequency's unit (days for time-based routines, hours for
# component-hour routines) and is negative when overdue, so `--value-lte 0` = everything due or overdue.
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
mobileops part-requests list --json

# Get a specific part request
mobileops part-requests get <id> --json
```

**API: GET /api/part-requests, GET /api/part-requests/:id (read-only; no list filters)**

---

### Terms

```bash
# List terms for a supplier
mobileops terms list --supplier-id <supplier_id> --json

# Get a specific term
mobileops terms get <id> --supplier-id <supplier_id> --json

# Create a term (price is stored as a whole number; decimals are truncated)
mobileops terms create --supplier-id <id> --model-id <id> --price "1500" --date-beginning "2026-01-01" --date-ending "2026-12-31" --json
```

**API: GET /api/terms, GET /api/terms/:id, POST /api/terms (no update through the API)**

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
```

**API: GET /api/cargo-types, GET /api/cargo-types/:id (read-only; cargo types are managed in the app)**

---

### Audits

```bash
# List audits, newest first (SIRE inspections, internal/external audits, surveys, dry dock surveys)
mobileops audits list --json
mobileops audits list --limit 50 --json

# Then drill into the observations of one audit
mobileops observations list --audit-id <audit_id> --json
```

**API: GET /api/audits (read-only, no filters)**

Each audit embeds its `observations`, `deficiencies` and `nonconformities`, so larger pages cost more. To find every SIRE finding across audits use `observations list --audit-type "SIRE Inspection"` instead.

---

### Work Rests (Hours of Rest)

```bash
# List work/rest logs, newest first (one record per user per day)
mobileops work-rests list --json
mobileops work-rests list --user-id <user_id> --from 2026-03-01 --to 2026-03-31 --json
mobileops work-rests list --vessel-ids <id1>,<id2> --json

# Filters: --from, --to, --user-id, --vessel-ids, --employee-position-ids, --employee-rate-ids (all ID filters comma-separated)
```

**API: GET /api/work-rests (read-only)**

Each record holds that day's `periods`, `reviews`, `attestations`, `violations` and `violation_count`. The ID filters match a record when any of its periods names one of the IDs.

---

### Events

```bash
# Search events, newest first (all filters optional and comma-separated)
mobileops events search --json
mobileops events search --job-ids <job_id> --json
mobileops events search --job-ref-numbers "J-1001,J-1002" --json
mobileops events search --order-ids <order_id> --json
mobileops events search --vessel-ids <vessel_id> --types "Depart,Arrive" --json
mobileops events search --work-type-ids <id> --json

# Filters: --job-ref-numbers (wins over --job-ids), --job-ids, --order-ref-numbers (wins over --order-ids), --order-ids,
#          --types, --event-type-ids (wins over --work-type-ids), --work-type-ids, --vessel-ids
```

**API: POST /api/events/search (read-only; events are created and edited in the app)**

The JSON response also carries `wheelhouse_template_event_fields`, the field definitions needed to read each event's values. To fetch a single event, search with a narrow filter such as `--job-ids` and pick it from the results.

---

### Codes

```bash
# List codes (archived codes are not returned)
mobileops codes list --json

# Get a specific code
mobileops codes get <id> --json

# Create a company-wide code, or one scoped to a vessel
mobileops codes create --name "FUEL" --type "Accounting" --description "Fuel purchases" --json
mobileops codes create --name "DECK-01" --type "Job" --vessel-id <vessel_id> --json

# Update a code; --company-wide clears the vessel scope
mobileops codes update <id> --description "Bunker fuel purchases" --json
mobileops codes update <id> --company-wide --json
```

**API: GET /api/codes, GET /api/codes/:id, POST /api/codes, PUT /api/codes/:id**

Writable fields: `--name`, `--description`, `--type` (the app offers Accounting, Billing, Job, Other), `--vessel-id`, `--company-wide` (update only, bool). Code IDs are what `--code-id` expects on parts, work requests and deficiencies.

---

### Fuel Level Readings

```bash
# Tank level readings for a vessel, oldest first
mobileops fuel-level-readings list --vessel-id <vessel_id> --json
mobileops fuel-level-readings list --vessel-id <vessel_id> --from 2026-03-01T00:00:00Z --to 2026-03-31T23:59:59Z --json
mobileops fuel-level-readings list --job-id <job_id> --json

# Filters: --vessel-id, --job-id, --from (ISO 8601), --to (ISO 8601)
```

**API: GET /api/fuel-level-readings (read-only)**

Each reading carries `reading_at`, `reading_type`, `level`, the fuel type, the measurement unit and the tank `location_*` it was taken at.

---

### Bunker Partitions (Fuel Lifts and Consumption)

```bash
# Fuel lifts for a vessel with the consumption drawn from each, oldest lift first
mobileops bunker-partitions list --vessel-id <vessel_id> --json
mobileops bunker-partitions list --vessel-id <vessel_id> --from 2026-03-01T00:00:00 --to 2026-03-31T23:59:59 --json

# Partitions that fed a particular job
mobileops bunker-partitions list --job-id <job_id> --json

# Filters: --vessel-id, --job-id, --from, --to (ISO 8601, read in the company's time zone)
```

**API: GET /api/bunker-partitions (read-only)**

A partition is one fuel lift (`quantity`, `cost`, `vendor_*`, `timestamp`) plus its `consumed` total and a `consumption` array attributing burn to jobs (FIFO across lifts). `active` marks lifts that still hold fuel.

---

## Write Operations Summary

| Resource | Create | Update | Delete |
|----------|--------|--------|--------|
| vessels | ✅ | ✅ | — |
| jobs | ✅ | ✅ | — (launched jobs are cancelled with `--cancel`) |
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
| makes | ✅ | — | — |
| models | ✅ | — | — |
| terms | ✅ | — | — |
| divisions | ✅ | ✅ | — |
| employee-positions | ✅ | ✅ | — |
| vessel-types | ✅ | ✅ | — |
| work-types | ✅ | ✅ | — |
| locations | ✅ | ✅ | — |
| codes | ✅ | ✅ | — |
| invoice-statements | — | ✅ (integration-status only) | — |
| position-reports | ✅ | — | — |

Read-only resources: suppliers, cargo-types, part-requests, vessel-spec-templates, forms, form-instances, routine-templates, routine-calculations, component-logs, purchase-orders, wheelhouse-logs, audits, work-rests, events, fuel-level-readings, bunker-partitions

Not available through the API at all (manage them in the MobileOps app): vendors, fuel types, fuel removals, event create/edit/delete.

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

### Pull a SIRE inspection report

```bash
# 1. Find the audits (SIRE inspections carry type "SIRE Inspection")
mobileops audits list --limit 50 --json | jq '.data[] | select(.type == "SIRE Inspection") | {id, name, date, vessel_name, status}'

# 2. All observations across every SIRE inspection, or for one audit
mobileops observations list --audit-type "SIRE Inspection" --limit 100 --json
mobileops observations list --audit-id <audit_id> --json
```

### Reconcile fuel for a vessel

```bash
# 1. Lifts and the consumption attributed to each job (FIFO)
mobileops bunker-partitions list --vessel-id <vessel_id> --from 2026-03-01T00:00:00 --to 2026-03-31T23:59:59 --json

# 2. Tank readings over the same window
mobileops fuel-level-readings list --vessel-id <vessel_id> --from 2026-03-01T00:00:00Z --to 2026-03-31T23:59:59Z --json

# 3. Engine hours and fuel burn summed by component
mobileops component-logs list --vessel-id <vessel_id> --sum --json
```

### Check whether a vessel is free before booking

```bash
mobileops jobs vessel-availability --vessel-id <vessel_id> --from 2026-04-01 --to 2026-04-15 --json
# empty data = free; otherwise the overlapping jobs are listed
mobileops jobs create --from 2026-04-01 --to 2026-04-15 --vessels <vessel_id> --customers <customer_id> --json
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
mobileops update                              # Update now (binary + this skill); the CLI also self-updates daily on its own
                                              # set MOBILEOPS_AUTO_UPDATE=0 to pin a version
mobileops tree                                # Print full command tree
mobileops --help --agent                      # Machine-readable manifest (commands, flags, API operation per command)
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
