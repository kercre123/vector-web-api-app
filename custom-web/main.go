package main

import (
	"bytes"
	"context"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/digital-dream-labs/vector-go-sdk/pkg/vector"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	"github.com/fogleman/gg"
	"github.com/gorilla/websocket"
	"hz.tools/mjpeg"
)

const serverFiles string = "/var/www"
const sdkAddress string = "localhost:443"
const vizAddress string = "localhost:8888"

var robot *vector.Vector
var bcAssumption bool = false
var ctx context.Context
var camStreamEnable bool = false
var camStreamClient vectorpb.ExternalInterface_CameraFeedClient

var transCfg = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore SSL warnings
}

func initSDK() {
	var err error
	robot, err = vector.New(
		vector.WithTarget(sdkAddress),
		vector.WithToken(getGUID()),
	)
	if err != nil {
		log.Fatal(err)
	}
	ctx = context.Background()
}

func assumeBehaviorControl(priority string) {
	var controlRequest *vectorpb.BehaviorControlRequest
	if priority == "high" {
		controlRequest = &vectorpb.BehaviorControlRequest{
			RequestType: &vectorpb.BehaviorControlRequest_ControlRequest{
				ControlRequest: &vectorpb.ControlRequest{
					Priority: vectorpb.ControlRequest_OVERRIDE_BEHAVIORS,
				},
			},
		}
	} else {
		controlRequest = &vectorpb.BehaviorControlRequest{
			RequestType: &vectorpb.BehaviorControlRequest_ControlRequest{
				ControlRequest: &vectorpb.ControlRequest{
					Priority: vectorpb.ControlRequest_DEFAULT,
				},
			},
		}
	}
	go func() {
		start := make(chan bool)
		stop := make(chan bool)
		bcAssumption = true
		go func() {
			// * begin - modified from official vector-go-sdk
			r, err := robot.Conn.BehaviorControl(
				ctx,
			)
			if err != nil {
				log.Println(err)
				return
			}

			if err := r.Send(controlRequest); err != nil {
				log.Println(err)
				return
			}

			for {
				ctrlresp, err := r.Recv()
				if err != nil {
					log.Println(err)
					return
				}
				if ctrlresp.GetControlGrantedResponse() != nil {
					start <- true
					break
				}
			}

			for {
				select {
				case <-stop:
					if err := r.Send(
						&vectorpb.BehaviorControlRequest{
							RequestType: &vectorpb.BehaviorControlRequest_ControlRelease{
								ControlRelease: &vectorpb.ControlRelease{},
							},
						},
					); err != nil {
						log.Println(err)
						return
					}
					return
				default:
					continue
				}
			}
			// * end - modified from official vector-go-sdk
		}()
		for {
			select {
			case <-start:
				for {
					if bcAssumption {
						time.Sleep(time.Millisecond * 500)
					} else {
						break
					}
				}
				stop <- true
				return
			}
		}
	}()
}

func sayText(text string) {
	_, _ = robot.Conn.SayText(
		ctx,
		&vectorpb.SayTextRequest{
			Text:           text,
			UseVectorVoice: true,
			DurationScalar: 1.0,
		},
	)
}

func driveWheelsForward(lw float32, rw float32, lwtwo float32, rwtwo float32) {
	_, _ = robot.Conn.DriveWheels(
		ctx,
		&vectorpb.DriveWheelsRequest{
			LeftWheelMmps:   lw,
			RightWheelMmps:  rw,
			LeftWheelMmps2:  lwtwo,
			RightWheelMmps2: rwtwo,
		},
	)
}

func moveLift(speed float32) {
	_, _ = robot.Conn.MoveLift(
		ctx,
		&vectorpb.MoveLiftRequest{
			SpeedRadPerSec: speed,
		},
	)
}

func moveHead(speed float32) {
	_, _ = robot.Conn.MoveHead(
		ctx,
		&vectorpb.MoveHeadRequest{
			SpeedRadPerSec: speed,
		},
	)
}

func releaseBehaviorControl() {
	bcAssumption = false
}

func convertPixesTo16BitRGB(r uint32, g uint32, b uint32, a uint32) uint16 {
	R, G, B := int(r/257), int(g/257), int(b/257)

	return uint16((int(R>>3) << 11) |
		(int(G>>2) << 5) |
		(int(B>>3) << 0))
}

func convertPixelsToRawBitmap(image image.Image) []uint16 {
	imgHeight, imgWidth := image.Bounds().Max.Y, image.Bounds().Max.X
	bitmap := make([]uint16, 184*96)

	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			bitmap[(y)*184+(x)] = convertPixesTo16BitRGB(image.At(x, y).RGBA())
		}
	}
	return bitmap
}

func TextOnImg(text string, size float64) []byte {
	bgImage := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 184, Y: 96},
	})
	imgWidth := bgImage.Bounds().Dx()
	imgHeight := bgImage.Bounds().Dy()
	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(bgImage, 0, 0)

	if err := dc.LoadFontFace("/data/test.ttf", size); err != nil {
		fmt.Println(err)
		return nil
	}

	x := float64(imgWidth / 2)
	y := float64((imgHeight / 2))
	maxWidth := float64(imgWidth) - 35.0
	dc.SetColor(color.White)
	dc.DrawStringWrapped(text, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)
	buf := new(bytes.Buffer)
	bitmap := convertPixelsToRawBitmap(dc.Image())
	for _, ui := range bitmap {
		binary.Write(buf, binary.LittleEndian, ui)
	}
	os.WriteFile("/tmp/test.raw", buf.Bytes(), 0644)
	return buf.Bytes()
}

func imgOnFace(text string, size float64) {
	faceBytes := TextOnImg(text, size)
	_, _ = robot.Conn.DisplayFaceImageRGB(
		ctx,
		&vectorpb.DisplayFaceImageRGBRequest{
			FaceData:         faceBytes,
			DurationMs:       5000,
			InterruptRunning: true,
		},
	)
}

func playSound(buf []byte, filename string) string {
	os.WriteFile("/tmp/"+strings.TrimSpace(filename), buf, 0644)
	var pcmFile []byte
	if strings.Contains(filename, ".pcm") || strings.Contains(filename, ".raw") {
		fmt.Println("Assuming already pcm")
		pcmFile, _ = os.ReadFile("/tmp/" + strings.TrimSpace(filename))
	} else {
		conOutput, conError := exec.Command("ffmpeg", "-y", "-i", "/tmp/"+strings.TrimSpace(filename), "-f", "s16le", "-acodec", "pcm_s16le", "-ar", "16000", "-ac", "1", "/tmp/output.pcm").Output()
		if conError != nil {
			fmt.Println(conError)
			return conError.Error()
		}
		fmt.Println("FFMPEG output: " + string(conOutput))
		pcmFile, _ = os.ReadFile("/tmp/output.pcm")
	}
	var audioChunks [][]byte
	for len(pcmFile) >= 1024 {
		audioChunks = append(audioChunks, pcmFile[:1024])
		pcmFile = pcmFile[1024:]
	}
	var audioClient vectorpb.ExternalInterface_ExternalAudioStreamPlaybackClient
	audioClient, _ = robot.Conn.ExternalAudioStreamPlayback(
		ctx,
	)
	audioClient.SendMsg(&vectorpb.ExternalAudioStreamRequest{
		AudioRequestType: &vectorpb.ExternalAudioStreamRequest_AudioStreamPrepare{
			AudioStreamPrepare: &vectorpb.ExternalAudioStreamPrepare{
				AudioFrameRate: 16000,
				AudioVolume:    100,
			},
		},
	})
	fmt.Println(len(audioChunks))
	for _, chunk := range audioChunks {
		audioClient.SendMsg(&vectorpb.ExternalAudioStreamRequest{
			AudioRequestType: &vectorpb.ExternalAudioStreamRequest_AudioStreamChunk{
				AudioStreamChunk: &vectorpb.ExternalAudioStreamChunk{
					AudioChunkSizeBytes: 1024,
					AudioChunkSamples:   chunk,
				},
			},
		})
		time.Sleep(time.Millisecond * 30)
	}
	audioClient.SendMsg(&vectorpb.ExternalAudioStreamRequest{
		AudioRequestType: &vectorpb.ExternalAudioStreamRequest_AudioStreamComplete{
			AudioStreamComplete: &vectorpb.ExternalAudioStreamComplete{},
		},
	})
	os.Remove("/tmp/" + strings.TrimSpace(filename))
	os.Remove("/tmp/output.pcm")
	return "success"
}

func skipOnboarding() {
	url := "http://" + sdkAddress + "/consolefunccall"
	var form = []byte("func=Exit Onboarding - Mark Complete&args=")
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Transport: transCfg}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func sdkAuth(username string, password string) string {
	cmd1 := exec.Command("/bin/rm", "-rf", "/data/protected")
	cmd2 := exec.Command("/bin/mkdir", "-p", "/data/protected")
	cmd1.Run()
	cmd2.Run()
	url := "https://accounts.api.anki.com/1/sessions"
	var credsForm = []byte("username=" + username + "&password=" + password)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(credsForm))
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
		log.Println(sessionToken)
		tokenEnc := b64.StdEncoding.EncodeToString([]byte(sessionToken))
		url2 := "https://" + sdkAddress + "/v1/user_authentication"
		var tokenJSON = []byte(`{"user_session_id": "` + tokenEnc + `"}`)
		req, _ := http.NewRequest("POST", url2, bytes.NewBuffer(tokenJSON))
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
			os.WriteFile("/data/protected/authStatus", []byte("noguid"), 0644)
			return "error2"
		}
		if strings.Contains(resp.Status, "403") {
			cmd1.Run()
			os.WriteFile("/data/protected/authStatus", []byte("noguid"), 0644)
			return "error2"
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
			os.WriteFile("/data/protected/authStatus", []byte("noguid"), 0644)
			return "error2"
		}
		clientGUIDenc := guid.ClientTokenGUID
		clientGUIDdec, _ := b64.StdEncoding.DecodeString(clientGUIDenc)
		clientGUID := string(clientGUIDdec)
		url3 := "https://" + sdkAddress + "/v1/pull_jdocs"
		var jdocJSON = []byte(`{"jdoc_types": [0, 1, 2, 3]}`)
		req2, _ := http.NewRequest("POST", url3, bytes.NewBuffer(jdocJSON))
		req2.Header.Set("Authorization", "Bearer "+clientGUID)
		req2.Header.Set("Content-Type", "application/json")
		client2 := &http.Client{Transport: transCfg}
		resp2, err := client2.Do(req2)
		if err != nil {
			panic(err)
		}
		defer resp2.Body.Close()
		os.WriteFile("/data/protected/client.guid", clientGUIDdec, 0644)
		os.WriteFile("/data/protected/authStatus", []byte("success"), 0644)
		skipOnboarding()
		return "success"
	} else {
		cmd1.Run()
		return "unknown"
	}
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
		url := "https://" + sdkAddress + "/v1/update_settings"
		var updateJSON = []byte(`{"update_settings": true, "settings": {"custom_eye_color": {"enabled": true, "hue": ` + hue + `, "saturation": ` + sat + `} } }`)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(updateJSON))
		req.Header.Set("Authorization", "Bearer "+clientGUID)
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

func getSDKSettings() []byte {
	clientGUID := getGUID()
	url := "https://" + sdkAddress + "/v1/update_settings"
	var updateJSON = []byte(`{"update_settings": false}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(updateJSON))
	req.Header.Set("Authorization", "Bearer "+clientGUID)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Transport: transCfg}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	settingsResponse := body
	return settingsResponse
}

func setPresetEyeColor(value string) {
	clientGUID := getGUID()
	if !strings.Contains(clientGUID, "error") {
		url := "https://" + sdkAddress + "/v1/update_settings"
		var updateJSON = []byte(`{"update_settings": true, "settings": {"custom_eye_color": {"enabled": false}, "eye_color": ` + value + `} }`)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(updateJSON))
		req.Header.Set("Authorization", "Bearer "+clientGUID)
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
		url := "https://" + sdkAddress + "/v1/update_settings"
		var updateJSON = []byte(`{"update_settings": true, "settings": {"` + setting + `": "` + value + `" } }`)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(updateJSON))
		req.Header.Set("Authorization", "Bearer "+clientGUID)
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
		url := "https://" + sdkAddress + "/v1/update_settings"
		var updateJSON = []byte(`{"update_settings": true, "settings": {"` + setting + `": ` + value + ` } }`)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(updateJSON))
		req.Header.Set("Authorization", "Bearer "+clientGUID)
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
	if _, err := os.Stat("/wirefiles/escape"); err == nil {
		return "escapepod"
	}
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

func getCustomSettings() string {
	var snore string
	var rainbowEyes string
	var freqStatus string
	var serverStatus string
	var jsonResponse string
	var alexaStatus string
	var soundStatus string
	var vicosVersion string
	var robotName string
	var serialNumber string
	var snowglobeStatus string
	var robotBranch string
	if _, err := os.Stat("/data/data/snore_disable"); err == nil {
		snore = "off"
	} else {
		snore = "on"
	}
	if _, err := os.Stat("/data/data/rainboweyes"); err == nil {
		rainbowEyes = "on"
	} else {
		rainbowEyes = "off"
	}
	if _, err := os.Stat("/data/data/freqStatus"); err == nil {
		fileBytes, err := ioutil.ReadFile("/data/data/freqStatus")
		if err != nil {
			log.Println("no freq status")
		}
		freqStatus = string(fileBytes)
	} else {
		freqStatus = "balanced"
	}
	if _, err := os.Stat("/wirefiles/escape"); err == nil {
		serverStatus = "escape"
	} else {
		serverStatus = "prod"
	}
	if _, err := os.Stat("/data/data/com.anki.victor/persistent/alexa/optedIn"); err == nil {
		alexaStatus = "on"
	} else {
		alexaStatus = "off"
	}
	if _, err := os.Stat("/anki/data/assets/cozmo_resources/sound/version"); err == nil {
		fileBytes, err := ioutil.ReadFile("/anki/data/assets/cozmo_resources/sound/version")
		if err != nil {
			log.Println("no sound status")
		}
		soundStatus = strings.TrimSpace(string(fileBytes))
	} else {
		soundStatus = "1.8.0.6051"
	}
	cmd1 := exec.Command("/bin/bash", "/sbin/vector-ctrl", "info_print")
	cmd1.Run()
	versionBytes, err := ioutil.ReadFile("/data/data/vicosVersion")
	if err != nil {
		log.Println("no version string")
	}
	vicosVersion = strings.TrimSpace(string(versionBytes))
	serialBytes, err := ioutil.ReadFile("/data/data/serialNumber")
	if err != nil {
		log.Println("no serial string")
	}
	serialNumber = strings.TrimSpace(string(serialBytes))
	nameBytes, err := ioutil.ReadFile("/data/data/robotName")
	if err != nil {
		log.Println("no name string")
	}
	robotName = strings.TrimSpace(string(nameBytes))
	branchBytes, err := ioutil.ReadFile("/data/data/robotBranch")
	if err != nil {
		log.Println("no branch string")
	}
	robotBranch = strings.TrimSpace(string(branchBytes))
	if _, err := os.Stat("/data/data/snowglobe"); err == nil {
		snowglobeStatus = "on"
	} else {
		snowglobeStatus = "off"
	}
	jsonResponse = `{"snore_status": "` + snore + `", "rainboweyes_status": "` + rainbowEyes + `", "freq_status": "` + freqStatus + `", "server_status": "` + serverStatus + `", "alexa_status": "` + alexaStatus + `", "sound_status": "` + soundStatus + `", "vicos_version": "` + vicosVersion + `", "robot_esn": "` + serialNumber + `", "robot_name": "` + robotName + `", "robot_branch": "` + robotBranch + `", "snowglobe_status": "` + snowglobeStatus + `"}`
	return jsonResponse
}

func sendSocketMessage(message string) {
	socketUrl := "ws://" + vizAddress + "/socket"
	conn, _, err1 := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err1 != nil {
		log.Fatal("Error connecting to Websocket Server:", err1)
	}
	defer conn.Close()
	err2 := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err2 != nil {
		log.Println("Error during writing to websocket:", err2)
	}
	err3 := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err3 != nil {
		log.Println("Error during closing websocket:", err3)
		return
	}
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
		fmt.Fprint(w, authStatus)
		return
	case r.URL.Path == "/api/cloud_intent":
		intent := r.FormValue("intent")
		sendSocketMessage(`{"type":"data","module":"intents","data":{"intentType":"cloud","request":"` + intent + `"}}`)
		fmt.Fprintf(w, "done")
		return
	case r.URL.Path == "/api/set_timer":
		secs := r.FormValue("secs")
		sendSocketMessage(`{"type":"data","module":"intents","data":{"intentType":"cloud","request":"{ \"intent\" : \"intent_clock_settimer_extend\", \"parameters\" : \"{\\\"timer_duration\\\":\\\"` + secs + `'\\\",\\\"unit\\\":\\\"s\\\"}\\n\"}"}}`)
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
		fmt.Fprintf(w, hue+sat)
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
		location := r.FormValue("location")
		setSettingSDKstring("default_location", location)
		fmt.Fprintf(w, "done")
		return
	case r.URL.Path == "/api/timezone":
		timezone := r.FormValue("timezone")
		setSettingSDKstring("time_zone", timezone)
		fmt.Fprintf(w, "done")
		return
	case r.URL.Path == "/api/stop_timer":
		sendSocketMessage(`{"type":"data","module":"intents","data":{"intentType":"cloud","request":"{ \"intent\" : \"intent_global_stop_extend\", \"metadata\" : \"text: stop the timer  confidence: 0.000000  handler: HOUNDIFY\", \"parameters\" : \"{\\\"entity_behavior_stoppable\\\":\\\"timer\\\"}\\n\", \"time\" : 1649608984, \"type\" : \"result\" }"}}`)
		fmt.Fprintf(w, "done")
		return
	case r.URL.Path == "/api/get_auth_status":
		authStatus := getAuthStatus()
		fmt.Fprint(w, authStatus)
		return
	case r.URL.Path == "/api/get_sdk_settings":
		settings := getSDKSettings()
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(settings)
		return
	case r.URL.Path == "/api/get_custom_settings":
		settings := getCustomSettings()
		fmt.Fprint(w, settings)
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
		skipOnboarding()
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
	case r.URL.Path == "/api/snowglobe":
		fmt.Fprintf(w, "executing")
		cmd := exec.Command("/bin/bash", "/sbin/vector-ctrldd", "snowglobe")
		cmd.Run()
		return
	case r.URL.Path == "/api/initSDK":
		initSDK()
		fmt.Fprintf(w, "done")
		return
	case r.URL.Path == "/api/assume_behavior_control":
		fmt.Fprintf(w, "done")
		assumeBehaviorControl(r.FormValue("priority"))
	case r.URL.Path == "/api/release_behavior_control":
		fmt.Fprintf(w, "done")
		releaseBehaviorControl()
		return
	case r.URL.Path == "/api/say_text":
		sayText(r.FormValue("text"))
		fmt.Fprintf(w, "said")
		return
	case r.URL.Path == "/api/move_wheels":
		lw, _ := strconv.Atoi(r.FormValue("lw"))
		rw, _ := strconv.Atoi(r.FormValue("rw"))
		driveWheelsForward(float32(lw), float32(rw), float32(lw), float32(rw))
		fmt.Fprintf(w, "")
		return
	case r.URL.Path == "/api/begin_cam_stream":
		camStreamClient, _ = robot.Conn.CameraFeed(ctx, &vectorpb.CameraFeedRequest{})
		camStreamEnable = true
		fmt.Fprintf(w, "success")
		return
	case r.URL.Path == "/api/stop_cam_stream":
		camStreamEnable = false
		camStreamClient = nil
		fmt.Fprintf(w, "success")
		return
	case r.URL.Path == "/api/move_lift":
		speed, _ := strconv.Atoi(r.FormValue("speed"))
		moveLift(float32(speed))
		fmt.Fprintf(w, "")
		return
	case r.URL.Path == "/api/move_head":
		speed, _ := strconv.Atoi(r.FormValue("speed"))
		moveHead(float32(speed))
		fmt.Fprintf(w, "")
		return
	case r.URL.Path == "/api/img_on_face":
		text := r.FormValue("text")
		sizeInt, _ := strconv.Atoi(r.FormValue("size"))
		if text == "" {
			text = "test"
		}
		size := float64(sizeInt)
		fmt.Println(size)
		imgOnFace(text, size)
		fmt.Fprintf(w, "done :)")
		return
	case r.URL.Path == "/api/play_sound":
		var buf bytes.Buffer
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Fprintf(w, "error")
			return
		}
		io.Copy(&buf, file)
		playSound(buf.Bytes(), header.Filename)
		fmt.Fprintf(w, "uploaded")
		return
	case r.URL.Path == "/api/sound_version":
		version := r.FormValue("version")
		cmd := exec.Command("/bin/bash", "/sbin/vector-ctrl", "pingtest")
		cmd.Run()
		var versions string = "1.8.1.6051 1.8.0.6021 1.7.0.3412 1.6.0.3331 1.5.0.3009 1.4.1.2806 1.3.0.2510 1.2.3.2506 1.2.2.2353 1.2.1.2343 1.1.1.2107 1.1.0.2106 1.0.2.1804 1.0.1.1768 1.0.0.1741"
		if strings.Contains(versions, version) {
			if _, err := os.Stat("/tmp/testPing"); err == nil {
				testBytes, err := ioutil.ReadFile("/tmp/testPing")
				if err != nil {
					log.Println("no test string")
					fmt.Fprintf(w, "error")
				}
				if strings.Contains(string(testBytes), "success") {
					cmd1 := exec.Command("/bin/rm", "-f", "/tmp/testPing")
					cmd1.Run()
					fmt.Fprintf(w, "executing")
					cmd2 := exec.Command("/bin/bash", "/sbin/vector-ctrldd", "sound_version", version)
					cmd2.Run()
				}
			} else {
				fmt.Fprintf(w, "error")
			}
		} else {
			fmt.Fprintf(w, "error")
		}
		return
	case r.URL.Path == "/api/freq":
		perfPreset := r.FormValue("freq")
		if strings.Contains(perfPreset, "performance") {
			cmd := exec.Command("/bin/bash", "/sbin/vector-ctrl", "freq", "1267200", "800000")
			cmd.Run()
			os.WriteFile("/data/data/freqStatus", []byte("performance"), 0644)
			fmt.Fprintf(w, "done")
		} else if strings.Contains(perfPreset, "balanced") {
			cmd := exec.Command("/bin/bash", "/sbin/vector-ctrl", "freq", "733333", "500000")
			cmd.Run()
			os.WriteFile("/data/data/freqStatus", []byte("balanced"), 0644)
			fmt.Fprintf(w, "done")
		} else if strings.Contains(perfPreset, "stock") {
			cmd := exec.Command("/bin/bash", "/sbin/vector-ctrl", "freq", "533333", "400000")
			cmd.Run()
			os.WriteFile("/data/data/freqStatus", []byte("stock"), 0644)
			fmt.Fprintf(w, "done")
		} else {
			fmt.Fprintf(w, "must be performance, balanced, or stock")
		}
		return
	}
}

func main() {
	camStream := mjpeg.NewStream()
	i := image.NewGray(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 640, Y: 360},
	})
	go func() {
		for {
			if camStreamEnable {
				response, _ := camStreamClient.Recv()
				imageBytes := response.GetData()
				img, _, _ := image.Decode(bytes.NewReader(imageBytes))
				camStream.Update(img)
			} else {
				for j := range i.Pix {
					i.Pix[j] = uint8(rand.Uint32())
				}

				time.Sleep(time.Second)
				camStream.Update(i)
			}
		}
	}()
	http.HandleFunc("/api/", apiHandler)
	fileServer := http.FileServer(http.Dir(serverFiles))
	http.Handle("/", fileServer)
	http.Handle("/stream", camStream)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
