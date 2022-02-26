#!/usr/bin/ruby

require 'erb'

num_endpoints = 4
base_port     = 10081

endpoints = []
next_port = base_port
num_endpoints.times do
  endpoints << next_port
  next_port += 1
end

# Generate Envoy config with right number of endpoints
template = ERB.new(File.read('envoy.yml.erb'), nil, '-')
File.write('envoy.yml', template.result(binding))

# Write out the ports to a config file to be consumed by backend service
File.write('endpoints.txt', endpoints.join(','))
