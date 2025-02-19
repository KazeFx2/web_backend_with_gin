package HomeworkUpload

import (
	"fmt"
	"main/Msg"
	"strconv"
	"time"
)

var AliveTimers = make([]*time.Timer, 0)

var AlarmHours = [3]time.Duration{1, 3, 6}

func GenExpireTime(end time.Time, hour time.Duration) (time.Duration, bool) {
	eTime := end.Add(time.Hour * (-1 * hour))
	if eTime.Before(time.Now()) {
		return 0, false
	}
	return eTime.Sub(time.Now()), true
}

func GenTimer() {
	for _, i := range AliveTimers {
		i.Stop()
	}
	AliveTimers = make([]*time.Timer, 0)
	for _, i := range HwControl.Control.ExpControls {
		var exp = i.Exp
		for _, hour := range AlarmHours {
			dur, ava := GenExpireTime(i.End, hour)
			if !ava {
				break
			}
			var h = hour
			timer := time.AfterFunc(dur, func() {
				msg := Msg.Message{
					Stu: Msg.Student{
						StudentName: "",
						StudentNum:  strconv.Itoa(exp),
					},
					Msg: fmt.Sprintf("距离实验%d提交结束还有%d小时，请注意提交报告", exp, h),
				}
				Msg.MessageChan <- msg
			})
			AliveTimers = append(AliveTimers, timer)
		}
	}
}
