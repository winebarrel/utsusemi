# utsusemi

utsusemi is a proxy server that automatically switches backend.

When the first backend returns 404, it goes to the next backend. If the next backend returns 404, go to the third backend...

## Getting Started

```sh
cp utsusemi.toml.sample utsusemi.toml
make go-get
make
./utsusemi
```
```sh
curl localhost:11080/ping
open http://localhost:11080/images/top/sp2/cmn/logo-170307.png
open http://localhost:11080/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png
```

## Configuration

```toml
#port = 11080

[[backend]]
target = "http://httpstat.us/404"
#ok = [200]

[[backend]]
target = "https://www.google.co.jp"
#ok = [200]

[[backend]]
target = "https://s.yimg.jp"
#ok = [200]
```
