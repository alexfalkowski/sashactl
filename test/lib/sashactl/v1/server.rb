# frozen_string_literal: true

module Sashactl
  module V1
    class OperationalBucket < Sinatra::Application
      set :public_folder, 'server'
    end

    class MissingBucket < Sinatra::Application
      get '/articles.yml' do
        status 404
      end
    end

    class ErroneousBucket < Sinatra::Application
      get '/articles.yml' do
        status 500
      end
    end

    class Server < Nonnative::HTTPServer
      def initialize(service)
        if Sashactl::V1.bucket_erroneous?
          super(Sinatra.new(ErroneousBucket), service)

          return
        end

        if Sashactl::V1.bucket_missing?
          super(Sinatra.new(MissingBucket), service)

          return
        end

        super(Sinatra.new(OperationalBucket), service)
      end
    end
  end
end
