package goffmpeg

/*
需求:
	1. 在二进制里内置静态编译后的ffmpeg, 使用时候直接释放到临时目录, 以便于支持windows/linux/darwin平台的视频转换
	2. 不依赖cgo, 方便跨平台编译
	3. 每个平台embed自己平台的ffmpeg二进制, 不放其他平台的
	4. 提供简单的接口, 实现目前遇到的通用需求
	5. 需要兼容单机上多个进程导入本package然后同时使用的情况
*/

