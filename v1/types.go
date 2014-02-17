package egonest

// This file contains useful types for unmarshalling responses

type Genre struct {
	Name string `json:"name"`
}

type Term struct {
	Name      string  `json:"name"`
	Frequency float64 `json:"frequency"`
	Weight    float64 `json:"weight"`
}

type License struct {
	Type            string `json:"type"`
	Attribution     string `json:"attribution"`
	URL             string `json:"url"`
	Version         string `json:"version"`
	Attribution_URL string `json:"attribution-url"`
}

type Bio struct {
	Text      string `json:"text"`
	Site      string `json:"site"`
	URL       string `json:"url"`
	License   `json:"license"`
	Truncated bool `json:"truncated"`
}

type Blog struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	Summary     string `json:"summary"`
	Id          string `json:"id"`
	Date_posted string `json:"date_posted"`
	Date_found  string `json:"date_found"`
}

type Review struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	Summary       string `json:"summary"`
	Image_URL     string `json:"image_url"`
	Release       string `json:"release"`
	Id            string `json:"id"`
	Date_reviewed string `json:"date_reviewed"`
	Date_found    string `json:"date_found"`
}

type Years_Active struct {
	Start *int `json:"start"`
	End   *int `json:"end"`
}

type Video struct {
	Title      string `json:"title"`
	URL        string `json:"url"`
	Site       string `json:"site"`
	Image_URL  string `json:"image_url"`
	Id         string `json:"id"`
	Date_found string `json:"date_found"`
}

type Image struct {
	URL     string `json:"url"`
	License `json:"license"`
}

type News struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Id         string `json:"id"`
	Summary    string `json:"summary"`
	Date_found string `json:"date_found"`
}

type Location struct {
	City     string `json:"city"`
	Region   string `json:"region"`
	Location string `json:"location"`
	Country  string `json:"country"`
}

// type Song struct {
// 	Id   string `json:"id"`
// 	Name string `json:"name"`
// }

// type Term struct {
// 	Name string `json:"name"`
// }

type Foreign_ID struct {
	Catalog    string `json:"catalog"`
	Foreign_id string `json:"foreign_id"`
}

type Track struct {
	Foreign_ID
	Foreign_release_id string `json:"foreign_release_id"`
	Id                 string `json:"id"`
}

type TrackAnalysisResponse struct {
}

type Artist struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	Genres          []Genre  `json:"genres"`
	Terms           []Term   `json:"terms"`
	Biographies     []Bio    `json:"biographies"`
	Blogs           []Blog   `json:"blogs"`
	Familiarity     float64  `json:"familiarity"`
	Hotttnesss      float64  `json:"hotttnesss"`
	Reviews         []Review `json:"reviews"`
	Years_Active    []Years_Active
	Video           []Video           `json:"video"`
	Urls            map[string]string `json:"urls"`
	Images          []Image           `json:"images"`
	News            []News            `json:"news"`
	Doc_counts      map[string]int    `json:"doc_counts"`
	Artist_location Location          `json:"artist_location"`
	Songs           []Song            `json:"songs"`
	Foreign_ids     []Foreign_ID      `json:"foreign_ids"`
	Twitter         string            `json:"twitter"`
}

type Song struct {
	Id                 string   `json:"id"`
	Title              string   `json:"title"`
	Artist_id          string   `json:"artist_id"`
	Artist_hotttnesss  float64  `json:"artist_hotttnesss"`
	Artist_name        string   `json:"artist_name"`
	Song_type          []string `json:"string_type"`
	Tracks             []Track  `json:"tracks"`
	Artist_location    Location `json:"artist_location"`
	Song_hotttnesss    float64  `json:"song_hotttnesss"`
	Artist_familiarity float64  `json:"artist_familiarity"`
	Audio_summary      `json:"audio_summary"`
}

type Audio_summary struct {
	Key            int     `json:"key"`
	Analysis_url   string  `json:"analysis_url"`
	Energy         float64 `json:"energy"`
	Liveness       float64 `json:"liveness"`
	Tempo          float64 `json:"tempo"`
	Speechiness    float64 `json:"speechiness"`
	Acousticness   float64 `json:"acousticness"`
	Mode           int     `json:"mode"`
	Time_signature int     `json:"time_signature"`
	Duration       float64 `json:"duration"`
	Loudness       float64 `json:"loudness"`
	Audio_md5      string  `json:"audio_md5"`
	Valence        float64 `json:"valence"`
	Danceability   float64 `json:"danceability"`
}

type TimeRange struct {
	Start      float64 `json:"start"`
	Duration   float64 `json:"duration"`
	Confidence float64 `json:"confidence"`
}

type Section struct {
	TimeRange
	Loudness                  float64 `json:"loudness"`
	Tempo                     float64 `json:"tempo"`
	Tempo_confidence          float64 `json:"tempo_confidence"`
	Key                       int     `json:"key"`
	Key_confidence            float64 `json:"key_confidence"`
	Mode                      int     `json:"mode"`
	Mode_confidence           float64 `json:"mode_confidence"`
	Time_signature            int     `json:"time_signature"`
	Time_signature_confidence float64 `json:"time_signature_confidence"`
}

type Segment struct {
	TimeRange
	Loudness_start    float64   `json:"loudness_start"`
	Loudness_max      float64   `json:"loudness_max"`
	Loudness_max_time float64   `json:"loudness_max_time"`
	Pitches           []float64 `json:"pitches"`
	Timbre            []float64 `json:"timbre"`
}

type Analysis struct {
	Meta struct {
		Analyzer_version string  `json:"analyzer_version"`
		Platform         string  `json:"platform"`
		Detailed_status  string  `json:"detailed_status"`
		Filename         string  `json:"filename"`
		Artist           string  `json:"artist"`
		Album            string  `json:"album"`
		Title            string  `json:"title"`
		Genre            string  `json:"genre"`
		Bitrate          int     `json:"bitrate"`
		Sample_rate      int     `json:"sample_rate"`
		Seconds          int     `json:"seconds"`
		Status_code      int     `json:"status_code"`
		Timestamp        int     `json:"timestamp"`
		Analysis_time    float64 `json:"analysis_time"`
	} `json:"meta"`
	Track struct {
		Num_samples               int     `json:"num_samples"`
		Duration                  float64 `json:"duration"`
		Sample_md5                string  `json:"sample_md5"`
		Decoder                   string  `json:"decoder"`
		Decoder_version           string  `json:"decoder_version"`
		Offset_seconds            float64 `json:"offset_seconds"`
		Window_seconds            float64 `json:"window_seconds"`
		Analysis_sample_rate      int     `json:"analysis_sample_rate"`
		Analysis_channels         int     `json:"analysis_channels"`
		End_of_fade_in            float64 `json:"end_of_fade_in"`
		Start_of_fade_out         float64 `json:"start_of_fade_out"`
		Codestring                string  `json:"codestring"`
		Code_version              string  `json:"code_version"`
		Echoprintstring           string  `json:"echoprintstring"`
		Echoprint_version         string  `json:"echoprint_version"`
		Synchstring               string  `json:"synchstring"`
		Synch_version             string  `json:"synch_version"`
		Rhythmstring              string  `json:"rhythmstring"`
		Rhythm_version            string  `json:"rhythm_version"`
		Loudness                  float64 `json:"loudness"`
		Tempo                     float64 `json:"tempo"`
		Tempo_confidence          float64 `json:"tempo_confidence"`
		Key                       int     `json:"key"`
		Key_confidence            float64 `json:"key_confidence"`
		Mode                      int     `json:"mode"`
		Mode_confidence           float64 `json:"mode_confidence"`
		Time_signature            int     `json:"time_signature"`
		Time_signature_confidence float64 `json:"time_signature_confidence"`
	} `json:"track"`
	Bars     []TimeRange `json:"bars"`
	Beats    []TimeRange `json:"beats"`
	Tatums   []TimeRange `json:"tatums"`
	Segments []Segment   `json:"segments"`
	Sections []Section   `json:"sections"`
}
