package handlers

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/mjibson/go-dsp/fft"
	"log"
	"math"
	"math/cmplx"
	"net/http"
	"time"
)

type FrequencyResponse struct {
	Frequency float64 `json:"frequency"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("failed to read message:", err)
			break
		}

		if messageType == websocket.BinaryMessage {
			audioData := bytesToFloat64Array(message)

			// apply FFT on data
			fftResult := fft.FFTReal(audioData)

			// compute spectre magnitude FFT
			magnitude := make([]float64, len(fftResult)/2)
			for i := range magnitude {
				magnitude[i] = cmplx.Abs(fftResult[i])
			}

			// search main frequency
			maxFreqIndex := findMaxFrequencyIndex(magnitude)

			//TODO use real sample rate 44100 hz
			sampleRate := 44100.
			freq := float64(maxFreqIndex) * sampleRate / float64(len(fftResult))

			log.Printf("found main frequency : %.2f Hz", freq)
			response := FrequencyResponse{Frequency: freq}
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				log.Println("failed to marshall response:", err)
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, jsonResponse); err != nil {
				log.Println("failed to send message to client:", err)
				break
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func findMaxFrequencyIndex(data []float64) int {
	maxIdx := 0
	maxVal := 0.0

	for i, val := range data {
		if val > maxVal {
			maxVal = val
			maxIdx = i
		}
	}

	return maxIdx
}

func bytesToFloat64Array(data []byte) []float64 {
	floatData := make([]float64, len(data)/8)
	for i := 0; i < len(data)/8; i++ {
		floatData[i] = math.Float64frombits(
			uint64(data[8*i]) |
				uint64(data[8*i+1])<<8 |
				uint64(data[8*i+2])<<16 |
				uint64(data[8*i+3])<<24 |
				uint64(data[8*i+4])<<32 |
				uint64(data[8*i+5])<<40 |
				uint64(data[8*i+6])<<48 |
				uint64(data[8*i+7])<<56,
		)
	}

	return floatData
}
