package main

import (
	"fmt"
	"net/http"
)


type eventBase struct {
	eventMessage chan []byte
}

//eventBase implementing Handler interface to be placed in the call to 'ListenAndServe(port,Handler)'
func (eb *eventBase) ServeHTTP(rw http.ResponseWriter, r *http.Request){		
	
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	
	select {
    case msg := <-eb.eventMessage:
		rw.Header().Set("Content-Type", "text/event-stream")
		rw.Header().Set("Cache-Control", "no-cache")
		rw.Header().Set("Connection", "keep-alive")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
        
		for {
			fmt.Println("event found, transmitting..")
			//w.Write(msg)
			fmt.Fprintf(rw, "data: %s\n\n", msg)
			
			flusher.Flush()
		}
    default:
		nomsg := "nah"
		fmt.Fprintf(rw, "data: %s\n\n", nomsg)
        fmt.Println("no event recieved")
    }
}

func dummySend(e *eventBase) {
	var msg []byte
	a := "akash"
	msg = []byte(a)
	//TODO feed data in msg
	select{
		case e.eventMessage <- msg:
			fmt.Println("dummy message sent")
		default:
			fmt.Println("nothing to dummy send")
	}
	
}

func beginRealTimeProc() (eBase *eventBase){
	eBase = &eventBase{
		eventMessage: make(chan []byte),	
	}
	return 
	
}

func (e *eventBase) listen() {
	for {
		select {
			case <-e.eventMessage:
				fmt.Println("I heard something...")
				
			default:
				fmt.Println("nothing")
		}
	}
}

func main() {
	eventGuy := beginRealTimeProc()
	dummySend(eventGuy)
	http.ListenAndServe(":6713", eventGuy)
}