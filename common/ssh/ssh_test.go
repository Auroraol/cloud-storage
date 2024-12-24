package sshx

import (
	"context"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_ssh(t *testing.T) {
	// 创建鉴权秘钥
	privateKeyConf := Credential{User: "root", Password: "-+66..[]l"}

	// 创建 sshx 客户端
	rs, err := NewClient("101.37.165.220:22", privateKeyConf, SetEstablishTimeout(10*time.Second), SetLogger(DefaultLogger{}))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rs.Handle(func(sub EnhanceClient) error {
		// if _, err := sub.ReceiveFile("/tmp/xxx", "/etc/passwd", false, true); err != nil {
		// 	panic(err)
		// }

		path := "/opt/goTest/log"

		data, err := sub.ReadFile(path)
		if err != nil {
			panic(err)
		}

		log.Printf("filebeat yml: %s", string(data))
		//
		//if err := sub.TempSendFile("./bin/cmd-simulator", func(tempFilepath string) error {
		//	res, err := sub.Command(ctx, fmt.Sprintf("chmod +x ./%s && ./%s -stdout Hello -stderr 'This is an error' -return-code 3", tempFilepath, tempFilepath))
		//	if err != nil {
		//		return fmt.Errorf("%v: %v", err, string(res))
		//	}
		//
		//	log.Printf("cmd-simulator: %v", string(res))
		//	return nil
		//}); err != nil {
		//	return err
		//}
		cmd := "wc -l " + path
		output, err := sub.Command(ctx, cmd)
		if err != nil {
			return err
		}
		lines := strings.Split(string(output), " ")
		totalLines, err := strconv.Atoi(strings.TrimSpace(lines[0]))
		if err != nil {
			return err
		}
		log.Printf("totalLines: %d", totalLines)
		//
		//log.Printf("whoami: %s", whoami)
		//
		//dataReader := strings.NewReader("Hello, world")
		//_, err = sub.WriteFileOverride("/root/test.txt", dataReader)
		//if err != nil {
		//	return err
		//}
		//
		//hello, err := sub.Command(ctx, "cat /root/test.txt")
		//if err != nil {
		//	return err
		//}
		//
		//log.Printf("res: %s", hello)
		//
		//if err := sub.Remove("/root/test.txt"); err != nil {
		//	return err
		//}
		//
		//psef, err := sub.Command(ctx, "ps -ef", RequestPty(120, 100))
		//if err != nil {
		//	return err
		//}
		//log.Printf("ps ef: %s", string(psef))
		//
		//if err := sub.TempWriteFile(strings.NewReader("Yes!"), func(tempFilepath string) error {
		//	log.Printf("temp file: %s", tempFilepath)
		//	res, err := sub.Command(context.TODO(), "cat "+tempFilepath)
		//	if err != nil {
		//		return err
		//	}
		//
		//	log.Printf("temp file content: %s", string(res))
		//
		//	return nil
		//}); err != nil {
		//	return err
		//}

		//if err := sub.SendDirectory("/root/temp", "/Users/mylxsw/codes/github/sshx"); err != nil {
		//	return err
		//}

		return nil
	}); err != nil {
		panic(err)
	}
}
