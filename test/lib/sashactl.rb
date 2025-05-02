# frozen_string_literal: true

require 'securerandom'
require 'yaml'
require 'base64'

require 'sashactl/v1/server'

module Sashactl
  module V1
    class << self
      def bucket_operational?
        bucket_state == 'operational'
      end

      def bucket_missing?
        bucket_state == 'missing'
      end

      def bucket_erroneous?
        bucket_state == 'erroneous'
      end

      def apply_bucket_state(state)
        ENV['BUCKET_STATE'] = state
      end

      def bucket_state
        ENV.fetch('BUCKET_STATE', nil)
      end
    end
  end
end
