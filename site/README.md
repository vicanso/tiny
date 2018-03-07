# tiny-site

## docker

### docker build

```bash
docker build -t vicanso/tiny-site .
```

### docker run
```bash
docker run -d --restart=always \
  -p 5031:5018 \
  -e GRPC=172.17.0.1:3016 \
  -e NODE_ENV=production \
  vicanso/tiny-site
```