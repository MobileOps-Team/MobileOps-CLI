# frozen_string_literal: true

require "faraday"
require "faraday/retry"
require "json"

module Mobileops
  class Client
    def initialize(host: nil, access_key: nil, secret_key: nil, environment: nil)
      env = environment || ENV["MOBILEOPS_ENV"] || "production"
      config = Config.load(environment: env) rescue nil

      @host       = host       || config&.dig(:host)       || Config.host_for(env)
      @access_key = access_key || config&.dig(:access_key)
      @secret_key = secret_key || config&.dig(:secret_key)

      raise AuthError unless @access_key && @secret_key
    end

    def get(path, params = {})
      response = connection.get(api_path(path), params)
      handle_response(response)
    end

    def post(path, body = {})
      response = connection.post(api_path(path)) do |req|
        req.body = body.to_json
      end
      handle_response(response)
    end

    def patch(path, body = {})
      response = connection.patch(api_path(path)) do |req|
        req.body = body.to_json
      end
      handle_response(response)
    end

    private

    def connection
      @connection ||= Faraday.new(url: @host) do |f|
        f.request :retry, max: 2, interval: 0.5, backoff_factor: 2
        f.headers["X-Api-Access-Key"] = @access_key
        f.headers["X-Api-Secret-Key"] = @secret_key
        f.headers["Content-Type"]     = "application/json"
        f.headers["Accept"]           = "application/json"
        f.headers["User-Agent"]       = "mobileops-cli/#{VERSION}"
        f.response :json, content_type: /\bjson$/
        f.adapter Faraday.default_adapter
      end
    end

    def api_path(path)
      "/api/#{path.sub(%r{^/}, '')}"
    end

    def handle_response(response)
      case response.status
      when 200..299
        body = response.body

        # The MobileOps API returns some errors as HTTP 200 with an error body
        if body.is_a?(Hash) && body["error"]
          error_msg = body["error_message"] || body["error"]
          case body["error"]
          when "api_connection_error"
            raise AuthError, error_msg
          when "api_key_write_error", "api_key_delete_error"
            raise PermissionError, error_msg
          when "resource_not_found"
            raise NotFoundError, error_msg
          else
            raise ApiError.new(error_msg, status: response.status)
          end
        end

        body
      when 401
        raise AuthError
      when 403
        raise PermissionError
      when 404
        error_msg = extract_error_message(response.body) || "Resource not found."
        raise NotFoundError, error_msg
      else
        error_msg = extract_error_message(response.body) || "API request failed (HTTP #{response.status})."
        raise ApiError.new(error_msg, status: response.status)
      end
    end

    def extract_error_message(body)
      return nil unless body.is_a?(Hash)

      body["error_message"] || body["error"] || body["errors"]&.to_s
    end
  end
end
