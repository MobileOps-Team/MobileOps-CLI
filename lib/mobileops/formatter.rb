# frozen_string_literal: true

require "terminal-table"
require "json"

module Mobileops
  class Formatter
    VESSEL_COLUMNS   = %w[_id name status division_id code archived].freeze
    JOB_COLUMNS      = %w[_id date_beginning date_ending status customer_name].freeze
    CREW_COLUMNS     = %w[_id first_name last_name email employee_number archived].freeze
    COMPONENT_COLUMNS = %w[id name vessel_id hours critical].freeze
    WORK_REQUEST_COLUMNS = %w[_id description status priority vessel_id date].freeze

    class << self
      def output(envelope, json: false)
        json = true if ARGV.include?("--json")
        if json
          puts JSON.pretty_generate(envelope)
          return
        end

        unless envelope[:ok]
          $stderr.puts "Error: #{envelope[:error]}"
          return
        end

        data = envelope[:data]

        if data.is_a?(Array)
          print_table(data)
          puts
          puts envelope[:summary]
        else
          print_record(data)
        end

        print_breadcrumbs(envelope[:breadcrumbs]) if envelope[:breadcrumbs]&.any?
      end

      private

      def print_table(records)
        return puts("No records found.") if records.empty?

        # Use first record's keys as columns
        columns = records.first.keys.first(8)

        table = Terminal::Table.new do |t|
          t.headings = columns.map { |c| c.to_s.gsub("_", " ").upcase }
          records.each do |record|
            t.add_row(columns.map { |c| truncate(record[c].to_s, 30) })
          end
        end

        puts table
      end

      def print_record(record)
        return puts("No record found.") if record.nil? || record.empty?

        max_key_length = record.keys.map(&:to_s).map(&:length).max

        record.each do |key, value|
          label = key.to_s.ljust(max_key_length)
          display_value = format_value(value)
          puts "  #{label}  #{display_value}"
        end
      end

      def print_breadcrumbs(breadcrumbs)
        puts
        puts "Next:"
        breadcrumbs.each do |cmd|
          puts "  → #{cmd}"
        end
      end

      def truncate(str, max)
        str.length > max ? "#{str[0..max - 3]}..." : str
      end

      def format_value(value)
        case value
        when Array
          value.empty? ? "[]" : value.first(5).map(&:to_s).join(", ")
        when Hash
          JSON.generate(value)
        when nil
          "-"
        else
          value.to_s
        end
      end
    end
  end
end
