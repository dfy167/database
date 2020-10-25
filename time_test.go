package database

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestString2Time(t *testing.T) {
	type args struct {
		timeString string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test1", args{"202001"}, time.Now(), false},
	}
	startTime, _ := String2Time("20200917")
	endTime, _ := String2Time("20200919")

	for _, t := range TimeBetween(startTime, endTime, Hour) {
		fmt.Println(Time2String(t, Hour))
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := String2Time(tt.args.timeString)
			if (err != nil) != tt.wantErr {
				t.Errorf("String2Time() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String2Time() = %v, want %v", got, tt.want)
			}
		})
	}
}
