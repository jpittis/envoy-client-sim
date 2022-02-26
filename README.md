A playground to experiment with how client-side Envoy config can affect the performance of
gRPC applications (I'm mostly interested in latency).

### Usage

Adjust the number of backend endpoints by editing the config at the top of `config/erb.rb`.

Run the generate script:

```
cd config && ./erb.rb
```

Then start the sim with docker-compose:

```
docker-compose up
```

The default client behavior is to round robin across each endpoint:

```
envoy-client-sim-backend-1  | 2022/02/26 01:10:55 Success! (10085)
envoy-client-sim-backend-1  | 2022/02/26 01:10:56 Success! (10081)
envoy-client-sim-backend-1  | 2022/02/26 01:10:57 Success! (10082)
envoy-client-sim-backend-1  | 2022/02/26 01:10:58 Success! (10083)
envoy-client-sim-backend-1  | 2022/02/26 01:10:59 Success! (10084)
envoy-client-sim-backend-1  | 2022/02/26 01:11:00 Success! (10085)
envoy-client-sim-backend-1  | 2022/02/26 01:11:01 Success! (10081)
```
