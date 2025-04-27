# frozen_string_literal: true

When('we run the client') do
  cmd = Nonnative.go_executable(%w[cover], 'reports', '../go-client-template', 'client', '-i', 'file:.config/client.yml')
  pid = spawn({}, cmd, %i[out err] => ['reports/client.log', 'a'])

  _, @status = Process.waitpid2(pid)
end

Then('it should run sucessfully') do
  expect(@status.exitstatus).to eq(0)
end
