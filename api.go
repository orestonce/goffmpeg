package goffmpeg

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
)

func SetupFfmpeg() (ffmpegExePath string, err error) {
	ffmpegExePath = filepath.Join(os.TempDir(), "ffmpeg-static")
	if runtime.GOOS == "windows" {
		ffmpegExePath = ffmpegExePath + ".exe"
	}
	info, err := os.Stat(ffmpegExePath)
	if err == nil && info.Mode().IsRegular() {
		if runtime.GOOS != "windows" {
			err = os.Chmod(ffmpegExePath, 0777)
			if err != nil {
				return "", err
			}
		}
		return ffmpegExePath, nil
	}
	reader, err := gzip.NewReader(bytes.NewReader(gEmbedFfmpegGz))
	if err != nil {
		return "", err
	}
	tmpName := ffmpegExePath + "." + strconv.Itoa(os.Getpid()) + ".tmp"
	fout, err := os.OpenFile(tmpName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(fout, reader)
	if err != nil {
		_ = fout.Close()
		_ = os.Remove(tmpName)
		return "", err
	}
	err = fout.Close()
	if err != nil {
		_ = os.Remove(tmpName)
		return "", err
	}
	err = os.Rename(tmpName, ffmpegExePath)
	if err != nil {
		_ = os.Remove(tmpName)
		return "", err
	}
	return ffmpegExePath, nil
}

type MergeMultiToSingleMp4_Req struct {
	FfmpegExePath string
	TsFileList    []string
	OutputMp4     string
	ProgressCh    chan<- int // 百分比, [0, 100]
}

func MergeMultiToSingleMp4(req MergeMultiToSingleMp4_Req) (err error) {
	outputMp4Temp := req.OutputMp4 + ".tmp"
	cmd := exec.Command(req.FfmpegExePath, "-i", "-", "-acodec", "copy", "-vcodec", "copy", "-f", "mp4", "-y", outputMp4Temp)
	setupCmd(cmd)
	ip, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if req.ProgressCh != nil {
		defer close(req.ProgressCh)
	}
	cmd.Stdout = os.Stdout
	//setupCmd(cmd)
	err = cmd.Start()
	if err != nil {
		return err
	}
	for idx, name := range req.TsFileList {
		fin, err := os.Open(name)
		if err != nil {
			cmd.Process.Kill()
			cmd.Wait()
			os.Remove(outputMp4Temp)
			return errors.New("read error: " + err.Error())
		}
		_, err = io.Copy(ip, fin)
		if err != nil {
			cmd.Process.Kill()
			cmd.Wait()
			fin.Close()
			os.Remove(outputMp4Temp)
			return errors.New("write error: " + err.Error())
		}
		fin.Close()
		if req.ProgressCh != nil {
			req.ProgressCh <- idx * 100 / len(req.TsFileList)
		}
	}
	err = ip.Close()
	if err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		os.Remove(outputMp4Temp)
		return errors.New("ip.Close error: " + err.Error())
	}
	err = cmd.Wait()
	if err != nil {
		os.Remove(outputMp4Temp)
		return err
	}
	err = os.Rename(outputMp4Temp, req.OutputMp4)
	if err != nil {
		os.Remove(outputMp4Temp)
		return err
	}
	return nil
}

func MustShowHelp(exePath string) {
	cmd := exec.Command(exePath, "-h")
	setupCmd(cmd)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
