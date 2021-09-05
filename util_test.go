package urlutil_test

import (
	"reflect"
	"testing"

	urlutil "github.com/mkfsn/googlemap-urlutil"
)

func TestParseEmbed(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    urlutil.Coordinate
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				url: "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3681.5759829131234!2d120.30009731496122!3d22.66959258513296!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x0%3A0x1a7a18fa2156ded8!2zR1Ug6auY6ZuE5ryi56We5beo6JuL6LO854mp5buj5aC05bqX!5e0!3m2!1szh-TW!2stw!4v1559715294783!5m2!1szh-TW!2stw",
			},
			want: urlutil.Coordinate{
				Lat: 22.66959258513296,
				Lng: 120.30009731496122,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := urlutil.ParseEmbed(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEmbed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseEmbed() got = %v, want %v", got, tt.want)
			}
		})
	}
}
