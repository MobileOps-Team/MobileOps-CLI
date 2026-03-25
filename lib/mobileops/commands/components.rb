# frozen_string_literal: true

module Mobileops
  module Commands
    class Components < Thor
      include Pagination

      class_option :json, type: :boolean, default: false, desc: "Output as JSON"

      desc "list", "List components"
      option :vessel_id, type: :string, desc: "Filter by vessel ID"
      def list
        client = Client.new
        params = pagination_params
        params[:vessel_id] = options[:vessel_id] if options[:vessel_id]

        response = client.get("components", params)

        envelope = Envelope.wrap_collection(
          response,
          resource_name: "components",
          breadcrumbs: response["data"]&.first(3)&.map { |c|
            id = c["_id"] || c["id"]
            "mobileops components get #{id}"
          } || []
        )

        Formatter.output(envelope, json: options[:json])
      end

      desc "get ID", "Get a single component by ID"
      def get(id)
        client = Client.new
        response = client.get("components/#{id}")

        vessel_id = response["vessel_id"]
        envelope = Envelope.wrap_record(
          response,
          resource_name: "Component",
          breadcrumbs: [
            ("mobileops vessels get #{vessel_id}" if vessel_id),
            "mobileops components list"
          ].compact
        )

        Formatter.output(envelope, json: options[:json])
      end
    end
  end
end
