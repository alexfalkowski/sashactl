# frozen_string_literal: true

Before do
  Sashactl.s3.delete
  Sashactl.s3.create
end
