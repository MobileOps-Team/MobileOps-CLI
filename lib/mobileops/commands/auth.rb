# frozen_string_literal: true

module Mobileops
  module Commands
    class Auth < Thor
      desc "login", "Authenticate with the MobileOps API"
      option :host, type: :string, desc: "API host URL (overrides environment default)"
      option :env, type: :string, desc: "Environment: production, staging, development, test"
      def login
        env = options[:env] || ENV["MOBILEOPS_ENV"] || "production"
        host = options[:host] || Config.host_for(env)

        puts "Environment: #{env} (#{host})"
        puts

        print "Access Key: "
        access_key = $stdin.gets&.chomp
        print "Secret Key: "
        secret_key = $stdin.gets&.chomp

        if access_key.nil? || access_key.empty? || secret_key.nil? || secret_key.empty?
          $stderr.puts "Error: Access key and secret key are required."
          exit 1
        end

        # Validate credentials with a lightweight API call
        print "Verifying credentials... "
        client = Client.new(host: host, access_key: access_key, secret_key: secret_key)

        begin
          client.get("assets", { limit: 1 })
          Config.save(host: host, access_key: access_key, secret_key: secret_key, environment: env)
          puts "OK"
          puts
          puts "Authenticated successfully!"
          puts "Credentials saved to #{Config.path(environment: env)}"
        rescue AuthError
          puts "FAILED"
          $stderr.puts "Error: Invalid credentials. Please check your access key and secret key."
          exit 1
        rescue StandardError => e
          puts "FAILED"
          $stderr.puts "Error: Could not connect to #{host} (#{e.message})"
          exit 1
        end
      end

      desc "logout", "Remove stored credentials"
      option :env, type: :string, desc: "Environment to logout from"
      def logout
        env = options[:env] || Config.current_environment

        if Config.exists?(environment: env)
          Config.delete(environment: env)
          puts "Credentials removed for #{env}."
        else
          puts "No credentials found for #{env}."
        end
      end

      desc "status", "Check authentication status"
      option :json, type: :boolean, default: false, desc: "Output as JSON"
      option :env, type: :string, desc: "Environment to check"
      def status
        env = options[:env] || Config.current_environment

        unless Config.exists?(environment: env)
          if options[:json]
            puts JSON.pretty_generate(Envelope.error(
              message: "Not authenticated for #{env}",
              breadcrumbs: ["mobileops auth login --env #{env}"]
            ))
          else
            puts "Not authenticated for #{env}. Run `mobileops auth login --env #{env}` to get started."
          end
          return
        end

        config = Config.load(environment: env)
        key_preview = "#{config[:access_key][0..7]}..."

        # Test connectivity
        connected = begin
          client = Client.new(environment: env)
          client.get("assets", { limit: 1 })
          true
        rescue StandardError
          false
        end

        data = {
          authenticated: true,
          environment: env,
          host: config[:host],
          access_key: key_preview,
          api_connected: connected
        }

        envelope = Envelope.wrap(
          data: data,
          summary: connected ? "Authenticated and connected (#{env})" : "Authenticated but API unreachable (#{env})",
          breadcrumbs: ["mobileops vessels list", "mobileops jobs list"]
        )

        if options[:json]
          puts JSON.pretty_generate(envelope)
        else
          puts "Environment: #{env}"
          puts "Status:      #{connected ? '✓ Connected' : '⚠ API unreachable'}"
          puts "Host:        #{config[:host]}"
          puts "Access Key:  #{key_preview}"
          puts "Config:      #{Config.path(environment: env)}"
        end
      end
    end
  end
end
