# frozen_string_literal: true

Nonnative.configure do |config|
  config.load_file('nonnative.yml')
end

After('@clear_pool') do
  Nonnative.stop
  Nonnative.clear_pool
end
