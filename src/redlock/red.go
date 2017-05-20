package redlock

import (
	"errors"
	"redis"
	"time"
)

const (
	lockScript = `
		return redis.call('SET', KEYS[1], ARGV[1], 'NX', 'PX', ARGV[2])
	`
	unlockScript = `
		if redis.call("get",KEYS[1]) == ARGV[1] then
		    return redis.call("del",KEYS[1])
		else
		    return 0
		end
	`
)

// Lock attempts to put a lock on the key for a specified duration (in milliseconds).
// If the lock was successfully acquired, true will be returned.
func Lock(key, value string, timeout time.Duration) (bool, error) {
	cmd := aredis.Client.SetNX(key, value, timeout)
	if cmd.Err() != nil {
		return false, nil
	}
	return true, nil
}

// Unlock attempts to remove the lock on a key so long as the value matches.
// If the lock cannot be removed, either because the key has already expired or
// because the value was incorrect, an error will be returned.
func Unlock(key []string, value string) error {
	cmd := aredis.Client.Eval(unlockScript, key)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	res, err := cmd.Result()
	if err != nil {
		return err
	}
	if res != value {
		return errors.New("Unlock failed, key or secret incorrect")
	}
	// Success
	return nil
}
