# frozen_string_literal: true

require "json"

module Mobileops
  module AgentHelp
    COMMANDS = [
      {
        name: "auth login",
        description: "Authenticate with MobileOps API using access key and secret key",
        options: [
          { name: "--host", type: "string", default: "https://www.mobileops.at", description: "API host URL" }
        ],
        output: { type: "message" }
      },
      {
        name: "auth logout",
        description: "Remove stored credentials",
        options: [],
        output: { type: "message" }
      },
      {
        name: "auth status",
        description: "Check authentication status and API connectivity",
        options: [],
        output: { type: "envelope", data_shape: "auth_status" }
      },
      {
        name: "vessels list",
        description: "List all vessels/assets with pagination",
        options: [
          { name: "--page", type: "integer", default: 1 },
          { name: "--limit", type: "integer", default: 10, max: 100 },
          { name: "--json", type: "boolean", default: false }
        ],
        output: { type: "envelope", data_shape: "array<vessel>" }
      },
      {
        name: "vessels get ID",
        description: "Get a single vessel by ID",
        options: [
          { name: "--json", type: "boolean", default: false }
        ],
        output: { type: "envelope", data_shape: "vessel" }
      },
      {
        name: "vessels specs ID",
        description: "List specs for a vessel",
        options: [
          { name: "--json", type: "boolean", default: false }
        ],
        output: { type: "envelope", data_shape: "array<vessel_spec>" }
      },
      {
        name: "jobs list",
        description: "List jobs with optional filters",
        options: [
          { name: "--page", type: "integer", default: 1 },
          { name: "--limit", type: "integer", default: 10, max: 100 },
          { name: "--vessel-id", type: "string", description: "Filter by vessel ID" },
          { name: "--from", type: "string", description: "Start date (YYYY-MM-DD)" },
          { name: "--to", type: "string", description: "End date (YYYY-MM-DD)" },
          { name: "--active-only", type: "boolean", default: false },
          { name: "--json", type: "boolean", default: false }
        ],
        output: { type: "envelope", data_shape: "array<job>" }
      },
      {
        name: "jobs get ID",
        description: "Get a single job by ID",
        options: [
          { name: "--json", type: "boolean", default: false }
        ],
        output: { type: "envelope", data_shape: "job" }
      },
      {
        name: "crew list",
        description: "List crew members/users",
        options: [
          { name: "--page", type: "integer", default: 1 },
          { name: "--limit", type: "integer", default: 10, max: 100 },
          { name: "--json", type: "boolean", default: false }
        ],
        output: { type: "envelope", data_shape: "array<user>" }
      },
      {
        name: "crew get ID",
        description: "Get a single crew member by ID",
        options: [
          { name: "--json", type: "boolean", default: false }
        ],
        output: { type: "envelope", data_shape: "user" }
      },
      {
        name: "crew find-by-employee-number NUMBER",
        description: "Find a crew member by employee number",
        options: [
          { name: "--json", type: "boolean", default: false }
        ],
        output: { type: "envelope", data_shape: "user" }
      },
      {
        name: "components list",
        description: "List components with optional vessel filter",
        options: [
          { name: "--page", type: "integer", default: 1 },
          { name: "--limit", type: "integer", default: 10, max: 100 },
          { name: "--vessel-id", type: "string", description: "Filter by vessel ID" },
          { name: "--json", type: "boolean", default: false }
        ],
        output: { type: "envelope", data_shape: "array<component>" }
      },
      {
        name: "components get ID",
        description: "Get a single component by ID",
        options: [
          { name: "--json", type: "boolean", default: false }
        ],
        output: { type: "envelope", data_shape: "component" }
      },
      {
        name: "work-requests list",
        description: "List work requests with optional vessel filter",
        options: [
          { name: "--page", type: "integer", default: 1 },
          { name: "--limit", type: "integer", default: 10, max: 100 },
          { name: "--vessel-id", type: "string", description: "Filter by vessel ID" },
          { name: "--json", type: "boolean", default: false }
        ],
        output: { type: "envelope", data_shape: "array<work_request>" }
      },
      {
        name: "work-requests get ID",
        description: "Get a single work request by ID",
        options: [
          { name: "--json", type: "boolean", default: false }
        ],
        output: { type: "envelope", data_shape: "work_request" }
      }
    ].freeze

    module_function

    def generate
      output = {
        name: "mobileops",
        version: VERSION,
        description: "CLI for managing vessels, jobs, crew, and operations via the MobileOps API",
        auth: {
          type: "api_key_pair",
          setup: "mobileops auth login",
          headers: %w[X-Api-Access-Key X-Api-Secret-Key]
        },
        commands: COMMANDS
      }

      puts JSON.pretty_generate(output)
    end
  end
end
