# frozen_string_literal: true

module Mobileops
  module Envelope
    module_function

    def wrap(data:, summary:, breadcrumbs: [])
      {
        ok: true,
        data: data,
        summary: summary,
        breadcrumbs: breadcrumbs
      }
    end

    def error(message:, breadcrumbs: [])
      {
        ok: false,
        error: message,
        breadcrumbs: breadcrumbs
      }
    end

    # Wrap an index response (API returns { data: [...], total: N })
    def wrap_collection(response, resource_name:, breadcrumbs: [])
      data  = response["data"] || response
      total = response["total"]

      summary = if total
                  "Showing #{data.size} of #{total} #{resource_name}"
                else
                  "Found #{data.size} #{resource_name}"
                end

      wrap(data: data, summary: summary, breadcrumbs: breadcrumbs)
    end

    # Wrap a show response (API returns flat object)
    def wrap_record(record, resource_name:, breadcrumbs: [])
      id = record["_id"] || record["id"]
      summary = "#{resource_name} #{id}"

      wrap(data: record, summary: summary, breadcrumbs: breadcrumbs)
    end
  end
end
