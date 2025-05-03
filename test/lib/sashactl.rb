# frozen_string_literal: true

require 'securerandom'
require 'yaml'
require 'base64'

require 'aws-sdk-s3'

require 'sashactl/s3'

module Sashactl
  class << self
    def s3
      @s3 ||= Sashactl::S3.new
    end
  end
end
