package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main()  {
	mux:=http.NewServeMux()
	mux.Handle("/",&myHandler{})
	//Handle和HandleFunc的区别，第一个参数无异，HandleFunc第二个参数是一个函数，Handle是一个类型

	server:=&http.Server{
		Addr:":4000",
		WriteTimeout:4*time.Second,
	}
	server.Handler=mux   //mux实际上是实现了handler接口的一个变量

	quit:=make(chan os.Signal)     //创建通道，存放系统信号
	signal.Notify(quit,os.Interrupt)   //向通道发送退出信息

	go func() {
		<-quit    //从通道里接收通知，接收到退出通知后，服务器主动退出，一旦quit取到值，后面代码开始执行
		if err:=server.Close();err!=nil{    //服务器关闭，关闭之后，会抛出err
			log.Fatal("Close server:",err)
		}
	}()

	//路由注册：注册一个函数用于响应某一个路由，简单来说，就是后面参数的匿名函数用于响应路由地址"/"
	//http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	//	writer.Write([]byte("hello,this is version 1!"))   //将字符串转换为字节
	//})
	//"/"路由包含了所有未注册的路由，所以只要有路由没有注册，它就会被归类到最近的"/"路由下
	
	//http.HandleFunc("/bye",sayBye)
	mux.HandleFunc("/bye",sayBye)

	//log.Println("Starting server...v1")   //打印日志
	//log.Println("Starting server...v2")
	log.Println("Starting server...v3")

	//log.Fatal(http.ListenAndServe(":4000",nil))
	//log.Fatal(http.ListenAndServe(":4000",mux))  //监听4000端口
	//log.Fatal(server.ListenAndServe())
	err:=server.ListenAndServe()    //监听端口后，比如我们在命令行输入ctr+c，停止服务器，
	                                //则signal.Notify会向通道传送Interrupt信息，被线程拦截到后，则服务器停止，为人为停止。
	if err!=nil{  //若关闭错误不为空，判断其类型
		if err==http.ErrServerClosed{   //响应关闭信号
			log.Print("Server closed under request!")
		}else {   //异常关闭错误
			log.Fatal("Server closed unexpected!")
		}
	}
	log.Println("Server exit!")

	//注：如果是正常关闭，则会打印
	//Server closed under request!   表示服务器是响应人为操作正常退出
	//Server exit!
}

type myHandler struct {}

//go对大小写敏感
func (*myHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request)  {
	//writer.Write([]byte("hello,this is version 2!"))
	writer.Write([]byte("hello,this is version 3!"+request.URL.String()))
}

func sayBye(writer http.ResponseWriter, request *http.Request)  {
	//writer.Write([]byte("Bye bye,this is version 1!"))
	//writer.Write([]byte("Bye bye,this is version 2!"))
	time.Sleep(3*time.Second)
	writer.Write([]byte("Bye bye,this is version 3!"))
}
