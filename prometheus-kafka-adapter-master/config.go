// Copyright 2018 Telefónica
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"strings"
	"text/template"

	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
)

var (
	kafkaBrokerList        = "kafka:9092"
	kafkaTopic             = "metrics"
	topicTemplate          *template.Template
	basicauth              = false
	basicauthUsername      = ""
	basicauthPassword      = ""
	kafkaCompression       = "none"
	kafkaBatchNumMessages  = "10000"
	kafkaSslClientCertFile = ""
	kafkaSslClientKeyFile  = ""
	kafkaSslClientKeyPass  = ""
	kafkaSslCACertFile     = ""
	kafkaSecurityProtocol  = ""
	kafkaSaslMechanism     = ""
	kafkaSaslUsername      = ""
	kafkaSaslPassword      = ""
	serializer             Serializer
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	if value := os.Getenv("LOG_LEVEL"); value != "" {
		log.Info("identifying LOG_LEVEL")
		logrus.SetLevel(parseLogLevel(value))
	}

	if value := os.Getenv("KAFKA_BROKER_LIST"); value != "" {
		log.Info("identifying KAFKA_BROKER_LIST")
		kafkaBrokerList = value
	}

	if value := os.Getenv("KAFKA_TOPIC"); value != "" {
		log.Info("identifying KAFKA_TOPIC")
		kafkaTopic = value
	}

	if value := os.Getenv("BASIC_AUTH_USERNAME"); value != "" {
		log.Info("identifying BASIC_AUTH_USERNAME")
		basicauth = true
		basicauthUsername = value
	}

	if value := os.Getenv("BASIC_AUTH_PASSWORD"); value != "" {
		log.Info("identifying BASIC_AUTH_PASSWORD")
		basicauthPassword = value
	}

	if value := os.Getenv("KAFKA_COMPRESSION"); value != "" {
		log.Info("identifying KAFKA_COMPRESSION")
		kafkaCompression = value
	}

	if value := os.Getenv("KAFKA_BATCH_NUM_MESSAGES"); value != "" {
		log.Info("identifying KAFKA_BATCH_NUM_MESSAGES")
		kafkaBatchNumMessages = value
	}

	if value := os.Getenv("KAFKA_SSL_CLIENT_CERT_FILE"); value != "" {
		log.Info("identifying KAFKA_SSL_CLIENT_CERT_FILE")
		kafkaSslClientCertFile = value
	}

	if value := os.Getenv("KAFKA_SSL_CLIENT_KEY_FILE"); value != "" {
		log.Info("identifying KAFKA_SSL_CLIENT_KEY_FILE")
		kafkaSslClientKeyFile = value
	}

	if value := os.Getenv("KAFKA_SSL_CLIENT_KEY_PASS"); value != "" {
		log.Info("identifying KAFKA_SSL_CLIENT_KEY_PASS")
		kafkaSslClientKeyPass = value
	}

	if value := os.Getenv("KAFKA_SSL_CA_CERT_FILE"); value != "" {
		log.Info("identifying KAFKA_SSL_CA_CERT_FILE")
		kafkaSslCACertFile = value
	}

	if value := os.Getenv("KAFKA_SECURITY_PROTOCOL"); value != "" {
		log.Info("identifying KAFKA_SECURITY_PROTOCOL")
		kafkaSecurityProtocol = strings.ToLower(value)
	}

	if value := os.Getenv("KAFKA_SASL_MECHANISM"); value != "" {
		log.Info("identifying KAFKA_SASL_MECHANISM")
		kafkaSaslMechanism = value
	}

	if value := os.Getenv("KAFKA_SASL_USERNAME"); value != "" {
		log.Info("identifying KAFKA_SASL_USERNAME")
		kafkaSaslUsername = value
	}

	if value := os.Getenv("KAFKA_SASL_PASSWORD"); value != "" {
		log.Info("identifying KAFKA_SASL_PASSWORD")
		kafkaSaslPassword = value
	}

	var err error
	serializer, err = parseSerializationFormat(os.Getenv("SERIALIZATION_FORMAT"))
	if err != nil {
		logrus.WithError(err).Fatalln("couldn't create a metrics serializer")
	}

	topicTemplate, err = parseTopicTemplate(kafkaTopic)
	if err != nil {
		logrus.WithError(err).Fatalln("couldn't parse the topic template")
	}
}

func parseLogLevel(value string) logrus.Level {
	level, err := logrus.ParseLevel(value)

	if err != nil {
		logrus.WithField("log-level-value", value).Warningln("invalid log level from env var, using info")
		return logrus.InfoLevel
	}

	return level
}

func parseSerializationFormat(value string) (Serializer, error) {
	switch value {
	case "json":
		return NewJSONSerializer()
	case "avro-json":
		return NewAvroJSONSerializer("schemas/metric.avsc")
	default:
		logrus.WithField("serialization-format-value", value).Warningln("invalid serialization format, using json")
		return NewJSONSerializer()
	}
}

func parseTopicTemplate(tpl string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"replace": func(old, new, src string) string {
			return strings.Replace(src, old, new, -1)
		},
		"substring": func(start, end int, s string) string {
			if start < 0 {
				start = 0
			}
			if end < 0 || end > len(s) {
				end = len(s)
			}
			if start >= end {
				panic("template function - substring: start is bigger (or equal) than end. That will produce an empty string.")
			}
			return s[start:end]
		},
	}
	return template.New("topic").Funcs(funcMap).Parse(tpl)
}
