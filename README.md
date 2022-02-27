A playground to experiment with how client-side Envoy config can affect the performance of
gRPC applications.

### Usage

All config options are exposed in the top level `config.yml` file. Adjust a config option,
and then run `bin/gen.rb` to have the options take affect. You can also run `bin/delay.sh`
to inject a latency of 100ms between Envoy and the backends (to simulate a real network).

Once you've generated the config, you can start the sim with `docker-compose up`. A
grafana dashboard is exposed at `localhost:3000`, and Envoy's admin interface is available
at `localhost:9901`.

For simplicity, containers are running using host networking (so beware of port
collisions). Linux is required for the use of traffic control and cadvisor.
