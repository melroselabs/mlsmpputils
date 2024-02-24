// Simple SMS API bridge code
//
// Alternative to using SMS API Bridge service from Melrose Labs (melroselabs.com)
//
// Code enables applications supporting LINK Mobility SMS API to switch to using Esendex Conversations API without changing existing code.
// Intended as example code that can be used as a starting point for more capable bridging between APIs.
//
// 24 Feb 2024 Mark Hay
//
// Melrose Labs
// melroselabs.com

package main

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
)

const targetURL = "https://conversations.esendex.com/v1/messages"
const targetAPIKey = "YOUR_API_KEY"
const targetAccountReference = "EX000000"

// Define structures for the Format A request and response
type FormatARequest struct {
    Source            string `json:"source"`
    Destination       string `json:"destination"`
    UserData          string `json:"userData"`
    PlatformId        string `json:"platformId"`
    PlatformPartnerId string `json:"platformPartnerId"`
    UseDeliveryReport bool   `json:"useDeliveryReport"`
}

type FormatAResponse struct {
    MessageId   string `json:"messageId"`
    ResultCode  int    `json:"resultCode"`
    Description string `json:"description"`
}

// Define structures for the Format B request and response
type FormatBRequest struct {
    Channel           string            `json:"channel"`
    Metadata          map[string]string `json:"metadata"`
    RecipientVariables map[string]string `json:"recipientVariables"`
    To                string            `json:"to"`
    Body              struct {
        Text struct {
            Value string `json:"value"`
        } `json:"text"`
    } `json:"body"`
    From          string `json:"from"`
    CharacterSet  string `json:"characterSet"`
    Expiry        int    `json:"expiry"`
}

type FormatBResponse struct {
    Id string `json:"id"`
}

func smsSendHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req FormatARequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Transform Format A request to Format B
    formatB := FormatBRequest{
        Channel: "sms",
        To: req.Destination,
        Body: struct {
            Text struct {
                Value string `json:"value"`
            } `json:"text"`
        }{
            Text: struct {
                Value string `json:"value"`
            }{
                Value: req.UserData,
            },
        },
        From: req.Source,
        CharacterSet: "GSM",
        Expiry: 30,
    }
    
    requestBody, err := json.Marshal(formatB)
    if err != nil {
        log.Print("Failed to marshal")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Send the Format B request to the external service
    client := &http.Client{}
    request, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(requestBody))
    if err != nil {
        log.Print("Failed to call target URL")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    request.Header.Set("Content-Type", "application/json")
    // Add authorization header if required
    request.Header.Set("Account-Reference", targetAccountReference)
    request.Header.Set("Api-Key", targetAPIKey)
    
    response, err := client.Do(request)
    if err != nil {
        log.Print("Failed to send request")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    var respB FormatBResponse
    if err := json.Unmarshal(body, &respB); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Transform Format B response to Format A response
    respA := FormatAResponse{
        MessageId:   respB.Id,
        ResultCode:  1005, // Assuming 1005 is "Queued"
        Description: "Queued", // Adjust based on actual response or mapping logic
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(respA)
}

func main() {
    http.HandleFunc("/sms/send", smsSendHandler)

    log.Println("Server starting on port 8081...")
    if err := http.ListenAndServe(":8081", nil); err != nil {
        log.Fatal(err)
    }
}
