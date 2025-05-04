# frozen_string_literal: true

Before('@delete') do
  FileUtils.rm_rf('reports/delete/articles')
  FileUtils.mkdir_p 'reports/delete/articles'
  FileUtils.cp_r 'reports/publish/articles/.', 'reports/delete/articles'
end

Given('we have a published article with slug {string}') do |slug|
  cmd = Nonnative.go_executable(%w[cover], 'reports', '../sashactl', 'publish', '-s', slug, '-i', 'file:.config/delete.yml')
  pid = spawn({}, cmd, %i[out err] => ['reports/delete.log', 'a'])
  _, status = Process.waitpid2(pid)

  expect(status.exitstatus).to eq(0)
end

When('we delete an article with slug {string}') do |slug|
  cmd = Nonnative.go_executable(%w[cover], 'reports', '../sashactl', 'delete', '-s', slug, '-i', 'file:.config/delete.yml')
  pid = spawn({}, cmd, %i[out err] => ['reports/delete.log', 'a'])

  _, @status = Process.waitpid2(pid)
end

Then('I should have a deleted article with slug {string}') do |slug|
  expect(Sashactl.s3.exists?("#{slug}/article.yml")).to be false
  expect(Sashactl.s3.exists?("#{slug}/images/1984.jpeg")).to be false
  expect(File.read('reports/delete/articles/articles.yml')).not_to include slug
end

Then('the article with slug {string} should be removed from the file system') do |slug|
  expect(File).to_not exist("reports/delete/articles/#{slug}")
end
