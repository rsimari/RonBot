package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "bytes"
    "time"
)

//************** IRON.IO STUFF *******************//

type Task struct {
        CodeName string `json:"code_name"`
        Payload string `json:"payload"`
        StartAt string `json:"start_at"`
}

type ReqData struct {
        Schedules []*Task `json:"schedules"`
}

func queueTwilio(execTime string, msg string) string {

        //IRON.IO Credentials
        const token = "IpsW3ZNlFsLl42T2vSqw"
        const project = "58cf2ffd0e2c7300061ad812"


        if len(execTime) == 8 {
          //add todays calendar date etc.
          loc, _ := time.LoadLocation("America/New_York") //Eventually get settings from file 
          
          date := time.Now().In(loc)
          dateString := date.Format("2006-01-02")   
          execTime = string(dateString) + "T" + execTime + "-04:00" //Default to EST
          //newTime, _ := time.Parse("2006-01-02T15:04:05-07:00", execTime)
          //newUTCTime := newTime.UTC()
  
        } else {

          execTime := execTime + "-04:00" //Default to EST

          newTime, _ := time.Parse("2006-01-02T15:04:05Z-07:00", execTime)
          newUTCTime := newTime.UTC()

          fmt.Println(newUTCTime.Format("2006-01-02T15:04:05Z"))	

        }


        // Insert our project ID and token into the API endpoint
        target := fmt.Sprintf("http://worker-us-east.iron.io/2/projects/%s/schedules?oauth=%s", project, token)

        // Build the payload
        payload := map[string]interface{} {
                "Message" : msg,
                "Phone" : "+14127601315", //should be phone from user defaults 
        }

        payload_bytes, err := json.Marshal(payload) //JSON Encode payload

        if err != nil {
                panic(err.Error())
                return "Oops. I was unable to set your reminder."
        }

        payload_str := string(payload_bytes) //Convert to string 

        // Build the task to be executed 
        task := &Task {
                CodeName: "send_twilio",
                Payload: payload_str,
                StartAt: execTime,
        }

        // Build a request containing the task
        json_data := &ReqData {
                Schedules: []*Task { task },
        }

        json_bytes, err := json.Marshal(json_data)

        if err != nil {
                panic(err.Error())
                return "Oops. I was unable to set your reminder."
        }

        json_str := string(json_bytes)

        // Post expects a Reader
        json_buf := bytes.NewBufferString(json_str)

        // Make the request
        resp, err := http.Post(target, "application/json", json_buf)

        if err != nil {
                panic(err.Error())
                return "Oops. I was unable to set your reminder."
        }

        defer resp.Body.Close()

        // Read the response
        resp_body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                panic(err.Error())
                return "Oops. I was unable to set your reminder."
        }

        // Print the response to STDOUT
        fmt.Println(string(resp_body))

        return "Okay. I will send you a text message to remind you!"
}

//************** END IRON.IO *******************//