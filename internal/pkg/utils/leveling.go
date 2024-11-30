package utils

import (
	"errors"
	"sync"
)

type UserData struct {
	XP    int
	Level int
}

var (
	userData      = make(map[string]map[string]*UserData)
	userDataMutex = &sync.RWMutex{}
)

func GetUserLevelAndXP(guildID, userID string) (int, int, error) {
	userDataMutex.RLock()
	defer userDataMutex.RUnlock()

	if guild, ok := userData[guildID]; ok {
		if user, ok := guild[userID]; ok {
			return user.Level, user.XP, nil
		}
	}
	return 0, 0, errors.New("user data not found")
}

func AddXP(guildID, userID string, xp int) error {
	userDataMutex.Lock()
	defer userDataMutex.Unlock()

	if _, ok := userData[guildID]; !ok {
		userData[guildID] = make(map[string]*UserData)
	}

	if _, ok := userData[guildID][userID]; !ok {
		userData[guildID][userID] = &UserData{XP: 0, Level: 0}
	}

	userData[guildID][userID].XP += xp
	return nil
}

func CheckLevelUp(guildID, userID string) (bool, int, error) {
	userDataMutex.Lock()
	defer userDataMutex.Unlock()

	if guild, ok := userData[guildID]; ok {
		if user, ok := guild[userID]; ok {
			requiredXP := (user.Level + 1) * 100
			if user.XP >= requiredXP {
				user.Level++
				return true, user.Level, nil
			}
			return false, user.Level, nil
		}
	}
	return false, 0, errors.New("user data not found")
}
