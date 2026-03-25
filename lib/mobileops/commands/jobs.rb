# frozen_string_literal: true

module Mobileops
  module Commands
    class Jobs < Thor
      include Pagination

      class_option :json, type: :boolean, default: false, desc: "Output as JSON"

      desc "list", "List jobs"
      option :vessel_id, type: :string, desc: "Filter by vessel ID"
      option :from, type: :string, desc: "Start date (YYYY-MM-DD)"
      option :to, type: :string, desc: "End date (YYYY-MM-DD)"
      option :active_only, type: :boolean, default: false, desc: "Only show active jobs"
      def list
        client = Client.new
        params = pagination_params

        params[:vessel_id] = options[:vessel_id] if options[:vessel_id]
        params[:date_beginning] = options[:from] if options[:from]
        params[:date_ending] = options[:to] if options[:to]
        params[:scope_by_active_event_list_status] = true if options[:active_only]

        response = client.get("jobs", params)

        envelope = Envelope.wrap_collection(
          response,
          resource_name: "jobs",
          breadcrumbs: response["data"]&.first(3)&.map { |j|
            id = j["_id"] || j["id"]
            "mobileops jobs get #{id}"
          } || []
        )

        Formatter.output(envelope, json: options[:json])
      end

      desc "get ID", "Get a single job by ID"
      def get(id)
        client = Client.new
        response = client.get("jobs/#{id}")

        envelope = Envelope.wrap_record(
          response,
          resource_name: "Job",
          breadcrumbs: [
            "mobileops jobs list",
            "mobileops vessels list"
          ]
        )

        Formatter.output(envelope, json: options[:json])
      end
    end
  end
end
