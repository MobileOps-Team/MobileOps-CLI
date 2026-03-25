#!/bin/bash
# MobileOps CLI installer
# Usage: curl -fsSL https://www.mobileops.at/install-cli | bash

set -e

echo "Installing MobileOps CLI..."

# Check Ruby is available
if ! command -v ruby &> /dev/null; then
  echo "Error: Ruby is required but not installed."
  echo "Install Ruby 3.1+ from https://www.ruby-lang.org/en/downloads/"
  exit 1
fi

# Check Ruby version
RUBY_VERSION=$(ruby -e 'puts RUBY_VERSION')
MAJOR=$(echo "$RUBY_VERSION" | cut -d. -f1)
MINOR=$(echo "$RUBY_VERSION" | cut -d. -f2)

if [ "$MAJOR" -lt 3 ] || ([ "$MAJOR" -eq 3 ] && [ "$MINOR" -lt 1 ]); then
  echo "Error: Ruby 3.1+ is required (found $RUBY_VERSION)"
  exit 1
fi

# Install the gem
gem install mobileops-cli

echo ""
echo "MobileOps CLI installed successfully!"
echo ""
echo "Get started:"
echo "  mobileops auth login"
echo "  mobileops vessels list"
echo ""
echo "For AI agents:"
echo "  mobileops --help --agent"
