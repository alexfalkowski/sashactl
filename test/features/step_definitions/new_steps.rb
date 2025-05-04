# frozen_string_literal: true

Before('@new') do
  FileUtils.rm_rf('reports/new/articles')
end

When('we create an article with name {string}') do |name|
  cmd = Nonnative.go_executable(%w[cover], 'reports', '../sashactl', 'new', '-n', "\"#{name}\"", '-i', 'file:.config/new.yml')
  pid = spawn({}, cmd, %i[out err] => ['reports/new.log', 'a'])

  _, @status = Process.waitpid2(pid)
end

Then('it should run successfully') do
  expect(@status.exitstatus).to eq(0)
end

Then('I should have an article named {string}') do |name|
  name = name.split.join('-')

  expect(File).to exist("reports/new/articles/#{name}")
end

Then('it should not run successfully') do
  expect(@status.exitstatus).to eq(1)
end

Then('I should not have an article named {string}') do |name|
  name = name.split.join('-')

  expect(File).to_not exist("reports/new/articles/#{name}")
end
