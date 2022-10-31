
# 하모니카 리눅스

macmini 구형의 경우 하모니카 리눅스를 설치하고 운용하면 편리하게 운용할 수 있다.


```bash
apt-get install ffmpeg
apt-get install mongodb
apt-get install golang
apt-get install curl
apt-get install openimageio
```


방화벽 열기

```bash
sudo ufw allow 80/tcp comment ‘http’
sudo ufw allow 443/tcp comment ‘https’
```