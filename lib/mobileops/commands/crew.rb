# frozen_string_literal: true

module Mobileops
  module Commands
    class Crew < Thor
      include Pagination

      class_option :json, type: :boolean, default: false, desc: "Output as JSON"

      desc "list", "List crew members"
      def list
        client = Client.new
        response = client.get("users", pagination_params)

        envelope = Envelope.wrap_collection(
          response,
          resource_name: "crew members",
          breadcrumbs: response["data"]&.first(3)&.map { |u|
            id = u["_id"] || u["id"]
            "mobileops crew get #{id}"
          } || []
        )

        Formatter.output(envelope, json: options[:json])
      end

      desc "get ID", "Get a single crew member by ID"
      def get(id)
        client = Client.new
        response = client.get("users/#{id}")

        user_id = response["_id"] || response["id"]
        name = [response["first_name"], response["last_name"]].compact.join(" ")

        envelope = Envelope.wrap_record(
          response,
          resource_name: "Crew member",
          breadcrumbs: [
            "mobileops crew list",
            "mobileops jobs list"
          ]
        )

        Formatter.output(envelope, json: options[:json])
      end

      desc "find-by-employee-number NUMBER", "Find a crew member by employee number"
      map "find-by-employee-number" => :find_by_employee_number
      def find_by_employee_number(number)
        client = Client.new
        response = client.get("users/by_employee_number", { employee_number: number })

        envelope = Envelope.wrap_record(
          response,
          resource_name: "Crew member",
          breadcrumbs: [
            "mobileops crew list"
          ]
        )

        Formatter.output(envelope, json: options[:json])
      end
    end
  end
end
