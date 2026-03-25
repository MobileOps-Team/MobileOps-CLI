# frozen_string_literal: true

module Mobileops
  class CLI < Thor
    class_option :json, type: :boolean, default: false, desc: "Output as JSON"
    class_option :agent, type: :boolean, default: false, desc: "Machine-readable help (use with --help)"
    class_option :env, type: :string, desc: "Environment: production, staging, development, test"

    desc "auth SUBCOMMAND", "Manage authentication"
    subcommand "auth", Commands::Auth

    desc "vessels SUBCOMMAND", "Manage vessels/assets"
    subcommand "vessels", Commands::Vessels

    desc "jobs SUBCOMMAND", "Manage jobs"
    subcommand "jobs", Commands::Jobs

    desc "crew SUBCOMMAND", "Manage crew members"
    subcommand "crew", Commands::Crew

    desc "components SUBCOMMAND", "Manage vessel components"
    subcommand "components", Commands::Components

    desc "work-requests SUBCOMMAND", "Manage work requests"
    subcommand "work-requests", Commands::WorkRequests

    desc "version", "Print CLI version"
    def version
      if options[:json]
        puts JSON.pretty_generate({ version: VERSION })
      else
        puts "mobileops-cli #{VERSION}"
      end
    end

    def self.help(shell, subcommand = false)
      if ARGV.include?("--agent")
        AgentHelp.generate
      else
        super
      end
    end

    # Friendly error handling
    def self.exit_on_failure?
      true
    end
  end
end
