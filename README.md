# utsusemi

utsusemi is a proxy server that automatically switches backend.

When the first backend returns 404, it goes to the next backend. If the next backend returns 404, go to the third backend...

## Getting Started

```sh
cp utsusemi.toml.sample utsusemi.toml
make
./utsusemi
```
```sh
curl localhost:11080/ping
curl localhost:11080
```

## Configuration

```toml
#port = 11080

[[backend]]
target = "http://httpstat.us/404"
#ok = [200]

[[backend]]
target = "https://httpstat.us/404"
#ok = [200]

[[backend]]
target = "https://httpstat.us/200"
#ok = [200]
```
