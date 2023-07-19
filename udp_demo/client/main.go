package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var change chan string = make(chan string)

func main() {
	//建立链接
	c, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 30000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	//监听来自服务端的消息
	go print(c)

	//登录
	var name string
	for {
		flag := 0
		fmt.Println("按1登录\n按2注册\n按3离开")
		i := bufio.NewReader(os.Stdin)
		number, _ := i.ReadString('\n')
		if number == "1\r\n" {
			for {
				name = land(c)
				if name == "/e" {
					break
				} else if name == "" {
					continue
				} else {
					flag = 1
					break
				}
			}
		} else if number == "3\r\n" {
			return
		} else if number == "2\r\n" {
			register(c)
		}
		if flag == 1 {
			break
		}
	}

	//告知服务端上线
	uplimit(name, c)
	//选择链接的用户
	for {

		fmt.Println("按@选择聊天对象\n按/check查看在线人员\n按/q退出\n按/pe查看个人中心")
		input := bufio.NewReader(os.Stdin)
		sel, _ := input.ReadString('\n')
		if sel == "@\r\n" {
			selec(c)
		} else if sel == "/check\r\n" {
			checkid(c)
		} else if sel == "/q\r\n" {
			fmt.Print("再见")
			_, err := c.Write([]byte(sel))
			if err != nil {
				fmt.Println(err)
				return
			}
			os.Exit(0)
		} else if sel == "/pe\r\n" {
			pepor(&name, c)
		}

	}

}

// 用户中心
func pepor(name *string, c *net.UDPConn) {
	fmt.Println("用户名：" + *name)
	fmt.Println("按/rename修改用户名:")
	fmt.Println("按/repass修改密码")
	input := bufio.NewReader(os.Stdin)
	sel, _ := input.ReadString('\n')
	if sel == "/rename\r\n" {
		*name = rena(c, *name)
	} else if sel == "/repass\r\n" {
		repassword(c)
	}
}

// 修改密码
func repassword(c *net.UDPConn) {
	fmt.Println("输入新密码")
	input := bufio.NewReader(os.Stdin)
	sel, _ := input.ReadString('\n')
	_, err := c.Write([]byte("/3" + sel))
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 修改用户名
func rena(c *net.UDPConn, name string) string {
	fmt.Println("输入新用户名")
	input := bufio.NewReader(os.Stdin)
	sel, _ := input.ReadString('\n')
	_, err := c.Write([]byte("/2" + sel))
	if err != nil {
		fmt.Println(err)
		return name
	}
	flag := string([]byte(<-change)[:1])
	if flag == "a" {
		fmt.Println("修改成功")
		return sel
	} else if flag == "b" {
		fmt.Println("修改失败，用户已存在")
		return name
	}
	return name
}

// 注册
func register(c *net.UDPConn) {
	var pass string
	var sel string
	for {
		fmt.Println("账号，密码不能由/，@等特殊符号开头")
		fmt.Println("输入注册的新账号")
		input := bufio.NewReader(os.Stdin)
		sel, _ = input.ReadString('\n')
		fmt.Println("输入注册的新密码")
		pass, _ = input.ReadString('\n')
		if []byte(sel)[0] == '/' || []byte(sel)[0] == '@' || []byte(pass)[0] == '/' || []byte(pass)[0] == '@' {
			fmt.Println("有非法输入")
		} else {
			break
		}
	}
	_, err := c.Write([]byte("/1" + sel + "//" + pass))
	if err != nil {
		fmt.Println(err)
		return
	}
	flag := string([]byte(<-change)[:1])
	if flag == "9" {
		fmt.Println("注册成功,请重新登录")
	} else if flag == "8" {
		fmt.Println("注册失败，用户已存在")
	}
}

// 监听来自服务端的消息
func print(c *net.UDPConn) {
	for {
		var buf [1024]byte
		n, adder, err := c.ReadFromUDP(buf[:])
		if err != nil {
			fmt.Println(err)
			return
		}
		if "/p" == string(buf[:2]) {
			change <- string(buf[2:])
			continue
		} else if "/f" == string(buf[:2]) {
			showfile(string(buf[2:]))
			continue
		}
		fmt.Printf("read from %v,mag:%s\n", adder, buf[:n])
	}
}

// 输出给客户端
func output(c *net.UDPConn, change string) string {
	input := bufio.NewReader(os.Stdin)
	for {
		s, _ := input.ReadString('\n')
		if s == "/e\r\n" || s == "@\r\n" {
			return s
		} else if s == "/file\r\n" {
			sandfile(c, string([]byte(change)[2:]))
			s = change + "我给你发送了一个文件，快去看看吧"
			_, _ = c.Write([]byte(s))
			continue
		}
		s = change + s
		_, err := c.Write([]byte(s))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if s == "/q\r\n" {
			fmt.Print("再见")
			//告知服务端下线
			os.Exit(0)
		}

	}
	return ""
}

// 告知服务端上线
func uplimit(name string, c *net.UDPConn) {
	input := "/s " + name
	s := input
	_, err := c.Write([]byte(s))
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 查询在线人数
func checkid(c *net.UDPConn) {
	input := "/c "
	s := input
	_, err := c.Write([]byte(s))
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 选择链接对象
func selec(c *net.UDPConn) {
	for {
		fmt.Println("按用户名选择链接对象\n按@a进行广播\n按/e返回")
		input := bufio.NewReader(os.Stdin)
		var change string
		sel, _ := input.ReadString('\n')
		if sel == "@a\r\n" {
			fmt.Println("切换到广播模式")
			change = "@a "
			change = output(c, change)
			if change == "/e\r\n" {
				return
			}
		} else if sel == "/e\r\n" {
			return
		} else if sel == "/q\r\n" {
			fmt.Print("再见")
			_, err := c.Write([]byte(sel))
			if err != nil {
				fmt.Println(err)
				return
			}
			os.Exit(0)
		} else if sel == "/check\r\n" {
			checkid(c)
		} else {
			if alone(c, sel) == "/e\r\n" {
				return
			}
		}
	}
}

// 私聊
func alone(c *net.UDPConn, sel string) string {
	number := ""
	s := "//" + sel
	_, err := c.Write([]byte(s))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	number = <-change
	if []byte(number)[0] == '0' {
		return ""
	}
	number = "/a" + string([]byte(number)[:5]) + " "
	number = output(c, number)
	return number
}

// 发送文件
func sandfile(c *net.UDPConn, sel string) {
	fmt.Println("输入文件地址")
	input := bufio.NewReader(os.Stdin)
	srcfile, _ := input.ReadString('\n')
	srcfile = string([]byte(srcfile)[:len(srcfile)-2])
	destfile := srcfile[strings.LastIndex(srcfile, "/")+1:]
	tmpfile := destfile + "tep.txt"
	file1, err := os.Open(srcfile)
	file3, err := os.OpenFile(tmpfile, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	defer file1.Close()
	file3.Seek(0, io.SeekStart)
	bs := make([]byte, 100, 100)
	n1, err := file3.Read(bs)
	countstr := string(bs[:n1])
	count, err := strconv.Atoi(countstr)
	count64 := int64(count)
	file1.Seek(count64, io.SeekStart)
	data := make([]byte, 500, 500)
	n2, n3 := -1, -1
	total := count
	for {
		time.Sleep(time.Microsecond * 100)
		n2, err = file1.Read(data)
		if err == io.EOF || n2 == 0 {
			file3.Close()
			os.Remove(tmpfile)
			break
		}

		s := "/F" + sel + "/" + srcfile + "%" + string(data[:n2])
		_, err := c.Write([]byte(s))
		if err != nil {
			fmt.Println(err)
		}
		//n3, err = file2.Write(data[:n2])
		total += n3

		file3.Seek(0, io.SeekStart)
		file3.WriteString(strconv.Itoa(total))
	}
}

// 接收文件
func showfile(flag string) {
	i := 0
	for ; ; i++ {
		if []byte(flag)[i] == 0 {
			break
		}
	}
	destfile := flag[:strings.Index(flag, "/")]
	file2, err := os.OpenFile(destfile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	defer file2.Close()
	file2.Seek(0, io.SeekEnd)
	file2.Write([]byte(flag[strings.Index(flag, "/")+1 : i]))
}

// 登录
func land(c *net.UDPConn) string {
	fmt.Println("输入账号(/e返回)")
	input := bufio.NewReader(os.Stdin)
	number, _ := input.ReadString('\n')
	if number == "/e\r\n" {
		return "/e"
	}
	fmt.Println("输入密码")
	pass, _ := input.ReadString('\n')
	_, err := c.Write([]byte("/p" + number + "//" + pass))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	flag := string([]byte(<-change)[:1])
	if flag == "3" {
		fmt.Println("登录成功")
		return string([]byte(number)[:len(number)-2])
	} else if flag == "2" {
		fmt.Println("已登陆，不能重复登录")
	} else if flag == "1" {
		fmt.Println("密码错误")
	} else if flag == "0" {
		fmt.Println("未注册")
	}
	return ""
}
