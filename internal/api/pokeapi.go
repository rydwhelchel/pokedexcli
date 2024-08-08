package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	baseUrl   = "https://pokeapi.co/api/v2"
	cacheTime = 30 * time.Second
	// Rolls 1-pokeBaseExperience, if you roll below catchThresh, you catch
	// For ref: Caterpie is 39 and Mewtwo is 340
	catchThresh = 30
)

var cache Cache

func GetURL(url string) ([]byte, error) {
	cached := checkCache(url)
	var body []byte
	if cached == nil {
		res, err := http.Get(url)
		if err != nil {
			return nil, errors.New("error making GET to " + url + " and error was " + err.Error())
		}
		b, err := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			return nil, errors.New(fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body))
		}
		if err != nil {
			return nil, errors.New("error reading body " + err.Error())
		}
		body = b
		cache.Add(url, b)
	} else {
		// If there was a cached value, use that instead
		body = cached
		cache.Update(url)
	}
	return body, nil
}

func checkCache(key string) []byte {
	// Cache duration will only be zero if we have not instantiated Cache
	if cache.duration == 0 {
		log.Println("Creating new cache")
		cache = NewCache(cacheTime)
		return nil
	}
	if val, ok := cache.cache[key]; ok {
		return val.contents
	}
	return nil
}

func GetNextLocations(conf *Config) []Location {
	if conf.next == "" {
		conf.next = baseUrl + "/location-area"
	}

	body, err := GetURL(conf.next)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error making request to %s - %s", conf.next, err))
	}

	batch := LocationBatch{}
	err = json.Unmarshal(body, &batch)
	if err != nil {
		log.Fatalf("Error unmarshalling body into LocationBatch: %s", err)
	}

	conf.prev = batch.Previous
	conf.next = batch.Next

	return batch.Results
}

func GetPrevLocations(conf *Config) []Location {
	if conf.prev == "" {
		return []Location{}
	}

	body, err := GetURL(conf.prev)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error making request to %s - %s", conf.prev, err))
	}

	batch := LocationBatch{}
	err = json.Unmarshal(body, &batch)
	if err != nil {
		log.Fatalf("Error unmarshalling body into LocationBatch: %s", err)
	}

	conf.prev = batch.Previous
	conf.next = batch.Next

	return batch.Results
}

func GetExploreArea(area string) []string {
	url := baseUrl + "/location-area/" + area
	body, err := GetURL(url)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error making request to %s - %s", url, err))
	}

	areaDetails := LocationAreaDetails{}
	err = json.Unmarshal(body, &areaDetails)
	if err != nil {
		log.Fatalf("Error unmarshalling body into LocationAreaDetails: %s", err)
	}

	pokemon := []string{}
	for _, enc := range areaDetails.PokemonEncounters {
		pokemon = append(pokemon, enc.Pokemon.Name)
	}

	return pokemon
}

func GetPokemon(pokemonName string) Pokemon {
	url := baseUrl + "/pokemon/" + pokemonName
	body, err := GetURL(url)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error making request to %s - %s", url, err))
	}

	pokemon := Pokemon{}
	err = json.Unmarshal(body, &pokemon)
	return pokemon
}
