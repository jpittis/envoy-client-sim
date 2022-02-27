#!/usr/bin/ruby

require 'erb'

num_endpoints        = 5
base_port            = 10081
lb_policy            = 'ROUND_ROBIN'
connect_timeout      = '1s'
idle_timeout         = '60s' # Defaults to 1h
enable_exact_balance = true
concurrency          = 4

endpoints = []
next_port = base_port
num_endpoints.times do
  endpoints << next_port
  next_port += 1
end

# Generate Envoy config with right number of endpoints
template = ERB.new(File.read('config/envoy.yml.erb'), nil, '-')
File.write('config/envoy.yml', template.result(binding))

# Write out the ports to a config file to be consumed by backend service
File.write('config/endpoints.txt', endpoints.join(','))

# Controls number of Envoy workers on process start
File.write('config/concurrency.txt', concurrency)
