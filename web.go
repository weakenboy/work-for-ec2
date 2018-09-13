package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kylelemons/go-gypsy/yaml"
	"log"
	"net/http"
	"html/template"
	"strconv"
)

var (
	svc *ec2.EC2
	owner string
)

func init(){
	var disablessl *bool
	var True = true
	config, err := yaml.ReadFile("conf.yaml")
	if err != nil {
		fmt.Println(err)

	}
	AccessKeyID,_ := config.Get("AccessKeyID")
	SecretAccessKey,_ := config.Get("SecretAccessKey")
	owner,_ = config.Get("name")
	url,_ := config.Get("url")
	cred := credentials.NewStaticCredentialsFromCreds(credentials.Value{AccessKeyID: AccessKeyID, SecretAccessKey: SecretAccessKey})
	disablessl = &True
	sess0 := session.New(&aws.Config{
		Endpoint:    aws.String(url),
		Region:      aws.String(url),
		Credentials: cred,
		DisableSSL:  disablessl,
	})
	svc = ec2.New(sess0)

}
func sayhello(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t,_:=template.ParseFiles("index.html")
		data := struct {
			Title string
		}{
			Title: "index",
		}
		t.Execute(w, data)
	}
	//r.ParseForm()  //解析参数，默认是不会解析的
	//fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	//fmt.Println("path", r.URL.Path)
	//fmt.Println("scheme", r.URL.Scheme)
	//fmt.Println(r.Form.Get("name"))
	//m := map[string]string{
	//	"message":"Hello 我的朋友!",
	//}
	//data,_ :=json.Marshal(m)
	//fmt.Fprintf(w, string(data)) //这个写入到w的是输出到客户端的
}

func listInstances(w http.ResponseWriter, r *http.Request) {//列出虚拟机实例列表
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	var insid string = r.Form.Get("insid")
	if(insid != ""){

		output, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
			InstanceIds:[]*string{&insid},
		})
		if err != nil {
			fmt.Println(err)
		}
		msg,_:=json.Marshal(output)
		fmt.Fprintf(w,string(msg))
	}else{
		name:="owner-id"
		value:=owner
		output, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
			Filters:[]*ec2.Filter{
				&ec2.Filter{
					Name:&name,
					Values:[]*string{&value},
				},
			},
		})
		if err != nil {
			fmt.Println(err)
		}
		msg,_:=json.Marshal(output)
		fmt.Fprintf(w,string(msg))
	}
}
func runInstances(w http.ResponseWriter, r *http.Request){//创建虚拟机实例
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息

	//Reserved.
	device := "/dev/vda"
	stgid:=r.Form.Get("stg")
	stgsize,_:=strconv.Atoi(r.Form.Get("stg_size"))
	stgsizeint:=int64(stgsize)
	var blockDeviceMappings []*ec2.BlockDeviceMapping=[]*ec2.BlockDeviceMapping{
		&ec2.BlockDeviceMapping{
			DeviceName:&device,
			Ebs:&ec2.EbsBlockDevice{
				VolumeSize:&stgsizeint,
			},
			StorageId:&stgid,
		},
	}//磁盘
	var imageId string=r.Form.Get("imgid")//镜像id
	var instanceType string=r.Form.Get("typeid")//实例类型
	var keyName string="default"//关键名
	num,_ :=strconv.Atoi(r.Form.Get("num"))//最大实例创建数量
	maxCount:=int64(num)
	minCount:=maxCount//最小实例创建数量
	//var ramdiskId string =r.Form.Get("")//内存磁盘id
	securityGroupIds:=make([]*string,1,1)//安全组id
	sgid:=r.Form.Get("sgid")//安全组id
	errx:=append(securityGroupIds,&sgid)
	if errx != nil {
		fmt.Println(errx)
	}
	var subnetId string = r.Form.Get("subnetid")//子网id

	output,err :=svc.RunInstances(&ec2.RunInstancesInput{
		BlockDeviceMappings:blockDeviceMappings,
		ImageId:&imageId,
		InstanceType:&instanceType,
		KeyName:&keyName,
		MaxCount:&maxCount,
		MinCount:&minCount,
		//RamdiskId:&ramdiskId,
		SecurityGroupIds:securityGroupIds,
		SubnetId:&subnetId,

	})
	if err != nil {
		fmt.Println(err)
	}
	msg,_:=json.Marshal(output)
	fmt.Fprintf(w,string(msg))
	//fmt.Fprintf(w,"{'states':'success','message':''}")
}
func stopInstances(w http.ResponseWriter, r *http.Request){//停止虚拟机实例
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	var force bool=true
	insids := make([]*string,1,1)
	insid:=r.Form.Get("insid")
	errx:=append(insids,&insid)
	if errx != nil {
		fmt.Println(errx)
	}
	output,err :=svc.StopInstances(&ec2.StopInstancesInput{
		Force:&force,
		InstanceIds:insids,
	})
	if err != nil {
		fmt.Println(err)
	}

	msg,_:=json.Marshal(output)
	fmt.Fprintf(w,string(msg))
}
func distoryInstances(w http.ResponseWriter, r *http.Request){//销毁虚拟机实例
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息

	//insids := make([]*string,1,1)
	insid:=r.Form.Get("insid")
	fmt.Println(insid)
	//errx:=append(insids,&insid)
	//if errx != nil {
	//	fmt.Println(errx)
	//}
	output,err := svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{&insid},
	})
	if err != nil {
		fmt.Println(err)
	}
	msg,_:=json.Marshal(output)
	fmt.Fprintf(w,string(msg))
}
func createVolume(w http.ResponseWriter, r *http.Request){//创建存储卷
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	var AvailabilityZone string
	var	DryRun bool
	var	Encrypted bool
	var	Iops int64
	var	KmsKeyId string
	var	Size int64
	var	SnapshotId string
	var	VolumeType string
	output,err:=svc.CreateVolume(&ec2.CreateVolumeInput{
		AvailabilityZone:&AvailabilityZone,
		DryRun:&DryRun,
		Encrypted:&Encrypted,
		Iops:&Iops,
		KmsKeyId:&KmsKeyId,
		Size:&Size,
		SnapshotId:&SnapshotId,
		VolumeType:&VolumeType,

	})
	if err != nil {
		fmt.Println(err)
	}
	msg,_:=json.Marshal(output)
	fmt.Fprintf(w,string(msg))
}
func attachVolume(w http.ResponseWriter, r *http.Request){//实例加载卷
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	var device string
	var dryrun bool
	var insid string
	var volumeid string
	output,err:=svc.AttachVolume(&ec2.AttachVolumeInput{
		Device : &device,
		DryRun : &dryrun,
		InstanceId : &insid,
		VolumeId : &volumeid,
	})
	if err != nil {
		fmt.Println(err)
	}
	msg,_:=json.Marshal(output)
	fmt.Fprintf(w,string(msg))
}
func listVolume(w http.ResponseWriter, r *http.Request){//实例卸载卷
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	output, err := svc.DescribeVolumes(&ec2.DescribeVolumesInput{})
	if err != nil {
		fmt.Println(err)
	}
	msg,_:=json.Marshal(output)
	fmt.Fprintf(w,string(msg))
}
func detachVolume(w http.ResponseWriter, r *http.Request){//实例卸载卷
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	var device string
	var dryrun bool
	var insid string
	var volumeid string
	var force bool
	output,err:=svc.DetachVolume(&ec2.DetachVolumeInput{
		Device : &device,
		DryRun : &dryrun,
		Force : &force,
		InstanceId : &insid,
		VolumeId : &volumeid,
	})
	if err != nil {
		fmt.Println(err)
	}
	msg,_:=json.Marshal(output)
	fmt.Fprintf(w,string(msg))
}
func attachNetworkInterfaces(w http.ResponseWriter, r *http.Request){//实例加载卷
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	var deviceIndex int64
	var dryRun bool
	var instanceId string
	var networkInterfaceId string
	output,err:=svc.AttachNetworkInterface(&ec2.AttachNetworkInterfaceInput{
		DeviceIndex:&deviceIndex,
		DryRun:&dryRun,
		InstanceId:&instanceId,
		NetworkInterfaceId:&networkInterfaceId,
	})
	if err != nil {
		fmt.Println(err)
	}
	msg,_:=json.Marshal(output)
	fmt.Fprintf(w,string(msg))
}
func detachNetworkInterfaces(w http.ResponseWriter, r *http.Request){//实例卸载卷
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	var attachmentId string
	var dryRun bool
	var force bool
	output,err:=svc.DetachNetworkInterface(&ec2.DetachNetworkInterfaceInput{
		AttachmentId:&attachmentId,
		DryRun:&dryRun,
		Force:&force,
	})
	if err != nil {
		fmt.Println(err)
	}
	msg,_:=json.Marshal(output)
	fmt.Fprintf(w,string(msg))
}
func listNetworkInterfaces(w http.ResponseWriter, r *http.Request){//实例卸载卷
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	output, err := svc.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{})
	if err != nil {
		fmt.Println(err)
	}
	msg,_:=json.Marshal(output)
	fmt.Fprintf(w,string(msg))
}
func listImages(w http.ResponseWriter, r *http.Request){//列出镜像列表
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	output, err := svc.DescribeImages(&ec2.DescribeImagesInput{})
	if err != nil {
		fmt.Println(err)
	}
	msg,_:=json.Marshal(output)
	fmt.Fprintf(w,string(msg))
}
func listSecurityGroups(w http.ResponseWriter, r *http.Request){//列出安全组
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	name:="owner-id"
	values:=make( []*string,1,1)
	val:=append(values,&owner)

	filter := ec2.Filter{
		Name:&name,
		Values:val,
	}
	fs:=make([]*ec2.Filter,1,1)
	fs0:=append(fs,&filter)
	output, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: fs0,
	})
	if err != nil {
		fmt.Println(err)
	}
	msg,_:=json.Marshal(output)
	//fmt.Println(output)
	fmt.Fprintf(w,string(msg))
}
func listVpc(w http.ResponseWriter, r *http.Request){//列出网络列表
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	output, err := svc.DescribeVpcs(&ec2.DescribeVpcsInput{})
	if err != nil {
		fmt.Println(err)
	}
	msg,_:=json.Marshal(output)
	fmt.Fprintf(w,string(msg))
}
func listSubnet(w http.ResponseWriter, r *http.Request){//列出子网
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	name:="vpc-id"
	value:=r.Form.Get("VpcId")
	values:=make( []*string,1,1)
	val:=append(values,&value)

	filter := ec2.Filter{
		Name:&name,
		Values:val,
	}
	fs:=make([]*ec2.Filter,1,1)
	fs0:=append(fs,&filter)
	output, err := svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters:fs0,
	})
	if err != nil {
		fmt.Println(err)
	}

	msg,_:=json.Marshal(output)
	fmt.Fprintf(w,string(msg))
}
func main() {
	http.HandleFunc("/Instances/list", listInstances) //设置访问的路由
	http.HandleFunc("/Instances/run", runInstances) //设置访问的路由
	http.HandleFunc("/Instances/stop", stopInstances) //设置访问的路由
	http.HandleFunc("/Instances/distory", distoryInstances) //设置访问的路由
	http.HandleFunc("/Volume/create", createVolume) //设置访问的路由
	http.HandleFunc("/Volume/attach", attachVolume) //设置访问的路由
	http.HandleFunc("/Volume/list", listVolume) //设置访问的路由
	http.HandleFunc("/Volume/detach", detachVolume) //设置访问的路由
	http.HandleFunc("/NetworkInterfaces/attach", attachNetworkInterfaces) //设置访问的路由
	http.HandleFunc("/NetworkInterfaces/detach", detachNetworkInterfaces) //设置访问的路由
	http.HandleFunc("/NetworkInterfaces/list", listNetworkInterfaces) //设置访问的路由
	http.HandleFunc("/Images/list", listImages) //设置访问的路由
	http.HandleFunc("/SecurityGroups/list", listSecurityGroups) //设置访问的路由
	http.HandleFunc("/Vpc/list", listVpc) //设置访问的路由
	http.HandleFunc("/Subnet/list", listSubnet) //设置访问的路由
	//http.HandleFunc("/hello", sayhello) //设置访问的路由
	h := http.FileServer(http.Dir("html"))
	http.DetectContentType([]byte("application/json/;charset=utf-8"))
	http.Handle("/html/", http.StripPrefix("/html/", h)) // 启动静态文件服务
	err := http.ListenAndServe(":9091", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}