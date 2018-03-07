# tiny-site

## docker

### docker build

```bash
docker build -t vicanso/tiny-site .
```

### docker run
```bash
docker run -d --restart=always \
  -p 5018:5018 \
  -e GRPC=172.17.0.1:3016 \
  -e NODE_ENV=production \
  vianso/tiny-site
```