Add this for graceful close our process 
--------------------------------------


- restart with syscall.SIGQUIT/SIGTERM/SIGINT signal

- internal
http server step
-> graceful shutdown(stop listen && wait all conns done) \
-> graceful shutdown, wait all coroutine done, \
-> main.go other before shutdown, \
-> server close

- tips:
AddOne():Done() = 1:1 \
-> set time out for better close coroutine \
-> insure no goroutine leak can use context.Background().

- example:
	```go
	 main.go 
        
        //when deregister server before: 
            defer func() {
    	  	    ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
    	  	    defer cancel()
    	  	    graceful.Shutdown(ctx)
    	  }()
    
        //use goroutine do something:
    		graceful.AddOne()
    		go func(){
    			defer graceful.Done()
    		....
    		}()
    
        //do something overtime process:
    		graceful.AddOne()
    		go func() {
    			defer graceful.Done()
    			for i := 0; i < 10; i++ {
    				if graceful.IsShutting() {
    					logger.Info("server is shutting...")
    					break
    				}
    			.... process
    			}
    		}()
```