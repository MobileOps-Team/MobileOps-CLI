# frozen_string_literal: true

require "tmpdir"

RSpec.describe Mobileops::Config do
  let(:tmp_dir) { Dir.mktmpdir }
  let(:credentials_file) { File.join(tmp_dir, "credentials.json") }

  before do
    stub_const("Mobileops::Config::CONFIG_DIR", tmp_dir)
    stub_const("Mobileops::Config::CREDENTIALS_FILE", credentials_file)
  end

  after do
    FileUtils.remove_entry(tmp_dir)
  end

  describe ".save and .load" do
    it "round-trips credentials" do
      described_class.save(
        host: "https://test.mobileops.com",
        access_key: "ak-123",
        secret_key: "sk-456"
      )

      config = described_class.load
      expect(config[:host]).to eq("https://test.mobileops.com")
      expect(config[:access_key]).to eq("ak-123")
      expect(config[:secret_key]).to eq("sk-456")
    end

    it "sets restrictive file permissions" do
      described_class.save(host: "https://test.com", access_key: "ak", secret_key: "sk")

      mode = File.stat(credentials_file).mode & 0o777
      expect(mode).to eq(0o600)
    end
  end

  describe ".exists?" do
    it "returns false when no credentials exist" do
      expect(described_class.exists?).to be false
    end

    it "returns true after saving" do
      described_class.save(host: "https://test.com", access_key: "ak", secret_key: "sk")
      expect(described_class.exists?).to be true
    end
  end

  describe ".delete" do
    it "removes credentials file" do
      described_class.save(host: "https://test.com", access_key: "ak", secret_key: "sk")
      described_class.delete
      expect(described_class.exists?).to be false
    end
  end

  describe ".load" do
    it "raises ConfigError when no file exists" do
      expect { described_class.load }.to raise_error(Mobileops::ConfigError)
    end
  end
end
