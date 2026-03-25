# frozen_string_literal: true

module Mobileops
  module Commands
    class Vessels < Thor
      include Pagination

      class_option :json, type: :boolean, default: false, desc: "Output as JSON"

      desc "list", "List all vessels/assets"
      def list
        client = Client.new
        response = client.get("assets", pagination_params)

        envelope = Envelope.wrap_collection(
          response,
          resource_name: "vessels",
          breadcrumbs: response["data"]&.first(3)&.map { |v|
            id = v["_id"] || v["id"]
            "mobileops vessels get #{id}"
          } || []
        )

        Formatter.output(envelope, json: options[:json])
      end

      desc "get ID", "Get a single vessel by ID"
      def get(id)
        client = Client.new
        response = client.get("assets/#{id}")

        vessel_id = response["_id"] || response["id"]
        envelope = Envelope.wrap_record(
          response,
          resource_name: "Vessel",
          breadcrumbs: [
            "mobileops vessels specs #{vessel_id}",
            "mobileops components list --vessel-id #{vessel_id}",
            "mobileops jobs list --vessel-id #{vessel_id}",
            "mobileops work-requests list --vessel-id #{vessel_id}"
          ]
        )

        Formatter.output(envelope, json: options[:json])
      end

      desc "specs VESSEL_ID", "List specs for a vessel"
      def specs(vessel_id)
        client = Client.new
        response = client.get("vessel-specs", { vessel_id: vessel_id })

        envelope = Envelope.wrap_collection(
          response,
          resource_name: "vessel specs",
          breadcrumbs: [
            "mobileops vessels get #{vessel_id}",
            "mobileops components list --vessel-id #{vessel_id}"
          ]
        )

        Formatter.output(envelope, json: options[:json])
      end
    end
  end
end
