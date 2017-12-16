package main

import (
	"context"
	"time"
	//"encoding/json"
	"io"
	"strings"

	"cloud.google.com/go/datastore"
	log "github.com/Sirupsen/logrus"

	"github.com/jedynykaban/testkeyholder/model"
	"github.com/jedynykaban/testkeyholder/services"

	"github.com/jinzhu/now"
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

	zulu := "2017-12-11T10:25:49Z"
	x, _ := now.Parse(zulu)
	log.Infof("Zulu: %v, parsed: %v", zulu, x)

	y, _ := time.Parse(time.RFC3339, zulu)
	log.Infof("Zulu: %v, parsed: %v", zulu, y)

	log.Info("application completed")
}

func mainZ() {
	log.Info("application started")

	zulu := "Mon, 11 Dec 2017 10:25:49 +0000"
	x, _ := now.Parse(zulu)
	log.Infof("Zulu: %v, parsed: %v", zulu, x)

	y, _ := time.Parse(time.RFC1123Z, zulu)
	log.Infof("Zulu: %v, parsed: %v", zulu, y)

	log.Info("application completed")
}

func mainX() {
	log.Info("application started")

	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "mosaiqio-dev")

	dbMitem := &model.DatabaseMitem{}
	//query := datastore.NewQuery("Mitem")
	id := "EhAKBU1pdGVtEICAgPb6mI4J"
	key, err := datastore.DecodeKey(id)
	_ = client.Get(ctx, key, dbMitem)

	// for idx := range pubs {
	// 	log.Infof("Publisher: %v (key: %v, decoded: %v)\n", pubs[idx], keys[idx], keys[idx].Encode())
	// }

	eKey := "EhAKBU1pdGVtEICAgPb2z4IJ"
	dKey := datastore.Key{}

	if len(eKey) > 0 {
		key, _ := datastore.DecodeKey(eKey)
		log.Infof("Encoded key: %v;  Decoded key: %v)\n", eKey, key)
	} else {
		log.Infof("Decoded key: %v;  Encoded key: %v)\n", dKey, dKey.Encode())
	}

	jsonRaw := dbMitem.Data

	if err != nil {
		log.Errorf("Will not work...: %v", err)
	}
	kojo := services.NewKojo()
	cd, err := kojo.GetCreationDate(jsonRaw)
	if err != nil {
		log.Error(err)
	}
	log.Infof("mitem creation date: %v\n", cd)

	log.Info("application completed")
}
