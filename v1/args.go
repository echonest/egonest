package egonest

// This file contains useful constants and utility functions.

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// License options for artist/biographies, artist/images, etc.
const (
	LicenseEchoSource        = "echo-source"
	LicenseAllRightsReserved = "all-rights-reserved"
	LicenseCCBYSA            = "cc-by-sa"
	LicenseCCBYNC            = "cc-by-nc"
	LicenseCCBYNCND          = "cc-by-nc-nd"
	LicenseCCBYNCDA          = "cc-by-nc-da"
	LicenseCCBYND            = "cc-by-nd"
	LicenseCCBY              = "cc-by"
	LicensePublicDomain      = "public-domain"
	LicenseUnknown           = "unknown"
)

// Hotttnesss types
const (
	HotttnesssOverall    = "overall"
	HotttnesssSocial     = "social"
	HotttnesssMainstream = "mainstream"
	HotttnesssReviews    = "reviews"
)

// Terms types
const (
	TermsStyle = "style"
	TermsMood  = "mood"
)

// The set of all buckets. Not all of these are valid for each call that accepts a bucket parameter.
// Consult developer.echonest.com for specifics.
const (
	BucketBios              = "biographies"
	BucketBlogs             = "blogs"
	BucketDocCounts         = "doc_counts"
	BucketFamiliarity       = "familiarity"
	BucketArtistFamiliarity = "artist_familiarity"
	BucketHotttnesss        = "hotttnesss"
	BucketSongHotttnesss    = "song_hotttnesss"
	BucketArtistHotttnesss  = "artist_hotttnesss"
	BucketSongType          = "song_type"
	BucketTracks            = "tracks"
	BucketImages            = "images"
	BucketArtistLocation    = "artist_location"
	BucketNews              = "news"
	BucketReviews           = "reviews"
	BucketSongs             = "songs"
	BucketTerms             = "terms"
	BucketURLs              = "urls"
	BucketVideo             = "video"
	BucketYearsActive       = "years_active"
	BucketAudioSummary      = "audio_summary"
	BucketGenre             = "genre"
)

type RosettaInfo struct {
	EntityTypes, Regions []string
}

// A Rosetta ID space maps to a slice of available entity types and a slice of available regions.
// If the slice of available regions is not empty, then you must append an available region to the catalog
// identifier when using that ID space as a Rosetta ID, e.g. rdio-NL:track:t10231
var Rosetta = map[string]RosettaInfo{
	"deezer":      {[]string{"artist", "track", "release"}, []string{}},
	"discogs":     {[]string{"artist"}, []string{}},
	"fma":         {[]string{"artist", "track", "release"}, []string{}},
	"jambase":     {[]string{"artist"}, []string{}},
	"lyricfind":   {[]string{"song"}, []string{}},
	"musicbrainz": {[]string{"artist"}, []string{}},
	"musixmatch":  {[]string{"song"}, []string{"WW"}},
	"rhapsody":    {[]string{"artist", "track"}, []string{"US"}},
	"rdio": {[]string{"artist", "song", "track"}, []string{"AT", "AU", "BR", "CA",
		"CH", "DE", "DK", "ES", "FI", "FR", "IE", "IT", "NL", "NO", "NZ", "PT", "SE", "UK", "US"}},
	"seatgeek":     {[]string{"artist"}, []string{}},
	"seatwave":     {[]string{"artist"}, []string{}},
	"songkick":     {[]string{"artist"}, []string{}},
	"songmeanings": {[]string{"artist", "song"}, []string{}},
	"spotify":      {[]string{"artist", "song", "track"}, []string{"US", "WW"}},
	"whosampled":   {[]string{"artist"}, []string{}},
	"7digital":     {[]string{"artist", "track"}, []string{"US", "UK"}},
}

// Location specifiers for calls like artist/search.
const (
	LocationCity    = "city"
	LocationCountry = "country"
	LocationRegion  = "region"
)

// Rankings
const (
	RankingFamiliarity = "familiarity"
	RankingRelevance   = "relevance"
)

// SearchTermBoost will format a search term argument in the manner needed by search calls.
func SearchTermBoost(term string, boostfactor float64) string {
	return fmt.Sprintf("%s^%f", term, boostfactor)
}

// SearchTermBan returns a search term formatted to be banned from a search.
func SearchTermBan(term string) string {
	return "-" + term
}

// SearchTermRequire returns a search term formatted to be required in a search.
func SearchTermRequire(term string) string {
	return "^" + term
}

// Valid arguments for sort parameters.
const (
	SortFamiliarity       = "familiarity"
	SortHotttnesss        = "hotttnesss"
	SortArtistStartYear   = "artist_start_year"
	SortArtistEndYear     = "artist_end_year"
	SortTempo             = "tempo"
	SortDuration          = "duration"
	SortLoudness          = "loudness"
	SortSpeechiness       = "speechiness"
	SortAcousticness      = "acousticness"
	SortLiveness          = "liveness"
	SortArtistFamiliarity = "artist_familiarity"
	SortArtistHotttnesss  = "artist_hotttnesss"
	SortSongHotttnesss    = "song_hotttnesss"
	SortLatitude          = "latitude"
	SortLongitude         = "longitude"
	SortMode              = "mode"
	SortKey               = "key"
	SortEnergy            = "energy"
	SortDanceability      = "danceability"
	SortWeight            = "weight"
	SortFrequency         = "frequency"
)

// Valid directions for sort parameters.
const (
	SortOrderAsc  = true
	SortOrderDesc = false
)

// SortOrder correctly formats a sort order argument to the API in the specified direction.
func SortOrder(attribute string, direction bool) string {
	if direction {
		return attribute + "-asc"
	} else {
		return attribute + "-desc"
	}
}

// Valid values for the "mode" audio attribute.
const (
	ModeMinor = 0
	ModeMajor = 1
)

// Valid values for the "key" audio attribute.
const (
	KeyC = iota
	KeyCSharp
	KeyD
	KeyEFlat
	KeyE
	KeyF
	KeyFSharp
	KeyG
	KeyAFlat
	KeyA
	KeyBFlat
	KeyB
)

// Known song types.
const (
	SongTypeChristmas = "christmas"
	SongTypeLive      = "live"
	SongTypeStudio    = "studio"
	SongTypeAcoustic  = "acoustic"
	SongTypeElectric  = "electric"
)

// Constants for use in song search.
const (
	SongTypeStateTrue = iota
	SongTypeStateFalse
	SongTypeStateAny
)

// Decorates a song type for use in song/search.
func SongTypeState(songtype string, state int) string {
	switch state {
	default:
		return songtype
	case SongTypeStateFalse:
		return songtype + ":false"
	case SongTypeStateAny:
		return songtype + ":any"
	}

}

// Known error codes
const (
	UnknownError  = -1
	InvalidKey    = 1
	KeyNotAllowed = 2
	RateLimit     = 3
	MissingArgs   = 4
	BadArgs       = 5
)
