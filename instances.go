package main

import "fmt"

type UID string
type Instance string

func buildUid(user UserJwt) UID {
	return UID(user.Username + "@" + user.Provider)
}

func nextFreeInstance(allInstances *map[Instance]bool) Instance {
	for instance, value := range *allInstances {
		if value {
			return instance
		}
	}
	return ""
}

func (uid UID) getInstance(userInstances *map[UID]Instance, allInstances *map[Instance]bool) (Instance, error) {
	instance, ok := (*userInstances)[uid]
	if ok {
		return instance, nil
	}

	instance = nextFreeInstance(allInstances)
	if instance == "" {
		return "", fmt.Errorf("no instances available")
	}
	(*allInstances)[instance] = false
	(*userInstances)[uid] = instance
	return instance, nil
}

func buildInstanceAvailability(allInstances *map[Instance]bool, backendMap BackendMap) {
	for _, backend := range backendMap.Backends {
		instance := Instance(backend.ID)
		(*allInstances)[instance] = true
	}
}
