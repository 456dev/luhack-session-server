package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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

func buildInstanceAvailability(allInstances *map[Instance]bool, backendMap BackendMap) {
	for _, backend := range backendMap.Backends {
		instance := Instance(backend.ID)
		(*allInstances)[instance] = true
	}
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
	// TODO implement instance storage
}

func loadInstances(backendMap BackendMap) (map[UID]Instance, map[Instance]bool) {
	// TODO implement instance loading
	// TODO buildInstanceAvailability(allInstances, backendMap)
	allInstances := make(map[Instance]bool)
	buildInstanceAvailability(&allInstances, backendMap)
	return make(map[UID]Instance), allInstances
}
