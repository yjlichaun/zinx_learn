package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

var idmap map[int]string = make(map[int]string)
var lock *sync.RWMutex

func main() {
	lock = new(sync.RWMutex)
	//lock.RLock()
	//lock.RUnlock()
	//lock.Lock()
	//lock.Unlock()
	//创建服务端
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 30000,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer listen.Close()
	//监听并建立新协程
	for {
		var buf [1024]byte
		n, addr, err := listen.ReadFromUDP(buf[:])
		if err != nil {
			fmt.Println(err)
			return
		}
		go server(n, addr, buf, listen)
	}
}

// 服务端业务处理
func server(n int, addr *net.UDPAddr, buf [1024]byte, listen *net.UDPConn) {
	fmt.Println("收到", addr, "的数据：", string(buf[:]))
	//_, err := listen.WriteToUDP(buf[:n], addr)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//上线广播
	if "/s" == string(buf[:2]) {
		lock.Lock()
		idmap[addr.Port] = string(buf[3:])
		len := len(idmap)
		prin := "    [系统消息]" + string([]byte(idmap[addr.Port])[:n]) + "已上线:当前人数：" + strconv.Itoa(len)
		lock.Unlock()
		//广播
		broadcast(listen, prin)
		//下线广播
	} else if "/q" == string(buf[:2]) {
		len := len(idmap)
		prin := "    [系统消息]" + string([]byte(idmap[addr.Port])[:n]) + "已下线:当前人数：" + strconv.Itoa(len-1)
		//广播
		below(idmap[addr.Port])
		broadcast(listen, prin)
		delete(idmap, addr.Port)
		//查看在线id
	} else if "/c" == string(buf[:2]) {
		checkids(addr, listen)
		//用户广播
	} else if "@a" == string(buf[:2]) {
		lock.RLock()
		prin := "[来自用户:" + string([]byte(idmap[addr.Port])[:100]) + " 的广播]:" + string(buf[3:n-1])
		lock.RUnlock()
		broadcast(listen, prin)
		//用户私聊
	} else if "//" == string(buf[:2]) {
		alones(listen, buf, addr, n)
	} else if "/a" == string(buf[:2]) {
		alonebroad(listen, buf, addr, n)
	} else if "/p" == string(buf[:2]) {
		land(listen, buf, addr, n)
	} else if "/1" == string(buf[:2]) {
		register(listen, buf, addr, n)
	} else if "/2" == string(buf[:2]) {
		rename(listen, buf, addr, n)
	} else if "/3" == string(buf[:2]) {
		repass(listen, buf, addr, n)
	} else if "/F" == string(buf[:2]) {
		file(listen, buf[2:], n)
	}
}

// 文件中转
func file(listen *net.UDPConn, buf []byte, n int) {
	port, _ := strconv.Atoi(string(buf[:5]))
	adder := &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: port,
	}
	filename := buf[6:strings.Index(string(buf), "%")]
	filename = filename[strings.LastIndex(string(filename), "/")+1:]
	filename = filename[:len(filename)]
	i := 0
	for ; ; i++ {
		if buf[i] == 0 {
			break
		}
	}
	s := "/f" + string(filename) + "/" + string(buf[strings.Index(string(buf), "%")+1:i])
	_, err := listen.WriteToUDP([]byte(s), adder)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 修改密码
func repass(listen *net.UDPConn, buf [1024]byte, addr *net.UDPAddr, n int) {
	db, _ := sql.Open("mysql", "root:12345678qaz@tcp(127.0.0.1:3306)/test")
	err := db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	var (
		name     string
		t        string
		password string
		tatu     int
	)
	password = string(buf[2 : n-2])
	sqlstr := "select name, password ,status from id where name=?"
	lock.RLock()
	err = db.QueryRow(sqlstr, idmap[addr.Port]).Scan(&name, &t, &tatu)

	sqlstr = "delete from id where name=?"
	_, err = db.Exec(sqlstr, idmap[addr.Port])
	lock.RUnlock()
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlstr = "insert into id(name,password,status)values(?,?,?)"
	_, err = db.Exec(sqlstr, name, password, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = listen.WriteToUDP([]byte("修改成功"), addr)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 修改名字
func rename(listen *net.UDPConn, buf [1024]byte, addr *net.UDPAddr, n int) {
	db, _ := sql.Open("mysql", "root:12345678qaz@tcp(127.0.0.1:3306)/test")
	err := db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	var (
		name     string
		t        string
		password string
		tatu     int
	)
	name = string(buf[2 : n-2])
	sqlstr := "select name, password ,status from id where name=?"
	err = db.QueryRow(sqlstr, name).Scan(&name, &password, &tatu)
	if err != nil {
		sqlstr := "select name, password ,status from id where name=?"
		lock.RLock()
		err = db.QueryRow(sqlstr, idmap[addr.Port]).Scan(&t, &password, &tatu)

		sqlstr = "delete from id where name=?"
		_, err = db.Exec(sqlstr, idmap[addr.Port])
		if err != nil {
			fmt.Println(err)
			return
		}
		lock.RUnlock()
		sqlstr = "insert into id(name,password,status)values(?,?,?)"
		_, err = db.Exec(sqlstr, name, password, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		lock.Lock()
		idmap[addr.Port] = name
		lock.Unlock()
		_, err = listen.WriteToUDP([]byte("/pa"), addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		_, err = listen.WriteToUDP([]byte("/pb"), addr)
	}

}

// 注册
func register(listen *net.UDPConn, buf [1024]byte, addr *net.UDPAddr, n int) {
	db, _ := sql.Open("mysql", "root:12345678qaz@tcp(127.0.0.1:3306)/test")
	err := db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	var i int
	for i = 2; ; i++ {
		if buf[i] == '\r' {
			break
		}
	}
	name := string(buf[2:i])
	password := string(buf[i+4 : n-2])
	sqlstr := "select name, password ,status from id where name=?"
	err = db.QueryRow(sqlstr, name).Scan(&name, &password, &n)
	if err != nil {
		sqlstr = "insert into id(name,password,status)values(?,?,?)"
		_, err = db.Exec(sqlstr, name, password, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("注册成功")
		_, err := listen.WriteToUDP([]byte("/p9"), addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("注册失败")
		_, err := listen.WriteToUDP([]byte("/p8"), addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

// 广播
func broadcast(listen *net.UDPConn, prin string) {
	lock.RLock()
	for a, _ := range idmap {
		adder := &net.UDPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: a,
		}
		_, err := listen.WriteToUDP([]byte(prin), adder)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	lock.RUnlock()
}

// 查询在线id
func checkids(addr *net.UDPAddr, listen *net.UDPConn) {
	lock.RLock()
	for _, name := range idmap {
		name = "当前在线用户" + name
		_, err := listen.WriteToUDP([]byte(string([]byte(name)[:100])), addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	lock.RUnlock()
}

// 私聊链接
func alones(listen *net.UDPConn, buf [1024]byte, addr *net.UDPAddr, n int) {
	lock.RLock()
	for a, v := range idmap {
		if com(v, string(buf[2:n])) && a != addr.Port {
			_, err := listen.WriteToUDP([]byte("链接成功"), addr)
			_, err = listen.WriteToUDP([]byte("/p"+strconv.Itoa(a)), addr)
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}
	}
	lock.RUnlock()
	_, err := listen.WriteToUDP([]byte("链接失败"), addr)
	_, err = listen.WriteToUDP([]byte("/p"+"0"), addr)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 字符串比较
func com(a string, b string) bool {
	for i := 0; ; i++ {
		if a[i] != b[i] {
			if b[i] == '\r' && a[i] == 0 {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

// 私聊发送
func alonebroad(listen *net.UDPConn, buf [1024]byte, addr *net.UDPAddr, n int) {
	change, _ := strconv.Atoi(string(buf[2:7]))
	lock.RLock()
	_, ok := idmap[change]
	lock.RUnlock()
	if !ok {
		_, err := listen.WriteToUDP([]byte("链接中断,请按“@”重新选择链接对象，或者按“/e”返回"), addr)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	s := &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: change,
	}
	lock.RLock()
	prin := fmt.Sprint("[来自" + string([]byte(idmap[addr.Port])[:17]) + "]:" + string(buf[7:n-1]))
	princome := fmt.Sprint("to->" + string([]byte(idmap[change])[:17]) + ":" + string(buf[7:n-1]))
	lock.RUnlock()
	_, err := listen.WriteToUDP([]byte(prin), s)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = listen.WriteToUDP([]byte(princome), addr)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 登录
func land(listen *net.UDPConn, buf [1024]byte, addr *net.UDPAddr, n int) {
	db, _ := sql.Open("mysql", "root:12345678qaz@tcp(127.0.0.1:3306)/test")
	err := db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	var i int
	for i = 2; ; i++ {
		if buf[i] == '\r' {
			break
		}
	}
	name := string(buf[2:i])
	password := string(buf[i+4 : n-2])
	var (
		name1     string
		password1 string
		status    int
	)
	sqlstr := "select name, password ,status from id where name=?"
	err = db.QueryRow(sqlstr, name).Scan(&name1, &password1, &status)
	if err != nil {
		fmt.Println("未注册")
		_, err := listen.WriteToUDP([]byte("/p0"), addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if password != password1 {
		fmt.Println("密码错误")
		_, err := listen.WriteToUDP([]byte("/p1"), addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if password == password1 && status == 1 {
		fmt.Println("已登陆，不能重复登录")
		_, err := listen.WriteToUDP([]byte("/p2"), addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("登录成功")
		_, err := listen.WriteToUDP([]byte("/p3"), addr)
		if err != nil {
			fmt.Println(err)
			return
		}
		sqlstr = "delete from id where name=?"
		_, err = db.Exec(sqlstr, name)
		if err != nil {
			fmt.Println(err)
			return
		}
		sqlstr = "insert into id(name,password,status)values(?,?,?)"
		_, err = db.Exec(sqlstr, name, password, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// 下线
func below(name string) {
	var (
		password string
		status   int
	)
	db, _ := sql.Open("mysql", "root:12345678qaz@tcp(127.0.0.1:3306)/test")
	err := db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	sqlstr := "select name, password ,status from id where name=?"
	err = db.QueryRow(sqlstr, name).Scan(&name, &password, &status)
	if err != nil {
		log.Fatal("error{}", err)
		return
	}
	sqlstr = "delete from id where name=?"
	_, err = db.Exec(sqlstr, name)
	if err != nil {
		fmt.Println(err)
		return
	}
	sqlstr = "insert into id(name,password,status)values(?,?,?)"
	_, err = db.Exec(sqlstr, name, password, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
}
