# Brew

- <https://brew.sh>

```bash
brew uninstall mongodb
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb-community
brew install curl
brew install ffmpeg
brew install openimageio
brew install go
```

# Macport

구형 macOS 에서는 brew가 지원되지 않습니다. macport를 이용해서 필요한 라이브러리와 SW를 설치합니다.
하지만 시에라에서는 mongodb가 설치되지 않습니다. 하이시에라 이전버전이라면 리눅스를 권장해드립니다.

- <https://www.macports.org/install.php>

```bash
sudo port install mongodb
sudo port install ffmpeg
sudo port install openimageio
sudo port install curl
sudo port install go
```
