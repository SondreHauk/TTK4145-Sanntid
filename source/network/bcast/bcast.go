package bcast

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"source/network/conn"
	"sync"
	"time"
)

// Increased to 4096. Initially 1024
const bufSize = 4096

// Encodes received values from `chans` into type-tagged JSON, then broadcasts
// it on `port`
func Transmitter(port int, chans ...interface{}) {
	checkArgs(chans...)
	typeNames := make([]string, len(chans))
	selectCases := make([]reflect.SelectCase, len(typeNames))
	
	for i, ch := range chans {
		selectCases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
		typeNames[i] = reflect.TypeOf(ch).Elem().String()
	}

	conn := conn.DialBroadcastUDP(port)
	netAddr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))
	localAddr,_ := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", port))

	for {
		chosen, value, _ := reflect.Select(selectCases)
		jsonstr, _ := json.Marshal(value.Interface())

		msgID := fmt.Sprintf("%d", time.Now().UnixNano()) // generate unique ID

		ttj, _ := json.Marshal(typeTaggedJSON{
			ID: msgID,
			TypeId: typeNames[chosen],
			JSON:   jsonstr,
		})
		if len(ttj) > 2000 {fmt.Printf("WV size: %d\n", len(ttj))}
		if len(ttj) > bufSize {
		    panic(fmt.Sprintf(
		        "Tried to send a message longer than the buffer size (length: %d, buffer size: %d)\n\t'%s'\n"+
		        "Either send smaller packets, or go to network/bcast/bcast.go and increase the buffer size",
		        len(ttj), bufSize, string(ttj)))
		}
		conn.WriteTo(ttj, netAddr)
		conn.WriteTo(ttj, localAddr)
    		
	}
}

//Functionality for deduplication of messages

var (
	messageCache = make(map[string]time.Time)
	cacheMutex = sync.Mutex{}
)

func cleanUpCache() {
	for {
		time.Sleep(5 * time.Second) // Clean cache every 5 seconds
		cacheMutex.Lock()
		for id, timestamp := range messageCache {
			if time.Since(timestamp) > 5 *time.Second {
				delete(messageCache, id) // Remove messages older than 5 sec
			}
		}
		cacheMutex.Unlock()
	}
}


// Matches type-tagged JSON received on `port` to element types of `chans`, then
// sends the decoded value on the corresponding channel
func Receiver(port int, chans ...interface{}) {
	checkArgs(chans...)
	chansMap := make(map[string]interface{})
	for _, ch := range chans {
		chansMap[reflect.TypeOf(ch).Elem().String()] = ch
	}

	var buf [bufSize]byte
	conn := conn.DialBroadcastUDP(port)

	go cleanUpCache()

	for {
		conn.SetReadDeadline(time.Now().Add(2 * time.Second)) // Prevent blocking
		n, _, e := conn.ReadFrom(buf[0:])
		if netErr, ok := e.(net.Error); ok && netErr.Timeout() {
			continue // Just a timeout, retry
		}

		if e != nil{
			fmt.Printf("bcast.Receiver(%d, ...):ReadFrom() failed: \"%+v\"\n", port, e)
			time.Sleep(time.Second)
			continue
		}

		var ttj typeTaggedJSON
		json.Unmarshal(buf[0:n], &ttj)

		cacheMutex.Lock()
        if _, exists := messageCache[ttj.ID]; exists {
            cacheMutex.Unlock()
            continue // Duplicate message, ignore it
        }
        messageCache[ttj.ID] = time.Now() // Store message ID
        cacheMutex.Unlock()

		ch, ok := chansMap[ttj.TypeId]
		if !ok {
			continue
		}
		v := reflect.New(reflect.TypeOf(ch).Elem())
		json.Unmarshal(ttj.JSON, v.Interface())
		reflect.Select([]reflect.SelectCase{{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(ch),
			Send: reflect.Indirect(v),
		}})
	}
}

type typeTaggedJSON struct {
	ID 	   string 	// Unique ID for deduplication
	TypeId string
	JSON   []byte
}

// Checks that args to Tx'er/Rx'er are valid:
//  All args must be channels
//  Element types of channels must be encodable with JSON
//  No element types are repeated
// Implementation note:
//  - Why there is no `isMarshalable()` function in encoding/json is a mystery,
//    so the tests on element type are hand-copied from `encoding/json/encode.go`
func checkArgs(chans ...interface{}) {
	n := 0
	for range chans {
		n++
	}
	elemTypes := make([]reflect.Type, n)

	for i, ch := range chans {
		// Must be a channel
		if reflect.ValueOf(ch).Kind() != reflect.Chan {
			panic(fmt.Sprintf(
				"Argument must be a channel, got '%s' instead (arg# %d)",
				reflect.TypeOf(ch).String(), i+1))
		}

		elemType := reflect.TypeOf(ch).Elem()

		// Element type must not be repeated
		for j, e := range elemTypes {
			if e == elemType {
				panic(fmt.Sprintf(
					"All channels must have mutually different element types, arg# %d and arg# %d both have element type '%s'",
					j+1, i+1, e.String()))
			}
		}
		elemTypes[i] = elemType

		// Element type must be encodable with JSON
		checkTypeRecursive(elemType, []int{i+1})

	}
}


func checkTypeRecursive(val reflect.Type, offsets []int){
	switch val.Kind() {
	case reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		panic(fmt.Sprintf(
			"Channel element type must be supported by JSON, got '%s' instead (nested arg# %v)",
			val.String(), offsets))
	case reflect.Map:
		if val.Key().Kind() != reflect.String {
			panic(fmt.Sprintf(
				"Channel element type must be supported by JSON, got '%s' instead (map keys must be 'string') (nested arg# %v)",
				val.String(), offsets))
		}
		checkTypeRecursive(val.Elem(), offsets)
	case reflect.Array, reflect.Ptr, reflect.Slice:
		checkTypeRecursive(val.Elem(), offsets)
	case reflect.Struct:
		for idx := 0; idx < val.NumField(); idx++ {
			checkTypeRecursive(val.Field(idx).Type, append(offsets, idx+1))
		}
	}
}
