package logrus_proxy

import "github.com/Sirupsen/logrus"

type Hook struct {
	levels   map[logrus.Level]bool
	hook     logrus.Hook
	fireFunc func(entry *logrus.Entry)
}

func (hook *Hook) Levels() []logrus.Level {
	res := make([]logrus.Level, len(hook.levels))
	for lvl := range hook.levels {
		res = append(res, lvl)
	}
	return res
}

func (hook *Hook) SetPreFireFunc(fn func(entry *logrus.Entry)) *Hook {
	hook.fireFunc = fn
	return hook
}

func (hook *Hook) Fire(entry *logrus.Entry) (err error) {
	if enabled, exists := hook.levels[entry.Level]; exists && enabled {
		if hook.fireFunc != nil {
			hook.fireFunc(entry)
		}
		err = hook.hook.Fire(entry)
	}
	return
}

func (hook *Hook) DisableLevel(level logrus.Level) *Hook {
	if enabled, exists := hook.levels[level]; exists && enabled {
		hook.levels[level] = false
	}
	return hook
}

func (hook *Hook) EnableLevel(level logrus.Level) *Hook {
	if enabled, exists := hook.levels[level]; exists && enabled {
		return hook
	}
	hook.levels[level] = supportsLevel(hook, level)
	return hook
}

func NewHook(hook logrus.Hook, levels []logrus.Level) *Hook {
	lvls := make(map[logrus.Level]bool)
	for _, lvl := range levels {
		lvls[lvl] = supportsLevel(hook, lvl)
	}
	return &Hook{
		hook:   hook,
		levels: lvls,
	}
}

func supportsLevel(hook logrus.Hook, level logrus.Level) bool {
	for _, lvl := range hook.Levels() {
		if lvl == level {
			return true
		}
	}
	return false
}
