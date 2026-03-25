require_relative "lib/mobileops/version"

Gem::Specification.new do |spec|
  spec.name          = "mobileops-cli"
  spec.version       = Mobileops::VERSION
  spec.authors       = ["MobileOps"]
  spec.email         = ["support@mobileops.com"]
  spec.summary       = "Command-line interface for MobileOps"
  spec.description   = "CLI tool for managing vessels, jobs, crew, and operations via the MobileOps API. Designed for humans and AI agents alike."
  spec.homepage      = "https://www.mobileops.at"
  spec.license       = "MIT"
  spec.required_ruby_version = ">= 3.1.0"

  spec.metadata = {
    "homepage_uri"    => spec.homepage,
    "source_code_uri" => "https://github.com/MobileOps/mobileops-cli",
    "bug_tracker_uri" => "https://github.com/MobileOps/mobileops-cli/issues"
  }

  spec.files         = Dir["lib/**/*", "bin/*", "LICENSE.txt"]
  spec.bindir        = "bin"
  spec.executables   = ["mobileops"]
  spec.require_paths = ["lib"]

  spec.add_dependency "thor", "~> 1.3"
  spec.add_dependency "faraday", "~> 2.0"
  spec.add_dependency "faraday-retry", "~> 2.0"
  spec.add_dependency "terminal-table", "~> 3.0"

  spec.add_development_dependency "bundler", "~> 2.0"
  spec.add_development_dependency "rake", "~> 13.0"
  spec.add_development_dependency "rspec", "~> 3.0"
  spec.add_development_dependency "webmock", "~> 3.0"
end
