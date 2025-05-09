# frozen_string_literal: true

Before('@delete') do
  FileUtils.rm_rf('fixtures/delete/articles')
  FileUtils.mkdir_p 'fixtures/delete/articles'
  FileUtils.cp_r 'fixtures/publish/articles/.', 'fixtures/delete/articles'
end

When('we delete an article with slug {string}') do |slug|
  cmd = Nonnative.go_executable(%w[cover], 'reports', '../sashactl', 'delete', '-s', slug, '-i', 'file:.config/delete.yml')
  pid = spawn({}, cmd, %i[out err] => ['reports/delete.log', 'a'])

  _, @status = Process.waitpid2(pid)
end

Then('I should have a deleted article with slug {string}') do |slug|
  expect(File.read('fixtures/delete/articles/articles.yml')).not_to include slug
end

Then('the deleted article with slug {string} should be removed from the file system') do |slug|
  expect(File).to_not exist("fixtures/delete/articles/#{slug}")
end
