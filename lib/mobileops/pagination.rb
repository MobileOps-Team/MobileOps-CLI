# frozen_string_literal: true

module Mobileops
  module Pagination
    def self.included(base)
      base.class_option :page, type: :numeric, default: 1, desc: "Page number"
      base.class_option :limit, type: :numeric, default: 10, desc: "Items per page (max 100)"
    end

    private

    def pagination_params
      limit = [[options[:limit].to_i, 1].max, 100].min
      {
        page: [options[:page].to_i, 1].max,
        limit: limit
      }
    end
  end
end
