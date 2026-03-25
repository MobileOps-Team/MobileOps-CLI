# frozen_string_literal: true

RSpec.describe Mobileops::Envelope do
  describe ".wrap" do
    it "wraps data with ok status" do
      result = described_class.wrap(data: { name: "test" }, summary: "Test result")
      expect(result[:ok]).to be true
      expect(result[:data][:name]).to eq("test")
      expect(result[:summary]).to eq("Test result")
    end
  end

  describe ".error" do
    it "wraps error message" do
      result = described_class.error(message: "Something went wrong")
      expect(result[:ok]).to be false
      expect(result[:error]).to eq("Something went wrong")
    end
  end

  describe ".wrap_collection" do
    it "generates summary with total when available" do
      response = { "data" => [{ "name" => "V1" }, { "name" => "V2" }], "total" => 10 }
      result = described_class.wrap_collection(response, resource_name: "vessels")

      expect(result[:summary]).to eq("Showing 2 of 10 vessels")
      expect(result[:data].size).to eq(2)
    end

    it "generates summary without total" do
      response = { "data" => [{ "name" => "V1" }] }
      result = described_class.wrap_collection(response, resource_name: "vessels")

      expect(result[:summary]).to eq("Found 1 vessels")
    end
  end

  describe ".wrap_record" do
    it "generates summary with record ID" do
      record = { "_id" => "abc123", "name" => "MV Atlantic" }
      result = described_class.wrap_record(record, resource_name: "Vessel")

      expect(result[:summary]).to eq("Vessel abc123")
      expect(result[:data]["name"]).to eq("MV Atlantic")
    end
  end
end
