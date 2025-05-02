# frozen_string_literal: true

Before do
  Sashactl.s3.delete
  Sashactl.s3.create
end

When('we publish an article with slug {string}') do |slug|
  cmd = Nonnative.go_executable(%w[cover], 'reports', '../sashactl', 'publish', '-s', slug, '-i', 'file:.config/publish.yml')
  pid = spawn({}, cmd, %i[out err] => ['reports/publish.log', 'a'])

  _, @status = Process.waitpid2(pid)
end

Then('I should have a published article with slug {string}') do |slug|
  expect(Sashactl.s3.exists?('articles.yml')).to be true
  expect(Sashactl.s3.exists?("#{slug}/article.yml")).to be true
  expect(Sashactl.s3.exists?("#{slug}/images/1984.jpeg")).to be true
end

Then('I should not have a published article with slug {string}') do |slug|
  expect(Sashactl.s3.exists?('articles.yml')).to be false
  expect(Sashactl.s3.exists?("#{slug}/article.yml")).to be false
  expect(Sashactl.s3.exists?("#{slug}/images/1984.jpeg")).to be false
end
