package tstutils

import (
	"os/exec"
	"fmt"
)

func StartMySQL() error {

	cmdStr := "docker run -p 3306:3306 --name mysql_cabapi -e MYSQL_ROOT_PASSWORD=root -d mysql:5.7.22"

	out, err := exec.Command("/bin/sh", "-c", cmdStr).Output()

	if err != nil {
		return err
	}

	fmt.Printf("%s", out)

	return nil

}

func CreateDatabase() error {

	cmdStr := "mysql -hlocalhost -P3306 --protocol=tcp -uroot -proot -e \"create database cabapi\";"

	out, err := exec.Command("/bin/sh", "-c", cmdStr).Output()

	if err != nil {
		return err
	}

	fmt.Printf("%s", out)

	return nil
}

func StopMySQL() error {

	cmdStr := "docker stop mysql_cabapi"
	out, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		return err
	}
	fmt.Printf("%s", out)

	cmdStr = "docker rm mysql_cabapi"
	out, err = exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		return err
	}
	fmt.Printf("%s", out)

	return nil
}

func StartRedis() error {

	cmdStr := "docker run -p 6379:6379 --name redis_cabapi -d redis:4.0.10"

	out, err := exec.Command("/bin/sh", "-c", cmdStr).Output()

	if err != nil {
		return err
	}

	fmt.Printf("%s", out)

	return nil

}

func StopRedis() error {

	cmdStr := "docker stop redis_cabapi"
	out, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		return err
	}
	fmt.Printf("%s", out)

	cmdStr = "docker rm redis_cabapi"
	out, err = exec.Command("/bin/sh", "-c", cmdStr).Output()
	if err != nil {
		return err
	}
	fmt.Printf("%s", out)

	return nil
}
