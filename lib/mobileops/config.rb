# frozen_string_literal: true

require "json"
require "fileutils"

module Mobileops
  class Config
    CONFIG_DIR = File.join(Dir.home, ".config", "mobileops")

    ENVIRONMENTS = {
      "production"  => "https://www.mobileops.at",
      "staging"     => "https://staging.mobileops.at",
      "development" => "http://localhost:3000",
      "test"        => "http://localhost:3001"
    }.freeze

    class << self
      # Resolve environment: CLI flag > ENV var > default (production)
      def current_environment
        ENV["MOBILEOPS_ENV"] || "production"
      end

      def host_for(environment)
        ENVIRONMENTS[environment] || ENVIRONMENTS["production"]
      end

      def credentials_file(environment = nil)
        env = environment || current_environment
        if env == "production"
          File.join(CONFIG_DIR, "credentials.json")
        else
          File.join(CONFIG_DIR, "credentials.#{env}.json")
        end
      end

      def load(environment: nil)
        env = environment || current_environment
        file = credentials_file(env)

        raise ConfigError, "No credentials found for #{env}. Run `mobileops auth login --env #{env}` first." unless File.exist?(file)

        data = JSON.parse(File.read(file))
        {
          host: data["host"],
          access_key: data["access_key"],
          secret_key: data["secret_key"],
          environment: env
        }
      rescue JSON::ParserError
        raise ConfigError, "Credentials file is corrupted. Run `mobileops auth login --env #{env}` to reconfigure."
      end

      def save(host:, access_key:, secret_key:, environment: nil)
        env = environment || current_environment
        FileUtils.mkdir_p(CONFIG_DIR, mode: 0o700)

        data = {
          host: host,
          access_key: access_key,
          secret_key: secret_key,
          environment: env
        }

        file = credentials_file(env)
        File.write(file, JSON.pretty_generate(data))
        File.chmod(0o600, file)
      end

      def delete(environment: nil)
        env = environment || current_environment
        file = credentials_file(env)
        File.delete(file) if File.exist?(file)
      end

      def exists?(environment: nil)
        env = environment || current_environment
        File.exist?(credentials_file(env))
      end

      def path(environment: nil)
        env = environment || current_environment
        credentials_file(env)
      end
    end
  end
end
