#!/usr/bin/ruby

require 'erb'
require 'yaml'

config = YAML.load_file('config.yml')

num_endpoints                 = config['num_endpoints']
base_port                     = config['base_port']
lb_policy                     = config['lb_policy']
connect_timeout               = config['connect_timeout']
idle_timeout                  = config['idle_timeout']
enable_exact_balance          = config['enable_exact_balance']
concurrency                   = config['concurrency']
enable_preconnect             = config['enable_preconnect']
per_upstream_preconnect_ratio = config['per_upstream_preconnect_ratio']
predictive_preconnect_ratio   = config['predictive_preconnect_ratio']

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
