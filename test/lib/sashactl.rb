# frozen_string_literal: true

require 'base64'
require 'fileutils'
require 'securerandom'
require 'yaml'

require 'aws-sdk-s3'

require 'sashactl/s3'

module Sashactl
  class << self
    def s3
      @s3 ||= Sashactl::S3.new
    end
  end
end
