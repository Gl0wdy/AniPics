package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func RandomPicProxy(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	if tag == "" {
		http.Error(w, "tag is required", http.StatusBadRequest)
		return
	}

	tags := []string{"waifu", "neko", "kitsune", "husbando", "fox_girl"}

	if tag == "random" {
		randomIndex := rand.Intn(len(tags))
		tag = tags[randomIndex]
	}

	apiURL := fmt.Sprintf("https://nekos.best/api/v2/%s", tag)
	metaResp, err := http.Get(apiURL)
	if err != nil || metaResp.StatusCode != http.StatusOK {
		http.Error(w, "failed to get metadata from nekos.best", http.StatusBadGateway)
		return
	}
	defer metaResp.Body.Close()

	var upstream struct {
		Results []struct {
			URL string `json:"url"`
		} `json:"results"`
	}

	if err := json.NewDecoder(metaResp.Body).Decode(&upstream); err != nil {
		http.Error(w, "invalid metadata response", http.StatusInternalServerError)
		return
	}

	if len(upstream.Results) == 0 || upstream.Results[0].URL == "" {
		http.Error(w, "no image found for this tag", http.StatusNotFound)
		return
	}

	imageURL := upstream.Results[0].URL

	client := &http.Client{Timeout: 10 * time.Second}
	imgResp, err := client.Get(imageURL)
	if err != nil || imgResp.StatusCode != http.StatusOK {
		http.Error(w, "failed to fetch image", http.StatusBadGateway)
		return
	}
	defer imgResp.Body.Close()

	for k, vv := range imgResp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	w.Header().Set("Cache-Control", "public, max-age=86400") // 24 часа кэша
	w.Header().Set("Access-Control-Allow-Origin", "*")       // если фронт на другом домене

	w.WriteHeader(imgResp.StatusCode)
	_, err = io.Copy(w, imgResp.Body)
	if err != nil {
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
