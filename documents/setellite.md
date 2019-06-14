# Setellite

![Setellite](figures/setellite.jpg?raw=true)


디지털아이디어에서 VFX Onset 툴은 아이패드 어플리케이션 Setellite2(https://www.setellite.nl) 를 사용합니다.
이 툴은 현장정보를 기입하는 App중에서 2018년 현재 가장 유명한 툴이기도 합니다.
현장업무의 복잡함을 줄이기 위해서 디지털아이디어는 VFX 현장의 모든 작업을
아이패드로 진행될 수 있도록 개발할 예정입니다.

#### RollMedia란?

![Rollmedia](figures/rollmedia.png?raw=true)
![RollmediaPoint](figures/rollmediapoint.png?raw=true)

Alexa 카메라의 경우 파일 저장형태는 `00_E002C001_180309_R28L` 입니다.
카메라마다 조금씩 차이는 있지만 보통 비슷한 형태입니다.

- `00` 값은 Reel 넘버입니다. 존재하지 않을 수도 있습니다.
- `E002C001` 값은 우리에게 중요한 RollMedia입니다. 현장툴에서 꼭 기입해줘야합니다.
- `180309` 값은 촬영일입니다.
- `R28L` 값은 촬영시 고유하게 발급되는 ID라고 생각하면 됩니다. 중복값이 없도록 하기 위한 장치입니다.

RollMedia 는 현장에서 촬영시 저장되는 이미지 파일명에 사용됩니다.
촬영시점에는 회사가 작업에 사용하는 SS_0010 형태의 샷이름이 나오기 전이기 때문에,
촬영시 생성되는 이미지 파일명(RollMedia)이 나중에 샷 데이터와 동기화할 때 중요한 정보가 됩니다.
감독 및 스탭 모니터에도 표기되는 문자이기도 합니다.
또한 RollMedia 문자는 EXR / DPX 메타데이터에도 포함되어 있습니다.
우리는 PostProduction에 사용하는 샷이름과 RollMedia를 동기화 하는 방식으로
현장데이터를 확인하게 됩니다.
앱을 이용해서 현장정보를 수집할 때는 RollMedia 기입이 굉장히 중요합니다.
슈퍼바이저 어시스트분들은 RollMedia 기입이 누락되지 않도록 신경써주세요.

#### 준비물
- 아이패드프로
- VPN계정 발급 : IT팀을 통해서 VPN 계정을 생성합니다.(대표이사 결제필요)
- OpenVPN App 설치 : https://itunes.apple.com/kr/app/openvpn-connect/id590379981?mt=8
- OpenVPN에 VPN계정 등록 : IT팀 문의
- Setellite App 설치 : https://itunes.apple.com/kr/app/setellite-2/id956157212?mt=8
- Setellite 가입 : https://www.setellite.nl

#### 작동방식
Setellite를 이용해서 현장정보를 수집한 이후 해당정보를 여러 포멧으로 Export 할 수 있습니다.
회사는 Setellite 앱 내부에서 CSV, PDF 형식으로 저장하는 방식을 사용합니다.
CSV 파일은 CSV > Single 형식으로 여러분의 계정 iCloud에 저장해주세요.
PDF 파일은 해당 프로젝트 PM에게 전달해주세요.
촬영후 숙소에서 CSV, PDF파일을 Export해서 iCloud에 저장하면, 매일 백업이 되니 슈퍼바이저 어시스트 분들은 이 방법을 활용해주세요.
CSV 파일은 아무리 커도 3메가를 넘기 힘듭니다. 해외 느린 인터넷 환경에서도 VPN을 통한 CSI에 업로드 할 수 있습니다.

- 데이터 입력순서
	- 아이패드 : 현장데이터 입력
	- 아이패드 : VPN을 통한 CSI 현장데이터 업로드 http://10.0.90.251/uploadsetellite
	- IO팀(넘벳) : 스캔데이터 처리(스캔데이터가 RollMedia형태로 들어오면 자동으로 현장데이터를 연결해줍니다.)

#### 권장 아이패드 / 촬영시 소음 줄이는 방법

영화 "상해보루"를 진행하면서 데이터가 1000건 이상 아이패드에 쌓일 때
일반사양의 아이패드는 느려지는 현상이 발생했습니다.
권장하는 장비는 2017년 이후 발매된 iPad Pro 입니다

우리나라에서 발매되는 아이패드는 몰래카메라를 방지하기 위해서
촬영시 삑~ 소리가 나게 법으로 규정되어있습니다.
이 소리는 현장에서 작게 나지만,  VFX팀이 현장정보를 수집할 때
작은 소리라도 오디오 녹음팀을 방해할 수 있습니다.
홍콩에서 아이패드를 구매하면 촬영시 소리가 나지 않기때문에
홍콩 해외출장시 아이패드를 구매해서 사용하는 것이 장기적으로 좋습니다.
이 정보는 슈퍼바이저팀과 공유되었습니다.
만약 한국에서 구매한 기기라면 현장에서 사용전에 스피커 구멍부분을
테입으로 막고 촬영장에서 사용해주세요.

#### 데이터 입력시 축약된 약자

현장툴 입력시간 단축을 위해 약자를 많이 사용합니다.
아래는 슈퍼바이져 팀과 약속한 약자입니다.

- OK
- K : Keep
- NG
- CL : Clean Plate
- SRC : Source
- RF : Reference
- EL : Element
- NR : Normal

#### 리눅스 터미널에서 CSV 추가하는 방법

Setellite CSV파일은 CSI의 ![UploadSetellite](http://10.0.90.251/uploadsetellite) 메뉴를 통해서 업로드 할 수 있습니다.
하지만 사용자가 원한다면, 수동으로 터미널에서 setellite2csv 명령어를 이용할 수도 있습니다.
setellite2csv 명령어의 코드 레포지터리는 http://gogs.idea.co.kr/di/setellite2csv 입니다.

터미널 사용법

```
$ setellite2csi -project [프로젝트명] -csv [csvfilename.csv] -overwrite [true|false]
```

#### 소프트웨어 유지비
- 2018년 기준. 슈퍼바이저 20명 기준 년 약 400만원 선

#### 아직 구현되지 않는 항목
- Red 카메라 Rollmedia 인식

#### Link
- Setellite 트위터 : https://twitter.com/setellite_app
- 개발요청 및 버그보고 : info@setellite.nl

#### History
- 2018.01 ~ : Setellite2 와 CSI의 연동부분 개발시작
- 2018.01 : TakeD 프로젝트는 향후 발전을 위해서 OpenSource로 변경함. https://github.com/didev/taked
- 2017.12 : TakeD 개발자 퇴사로 Setellite로 변경.
- 2017.05 : 아이패드용 현장툴 프로젝트 TakeD 시작. http://gogs.idea.co.kr/legacy/taked
- 2017.03 : 슈퍼바이저팀 요청. 아이패드 현장툴
