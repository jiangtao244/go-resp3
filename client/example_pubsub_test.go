/*
Copyright 2019 Stefan Miller

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client_test

import (
	"fmt"
	"log"

	"github.com/stfnmllr/go-resp3/client"
)

func MsgCallback(pattern, channel, msg string) {
	fmt.Printf("Channel: %s Message: %s", channel, msg)
}

func Example_pubsub() {
	// Open connection with standard host and port:
	// - if environment variablen REDIS_HOST AND REDIS_PORT are set,
	//   these values are used
	// - otherwise localhost (127.0.0.1) and Redis default port 6379 are used
	conn, err := client.Dial("")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Subscribe channel.
	if err := conn.Subscribe([]string{"mychannel"}, MsgCallback).Err(); err != nil {
		log.Fatal(err)
	}
	// Connection allows to proceed with any other Redis commands.
	// Not limited to Subscribe, Psubscribe, Unsubscribe, Punsubscribe, Ping and Quit.
	if err := conn.Set("mykey", "Hello Redis").Err(); err != nil {
		log.Fatal(err)
	}
	value, err := conn.Get("mykey").ToString()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(value)
	// Publish message
	if err := conn.Publish("mychannel", "Hello Redis").Err(); err != nil {
		log.Fatal(err)
	}
	// Unsubscribe channel.
	if err := conn.Unsubscribe([]string{"mychannel"}).Err(); err != nil {
		log.Fatal(err)
	}
	// Output:
	// Hello Redis
	// Channel: mychannel Message: Hello Redis
}
