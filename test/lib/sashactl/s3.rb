# frozen_string_literal: true

module Sashactl
  class S3
    def initialize
      credentials = Aws::Credentials.new('access', 'secret')
      @client = Aws::S3::Client.new(endpoint: 'http://localhost:4566', credentials:, region: 'eu-west-1', force_path_style: true)
    end

    def create
      client.head_bucket(bucket: 'articles')
    rescue StandardError
      client.create_bucket(bucket: 'articles')
    end

    def delete
      bucket = Aws::S3::Bucket.new('articles', client: client)
      bucket.delete!

      true
    rescue StandardError
      false
    end

    def exists?(path)
      client.head_object(bucket: 'articles', key: path)

      true
    rescue StandardError
      false
    end

    private

    attr_reader :client
  end
end
