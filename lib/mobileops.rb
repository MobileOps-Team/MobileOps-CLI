# frozen_string_literal: true

require "thor"
require "json"

require_relative "mobileops/version"
require_relative "mobileops/errors"
require_relative "mobileops/config"
require_relative "mobileops/client"
require_relative "mobileops/envelope"
require_relative "mobileops/formatter"
require_relative "mobileops/pagination"
require_relative "mobileops/agent_help"
require_relative "mobileops/commands/auth"
require_relative "mobileops/commands/vessels"
require_relative "mobileops/commands/jobs"
require_relative "mobileops/commands/crew"
require_relative "mobileops/commands/components"
require_relative "mobileops/commands/work_requests"
require_relative "mobileops/cli"

module Mobileops
end
