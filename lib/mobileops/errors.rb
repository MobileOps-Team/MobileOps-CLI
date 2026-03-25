# frozen_string_literal: true

module Mobileops
  class Error < StandardError; end

  class AuthError < Error
    def initialize(msg = "Authentication failed. Run `mobileops auth login` to configure credentials.")
      super
    end
  end

  class NotFoundError < Error
    def initialize(msg = "Resource not found.")
      super
    end
  end

  class ApiError < Error
    attr_reader :status, :code

    def initialize(msg = "API request failed.", status: nil, code: nil)
      @status = status
      @code = code
      super(msg)
    end
  end

  class ConfigError < Error
    def initialize(msg = "No credentials found. Run `mobileops auth login` first.")
      super
    end
  end

  class PermissionError < Error
    def initialize(msg = "API key does not have permission for this operation.")
      super
    end
  end
end
