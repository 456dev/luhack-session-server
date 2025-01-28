package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type UID string
type Instance string

func buildUid(user UserJwt) UID {
	return UID(user.Username + "@" + user.Provider)
}

func nextFreeInstance(allInstances *map[Instance]bool, backends []Backend) Instance {
	for instance, value := range *allInstances {
		isHealthy := instance.isHealthy(backends)
		if !isHealthy {
			log.Println("Instance ", instance, " is unhealthy")
		}
		if value && isHealthy {
			return instance
		}
	}
	return ""
}

func (uid UID) getInstance(userInstances *map[UID]Instance, allInstances *map[Instance]bool, backends []Backend) (Instance, error) {
	instance, ok := (*userInstances)[uid]
	if ok {
		return instance, nil
	}

	instance = nextFreeInstance(allInstances, backends)
	if instance == "" {
		return "", fmt.Errorf("no instances available")
	}
	(*allInstances)[instance] = false
	(*userInstances)[uid] = instance
	storeInstances(*userInstances, *allInstances)
	log.Println("Assigned instance ", instance, " to ", uid)
	return instance, nil
}

func (uid UID) releaseInstance(userInstances *map[UID]Instance, allInstances *map[Instance]bool) error {
	instance, ok := (*userInstances)[uid]
	if !ok {
		return fmt.Errorf("no instance to release")
	}
	delete(*userInstances, uid)
	//NOTE we don't make the instance available again as it is potentially a security risk
	storeInstances(*userInstances, *allInstances)
	log.Println("Released instance ", instance, " from ", uid)
	return nil
}

func pokeHTTP(host string) bool {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(host)
	if err != nil {
		return false
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return true
}

func (instanceID *Instance) isHealthy(instances []Backend) bool {
	found := false
	var instance Backend
	for _, instance = range instances {
		if instance.ID == string(*instanceID) {
			found = true
			break
		}
	}
	if !found {
		return false
	}
	for _, service := range instance.Services {
		if service.Host == "" {
			return false
		}
		if !pokeHTTP("http://" + service.Host) {
			return false
		}
	}
	return true
}

func storeInstances(userInstances map[UID]Instance, allInstances map[Instance]bool) {
	file, err := os.Create("instances.gob")
	if err != nil {
		log.Fatal("Failed to create file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Failed to close file:", err)
		}
	}(file)

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(userInstances)
	if err != nil {
		log.Fatal("Failed to encode userInstances:", err)
		return
	}

	err = encoder.Encode(allInstances)
	if err != nil {
		log.Fatal("Failed to encode allInstances:", err)
		return
	}
}

func buildInstanceAvailability(allInstances *map[Instance]bool, backendMap BackendMap) {
	for _, backend := range backendMap.Backends {
		instance := Instance(backend.ID)
		(*allInstances)[instance] = true
	}
}

func loadInstances(backendMap BackendMap) (map[UID]Instance, map[Instance]bool) {
	backendMapInstances := make(map[Instance]bool)
	buildInstanceAvailability(&backendMapInstances, backendMap)

	file, err := os.Open("instances.gob")
	if err != nil {
		log.Println("Failed to open file:", err)
		return make(map[UID]Instance), backendMapInstances
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Failed to close file:", err)
		}
	}(file)

	decoder := gob.NewDecoder(file)

	var userInstances map[UID]Instance
	err = decoder.Decode(&userInstances)
	if err != nil {
		log.Println("Failed to decode userInstances:", err)
		userInstances = make(map[UID]Instance)
	}

	var allInstances map[Instance]bool
	err = decoder.Decode(&allInstances)
	if err != nil {
		log.Println("Failed to decode allInstances:", err)
		allInstances = make(map[Instance]bool)
	}

	for instance, _ := range allInstances {
		if _, ok := backendMapInstances[instance]; !ok {
			log.Println("Removing invalid instance from allInstances: ", instance)
			delete(allInstances, instance)
		}
	}

	for _, backend := range backendMap.Backends {
		instance := Instance(backend.ID)
		if _, ok := allInstances[instance]; !ok {
			log.Println("Adding missing instance to allInstances: ", instance)
			allInstances[instance] = true
		}
	}

	for uid, instance := range userInstances {
		if _, ok := allInstances[instance]; !ok {
			log.Println("Removing invalid instance from userInstances: ", instance)
			delete(userInstances, uid)
		} else {
			log.Println("Setting existing instance to used: ", instance, "(belongs to ", uid, ")")
			allInstances[instance] = false
		}

	}

	return userInstances, allInstances
}
