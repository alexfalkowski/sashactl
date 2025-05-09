# frozen_string_literal: true

Before('@unpublish') do
  FileUtils.rm_rf('reports/unpublish/articles')
  FileUtils.mkdir_p 'reports/unpublish/articles'
  FileUtils.cp_r 'reports/publish/articles/.', 'reports/unpublish/articles'
end

When('we unpublish an article with slug {string}') do |slug|
  cmd = Nonnative.go_executable(%w[cover], 'reports', '../sashactl', 'unpublish', '-s', slug, '-i', 'file:.config/unpublish.yml')
  pid = spawn({}, cmd, %i[out err] => ['reports/unpublish.log', 'a'])

  _, @status = Process.waitpid2(pid)
end

Then('I should have an unpublished article with slug {string}') do |slug|
  expect(Sashactl.s3.exists?("#{slug}/article.yml")).to be false
  expect(Sashactl.s3.exists?("#{slug}/images/1984.jpeg")).to be false
  expect(File.read('reports/unpublish/articles/articles.yml')).not_to include slug
end

Then('the unpublished article with slug {string} should be removed from the file system') do |slug|
  expect(File).to_not exist("reports/unpublish/articles/#{slug}")
end
