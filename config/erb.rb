#!/usr/bin/ruby

require 'erb'

endpoints = [10081, 10082, 10083]

template = ERB.new(File.read('envoy.yml.erb'), nil, '-')

File.write('envoy.yml', template.result(binding))
