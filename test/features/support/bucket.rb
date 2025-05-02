# frozen_string_literal: true

Before('@operational') do
  Sashactl::V1.apply_bucket_state 'operational'
end

Before('@missing') do
  Sashactl::V1.apply_bucket_state 'missing'
end

Before('@erroneous') do
  Sashactl::V1.apply_bucket_state 'erroneous'
end
