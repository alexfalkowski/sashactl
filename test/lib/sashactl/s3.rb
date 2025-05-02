# frozen_string_literal: true

module Sashactl
  class S3
    def initialize
      credentials = Aws::Credentials.new('access', 'secret')
      @client = Aws::S3::Client.new(endpoint: 'http://localhost:4566', credentials:, region: 'eu-west-1', force_path_style: true)
    end

    def create
      client.head_bucket(bucket: 'sasha-cms')
    rescue StandardError
      client.create_bucket(bucket: 'sasha-cms')
    end

    def delete
      bucket = Aws::S3::Bucket.new('sasha-cms', client: client)
      bucket.delete!

      true
    rescue StandardError
      false
    end

    def exists?(path)
      client.head_object(bucket: 'sasha-cms', key: path)

      true
    rescue StandardError
      false
    end

    private

    attr_reader :client
  end
end
