# 第五届字节跳动青训营“抖声”项目

gogogo 团队

2023.2.10 jijinkang 简单说明：完成 Feed、publish/list 两个接口业务逻辑

- config：SQL 等配置
- data：数据声明
- initDAO：数据库初始化及建表

  2023.2.11 jijinkang

- FFmpeg: 安装下载地址：https://www.gyan.dev/ffmpeg/builds/#release-builds [ffmpeg-5.0.1-essentials_build.zip ] (https://www.gyan.dev/ffmpeg/builds/packages/ffmpeg-5.0.1-essentials_build.zip)77 MB [.sha256](https://www.gyan.dev/ffmpeg/builds/packages/ffmpeg-5.0.1-essentials_build.zip.sha256)

- MinIO:windows 服务端下载地址： https://dl.minio.io/server/minio/release/windows-amd64/minio.exe 命令` .\minio.exe server [存储目录]`

- **省略中途本地存储的步骤，视频流直接转存到 MinIO， 使用 URL 取帧获得 cover**，不知道能否使用视频流获得 cover

  2023.2.13 wangmingxian 简单说明：完成 Register, Login, User 三个部分及 jwt 鉴权

  2023.2.18 zk 完成/douyin/favorite/action/接口，继续写/douyin/favorite/list/接口需要 Video 和 User 的完整的数据库

补充：config MinIO配置Endpoint应改为服务器（即本机）ip，即可正常播放视频

