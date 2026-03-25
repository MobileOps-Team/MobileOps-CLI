# frozen_string_literal: true

module Mobileops
  module Commands
    class WorkRequests < Thor
      include Pagination

      class_option :json, type: :boolean, default: false, desc: "Output as JSON"

      desc "list", "List work requests"
      option :vessel_id, type: :string, desc: "Filter by vessel ID"
      def list
        client = Client.new
        params = pagination_params
        params[:vessel_id] = options[:vessel_id] if options[:vessel_id]

        response = client.get("work-requests", params)

        envelope = Envelope.wrap_collection(
          response,
          resource_name: "work requests",
          breadcrumbs: response["data"]&.first(3)&.map { |wr|
            id = wr["_id"] || wr["id"]
            "mobileops work-requests get #{id}"
          } || []
        )

        Formatter.output(envelope, json: options[:json])
      end

      desc "get ID", "Get a single work request by ID"
      def get(id)
        client = Client.new
        response = client.get("work-requests/#{id}")

        vessel_id = response["vessel_id"]
        envelope = Envelope.wrap_record(
          response,
          resource_name: "Work request",
          breadcrumbs: [
            ("mobileops vessels get #{vessel_id}" if vessel_id),
            "mobileops work-requests list"
          ].compact
        )

        Formatter.output(envelope, json: options[:json])
      end
    end
  end
end
