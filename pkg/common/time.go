package common

import "time"

// Time 生成当前时间戳
func Time() int64 {
	return time.Now().Unix()
}

// Strtotime strtotime()
// Strtotime("02/01/2006 15:04:05", "02/01/2016 15:04:05") == 1451747045
// Strtotime("3 04 PM", "8 41 PM") == -62167144740
func Strtotime(format, strtime string) (int64, error) {
	loc, err := time.LoadLocation("Local") //获取时区
	if err != nil {
		return 0, err
	}
	//t, err := time.Parse(format, strtime)
	t, err := time.ParseInLocation(format, strtime, loc)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// Date Date("02/01/2006 15:04:05 PM", 1524799394)
func Date(format string, timestamp int64) string {
	return time.Unix(timestamp, 0).Format(format)
}

// Sleep 休眠N秒
func Sleep(t int64) {
	time.Sleep(time.Duration(t) * time.Second)
}
