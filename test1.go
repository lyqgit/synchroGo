package main


/**
 * @api {get} /user/:id Request User information
 * @apiName GetUser
 * @apiGroup User
 *
 * @apiParam {Number} id Users unique ID.
 *
 * @apiExample
 * {"id":1, "name": "n1"}
 */

import (
	"encoding/json"
	"os"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"io/ioutil"
	"io"
	"github.com/fsnotify/fsnotify"
	"time"
)

var change string = "0"

var dir string

type cssConfig struct{
	File string `json:"file"`
	Url string `json:"url"`
}

type jsConfig struct{
	File string `json:"file"`
	Url string `json:"url"`
}

type config struct{
	Port string `json:"port"`
	Custom bool `json:"csutom"`
	Css cssConfig `json:"css"`
	Js jsConfig `json:"js"`
}

var configData config

func configInit(){//读取配置文件
	jsonStr,err := ioutil.ReadFile("./webtest/config.json")
	if err != nil{
		fmt.Println(err)
	}

	json.Unmarshal(jsonStr,&configData)

	fmt.Println(configData.Css.Url)

}

var watcher *fsnotify.Watcher

func main(){

	configInit()
	watchInit(&configData)
	// go watch(&configData)//执行监听文件

	dir, err := filepath.Abs(filepath.Dir("./"))
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(strings.Replace(dir, "\\", "/", -1))

	http.Handle("/css/",http.StripPrefix("/css/",http.FileServer(http.Dir("./webtest/css/"))))
	http.HandleFunc("/messageEvent",messageEvent)
	http.HandleFunc("/",serverH)

	fmt.Println(http.ListenAndServe(":"+configData.Port,nil))
}

func watchInit(conf *config){
	//创建一个监控对象
	watch, err := fsnotify.NewWatcher();
	if err != nil {
			fmt.Println(err);
	}
	// defer watch.Close();

	if conf.Custom{
		dir = conf.Css.File
	}else{
		dir = "./webtest"
	}

	
	//添加要监控的对象，文件或文件夹
	err = watch.Add(dir);

	if err != nil {
		fmt.Println(err);
	}

	file,err := os.Stat(dir)

	if file.IsDir(){
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			//这里判断是否为目录，只需监控目录即可
			//目录下的文件也在监控范围内，不需要我们一个一个加
			if info.IsDir() {
					path, err := filepath.Abs(path);
					if err != nil {
							return err;
					}
					err = watch.Add(path);
					if err != nil {
							return err;
					}
					fmt.Println("监控 : ", path);
			}
			return nil;
		});
	}

	

	if err != nil {
		fmt.Println(err);
	}
	watcher = watch
}

func watchOnce(){
	
	ev := <-watcher.Events

	timeStr:=time.Now().Format("2006-01-02 15:04:05")
	//判断事件发生的类型，如下5种
	// Create 创建
	// Write 写入
	// Remove 删除
	// Rename 重命名
	// Chmod 修改权限
	if ev.Op == fsnotify.Create {
		change = "1"
		fmt.Println(timeStr+" 创建文件 : ", ev.Name);
	}
	if ev.Op == fsnotify.Write {
		change = "1"
		fmt.Println(timeStr+" 写入文件 : ", ev.Name);
	}
	if ev.Op == fsnotify.Remove {
		change = "1"
		fmt.Println(timeStr+" 删除文件 : ", ev.Name);
	}
	if ev.Op == fsnotify.Rename {
		change = "1"
		fmt.Println(timeStr+" 重命名文件 : ", ev.Name);
	}
	if ev.Op == fsnotify.Chmod {
		change = "1"
		fmt.Println(timeStr+" 修改权限 : ", ev.Name);
	}
		
				
			
		
}

func watch(conf *config){
	//创建一个监控对象
	watch, err := fsnotify.NewWatcher();
	if err != nil {
			fmt.Println(err);
	}
	defer watch.Close();

	if conf.Custom{
		dir = conf.Css.File
	}else{
		dir = "./webtest"
	}

	
	//添加要监控的对象，文件或文件夹
	err = watch.Add(dir);

	if err != nil {
		fmt.Println(err);
	}

	file,err := os.Stat(dir)

	if file.IsDir(){
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			//这里判断是否为目录，只需监控目录即可
			//目录下的文件也在监控范围内，不需要我们一个一个加
			if info.IsDir() {
					path, err := filepath.Abs(path);
					if err != nil {
							return err;
					}
					err = watch.Add(path);
					if err != nil {
							return err;
					}
					fmt.Println("监控 : ", path);
			}
			return nil;
		});
	}

	

	if err != nil {
		fmt.Println(err);
	}
	//判断事件发生的类型，如下5种
	// Create 创建
	// Write 写入
	// Remove 删除
	// Rename 重命名
	// Chmod 修改权限

	
		for{
			select{
				case ev := <-watch.Events:{
					if ev.Op == fsnotify.Create {
						change = "1"
						fmt.Println("创建文件 : ", ev.Name);
					}
					if ev.Op == fsnotify.Write {
						change = "1"
						fmt.Println("写入文件 : ", ev.Name);
					}
					if ev.Op == fsnotify.Remove {
						change = "1"
						fmt.Println("删除文件 : ", ev.Name);
					}
					if ev.Op == fsnotify.Rename {
						change = "1"
						fmt.Println("重命名文件 : ", ev.Name);
					}
					if ev.Op == fsnotify.Chmod {
						change = "1"
						fmt.Println("修改权限 : ", ev.Name);
					}
					// if (ev.Op&fsnotify.Create == fsnotify.Create) || (ev.Op&fsnotify.Write == fsnotify.Write) || (ev.Op&fsnotify.Remove == fsnotify.Remove) || (ev.Op&fsnotify.Rename == fsnotify.Rename) || (ev.Op&fsnotify.Chmod == fsnotify.Chmod) {
					// 	change = "1"
					// }
				}
					
			}
		}
	
}

func messageEvent(w http.ResponseWriter,r *http.Request){
	change = "0"
	watchOnce()
	fmt.Println(change)
	w.Header().Set("Content-type","text/event-stream;charset=utf-8")
	w.Header().Set("Connection","keep-alive")
	w.Write([]byte("data:"+change+"\n\n"))
}

func serverH(w http.ResponseWriter,r *http.Request){
	by,_ := ioutil.ReadFile("."+r.URL.String())
	
	w.Header().Set("Content-type","text/html;charset=utf-8")
	w.Header().Set("Cache-control","no-cache")
	w.Write(by)
	io.WriteString(w,`<script>
		var es = new EventSource('/messageEvent');
		es.onmessage=function(e){console.log(e.data);
		if(e.data == 1){
			window.location.reload();
		}};
		var link=document.getElementsByTagName('link');
		if(link.length != 0){
			for(var i=0;i<link.length;i++){
				link[i].href+='?_='+Math.floor(Math.random()*10000);
			}
		};
		var script=document.getElementsByTagName('script');
		if(script.length != 0){
			for(var i=0;i<script.length;i++){
				if(script[i].src != ''){
					script[i].src+='?_='+Math.floor(Math.random()*10000);
				}		
			}
		}
		</script>`)
}