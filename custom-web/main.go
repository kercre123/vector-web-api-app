package main


import (
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "os"
    "io/ioutil"
    "bytes"
    "time"
    "encoding/json"
    "strings"
    b64 "encoding/base64"
    "crypto/tls"
    "errors"
    "github.com/rgamba/evtwebsocket"

)

var transCfg = &http.Transport{
 TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore SSL warnings
}

var serverFiles = "/var/www"
var interfaceLocation = "/sbin/custom-web-interface"
const address string = "localhost:8888"
var c = evtwebsocket.Conn{}

func sdkAuth(username string, password string) string {
    cmd1 := exec.Command("/bin/rm", "-rf", "/data/protected")
    cmd2 := exec.Command("/bin/mkdir", "-p", "/data/protected")
    cmd1.Run()
    cmd2.Run()
    url := "https://accounts.api.anki.com/1/sessions"
    var credsForm = []byte("username=" + username + "&password=" + password)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(credsForm))
    req.Header.Set("Anki-App-Key", "luyain9ep5phahP8aph8xa")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    client := &http.Client{Transport: transCfg}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    sessionsResponse := string(body)
    if strings.Contains(sessionsResponse, "invalid_username_or_password") {
        cmd1.Run()
        authStatus := "error"
        return authStatus
    } else if strings.Contains(sessionsResponse, "session_token") {
       type SessionsResponses struct {
        Session struct {
            SessionToken string    `json:"session_token"`
            UserID       string    `json:"user_id"`
            Scope        string    `json:"scope"`
            TimeCreated  time.Time `json:"time_created"`
            TimeExpires  time.Time `json:"time_expires"`
        } `json:"session"`
        User struct {
            UserID               string      `json:"user_id"`
            DriveGuestID         string      `json:"drive_guest_id"`
            PlayerID             string      `json:"player_id"`
            CreatedByAppName     string      `json:"created_by_app_name"`
            CreatedByAppVersion  string      `json:"created_by_app_version"`
            CreatedByAppPlatform string      `json:"created_by_app_platform"`
            Dob                  string      `json:"dob"`
            Email                string      `json:"email"`
            FamilyName           interface{} `json:"family_name"`
            Gender               interface{} `json:"gender"`
            GivenName            interface{} `json:"given_name"`
            Username             string      `json:"username"`
            EmailIsVerified      bool        `json:"email_is_verified"`
            EmailFailureCode     interface{} `json:"email_failure_code"`
            EmailLang            string      `json:"email_lang"`
            PasswordIsComplex    bool        `json:"password_is_complex"`
            Status               string      `json:"status"`
            TimeCreated          time.Time   `json:"time_created"`
            DeactivationReason   interface{} `json:"deactivation_reason"`
            PurgeReason          interface{} `json:"purge_reason"`
            EmailIsBlocked       bool        `json:"email_is_blocked"`
            NoAutodelete         bool        `json:"no_autodelete"`
            IsEmailAccount       bool        `json:"is_email_account"`
        } `json:"user"`
    }
    var sessionss SessionsResponses
    json.Unmarshal([]byte(sessionsResponse), &sessionss)
    sessionToken := sessionss.Session.SessionToken
    tokenEnc := b64.StdEncoding.EncodeToString([]byte(sessionToken))
    url2 := "https://localhost:443/v1/user_authentication"
    var tokenJSON = []byte(`{"user_session_id": "` + tokenEnc + `"}`)
    req, err := http.NewRequest("POST", url2, bytes.NewBuffer(tokenJSON))
    req.Header.Set("Accept", "/")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{Transport: transCfg}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    log.Println(resp.Status)
    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    if strings.Contains(resp.Status, "401") {
        cmd1.Run()
        return "error2"
        os.WriteFile("/data/protected/authStatus", []byte("noguid"), 0644)
    }
    if strings.Contains(resp.Status, "403") {
        cmd1.Run()
        return "error2"
        os.WriteFile("/data/protected/authStatus", []byte("noguid"), 0644)
    }
    body, _ := ioutil.ReadAll(resp.Body)
    guidResponse := string(body)
    fmt.Println("response Body:", guidResponse)
    type GUIDRJson struct {
        Status struct {
            Code int `json:"code"`
        } `json:"status"`
        Code            int    `json:"code"`
        ClientTokenGUID string `json:"client_token_guid"`
    }
    var guid GUIDRJson
    json.Unmarshal([]byte(guidResponse), &guid)
    if guid.Code == 0 {
        return "error2"
        os.WriteFile("/data/protected/authStatus", []byte("noguid"), 0644)
    }
    clientGUIDenc := guid.ClientTokenGUID
    clientGUIDdec, _ := b64.StdEncoding.DecodeString(clientGUIDenc)
    clientGUID := string(clientGUIDdec)
    url3 := "https://localhost:443/v1/pull_jdocs"
    var jdocJSON = []byte(`{"jdoc_types": [0, 1, 2, 3]}`)
    req2, err := http.NewRequest("POST", url3, bytes.NewBuffer(jdocJSON))
    req2.Header.Set("Authorization", "Bearer " + clientGUID)
    req2.Header.Set("Content-Type", "application/json")
    client2 := &http.Client{Transport: transCfg}
    resp2, err := client2.Do(req2)
    if err != nil {
        panic(err)
    }
    defer resp2.Body.Close()
    os.WriteFile("/data/protected/client.guid", clientGUIDdec, 0644)
    os.WriteFile("/data/protected/authStatus", []byte("success"), 0644)
    return "success"
    } else {
        cmd1.Run()
        return "unknown"
    }
    return "unknown"
}

func getGUID() string {
    fileBytes, err := ioutil.ReadFile("/data/protected/client.guid")
    if err != nil {
        return "error"
    }
    clientGUID := string(fileBytes)
    return clientGUID
}

func setCustomEyeColor(hue string, sat string) {
    clientGUID := getGUID()
    if !strings.Contains(clientGUID, "error") {
        url := "https://localhost:443/v1/update_settings"
        var updateJSON = []byte(`{"update_settings": true, "settings": {"custom_eye_color": {"enabled": true, "hue": ` + hue + `, "saturation": ` + sat + `} } }`)
        req, err := http.NewRequest("POST", url, bytes.NewBuffer(updateJSON))
        req.Header.Set("Authorization", "Bearer " + clientGUID)
        req.Header.Set("Content-Type", "application/json")
        client := &http.Client{Transport: transCfg}
        resp, err := client.Do(req)
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()
    } else {
        log.Println("GUID not there")
    }
}

func setPresetEyeColor(value string) {
    clientGUID := getGUID()
    if !strings.Contains(clientGUID, "error") {
        url := "https://localhost:443/v1/update_settings"
        var updateJSON = []byte(`{"update_settings": true, "settings": {"custom_eye_color": {"enabled": false}, "eye_color": ` + value + `} }`)
        req, err := http.NewRequest("POST", url, bytes.NewBuffer(updateJSON))
        req.Header.Set("Authorization", "Bearer " + clientGUID)
        req.Header.Set("Content-Type", "application/json")
        client := &http.Client{Transport: transCfg}
        resp, err := client.Do(req)
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()
    } else {
        log.Println("GUID not there")
    }
}

func setSettingSDKstring(setting string, value string) {
    clientGUID := getGUID()
    if !strings.Contains(clientGUID, "error") {
        url := "https://localhost:443/v1/update_settings"
        var updateJSON = []byte(`{"update_settings": true, "settings": {"` + setting + `": "` + value + `" } }`)
        req, err := http.NewRequest("POST", url, bytes.NewBuffer(updateJSON))
        req.Header.Set("Authorization", "Bearer " + clientGUID)
        req.Header.Set("Content-Type", "application/json")
        client := &http.Client{Transport: transCfg}
        resp, err := client.Do(req)
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()
    } else {
        log.Println("GUID not there")
    }
}

func setSettingSDKintbool(setting string, value string) {
    clientGUID := getGUID()
    if !strings.Contains(clientGUID, "error") {
        url := "https://localhost:443/v1/update_settings"
        var updateJSON = []byte(`{"update_settings": true, "settings": {"` + setting + `": ` + value + ` } }`)
        req, err := http.NewRequest("POST", url, bytes.NewBuffer(updateJSON))
        req.Header.Set("Authorization", "Bearer " + clientGUID)
        req.Header.Set("Content-Type", "application/json")
        client := &http.Client{Transport: transCfg}
        resp, err := client.Do(req)
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()
    } else {
        log.Println("GUID not there")
    }
}

func getAuthStatus() string {
    if _, err := os.Stat("/data/protected/authStatus"); err == nil {
        fileBytes, err := ioutil.ReadFile("/data/protected/authStatus")
        if err != nil {
            return "unknown"
        }
        authStatusFileString := string(fileBytes)
        if strings.Contains(authStatusFileString, "success") {
            return "authorized"
        } else if strings.Contains(authStatusFileString, "noguid") {
            return "notauthorized2"
        } else {
            return "unknown"
        }
    } else if errors.Is(err, os.ErrNotExist) {
      return "notauthorized1"
  } else {
    return "notauthorized1"
}
}

func launchIntent(intent string) {
    msg := evtwebsocket.Msg{
        Body: []byte(`{"type":"data","module":"intents","data":{"intentType":"cloud","request":"` + intent + `"}}`),
    }
    c.Send(msg)
}

func setTimer(seconds string) {
    msg := evtwebsocket.Msg{
        Body: []byte(`{"type":"data","module":"intents","data":{"intentType":"cloud","request":"{ \"intent\" : \"intent_clock_settimer_extend\", \"parameters\" : \"{\\\"timer_duration\\\":\\\"` + seconds + `'\\\",\\\"unit\\\":\\\"s\\\"}\\n\"}"}}`),
    }
    c.Send(msg)
}

func stopTimer() {
    msg := evtwebsocket.Msg{
        Body: []byte(`{"type":"data","module":"intents","data":{"intentType":"cloud","request":"{ \"intent\" : \"intent_global_stop_extend\", \"metadata\" : \"text: stop the timer  confidence: 0.000000  handler: HOUNDIFY\", \"parameters\" : \"{\\\"entity_behavior_stoppable\\\":\\\"timer\\\"}\\n\", \"time\" : 1649608984, \"type\" : \"result\" }"}}`),
    }
    c.Send(msg)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    switch {
    default: 
        http.Error(w, "not found", http.StatusNotFound)
        return
    case r.URL.Path == "/api/sdk_auth":
        username := r.FormValue("username")
        password := r.FormValue("password")
        authStatus := sdkAuth(username, password)
        fmt.Fprintf(w, authStatus)
        return
    case r.URL.Path == "/api/cloud_intent":
        intent := r.FormValue("intent")
        launchIntent(intent)
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/set_timer":
        secs := r.FormValue("secs")
        setTimer(secs)
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/eye_color":
        eye_color := r.FormValue("color")
        setPresetEyeColor(eye_color)
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/custom_eye_color":
        hue := r.FormValue("hue")
        sat := r.FormValue("sat")
        setCustomEyeColor(hue, sat)
        fmt.Fprintf(w, hue + sat)
        return
    case r.URL.Path == "/api/volume":
        volume := r.FormValue("volume")
        setSettingSDKintbool("master_volume", volume)
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/locale":
        locale := r.FormValue("locale")
        setSettingSDKstring("locale", locale)
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/location":
        text := r.FormValue("text")
        setSettingSDKstring("default_location", text)
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/timezone":
        text := r.FormValue("text")
        setSettingSDKstring("time_zone", text)
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/stop_timer":
        stopTimer()
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/snap_pic":
        cmd := exec.Command("/bin/bash", "/anki/bin/vector-ctrl", "-pic")
        cmd.Run()
        fmt.Fprintf(w, "pic snapped, at /tmp/img.jpg, use /api/get_pic")
        return
    case r.URL.Path == "/api/get_auth_status":
        authStatus := getAuthStatus()
        fmt.Fprintf(w, authStatus)
        return
    case r.URL.Path == "/api/snore_status":
        if _, err := os.Stat("/data/protected/authStatus"); err == nil {
            fmt.Fprintf(w, "off")
        } else {
            fmt.Fprintf(w, "on")
        }
        return
    case r.URL.Path == "/api/get_current_settings":
        fileBytes, err := ioutil.ReadFile("/data/data/com.anki.victor/persistent/jdocs/vic.RobotSettings.json")
        if err != nil {
            fmt.Fprintf(w, "error getting settings (file not there)")
        }
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/octet-stream")
        w.Write(fileBytes)
        return
    case r.URL.Path == "/api/rainbow_status":
        if _, err := os.Stat("/data/data/rainboweyes"); err == nil {
            fmt.Fprintf(w, "on")
        } else {
            fmt.Fprintf(w, "off")
        }
        return
    case r.URL.Path == "/api/server_status":
        if _, err := os.Stat("/wirefiles/escape"); err == nil {
            fmt.Fprintf(w, "escape")
        } else {
            fmt.Fprintf(w, "prod")
        }
        return
    case r.URL.Path == "/api/rainbow_on":
        cmd := exec.Command("/bin/bash", "/sbin/vector-ctrl", "rainbowon")
        cmd.Run()
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/rainbow_off":
        cmd := exec.Command("/bin/bash", "/sbin/vector-ctrl", "rainbowoff")
        cmd.Run()
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/snore_enable":
        fmt.Fprintf(w, "executing")
        cmd := exec.Command("/bin/bash", "/sbin/vector-ctrldd", "snore_enable")
        cmd.Run()
        return
    case r.URL.Path == "/api/snore_disable":
        fmt.Fprintf(w, "executing")
        cmd := exec.Command("/bin/bash", "/sbin/vector-ctrldd", "snore_disable")
        cmd.Run()
        return
    case r.URL.Path == "/api/time_format_12":
        setSettingSDKintbool("clock_24_hour", "false")
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/time_format_24":
        setSettingSDKintbool("clock_24_hour", "true")
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/skip_onboarding":
        cmd := exec.Command("/bin/bash", interfaceLocation, "skip_onboarding")
        cmd.Run()
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/temp_c":
        setSettingSDKintbool("temp_is_fahrenheit", "false")
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/temp_f":
        setSettingSDKintbool("temp_is_fahrenheit", "true")
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/button_hey_vector":
        setSettingSDKintbool("button_wakeword", "0")
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/button_alexa":
        setSettingSDKintbool("button_wakeword", "1")
        fmt.Fprintf(w, "done")
        return
    case r.URL.Path == "/api/server_escape":
        fmt.Fprintf(w, "executing")
        cmd := exec.Command("/bin/bash", "/sbin/vector-ctrldd", "server_escape")
        cmd.Run()
        return
    case r.URL.Path == "/api/server_prod":
        fmt.Fprintf(w, "executing")
        cmd := exec.Command("/bin/bash", "/sbin/vector-ctrldd", "server_prod")
        cmd.Run()
        return
    case r.URL.Path == "/api/get_pic":
        fileBytes, err := ioutil.ReadFile("/tmp/img.jpg")
        if err != nil {
            panic(err)
        }
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/octet-stream")
        w.Write(fileBytes)
        cmd := exec.Command("/bin/rm", "/tmp/img.jpg")
        cmd.Run()
        return
    }
}


func main() {
    c.Dial("ws://localhost:8888/socket", "")
    http.HandleFunc("/api/", apiHandler)
    fileServer := http.FileServer(http.Dir(serverFiles))
    http.Handle("/", fileServer)


    fmt.Printf("Starting server at port 8080\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
