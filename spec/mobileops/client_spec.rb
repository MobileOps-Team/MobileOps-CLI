# frozen_string_literal: true

RSpec.describe Mobileops::Client do
  let(:host) { "https://www.mobileops.at" }
  let(:access_key) { "test-access-key" }
  let(:secret_key) { "test-secret-key" }
  let(:client) { described_class.new(host: host, access_key: access_key, secret_key: secret_key) }

  describe "#get" do
    it "returns parsed JSON on success" do
      stub_request(:get, "#{host}/api/assets")
        .with(headers: {
          "X-Api-Access-Key" => access_key,
          "X-Api-Secret-Key" => secret_key
        })
        .to_return(
          status: 200,
          body: { data: [{ name: "MV Test" }], total: 1 }.to_json,
          headers: { "Content-Type" => "application/json" }
        )

      result = client.get("assets")
      expect(result["data"].first["name"]).to eq("MV Test")
    end

    it "raises AuthError on 401" do
      stub_request(:get, "#{host}/api/assets")
        .to_return(status: 401, body: "", headers: {})

      expect { client.get("assets") }.to raise_error(Mobileops::AuthError)
    end

    it "raises PermissionError on 403" do
      stub_request(:get, "#{host}/api/assets")
        .to_return(status: 403, body: "", headers: {})

      expect { client.get("assets") }.to raise_error(Mobileops::PermissionError)
    end

    it "raises NotFoundError on 404" do
      stub_request(:get, "#{host}/api/assets/bad-id")
        .to_return(
          status: 404,
          body: { error: "Not found" }.to_json,
          headers: { "Content-Type" => "application/json" }
        )

      expect { client.get("assets/bad-id") }.to raise_error(Mobileops::NotFoundError)
    end

    it "raises ApiError on 500" do
      stub_request(:get, "#{host}/api/assets")
        .to_return(status: 500, body: "", headers: {})

      expect { client.get("assets") }.to raise_error(Mobileops::ApiError)
    end
  end

  describe "#initialize" do
    it "raises AuthError when no credentials are provided" do
      allow(Mobileops::Config).to receive(:load).and_raise(Mobileops::ConfigError)

      expect { described_class.new }.to raise_error(Mobileops::AuthError)
    end
  end
end
