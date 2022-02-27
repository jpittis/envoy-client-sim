A playground to experiment with how client-side Envoy config can affect the performance of
gRPC applications (I'm mostly interested in latency).

### Usage

Adjust the number of backend endpoints by editing the config at the top of `config/gen.rb`.

```
num_endpoints   = 5
base_port       = 10081
lb_policy       = 'ROUND_ROBIN'
connect_timeout = '1s'
idle_timeout    = '60s' # Defaults to 1h
```

Run the generate script:

```
./config/gen.rb
```

Then start the sim with docker-compose:

```
docker-compose up
```

The default client behavior is to round robin across each endpoint:

```
Success! (name=10084, duration=410.644349ms)
Success! (name=10085, duration=407.237265ms)
Success! (name=10081, duration=404.425664ms)
Success! (name=10082, duration=407.27783ms)
Success! (name=10083, duration=406.061373ms)
Success! (name=10084, duration=203.802976ms)
```

You can also inject latency between the client side Envoy and the backends using a sketchy
tc script I cooked up (which currently defaults to 100ms):

```
./bin/delay.sh
```

Envoy worker concurrency can be configured in `config/Docker-envoy`.

A grafana dashboard is available at `localhost:3000`.

The usage of cadvisor means this project will likely not work on OSX, only Linux.
