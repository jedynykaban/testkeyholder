package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/fake-gcs-server/fakestorage"
	slugify "github.com/metal3d/go-slugify"

	"bytes"
	"encoding/base64"
	"encoding/gob"
)

type x struct {
	ID    string
	Name  string
	Value int
}

var config Config

func init() {
	config = getConfig()
	setupLogging(config.Service.LogOutput, config.Service.LogLevel, config.Service.LogFormat)
}

func setupLogging(output io.Writer, level log.Level, format string) {
	log.SetOutput(output)
	log.SetLevel(level)
	if strings.EqualFold(format, "json") {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

func main() {
	log.Info("application started")

	s := "ALA MA KOTkA"
	log.Info(strings.Title(s))
	s = strings.ToLower(s)
	log.Info(strings.Title(s))

	log.Info("t: ", time.Duration(rand.Intn(int(time.Second))))

	toSlug := "A > B"
	log.Info("original: ", toSlug)
	slugged := slugify.Marshal(toSlug, true)
	log.Info("slug(1): ", slugged)
	log.Info("slug(2): ", slugify.Marshal(slugged))

	serverSimpleTest()
	serverMoreComplexTest()

	log.Info("application completed")
}

func serverSimpleTest() {
	server := fakestorage.NewServer([]fakestorage.Object{
		{
			BucketName: "some-bucket",
			Name:       "some/object/file.txt",
			Content:    []byte("inside the file"),
		},
	})
	defer server.Stop()
	client := server.Client()
	object := client.Bucket("some-bucket").Object("some/object/file.txt")
	reader, err := object.NewReader(context.Background())
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", data)
}

func serverMoreComplexTest() {
	xb := x{ID: "id", Name: "test", Value: 13}
	fmt.Printf("%v\n", xb)
	encoded := toGOB(xb)

	server := fakestorage.NewServer([]fakestorage.Object{})
	defer server.Stop()

	obj := fakestorage.Object{
		BucketName: "another-bucket",
		Name:       xb.ID,
		Content:    encoded,
	}
	server.CreateObject(obj)

	fmt.Println("-------------------")

	client := server.Client()
	object := client.Bucket(obj.BucketName).Object(obj.Name)
	reader, err := object.NewReader(context.Background())
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	xa := fromGOB(data)
	fmt.Printf("%v\n", xa)
}

// go binary encoder
func toGOB(m x) []byte {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	return b.Bytes()
}

func toGOB64(m x) string {
	return base64.StdEncoding.EncodeToString(toGOB(m))
}

// go binary decoder
func fromGOB(by []byte) x {
	m := x{}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err := d.Decode(&m)
	if err != nil {
		fmt.Println(`failed gob Decode`, err)
	}
	return m
}

// go binary decoder
func fromGOB64(str string) x {
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println(`failed base64 Decode`, err)
	}
	return fromGOB(by)
}
